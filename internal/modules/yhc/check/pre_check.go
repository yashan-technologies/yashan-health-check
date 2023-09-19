package check

import (
	"fmt"

	"yhc/defs/bashdef"
	"yhc/log"
	"yhc/utils/execerutil"
)

func (c *YHCChecker) CheckSarAccess() error {
	cmd := []string{
		"-c",
		bashdef.CMD_SAR,
		"-V",
	}
	exe := execerutil.NewExecer(log.Module)
	ret, _, stderr := exe.Exec(bashdef.CMD_BASH, cmd...)
	if ret != 0 {
		return fmt.Errorf("failed to check sar command, err: %s", stderr)
	}
	return nil
}
