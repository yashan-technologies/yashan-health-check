package check

import (
	"yhc/internal/modules/yhc/check/define"
)

func (c *YHCChecker) GetYasdbForeignKeysWithoutIndexes() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_FOREIGN_KEYS_WITHOUT_INDEXES)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbForeignKeysWithImplicitDataTypeConversion() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_FOREIGN_KEYS_WITH_IMPLICIT_DATA_TYPE_CONVERSION)
	defer c.fillResult(data)
	return
}
