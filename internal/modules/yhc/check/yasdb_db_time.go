package check

import (
	"time"

	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	KEY_DB_TIME_MS = "DB_TIME_MS"
	KEY_SNAP_TIME  = "SNAP_TIME"
)

func (c *YHCChecker) GetYasdbHistoryDBTime() (err error) {
	data := &define.YHCItem{Name: define.METRIC_YASDB_HISTORY_DB_TIME}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_YASDB_HISTORY_DB_TIME))
	yasdb := yasdbutil.NewYashanDB(logger, c.base.DBInfo)
	dbTimes, err := yasdb.QueryMultiRows(define.SQL_QUERY_HISTORY_DB_TIME, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		logger.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_HISTORY_DB_TIME, err)
		data.Error = err.Error()
		return
	}
	content := make(define.WorkloadOutput)
	for _, row := range dbTimes {
		t, err := time.Parse(timedef.TIME_FORMAT, row[KEY_SNAP_TIME])
		if err != nil {
			logger.Errorf("parse time %s failed: %s", row[KEY_SNAP_TIME], err)
			continue
		}
		content[t.Unix()] = define.WorkloadItem{KEY_DB_TIME_MS: row}
	}
	data.Details = content
	return
}
