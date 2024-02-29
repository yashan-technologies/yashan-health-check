package check

import (
	"errors"
	"fmt"
	"strings"

	"yhc/defs/bashdef"
	"yhc/defs/regexpdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/execerutil"

	"git.yasdb.com/go/yaserr"
)

const (
	KEY_BIOS_INFORMATION = "BIOS Information"
)

func (c *YHCChecker) GetHostBIOSInfo(name string) (err error) {
	data := &define.YHCItem{
		Name: define.METRIC_HOST_BIOS_INFO,
	}
	defer c.fillResults(data)

	log := log.Module.M(string(define.METRIC_HOST_BIOS_INFO))
	execer := execerutil.NewExecer(log)
	ret, stdout, stderr := execer.EnvExec(_envs, bashdef.CMD_DMIDECODE)
	if ret != 0 {
		err = fmt.Errorf("failed to get host bios info, err: %s", stderr)
		log.Error(err)
		data.Error = err.Error()
		return
	}
	detail, err := c.parseDmidecodeOutput(stdout)
	if err != nil {
		err = yaserr.Wrap(err)
		log.Error(err)
		data.Error = err.Error()
		return err
	}
	data.Details = strings.TrimSpace(detail)
	return
}

func (c *YHCChecker) parseDmidecodeOutput(input string) (output string, err error) {
	blocks := regexpdef.MultiLineRegexp.Split(input, -1)
	for _, block := range blocks {
		if strings.Contains(block, KEY_BIOS_INFORMATION) {
			output = block
			return
		}
	}
	err = errors.New("failed to get BIOS Information")
	return
}
