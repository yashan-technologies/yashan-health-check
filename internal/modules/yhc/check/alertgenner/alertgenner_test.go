package alertgenner_test

import (
	"encoding/json"
	"os"
	"testing"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/alertgenner"
	"yhc/internal/modules/yhc/check/define"

	"git.yasdb.com/go/yaslog"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

var (
	metricResultStr = `{"yasdb_table_with_row_size_exceeds_block_size": {
    "details": [
        {
            "OWNER": "TEST",
            "TABLE_NAME": "T14"
        }
    ],
    "alerts": {
        "waring": [
            {
                "level": "waring",
                "value": "T14",
                "labels": {
                    "OWNER": "TEST",
                    "TABLE_NAME": "T14"
                },
                "Expression": "name_of_table_with_row_size_exceeds_block_size != ''",
                "Description": "符合条件的表的名称",
                "Suggestion": "当表的行大小大于数据库或表空间的块大小时，则每行需要两个IOs，建议使用更大的数据块来存储该表。您可能需要增加数据库块大小(重组练习)或将表移动到块大小更大的表空间(重新定位)。"
            }
        ]
    }
}}`

	detailStr = `[{
        "OWNER": "TEST",
        "TABLE_NAME": "T14"
    }]`

	metricsStr = `
[[metrics]]
  name = "yasdb_table_with_row_size_exceeds_block_size"
  name_alias = "行大小超过块大小的表"
  module_name = "object_check"
  default = true
  enabled = true
  sql = '''SELECT a.OWNER, a.TABLE_NAME
    FROM (
        SELECT OWNER, TABLE_NAME, SUM(DATA_LENGTH) AS MAX_DL
        FROM DBA_TAB_COLUMNS
        WHERE  OWNER <> 'SYS' AND DATA_TYPE NOT LIKE '%LOB'
        GROUP BY OWNER, TABLE_NAME
    ) a, (
        SELECT to_number(decode(value, '8K','8192','16K','16384','32K','32768',value)) as VALUE 
        FROM v$parameter 
        WHERE NAME = 'DB_BLOCK_SIZE'
        ) b
    WHERE a.max_dl > b.value;'''
  labels = ["OWNER", "TABLE_NAME"]
  [metrics.column_alias]
    OWNER = "表所属用户"
    TABLE_NAME = "表名称"

  [metrics.item_names]
    TABLE_NAME = "name_of_table_with_row_size_exceeds_block_size"

  [metrics.alert_rules]

    [[metrics.alert_rules.waring]]
      expression = "name_of_table_with_row_size_exceeds_block_size != ''"
      description = "符合条件的表的名称"
      suggestion = "当表的行大小大于数据库或表空间的块大小时，则每行需要两个IOs，建议使用更大的数据块来存储该表。您可能需要增加数据库块大小(重组练习)或将表移动到块大小更大的表空间(重新定位)。"
`
	logger                                                   = yaslog.NewDefaultConsoleLogger()
	metricResult       map[define.MetricName]*define.YHCItem = make(map[define.MetricName]*define.YHCItem)
	metricsConf        confdef.YHCMetricConfig               = confdef.YHCMetricConfig{}
	exceedsBlockDetail []map[string]string                   = make([]map[string]string, 0)
)

func preTest() {
	_, _ = toml.Decode(metricsStr, &metricsConf)
	_ = json.Unmarshal([]byte(metricResultStr), &metricResult)
	_ = json.Unmarshal([]byte(detailStr), &exceedsBlockDetail)
	metricResult[define.MetricName("yasdb_table_with_row_size_exceeds_block_size")].Details = exceedsBlockDetail
}

func TestMain(m *testing.M) {
	preTest()
	os.Exit(m.Run())
}

func TestGenAlerts(t *testing.T) {
	genner := alertgenner.NewAlterGenner(logger, metricsConf.Metrics, metricResult)
	genner.GenAlerts()
	assert.NotEqual(t, 0, metricResult[define.MetricName("yasdb_table_with_row_size_exceeds_block_size")].Alerts)
}
