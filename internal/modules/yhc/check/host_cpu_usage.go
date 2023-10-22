package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
)

func (c *YHCChecker) GetHostHistoryCPUUsage() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_CPU_USAGE,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_CPU_USAGE))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_CPU_USAGE)
	if err != nil {
		log.Error("failed to get host cpu usage info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentCPUUsage() (err error) {
	data := &define.YHCItem{
		Name:     define.METRIC_HOST_CURRENT_CPU_USAGE,
		DataType: define.DATATYPE_SAR,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_CPU_USAGE))
	hasSar := c.CheckSarAccess() == nil
	if !hasSar {
		data.DataType = define.DATATYPE_GOPSUTIL
	}
	resp, err := c.hostCurrentWorkload(log, define.METRIC_HOST_HISTORY_CPU_USAGE, hasSar)
	if err != nil {
		log.Error("failed to get host cpu usage info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}
