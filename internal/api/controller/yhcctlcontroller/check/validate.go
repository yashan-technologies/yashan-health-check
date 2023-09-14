package checkcontroller

import (
	"os"
	"path"
	"time"

	"yhc/defs/confdef"
	"yhc/defs/errdef"
	"yhc/defs/regexpdef"
	"yhc/defs/runtimedef"
	"yhc/log"
	"yhc/utils/fileutil"
	"yhc/utils/jsonutil"
	"yhc/utils/stringutil"
	"yhc/utils/timeutil"
	"yhc/utils/userutil"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	f_type   = "type"
	f_range  = "range"
	f_start  = "start"
	f_end    = "end"
	f_output = "output"

	range_help = "you must ensure that the number before (M|d|h|m) is greater than 0"
)

var (
	_examplesTime = []string{
		"yyyy-MM-dd",
		"yyyy-MM-dd-hh",
		"yyyy-MM-dd-hh-mm",
	}

	_examplesRange = []string{
		"1M",
		"1d",
		"1h",
		"1m",
	}
)

func (c *CheckCmd) validate() error {

	if err := c.validateRange(); err != nil {
		return err
	}
	if err := c.validateStartAndEnd(); err != nil {
		return err
	}
	if err := c.validateOutput(); err != nil {
		return err
	}
	return nil
}

func (c *CheckCmd) validateRange() error {
	conf := confdef.GetYHCConf()
	log.Controller.Debugf("conf: %s\n", jsonutil.ToJSONString(conf))
	log.Controller.Debugf("cmd: %s", jsonutil.ToJSONString(c))
	if stringutil.IsEmpty(c.Range) {
		return nil
	}
	if !regexpdef.RangeRegex.MatchString(c.Range) {
		return errdef.NewErrYHCFlag(f_range, c.Range, _examplesRange, range_help)
	}
	minDuration, maxDuration, err := conf.GetMinAndMaxDuration()
	if err != nil {
		log.Controller.Errorf("get duration err: %s", err.Error())
		return err
	}
	log.Controller.Debugf("get min %s max %s", minDuration.String(), maxDuration.String())
	r, err := timeutil.GetDuration(c.Range)
	if err != nil {
		return err
	}
	if r > maxDuration {
		return errdef.NewGreaterMaxDur(conf.MaxDuration)
	}
	if r < minDuration {
		return errdef.NewLessMinDur(conf.MinDuration)
	}
	return nil
}

func (c *CheckCmd) validateStartAndEnd() error {
	conf := confdef.GetYHCConf()
	var (
		startNotEmpty, endNotEmpty bool
		start, end                 time.Time
		err                        error
	)
	if !stringutil.IsEmpty(c.Start) {
		if !regexpdef.TimeRegex.MatchString(c.Start) {
			return errdef.NewErrYHCFlag(f_start, c.Start, _examplesTime, "")
		}
		start, err = timeutil.GetTimeDivBySepa(c.Start, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		now := time.Now()
		if start.After(now) {
			return errdef.ErrStartShouldLessCurr
		}
		startNotEmpty = true
	}
	if !stringutil.IsEmpty(c.End) {
		if !regexpdef.TimeRegex.MatchString(c.End) {
			return errdef.NewErrYHCFlag(f_end, c.End, _examplesTime, "")
		}
		end, err = timeutil.GetTimeDivBySepa(c.End, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		endNotEmpty = true
	}
	if startNotEmpty && endNotEmpty {
		minDuration, maxDuration, err := conf.GetMinAndMaxDuration()
		if err != nil {
			log.Controller.Errorf("get duration err: %s", err.Error())
			return err
		}
		if end.Before(start) {
			return errdef.ErrEndLessStart
		}
		r := end.Sub(start)
		if r > maxDuration {
			return errdef.NewGreaterMaxDur(conf.MaxDuration)
		}
		if r < minDuration {
			return errdef.NewLessMinDur(conf.MaxDuration)
		}
	}
	return nil
}

func (c *CheckCmd) validateOutput() error {
	output := c.Output
	if !regexpdef.PathRegex.Match([]byte(output)) {
		return errdef.ErrPathFormat
	}
	if !path.IsAbs(output) {
		output = path.Join(runtimedef.GetYHCHome(), output)
	}
	_, err := os.Stat(output)
	if err != nil {
		if os.IsPermission(err) {
			return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
		}
		if !os.IsNotExist(err) {
			return err
		}
		if err := fs.Mkdir(output); err != nil {
			log.Controller.Errorf("create output err: %s", err.Error())
			if os.IsPermission(err) {
				return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
			}
			return err
		}
	}
	return fileutil.CheckUserWrite(output)

}

func (c *CheckCmd) fillDefault() {
	if stringutil.IsEmpty(c.Output) {
		c.Output = confdef.GetYHCConf().Output
	}
	if !path.IsAbs(c.Output) {
		c.Output = path.Join(runtimedef.GetYHCHome(), c.Output)
	}
	c.Output = path.Clean(c.Output)
}
