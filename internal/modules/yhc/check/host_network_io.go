package check

import (
	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
)

func (c *YHCChecker) GetHostHistoryNetworkIO() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_HISTORY_NETWORK_IO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_HISTORY_NETWORK_IO))
	resp, err := c.hostHistoryWorkload(log, define.METRIC_HOST_HISTORY_NETWORK_IO)
	if err != nil {
		log.Error("failed to get host network io info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	data.Details = resp
	return
}

func (c *YHCChecker) GetHostCurrentNetworkIO() (err error) {
	data := &define.YHCItem{
		Name:     define.METRIC_HOST_CURRENT_NETWORK_IO,
		DataType: define.DATATYPE_SAR,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_CURRENT_NETWORK_IO))
	hasSar := c.CheckSarAccess() == nil
	if !hasSar {
		data.DataType = define.DATATYPE_GOPSUTIL
	}
	resp, err := c.hostCurrentWorkload(log, define.METRIC_HOST_CURRENT_NETWORK_IO, hasSar)
	if err != nil {
		log.Error("failed to get host network io info, err: %s", err.Error())
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
