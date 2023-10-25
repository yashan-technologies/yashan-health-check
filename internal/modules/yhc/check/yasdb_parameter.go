package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/yasdbutil"

	"git.yasdb.com/go/yaserr"
)

const (
	KEY_PARAMETER_NAME  = "NAME"
	KEY_PARAMETER_VALUE = "VALUE"
)

func (c *YHCChecker) GetYasdbParameter(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_YASDB_PARAMETER,
	}
	defer c.fillResult(data)
	log := log.Module.M(string(define.METRIC_YASDB_PARAMETER))
	sql, err := c.getSQL(define.METRIC_YASDB_PARAMETER)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return err
	}
	metric, err := c.getMetric(define.METRIC_YASDB_PARAMETER)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return err
	}
	yasdb := yasdbutil.NewYashanDB(log, c.base.DBInfo)
	res, err := yasdb.QueryMultiRows(sql, confdef.GetYHCConf().SqlTimeout)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
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
