package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
)

func (c *YHCChecker) GetHostHistoryNetworkIO(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_NETWORK_IO,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_NETWORK_IO))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_NETWORK_IO)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentNetworkIO(name string) (err error) {
	data := &define.YHCItem{
		Name:     define.METRIC_HOST_CURRENT_NETWORK_IO,
		DataType: define.DATATYPE_SAR,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_CURRENT_NETWORK_IO))
	hasSar := c.CheckSarAccess() == nil
	if !hasSar {
		data.DataType = define.DATATYPE_GOPSUTIL
	}
	resp, err := c.hostCurrentWorkload(log, define.METRIC_HOST_CURRENT_NETWORK_IO, hasSar)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = c.filterDiscardNetwork(resp)
	return
}

func (c *YHCChecker) filterDiscardNetwork(input define.WorkloadOutput) (res define.WorkloadOutput) {
	res = make(define.WorkloadOutput)
	for timestamp, item := range input {
		for key := range item {
			if confdef.IsDiscardNetwork(key) {
				delete(item, key)
			}
		}
		res[timestamp] = item
	}
	return
}
