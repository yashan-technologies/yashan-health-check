package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbObjectCount() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_COUNT)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbOwnerObject() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_OWNER)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTablespaceObject() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_TABLESPACE)
	defer c.fillResult(data)
	return
}
