package check

import (
	"fmt"
	"time"

	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
)

const (
	KEY_HIT_RATE = "HIT_RATE"
)

func (c *YHCChecker) GetYasdbHistoryBufferHitRate(name string) (err error) {
	var datas []*define.YHCItem
	defer c.fillResults(datas...)

	logger := log.Module.M(string(define.METRIC_YASDB_HISTORY_BUFFER_HIT_RATE))
	for _, yasdb := range c.GetCheckNodes(logger) {
		data := &define.YHCItem{Name: define.METRIC_YASDB_HISTORY_BUFFER_HIT_RATE, NodeID: yasdb.NodeID}
		datas = append(datas, data)

		var dbTimes []map[string]string
		dbTimes, err = yasdb.QueryMultiRows(
			fmt.Sprintf(define.SQL_QUERY_HISTORY_BUFFER_HIT_RATE,
				c.formatFunc(c.base.Start),
				c.formatFunc(c.base.End)),
			confdef.GetYHCConf().SqlTimeout)
		if err != nil {
			err = yaserr.Wrap(err)
			logger.Error(err)
			data.Error = err.Error()
			continue
		}
		content := make(define.WorkloadOutput)
		for _, row := range dbTimes {
			t, err := time.ParseInLocation(timedef.TIME_FORMAT, row[KEY_SNAP_TIME], time.Local)
			if err != nil {
				logger.Errorf("parse time %s failed: %s", row[KEY_SNAP_TIME], err)
				continue
			}
			content[t.Unix()] = define.WorkloadItem{KEY_HIT_RATE: row}
		}
		data.Details = content
	}
	return
}
