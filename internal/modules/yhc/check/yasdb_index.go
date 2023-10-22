package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbIndexBlevel() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_BLEVEL)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbIndexColumn() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_COLUMN)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbIndexInvisible() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_INVISIBLE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbIndexOversized() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_OVERSIZED)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTableIndexNotTogether() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_TABLE_INDEX_NOT_TOGETHER)
	defer c.fillResult(data)
	return
}
