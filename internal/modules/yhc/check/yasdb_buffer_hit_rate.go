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
	KEY_HIT_RATE = "HIT_RATE"
)

func (c *YHCChecker) GetYasdbBufferHitRate() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_BUFFER_HIT_RATE)
	defer c.fillResult(data)
	return
}

func (c *YHCChecker) GetYasdbHistoryBufferHitRate() (err error) {
	data := &define.YHCItem{Name: define.METRIC_YASDB_HISTORY_BUFFER_HIT_RATE}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_YASDB_HISTORY_BUFFER_HIT_RATE))
	yasdb := yasdbutil.NewYashanDB(logger, c.base.DBInfo)
	dbTimes, err := yasdb.QueryMultiRows(define.SQL_QUERY_HISTORY_BUFFER_HIT_RATE, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		logger.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_HISTORY_BUFFER_HIT_RATE, err)
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
		content[t.Unix()] = define.WorkloadItem{KEY_HIT_RATE: row}
	}
	data.Details = content
	return
}
