package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_CONTROLFILE = "select * from v$controlfile;"
)

func (c *YHCChecker) GetYasdbControlFile() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_CONTROLFILE, SQL_QUERY_CONTROLFILE)
	defer c.fillResult(data)
	return
}
