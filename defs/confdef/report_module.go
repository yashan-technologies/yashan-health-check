package confdef

import (
	"path"
	"strings"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

var _moduleConfig *YHCModuleConfig

type YHCModuleConfig struct {
	Modules          []*YHCModuleNode `toml:"modules"`
	moduleAliasMap   map[string]string
	metricModulesMap map[string][]string
	metricOrder      []string
}

type YHCModuleNode struct {
	Name        string           `toml:"name"`
	NameAlias   string           `toml:"name_alias"`
	Children    []*YHCModuleNode `toml:"children"`
	MetricNames []string         `toml:"metric_names"`
}

func initModuleConf(p string) error {
	if !path.IsAbs(p) {
		p = path.Join(runtimedef.GetYHCHome(), p)
	}
	conf, err := loadModuleConf(p)
	if err != nil {
		return err
	}
	conf.moduleAliasMap = genModuleAliasMap(conf.Modules)
	conf.metricModulesMap = genMetricModulesMap(conf.Modules)
	conf.metricOrder = genMetricOrder(conf.Modules)
	_moduleConfig = conf
	return nil
}

func genModuleAliasMap(modules []*YHCModuleNode) map[string]string {
	res := make(map[string]string)

	var fn func(res map[string]string, node *YHCModuleNode)
	fn = func(res map[string]string, node *YHCModuleNode) {
		if node == nil {
			return
		}
		if stringutil.IsEmpty(node.NameAlias) {
			node.NameAlias = node.Name
		}
		res[node.Name] = node.NameAlias
		for _, child := range node.Children {
			fn(res, child)
		}
	}

	for _, module := range modules {
		fn(res, module)
	}
	return res
}

func genMetricModulesMap(nodes []*YHCModuleNode) map[string][]string {
	var fn func(node *YHCModuleNode, path []string, index map[string][]string)
	fn = func(node *YHCModuleNode, path []string, index map[string][]string) {
		path = append(path, node.Name)
		for _, metricName := range node.MetricNames {
			index[metricName] = append([]string{}, path...)
		}
		for _, child := range node.Children {
			fn(child, path, index)
		}
	}

	index := make(map[string][]string)
	for _, node := range nodes {
		fn(node, []string{}, index)
	}
	return index
}

func genMetricOrder(modules []*YHCModuleNode) []string {
	var result []string

	var fn func(node *YHCModuleNode, result *[]string)
	fn = func(node *YHCModuleNode, result *[]string) {
		*result = append(*result, node.MetricNames...)
		for _, child := range node.Children {
			fn(child, result)
		}
	}

	for _, module := range modules {
		fn(module, &result)
	}
	return result
}

func loadModuleConf(p string) (*YHCModuleConfig, error) {
	conf := &YHCModuleConfig{}
	if !fs.IsFileExist(p) {
		return conf, &errdef.ErrFileNotFound{FName: p}
	}
	if _, err := toml.DecodeFile(p, conf); err != nil {
		return conf, &errdef.ErrFileParseFailed{FName: p, Err: err}
	}
	return conf, nil
}

func GetModuleConf() *YHCModuleConfig {
	return _moduleConfig
}

func GetModuleAliasMap() map[string]string {
	return _moduleConfig.moduleAliasMap
}

func GetModuleAlias(name string) string {
	name = strings.TrimSpace(name)
	return _moduleConfig.moduleAliasMap[name]
}

func GetMetricModules(metricName string) []string {
	return _moduleConfig.metricModulesMap[metricName]
}

func GetMetricOrder() []string {
	return _moduleConfig.metricOrder
}
