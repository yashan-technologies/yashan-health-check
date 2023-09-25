package check

import "yhc/internal/modules/yhc/check/define"

// todo:sql有问题，已反馈，待给出正确的sql
func (c *YHCChecker) GetYasdbWaitEvent() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_WAIT_EVENT)
	defer c.fillResult(data)
	return
}
