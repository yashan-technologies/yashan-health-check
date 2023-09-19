package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
)

func (c *YHCChecker) GetHostCPUUsage() (err error) {
	data := define.YHCItem{
		Name:     define.METRIC_HOST_CPU_USAGE,
		Children: make(map[define.MetricName]define.YHCItem),
	}
	defer c.fillResult(&data)

	log := log.Module.M(string(define.METRIC_HOST_CPU_USAGE))
	resp, err := c.hostWorkload(log, define.METRIC_HOST_CPU_USAGE)
	if err != nil {
		log.Error("failed to get host cpu usage info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Children[KEY_HISTORY] = define.YHCItem{
		Error:    resp.Errors[KEY_HISTORY],
		Details:  resp.Data[KEY_HISTORY],
		DataType: resp.DataType,
	}
	data.Children[KEY_CURRENT] = define.YHCItem{
		Error:    resp.Errors[KEY_CURRENT],
		Details:  resp.Data[KEY_CURRENT],
		DataType: resp.DataType,
	}
	return
}
