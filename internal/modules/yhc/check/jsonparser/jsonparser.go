package jsonparser

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"yhc/defs/compiledef"
	"yhc/defs/confdef"
	"yhc/defs/timedef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/utils/stringutil"

	"git.yasdb.com/go/yaslog"
)

const (
	_REPORT_TITLE = "YashanDB 深度巡检报告"
	_FILE_CONTROL = "此文档仅供崖山科技有限公司与最终用户审阅，不得向与此无关的个人或机构传阅或复制。"
	_AUTHOR       = "Yashan Health Check"
	_CHANGE_LOG   = "生成巡检报告"
)

var _hostOverviewRelated = []define.MetricName{
	define.METRIC_HOST_INFO,
	define.METRIC_HOST_CPU_INFO,
}

var _yasdbOverviewRelated = []define.MetricName{
	define.METRIC_YASDB_DATABASE,
	define.METRIC_YASDB_INSTANCE,
	define.METRIC_YASDB_LISTEN_ADDR,
}

var _mergeMetricMap = map[define.MetricName][]define.MetricName{
	define.METRIC_HOST_INFO:      _hostOverviewRelated,
	define.METRIC_YASDB_DATABASE: _yasdbOverviewRelated,
}

type MetricParseFunc func(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error

type JsonParser struct {
	log     yaslog.YasLog
	base    define.CheckerBase
	metrics []*confdef.YHCMetric
	results map[define.MetricName]*define.YHCItem
}

func NewJsonParser(log yaslog.YasLog, base define.CheckerBase, metrics []*confdef.YHCMetric, results map[define.MetricName]*define.YHCItem) *JsonParser {
	parser := &JsonParser{
		log:     log,
		metrics: metrics,
		results: results,
		base:    base,
	}
	return parser
}

func (j *JsonParser) Parse() *define.PandoraReport {
	report := &define.PandoraReport{
		ReportTitle: _REPORT_TITLE,
		FileControl: _FILE_CONTROL,
		Author:      _AUTHOR,
		ChangeLog:   _CHANGE_LOG,
		Time:        j.base.Start.Format(timedef.TIME_FORMAT),
		CostTime:    int(j.base.End.Sub(j.base.Start).Seconds()),
		Version:     compiledef.GetAPPVersion(),
	}
	j.MergeMetrics()
	for _, module := range confdef.GetModuleConf().Modules {
		menu := &define.PandoraMenu{IsMenu: true, Title: confdef.GetModuleAlias(module.Name)}
		report.ReportData = append(report.ReportData, menu)
		j.dealYHCModule(module, menu)
	}
	return report
}

func (j *JsonParser) dealYHCModule(module *confdef.YHCModuleNode, menu *define.PandoraMenu) {
	if module == nil {
		return
	}
	for _, metricName := range module.MetricNames {
		result, ok := j.results[define.MetricName(metricName)]
		if !ok {
			continue
		}
		metric, err := j.getMetric(metricName)
		if err != nil {
			continue
		}
		fn, err := j.genMetricParseFunc(metric)
		if err != nil {
			j.log.Errorf("failed to gen parse func of metric %s", metricName)
			continue
		}
		childMenu := &define.PandoraMenu{Title: metric.NameAlias}
		if err := fn(childMenu, result, metric); err != nil {
			j.log.Errorf("failed to parse metric %s, err: %v", metricName, err)
			continue
		}
		menu.Children = append(menu.Children, childMenu)
	}
	for _, child := range module.Children {
		j.dealYHCModule(child, menu)
	}
}

func (j *JsonParser) genMetricParseFunc(metric *confdef.YHCMetric) (MetricParseFunc, error) {
	if !metric.Default {
		switch metric.MetricType {
		case confdef.MT_SQL:
			return j.genCustomSqlParseFunc(metric)
		case confdef.MT_BASH:
			return j.genCustomBashParseFunc(metric)
		default:
			return nil, fmt.Errorf("invalid metric type %s", metric.MetricType)
		}
	}
	return j.genDefaultMetricParseFunc(metric)
}

func (j *JsonParser) genDefaultMetricParseFunc(metric *confdef.YHCMetric) (MetricParseFunc, error) {
	parseFuncMap := map[define.MetricName]MetricParseFunc{
		define.METRIC_YASDB_INSTANCE:           j.parseMap,
		define.METRIC_YASDB_DATABASE:           j.parseMap,
		define.METRIC_YASDB_FILE_PERMISSION:    j.parseTable,
		define.METRIC_YASDB_LISTEN_ADDR:        j.parseMap,
		define.METRIC_HOST_INFO:                j.parseMap,
		define.METRIC_HOST_CPU_INFO:            j.parseMap,
		define.METRIC_HOST_HISTORY_CPU_USAGE:   j.parseHostCPUUsage,
		define.METRIC_HOST_CURRENT_CPU_USAGE:   j.parseHostCPUUsage,
		define.METRIC_YASDB_REPLICATION_STATUS: j.parseTable,
		define.METRIC_YASDB_PARAMETER:          j.parseMap,
		define.METRIC_YASDB_TABLESPACE:         j.parseTable,
		define.METRIC_YASDB_CONTROLFILE_COUNT:  j.parseMap,
		define.METRIC_YASDB_CONTROLFILE:        j.parseTable,
		define.METRIC_YASDB_SESSION:            j.parseTable,
		define.METRIC_YASDB_WAIT_EVENT:         j.parseTable,
		define.METRIC_YASDB_OBJECT_COUNT:       j.parseMap,
		define.METRIC_YASDB_OBJECT_OWNER:       j.parseTable,
		define.METRIC_YASDB_OBJECT_TABLESPACE:  j.parseTable,
		define.METRIC_YASDB_INDEX_BLEVEL:       j.parseTable,
		define.METRIC_YASDB_INDEX_COLUMN:       j.parseTable,
		define.METRIC_YASDB_INDEX_INVISIBLE:    j.parseTable,
		define.METRIC_YASDB_REDO_LOG:           j.parseTable,
		define.METRIC_YASDB_REDO_LOG_COUNT:     j.parseTable,
		define.METRIC_YASDB_RUN_LOG_ERROR:      j.parseText,
	}
	fn, ok := parseFuncMap[define.MetricName(metric.Name)]
	if !ok {
		return nil, fmt.Errorf("failed to find parse func of metric %s", metric.Name)
	}
	return fn, nil
}

func (j *JsonParser) getMetric(name string) (*confdef.YHCMetric, error) {
	for _, metric := range j.metrics {
		if metric.Name == name {
			return metric, nil
		}
	}
	return nil, fmt.Errorf("failed to found metric by %s, may be it does not check", name)
}

func (j *JsonParser) parseTable(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse table of %s because the details is nil", item.Name)
	}
	attributes := define.TableAttributes{
		Title: define.CustomOptionTitle{
			Text: metric.NameAlias,
		},
	}
	switch item.Details.(type) {
	case map[string]string:
		j.dealTableStringRow(&attributes, metric, item.Details.(map[string]string))
	case map[string]interface{}:
		j.dealTableAnyRow(&attributes, metric, item.Details.(map[string]interface{}))
	case []map[string]string:
		datas := item.Details.([]map[string]string)
		for _, data := range datas {
			j.dealTableStringRow(&attributes, metric, data)
		}
	case []map[string]interface{}:
		datas := item.Details.([]map[string]interface{})
		for _, data := range datas {
			j.dealTableAnyRow(&attributes, metric, data)
		}
	default:
		return fmt.Errorf("failed to parse table, unsupport data type %T", item.Details)
	}
	element := &define.PandoraElement{
		ElementType: define.ET_TABLE,
		Attributes:  attributes,
	}
	menu.Elements = append(menu.Elements, element)
	return j.parseAlert(menu, item, metric)
}

func (j *JsonParser) dealTableStringRow(attributes *define.TableAttributes, metric *confdef.YHCMetric, data map[string]string) {
	if len(attributes.TableColumns) == 0 {
		columnsMap := make(map[string]*define.TableColumn)
		for key := range data {
			title := j.getColumnAlias(metric, key)
			column := &define.TableColumn{
				Title:     title,
				DataIndex: key,
			}
			columnsMap[key] = column
		}
		columns := []*define.TableColumn{}
		for _, column := range columnsMap {
			columns = append(columns, column)
		}
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].DataIndex < columns[j].DataIndex
		})
		attributes.TableColumns = columns
	}
	dataSource := make(map[string]interface{})
	for key, value := range data {
		dataSource[key] = value
	}
	attributes.DataSource = append(attributes.DataSource, dataSource)
}

func (j *JsonParser) SortTableColumns(metric *confdef.YHCMetric, columns []*define.TableColumn) []*define.TableColumn {
	// todo:
	return nil
}

func (j *JsonParser) dealTableAnyRow(attributes *define.TableAttributes, metric *confdef.YHCMetric, data map[string]interface{}) {
	if len(attributes.TableColumns) == 0 {
		columnsMap := make(map[string]*define.TableColumn)
		for key := range data {
			title := j.getColumnAlias(metric, key)
			column := &define.TableColumn{
				Title:     title,
				DataIndex: key,
			}
			columnsMap[key] = column
		}
		columns := []*define.TableColumn{}
		for _, column := range columnsMap {
			columns = append(columns, column)
		}
		sort.Slice(columns, func(i, j int) bool {
			return columns[i].DataIndex < columns[j].DataIndex
		})
		attributes.TableColumns = columns
	}
	dataSource := make(map[string]interface{})
	for key, value := range data {
		dataSource[key] = value
	}
	attributes.DataSource = append(attributes.DataSource, dataSource)
}

func (j *JsonParser) parseCode(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse code of %s because the details is nil", item.Name)
	}
	attributes := define.CodeAttributes{
		Language: "shell",
	}
	switch item.Details.(type) {
	case string:
		code := item.Details.(string)
		attributes.Code = code
	default:
		return fmt.Errorf("failed to parse code, unsupport type %T", item.Details)
	}
	menu.Elements = append(menu.Elements, &define.PandoraElement{
		ElementType: define.ET_CODE,
		Attributes:  attributes,
	})
	return nil
}

func (j *JsonParser) parseMap(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse map of %s because the details is nil", item.Name)
	}
	element := &define.PandoraElement{
		ElementType: define.ET_DESCRIPTION,
	}
	attributes := define.DescriptionAttributes{}
	switch item.Details.(type) {
	case map[string]string:
		datas := item.Details.(map[string]string)
		for key, value := range datas {
			label := j.getColumnAlias(metric, key)
			attributes.Data = append(attributes.Data, &define.DescriptionData{Label: label, Value: value})
		}
	case map[string]interface{}:
		datas := item.Details.(map[string]interface{})
		for key, value := range datas {
			label := j.getColumnAlias(metric, key)
			attributes.Data = append(attributes.Data, &define.DescriptionData{Label: label, Value: value})
		}
	default:
		return fmt.Errorf("failed to parse map, unsupport data type %T", item.Details)
	}
	sort.Slice(attributes.Data, func(i, j int) bool {
		return attributes.Data[i].Label < attributes.Data[j].Label
	})
	element.Attributes = attributes
	menu.Elements = append(menu.Elements, element)
	return j.parseAlert(menu, item, metric)
}

func (j *JsonParser) parseText(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse text of %s because the details is nil", item.Name)
	}
	element := define.PandoraElement{
		ElementType: define.ET_PRE,
	}
	switch item.Details.(type) {
	case string:
		text := item.Details.(string)
		element.InnerText = text
	case []string:
		texts := item.Details.([]string)
		element.InnerText = strings.Join(texts, stringutil.STR_NEWLINE)
	default:
		return fmt.Errorf("failed to parse code, unsupport type %T", item.Details)
	}
	menu.Elements = append(menu.Elements, &element)
	return nil
}

func (j *JsonParser) parseAlert(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if len(item.Alerts) == 0 {
		return nil
	}
	for _, alerts := range item.Alerts {
		for _, alert := range alerts {
			element := define.PandoraElement{
				ElementType: define.ET_ALERT,
				Attributes: define.AlertAttributes{
					Message:     alert.Description,
					AlertType:   define.AlertType(alert.Level),
					Description: j.genAlertDescription(metric, alert),
				},
			}
			menu.Elements = append(menu.Elements, &element)
		}
	}
	return nil
}

func (j *JsonParser) genAlertDescription(metric *confdef.YHCMetric, alert *define.YHCAlert) string {
	desc := fmt.Sprintf("表达式：%s，值：%v，告警建议：%s", alert.Expression, alert.Value, alert.Suggestion)
	if len(alert.Labels) != 0 {
		labelArr := []string{}
		for k, v := range alert.Labels {
			labelAlias := j.getColumnAlias(metric, k)
			labelArr = append(labelArr, fmt.Sprintf("%s：%s", labelAlias, v))
		}
		desc = fmt.Sprintf("%s，标签：{%s}", desc, strings.Join(labelArr, "；"))
	}
	return desc
}

// 部分指标由于sql限制，分开采集，生成报告的时候需要合并到同一张表格中
func (j *JsonParser) MergeMetrics() {
	for to, from := range _mergeMetricMap {
		j.mergeMetric(to, from)
	}
}

func (j *JsonParser) getColumnAlias(metric *confdef.YHCMetric, columnName string) string {
	relatedMetrics := j.getRelatedMetrics(metric)
	for _, metricName := range relatedMetrics {
		metric, err := j.getMetric(string(metricName))
		if err != nil {
			j.log.Errorf("failed to get metric by name %s", metricName)
			continue
		}
		alias, ok := metric.ColumnAlias[columnName]
		if ok {
			return alias
		}
	}
	return columnName
}

// 部分指标在展示的时候需要合并信息，此函数返回当前指标的关联指标
func (j *JsonParser) getRelatedMetrics(metric *confdef.YHCMetric) []define.MetricName {
	var relatedMap = map[define.MetricName][]define.MetricName{}
	for _, metric := range _yasdbOverviewRelated {
		relatedMap[metric] = _yasdbOverviewRelated
	}
	for _, metric := range _hostOverviewRelated {
		relatedMap[metric] = _hostOverviewRelated
	}
	related, ok := relatedMap[define.MetricName(metric.Name)]
	if !ok {
		related = []define.MetricName{define.MetricName(metric.Name)}
	}
	return related
}

func (j *JsonParser) mergeMetric(to define.MetricName, from []define.MetricName) {
	resDetail := make(map[string]interface{})
	resAlerts := make(map[string][]*define.YHCAlert)
	var merge bool
	for _, m := range from {
		fromResult, ok := j.results[m]
		if !ok {
			continue
		}
		detail := fromResult.Details
		switch detailType := detail.(type) {
		case map[string]interface{}:
			data := detail.(map[string]interface{})
			for k, v := range data {
				resDetail[k] = v
			}
		case map[string]string:
			data := detail.(map[string]string)
			for k, v := range data {
				resDetail[k] = v
			}
		default:
			j.log.Errorf("failed to merge metrics, unsupport data type %T", detailType)
			continue
		}

		for level, alerts := range fromResult.Alerts {
			a, ok := resAlerts[level]
			if !ok {
				a = []*define.YHCAlert{}
			}
			a = append(a, alerts...)
			resAlerts[level] = a
		}
		delete(j.results, m)
		merge = true
	}
	if merge {
		j.results[to] = &define.YHCItem{
			Name:    to,
			Details: resDetail,
			Alerts:  resAlerts,
		}
	}
}

func (j *JsonParser) genCustomBashParseFunc(metric *confdef.YHCMetric) (MetricParseFunc, error) {
	fn := func(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
		if len(item.Error) == 0 {
			return fmt.Errorf("failed to gen parse func because the metric %s check failed", metric.Name)
		}
		if err := j.parseCode(menu, item, metric); err != nil {
			return err
		}
		if err := j.parseAlert(menu, item, metric); err != nil {
			return err
		}
		return nil
	}
	return fn, nil
}

func (j *JsonParser) genCustomSqlParseFunc(metric *confdef.YHCMetric) (MetricParseFunc, error) {
	fn := func(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
		if len(item.Error) == 0 {
			return fmt.Errorf("failed to gen parse func because the metric %s check failed", metric.Name)
		}
		if err := j.parseTable(menu, item, metric); err != nil {
			return err
		}
		if err := j.parseAlert(menu, item, metric); err != nil {
			return err
		}
		return nil
	}
	return fn, nil
}

func (j *JsonParser) parseHostCPUUsage(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	cpuExcludeMap := map[string]struct{}{
		"cpu": {},
	}
	return j.parseHostWorkload(menu, item, metric, cpuExcludeMap)
}

func (j *JsonParser) parseHostWorkload(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric, excludeMap map[string]struct{}) error {
	if len(item.Error) != 0 {
		return fmt.Errorf("failed to gen parse func because the metric %s check failed", metric.Name)
	}
	data, ok := item.Details.(define.WorkloadOutput)
	if !ok {
		return fmt.Errorf("invalid data type %T", item.Details)
	}
	if len(data) == 0 {
		return nil
	}
	timeArray := []int64{}
	for time := range data {
		timeArray = append(timeArray, time)
	}
	sort.Slice(timeArray, func(i, j int) bool { return timeArray[i] < timeArray[j] })

	// create attributes map to store all attribute
	attributes := make(map[string]define.ChartAttributes)
	for name := range data[timeArray[0]] {
		attribute := define.ChartAttributes{
			CustomOptions: define.ChartCustomOptions{
				ChartType: define.CT_LINE,
				Title: define.CustomOptionTitle{
					Text: name,
				},
				Data: []*define.ChartData{},
			},
		}
		attributes[name] = attribute
	}

	// fill chart data from origin data
	for _, t := range timeArray {
		timeStr := time.Unix(t, 0).Format(timedef.TIME_FORMAT)
		for name, value := range data[t] {
			// parse data to map
			m, err := j.convertObjectToMap(value)
			if err != nil {
				j.log.Errorf("failed to parse object %T, err: %v", value, err)
				continue
			}
			// use map to record data
			attribute := attributes[name]
			chartDataMap := make(map[string]*define.ChartData)
			for _, d := range attribute.CustomOptions.Data {
				chartDataMap[d.Name] = d
			}
			for lineName, lineValue := range m {
				if _, ok := excludeMap[lineName]; ok {
					j.log.Debugf("column %s in exclude map, skip", lineName)
					continue
				}
				alias := j.getColumnAlias(metric, lineName)
				chartData, ok := chartDataMap[lineName]
				if !ok {
					chartData = &define.ChartData{
						Name: alias,
					}
				}
				chartData.Value = append(chartData.Value, &define.ChartCoordinate{X: timeStr, Y: lineValue})
				chartDataMap[lineName] = chartData
			}
			chartDatas := []*define.ChartData{}
			for _, d := range chartDataMap {
				chartDatas = append(chartDatas, d)
			}
			attribute.CustomOptions.Data = chartDatas
			attributes[name] = attribute
		}
	}
	for _, attribute := range attributes {
		menu.Elements = append(menu.Elements, &define.PandoraElement{
			ElementType: define.ET_CHART,
			Attributes:  attribute,
		})
	}
	return nil
}

func (j *JsonParser) convertObjectToMap(object interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	bytes, err := json.Marshal(object)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(bytes, &m); err != nil {
		return m, err
	}
	return m, nil
}
