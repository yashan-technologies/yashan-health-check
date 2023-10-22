package check

import (
	"time"

	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"github.com/shirou/gopsutil/host"
)

const (
	KEY_BOOT_TIME             = "bootTime"
	KEY_UP_TIME               = "uptime"
	KEY_HOST_ID               = "hostid"
	KEY_VIRTUALIZATION_SYSTEM = "virtualizationSystem"
	KEY_VIRTUALIZATION_ROLE   = "virtualizationRole"
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
	detail, err := c.convertObjectData(host)
	if err != nil {
		log.Errorf("failed to covert host info, err: %v", err)
		data.Error = err.Error()
		return
	}
	data.Details = c.dealHostInfo(detail)
	return
}

func (c *YHCChecker) dealHostInfo(res map[string]interface{}) map[string]interface{} {
	delete(res, KEY_VIRTUALIZATION_ROLE)
	delete(res, KEY_VIRTUALIZATION_SYSTEM)
	delete(res, KEY_HOST_ID)
	bootTime := res[KEY_BOOT_TIME].(float64)
	res[KEY_BOOT_TIME] = time.Unix(int64(bootTime), 0).Format(timedef.TIME_FORMAT)
	upTime := res[KEY_UP_TIME].(float64)
	res[KEY_UP_TIME] = (time.Duration(int64(upTime)) * time.Second).String()
	return res
}
