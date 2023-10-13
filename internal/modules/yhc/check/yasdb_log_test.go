package check_test

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"

	"yhc/commons/constants"
	"yhc/defs/timedef"

	"github.com/shirou/gopsutil/host"
)

func TestMatch(t *testing.T) {
	// logStr := "[    0.000000] MTRR default type: write-back"
	logStr := "[1453203.223876] IPv6: ADDRCONF(NETDEV_UP): veth0dfbd09: link is not ready"
	// logStr1 := "[    0.359095] pid_max: default: 32768 minimum: 301"

	// 定义匹配正则表达式的模式
	pattern := `^\[\s*(\d+\.\d+)\]`
	regex := regexp.MustCompile(pattern)
	info, _ := host.Info()

	// 使用正则表达式匹配时间戳
	matches := regex.FindStringSubmatch(logStr)
	if len(matches) > 1 {
		timestamp := matches[1]
		log.Printf("时间戳：%s\n", timestamp)
		secondFromBoot, err := strconv.ParseFloat(matches[1], constants.BIT_SIZE_64)
		if err != nil {
			err = fmt.Errorf("dmesg log time: %s format err: %s, skip", matches[1], err.Error())
			log.Fatal(err)
			return
		}
		time := time.Unix(int64(info.BootTime+uint64(secondFromBoot)), 0)
		log.Println(time.Format(timedef.TIME_FORMAT))
	} else {
		log.Println("无法提取时间戳")
	}
}
