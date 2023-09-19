package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_INSTANCE = "select status as instance_status, version, startup_time from v$instance;"
)

func (c *YHCChecker) GetYasdbInstance() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_INSTANCE, SQL_QUERY_INSTANCE)
	defer c.fillResult(data)
	return
}
