package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbRunningTask() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TASK_RUNNING)
	defer c.fillResult(data)
	return
}
