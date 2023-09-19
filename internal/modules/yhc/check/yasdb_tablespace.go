package check

import (
	"fmt"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	SQL_QUERY_TABLESPACE = `SELECT TABLESPACE_NAME, CONTENTS, STATUS, ALLOCATION_TYPE AS ALLOCATIONTYPE
    , TOTAL_BYTES - USER_BYTES AS USEDBYTES, TOTAL_BYTES AS TOTALBYTES
    , (TOTAL_BYTES - USER_BYTES) / TOTAL_BYTES * 100 AS RATE
    FROM SYS.DBA_TABLESPACES;`
	SQL_QUERY_TABLESPACE_DATA_PERCENTAGE_FORMATER = `SELECT A.TABLESPACE_NAME, A.B1/B.B2*100 AS DATA_PERCENTAGE FROM 
    (SELECT TABLESPACE_NAME,SUM(BYTES) AS B1 FROM dba_segments WHERE SEGMENT_TYPE LIKE 'TABLE%%' GROUP BY TABLESPACE_NAME ) A,
    (SELECT TABLESPACE_NAME,TOTAL_BYTES AS B2 FROM DBA_TABLESPACES) B WHERE (A.TABLESPACE_NAME=B.TABLESPACE_NAME AND A.TABLESPACE_NAME ='%s');`

	KEY_TABLESPACE_DATA_PERCENTAGE = "DATA_PERCENTAGE"
)

func (c *YHCChecker) GetYasdbTablespace() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_TABLESPACE,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_TABLESPACE))
	yasdb := yasdbutil.NewYashanDB(log, &c.Yasdb)
	tablespaces, err := yasdb.QueryMultiRows(SQL_QUERY_TABLESPACE, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql %s, err: %v", SQL_QUERY_SESSION, err)
		data.Error = err.Error()
		return
	}
	for _, tablespace := range tablespaces {
		tablespaceName := tablespace["TABLESPACE_NAME"]
		queryDataPercentageSQL := fmt.Sprintf(SQL_QUERY_TABLESPACE_DATA_PERCENTAGE_FORMATER, tablespaceName)
		dataPercentage, e := yasdb.QueryMultiRows(queryDataPercentageSQL, confdef.GetYHCConf().SqlTimeout)
		if e != nil {
			log.Errorf("failed to get data with sql %s, err: %v", queryDataPercentageSQL, err)
			continue
		}
		if len(dataPercentage) <= 0 {
			continue
		}
		tablespace[KEY_TABLESPACE_DATA_PERCENTAGE] = dataPercentage[0][KEY_TABLESPACE_DATA_PERCENTAGE]
	}
	data.Details = tablespaces
	return
}
