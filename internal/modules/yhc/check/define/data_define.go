package define

import (
	"time"

	"yhc/commons/yasdb"
	"yhc/defs/confdef"
)

const (
	DATATYPE_SAR      DataType = "sar"
	DATATYPE_GOPSUTIL DataType = "gopstuil"
)

const (
	WT_CPU     WorkloadType = "cpu"
	WT_NETWORK WorkloadType = "network"
	WT_MEMORY  WorkloadType = "memory"
	WT_DISK    WorkloadType = "disk"
)

type WorkloadType string

type WorkloadItem map[string]interface{}

type WorkloadOutput map[int64]WorkloadItem

type DataType string

type HostWorkResponse struct {
	Data     map[string]interface{}
	Errors   map[string]string
	DataType DataType
}

type YHCItem struct {
	Name     MetricName             `json:"-"` // 检查项名称
	Error    string                 `json:"error,omitempty"`
	Details  interface{}            `json:"details,omitempty"`  // 每个检查项包含的数据
	DataType DataType               `json:"datatype,omitempty"` // 数据类型，在Details可能使用多种数据时使用
	Alerts   map[string][]*YHCAlert `json:"alerts,omitempty"`
}

type YHCAlert struct {
	Level  string            `json:"level"`
	Value  any               `json:"value"`
	Labels map[string]string `json:"labels"`
	confdef.AlertDetails
}

type NoNeedCheckMetric struct {
	Name        string
	Description string
	Error       error
}

type CheckerBase struct {
	DBInfo *yasdb.YashanDB
	Start  time.Time
	End    time.Time
	Output string
	// TODO: add other struct which checker needed
}
