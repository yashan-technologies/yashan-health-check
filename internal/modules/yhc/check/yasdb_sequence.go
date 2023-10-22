package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbSequenceNoAvailable() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SEQUENCE_NO_AVAILABLE)
	defer c.fillResult(data)
	return
}
