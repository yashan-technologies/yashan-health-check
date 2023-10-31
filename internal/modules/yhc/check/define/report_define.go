package define

const (
	ET_PRE         ElementType = "pre" // 可换行的文本
	ET_DIV         ElementType = "div"
	ET_ALERT       ElementType = "custom-alert"
	ET_EMPTY       ElementType = "a-empty"
	ET_CODE        ElementType = "custom-code"
	ET_TABLE       ElementType = "custom-table"
	ET_CHART       ElementType = "custom-chart"
	ET_DESCRIPTION ElementType = "custom-description"
	ET_H3          ElementType = "h3"
)

const (
	CT_BAR  ChartType = "bar"
	CT_PIE  ChartType = "pie"
	CT_LINE ChartType = "line"
)

const (
	AT_SUCCESS  AlertType = "success"
	AT_INFO     AlertType = "info"
	AT_WARNING  AlertType = "warning"
	AT_CRITICAL AlertType = "critical"
	AT_ERROR    AlertType = "error"
)

var AlertTypeAliasMap = map[AlertType]string{
	AT_SUCCESS:  "成功",
	AT_INFO:     "提示",
	AT_WARNING:  "警告",
	AT_CRITICAL: "严重",
	AT_ERROR:    "错误",
}

type ElementType string

type ChartType string

type AlertType string

type PandoraReport struct {
	ReportTitle    string         `json:"reportTitle,omitempty"`
	ReportSubTitle string         `json:"reportSubTitle,omitempty"`
	FileControl    string         `json:"fileControl,omitempty"`
	Author         string         `json:"author,omitempty"`
	Version        string         `json:"version,omitempty"`
	Time           string         `json:"time,omitempty"`
	CostTime       int            `json:"costTime,omitempty"`
	ChangeLog      string         `json:"changeLog,omitempty"`
	ReportData     []*PandoraMenu `json:"reportData,omitempty"`
}

type PandoraMenu struct {
	IsMenu       bool              `json:"isMenu,omitempty"`
	IsChapter    bool              `json:"isChapter,omitempty"`
	Title        string            `json:"title,omitempty"`
	TitleEn      string            `json:"-"`
	WarningCount int               `json:"warningCount,omitempty"`
	Children     []*PandoraMenu    `json:"children,omitempty"`
	MenuIndex    int               `json:"menuIndex"`
	Elements     []*PandoraElement `json:"elements,omitempty"`
}

type PandoraElement struct {
	MetricName   string      `json:"metricName,omitempty"`
	ElementTitle string      `json:"elementTitle,omitempty"`
	ElementType  ElementType `json:"element,omitempty"`
	InnerText    string      `json:"innerText,omitempty"`
	Attributes   interface{} `json:"attributes,omitempty"`
	Solts        interface{} `json:"solts,omitempty"`
	SoltName     string      `json:"soltName,omitempty"`
	Config       interface{} `json:"config,omitempty"`
	Extend       interface{} `json:"extend,omitempty"`
}

type CustomOptionTitle struct {
	Text    string `json:"text,omitempty"`
	SubText string `json:"subText,omitempty"`
}

type AlertAttributes struct {
	AlertType   AlertType `json:"type,omitempty"`
	Message     string    `json:"message,omitempty"`
	Description string    `json:"description,omitempty"`
}

type TableAttributes struct {
	Title        string                   `json:"title,omitempty"`
	DataSource   []map[string]interface{} `json:"dataSource"`
	TableColumns []*TableColumn           `json:"columns,omitempty"`
}

type TableColumn struct {
	Title     string `json:"title,omitempty"`
	DataIndex string `json:"dataIndex,omitempty"`
}

type CodeAttributes struct {
	Title    string `json:"title,omitempty"`
	Language string `json:"language,omitempty"`
	Code     string `json:"code,omitempty"`
}

type PAttributes struct {
	Title CustomOptionTitle `json:"title,omitempty"`
}

type ChartCoordinate struct {
	X interface{} `json:"x,omitempty"`
	Y interface{} `json:"y,omitempty"`
}

type ChartData struct {
	Name  string             `json:"name,omitempty"`
	Value []*ChartCoordinate `json:"value,omitempty"`
}

type ChartCustomOptions struct {
	Title     CustomOptionTitle `json:"title,omitempty"`
	ChartType ChartType         `json:"type,omitempty"`
	Data      []*ChartData      `json:"data,omitempty"`
}

type ChartAttributes struct {
	CustomOptions ChartCustomOptions `json:"customOptions,omitempty"`
}

type ChartSeries struct {
	Name  string        `json:"name,omitempty"`
	Value []interface{} `json:"value,omitempty"`
}

type DescriptionAttributes struct {
	Title string             `json:"title,omitempty"`
	Data  []*DescriptionData `json:"data,omitempty"`
}

type DescriptionData struct {
	Label string      `json:"label,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
