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
	Name           string                    `toml:"name"`
	NameAlias      string                    `toml:"name_alias,omitempty"`
	ModuleName     string                    `toml:"module_name"`
	MetricType     MetricType                `toml:"metric_type"`
	Hidden         bool                      `toml:"hidden"`
	Default        bool                      `toml:"default"`
	Enabled        bool                      `toml:"enabled"`
	ColumnAlias    map[string]string         `toml:"column_alias,omitempty"`
	ColumnOrder    []string                  `toml:"column_order,omitempty"`
	HiddenColumns  []string                  `toml:"hidden_columns,omitempty"`  // hide column in table, only used in alert expression
	ByteColumns    []string                  `toml:"byte_columns,omitempty"`    // convert byte columns to human readable size
	PercentColumns []string                  `toml:"percent_columns,omitempty"` // convert percent columns to number + '%'
	ItemNames      map[string]string         `toml:"item_names,omitempty"`
	NumberColumns  []string                  `toml:"number_columns,omitempty"`
	Labels         []string                  `toml:"labels,omitempty"`
	AlertRules     map[string][]AlertDetails `toml:"alert_rules,omitempty"`
	SQL            string                    `toml:"sql,omitempty"`     // SQL类型的指标的sql语句
	Command        string                    `toml:"command,omitempty"` // bash类型指标的bash命令
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
	AL_WARNING  = "warning"
	AL_CRITICAL = "critical"
)

var AlertLevelMap = map[string]string{
	AL_INVALID:  "无效",
	AL_INFO:     "提示",
	AL_WARNING:  "警告",
	AL_CRITICAL: "严重",
}

func InitMetricConf(paths []string) error {
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
