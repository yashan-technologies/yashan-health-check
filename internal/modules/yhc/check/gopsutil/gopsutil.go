package gopsutil

import (
	"errors"
	"time"

	"yhc/internal/modules/yhc/check/define"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var _typeToFuncMap = map[define.WorkloadType]collectWorkloadFunc{
	define.WT_CPU:     collectCPUUasge,
	define.WT_DISK:    collectDiskIO,
	define.WT_MEMORY:  collectMemoryUsage,
	define.WT_NETWORK: collectNetworkIO,
}

func Collect(t define.WorkloadType, scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error) {
	collectFunc, ok := _typeToFuncMap[t]
	if !ok {
		return nil, errors.New("invalid workload type, could not found collect function")
	}
	return collectFunc(scrapeInterval, scrapeTimes)
}

func collectCPUUasge(scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error) {
	res := make(define.WorkloadOutput)
	for i := 0; i < scrapeTimes; i++ {
		m := make(define.WorkloadItem, 0)
		cpuInfos, err := cpu.Times(false)
		if err != nil {
			return res, err
		}
		for _, cpuInfo := range cpuInfos {
			m[cpuInfo.CPU] = CpuUsage{cpuInfo}
		}
		now := time.Now().Unix()
		res[now] = m
		time.Sleep(time.Second * time.Duration(scrapeInterval))
	}
	return res, nil
}

func collectNetworkIO(scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error) {
	res := make(define.WorkloadOutput)
	// time -> iface -> value
	datas := make(map[int64]map[string]net.IOCountersStat)
	// used to get value from res map
	timeArray := []int64{}
	// get data first
	for i := 0; i < scrapeTimes+1; i++ {
		m := make(map[string]net.IOCountersStat, 0)
		if ioCounters, err := net.IOCounters(true); err != nil {
			return res, err
		} else {
			for _, io := range ioCounters {
				m[io.Name] = io
			}
		}
		now := time.Now().Unix()
		timeArray = append(timeArray, now)
		datas[now] = m
		time.Sleep(time.Second * time.Duration(scrapeInterval))
	}
	// calculate last
	for i := 1; i < len(timeArray); i++ {
		newData, oldData := datas[timeArray[i]], datas[timeArray[i-1]]
		m := make(define.WorkloadItem, 0)
		for iface, newIO := range newData {
			if oldIO, ok := oldData[iface]; ok {
				netIO := NetworkIO{
					Iface:   iface,
					Rxpck:   float64((newIO.PacketsRecv - oldIO.PacketsRecv) / uint64(scrapeInterval)),
					Txpck:   float64((newIO.PacketsSent - oldIO.PacketsSent) / uint64(scrapeInterval)),
					RxkB:    float64((newIO.BytesRecv - oldIO.BytesRecv) / uint64(scrapeInterval) / 1024),
					TxkB:    float64((newIO.BytesSent - oldIO.BytesSent) / uint64(scrapeInterval) / 1024),
					Errin:   float64((newIO.Errin - oldIO.Errin) / uint64(scrapeInterval)),
					Errout:  float64((newIO.Errout - oldIO.Errout) / uint64(scrapeInterval)),
					Dropin:  float64((newIO.Dropin - oldIO.Dropin) / uint64(scrapeInterval)),
					Dropout: float64((newIO.Dropout - oldIO.Dropout) / uint64(scrapeInterval)),
					Fifoin:  float64((newIO.Fifoin - oldIO.Fifoin) / uint64(scrapeInterval)),
					Fifoout: float64((newIO.Fifoout - oldIO.Fifoout) / uint64(scrapeInterval)),
				}
				m[iface] = netIO
			}
		}
		res[timeArray[i]] = m
	}
	return res, nil
}

func collectDiskIO(scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error) {
	res := make(define.WorkloadOutput)
	datas := make(map[int64]map[string]disk.IOCountersStat)
	timeArray := []int64{}
	// get data
	for i := 0; i < scrapeTimes+1; i++ {
		m := make(map[string]disk.IOCountersStat, 0)
		partitions, err := disk.Partitions(false)
		if err != nil {
			return res, err
		}
		// collect disk names
		diskNames := []string{}
		for _, partition := range partitions {
			diskNames = append(diskNames, partition.Device)
		}
		ioCounters, err := disk.IOCounters(diskNames...)
		if err != nil {
			return res, err
		}
		for _, io := range ioCounters {
			m[io.Name] = io
		}
		now := time.Now().Unix()
		timeArray = append(timeArray, now)
		datas[now] = m
		time.Sleep(time.Second * time.Duration(scrapeInterval))
	}
	// calculate data
	for i := 1; i < len(timeArray); i++ {
		newData, oldData := datas[timeArray[i]], datas[timeArray[i-1]]
		m := make(define.WorkloadItem, 0)
		for name, newIO := range newData {
			if oldIO, ok := oldData[name]; ok {
				netIO := DiskIO{
					Dev:                name,
					Iops:               newIO.IopsInProgress,
					SerialNumber:       newIO.SerialNumber,
					Label:              newIO.Label,
					KBReadSec:          float64((newIO.ReadBytes - oldIO.ReadBytes) / uint64(scrapeInterval) / 1024),
					KBWriteSec:         float64((newIO.WriteBytes - oldIO.WriteBytes) / uint64(scrapeInterval) / 1024),
					ReadCountSec:       float64((newIO.ReadCount - oldIO.ReadCount) / uint64(scrapeInterval)),
					WriteCountSec:      float64((newIO.WriteCount - oldIO.WriteCount) / uint64(scrapeInterval)),
					MergeReadCountSec:  float64((newIO.MergedReadCount - oldIO.MergedWriteCount) / uint64(scrapeInterval)),
					MergeWriteCountSec: float64((newIO.MergedWriteCount - oldIO.MergedWriteCount) / uint64(scrapeInterval)),
				}
				m[name] = netIO
			}
		}
		res[timeArray[i]] = m
	}
	return res, nil
}

func collectMemoryUsage(scrapeInterval, scrapeTimes int) (define.WorkloadOutput, error) {
	res := make(define.WorkloadOutput)
	for i := 0; i < scrapeTimes; i++ {
		m := make(define.WorkloadItem, 0)
		if memInfo, err := mem.VirtualMemory(); err != nil {
			return res, err
		} else {
			m["mem"] = memInfo
		}
		now := time.Now().Unix()
		res[now] = m
		time.Sleep(time.Second * time.Duration(scrapeInterval))
	}
	return res, nil
}
