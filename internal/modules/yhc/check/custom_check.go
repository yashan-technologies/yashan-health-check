package check

import (
	"fmt"

	"yhc/defs/bashdef"
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
)

func (c *YHCChecker) GenCustomCheckFunc(metric *confdef.YHCMetric) (fn func() error, err error) {
	if metric.Default {
		return nil, fmt.Errorf("metric %s is not a custom metric", metric.Name)
	}
	switch metric.MetricType {
	case confdef.MT_BASH:
		return c.genCustomBashFunc(metric), nil
	case confdef.MT_SQL:
		return c.genCustomSQLFunc(metric), nil
	default:
		return nil, fmt.Errorf("unsupport metric type %s", metric.MetricType)
	}
}

func (c *YHCChecker) genCustomBashFunc(metric *confdef.YHCMetric) (fn func() error) {
	fn = func() (err error) {
		data := &define.YHCItem{
			Name: define.MetricName(metric.Name),
		}
		defer c.fillResult(data)

		log := log.Module.M(metric.Name)
		execer := execerutil.NewExecer(log)

		ret, stdout, stderr := execer.Exec(bashdef.CMD_BASH, "-c", metric.Command)
		if ret != 0 {
			err = fmt.Errorf("failed to exec command %s, err: %v", metric.Command, stderr)
			log.Error(err)
			return
		}
		data.Details = stdout
		return
	}
	return
}

func (c *YHCChecker) genCustomSQLFunc(metric *confdef.YHCMetric) (fn func() error) {
	fn = func() (err error) {
		data, err := c.queryMultiRows(define.MetricName(metric.Name))
		defer c.fillResult(data)
		return
	}
	return
}
