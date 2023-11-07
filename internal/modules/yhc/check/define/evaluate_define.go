package define

import "yhc/defs/confdef"

type EvaluateResult struct {
	EvaluateModel *confdef.EvaluateModel `json:"evaluateModel"`
	Score         float64                `json:"score"`
	HealthStatus  string                 `json:"healthStatus"`
	AlertSummary  *AlertSummary          `json:"alertSummary"`
}

type AlertSummary struct {
	InfoCount     int `json:"infoCount"`
	WarningCount  int `json:"warningCount"`
	CriticalCount int `json:"criticalCount"`
}
