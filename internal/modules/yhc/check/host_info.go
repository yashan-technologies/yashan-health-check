package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"github.com/shirou/gopsutil/host"
)

func (c *YHCChecker) GetHostInfo() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_INFO))
	host, err := host.Info()
	if err != nil {
		log.Errorf("failed to get host info, err: %v", err)
		data.Error = err.Error()
		return
	}
	data.Details = host
	return
}
