package sar

import "yhc/internal/modules/yhc/check/define"

const (
	PREFIX_LINUX   = "Linux"
	PREFIX_AVERAGE = "Average"
)

var _envs = []string{"LANG=en_US.UTF-8", "LC_TIME=en_US.UTF-8"}

var WorkloadTypeToSarArgMap = map[define.WorkloadType]string{
	define.WT_CPU:     "-u",
	define.WT_DISK:    "-d",
	define.WT_MEMORY:  "-r",
	define.WT_NETWORK: "-n DEV",
}

type SarParseFunc func(m define.WorkloadItem, values []string) define.WorkloadItem

type SarCheckTitleFunc func(line string) bool

type CPUUsage struct {
	CPU    string  `json:"cpu"`    // cpu name, it will be 'all'
	User   float64 `json:"user"`   // percentage of CPU time spent in user space
	Nice   float64 `json:"nice"`   // percentage of CPU time spent on low priority tasks (niceness)
	System float64 `json:"system"` // percentage of CPU time spent in kernel space (system calls and kernel threads)
	IOWait float64 `json:"iowait"` // percentage of CPU time waiting for I/O operations to complete
	Steal  float64 `json:"steal"`  // percentage of CPU time stolen by other virtual machines or hypervisor in a virtualized environment
	Idle   float64 `json:"idle"`   // percentage of CPU time idle or not utilized
}

type NetworkIO struct {
	Iface  string  `json:"iface"`  // interface name
	Rxpck  float64 `json:"rxpck"`  // rxpck/s
	Txpck  float64 `json:"txpck"`  // txpck/s
	RxkB   float64 `json:"rxkB"`   // rxkB/s
	TxkB   float64 `json:"txkB"`   // txkB/s
	Rxcmp  float64 `json:"rxcmp"`  // rxcmp/s
	Txcmp  float64 `json:"excmp"`  // txcmp/s
	Rxmcst float64 `json:"rxmcst"` // rxmcst/s
	Ifutil float64 `json:"ifutil"` // %ifutil, ubuntu or kylin
}

type DiskIO struct {
	Dev     string  `json:"dev"`     // dev name.
	Tps     float64 `json:"tps"`     // tps, number of transfers per second
	RdSec   float64 `json:"rdSec"`   // rd_sec/s, number of sectors read from the device per second. centos
	WrSec   float64 `json:"wrSec"`   // wr_sec/s, number of sectors written to the device per second. centos
	AvgrqSz float64 `json:"avgrqSz"` // avgrq-sz, average size (in sectors) of the requests sent to the device.
	AvgquSz float64 `json:"avgquSz"` // avgqu-sz, average queue length (number of requests) waiting for service. get by sar
	Await   float64 `json:"await"`   // await, average time (in milliseconds) that a request waits in the queue and service time.
	Svctm   float64 `json:"svctm"`   // svctm, average service time (in milliseconds) for a request. centos or ubuntu
	Util    float64 `json:"util"`    // %util, percentage of time that the device was busy with I/O operations.
	RKBSec  float64 `json:"rkBSec"`  // rkB/s, ubuntu or kylin
	WKBSec  float64 `json:"wkBSec"`  // wkB/s, ubuntu or kylin
	DKBSec  float64 `json:"dkBSec"`  // dkB/s, kylin
}

type MemoryUsage struct {
	KBMemFree   int64   `json:"kBMemFree"`   // kbmemfree
	KBAvail     int64   `json:"kBAvail"`     // kbavail, centos or lylin
	KBmemUsed   int64   `json:"kBMemUsed"`   // kbmemused
	MemUsed     float64 `json:"memUsed"`     // %memused
	KBBuffers   int64   `json:"kBBuffers"`   // kbbuffers
	KBCached    int64   `json:"kBCached"`    // kbcached
	KBCommit    int64   `json:"kBCommit"`    // kbcommit
	Commit      float64 `json:"commit"`      // %commit
	KBActive    int64   `json:"kBActive"`    // kbactive
	KBInact     int64   `json:"kBInact"`     // kbinact
	KBDirty     int64   `json:"kBDirty"`     // kbdirty
	RealMemUsed float64 `json:"realMemUsed"` // real mem used percent
}
