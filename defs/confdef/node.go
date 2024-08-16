package confdef

import (
	"path"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

var _nodesConfig *NodesConfig

type NodeConfig struct {
	ListenAddr string `toml:"listen_addr"`
	User       string `toml:"user,omitempty"`
	Password   string `toml:"password,omitempty"`
}

type NodesConfig struct {
	Nodes []NodeConfig `toml:"nodes"`
}

func GetNodesConfig() *NodesConfig {
	return _nodesConfig
}

func initNodesConfig(p string) error {
	conf := &NodesConfig{}
	if !path.IsAbs(p) {
		p = path.Join(runtimedef.GetYHCHome(), p)
	}
	if !fs.IsFileExist(p) {
		return &errdef.ErrFileNotFound{FName: p}
	}
	if _, err := toml.DecodeFile(p, conf); err != nil {
		return &errdef.ErrFileParseFailed{FName: p, Err: err}
	}
	_nodesConfig = conf
	return nil
}
