package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbReplicationStatus() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_REPLICATION_STATUS)
	defer c.fillResult(data)
	return
}
