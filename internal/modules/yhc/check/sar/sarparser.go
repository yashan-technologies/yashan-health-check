package sar

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"yhc/commons/constants"
	"yhc/defs/regexpdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/fs"
)

const (
	memoryUsageKey = "memory"
)

// dd:mm:ss AM CPU     %user     %nice   %system   %iowait    %steal     %idle
const (
	_base_cpu_date_index OutputIndex = iota
	_base_cpu_period_index
	_base_cpu_cpu_index
	_base_cpu_user_index
	_base_cpu_nice_index
	_base_cpu_system_index
	_base_cpu_iowait_index
	_base_cpu_steal_index
	_base_cpu_idle_index
	_base_cpu_length
)

// dd:mm:ss AM IFACE   rxpck/s   txpck/s    rxkB/s    txkB/s   rxcmp/s   txcmp/s  rxmcst/s
const (
	_base_network_date_index OutputIndex = iota
	_base_network_period_index
	_base_network_iface_index
	_base_network_rxpck_index
	_base_network_txpck_index
	_base_network_rxkb_index
	_base_network_txkb_index
	_base_network_rxcmp_index
	_base_network_txcmp_index
	_base_network_rxmcst_index
	_base_network_length
)

// dd:mm:ss AM kbmemfree kbmemused  %memused kbbuffers  kbcached  kbcommit   %commit  kbactive   kbinact   kbdirty
const (
	_base_memory_date_index OutputIndex = iota
	_base_memory_period_index
	_base_memory_kbmemfree_index
	_base_memory_kbmemused_index
	_base_memory_memused_index
	_base_memory_kbbuffers_index
	_base_memory_kbcached_index
	_base_memory_kbcommit_index
	_base_memory_commit_index
	_base_memory_kbactive_index
	_base_memory_kbinact_index
	_base_memory_kbdirty_index
	_base_memory_length
)

// dd:mm:ss AM DEV       tps  rd_sec/s  wr_sec/s  avgrq-sz  avgqu-sz     await     svctm     %util
const (
	_base_disk_date_index OutputIndex = iota
	_base_disk_period_index
	_base_disk_dev_index
	_base_disk_tps_index
	_base_disk_rdsec_index
	_base_disk_wrsec_index
	_base_disk_avgrqsz_index
	_base_disk_avgqusz_index
	_base_disk_await_index
	_base_disk_svctm_index
	_base_disk_util_index
	_base_disk_length
)

type OutputIndex int

type SarParser interface {
	GetParserFunc(t define.WorkloadType) (SarParseFunc, SarCheckTitleFunc)
	ParseCpu(m define.WorkloadItem, values []string) define.WorkloadItem
	IsCpuTitle(line string) bool
	ParseNetwork(m define.WorkloadItem, values []string) define.WorkloadItem
	IsNetworkTitle(line string) bool
	ParseMemory(m define.WorkloadItem, values []string) define.WorkloadItem
	IsMemoryTitle(line string) bool
	ParseDisk(m define.WorkloadItem, values []string) define.WorkloadItem
	IsDiskTitle(line string) bool
	GetSarDir() string
}

// use centos as base parser
type baseParser struct {
	log yaslog.YasLog
}

func NewBaseParser(yaslog yaslog.YasLog) *baseParser {
	return &baseParser{
		log: yaslog,
	}
}

// [Interface Func]
func (b *baseParser) GetParserFunc(t define.WorkloadType) (SarParseFunc, SarCheckTitleFunc) {
	switch t {
	case define.WT_DISK:
		return b.ParseDisk, b.IsDiskTitle
	case define.WT_MEMORY:
		return b.ParseMemory, b.IsMemoryTitle
	case define.WT_NETWORK:
		return b.ParseNetwork, b.IsNetworkTitle
	case define.WT_CPU:
		return b.ParseCpu, b.IsCpuTitle
	default:
		return nil, nil
	}
}

// [Interface Func]

func (b *baseParser) ParseCpu(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -u
	if len(values) < int(_base_cpu_length) {
		// not enough data, skip
		b.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	cpuUsage := CPUUsage{}
	var err error
	cpuUsage.CPU = values[_base_cpu_cpu_index]
	if cpuUsage.User, err = strconv.ParseFloat(values[_base_cpu_user_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if cpuUsage.Nice, err = strconv.ParseFloat(values[_base_cpu_nice_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if cpuUsage.System, err = strconv.ParseFloat(values[_base_cpu_system_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if cpuUsage.IOWait, err = strconv.ParseFloat(values[_base_cpu_iowait_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if cpuUsage.Steal, err = strconv.ParseFloat(values[_base_cpu_steal_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if cpuUsage.Idle, err = strconv.ParseFloat(values[_base_cpu_idle_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	m[cpuUsage.CPU] = cpuUsage
	return m
}

// [Interface Func]
func (b *baseParser) IsCpuTitle(line string) bool {
	return strings.Contains(line, "CPU")
}

// [Interface Func]
func (b *baseParser) ParseNetwork(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -n DEV
	if len(values) < int(_base_network_length) {
		b.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	networkIO := NetworkIO{}
	var err error
	networkIO.Iface = values[_base_network_iface_index] // fill data
	if networkIO.Rxpck, err = strconv.ParseFloat(values[_base_network_rxpck_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.Txpck, err = strconv.ParseFloat(values[_base_network_txpck_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.RxkB, err = strconv.ParseFloat(values[_base_network_rxkb_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.TxkB, err = strconv.ParseFloat(values[_base_network_txkb_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.Rxcmp, err = strconv.ParseFloat(values[_base_network_rxcmp_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.Txcmp, err = strconv.ParseFloat(values[_base_network_txcmp_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if networkIO.Rxmcst, err = strconv.ParseFloat(values[_base_network_rxmcst_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	m[networkIO.Iface] = networkIO
	return m
}

// [Interface Func]
func (b *baseParser) IsNetworkTitle(line string) bool {
	return strings.Contains(line, "IFACE")
}

// [Interface Func]
func (b *baseParser) ParseMemory(m define.WorkloadItem, values []string) define.WorkloadItem {
	// commadn: sar -r
	if len(values) < int(_base_memory_length) {
		b.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	memoryUsage := MemoryUsage{}
	var err error
	if memoryUsage.KBMemFree, err = strconv.ParseInt(values[_base_memory_kbmemfree_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBmemUsed, err = strconv.ParseInt(values[_base_memory_kbmemused_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.MemUsed, err = strconv.ParseFloat(values[_base_memory_memused_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBBuffers, err = strconv.ParseInt(values[_base_memory_kbbuffers_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBCached, err = strconv.ParseInt(values[_base_memory_kbcached_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBCommit, err = strconv.ParseInt(values[_base_memory_kbcommit_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.Commit, err = strconv.ParseFloat(values[_base_memory_commit_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBActive, err = strconv.ParseInt(values[_base_memory_kbactive_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBInact, err = strconv.ParseInt(values[_base_memory_kbinact_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if memoryUsage.KBDirty, err = strconv.ParseInt(values[_base_memory_kbdirty_index], constants.BASE_DECIMAL, constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	m[memoryUsageKey] = b.calculateRealMemUsed(memoryUsage)
	return m
}

func (b *baseParser) calculateRealMemUsed(m MemoryUsage) MemoryUsage {
	if m.KBMemFree+m.KBmemUsed != 0 {
		m.RealMemUsed = 100 * (1 - float64(m.KBMemFree+m.KBBuffers+m.KBCached)/float64(m.KBMemFree+m.KBmemUsed))
	}
	return m
}

// [Interface Func]
func (b *baseParser) IsMemoryTitle(line string) bool {
	return strings.Contains(line, "kbmemfree")
}

// [Interface Func]
func (b *baseParser) ParseDisk(m define.WorkloadItem, values []string) define.WorkloadItem {
	// command: sar -d
	if len(values) < int(_base_disk_length) {
		b.log.Warnf("not enough data, skip line: %s", strings.Join(values, stringutil.STR_BLANK_SPACE))
		return m
	}
	diskIO := DiskIO{}
	var err error
	diskIO.Dev = values[_base_disk_dev_index]
	if diskIO.Tps, err = strconv.ParseFloat(values[_base_disk_tps_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.RdSec, err = strconv.ParseFloat(values[_base_disk_rdsec_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.WrSec, err = strconv.ParseFloat(values[_base_disk_wrsec_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.AvgrqSz, err = strconv.ParseFloat(values[_base_disk_avgrqsz_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.AvgquSz, err = strconv.ParseFloat(values[_base_disk_avgqusz_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.Await, err = strconv.ParseFloat(values[_base_disk_await_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.Svctm, err = strconv.ParseFloat(values[_base_disk_svctm_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	if diskIO.Util, err = strconv.ParseFloat(values[_base_disk_util_index], constants.BIT_SIZE_64); err != nil {
		b.log.Error(err)
	}
	m[diskIO.Dev] = diskIO
	return m
}

// [Interface Func]
func (b *baseParser) IsDiskTitle(line string) bool {
	return strings.Contains(line, "DEV") && strings.Contains(line, "tps")
}

// [Interface Func]
func (b *baseParser) GetSarDir() string {
	defaultConfigPath := "/etc/sysconfig/sysstat"
	defaultSarDir := "/var/log/sa"
	currentFilePath := b.getSarDirFromConfig(defaultConfigPath)
	if stringutil.IsEmpty(currentFilePath) {
		return defaultSarDir
	}
	return currentFilePath
}

func (b *baseParser) getSarDirFromConfig(configPath string) string {
	if !fs.IsFileExist(configPath) {
		return ""
	}
	configMap := make(map[string]string)
	file, err := os.Open(configPath)
	if err != nil {
		b.log.Error(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, stringutil.STR_HASH) {
			// ignore line start with '#'
			continue
		}
		// key=value
		re := regexpdef.KeyValueRegexp
		match := re.FindStringSubmatch(line)
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			configMap[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		b.log.Error(err)
		return ""
	}
	return configMap["SAR_DIR"]
}
