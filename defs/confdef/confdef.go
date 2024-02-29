package confdef

import (
	"path"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

func InitYHCConf(yhcConf string) error {
	if err := initYHCConf(yhcConf); err != nil {
		return err
	}
	if err := initMetricConf(_yhcConf.MetricPaths); err != nil {
		return err
	}
	if err := initModuleConf(_yhcConf.DefaultModulePath); err != nil {
		return err
	}
	if err := initEvaluateModel(_yhcConf.EvaluateModelPath); err != nil {
		return err
	}
	if err := initNodesConfig(_yhcConf.NodesConfigPath); err != nil {
		return err
	}
	return nil
}

func initYHCConf(yhcConf string) error {
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
