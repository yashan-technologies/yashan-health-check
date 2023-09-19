package confdef

var _yhcMetricConfig YHCMetricConfig

type YHCMetricConfig struct {
	Metrics []YHCMetric `toml:"metrics"`
}

type YHCMetric struct {
	Name          string                  `toml:"name"`
	NameAlias     string                  `toml:"name_alias,omitempty"`
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
