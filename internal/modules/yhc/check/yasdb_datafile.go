package check

import (
	"os"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	KEY_FILE_NAME        = "FILE_NAME"
	KEY_FILE_PERMISSION  = "FILE_PERMISSION"
	KEY_INCREASE_PERCENT = "INCREASE_PERCENT"
)

func (c *YHCChecker) GetYasdbDataFile() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_DATAFILE,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_YASDB_DATAFILE))
	sql, err := c.getSQL(define.METRIC_YASDB_DATAFILE)
	if err != nil {
		log.Errorf("failed to get sql of %s, err: %v", define.METRIC_YASDB_DATAFILE, err)
		data.Error = err.Error()
		return
	}
	metric, err := c.getMetric(define.METRIC_YASDB_DATAFILE)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", define.METRIC_YASDB_DATAFILE, err)
		data.Error = err.Error()
		return
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql '%s', err: %v", sql, err)
		data.Error = err.Error()
		return
	}
	for _, r := range res {
		dataFile := r[KEY_FILE_NAME]
		if fileInfo, err := os.Stat(dataFile); err != nil {
			r[KEY_FILE_PERMISSION] = err.Error()
		} else {
			r[KEY_FILE_PERMISSION] = fileInfo.Mode().String()
		}
	}
	data.Details = c.convertMultiSqlData(metric, res)
	return
}
