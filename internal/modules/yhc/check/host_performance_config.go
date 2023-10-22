package check

import (
	"strings"

	"yhc/defs/bashdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"

	"github.com/shirou/gopsutil/mem"
)

const (
	KEY_HUGE_PAGE_ENABLED   = "hugePageEnabled"
	KEY_SWAP_MEMORY_ENABLED = "swapMemoryEnabled"

	HUGE_PAGE_DISABLED = "never"
)

func (c *YHCChecker) GetHugePageEnabled() (err error) {
	data := &define.YHCItem{Name: define.METRIC_HOST_HUGE_PAGE}
	defer c.fillResult(data)

	command := bashdef.CMD_GREP + ` -oP '(?<=\[).+?(?=\])' /sys/kernel/mm/transparent_hugepage/enabled`
	logger := log.Module.M(string(define.METRIC_HOST_HUGE_PAGE))
	execer := execerutil.NewExecer(logger)
	ret, stdout, stderr := execer.Exec(bashdef.CMD_BASH, "-c", command)
	if ret != 0 {
		data.Error = stderr
		return
	}
	enabled := STR_FALSE
	if !strings.Contains(stdout, HUGE_PAGE_DISABLED) {
		enabled = STR_TRUE
	}
	data.Details = map[string]string{
		KEY_HUGE_PAGE_ENABLED: enabled,
	}
	return
}

func (c *YHCChecker) GetSwapMemoryEnabled() (err error) {
	data := &define.YHCItem{Name: define.METRIC_HOST_SWAP_MEMORY}
	defer c.fillResult(data)

	logger := log.Module.M(string(define.METRIC_HOST_SWAP_MEMORY))

	swapMemory, err := mem.SwapMemory()
	if err != nil {
		data.Error = err.Error()
		logger.Errorf("get swap memory failed: %s", err)
		return
	}
	enabled := STR_FALSE
	if swapMemory.Total != 0 {
		enabled = STR_TRUE
	}
	data.Details = map[string]string{
		KEY_SWAP_MEMORY_ENABLED: enabled,
	}
	return
}
