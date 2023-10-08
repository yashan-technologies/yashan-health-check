package check

import (
	"fmt"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	KEY_TABLESPACE_DATA_PERCENTAGE = "DATA_PERCENTAGE"
)

func (c *YHCChecker) GetYasdbTablespace() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_TABLESPACE,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_TABLESPACE))
	metric, err := c.getMetric(define.METRIC_YASDB_TABLESPACE)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", define.METRIC_YASDB_TABLESPACE, err)
		data.Error = err.Error()
		return err
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	tablespaces, err := yasdb.QueryMultiRows(define.SQL_QUERY_TABLESPACE, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_TABLESPACE, err)
		data.Error = err.Error()
		return
	}
	for _, tablespace := range tablespaces {
		tablespaceName := tablespace["TABLESPACE_NAME"]
		queryDataPercentageSQL := fmt.Sprintf(define.SQL_QUERY_TABLESPACE_DATA_PERCENTAGE_FORMATER, tablespaceName)
		dataPercentage, e := yasdb.QueryMultiRows(queryDataPercentageSQL, confdef.GetYHCConf().SqlTimeout)
		if e != nil {
			log.Errorf("failed to get data with sql %s, err: %v", queryDataPercentageSQL, err)
			continue
		}
		if len(dataPercentage) <= 0 {
			tablespace[KEY_TABLESPACE_DATA_PERCENTAGE] = "0"
			continue
		}
		tablespace[KEY_TABLESPACE_DATA_PERCENTAGE] = dataPercentage[0][KEY_TABLESPACE_DATA_PERCENTAGE]
	}
	data.Details = c.convertMultiSqlData(metric, tablespaces)
	return
}
