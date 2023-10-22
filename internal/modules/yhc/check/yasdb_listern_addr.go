package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbListenAddr() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_LISTEN_ADDR)
	defer c.fillResult(data)
	return
}
