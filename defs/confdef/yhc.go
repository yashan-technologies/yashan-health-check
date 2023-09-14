package confdef

import (
	"time"

	"yhc/utils/timeutil"
)

var _yhcConf YHC

type YHC struct {
	LogLevel    string `toml:"log_level"`
	Range       string `toml:"range"`
	Output      string `toml:"output"`
	MaxDuration string `toml:"max_duration"`
	MinDuration string `toml:"min_duration"`
}

func GetYHCConf() YHC {
	return _yhcConf
}

func (c YHC) GetMaxDuration() (time.Duration, error) {
	if len(c.MaxDuration) == 0 {
		return time.Hour * 24, nil
	}
	maxDuration, err := timeutil.GetDuration(c.MaxDuration)
	if err != nil {
		return 0, err
	}
	return maxDuration, err
}

func (c YHC) GetMinDuration() (time.Duration, error) {
	if len(c.MinDuration) == 0 {
		return time.Minute * 1, nil
	}
	minDuration, err := timeutil.GetDuration(c.MinDuration)
	if err != nil {
		return 0, err
	}
	return minDuration, err
}

func (c YHC) GetMinAndMaxDuration() (min time.Duration, max time.Duration, err error) {
	min, err = c.GetMinDuration()
	if err != nil {
		return
	}
	max, err = c.GetMaxDuration()
	if err != nil {
		return
	}
	return
}

func (c YHC) GetRange() (r time.Duration) {
	r, err := timeutil.GetDuration(c.Range)
	if err != nil {
		return time.Hour * 24
	}
	return
}