package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbInstance() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_INSTANCE)
	defer c.fillResult(data)
	return
}
