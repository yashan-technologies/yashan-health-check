package check

import (
	"yhc/internal/modules/yhc/check/define"
)

func (c *YHCChecker) GetYasdbTableWithToMuchColumns() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TABLE_WITH_TOO_MUCH_COLUMNS)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTableWithToMuchIndexes() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TABLE_WITH_TOO_MUCH_INDEXES)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbPartitionedTableWithoutPartitionedIndexes() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_PARTITIONED_TABLE_WITHOUT_PARTITIONED_INDEXES)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTableWithRowSizeExceedsBlockSize() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbPartitionedTableWithNumberOfHashPartitionsIsNotAPowerOfTwo() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_PARTITIONED_TABLE_WITH_NUMBER_OF_HASH_PARTITIONS_IS_NOT_A_POWER_OF_TWO)
	defer c.fillResult(data)
	return
}
