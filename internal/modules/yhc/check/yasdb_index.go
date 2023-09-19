package check

import "yhc/internal/modules/yhc/check/define"

const (
	SQL_QUERY_INDEX_BLEVEL    = "select OWNER, INDEX_NAME, BLEVEL from dba_indexes where BLEVEL>3;"
	SQL_QUERY_INDEX_COLUMN    = "select INDEX_OWNER, INDEX_NAME, count(*) from dba_ind_columns group by INDEX_OWNER,INDEX_NAME having count(*) > 10;"
	SQL_QUERY_INDEX_INVISIBLE = "select OWNER, INDEX_NAME, TABLE_OWNER, TABLE_NAME FROM dba_indexes where owner<> 'SYS' and VISIBILITY <> 'VISIBLE'"
)

func (c *YHCChecker) GetYasdbIndexBlevel() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_BLEVEL, SQL_QUERY_INDEX_BLEVEL)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbIndexColumn() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_COLUMN, SQL_QUERY_INDEX_COLUMN)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbIndexInvisible() (err error) {
	data, err := c.queryMultiRows(define.METRIC_YASDB_INDEX_INVISIBLE, SQL_QUERY_INDEX_INVISIBLE)
	defer c.fillResult(data)
	return
}
