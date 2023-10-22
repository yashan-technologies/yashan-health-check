package alertgenner

import (
	"fmt"
	"strings"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/pandora/alertql"
	"git.yasdb.com/pandora/alertql/defs/metricdef"
)

type AlertGenner struct {
	log     yaslog.YasLog
	metrics []*confdef.YHCMetric
	result  map[define.MetricName]*define.YHCItem
}

func NewAlterGenner(log yaslog.YasLog, metrics []*confdef.YHCMetric, result map[define.MetricName]*define.YHCItem) *AlertGenner {
	genner := &AlertGenner{
		log:     log,
		metrics: metrics,
		result:  result,
	}
	return genner
}

func (a *AlertGenner) GenAlerts() map[define.MetricName]*define.YHCItem {
	data := a.genMetricsData()
	for _, metric := range a.metrics {
		item := a.result[define.MetricName(metric.Name)]
		for alertLevel, alertRules := range metric.AlertRules {
			for _, rule := range alertRules {
				expression := rule.Expression
				expr, err := alertql.NewExpression(expression, nil, true)
				if err != nil {
					a.log.Errorf("failed to gen alert exprersion by '%s', err: %v", expression, err)
					continue
				}
				alerts, err := expr.Execute(data)
				if err != nil {
					a.log.Errorf("failed to gen alert by expression: %s, err: %v", expression, err)
					continue
				}
				for _, alert := range alerts {
					yhcAlert := &define.YHCAlert{
						Level:        alertLevel,
						Value:        alert.Value,
						Labels:       alert.Labels,
						AlertDetails: rule,
					}
					if item.Alerts == nil {
						item.Alerts = make(map[string][]*define.YHCAlert)
					}
					item.Alerts[yhcAlert.Level] = append(item.Alerts[yhcAlert.Level], yhcAlert)
				}
			}
		}
	}
	return a.result
}

func (a *AlertGenner) genMetricsData() map[string]interface{} {
	pool := metricdef.MetricsPool{}
	for _, metric := range a.metrics {
		metricName := define.MetricName(metric.Name)
		a.dealItem(&pool, metric, a.result[metricName])
	}
	data := make(map[string]interface{})
	for k, v := range pool {
		data[k] = v
	}
	return data
}

func (a *AlertGenner) dealItem(pool *metricdef.MetricsPool, metric *confdef.YHCMetric, item *define.YHCItem) {
	if item == nil {
		return
	}
	details := item.Details
	switch detailsType := details.(type) {
	case []string:
		a.log.Debugf("unsupport alert type []string, skip")
	case []interface{}:
		a.log.Debugf("unsupport alert type []interface{}, skip")
	case map[string]string:
		data := details.(map[string]string)
		a.dealSingleStringRow(pool, metric, data)
	case []map[string]string:
		datas := details.([]map[string]string)
		for _, data := range datas {
			a.dealSingleStringRow(pool, metric, data)
		}
	case map[string]interface{}:
		data := details.(map[string]interface{})
		a.dealSingleAnyRow(pool, metric, data)
	case []map[string]any:
		datas := details.([]map[string]interface{})
		for _, data := range datas {
			a.dealSingleAnyRow(pool, metric, data)
		}
	default:
		a.log.Errorf("unsupport data type %T", detailsType)
	}
}

func (a *AlertGenner) dealSingleStringRow(pool *metricdef.MetricsPool, metric *confdef.YHCMetric, data map[string]string) {
	for key, value := range data {
		subMetricName, ok := metric.ItemNames[key]
		if !ok {
			subMetricName = fmt.Sprintf("%s_%s", metric.Name, strings.ToLower(key))
		}
		metrics, ok := (*pool)[subMetricName]
		if !ok {
			metrics = []metricdef.Metric{}
		}
		labelsMap := make(map[string]string)
		for _, label := range metric.Labels {
			labelsMap[label] = data[label]
		}
		metrics = append(metrics, metricdef.Metric{
			Value:  value,
			Labels: labelsMap,
		})
		(*pool)[subMetricName] = metrics
	}
}

func (a *AlertGenner) dealSingleAnyRow(pool *metricdef.MetricsPool, metric *confdef.YHCMetric, data map[string]interface{}) {
	for key, value := range data {
		subMetricName, ok := metric.ItemNames[key]
		if !ok {
			subMetricName = fmt.Sprintf("%s_%s", metric.Name, strings.ToLower(key))
		}
		metrics, ok := (*pool)[subMetricName]
		if !ok {
			metrics = []metricdef.Metric{}
		}
		labelsMap := make(map[string]string)
		for _, label := range metric.Labels {
			v, ok := data[label].(string)
			if !ok {
				a.log.Warnf("unsupport label type %T", value)
				continue
			}
			labelsMap[label] = v
		}
		metrics = append(metrics, metricdef.Metric{
			Value:  value,
			Labels: labelsMap,
		})
		(*pool)[subMetricName] = metrics
	}
}
