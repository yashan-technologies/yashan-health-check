package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbTopSQLByCPUTime() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TOP_SQL_BY_CPU_TIME)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTopSQLByBufferGets() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TOP_SQL_BY_BUFFER_GETS)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTopSQLByDiskReads() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TOP_SQL_BY_DISK_READS)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTopSQLByParseCalls() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_TOP_SQL_BY_PARSE_CALLS)
	defer c.fillResult(data)
	return
}
