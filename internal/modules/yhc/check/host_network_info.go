package check

import (
	"strings"

	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yasutil/netcli"
	"github.com/shirou/gopsutil/net"
)

const (
	KEY_NETWORK_NAME          = "name"
	KEY_NETWORK_MTU           = "mtu"
	KEY_NETWORK_FLAGS         = "flags"
	KEY_NETWORK_IPV4          = "ipv4"
	KEY_NETWORK_IPV6          = "ipv6"
	KEY_NETWORK_HARDWARE_ADDR = "hardwareAddr"
)

func (c *YHCChecker) GetHostNetworkInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_NETWORK_INFO,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_NETWORK_INFO))
	netInfo, err := net.Interfaces()
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = c.dealNetworkData(netInfo)
	return
}

func (c *YHCChecker) dealNetworkData(networks []net.InterfaceStat) (res []map[string]any) {
	for _, network := range networks {
		var ipv4, ipv6 []string
		for _, addr := range network.Addrs {
			ip := addr.Addr
			if netcli.IsIPv6(ip) {
				ipv6 = append(ipv6, ip)
				continue
			}
			ipv4 = append(ipv4, ip)
		}
		res = append(res, map[string]any{
			KEY_NETWORK_NAME:          network.Name,
			KEY_NETWORK_MTU:           network.MTU,
			KEY_NETWORK_HARDWARE_ADDR: network.HardwareAddr,
			KEY_NETWORK_IPV4:          strings.Join(ipv4, stringutil.STR_NEWLINE),
			KEY_NETWORK_IPV6:          strings.Join(ipv6, stringutil.STR_NEWLINE),
			KEY_NETWORK_FLAGS:         strings.Join(network.Flags, stringutil.STR_COMMA),
		})
	}
	return
}
