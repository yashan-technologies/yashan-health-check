package check

import (
	"errors"
	"fmt"
	"strings"

	"yhc/defs/bashdef"
	"yhc/defs/runtimedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"
	"yhc/utils/osutil"
	"yhc/utils/userutil"
)

const (
	_firewalld_inactive      = "inactive"
	_firewalld_active        = "active"
	_ubuntu_firewalld_active = "Status: active"

	KEY_FIREWALLD_STATUS = "firewalldStatus"
)

func (c *YHCChecker) GetHostFirewalldStatus(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_FIREWALLD,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_FIREWALLD))
	osRelease := runtimedef.GetOSRelease()
	execer := execerutil.NewExecer(log)
	// ubuntu
	if osRelease.Id == osutil.UBUNTU_ID {
		if !userutil.IsCurrentUserRoot() {
			err = errors.New("checking ubuntu firewall status requires sudo or root")
			data.Error = err.Error()
			return
		}
		_, stdout, _ := execer.EnvExec(_envs, bashdef.CMD_BASH, "-c", fmt.Sprintf("%s status", bashdef.CMD_UFW))
		data.Details = map[string]any{
			KEY_FIREWALLD_STATUS: strings.Contains(stdout, _ubuntu_firewalld_active),
		}
		return
	}
	// other os
	_, stdout, _ := execer.EnvExec(_envs, bashdef.CMD_BASH, "-c", fmt.Sprintf("%s is-active firewalld", bashdef.CMD_SYSTEMCTL))
	data.Details = map[string]any{
		KEY_FIREWALLD_STATUS: strings.Contains(stdout, _firewalld_active) && !strings.Contains(stdout, _firewalld_inactive),
	}
	return
}

func (c *YHCChecker) GetHostIPTables(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_IPTABLES,
	}
	defer c.fillResult(data)

	log := log.Module.M(string(define.METRIC_HOST_IPTABLES))
	execer := execerutil.NewExecer(log)
	ret, stdout, stderr := execer.EnvExec(_envs, bashdef.CMD_IPTABLES, "-L")
	if ret != 0 {
		err = fmt.Errorf("failed to get iptables, err: %v", stderr)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	data.Details = strings.TrimSpace(stdout)
	return
}
