package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbInvalidObject() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INVALID_OBJECT)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbInvisibleIndex() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INVISIBLE_INDEX)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbDisableConstraint() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_DISABLED_CONSTRAINT)
	defer c.fillResult(data)
	return
}
