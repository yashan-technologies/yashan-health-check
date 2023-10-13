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

// 将不同指标的数据合并到一个map中，只支持map之间的合并
var _mergeMetricMap = map[define.MetricName][]define.MetricName{
	define.METRIC_HOST_INFO: {
		define.METRIC_HOST_INFO,
		define.METRIC_HOST_CPU_INFO,
	},
	define.METRIC_YASDB_DATABASE: {
		define.METRIC_YASDB_DATABASE,
		define.METRIC_YASDB_INSTANCE,
		define.METRIC_YASDB_LISTEN_ADDR,
	},
}

// 将不同指标的element放在一个指标下
var _mergeElementMap = map[define.MetricName][]define.MetricName{
	define.METRIC_HOST_INFO: {
		define.METRIC_HOST_INFO,
		define.METRIC_HOST_CPU_INFO,
		define.METRIC_HOST_DISK_INFO,
		define.METRIC_HOST_DISK_BLOCK_INFO,
		define.METRIC_HOST_MEMORY_INFO,
		define.METRIC_HOST_NETWORK_INFO,
	},
	define.METRIC_YASDB_DATABASE: {
		define.METRIC_YASDB_DATABASE,
		define.METRIC_YASDB_FILE_PERMISSION,
	},
	define.METRIC_YASDB_OBJECT_COUNT: {
		define.METRIC_YASDB_OBJECT_COUNT,
		define.METRIC_YASDB_OBJECT_TABLESPACE,
		define.METRIC_YASDB_OBJECT_OWNER,
	},
	define.METRIC_YASDB_REDO_LOG: {
		define.METRIC_YASDB_REDO_LOG,
		define.METRIC_YASDB_REDO_LOG_COUNT,
	},
}

type MetricParseFunc func(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error

type JsonParser struct {
	log            yaslog.YasLog
	base           define.CheckerBase
	startCheckTime time.Time
	endCheckTime   time.Time
	metrics        []*confdef.YHCMetric
	results        map[define.MetricName]*define.YHCItem
}

func NewJsonParser(log yaslog.YasLog, base define.CheckerBase, startCheck, endCheck time.Time, metrics []*confdef.YHCMetric, results map[define.MetricName]*define.YHCItem) *JsonParser {
	parser := &JsonParser{
		log:            log,
		metrics:        metrics,
		results:        results,
		startCheckTime: startCheck,
		endCheckTime:   endCheck,
		base:           base,
	}
	return parser
}

// todo: 这个parse函数各个模块之间的关系处理有点问题，需要优化
// todo: 包括wordgenner的模块处理也有问题，后续优化！
func (j *JsonParser) Parse() *define.PandoraReport {
	report := &define.PandoraReport{
		ReportTitle: _REPORT_TITLE,
		FileControl: _FILE_CONTROL,
		Author:      _AUTHOR,
		ChangeLog:   _CHANGE_LOG,
		Time:        j.startCheckTime.Format(timedef.TIME_FORMAT),
		CostTime:    int(j.endCheckTime.Sub(j.startCheckTime).Seconds()),
		Version:     compiledef.GetAPPVersion(),
	}
	j.mergeMetrics()
	for _, module := range confdef.GetModuleConf().Modules {
		menu := &define.PandoraMenu{IsMenu: true, Title: confdef.GetModuleAlias(module.Name)}
		report.ReportData = append(report.ReportData, menu)
		j.dealYHCModule(module, menu)
	}
	j.mergeElements(report)
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
		define.METRIC_YASDB_INSTANCE:                                         j.parseMap,
		define.METRIC_YASDB_DATABASE:                                         j.parseMap,
		define.METRIC_YASDB_FILE_PERMISSION:                                  j.parseTable,
		define.METRIC_YASDB_LISTEN_ADDR:                                      j.parseMap,
		define.METRIC_YASDB_OS_AUTH:                                          j.parseMap,
		define.METRIC_HOST_INFO:                                              j.parseMap,
		define.METRIC_HOST_FIREWALLD:                                         j.parseMap,
		define.METRIC_HOST_IPTABLES:                                          j.parseCode,
		define.METRIC_HOST_CPU_INFO:                                          j.parseMap,
		define.METRIC_HOST_DISK_INFO:                                         j.parseTable,
		define.METRIC_HOST_DISK_BLOCK_INFO:                                   j.parseTable,
		define.METRIC_HOST_BIOS_INFO:                                         j.parseCode,
		define.METRIC_HOST_MEMORY_INFO:                                       j.parseTable,
		define.METRIC_HOST_NETWORK_INFO:                                      j.parseTable,
		define.METRIC_HOST_HISTORY_CPU_USAGE:                                 j.parseHostWorkload,
		define.METRIC_HOST_CURRENT_CPU_USAGE:                                 j.parseHostWorkload,
		define.METRIC_HOST_CURRENT_DISK_IO:                                   j.parseHostWorkload,
		define.METRIC_HOST_HISTORY_DISK_IO:                                   j.parseHostWorkload,
		define.METRIC_HOST_CURRENT_MEMORY_USAGE:                              j.parseHostWorkload,
		define.METRIC_HOST_HISTORY_MEMORY_USAGE:                              j.parseHostWorkload,
		define.METRIC_HOST_CURRENT_NETWORK_IO:                                j.parseHostWorkload,
		define.METRIC_HOST_HISTORY_NETWORK_IO:                                j.parseHostWorkload,
		define.METRIC_YASDB_REPLICATION_STATUS:                               j.parseTable,
		define.METRIC_YASDB_PARAMETER:                                        j.parseMap,
		define.METRIC_YASDB_TABLESPACE:                                       j.parseTable,
		define.METRIC_YASDB_CONTROLFILE_COUNT:                                j.parseMap,
		define.METRIC_YASDB_CONTROLFILE:                                      j.parseTable,
		define.METRIC_YASDB_DATAFILE:                                         j.parseTable,
		define.METRIC_YASDB_SESSION:                                          j.parseMap,
		define.METRIC_YASDB_WAIT_EVENT:                                       j.parseTable,
		define.METRIC_YASDB_OBJECT_COUNT:                                     j.parseMap,
		define.METRIC_YASDB_OBJECT_OWNER:                                     j.parseTable,
		define.METRIC_YASDB_OBJECT_TABLESPACE:                                j.parseTable,
		define.METRIC_YASDB_INDEX_BLEVEL:                                     j.parseTable,
		define.METRIC_YASDB_INDEX_COLUMN:                                     j.parseTable,
		define.METRIC_YASDB_INDEX_INVISIBLE:                                  j.parseTable,
		define.METRIC_YASDB_REDO_LOG:                                         j.parseTable,
		define.METRIC_YASDB_REDO_LOG_COUNT:                                   j.parseTable,
		define.METRIC_YASDB_RUN_LOG_ERROR:                                    j.parseText,
		define.METRIC_YASDB_INDEX_TABLE_INDEX_NOT_TOGETHER:                   j.parseTable,
		define.METRIC_YASDB_INDEX_OVERSIZED:                                  j.parseTable,
		define.METRIC_YASDB_SEQUENCE_NO_AVAILABLE:                            j.parseTable,
		define.METRIC_YASDB_TASK_RUNNING:                                     j.parseTable,
		define.METRIC_YASDB_PACKAGE_NO_PACKAGE_PACKAGE_BODY:                  j.parseTable,
		define.METRIC_YASDB_SECURITY_LOGIN_PASSWORD_STRENGTH:                 j.parseTable,
		define.METRIC_YASDB_SECURITY_LOGIN_MAXIMUM_LOGIN_ATTEMPTS:            j.parseTable,
		define.METRIC_YASDB_SECURITY_USER_NO_OPEN:                            j.parseTable,
		define.METRIC_YASDB_SECURITY_USER_WITH_SYSTEM_TABLE_PRIVILEGES:       j.parseTable,
		define.METRIC_YASDB_SECURITY_USER_WITH_DBA_ROLE:                      j.parseTable,
		define.METRIC_YASDB_SECURITY_USER_ALL_PRIVILEGE_OR_SYSTEM_PRIVILEGES: j.parseTable,
		define.METRIC_YASDB_SECURITY_USER_USE_SYSTEM_TABLESPACE:              j.parseTable,
		define.METRIC_YASDB_SECURITY_AUDIT_CLEANUP_TASK:                      j.parseTable,
		define.METRIC_YASDB_SECURITY_AUDIT_FILE_SIZE:                         j.parseTable,
		define.METRIC_YASDB_UNDO_LOG_SIZE:                                    j.parseTable,
		define.METRIC_YASDB_UNDO_LOG_TOTAL_BLOCK:                             j.parseTable,
		define.METRIC_YASDB_UNDO_LOG_RUNNING_TRANSACTIONS:                    j.parseTable,
		define.METRIC_YASDB_RUN_LOG_DATABASE_CHANGES:                         j.parseText,
		define.METRIC_YASDB_ALERT_LOG_ERROR:                                  j.parseText,
		define.METRIC_HOST_DMESG_LOG_ERROR:                                   j.parseText,
		define.METRIC_HOST_SYSTEM_LOG_ERROR:                                  j.parseText,
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
		Title: metric.NameAlias,
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
	attributes.TableColumns = j.sortTableColumns(metric, attributes.TableColumns)
	element := &define.PandoraElement{
		MetricName:   metric.Name,
		ElementTitle: metric.NameAlias,
		ElementType:  define.ET_TABLE,
		Attributes:   attributes,
	}
	menu.Elements = append(menu.Elements, element)
	return j.parseAlert(menu, item, metric)
}

func (j *JsonParser) sortTableColumns(metric *confdef.YHCMetric, columns []*define.TableColumn) []*define.TableColumn {
	columnMap := map[string]*define.TableColumn{}
	for _, column := range columns {
		columnMap[column.DataIndex] = column
	}
	var order, unorder []*define.TableColumn
	relatedMetric := j.getRelatedMetrics(metric)
	for _, metricName := range relatedMetric {
		m, err := j.getMetric(string(metricName))
		if err != nil {
			j.log.Error(err)
			continue
		}
		for _, o := range m.ColumnOrder {
			if column, ok := columnMap[o]; ok {
				order = append(order, column)
				delete(columnMap, o)
			}
		}
	}
	for _, column := range columnMap {
		unorder = append(unorder, column)
	}
	sort.Slice(unorder, func(i, j int) bool {
		return unorder[i].DataIndex < unorder[j].DataIndex
	})
	return append(order, unorder...)
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
		attributes.TableColumns = columns
	}
	dataSource := make(map[string]interface{})
	for key, value := range data {
		dataSource[key] = value
	}
	attributes.DataSource = append(attributes.DataSource, dataSource)
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
		Title:    confdef.GetModuleAlias(metric.Name),
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
		MetricName:   metric.Name,
		ElementTitle: metric.NameAlias,
		ElementType:  define.ET_CODE,
		Attributes:   attributes,
	})
	return nil
}

func (j *JsonParser) parseMap(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse map of %s because the details is nil", item.Name)
	}
	element := &define.PandoraElement{
		MetricName:   metric.Name,
		ElementTitle: metric.NameAlias,
		ElementType:  define.ET_DESCRIPTION,
	}
	attributes := define.DescriptionAttributes{}
	switch item.Details.(type) {
	case map[string]string:
		datas := item.Details.(map[string]string)
		for key, value := range datas {
			attributes.Data = append(attributes.Data, &define.DescriptionData{Label: key, Value: value})
		}
	case map[string]interface{}:
		datas := item.Details.(map[string]interface{})
		for key, value := range datas {
			attributes.Data = append(attributes.Data, &define.DescriptionData{Label: key, Value: value})
		}
	default:
		return fmt.Errorf("failed to parse map, unsupport data type %T", item.Details)
	}
	attributes.Data = j.sortMapData(metric, attributes.Data)
	element.Attributes = attributes
	menu.Elements = append(menu.Elements, element)
	return j.parseAlert(menu, item, metric)
}

func (j *JsonParser) sortMapData(metric *confdef.YHCMetric, datas []*define.DescriptionData) []*define.DescriptionData {
	dataMap := map[string]*define.DescriptionData{}
	for _, data := range datas {
		dataMap[data.Label] = data
	}
	var order, unorder []*define.DescriptionData
	relatedMetric := j.getRelatedMetrics(metric)
	for _, metricName := range relatedMetric {
		m, err := j.getMetric(string(metricName))
		if err != nil {
			j.log.Error(err)
			continue
		}
		for _, o := range m.ColumnOrder {
			if column, ok := dataMap[o]; ok {
				order = append(order, column)
				delete(dataMap, o)
			}
		}
	}
	for _, data := range dataMap {
		unorder = append(unorder, data)
	}
	sort.Slice(unorder, func(i, j int) bool {
		return unorder[i].Label < unorder[j].Label
	})
	order = append(order, unorder...)
	// replace with column alias
	for _, o := range order {
		o.Label = j.getColumnAlias(metric, o.Label)
	}
	return order
}

func (j *JsonParser) parseText(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	if item.Details == nil {
		return fmt.Errorf("failed to parse text of %s because the details is nil", item.Name)
	}
	element := define.PandoraElement{
		MetricName:   metric.Name,
		ElementTitle: metric.NameAlias,
		ElementType:  define.ET_PRE,
	}
	attributes := define.DescriptionAttributes{
		Title: metric.NameAlias,
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
	element.Attributes = attributes
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
				MetricName:  metric.Name,
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
func (j *JsonParser) mergeMetrics() {
	for to, from := range _mergeMetricMap {
		j.mergeMetric(to, from)
	}
}

func (j *JsonParser) mergeElements(report *define.PandoraReport) {
	for to, from := range _mergeElementMap {
		j.mergeElement(report, to, from)
	}
	j.cleanEmptyMenus(report.ReportData)
}

func (j *JsonParser) mergeElement(report *define.PandoraReport, to define.MetricName, from []define.MetricName) {
	toMenu, _ := j.findMenuAndElements(report, to)
	if toMenu == nil {
		return
	}
	for _, m := range from {
		if to == m {
			continue
		}
		menu, elements := j.findMenuAndElements(report, m)
		if len(elements) == 0 {
			continue
		}
		toMenu.Elements = append(toMenu.Elements, elements...)
		j.deleteElementsFromMenu(menu, m)
	}
}

func (j *JsonParser) findMenuAndElements(report *define.PandoraReport, metricName define.MetricName) (targetMenu *define.PandoraMenu, elements []*define.PandoraElement) {
	for _, menu := range report.ReportData {
		targetMenu, elements = j.findElementsInMenu(menu, metricName)
		if len(elements) > 0 {
			break
		}
	}
	return
}

func (j *JsonParser) findElementsInMenu(menu *define.PandoraMenu, metricName define.MetricName) (findMenu *define.PandoraMenu, findElements []*define.PandoraElement) {
	if menu == nil {
		return
	}
	for _, element := range menu.Elements {
		if element.MetricName == string(metricName) {
			findElements = append(findElements, element)
			findMenu = menu
		}
	}
	if len(findElements) > 0 {
		return
	}
	for _, childMenu := range menu.Children {
		findMenu, findElements = j.findElementsInMenu(childMenu, metricName)
		if len(findElements) > 0 {
			return
		}
	}
	return
}

func (j *JsonParser) deleteElementsFromMenu(menu *define.PandoraMenu, metricName define.MetricName) {
	var updatedElements []*define.PandoraElement
	for _, element := range menu.Elements {
		if element.MetricName != string(metricName) {
			updatedElements = append(updatedElements, element)
		}
	}
	menu.Elements = updatedElements
}

func (j *JsonParser) cleanEmptyMenus(menus []*define.PandoraMenu) []*define.PandoraMenu {
	var cleanedMenus []*define.PandoraMenu
	for _, menu := range menus {
		if len(menu.Elements) == 0 && len(menu.Children) == 0 {
			continue
		}
		menu.Children = j.cleanEmptyMenus(menu.Children)
		if len(menu.Elements) > 0 || len(menu.Children) > 0 {
			cleanedMenus = append(cleanedMenus, menu)
		}
	}
	return cleanedMenus
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
	for metricName, related := range _mergeMetricMap {
		if metricName == define.MetricName(metric.Name) {
			return related
		}
	}
	return []define.MetricName{define.MetricName(metric.Name)}
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

func (j *JsonParser) parseHostOtherWorkload(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric, includeFields map[string]struct{}) error {
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
	for _, value := range data[timeArray[0]] {
		m, err := j.convertObjectToMap(value)
		if err != nil {
			return err
		}
		for field := range m {
			if _, ok := includeFields[field]; !ok {
				continue
			}
			attribute := define.ChartAttributes{
				CustomOptions: define.ChartCustomOptions{
					ChartType: define.CT_LINE,
					Title: define.CustomOptionTitle{
						Text: j.getColumnAlias(metric, field),
					},
					Data: []*define.ChartData{},
				},
			}
			attributes[field] = attribute
		}
	}

	// fill chart data from origin data
	for _, t := range timeArray {
		timeStr := time.Unix(t, 0).Format(timedef.TIME_FORMAT)
		for name, obj := range data[t] {
			// parse data to map
			m, err := j.convertObjectToMap(obj)
			if err != nil {
				j.log.Errorf("failed to parse object %T, err: %v", obj, err)
				continue
			}
			for field, value := range m {
				if _, ok := includeFields[field]; !ok {
					continue
				}
				attribute := attributes[field]

				chartDataMap := make(map[string]*define.ChartData)
				for _, d := range attribute.CustomOptions.Data {
					chartDataMap[d.Name] = d
				}
				chartData, ok := chartDataMap[name]
				if !ok {
					chartData = &define.ChartData{Name: name}
				}
				chartData.Value = append(chartData.Value, &define.ChartCoordinate{X: timeStr, Y: value})
				chartDataMap[name] = chartData
				chartDatas := []*define.ChartData{}
				for _, d := range chartDataMap {
					chartDatas = append(chartDatas, d)
				}
				attribute.CustomOptions.Data = chartDatas
				attributes[field] = attribute
			}
		}
	}
	for _, attribute := range attributes {
		menu.Elements = append(menu.Elements, &define.PandoraElement{
			MetricName:  metric.Name,
			ElementType: define.ET_CHART,
			Attributes:  attribute,
		})
	}
	return nil
}

func (j *JsonParser) parseHostWorkload(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric) error {
	includeFields := map[string]struct{}{}
	for column := range metric.ColumnAlias {
		includeFields[column] = struct{}{}
	}
	switch item.Name {
	case define.METRIC_HOST_CURRENT_CPU_USAGE, define.METRIC_HOST_HISTORY_CPU_USAGE:
		return j.parseHostCPUUsage(menu, item, metric, includeFields)
	default:
		return j.parseHostOtherWorkload(menu, item, metric, includeFields)
	}
}

func (j *JsonParser) parseHostCPUUsage(menu *define.PandoraMenu, item *define.YHCItem, metric *confdef.YHCMetric, includeFields map[string]struct{}) error {
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
					Text: metric.NameAlias,
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
				if _, ok := includeFields[lineName]; !ok {
					continue
				}
				chartData, ok := chartDataMap[lineName]
				if !ok {
					chartData = &define.ChartData{
						Name: lineName,
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
		datas := attribute.CustomOptions.Data
		for _, data := range datas {
			data.Name = j.getColumnAlias(metric, data.Name)
		}
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
