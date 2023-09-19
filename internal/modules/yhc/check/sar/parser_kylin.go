package sar

import (
	"strconv"
	"strings"

	"yhc/internal/modules/yhc/check/define"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
)

// dd:mm:ss AM IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s   %ifutil
const (
	_kylin_network_date_index OutputIndex = iota
	_kylin_network_period_index
	_kylin_network_iface_index
	_kylin_network_rxpck_index
	_kylin_network_txpck_index
	_kylin_network_rxkb_index
	_kylin_network_txkb_index
	_kylin_network_rxcmp_index
	_kylin_network_txcmp_index
	_kylin_network_rxmcst_index
	_kylin_network_ifutil_index
	_kylin_network_length
)

// dd:mm:ss AM  DEV       tps     rkB/s     wkB/s     dkB/s   areq-sz    aqu-sz     await     %util
const (
	_kylin_disk_date_index OutputIndex = iota
	_kylin_disk_period_index
	_kylin_disk_dev_index
	_kylin_disk_tps_index
	_kylin_disk_rkb_index
	_kylin_disk_wkb_index
	_kylin_disk_dkb_index
	_kylin_disk_areqsz_index
	_kylin_disk_aqusz_index
	_kylin_disk_await_index
	_kylin_disk_util_index
	_kylin_disk_length
)

// dd:mm:ss AM kbmemfree   kbavail kbmemused  %memused kbbuffers  kbcached  kbcommit   %commit  kbactive   kbinact   kbdirty
const (
	_kylin_memory_date_index OutputIndex = iota
	_kylin_memory_period_index
	_kylin_memory_kbmemfree_index
	_kylin_memory_kbavail_index
	_kylin_memory_kbmemused_index
	_kylin_memory_memused_index
	_kylin_memory_kbbuffers_index
	_kylin_memory_kbcached_index
	_kylin_memory_kbcommit_index
	_kylin_memory_commit_index
	_kylin_memory_kbactive_index
	_kylin_memory_kbinact_index
	_kylin_memory_kbdirty_index
	_kylin_memory_length
)

type KylinParser struct {
	base *baseParser
}

func NewKylinParser(yaslog yaslog.YasLog) *KylinParser {
	return &KylinParser{
		base: NewBaseParser(yaslog),
	}
}

// [Interface Func]
func (k *KylinParser) GetParserFunc(t define.WorkloadType) (SarParseFunc, SarCheckTitleFunc) {
	switch t {
	case define.WT_DISK:
		return k.ParseDisk, k.IsDiskTitle
	case define.WT_MEMORY:
		return k.ParseMemory, k.IsMemoryTitle
	case define.WT_NETWORK:
		return k.ParseNetwork, k.IsNetworkTitle
	case define.WT_CPU:
		return k.ParseCpu, k.IsCpuTitle
	default:
		return nil, nil
	}
}

// [Interface Func]
func (k *KylinParser) ParseCpu(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -u
	return k.base.ParseCpu(m, values)
}

// [Interface Func]
func (k *KylinParser) IsCpuTitle(line string) bool {
	return k.base.IsCpuTitle(line)
}

// [Interface Func]
func (k *KylinParser) ParseNetwork(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -n DEV
	if len(values) < int(_kylin_network_length) {
		k.base.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	networkIO := NetworkIO{}
	var err error
	networkIO.Iface = values[_kylin_network_iface_index]
	if networkIO.Rxpck, err = strconv.ParseFloat(values[_kylin_network_rxpck_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.Txpck, err = strconv.ParseFloat(values[_kylin_network_txpck_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.RxkB, err = strconv.ParseFloat(values[_kylin_network_rxkb_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.TxkB, err = strconv.ParseFloat(values[_kylin_network_txkb_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.Rxcmp, err = strconv.ParseFloat(values[_kylin_network_rxcmp_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.Txcmp, err = strconv.ParseFloat(values[_kylin_network_txcmp_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.Rxmcst, err = strconv.ParseFloat(values[_kylin_network_rxmcst_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if networkIO.Ifutil, err = strconv.ParseFloat(values[_kylin_network_ifutil_index], 64); err != nil {
		k.base.log.Error(err)
	}
	m[networkIO.Iface] = networkIO
	return m
}

// [Interface Func]
func (k *KylinParser) IsNetworkTitle(line string) bool {
	return k.base.IsNetworkTitle(line)
}

// [Interface Func]
func (k *KylinParser) ParseDisk(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -d
	if len(values) < int(_kylin_disk_length) {
		k.base.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	diskIO := DiskIO{}
	var err error
	diskIO.Dev = values[_kylin_disk_dev_index]
	if diskIO.Tps, err = strconv.ParseFloat(values[_kylin_disk_tps_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.RKBSec, err = strconv.ParseFloat(values[_kylin_disk_rkb_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.WKBSec, err = strconv.ParseFloat(values[_kylin_disk_wkb_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.DKBSec, err = strconv.ParseFloat(values[_kylin_disk_dkb_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.AvgrqSz, err = strconv.ParseFloat(values[_kylin_disk_areqsz_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.AvgquSz, err = strconv.ParseFloat(values[_kylin_disk_aqusz_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.Await, err = strconv.ParseFloat(values[_kylin_disk_await_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if diskIO.Util, err = strconv.ParseFloat(values[_kylin_disk_util_index], 64); err != nil {
		k.base.log.Error(err)
	}
	m[diskIO.Dev] = diskIO
	return m
}

// [Interface Func]
func (k *KylinParser) IsDiskTitle(line string) bool {
	return k.base.IsDiskTitle(line)
}

// [Interface Func]
func (k *KylinParser) ParseMemory(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -r
	if len(values) < int(_kylin_memory_length) {
		k.base.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	memoryUsage := MemoryUsage{}
	var err error
	if memoryUsage.KBMemFree, err = strconv.ParseInt(values[_kylin_memory_kbmemfree_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBAvail, err = strconv.ParseInt(values[_kylin_memory_kbavail_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBmemUsed, err = strconv.ParseInt(values[_kylin_memory_kbmemused_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.MemUsed, err = strconv.ParseFloat(values[_kylin_memory_memused_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBBuffers, err = strconv.ParseInt(values[_kylin_memory_kbbuffers_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBCached, err = strconv.ParseInt(values[_kylin_memory_kbcached_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBCommit, err = strconv.ParseInt(values[_kylin_memory_kbcommit_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.Commit, err = strconv.ParseFloat(values[_kylin_memory_commit_index], 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBActive, err = strconv.ParseInt(values[_kylin_memory_kbactive_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBInact, err = strconv.ParseInt(values[_kylin_memory_kbinact_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	if memoryUsage.KBDirty, err = strconv.ParseInt(values[_kylin_memory_kbdirty_index], 10, 64); err != nil {
		k.base.log.Error(err)
	}
	memoryUsage.RealMemUsed = memoryUsage.MemUsed // no need to calculate
	m[memoryUsageKey] = memoryUsage
	return m
}

// [Interface Func]
func (k *KylinParser) IsMemoryTitle(line string) bool {
	return k.base.IsMemoryTitle(line)
}

// [Interface Func]
func (k *KylinParser) GetSarDir() string {
	return k.base.GetSarDir()
}
