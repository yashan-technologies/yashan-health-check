package check

import (
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/osutil"

	"git.yasdb.com/go/yaserr"
	"github.com/shirou/gopsutil/disk"
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
		details = append(details, detail)
	}
	data.Details = details
	return
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
		details = append(details, detail)
	}
	data.Details = details
	return
}
