package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbAuditCleanupTask() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_AUDIT_CLEANUP_TASK)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbAuditFileSize() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_SECURITY_AUDIT_FILE_SIZE)
	defer c.fillResult(data)
	return
}
