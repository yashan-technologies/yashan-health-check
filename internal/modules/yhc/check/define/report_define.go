package define

const (
	ET_PRE         ElementType = "pre" // 可换行的文本
	ET_DIV         ElementType = "div"
	ET_ALERT       ElementType = "custom-alert" // TODO: 支持custom-alert后更换成custom-alert
	ET_EMPTY       ElementType = "a-empty"
	ET_CODE        ElementType = "custom-code"
	ET_TABLE       ElementType = "custom-table"
	ET_CHART       ElementType = "custom-chart"
	ET_DESCRIPTION ElementType = "custom-description"
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
	IsMenu    bool              `json:"isMenu,omitempty"`
	IsChapter bool              `json:"isChapter,omitempty"`
	Title     string            `json:"title,omitempty"`
	Children  []*PandoraMenu    `json:"children,omitempty"`
	Elements  []*PandoraElement `json:"elements,omitempty"`
}

type PandoraElement struct {
	ElementType ElementType `json:"element,omitempty"`
	InnerText   string      `json:"innerText,omitempty"`
	Attributes  interface{} `json:"attributes,omitempty"`
	Solts       interface{} `json:"solts,omitempty"`
	SoltName    string      `json:"soltName,omitempty"`
	Config      interface{} `json:"config,omitempty"`
	Extend      interface{} `json:"extend,omitempty"`
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
	Title        CustomOptionTitle        `json:"title,omitempty"`
	DataSource   []map[string]interface{} `json:"dataSource,omitempty"`
	TableColumns []*TableColumn           `json:"columns,omitempty"`
}

type TableColumn struct {
	Title     string `json:"title,omitempty"`
	DataIndex string `json:"dataIndex,omitempty"`
}

type CodeAttributes struct {
	Title    CustomOptionTitle `json:"title,omitempty"`
	Language string            `json:"language,omitempty"`
	Code     string            `json:"code,omitempty"`
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
