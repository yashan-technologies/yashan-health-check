package confdef

import (
	"path"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

var _metricConfig *YHCMetricConfig

const (
	M_HOST     = "host"
	M_DATABASE = "database"
	M_OBJECTS  = "objects"
	M_SAFETY   = "safety"
	M_CUSTOM   = "custom"
)

const (
	MT_INVALID MetricType = "invalid"
	MT_SQL     MetricType = "sql"
	MT_BASH    MetricType = "bash"
)

type YHCMetricConfig struct {
	Metrics []*YHCMetric `toml:"metrics"`
}

type YHCMetric struct {
	Name          string                    `toml:"name"`
	NameAlias     string                    `toml:"name_alias,omitempty"`
	ModuleName    string                    `toml:"module_name"`
	MetricType    MetricType                `toml:"metric_type"`
	Hidden        bool                      `toml:"hidden"`
	Default       bool                      `toml:"default"`
	Enabled       bool                      `toml:"enabled"`
	ColumnAlias   map[string]string         `toml:"column_alias,omitempty"`
	ColumnOrder   []string                  `toml:"column_order,omitempty"`
	ItemNames     map[string]string         `toml:"item_names,omitempty"`
	NumberColumns []string                  `toml:"number_columns,omitempty"`
	Labels        []string                  `toml:"labels,omitempty"`
	AlertRules    map[string][]AlertDetails `toml:"alert_rules,omitempty"`
	SQL           string                    `toml:"sql,omitempty"`     // SQL类型的指标的sql语句
	Command       string                    `toml:"command,omitempty"` // bash类型指标的bash命令
}

type AlertDetails struct {
	Expression  string `toml:"expression"`
	Description string `toml:"description,omitempty"`
	Suggestion  string `toml:"suggestion,omitempty"`
}

type MetricType string

const (
	AL_INVALID  = "invalid"
	AL_INFO     = "info"
	AL_WARING   = "warning"
	AL_CRITICAL = "critical"
)

func initMetricConf(paths []string) error {
	conf := YHCMetricConfig{}
	for _, p := range paths {
		if !path.IsAbs(p) {
			p = path.Join(runtimedef.GetYHCHome(), p)
		}
		c, err := loadMetricConf(p)
		if err != nil {
			return err
		}
		for _, metric := range c.Metrics {
			if len(metric.ModuleName) == 0 {
				metric.ModuleName = M_CUSTOM
			}
			conf.Metrics = append(conf.Metrics, metric)
		}
	}
	_metricConfig = &conf
	return nil
}

func loadMetricConf(p string) (config *YHCMetricConfig, err error) {
	config = new(YHCMetricConfig)
	if !fs.IsFileExist(p) {
		return config, &errdef.ErrFileNotFound{FName: p}
	}
	if _, err := toml.DecodeFile(p, config); err != nil {
		return config, &errdef.ErrFileParseFailed{FName: p, Err: err}
	}
	return config, nil
}

func GetMetricConf() *YHCMetricConfig {
	return _metricConfig
}
