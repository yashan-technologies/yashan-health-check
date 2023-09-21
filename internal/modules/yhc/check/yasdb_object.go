package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_TOTAL_OBJECT = "select count(*) as total_count from dba_objects;"
	SQL_QUERY_OWNER_OBJECT = `SELECT owner,object_type, COUNT(*) AS object_count FROM dba_objects
    WHERE owner NOT IN ('SYS', 'SYSTEM') AND object_type NOT LIKE 'BIN$%'
    GROUP BY owner, object_type
    ORDER BY owner,object_type;`
	SQL_QUERY_TABLESPACE_OBJECT = `SELECT tablespace_name, COUNT(*) AS object_count FROM dba_segments
    WHERE segment_type IN ('TABLE', 'INDEX', 'VIEW', 'SEQUENCE')
    GROUP BY tablespace_name
    ORDER BY tablespace_name;`
)

func (c *YHCChecker) GetYasdbObjectCount() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_COUNT, SQL_QUERY_TOTAL_OBJECT)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbOwnerObject() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_OWNER, SQL_QUERY_OWNER_OBJECT)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbTablespaceObject() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_OBJECT_TABLESPACE, SQL_QUERY_TABLESPACE_OBJECT)
	defer c.fillResult(data)
	return
}
