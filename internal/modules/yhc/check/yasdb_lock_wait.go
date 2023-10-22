package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbTableLockWait() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_TABLE_LOCK_WAIT)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbRowLockWait() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_ROW_LOCK_WAIT)
	defer c.fillResult(data)
	return
}
