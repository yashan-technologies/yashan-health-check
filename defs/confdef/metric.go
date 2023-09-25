package confdef

import (
	"path"

	"yhc/defs/errdef"
	"yhc/defs/runtimedef"

	"git.yasdb.com/go/yasutil/fs"
	"github.com/BurntSushi/toml"
)

const (
	M_HOST     = "host"
	M_DATABASE = "database"
	M_OBJECTS  = "objects"
	M_SAFETY   = "safety"
	M_CUSTOM   = "custom"
)

type YHCMetricConfig struct {
	Metrics  []*YHCMetric `toml:"metrics"`
	Includes []string     `toml:"includes"`
}

type YHCMetric struct {
	Name          string                  `toml:"name"`
	NameAlias     string                  `toml:"name_alias,omitempty"`
	ModuleName    string                  `toml:"module_name"`
	MetricType    MetricType              `toml:"metric_type"`
	Hidden        bool                    `toml:"hidden"`
	Default       bool                    `toml:"default"`
	Enabled       bool                    `toml:"enabled"`
	ColumnAlias   map[string]string       `toml:"column_alias,omitempty"`
	ItemNames     map[string]string       `toml:"item_names,omitempty"`
	NumberColumns []string                `toml:"number_columns,omitempty"`
	Labels        []string                `toml:"labels,omitempty"`
	AlertRules    map[string]AlertDetails `toml:"alert_rules,omitempty"`
	SQL           string                  `toml:"sql,omitempty"`     // SQL类型的指标的sql语句
	Command       string                  `toml:"command,omitempty"` // bash类型指标的bash命令
}

type AlertDetails struct {
	Expression  string `toml:"expression"`
	Description string `toml:"description,omitempty"`
	Suggestion  string `toml:"suggestion,omitempty"`
}

const (
	MT_INVALID MetricType = "invalid"
	MT_SQL     MetricType = "sql"
	MT_BASH    MetricType = "bash"
	MT_UNION   MetricType = "union"
)

type MetricType string

const (
	OT_INVALID OutputType = "invalid"
	OT_TEXT    OutputType = "text"
	OT_MAP     OutputType = "map"
	OT_TABLE   OutputType = "table"
	OT_GRAPH   OutputType = "graph"
)

type OutputType string

const (
	AL_INVALID  = "invalid"
	AL_INFO     = "info"
	AL_WARING   = "waring"
	AL_CRITICAL = "critical"
)

type AlertLevel string

var _metricConfig *YHCMetricConfig

func initMetricConf(p string) error {
	if !path.IsAbs(p) {
		p = path.Join(runtimedef.GetYHCHome(), p)
	}
	def, err := loadMetricConf(p)
	if err != nil {
		return err
	}
	defer func() {
		_metricConfig = def
	}()
	if len(def.Includes) == 0 {
		return nil
	}
	parent := path.Dir(p)
	for _, include := range def.Includes {
		custom, err := loadCustomMetricConf(parent, include)
		if err != nil {
			return err
		}
		for _, metric := range custom.Metrics {
			if len(metric.ModuleName) == 0 {
				metric.ModuleName = M_CUSTOM
			}
			def.Metrics = append(def.Metrics, metric)
		}
	}
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

func loadCustomMetricConf(parent string, include string) (customMetric *YHCMetricConfig, err error) {
	target := include
	if !path.IsAbs(include) {
		target = path.Join(parent, include)
	}
	customMetric, err = loadMetricConf(target)
	return
}

func GetMetricConf() *YHCMetricConfig {
	return _metricConfig
}