package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbNoPackagePkgBody() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_PACKAGE_NO_PACKAGE_PACKAGE_BODY)
	defer c.fillResult(data)
	return
}
