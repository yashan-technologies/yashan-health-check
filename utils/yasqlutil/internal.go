package yasqlutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"yhc/commons/constants"
	"yhc/log"
	"yhc/utils/pwdutil"

	"github.com/shopspring/decimal"
)

const (
	// YASQL ENV
	YASQL_BIN  = "yasql"
	BIN_PATH   = "bin"
	LIB_PATH   = "lib"
	LIB_KEY    = "LD_LIBRARY_PATH"
	YASDB_DATA = "YASDB_DATA"

	// SQL EXECUTE
	YASQL_CONNECT_ERROR_CODE = "YASQL-00007"
	YAS_ERROR_PREFIX         = "YAS-"
	BLANK                    = " "
	COMMA                    = ","
	SELECT                   = "select" + BLANK
	FROM                     = BLANK + "from" + BLANK
	WHERE                    = BLANK + "where" + BLANK
	AND                      = BLANK + "and" + BLANK
	ORDER_BY                 = BLANK + "order by" + BLANK
	LIMIT                    = BLANK + "limit" + BLANK
	OFFSET                   = BLANK + "offset" + BLANK
	DEFAULT_ORDER_BY         = ORDER_BY + "1 desc"

	// yas err
	YAS_DB_NOT_OPEN                  = "YAS-02078"
	YAS_NO_DBUSER                    = "YAS-02010"
	YAS_INVALID_USER_OR_PASSWORD     = "YAS-02143"
	YAS_USER_LACK_LOGIN_AUTH         = "YAS-02245"
	YAS_USER_LACK_AUTH               = "YAS-02213"
	YAS_TABLE_OR_VIEW_DOES_NOT_EXIST = "YAS-02012"
	YAS_FAILED_CONNECT_SOCKET        = "YAS-00402"

	// yas err regexp
	_yasErrPattern   = ".*YAS-[0-9]{5}.*"
	_yasqlErrPattern = ".*YASQL-[0-9]{5}.*"
	_yasqlErrUser    = "please input user name"
	_yasqlErrValue   = "Enter value for"
)

var (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type Yasql struct {
	*sync.Mutex
	User                  string
	Password              string
	Ip                    string
	Port                  uint
	YasqlHome             string
	YasdbData             string
	ConnectLocal          bool // 是否使用本地环境变量连接数据库，无需ip端口
	SqlStatement          *SqlStatement
	err                   error
	HostId                string // carry host id
	YasdbId               string // carry yasdb id
	IsGetTableNotExistErr bool   // 如果为true 不将table view not exist err 转换为没有权限,如果为false  table view not exist err 转换为没有权限
}

type Select struct {
	Table    string
	Columns  []string
	ColTypes map[string]string // default type -> string
}

type SelectRaw struct {
	RawSql   string
	ColTypes map[string]string // default type -> string
}

type Order string

type Sort struct {
	ColName string
	Order   Order
}

func (tx *Yasql) copy() *Yasql {
	return &Yasql{
		Mutex:     &sync.Mutex{},
		User:      tx.User,
		Password:  tx.Password,
		Ip:        tx.Ip,
		Port:      tx.Port,
		YasqlHome: tx.YasqlHome,
		SqlStatement: &SqlStatement{
			Table:    tx.SqlStatement.Table,
			Select:   tx.SqlStatement.Select,
			Where:    tx.SqlStatement.Where,
			Sort:     tx.SqlStatement.Sort,
			Paging:   tx.SqlStatement.Paging,
			ColTypes: tx.SqlStatement.ColTypes,
		},
		IsGetTableNotExistErr: tx.IsGetTableNotExistErr,
		err:                   tx.err,
	}
}

func getRows(res string, isGetTableNotExistErr bool) ([]string, error) {
	res = strings.Replace(res, "\n\n", "\n", -1)
	res = strings.TrimSpace(res)
	rows := strings.Split(res, "\n")
	matchYasqlErr, _ := regexp.Match(_yasqlErrPattern, []byte(res))
	matchYasErr, _ := regexp.Match(_yasErrPattern, []byte(res))
	matchUserErr, _ := regexp.Match(_yasqlErrUser, []byte(res))
	matchValueErr, _ := regexp.Match(_yasqlErrValue, []byte(res))
	if matchYasqlErr || matchYasErr || matchUserErr || matchValueErr {
		log.Yasql.Errorf("Get Rows Match err: %s\n", res)
		return nil, NewYasErr(res)
	}
	return rows, nil
}

func getTitles(rows []string) []string {
	titleRow := rows[0]
	regPattern := "\\s+"
	reg, _ := regexp.Compile(regPattern)
	res := reg.ReplaceAllString(titleRow, COMMA)
	return strings.Split(res, COMMA)
}

func getElements(row string, lens []int) []string {
	elems := []string{}
	start := 0
	for i := 0; i < len(lens)-1; i++ {
		end := start + lens[i] + len(BLANK)
		if end > len(row)-1 { // 说明是最后一个
			continue
		}
		elem := row[start:end]
		elem = removeAllBlank(elem)
		elems = append(elems, elem)
		start = end
	}
	last := strings.TrimSpace(row[start:])
	elems = append(elems, last)
	return elems
}

func removeAllBlank(s string) string {
	// replace all blanks to ""
	return strings.TrimSpace(s)
}

func getColLens(secondRow string) []int {
	secondRow = strings.TrimSpace(secondRow)
	elems := strings.Split(secondRow, BLANK)
	lens := []int{}
	for _, elem := range elems {
		lens = append(lens, len(elem))
	}
	return lens
}

// after exec, sqlStatement will be reset
func (tx *Yasql) exec() (int, string, string) {
	sql := tx.SqlStatement.Sql.String()
	cmd := tx.genCmd()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	start := time.Now()
	args := strings.Join(cmd.Args, " ")
	// passwd protected
	args = pwdutil.PrintLogDesensitize(args)
	log.Yasql.Debug("exec: ", args, " starting")
	if err := cmd.Start(); err != nil {
		log.Yasql.Errorf("exec %s start failed, %s", cmd.Args[0], err)
		return -1, "", err.Error()
	}
	if err := cmd.Wait(); err != nil {
		log.Yasql.Errorf("exec %s wait failed, %s", cmd.Args[0], err)
		log.Yasql.Errorf("exec:%s, %s", stdout.String(), stderr.String())
	}
	log.Yasql.Info(FormatInfo(sql, start))
	code := cmd.ProcessState.ExitCode()
	return code, stdout.String(), stderr.String()
}

func (tx *Yasql) genCmd() *exec.Cmd {
	defer tx.resetSqlStatement()
	yasqlBin := path.Join(tx.YasqlHome, BIN_PATH, YASQL_BIN)
	env := []string{fmt.Sprintf("%s=%s", LIB_KEY, path.Join(tx.YasqlHome, LIB_PATH))}
	connectStr := fmt.Sprintf("%s/%s@%s:%d", tx.User, tx.Password, tx.Ip, tx.Port)
	if tx.ConnectLocal {
		env = append(env, fmt.Sprintf("%s=%s", YASDB_DATA, tx.YasdbData))
		connectStr = fmt.Sprintf("%s/%s", tx.User, tx.Password)
	}
	cmd := exec.Command(yasqlBin, connectStr, "-c", tx.SqlStatement.Sql.String())
	cmd.Env = env
	return cmd
}

// string contact
func str(args ...string) string {
	buf := &bytes.Buffer{}
	for _, s := range args {
		buf.WriteString(s)
	}
	return buf.String()
}

// parse string value
func getRealValue(colTypes map[string]string, key string, strValue string) interface{} {
	v, ok := colTypes[key]
	if !ok {
		return strValue
	}
	switch v {
	case reflect.Int.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String():
		if strings.Contains(strings.ToUpper(strValue), "E") {
			decimalNum, err := decimal.NewFromString(strValue)
			if err == nil {
				strValue = decimalNum.String()
			}
		}
		num, _ := strconv.ParseInt(strValue, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
		return num
	case reflect.Uint.String(), reflect.Uint8.String(), reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String():
		num, _ := strconv.ParseUint(strValue, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
		return num
	case reflect.Float32.String(), reflect.Float64.String():
		num, _ := strconv.ParseFloat(strValue, constants.BIT_SIZE_64)
		return num
	case reflect.Bool.String():
		val, _ := strconv.ParseBool(strValue)
		return val
	default:
		return strValue
	}
}

// use json to convert []map to []struct, thus column name needing be same with struct field name.
// case insensitive, example: "ID", "id", "iD" or "Id" is ok.
func setData(dataList interface{}, dest interface{}) error {
	ref := reflect.TypeOf(dest)
	if ref.Kind() != reflect.Ptr {
		return ErrInvalidDest
	}
	data, err := json.Marshal(dataList)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

type SqlStatement struct {
	Table      string
	Select     string
	Where      []string
	Sort       []string
	Paging     string
	ColTypes   map[string]string
	Sql        bytes.Buffer
	IsCountSql bool
}

// func buildSql() makes tx.Where().Sort().Select().Where().Sort() valid
func (s *SqlStatement) buildSql() {
	if s.Sql.String() != "" {
		return
	}
	s.Sql.WriteString(s.Select)
	if len(s.Where) > 0 {
		s.Sql.WriteString(WHERE)
		var whereCondition string
		for _, where := range s.Where {
			whereCondition += str(where, AND)
		}
		whereCondition = whereCondition[0 : len(whereCondition)-len(AND)]
		s.Sql.WriteString(whereCondition)
	}
	if !s.IsCountSql && len(s.Sort) > 0 {
		s.Sql.WriteString(ORDER_BY)
		var sortCondition string
		for _, sort := range s.Sort {
			sortCondition += str(sort, COMMA)
		}
		sortCondition = sortCondition[0 : len(sortCondition)-len(COMMA)]
		s.Sql.WriteString(sortCondition)
	} else {
		s.Sql.WriteString(DEFAULT_ORDER_BY)
	}
	s.Sql.WriteString(s.Paging)
}

func (tx *Yasql) resetSqlStatement() {
	tx.SqlStatement = new(SqlStatement)
	tx.SqlStatement.Sql.Reset()
}

func FormatInfo(sql string, start time.Time) string {
	desensitizeSQL := pwdutil.PrintLogDesensitize(sql)
	return fmt.Sprintf("[YASQL] [%s] %s", time.Since(start), desensitizeSQL)
}

func FormatError(err error) error {
	if err != nil {
		err = fmt.Errorf("[YASQL] FAILED: %s", err)
	}
	return err
}
