package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
)

func (c *YHCChecker) GetHostHistoryMemoryUsage(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_MEMORY_USAGE,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_MEMORY_USAGE))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_MEMORY_USAGE)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentMemoryUsage(name string) (err error) {
	data := &define.YHCItem{
		Name:     define.METRIC_HOST_CURRENT_MEMORY_USAGE,
		DataType: define.DATATYPE_SAR,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_CURRENT_MEMORY_USAGE))
	hasSar := c.CheckSarAccess() == nil
	if !hasSar {
		data.DataType = define.DATATYPE_GOPSUTIL
	}
	resp, err := c.hostCurrentWorkload(log, define.METRIC_HOST_CURRENT_MEMORY_USAGE, hasSar)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}
