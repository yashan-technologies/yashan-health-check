package yasqlutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"yhc/commons/constants"
	"yhc/defs/timedef"
	"yhc/log"
)

type DataType string

type QueryRangeType int8

const (
	_timestamp DataType = "timestamp"
	_int       DataType = "int"
)

const (
	_range_invalid QueryRangeType = iota
	_range_min_only
	_range_max_only
	_range_both
)

type QueryRange struct {
	Datatype DataType `json:"datatype"`
	Max      string   `json:"max"`
	Min      string   `json:"min"`
}

// if data from only one table -> use PaginationQ.Select, select.Table and select.Columns required
// else if data from muti tables -> use PaginationQ.RawSelect and CountSql
type PaginationQ struct {
	Select                *Select                              `json:"-"          form:"-"`
	SelectRaw             *SelectRaw                           `json:"-"          form:"-"`
	CountSql              string                               `json:"-"          form:"-"`
	Ok                    bool                                 `json:"-"`
	Size                  int                                  `json:"pageSize"   form:"pageSize"` // 页面记录条数
	Page                  int                                  `json:"pageNo"     form:"pageNo"`   // 第几页
	Sort                  string                               `json:"-"          form:"sortBy"`
	Order                 string                               `json:"-"          form:"orderBy"`
	Filter                string                               `json:"-"          form:"filter"`
	Search                string                               `json:"-"          form:"search"`
	Range                 string                               `json:"-"          form:"range"`
	Data                  interface{}                          `json:"data"       swaggertype:"object"` // data list
	Total                 int64                                `json:"totalCount"`                      // 总记录数
	TotalPage             int                                  `json:"totalPage"`                       // 总页面数
	AfterSearch           func(tx *Yasql, v interface{}) error `json:"-"`
	excludeHumpReplaceMap map[string]struct{}                  `json:"-"`
}

type WhereFunc func(tx *Yasql) *Yasql

// replace hump word. aBc -> A_BC
func (p *PaginationQ) humpReplace(k string) string {
	if k == "" {
		return k
	}
	if p.excludeHumpReplace(k) {
		return k
	}
	s := []string{}
	for i := 0; i < len(k); i++ {
		if i == 0 {
			lowS := strings.ToLower(string(k[i]))
			s = append(s, lowS)
			continue
		}
		if 64 < k[i] && k[i] < 91 {
			lowS := strings.ToLower(string(k[i]))
			s = append(s, fmt.Sprintf("_%s", lowS))
		} else {
			s = append(s, string(k[i]))
		}
	}
	return strings.ToUpper(strings.Join(s, ""))
}

func (p *PaginationQ) excludeHumpReplace(k string) bool {
	if p.excludeHumpReplaceMap == nil {
		return false
	}
	_, ok := p.excludeHumpReplaceMap[k]
	return ok
}

func (p *PaginationQ) SetExcludeHumpReplaceMap(excludeRows ...string) {
	if p.excludeHumpReplaceMap == nil {
		p.excludeHumpReplaceMap = map[string]struct{}{}
	}
	for _, excludeRow := range excludeRows {
		p.excludeHumpReplaceMap[excludeRow] = struct{}{}
	}
}

func (p *PaginationQ) Where(whereCondition string, args ...interface{}) WhereFunc {
	return func(tx *Yasql) *Yasql {
		return tx.Where(whereCondition, args...)
	}
}

func (p *PaginationQ) FilterBy(tx *Yasql) (*Yasql, error) {
	if p.Filter == "" {
		return tx, nil
	}
	filters := map[string]interface{}{}
	if err := json.Unmarshal([]byte(p.Filter), &filters); err != nil {
		return tx, err
	}
	query := tx
	for k, v := range filters {
		k := p.humpReplace(k)
		if value, ok := v.(string); ok {
			query = tx.Where(fmt.Sprintf("%s = ?", k), value)
		} else if value, ok := v.(bool); ok {
			query = tx.Where(fmt.Sprintf("%s = %v", k, value))
		} else if value, ok := v.([]interface{}); ok {
			if len(value) > 0 {
				elem := value[0]
				str := ""
				switch elem.(type) {
				case string:
					for _, val := range value {
						str += fmt.Sprintf("'%s',", val)
					}
				default:
					for _, val := range value {
						str += fmt.Sprintf("%v,", val)
					}
				}
				query = tx.Where(fmt.Sprintf("%s in (%s)", k, str[0:len(str)-1]))
			}
		} else {
			log.Yasql.Errorf("unknow filter args: %s: %v", k, v)
			continue
		}
	}
	return query, nil
}

func (p *PaginationQ) SearchBy(tx *Yasql) (*Yasql, error) {
	if p.Search == "" {
		return tx, nil
	}
	search := map[string]string{}
	if err := json.Unmarshal([]byte(p.Search), &search); err != nil {
		return tx, err
	}
	query := tx
	for k, v := range search {
		if v == "" {
			continue
		}
		k := p.humpReplace(k)
		query = tx.Where(fmt.Sprintf("(%s like ? or %s like ?)", k, k), fmt.Sprintf("%%%s%%", v), fmt.Sprintf("%%%s%%", strings.ToUpper(v)))
	}
	return query, nil
}

func (p *PaginationQ) SortBy(tx *Yasql) *Yasql {
	if p.Sort == "" {
		return tx
	}
	var order Order
	if strings.ToLower(p.Order) == "asc" {
		order = ASC
	} else {
		order = DESC
	}
	k := p.humpReplace(p.Sort)
	sort := &Sort{ColName: k, Order: order}
	return tx.SortBy(sort)
}

func (p *PaginationQ) RangeBy(tx *Yasql) error {
	if p.Range == "" {
		return nil
	}
	ranges := make(map[string]QueryRange)
	if err := json.Unmarshal([]byte(p.Range), &ranges); err != nil {
		return err
	}
	for field, r := range ranges {
		field = p.humpReplace(field)
		switch r.Datatype {
		case _timestamp, _int:
			if _, err := r.genRangeQuery(tx, field); err != nil {
				return err
			}
		default:
			log.Yasql.Errorf("not supported datatype %s: ", r.Datatype)
			return errors.New("datatype not supported")
		}
	}
	return nil
}

// SearchAll optimized pagination method for yasql
func (p *PaginationQ) Match(tx *Yasql, whereFuncs ...WhereFunc) error {
	if p.Select == nil {
		return ErrInvalidSelect
	}
	if p.Select.Table == "" {
		return ErrTableNotSet
	}
	if len(p.Select.Columns) == 0 {
		return ErrNoColumns
	}
	query := tx.Select(p.Select)
	return p.match(query, whereFuncs...)
}

func (p *PaginationQ) MatchRaw(tx *Yasql, whereFuncs ...WhereFunc) error {
	if p.SelectRaw == nil || p.SelectRaw.RawSql == "" {
		return ErrInvalidSelect
	}
	query := tx.SelectRaw(p.SelectRaw)
	return p.match(query, whereFuncs...)
}

func (p *PaginationQ) match(query *Yasql, whereFuncs ...WhereFunc) error {
	if p.Size == 0 {
		p.Size = 10
	}
	if p.Size < 0 {
		p.Size = math.MaxInt32
	}
	if p.Page < 1 {
		p.Page = 1
	}
	for _, where := range whereFuncs {
		query = where(query)
	}
	if _, err := p.FilterBy(query); err != nil {
		return err
	}
	if _, err := p.SearchBy(query); err != nil {
		return err
	}
	if err := p.RangeBy(query); err != nil {
		return err
	}
	p.SortBy(query)
	if p.Select == nil {
		if err := query.CountRaw(p.CountSql, &p.Total).Error(); err != nil {
			log.Yasql.Errorf("count raw err: %s", err.Error())
			return err
		}
	} else {
		if err := query.Count(p.Select, &p.Total).Error(); err != nil {
			log.Yasql.Errorf("count err: %s", err.Error())
			return err
		}
	}
	tp := int(p.Total) % p.Size
	if tp == 0 {
		p.TotalPage = int(p.Total) / p.Size
	} else {
		p.TotalPage = (int(p.Total) / p.Size) + 1
	}
	offset := p.Size * (p.Page - 1)
	err := query.Paging(p.Size, offset).Find(p.Data).Error()
	if err != nil {
		log.Yasql.Errorf("page err: %s", err.Error())
	}
	p.Ok = err == nil
	if err == nil && p.AfterSearch != nil {
		return p.AfterSearch(query, p.Data)
	}
	return err
}

func (r *QueryRange) queryRangeType() QueryRangeType {
	if r.Min == "" {
		if r.Max == "" {
			return _range_invalid
		}
		return _range_max_only
	}
	if r.Max == "" {
		return _range_min_only
	}
	return _range_both
}

func (r *QueryRange) genRangeQuery(tx *Yasql, field string) (*Yasql, error) {
	switch r.queryRangeType() {
	case _range_min_only:
		return r.genMinOnlyRange(tx, field)
	case _range_max_only:
		return r.genMaxOnlyRange(tx, field)
	case _range_both:
		return r.genMinMaxOnlyRange(tx, field)
	default:
		return tx, nil
	}
}

func (r *QueryRange) genMinOnlyRange(tx *Yasql, field string) (*Yasql, error) {
	min, err := strconv.ParseInt(r.Min, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		log.Yasql.Errorf("%s min not int", r.Min)
		return nil, errors.New("min not int")
	}
	minStr := r.Min
	if r.Datatype == _timestamp {
		timeFormat := time.Unix(min, 0).Format(timedef.TIME_FORMAT)
		minStr = fmt.Sprintf("TIMESTAMP('%s')", timeFormat)
	}
	// where := field + ">=" + minStr
	where := fmt.Sprintf("%s <= %s", minStr, field)
	tx.Where(where)
	return tx, nil
}

func (r *QueryRange) genMaxOnlyRange(tx *Yasql, field string) (*Yasql, error) {
	max, err := strconv.ParseInt(r.Max, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		return nil, errors.New("max not int")
	}
	maxStr := r.Max
	if r.Datatype == _timestamp {
		timeFormat := time.Unix(max, 0).Format(timedef.TIME_FORMAT)
		maxStr = fmt.Sprintf("TIMESTAMP('%s')", timeFormat)
	}
	// where := field + "<=" + maxStr
	where := fmt.Sprintf("%s <= %s", field, maxStr)
	tx.Where(where)
	return tx, nil
}

func (r *QueryRange) genMinMaxOnlyRange(tx *Yasql, field string) (*Yasql, error) {
	max, err := strconv.ParseInt(r.Max, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		log.Yasql.Errorf("prase %s to int err", r.Min)
		return nil, errors.New("max not int")
	}
	min, err := strconv.ParseInt(r.Min, constants.BASE_DECIMAL, constants.BIT_SIZE_64)
	if err != nil {
		log.Yasql.Errorf("prase %s to int err", r.Min)
		return nil, errors.New("min not int")
	}
	maxStr := r.Max
	minStr := r.Min
	if r.Datatype == _timestamp {
		minTimeFormat := time.Unix(min, 0).Format(timedef.TIME_FORMAT)
		maxTimeFormat := time.Unix(max, 0).Format(timedef.TIME_FORMAT)
		maxStr = fmt.Sprintf("TIMESTAMP('%s')", maxTimeFormat)
		minStr = fmt.Sprintf("TIMESTAMP('%s')", minTimeFormat)
	}
	// where := minStr + "<=" + field + " add " + field + "<=" + maxStr
	where := fmt.Sprintf("%s <=  %s AND %s <=  %s ", minStr, field, field, maxStr)
	tx.Where(where)
	return tx, nil
}
