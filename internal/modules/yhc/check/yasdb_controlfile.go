package check

import (
	"yhc/internal/modules/yhc/check/define"
)

func (c *YHCChecker) GetYasdbControlFile() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_CONTROLFILE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbControlFileCount() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_CONTROLFILE_COUNT)
	defer c.fillResult(data)
	return
}
