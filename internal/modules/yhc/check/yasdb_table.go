package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
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
	data := &define.YHCItem{Name: define.METRIC_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_YASDB_SESSION))
	yasdb := yasdbutil.NewYashanDB(logger, c.base.DBInfo)
	tablesFromDBATabColumns, err := yasdb.QueryMultiRows(define.SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TAB_COLUMNS, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		logger.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TAB_COLUMNS, err)
		data.Error = err.Error()
		return
	}
	tablesFromDBATables, err := yasdb.QueryMultiRows(define.SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TABLES, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		logger.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_YASDB_TABLE_WITH_ROW_SIZE_EXCEEDS_BLOCK_SIZE_FROM_DBA_TABLES, err)
		data.Error = err.Error()
		return
	}
	tables := removeDuplicateMaps(append(tablesFromDBATabColumns, tablesFromDBATables...))
	data.Details = tables
	return
}

func (c *YHCChecker) GetYasdbPartitionedTableWithNumberOfHashPartitionsIsNotAPowerOfTwo() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_PARTITIONED_TABLE_WITH_NUMBER_OF_HASH_PARTITIONS_IS_NOT_A_POWER_OF_TWO)
	defer c.fillResult(data)
	return
}

func genMapString(m map[string]string) (s string) {
	for k, v := range m {
		s += k + v
	}
	return s
}

func removeDuplicateMaps(maps []map[string]string) []map[string]string {
	uniqueMaps := []map[string]string{}
	checkMap := map[string]bool{}

	for _, m := range maps {
		key := genMapString(m)
		if !checkMap[key] {
			checkMap[key] = true
			uniqueMaps = append(uniqueMaps, m)
		}
	}

	return uniqueMaps
}
