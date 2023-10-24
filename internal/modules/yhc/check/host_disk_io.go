package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
)

func (c *YHCChecker) GetHostHistoryDiskIO(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_DISK_IO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_DISK_IO))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_DISK_IO)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentDiskIO(name string) (err error) {
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
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}
