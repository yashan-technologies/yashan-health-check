package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbCurrentUndoSize() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_UNDO_LOG_SIZE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbCurrentUndoBlock() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_UNDO_LOG_TOTAL_BLOCK)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbRunningTransactions() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_UNDO_LOG_RUNNING_TRANSACTIONS)
	defer c.fillResult(data)
	return
}
