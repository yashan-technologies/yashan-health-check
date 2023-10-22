package osutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"yhc/defs/regexpdef"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/execer"
)

const (
	column_str = "NAME,KNAME,MODEL,UUID,SIZE,ROTA,FSTYPE,TYPE,MOUNTPOINT,PKNAME" // lsblk 输出的列名
)

type Device struct {
	Name       string `json:"name"`       // 设备名称
	KName      string `json:"kname"`      // 内核分配给该设备的名称
	Model      string `json:"model"`      // 设备型号
	Uuid       string `json:"uuid"`       // 该设备分区的 UUID
	Size       string `json:"size"`       // 设备大小
	Rota       string `json:"rota"`       // 设备是否为旋转介质，为 0表示不是旋转介质，为 1则表示是
	Fstype     string `json:"fstype"`     // 文件系统类型
	Type       string `json:"type"`       // 设备类型，主要有 disk、part、rom 等几种
	MountPoint string `json:"mountpoint"` // 设备挂载点
	Pkname     string `json:"pkname"`     // 父设备名称
}

func Lsblk(log yaslog.YasLog) ([]*Device, error) {
	args := []string{"-P", "-b", "-o", column_str}
	execer := execer.NewExecer(log, execer.WithPrintResult())
	ret, stdout, stderr := execer.Exec("lsblk", args...)
	if ret != 0 {
		err := fmt.Errorf("failed to exec lsblk,err: %s", stderr)
		log.Error(err)
		return nil, err
	}
	disks, err := parserLsblk(log, strings.NewReader(stdout))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return disks, nil
}

func parserLsblk(log yaslog.YasLog, r io.Reader) ([]*Device, error) {
	scan := bufio.NewScanner(r)
	devices := []*Device{}
	for scan.Scan() {
		device := &Device{}
		deviceInfo := make(map[string]string)
		raw := scan.Text()
		sr := regexpdef.LsblkOutputRegexp.FindAllStringSubmatch(raw, -1)
		for _, k := range sr {
			if len(k) < 2 { // 格式有误
				continue
			}
			k[1] = strings.TrimSpace(strings.ToLower(k[1]))
			k[2] = strings.TrimSpace(strings.ToLower(k[2]))
			deviceInfo[k[1]] = k[2]
		}
		data, err := json.Marshal(deviceInfo)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, device); err != nil {
			return nil, err
		}
		if regexpdef.LsblkIgnoreDeviceRegexp.MatchString(device.KName) {
			log.Debugf("ignore device: %s", device.KName)
			continue
		}
		devices = append(devices, device)
	}
	return devices, nil
}
