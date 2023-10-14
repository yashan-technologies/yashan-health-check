package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbLongRunningTransaction() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_LONG_RUNNING_TRANSACTION)
	defer c.fillResult(data)
	return
}
