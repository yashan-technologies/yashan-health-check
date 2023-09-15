package define

import (
	"sync"
	"time"

	"yhc/commons/yasdb"
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
	Children map[MetricName]YHCItem `json:"children,omitempty"`
}

type YHCModule struct {
	Module string `json:"-"`
	mtx    sync.RWMutex
	items  map[MetricName]*YHCItem
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

func NewNoNeedCheckMetric(name string, err error, desc string) *NoNeedCheckMetric {
	return &NoNeedCheckMetric{
		Name:        name,
		Error:       err,
		Description: desc,
	}
}

func (c *YHCModule) Set(item *YHCItem) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.items == nil {
		c.items = make(map[MetricName]*YHCItem)
	}
	c.items[item.Name] = item
}

func (c *YHCModule) Items() map[MetricName]*YHCItem {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	items := make(map[MetricName]*YHCItem)
	for k, v := range c.items {
		tmp := *v
		items[k] = &tmp
	}
	return items
}
