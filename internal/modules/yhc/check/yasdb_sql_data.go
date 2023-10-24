package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbSingleRowData(name string) (err error) {
	data, err := c.querySingleRow(define.MetricName(name))
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbMultiRowData(name string) (err error) {
	data, err := c.queryMultiRows(define.MetricName(name))
	defer c.fillResult(data)
	return
}
