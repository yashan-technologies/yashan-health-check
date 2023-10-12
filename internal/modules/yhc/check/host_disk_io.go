package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
)

// todo: 处理chart图的时候要额外考虑图表怎么画

func (c *YHCChecker) GetHostHistoryDiskIO() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_DISK_IO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_DISK_IO))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_DISK_IO)
	if err != nil {
		log.Error("failed to get host disk io info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentDiskIO() (err error) {
	data := &define.YHCItem{
		Name:     define.METRIC_HOST_CURRENT_DISK_IO,
		DataType: define.DATATYPE_SAR,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_CURRENT_DISK_IO))
	hasSar := c.CheckSarAccess() == nil
	if !hasSar {
		data.DataType = define.DATATYPE_GOPSUTIL
	}
	resp, err := c.hostCurrentWorkload(log, define.METRIC_HOST_CURRENT_DISK_IO, hasSar)
	if err != nil {
		log.Error("failed to get host disk io info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}
