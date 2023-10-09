package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbHighFrequencySQL() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_HIGH_FREQUENCY_SQL)
	defer c.fillResult(data)
	return
}
