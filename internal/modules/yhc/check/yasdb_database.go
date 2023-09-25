package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbDatabase() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_DATABASE)
	defer c.fillResult(data)
	return
}
