package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_LISTEN_ADDR = `select VALUE as LISTEN_ADDR from v$parameter where name = 'LISTEN_ADDR';`
)

func (c *YHCChecker) GetYasdbListenAddr() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_LISTEN_ADDR, SQL_QUERY_LISTEN_ADDR)
	defer c.fillResult(data)
	return
}
