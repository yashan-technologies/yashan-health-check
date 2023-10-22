package check

import (
	"fmt"
	"strconv"
	"time"

	"yhc/commons/constants"
	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yaslog"
)

const (
	KEY_DB_TIME_MS = "DB_TIME_MS"
	KEY_SNAP_TIME  = "SNAP_TIME"
	KEY_DB_TIMES   = "DB_TIMES"
)

func (c *YHCChecker) GetYasdbHistoryDBTime() (err error) {
	data := &define.YHCItem{Name: define.METRIC_YASDB_HISTORY_DB_TIME}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_YASDB_HISTORY_DB_TIME))
	yasdb := yasdbutil.NewYashanDB(logger, c.base.DBInfo)
	dbTimes, err := yasdb.QueryMultiRows(fmt.Sprintf(define.SQL_QUERY_SNAP_DB_TIMES, c.formatFunc(c.base.Start), c.formatFunc(c.base.End)), confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		logger.Errorf("failed to get data with sql %s, err: %v", define.SQL_QUERY_SNAP_DB_TIMES, err)
		data.Error = err.Error()
		return
	}
	content, e := c.getDBTimes(logger, dbTimes)
	if e != nil {
		err = yaserr.Wrap(e)
		logger.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = content
	return
}

func (c *YHCChecker) getDBTimes(log yaslog.YasLog, dbTimes []map[string]string) (define.WorkloadOutput, error) {
	content := make(define.WorkloadOutput)
	for i := 1; i < len(dbTimes); i++ {
		current := dbTimes[i]
		previous := dbTimes[i-1]
		currentDBTime, currentSnapTime, err := c.dbTime(log, current)
		if err != nil {
			return nil, err
		}
		previousDBTime, previousSnapTime, err := c.dbTime(log, previous)
		if err != nil {
			return nil, err
		}
		t := currentDBTime - previousDBTime
		if t <= 0 {
			continue
		}
		midSnapTime := previousSnapTime.Add(currentSnapTime.Sub(previousSnapTime) / 2)
		content[midSnapTime.Unix()] = define.WorkloadItem{KEY_DB_TIME_MS: map[string]interface{}{
			KEY_SNAP_TIME: midSnapTime.Format(timedef.TIME_FORMAT),
			KEY_DB_TIMES:  strconv.FormatFloat(t, 'f', 0, constants.BIT_SIZE_64),
		}}
	}
	return content, nil
}

func (c *YHCChecker) formatFunc(t time.Time) string {
	return t.Format(timedef.TIME_FORMAT)
}

func (c *YHCChecker) dbTime(log yaslog.YasLog, dbTime map[string]string) (timeNum float64, t time.Time, err error) {
	parseTime := func(s string) (t time.Time, err error) {
		t, err = time.ParseInLocation(timedef.TIME_FORMAT, s, time.Local)
		return
	}
	parseFloat := func(s string) (num float64, err error) {
		num, err = strconv.ParseFloat(s, constants.BIT_SIZE_64)
		return
	}
	t, err = parseTime(dbTime[KEY_SNAP_TIME])
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		return
	}
	timeNum, err = parseFloat(dbTime[KEY_DB_TIMES])
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		return
	}
	return
}
