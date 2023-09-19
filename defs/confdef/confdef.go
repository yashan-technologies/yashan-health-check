package confdef

import (
	"path"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

func InitYHCConf(yhcConf string) error {
	if !path.IsAbs(yhcConf) {
		yhcConf = path.Join(runtimedef.GetYHCHome(), yhcConf)
	}
	if !fs.IsFileExist(yhcConf) {
		return &errdef.ErrFileNotFound{FName: yhcConf}
	}
	if _, err := toml.DecodeFile(yhcConf, &_yhcConf); err != nil {
		return &errdef.ErrFileParseFailed{FName: yhcConf, Err: err}
	}
	return nil
}

func InitYHCMetricConf(yhcMetricsConf string) error {
	if !path.IsAbs(yhcMetricsConf) {
		yhcMetricsConf = path.Join(runtimedef.GetYHCHome(), yhcMetricsConf)
	}
	if !fs.IsFileExist(yhcMetricsConf) {
		return &errdef.ErrFileNotFound{FName: yhcMetricsConf}
	}
	if _, err := toml.DecodeFile(yhcMetricsConf, &_yhcMetricConfig); err != nil {
		return &errdef.ErrFileParseFailed{FName: yhcMetricsConf, Err: err}
	}
	return nil
}
