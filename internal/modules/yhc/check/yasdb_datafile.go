package check

import (
	"os"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
)

const (
	KEY_FILE_NAME        = "FILE_NAME"
	KEY_FILE_PERMISSION  = "FILE_PERMISSION"
	KEY_INCREASE_PERCENT = "INCREASE_PERCENT"
)

func (c *YHCChecker) GetYasdbDataFile(name string) (err error) {
	log := log.Module.M(string(define.METRIC_YASDB_DATAFILE))
	sql, err := c.getSQL(define.METRIC_YASDB_DATAFILE)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		return
	}
	metric, err := c.getMetric(define.METRIC_YASDB_DATAFILE)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", define.METRIC_YASDB_DATAFILE, err)
		return
	}

	var datas []*define.YHCItem
	for _, yasdb := range c.GetCheckNodes(log) {
		data := &define.YHCItem{Name: define.METRIC_YASDB_DATAFILE, NodeID: yasdb.NodeID}
		var res []map[string]string
		res, err = yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
		if err != nil {
			err = yaserr.Wrap(err)
			log.Error(err)
			data.Error = err.Error()
			continue
		}
		for _, r := range res {
			if yasdb.IsLocal {
				dataFile := r[KEY_FILE_NAME]
				if fileInfo, err := os.Stat(dataFile); err != nil {
					r[KEY_FILE_PERMISSION] = err.Error()
				} else {
					r[KEY_FILE_PERMISSION] = fileInfo.Mode().String()
				}
			} else {
				r[KEY_FILE_PERMISSION] = "未知"
			}
		}
		data.Details = c.convertMultiSqlData(metric, res)
		datas = append(datas, data)
	}
	c.fillResults(datas...)
	return
}
