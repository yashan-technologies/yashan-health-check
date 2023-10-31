package check

import (
	"strconv"

	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/mathutil"
	"yhc/utils/osutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yasutil/size"
	"github.com/shirou/gopsutil/disk"
)

const (
	KEY_DISK_TATAL        = "total"
	KEY_DISK_FREE         = "free"
	KEY_DISK_USED         = "used"
	KEY_DISK_INODES_TOTAL = "inodesTotal"
	KEY_DISK_INODES_FREE  = "inodesFree"
	KEY_DISK_INODES_USED  = "inodesUsed"

	KEY_DISK_BLOCK_SIZE = "size"
)

type DiskUsage struct {
	Device       string `json:"device"`
	MountOptions string `json:"mountOptions"`
	disk.UsageStat
}

func (c *YHCChecker) GetHostDiskInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_DISK_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_DISK_INFO))
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Errorf("failed to get host disk info, err: %s", err.Error())
		data.Error = err.Error()
		return
	}
	details := []map[string]any{}
	for _, partition := range partitions {
		var usageStat *disk.UsageStat
		usageStat, err = disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Errorf("failed to get disk usage info, err: %s", err.Error())
			data.Error = err.Error()
			return
		}
		usageStat.InodesUsedPercent = mathutil.Round(usageStat.InodesUsedPercent, decimal)
		usageStat.UsedPercent = mathutil.Round(usageStat.UsedPercent, decimal)
		usage := DiskUsage{
			Device:       partition.Device,
			MountOptions: partition.Opts,
			UsageStat:    *usageStat,
		}
		var detail map[string]any
		detail, err = c.convertObjectData(usage)
		if err != nil {
			log.Errorf("failed to covert disk info, err: %v", err)
			continue
		}
		details = append(details, c.formatHostDiskInfo(detail))
	}
	data.Details = details
	return
}

func (c *YHCChecker) formatHostDiskInfo(detail map[string]any) map[string]any {
	bytesArr := []string{KEY_DISK_FREE, KEY_DISK_USED, KEY_DISK_TATAL}
	for _, field := range bytesArr {
		value, ok := detail[field]
		if !ok {
			continue
		}
		data, ok := value.(float64)
		if !ok {
			continue
		}
		detail[field] = size.GenHumanReadableSize(data, decimal)
	}
	numArr := []string{KEY_DISK_INODES_FREE, KEY_DISK_INODES_TOTAL, KEY_DISK_INODES_USED}
	for _, field := range numArr {
		value, ok := detail[field]
		if !ok {
			continue
		}
		data, ok := value.(float64)
		if !ok {
			continue
		}
		detail[field] = mathutil.GenHumanReadableNumber(data, decimal)
	}
	return detail
}

func (c *YHCChecker) GetHostDiskBlockInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_DISK_BLOCK_INFO,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_DISK_BLOCK_INFO))
	devices, err := osutil.Lsblk(log)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return err
	}
	details := []map[string]any{}
	for _, device := range devices {
		detail, err := c.convertObjectData(device)
		if err != nil {
			log.Errorf("failed to convert device data, err: %v", err)
			continue
		}
		details = append(details, c.formatHostDiskBlockInfo(detail))
	}
	data.Details = details
	return
}

func (c *YHCChecker) formatHostDiskBlockInfo(detail map[string]any) map[string]any {
	strArr := []string{KEY_DISK_BLOCK_SIZE}
	for _, field := range strArr {
		value, ok := detail[field]
		if !ok {
			continue
		}
		data, ok := value.(string)
		if !ok {
			continue
		}
		num, err := strconv.ParseFloat(data, 64)
		if err != nil {
			continue
		}
		detail[field] = size.GenHumanReadableSize(num, decimal)
	}
	return detail
}
