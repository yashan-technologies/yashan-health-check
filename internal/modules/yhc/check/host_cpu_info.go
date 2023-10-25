package check

import (
	"errors"
	"fmt"
	"strings"

	"yhc/defs/bashdef"
	"yhc/defs/runtimedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/osutil"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yaslog"
	"github.com/shirou/gopsutil/cpu"
)

const (
	KEY_CPU_PHYSICAL_CORES = "physicalCores"
	KEY_CPU_LOGICAL_CORES  = "logicalCores"
	KEY_CPU_MODEL_NAME     = "modelName"
	KEY_CPU_VERDOR_ID      = "vendorId"
	KEY_CPU_FLAGS          = "flags"
	KEY_CPU_GHZ            = "GHz"

	KY_PRODUCT_INFO  = "/etc/.productinfo"
	KEY_KY_MAX_SPEED = "Max Speed:"
)

func (c *YHCChecker) GetHostCPUInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_CPU_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_CPU_INFO))
	cpuInfos, err := cpu.Info()
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = c.countCPUInfo(log, cpuInfos)
	return
}

func (c *YHCChecker) countCPUInfo(log yaslog.YasLog, cpuInfos []cpu.InfoStat) (res map[string]interface{}) {
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
	if runtimedef.GetOSRelease().Id == osutil.KYLIN_ID {
		freq, err := c.getKyCPUFrequency(log)
		if err != nil {
			delete(res, KEY_CPU_GHZ)
			log.Error(err)
			return
		}
		res[KEY_CPU_GHZ] = freq
	}
	return res
}

func (c *YHCChecker) getKyCPUFrequency(log yaslog.YasLog) (string, error) {
	execer := execerutil.NewExecer(log)
	ret, stdout, stderr := execer.EnvExec(_envs, bashdef.CMD_DMIDECODE, "-t", "processor")
	if ret != 0 {
		err := fmt.Errorf("failed to get kylin CPU frequency info, err: %s", stderr)
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(stdout), stringutil.STR_NEWLINE)
	if len(lines) == 0 {
		err := errors.New("failed to get dmidecode info")
		return "", err
	}
	arr := []string{}
	for _, line := range lines {
		if strings.Contains(line, KEY_KY_MAX_SPEED) {
			arr = append(arr, strings.TrimSpace(line))
		}
	}
	if len(arr) == 0 {
		err := errors.New("failed to get kylin CPU frequency info")
		return "", err
	}
	res := strings.TrimPrefix(arr[0], KEY_KY_MAX_SPEED)
	return strings.TrimSpace(res), nil
}
