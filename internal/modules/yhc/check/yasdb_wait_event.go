package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_WAIT_EVENT = `SELECT count(s.WAIT_EVENT) current_waits FROM sys.v$system_event se, sys.v$session s WHERE se.EVENT = s.WAIT_EVENT
    AND se.event not in ('SQL*Net message from client',
    'SQL*Net more data from client',
    'pmon timer',
    'rdbms ipc message',
    'rdbms ipc reply',
    'smon timer');`
)

// todo:sql有问题，已反馈，待给出正确的sql
func (c *YHCChecker) GetYasdbWaitEvent() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_WAIT_EVENT, SQL_QUERY_WAIT_EVENT)
	defer c.fillResult(data)
	return
}
