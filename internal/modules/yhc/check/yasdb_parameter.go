package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"
)

const (
	KEY_PARAMETER_NAME  = "NAME"
	KEY_PARAMETER_VALUE = "VALUE"
)

func (c *YHCChecker) GetYasdbParameter() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_PARAMETER,
	}
	defer c.fillResult(data)
	log := log.Module.M(string(define.METRIC_YASDB_PARAMETER))
	sql, err := c.getSQL(define.METRIC_YASDB_PARAMETER)
	if err != nil {
		log.Errorf("failed to get sql of %s, err: %v", define.METRIC_YASDB_PARAMETER, err)
		data.Error = err.Error()
		return err
	}
	metric, err := c.getMetric(define.METRIC_YASDB_PARAMETER)
	if err != nil {
		log.Errorf("failed to get metric by name %s, err: %v", define.METRIC_YASDB_PARAMETER, err)
		data.Error = err.Error()
		return err
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		log.Errorf("failed to get data with sql '%s', err: %v", sql, err)
		data.Error = err.Error()
		return err
	}
	detail := map[string]string{}
	for _, v := range res {
		detail[v[KEY_PARAMETER_NAME]] = v[KEY_PARAMETER_VALUE]
	}
	data.Details = c.convertSqlData(metric, detail)
	return
}
