package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_PARAMETER = "select name, value from v$parameter where value is not null;"
)

func (c *YHCChecker) GetYasdbParameter() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_PARAMETER, SQL_QUERY_PARAMETER)
	defer c.fillResult(data)
	return
}
