package osutil

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"yhc/defs/regexpdef"
	"yhc/utils/stringutil"
)

const (
	OS_RELEASE_PATH = "/etc/os-release"
)

const (
	UBUNTU_ID = "ubuntu"
	CENTOS_ID = "centos"
	KYLIN_ID  = "kylin"
)

type OSRelease struct {
	Name       string `json:"NAME"`        // 操作系统的名称，例如 "Ubuntu"、"CentOS" 等
	Version    string `json:"VERSION"`     // 操作系统的版本号，可能包括次版本号和修订版本号
	Id         string `json:"ID"`          // 操作系统的唯一标识符，通常是小写字母，例如 "ubuntu"、"centos" 等
	PrettyName string `json:"PRETTY_NAME"` // 格式化的操作系统名称和版本号，用于人类可读性。
	VersionId  string `json:"VERSION_ID"`  // 操作系统版本的标识符，通常是数值
}

func GetOSRelease() (*OSRelease, error) {
	m := make(map[string]string)
	file, err := os.Open(OS_RELEASE_PATH)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = regexpdef.SpaceRegexp.ReplaceAllString(line, stringutil.STR_EMPTY)
		if strings.HasPrefix(line, stringutil.STR_HASH) {
			continue
		}
		fields := strings.Split(line, stringutil.STR_EQUAL_SIGN)
		if len(fields) != 2 {
			continue
		}
		m[fields[0]] = strings.Trim(fields[1], stringutil.STR_SINGEL_QUOTE)
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	osRelease := new(OSRelease)
	if err := json.Unmarshal(bytes, osRelease); err != nil {
		return nil, err
	}
	return osRelease, nil
}
