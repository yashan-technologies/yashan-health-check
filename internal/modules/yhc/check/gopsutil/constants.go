package gopsutil

import (
	"yhc/internal/modules/yhc/check/define"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type CpuUsage struct {
	cpu.TimesStat
}

type DiskIO struct {
	Dev                string  `json:"dev"`                          // dev name.
	Iops               uint64  `json:"iops,omitempty"`               // iops in process.
	SerialNumber       string  `json:"serialNumber,omitempty"`       // serial number.
	Label              string  `json:"label,omitempty"`              // label.
	KBReadSec          float64 `json:"kBReadSec,omitempty"`          // kBytes read per second.
	KBWriteSec         float64 `json:"kBWriteSec,omitempty"`         // kBytes write per second.
	ReadCountSec       float64 `json:"readCountSec,omitempty"`       // read count per second.
	WriteCountSec      float64 `json:"writeCountSec,omitempty"`      // write count per second.
	MergeReadCountSec  float64 `json:"mergeReadCountSec,omitempty"`  // merge read count per second.
	MergeWriteCountSec float64 `json:"mergeWriteCountSec,omitempty"` // merge write count per second
}

type MemoryUsage struct {
	mem.VirtualMemoryStat
}

type NetworkIO struct {
	Iface   string  `json:"iface"`             // interface name
	Rxpck   float64 `json:"rxpck"`             // rxpck/s
	Txpck   float64 `json:"txpck"`             // txpck/s
	RxkB    float64 `json:"rxkB"`              // rxkB/s
	TxkB    float64 `json:"txkB"`              // txkB/s
	Errin   float64 `json:"errin,omitempty"`   // err(in) number per second
	Errout  float64 `json:"errout,omitempty"`  // err(out) number per second
	Dropin  float64 `json:"dropin,omitempty"`  // drop package(in) per second
	Dropout float64 `json:"dropout,omitempty"` // drop package(out) per second
	Fifoin  float64 `json:"fifoin,omitempty"`  // FIFO buffers errors per second while receiving
	Fifoout float64 `json:"fifoout,omitempty"` // FIFO buffers errors per second while sending
}

type collectWorkloadFunc func(scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error)
