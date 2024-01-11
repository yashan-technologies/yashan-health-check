package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yasutil/size"
	"github.com/shirou/gopsutil/mem"
)

const (
	KEY_MEMORY_TYPE               = "type"
	KEY_MEMORY_TOTAL              = "total"
	KEY_MEMORY_USED               = "used"
	KEY_MEMORY_FREE               = "free"
	KEY_MEMORY_SHARED             = "shared"
	KEY_MEMORY_BUFFERS_AND_CACHED = "buffers_cached"
	KEY_MEMORY_AVAILABLE          = "available"

	SYSTEM_MEMORY_TYPE = "system"
	SWAP_MEMORY_TYPE   = "swap"
)

func (c *YHCChecker) GetHostMemoryInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_MEMORY_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_MEMORY_INFO))
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return err
	}
	data.Details = c.dealMemoryData(memInfo)
	return
}

func (c *YHCChecker) dealMemoryData(memory *mem.VirtualMemoryStat) (res []map[string]any) {
	res = append(res,
		map[string]any{
			KEY_MEMORY_TYPE:               SYSTEM_MEMORY_TYPE,
			KEY_MEMORY_TOTAL:              size.GenHumanReadableSize(float64(memory.Total), decimal),
			KEY_MEMORY_USED:               size.GenHumanReadableSize(float64(memory.Used), decimal),
			KEY_MEMORY_FREE:               size.GenHumanReadableSize(float64(memory.Free), decimal),
			KEY_MEMORY_SHARED:             size.GenHumanReadableSize(float64(memory.Shared), decimal),
			KEY_MEMORY_BUFFERS_AND_CACHED: size.GenHumanReadableSize(float64(memory.Buffers+memory.Cached), decimal),
			KEY_MEMORY_AVAILABLE:          size.GenHumanReadableSize(float64(memory.Available), decimal),
		},
		map[string]any{
			KEY_MEMORY_TYPE:               SWAP_MEMORY_TYPE,
			KEY_MEMORY_TOTAL:              size.GenHumanReadableSize(float64(memory.SwapTotal), decimal),
			KEY_MEMORY_USED:               size.GenHumanReadableSize(float64(memory.SwapTotal-memory.SwapFree-memory.SwapCached), decimal),
			KEY_MEMORY_FREE:               size.GenHumanReadableSize(float64(memory.SwapFree), decimal),
			KEY_MEMORY_SHARED:             "/",
			KEY_MEMORY_BUFFERS_AND_CACHED: size.GenHumanReadableSize(float64(memory.SwapCached), decimal),
			KEY_MEMORY_AVAILABLE:          "/",
		},
	)
	return
}
