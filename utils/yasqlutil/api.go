package yasqlutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"yhc/commons/constants"
	"yhc/log"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	// only used to check password
	_check_password_sql = "SELECT STATUS FROM V$INSTANCE"
)

func GetLocalInstance(user string, password string, yasqlHome string, yasdbData string) *Yasql {
	var err error
	if !fs.IsFileExist(yasqlHome) {
		err = ErrYasqlHomeNotExist
	}
	return &Yasql{
		Mutex:        &sync.Mutex{},
		User:         user,
		Password:     Quote(password),
		YasdbData:    yasdbData,
		YasqlHome:    yasqlHome,
		ConnectLocal: true,
		SqlStatement: &SqlStatement{},
		err:          err,
	}
}

func GetInstance(user string, password string, ip string, port uint, yasqlHome string) *Yasql {
	var err error
	if !fs.IsFileExist(yasqlHome) {
		err = ErrYasqlHomeNotExist
	}
	return &Yasql{
		Mutex:        &sync.Mutex{},
		User:         user,
		Password:     Quote(password),
		Ip:           ip,
		Port:         port,
		YasqlHome:    yasqlHome,
		SqlStatement: &SqlStatement{},
		err:          err,
	}
}

func (tx *Yasql) Error() error {
	return tx.err
}

func (tx *Yasql) ResetError() {
	tx.err = nil
}

// Limits:
// 1.s.Table & s.Columns must be set.
// 2.if table has a column, of which type is not string type -> use s.ColTypes
// 3.FIXME: only can select one too long col, and need put it at last position of select
func (tx *Yasql) Select(s *Select) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	if s != nil {
		if s.Table == "" {
			tx.err = ErrTableNotSet
			return tx
		}
		if len(s.Columns) == 0 {
			tx.err = ErrNoColumns
			return tx
		}
		var columns string
		for _, col := range s.Columns {
			columns += str(col, COMMA)
		}
		columns = columns[0 : len(columns)-len(COMMA)]
		tx.SqlStatement.Select = str(SELECT, columns, FROM, s.Table)
		tx.SqlStatement.ColTypes = s.ColTypes
		return tx
	}
	tx.err = ErrInvalidSelect
	return tx
}

// Limits:
// 1.sql & colTypes required.
// 2.if table has a column, of which type is not string type -> use s.ColTypes
// 3.FIXME: only can select one too long col, and need put it at last position of select
func (tx *Yasql) SelectRaw(s *SelectRaw) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	if s == nil || s.RawSql == "" {
		tx.err = ErrInvalidSelect
		return tx
	}
	tx.SqlStatement.Select = strings.ReplaceAll(s.RawSql, ";", "")
	tx.SqlStatement.ColTypes = s.ColTypes
	return tx
}

// Limits:
// 1.should be used with tx.Select() | tx.SelectRaw()
// 2.type of args should be string
//   - you can use like this: tx.Where("id = ? and user = ?","001","admin")
//   - if not -> use like this: tx.Where(fmt.Sprintf("id = %d",1))
func (tx *Yasql) Where(whereCondition string, args ...interface{}) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	whereStr := strings.Replace(whereCondition, "?", "'%s'", -1)
	//在使用上已经进行了约束，只有包含?的where条件才会有args可变参数，需要保证args中的字符串中不包含单个数量的连续单引号
	newArgs := tx.escapeQuotes(args...)
	whereStr = fmt.Sprintf(whereStr, newArgs...)
	tx.SqlStatement.Where = append(tx.SqlStatement.Where, whereStr)
	return tx
}

// Usage: tx.Raw().Execute()
// Limits:
// 1.type of args should be string
//   - you can use like this: tx.Raw("insert into table(id,name) values(?,?)","001","admin")
//   - if not -> use like this:  tx.Raw(fmt.Sprintf("insert into table(id,name) values(%d,?)",1),"admin")
func (tx *Yasql) Raw(optCondition string, args ...interface{}) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	tx.resetSqlStatement()
	optCondition = strings.Replace(optCondition, "?", "'%s'", -1)
	//在使用上已经进行了约束，只有包含?的where条件才会有args可变参数，需要保证args中的字符串中不包含单个数量的连续单引号
	newArgs := tx.escapeQuotes(args...)
	optCondition = fmt.Sprintf(optCondition, newArgs...)
	tx.SqlStatement.Sql.WriteString(optCondition)
	return tx
}

// Limits:
// 1.should be used with tx.Select() | tx.SelectRaw()
// 2.sort.ColName must be set. if sort.Order not set: use db default order -> ASC
func (tx *Yasql) SortBy(sorts ...*Sort) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	var sortStr string
	for _, sort := range sorts {
		sortStr += str(sort.ColName, BLANK, string(sort.Order), COMMA)
	}
	sortStr = sortStr[0 : len(sortStr)-len(COMMA)]
	tx.SqlStatement.Sort = append(tx.SqlStatement.Sort, sortStr)
	return tx
}

// Limits: should be used with tx.Select() | tx.SelectRaw()
func (tx *Yasql) Paging(limit int, offset int) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	tx.SqlStatement.Paging = fmt.Sprintf("%s%d%s%d", LIMIT, limit, OFFSET, offset)
	return tx
}

func (tx *Yasql) Count(s *Select, count *int64) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	if s == nil {
		tx.err = ErrInvalidSelect
		log.Yasql.Error(tx.err)
		return tx
	}
	if s.Table == "" {
		tx.err = ErrTableNotSet
		log.Yasql.Error(tx.err)
		return tx
	}
	// stmt will be reset after sql exec, thus needing a new instance when using tx.Count() with others
	query := tx.copy()
	query.SqlStatement.Select = str(SELECT, "count(*) as TOTAL", FROM, s.Table)
	query.SqlStatement.IsCountSql = true
	dataList := make([]map[string]interface{}, 0)
	if query.Find(&dataList).Error() != nil {
		return tx
	}
	str := dataList[0]["TOTAL"].(string)
	*count, _ = strconv.ParseInt(str, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	return tx
}

func (tx *Yasql) CountRaw(sql string, count *int64) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	if sql == "" {
		tx.err = ErrInvalidSelect
		return tx
	}
	// stmt will be reset after sql exec, thus needing a new instance when using tx.Count() with others
	query := tx.copy()
	query.SqlStatement.Select = sql
	dataList := make([]map[string]interface{}, 0)
	if query.Find(&dataList).Error() != nil {
		tx.err = query.err
		return tx
	}
	str := dataList[0]["TOTAL"].(string)
	*count, _ = strconv.ParseInt(str, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	return tx
}

// Limits: should be used with tx.Select() | tx.Count() | tx.SelectRaw() | tx.CountRaw()
func (tx *Yasql) Find(dest interface{}) *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	tx.SqlStatement.buildSql()
	colTypes := tx.SqlStatement.ColTypes
	ret, res, stderr := tx.exec()
	log.Yasql.Debug(res)
	if ret != 0 {
		tx.err = errors.New(stderr)
		log.Yasql.Error(FormatError(tx.err))
		return tx
	}
	if stderr != "" {
		tx.err = errors.New(stderr)
		log.Yasql.Error(FormatError(tx.err))
		return tx
	}
	rows, err := getRows(res, tx.IsGetTableNotExistErr)
	if err != nil {
		tx.err = err
		log.Yasql.Error(FormatError(tx.err))
		return tx
	}
	dataList := make([]map[string]interface{}, 0)
	if len(rows) >= 3 {
		// first row is title, the second is split line
		titles := getTitles(rows)
		colLens := getColLens(rows[1])
		for i := 2; i < len(rows)-1; i++ {
			elements := getElements(rows[i], colLens)
			data := map[string]interface{}{}
			for i := 0; i < len(elements); i++ {
				data[titles[i]] = getRealValue(colTypes, titles[i], elements[i])
			}
			dataList = append(dataList, data)
		}
		err := setData(dataList, dest)
		if err != nil {
			tx.err = err
			log.Yasql.Error(FormatError(tx.err))
			return tx
		}
	}
	return tx
}

// Limits: should be used with tx.Raw() or tx.Select()
func (tx *Yasql) Execute() *Yasql {
	tx.Lock()
	defer tx.Unlock()
	if tx.err != nil {
		return tx
	}
	tx.SqlStatement.buildSql()
	ret, res, stderr := tx.exec()
	if ret != 0 {
		tx.err = errors.New(stderr)
		log.Yasql.Error(FormatError(tx.err))
		return tx
	}
	if stderr != "" {
		tx.err = errors.New(stderr)
		log.Yasql.Error(FormatError(tx.err))
		return tx
	}
	_, err := getRows(res, tx.IsGetTableNotExistErr)
	if err != nil {
		tx.err = err
		log.Yasql.Error(FormatError(tx.err))
	}
	return tx
}

func (tx *Yasql) CheckPassword() error {
	err := tx.Raw(_check_password_sql).Execute().Error()
	tx.ResetError()
	return err
}

func (tx *Yasql) escapeQuotes(inputs ...interface{}) []interface{} {
	res := []interface{}{}
	for _, input := range inputs {
		value, ok := input.(string)
		if !ok {
			res = append(res, input)
		}
		escapedInput := strings.ReplaceAll(value, "'", "''")
		res = append(res, escapedInput)
	}
	return res
}
