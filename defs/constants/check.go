package constdef

import "yhc/defs/confdef"

type ModuleMetrics struct {
	Name    string
	Metrics []*confdef.YHCMetric
	Enabled bool
}
