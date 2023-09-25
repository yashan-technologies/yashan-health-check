package checkcontroller

import (
	"reflect"

	"yhc/commons/std"
	"yhc/commons/yasdb"
	"yhc/defs/confdef"
	constdef "yhc/defs/constants"
	"yhc/defs/errdef"
	"yhc/internal/modules/yhc/check"
	"yhc/internal/modules/yhc/check/define"
	yhcyasdb "yhc/internal/modules/yhc/yasdb"
	"yhc/log"
	"yhc/utils/jsonutil"

	"git.yasdb.com/go/yasutil/tabler"
	"git.yasdb.com/pandora/tview"
	"git.yasdb.com/pandora/yasqlgo"
	"github.com/gdamore/tcell/v2"
)

const (
	EXIT_CONTINUE = iota
	EXIT_NOT_CONTINUE
	EXIT_CONTROL_C
)

const (
	// yashan health check
	_header = "Yashan Health Check"

	// terminal view item chinese name
	_module               = "模块"
	_check_item           = "检查项"
	_detail               = "详情"
	_health_check_summary = "健康检查概览"
	_tips_header          = "以下检查项,不会进行检查,详细如下"
	_next_button_name     = "下一步"
	_exit_button_name     = "退出"

	// yashan health check page name
	_yasdb   = "yasdb"
	_tips    = "tips"
	_summary = "summary"

	// summary flex index
	_summary_metric_flex_index = 1
	_summary_table_flex_index  = 2

	_check_list_width          = 30
	_table_cell_max_width      = 50
	_validate_dba_sql          = check.SQL_QUERY_TOTAL_OBJECT
	_base_yasdb_process_format = `.*yasdb (?i:(nomount|mount|open))`
)

var (
	// terminal view exit code
	globalExitCode int

	// Filled in after yasdb page validation
	moduleNoNeedCheckMetrics = map[string]map[string]*define.NoNeedCheckMetric{}

	// Filled in after yasdb page validation
	globalFilterModule = []*constdef.ModuleMetrics{}

	alertRuleOrder = []string{
		confdef.AL_CRITICAL,
		confdef.AL_WARING,
		confdef.AL_INFO,
		confdef.AL_INVALID,
	}

	exitCodeMap = map[int]string{
		EXIT_CONTINUE:     "continue health check",
		EXIT_NOT_CONTINUE: "stop health check",
		EXIT_CONTROL_C:    "exit with control c",
	}

	tipsTableColumns = []string{"模块名称", "检查项名称", "原因"}

	alarmTableColumns = []string{"告警等级", "告警表达式", "描述", "建议"}
)

type PagePrimitive struct {
	Name      string          // page name
	Primitive tview.Primitive // page view
	Show      bool            // 是否展示
}

func StartTerminalView(modules []*constdef.ModuleMetrics, yasdb *yasdb.YashanDB) {
	app := tview.NewApplication()
	app.SetInputCapture(captureCtrlCFunc(app))
	if err := app.SetRoot(index(app, yasdb, modules), true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func captureCtrlCFunc(app *tview.Application) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			exitFunc(app, EXIT_CONTROL_C)
		}
		return event
	}
}

func index(app *tview.Application, yasdb *yasdb.YashanDB, modules []*constdef.ModuleMetrics) *tview.Flex {
	f := newFlex(_header, true, tview.FlexRow)
	yasdbPage := newYasdbPage(yasdb)
	pages := newPages(yasdbPage)
	f.AddItem(pages, 0, 1, true)
	f.AddItem(indexFooter(app, pages, f, modules), 3, 1, false)
	return f
}

func indexFooter(app *tview.Application, page *tview.Pages, index *tview.Flex, modules []*constdef.ModuleMetrics) *tview.Flex {
	f := newFlex("", false, tview.FlexColumn)
	next := newButton(_next_button_name, true)
	exit := newButton(_exit_button_name, true)
	next.SetSelectedFunc(nextClickFunc(app, page, next, index, modules))
	exit.SetSelectedFunc(func() { exitFunc(app, EXIT_NOT_CONTINUE) })
	f.AddItem(next, 0, 1, false)
	f.AddItem(exit, 0, 1, false)
	return f
}

func exitFunc(app *tview.Application, code int) {
	globalExitCode = code
	app.Stop()
}

// yashandb info form page
func dbInfoPage(yasdb *yasdb.YashanDB) *tview.Form {
	form := tview.NewForm()
	form.SetTitleAlign(tview.AlignCenter)
	form.AddInputField(constdef.YASDB_HOME, yasdb.YasdbHome, 100, nil, func(text string) { yasdb.YasdbHome = trimSpace(text) })
	form.AddInputField(constdef.YASDB_DATA, yasdb.YasdbData, 100, nil, func(text string) { yasdb.YasdbData = trimSpace(text) })
	form.AddInputField(constdef.YASDB_USER, yasdb.YasdbUser, 100, nil, func(text string) { yasdb.YasdbUser = trimSpace(text) })
	form.AddPasswordField(constdef.YASDB_PASSWORD, yasdb.YasdbPassword, 100, '*', func(text string) { yasdb.YasdbPassword = trimSpace(text) })
	return form
}

// yashan health check metric summary page
func summaryFlexPage(modules []*constdef.ModuleMetrics) *tview.Flex {
	flex := newFlex("", true, tview.FlexRow)
	header := newTextView(_health_check_summary, true, tview.AlignCenter, tcell.ColorYellow)
	body := summaryBody(modules)
	flex.AddItem(header, 3, 1, false)
	flex.AddItem(body, 0, 1, false)
	return flex
}

// body of summary page
func summaryBody(modules []*constdef.ModuleMetrics) *tview.Flex {
	flex := tview.NewFlex()
	fillSummaryBody(modules, flex)
	return flex
}

func fillSummaryBody(modules []*constdef.ModuleMetrics, flex *tview.Flex) {
	moduleList := newCheckedList(_module, true)
	itemList := newCheckedList(_check_item, true)
	table := summaryTable(modules)
	flex.AddItem(moduleList, _check_list_width, 1, false)
	flex.AddItem(itemList, _check_list_width, 1, false)
	flex.AddItem(table, 0, 1, false)
	moduleList.SetCheckedFunc(func(index int, checked bool) { modules[index].Enabled = checked })
	moduleList.SetChangedFunc(func(i int, s1, s2 string, r rune) { moduleListChangedFunc(modules[i], flex) })
	addModuleList(moduleList, modules)
	if len(modules) == 0 {
		return
	}
	addItemList(itemList, modules[0].Name, modules[0].Metrics)
	if len(modules[0].Metrics) == 0 {
		return
	}
	drawAlertRuleTable(table, modules[0].Metrics[0].AlertRules)
}

func addItemList(itemList *tview.CheckList, moduleName string, metrics []*confdef.YHCMetric) {
	for _, item := range metrics {
		itemList.AddItem(item.NameAlias, "", 0, nil, item.Enabled)
	}
}

func addModuleList(moduleList *tview.CheckList, modules []*constdef.ModuleMetrics) {
	for _, item := range modules {
		alias, err := define.GetModuleDefaultAlias(define.ModuleName(item.Name))
		if err != nil {
			log.Controller.Errorf("get module alias err: %s", err.Error())
			continue
		}
		moduleList.AddItem(alias, "", 0, nil, item.Enabled)
	}
}

func summaryTable(modules []*constdef.ModuleMetrics) *tview.Table {
	if len((modules)) == 0 {
		return nil
	}
	if len(modules[0].Metrics) == 0 {
		return nil
	}
	table := newTable(_detail, true, true)
	drawAlertRuleTable(table, modules[0].Metrics[0].AlertRules)
	return table
}

func moduleListChangedFunc(module *constdef.ModuleMetrics, flex *tview.Flex) {
	itemList := flex.GetItem(_summary_metric_flex_index)
	table := flex.GetItem(_summary_table_flex_index)
	newItemList := newCheckedList(_check_item, true)
	newTable := newTable(_detail, true, true)
	newItemList.SetCheckedFunc(func(index int, checked bool) { module.Metrics[index].Enabled = checked })
	newItemList.SetChangedFunc(func(i int, m, s string, sc rune) { itemListChangedFunc(module.Metrics[i], flex) })
	flex.RemoveItem(table)
	flex.RemoveItem(itemList)
	flex.AddItem(newItemList, _check_list_width, 1, false)
	flex.AddItem(table, 0, 1, false)
	addItemList(newItemList, module.Name, module.Metrics)
	if len(module.Metrics) == 0 {
		return
	}
	drawAlertRuleTable(newTable, module.Metrics[0].AlertRules)
}

func itemListChangedFunc(item *confdef.YHCMetric, flex *tview.Flex) {
	table := flex.GetItem(_summary_table_flex_index)
	newTable := newTable(_detail, true, true)
	drawAlertRuleTable(newTable, item.AlertRules)
	flex.RemoveItem(table)
	flex.AddItem(newTable, 0, 1, false)
}

func nextClickFunc(app *tview.Application, page *tview.Pages, button *tview.Button, index *tview.Flex, modules []*constdef.ModuleMetrics) func() {
	return func() {
		pageName, pageView := page.GetFrontPage()
		log.Controller.Infof("click next current page: %s", pageName)
		if pageName == _yasdb {
			// validate yasdb
			form, ok := pageView.(*tview.Form)
			if ok {
				yasdbEnv, err := yasdbValidate(form)
				if err != nil {
					modal := newModal(app, index, err)
					app.SetRoot(modal, true)
					return
				}
				metricValidate(yasdbEnv, modules)
				if len(moduleNoNeedCheckMetrics) != 0 {
					// write no need check metrics to console.log
					std.WriteToFile("the following metric will not be checked \n")
					noNeedStr := genNoNeedCheckMetricsStr()
					std.WriteToFile(noNeedStr)
					page.AddAndSwitchToPage(_tips, tipsPage(), true)
					return
				}
			}
		}
		if pageName == _summary {
			exitFunc(app, EXIT_CONTINUE)
			return
		}
		page.AddAndSwitchToPage(_summary, summaryFlexPage(globalFilterModule), true)
	}
}

func tipsPage() *tview.Flex {
	f := tview.NewFlex()
	f.SetDirection(tview.FlexRow)
	header := newTextView(_tips_header, true, tview.AlignCenter, tcell.ColorYellow)
	table := newTable(_detail, true, true)
	type moduleMetric struct {
		ModuleName  string
		MetricName  string
		Description string
	}
	tips := make([]interface{}, 0)
	for _, module := range define.Level1ModuleOrder {
		moduleStr := string(module)
		if _, ok := moduleNoNeedCheckMetrics[moduleStr]; !ok {
			continue
		}
		moduleAlias, _ := define.GetModuleDefaultAlias(module)
		for _, notCheck := range moduleNoNeedCheckMetrics[moduleStr] {
			tips = append(tips, &moduleMetric{
				ModuleName:  moduleAlias,
				MetricName:  notCheck.Name,
				Description: notCheck.Description,
			})
		}
	}
	fillTableCell(table, tipsTableColumns, tips)
	f.AddItem(header, 3, 1, false)
	f.AddItem(table, 0, 1, false)
	return f
}

func drawAlertRuleTable(table *tview.Table, alertRules map[string]confdef.AlertDetails) {
	type rule struct {
		Level       string
		Expression  string
		Description string
		Suggestion  string
	}
	rules := make([]interface{}, 0)
	for _, level := range alertRuleOrder {
		if _, ok := alertRules[level]; !ok {
			continue
		}
		rules = append(rules, rule{
			Level:       level,
			Expression:  alertRules[level].Expression,
			Description: alertRules[level].Description,
			Suggestion:  alertRules[level].Suggestion,
		})
	}
	fillTableCell(table, alarmTableColumns, rules)
}

func fillTableCell(table *tview.Table, columns []string, data []interface{}) {
	cols := len(columns)
	rows := len(data) + 1
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color, text := tcell.ColorYellow, columns[c]
			if r > 0 {
				color = tcell.ColorWhite
				value := reflect.ValueOf(data[r-1])
				if value.Kind() == reflect.Ptr {
					value = value.Elem()
				}
				fieldValue := value.Field(c)
				text = fieldValue.String()
			}
			cell := tview.NewTableCell(columns[c])
			cell.SetMaxWidth(_table_cell_max_width)
			cell.SetAlign(tview.AlignLeft)
			cell.SetText(text)
			cell.SetTextColor(color)
			table.SetCell(r, c, cell)
		}
	}
}

func yasdbValidate(form *tview.Form) (*yasdb.YashanDB, error) {
	res, err := getFormDataByLabels(form, constdef.YASDB_HOME, constdef.YASDB_DATA, constdef.YASDB_USER, constdef.YASDB_PASSWORD)
	if err != nil {
		log.Controller.Errorf("get form data err: %s", err.Error())
		return nil, err
	}
	yasdb := &yasdb.YashanDB{
		YasdbHome: res[constdef.YASDB_HOME],
		YasdbData: res[constdef.YASDB_DATA],
		YasdbUser: res[constdef.YASDB_USER],
	}
	std.WriteToFile("get yasdb info : \n")
	std.WriteToFile(jsonutil.ToJSONString(yasdb) + "\n")
	// after log fill password
	yasdb.YasdbPassword = res[constdef.YASDB_PASSWORD]
	if err := yasdb.ValidUserAndPwd(log.Controller); err != nil {
		return nil, err
	}
	if err := fillListenAddr(yasdb); err != nil {
		log.Controller.Errorf("fill listen addr err: %s", err.Error())
		return nil, err
	}

	return yasdb, nil

}

func metricValidate(env *yasdb.YashanDB, modules []*constdef.ModuleMetrics) {
	log := log.Controller.M("metric validate")
	for _, module := range modules {
		for _, metric := range module.Metrics {
			metricDefine := define.MetricName(metric.Name)
			if _, ok := check.NeedCheckMetricMap[metricDefine]; !ok {
				continue
			}
			if _, ok := check.NeedCheckMetricFuncMap[metricDefine]; !ok {
				log.Warnf("metric %s is defined in NeedCheckMetricMap, but NeedCheckMetricFuncMap is not defined", metric.Name)
				continue
			}
			if noNeedCheck := check.NeedCheckMetricFuncMap[metricDefine](log, env, metric); noNeedCheck != nil {
				if _, ok := moduleNoNeedCheckMetrics[module.Name]; !ok {
					moduleNoNeedCheckMetrics[module.Name] = make(map[string]*define.NoNeedCheckMetric)
				}
				moduleNoNeedCheckMetrics[module.Name][metric.Name] = noNeedCheck
			}
		}
	}
	globalFilterModule = filterNeedCheckMetric(modules)
}

func getFormData(form *tview.Form, label string) (string, error) {
	item := form.GetFormItemByLabel(label)
	if item == nil {
		return "", errdef.NewFormItemUnFound(label)
	}
	return item.(*tview.InputField).GetText(), nil
}

func getFormDataByLabels(form *tview.Form, labels ...string) (res map[string]string, err error) {
	res = make(map[string]string)
	for _, label := range labels {
		value, valueErr := getFormData(form, label)
		if valueErr != nil {
			err = valueErr
			return
		}
		res[label] = trimSpace(value)
	}
	return
}

func filterNeedCheckMetric(modules []*constdef.ModuleMetrics) (result []*constdef.ModuleMetrics) {
	result = make([]*constdef.ModuleMetrics, 0)
	for _, module := range modules {
		if _, ok := moduleNoNeedCheckMetrics[module.Name]; !ok {
			result = append(result, module)
			continue
		}
		metrics := make([]*confdef.YHCMetric, 0)
		for _, metric := range module.Metrics {
			if _, ok := moduleNoNeedCheckMetrics[module.Name][metric.Name]; ok {
				continue
			}
			metrics = append(metrics, metric)
			module.Metrics = metrics
		}
		if len(metrics) != 0 {
			result = append(result, &constdef.ModuleMetrics{
				Name:    module.Name,
				Metrics: metrics,
				Enabled: module.Enabled,
			})
		}
	}
	return result
}

func genNoNeedCheckMetricsStr() string {
	t := tabler.NewTable("",
		tabler.NewRowTitle("module", 25),
		tabler.NewRowTitle("metric", 25),
		tabler.NewRowTitle("description", 10),
		tabler.NewRowTitle("error", 10),
	)
	for _, module := range define.Level1ModuleOrder {
		moduleStr := string(module)
		if _, ok := moduleNoNeedCheckMetrics[moduleStr]; !ok {
			continue
		}
		moduleAlias, _ := define.GetModuleDefaultAlias(module)
		for _, metric := range moduleNoNeedCheckMetrics[moduleStr] {
			if err := t.AddColumn(moduleAlias, metric.Name, metric.Description, metric.Error.Error()); err != nil {
				log.Controller.Errorf("add columns err: %s", err.Error())
				continue
			}
		}
	}
	return t.String()
}

func fillListenAddr(db *yasdb.YashanDB) error {
	log := log.Controller.M("fill listen addr")
	tx := yasqlgo.NewLocalInstance(db.YasdbUser, db.YasdbPassword, db.YasdbHome, db.YasdbData, log)
	listenAddr, err := yhcyasdb.QueryParameter(tx, yhcyasdb.LISTEN_ADDR)
	if err != nil {
		return err
	}
	db.ListenAddr = trimSpace(listenAddr)
	return nil
}
