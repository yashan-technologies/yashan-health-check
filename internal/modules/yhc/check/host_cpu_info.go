package check

import (
	"fmt"
	"strings"

	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"github.com/shirou/gopsutil/cpu"
)

const (
	KEY_CPU_PHYSICAL_CORES = "physicalCores"
	KEY_CPU_LOGICAL_CORES  = "logicalCores"
	KEY_CPU_MODEL_NAME     = "modelName"
	KEY_CPU_VERDOR_ID      = "vendorId"
	KEY_CPU_FLAGS          = "flags"
	KEY_CPU_GHZ            = "GHz"
)

func (c *YHCChecker) GetHostCPUInfo() (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_CPU_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_CPU_INFO))
	cpuInfos, err := cpu.Info()
	if err != nil {
		log.Errorf("failed to get host cpu info, err: %v", err)
		data.Error = err.Error()
		return
	}
	data.Details = c.countCPUInfo(cpuInfos)
	return
}

func (c *YHCChecker) countCPUInfo(cpuInfos []cpu.InfoStat) (res map[string]interface{}) {
	res = make(map[string]interface{})
	var physicalCores, logicalCores int
	tmp := make(map[string]struct{})
	for _, c := range cpuInfos {
		tmp[c.PhysicalID] = struct{}{}
		logicalCores += int(c.Cores)
	}
	physicalCores = len(tmp)
	if physicalCores <= 0 {
		return res
	}
	res[KEY_CPU_PHYSICAL_CORES] = physicalCores
	res[KEY_CPU_LOGICAL_CORES] = logicalCores
	res[KEY_CPU_MODEL_NAME] = cpuInfos[0].ModelName
	res[KEY_CPU_VERDOR_ID] = cpuInfos[0].VendorID
	res[KEY_CPU_FLAGS] = strings.Join(cpuInfos[0].Flags, ",")
	res[KEY_CPU_GHZ] = fmt.Sprintf("@%.2fGHz", cpuInfos[0].Mhz/1000)
	return res
}
