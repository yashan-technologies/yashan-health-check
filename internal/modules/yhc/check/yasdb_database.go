package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_DATABASE = "select database_name, status as database_status, log_mode, open_mode, database_role, protection_mode, create_time from v$database;"
)

func (c *YHCChecker) GetYasdbDatabase() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_DATABASE, SQL_QUERY_DATABASE)
	defer c.fillResult(data)
	return
}
