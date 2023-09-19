package define

import "fmt"

var _DefaultModuleAlias = map[ModuleName]string{
	// level 1 modules
	MODULE_OVERVIEW: "概述",
	MODULE_HOST:     "主机检查",
	MODULE_YASDB:    "数据库检查",
	MODULE_OBJECT:   "对象检查",
	MODULE_SECURITY: "安全检查",
	MODULE_LOG:      "日志分析",
	MODULE_CUSTOM:   "自定义检查",

	// level 2 modules
	MODULE_OVERVIEW_HOST:  "主机概览",
	MODULE_OVERVIEW_YASDB: "数据库概览",

	MODULE_HOST_WORKLOAD: "主机负载检查",

	MODULE_YASDB_STANDBY:     "主备检查",
	MODULE_YASDB_CONFIG:      "数据库配置检查",
	MODULE_YASDB_TABLESPACE:  "表空间检查",
	MODULE_YASDB_CONTROLFILE: "控制文件检查",
	MODULE_YASDB_BACKUP:      "备份检查",
	MODULE_YASDB_WORKLOAD:    "负载检查",
	MODULE_YASDB_PERFORMANCE: "性能分析",

	MODULE_OBJECT_NUMBER:     "对象数量统计",
	MODULE_OBJECT_STATUS:     "对象状态检查",
	MODULE_OBJECT_TABLE:      "表",
	MODULE_OBJECT_CONSTRAINT: "约束",
	MODULE_OBJECT_INDEX:      "索引",
	MODULE_OBJECT_SEQUENCE:   "序列",
	MODULE_OBJECT_TASK:       "任务",
	MODULE_OBJECT_PACKAGE:    "包",

	MODULE_SECURITY_LOGIN:      "登录检查",
	MODULE_SECURITY_PERMISSION: "用户与权限检查",
	MODULE_SECURITY_AUDIT:      "审计日志分析",

	MODULE_LOG_RUN:   "数据库变更",
	MODULE_LOG_REDO:  "redo日志分析",
	MODULE_LOG_UNDO:  "undo日志分析",
	MODULE_LOG_ERROR: "错误日志分析",
}

var _DefaultMetricAlias = map[MetricName]string{
	METRIC_YASDB_INSTANCE:           "实例信息",
	METRIC_YASDB_DATABASE:           "数据库信息",
	METRIC_YASDB_FILE_PERMISSION:    "数据库文件权限",
	METRIC_YASDB_LISTEN_ADDR:        "数据库IP及端口",
	METRIC_HOST_INFO:                "主机信息",
	METRIC_HOST_CPU_INFO:            "CPU信息",
	METRIC_HOST_CPU_USAGE:           "CPU使用率",
	METRIC_YASDB_REPLICATION_STATUS: "数据库主备连接状态",
	METRIC_YASDB_PARAMETER:          "数据库参数检查",
	METRIC_YASDB_TABLESPACE:         "表空间",
	METRIC_YASDB_CONTROLFILE:        "控制文件",
	METRIC_YASDB_SESSION:            "会话数检查",
	METRIC_YASDB_WAIT_EVENT:         "等待事件",
	METRIC_YASDB_OBJECT_COUNT:       "对象总数",
	METRIC_YASDB_OBJECT_OWNER:       "各owner对象统计",
	METRIC_YASDB_OBJECT_TABLESPACE:  "各tablespace对象统计",
	METRIC_YASDB_INDEX_BLEVEL:       "超过三层的索引",
	METRIC_YASDB_INDEX_COLUMN:       "字段过多的索引",
	METRIC_YASDB_INDEX_INVISIBLE:    "不可见索引",
	METRIC_YASDB_REDO_LOG:           "redo日志分析",
	METRIC_YASDB_RUN_LOG_ERROR:      "run.log错误分析",
}

var ModuleTree = map[ModuleName]map[ModuleName]map[MetricName]struct{}{
	MODULE_OVERVIEW: {
		MODULE_OVERVIEW_HOST: {},
		MODULE_OVERVIEW_YASDB: {
			METRIC_YASDB_INSTANCE:        {},
			METRIC_YASDB_DATABASE:        {},
			METRIC_YASDB_FILE_PERMISSION: {},
			METRIC_YASDB_LISTEN_ADDR:     {},
		},
	},
	MODULE_HOST: {
		MODULE_HOST_WORKLOAD: {
			METRIC_HOST_INFO:     {},
			METRIC_HOST_CPU_INFO: {},
		},
	},
	MODULE_YASDB: {
		MODULE_YASDB_STANDBY: {
			METRIC_YASDB_REPLICATION_STATUS: {},
		},
		MODULE_YASDB_CONFIG: {
			METRIC_YASDB_PARAMETER: {},
		},
		MODULE_YASDB_TABLESPACE: {
			METRIC_YASDB_TABLESPACE: {},
		},
		MODULE_YASDB_CONTROLFILE: {
			METRIC_YASDB_CONTROLFILE: {},
		},
		MODULE_YASDB_BACKUP: {},
		MODULE_YASDB_WORKLOAD: {
			METRIC_YASDB_SESSION: {},
		},
		MODULE_YASDB_PERFORMANCE: {},
	},
	MODULE_OBJECT: {
		MODULE_OBJECT_NUMBER: {
			METRIC_YASDB_OBJECT_COUNT:      {},
			METRIC_YASDB_OBJECT_OWNER:      {},
			METRIC_YASDB_OBJECT_TABLESPACE: {},
		},
		MODULE_OBJECT_STATUS:     {},
		MODULE_OBJECT_TABLE:      {},
		MODULE_OBJECT_CONSTRAINT: {},
		MODULE_OBJECT_INDEX: {
			METRIC_YASDB_INDEX_BLEVEL:    {},
			METRIC_YASDB_INDEX_COLUMN:    {},
			METRIC_YASDB_INDEX_INVISIBLE: {},
		},
		MODULE_OBJECT_SEQUENCE: {},
		MODULE_OBJECT_TASK:     {},
		MODULE_OBJECT_PACKAGE:  {},
	},
	MODULE_SECURITY: {
		MODULE_SECURITY_LOGIN:      {},
		MODULE_SECURITY_PERMISSION: {},
		MODULE_SECURITY_AUDIT:      {},
	},
	MODULE_LOG: {
		MODULE_LOG_RUN: {},
		MODULE_LOG_REDO: {
			METRIC_YASDB_REDO_LOG: {},
		},
		MODULE_LOG_UNDO: {},
		MODULE_LOG_ERROR: {
			METRIC_YASDB_RUN_LOG_ERROR: {},
		},
	},
	MODULE_CUSTOM: {},
}

const (
	// level 1 modules
	MODULE_OVERVIEW ModuleName = "overview"
	MODULE_HOST     ModuleName = "host-check"
	MODULE_YASDB    ModuleName = "yasdb-check"
	MODULE_OBJECT   ModuleName = "object-check"
	MODULE_SECURITY ModuleName = "security-check"
	MODULE_LOG      ModuleName = "log-analysis"
	MODULE_CUSTOM   ModuleName = "custom-check"

	// the followings are level 2 modules
	// parent module: MN_OVERVIEW
	MODULE_OVERVIEW_HOST  ModuleName = "overview-host"
	MODULE_OVERVIEW_YASDB ModuleName = "overview-yasdb"

	// parent module: MN_HOST
	MODULE_HOST_WORKLOAD ModuleName = "host-workload-check"

	// parent module: MN_YASDB
	MODULE_YASDB_STANDBY     ModuleName = "yasdb-standby-check"
	MODULE_YASDB_CONFIG      ModuleName = "yasdb-config-check"
	MODULE_YASDB_TABLESPACE  ModuleName = "yasdb-tablespace-check"
	MODULE_YASDB_CONTROLFILE ModuleName = "yasdb-controlfile-check"
	MODULE_YASDB_BACKUP      ModuleName = "yasdb-backup-check"
	MODULE_YASDB_WORKLOAD    ModuleName = "yasdb-workload-check"
	MODULE_YASDB_PERFORMANCE ModuleName = "yasdb-performance-analysis"

	// parent module: MN_OBJECT
	MODULE_OBJECT_NUMBER     ModuleName = "object-number-count"
	MODULE_OBJECT_STATUS     ModuleName = "object-status-check"
	MODULE_OBJECT_TABLE      ModuleName = "object-table-check"
	MODULE_OBJECT_CONSTRAINT ModuleName = "object-constraint-check"
	MODULE_OBJECT_INDEX      ModuleName = "object-index-check"
	MODULE_OBJECT_SEQUENCE   ModuleName = "object-sequence-check"
	MODULE_OBJECT_TASK       ModuleName = "object-task-check"
	MODULE_OBJECT_PACKAGE    ModuleName = "object-package-check"

	MODULE_SECURITY_LOGIN      ModuleName = "security-login-config"
	MODULE_SECURITY_PERMISSION ModuleName = "security-permission-check"
	MODULE_SECURITY_AUDIT      ModuleName = "security-audit-check"

	// parent module: MN_LOG
	MODULE_LOG_RUN   ModuleName = "log-run-analysis"
	MODULE_LOG_REDO  ModuleName = "log-redo-analysis"
	MODULE_LOG_UNDO  ModuleName = "log-undo-analysis"
	MODULE_LOG_ERROR ModuleName = "log-error-analysis"
)

type ModuleName string

const (
	METRIC_YASDB_INSTANCE           MetricName = "yasdb-instance"
	METRIC_YASDB_DATABASE           MetricName = "yasdb-database"
	METRIC_YASDB_FILE_PERMISSION    MetricName = "yasdb-file-permission"
	METRIC_YASDB_LISTEN_ADDR        MetricName = "yasdb-listen-address"
	METRIC_HOST_INFO                MetricName = "host-info"
	METRIC_HOST_CPU_INFO            MetricName = "host-cpu-info"
	METRIC_HOST_CPU_USAGE           MetricName = "host-cpu-usage"
	METRIC_YASDB_REPLICATION_STATUS MetricName = "yasdb-replication-status"
	METRIC_YASDB_PARAMETER          MetricName = "yasdb-parameter"
	METRIC_YASDB_TABLESPACE         MetricName = "yasdb-tablespace"
	METRIC_YASDB_CONTROLFILE        MetricName = "yasdb-controlfile"
	METRIC_YASDB_SESSION            MetricName = "yasdb-session"
	METRIC_YASDB_WAIT_EVENT         MetricName = "yasdb-wait-event"
	METRIC_YASDB_OBJECT_COUNT       MetricName = "yasdb-object-count"
	METRIC_YASDB_OBJECT_OWNER       MetricName = "yasdb-object-owner"
	METRIC_YASDB_OBJECT_TABLESPACE  MetricName = "yasdb-object-tablespace"
	METRIC_YASDB_INDEX_BLEVEL       MetricName = "yasdb-index-blevel"
	METRIC_YASDB_INDEX_COLUMN       MetricName = "yasdb-index-column"
	METRIC_YASDB_INDEX_INVISIBLE    MetricName = "yasdb-index-invisible"
	METRIC_YASDB_REDO_LOG           MetricName = "yasdb-redo-log"
	METRIC_YASDB_RUN_LOG_ERROR      MetricName = "yasdb-run-log-error"
)

type MetricName string

type Metric struct {
	MetricName MetricName
	SubMetrics map[MetricName]*Metric
}

func GetMetricDefaultAlias(name MetricName) (string, error) {
	alias, ok := _DefaultMetricAlias[name]
	if !ok {
		return "", fmt.Errorf("failed to find default Metirc %s", name)
	}
	return alias, nil
}

func GetModuleDefaultAlias(name ModuleName) (string, error) {
	alias, ok := _DefaultModuleAlias[name]
	if !ok {
		return "", fmt.Errorf("failed to find default Module %s", name)
	}
	return alias, nil
}
