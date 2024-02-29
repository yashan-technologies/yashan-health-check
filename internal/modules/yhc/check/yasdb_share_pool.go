package check

import (
	"strconv"

	"yhc/defs/confdef"
	"yhc/internal/modules/yhc/check/define"
	"yhc/log"
	"yhc/utils/mathutil"

	"git.yasdb.com/go/yaserr"
	"git.yasdb.com/go/yasutil/size"
)

const (
	KEY_TOTAL_SIZE      = "TOTAL_SIZE"
	KYE_FREE_SIZE       = "FREE_SIZE"
	KEY_USED_SIZE       = "USED_SIZE"
	KEY_USED_PERCENTAGE = "USED_PERCENTAGE"

	COLUMN_NAME  = "NAME"
	COLUMN_BYTES = "BYTES"

	NAME_FREE_MEMORY = "free memory"
)

const decimal = 2

func (c *YHCChecker) GetYasdbSharePool(name string) (err error) {
	logger := log.Module.M(string(define.METRIC_YASDB_SHARE_POOL))

	var datas []*define.YHCItem
	for _, yasdb := range c.GetCheckNodes(logger) {
		data := &define.YHCItem{Name: define.METRIC_YASDB_SHARE_POOL, NodeID: yasdb.NodeID}
		var sharePoolData []map[string]string
		sharePoolData, err = yasdb.QueryMultiRows(define.SQL_QUERY_SHARE_POOL, confdef.GetYHCConf().SqlTimeout)
		if err != nil {
			err = yaserr.Wrap(err)
			logger.Error(err)
			data.Error = err.Error()
			continue
		}
		var totalBytes, freeBytes float64
		for _, row := range sharePoolData {
			bytes, e := strconv.ParseFloat(row[COLUMN_BYTES], 64)
			if err != nil {
				err = yaserr.Wrap(e)
				logger.Error(err)
				data.Error = err.Error()
				return
			}
			totalBytes += bytes
			if row[COLUMN_NAME] == NAME_FREE_MEMORY {
				freeBytes = bytes
			}
		}
		usedBytes := totalBytes - freeBytes
		usedPercentage := usedBytes / totalBytes * 100
		content := map[string]interface{}{
			KEY_TOTAL_SIZE:      size.GenHumanReadableSize(totalBytes, decimal),
			KYE_FREE_SIZE:       size.GenHumanReadableSize(freeBytes, decimal),
			KEY_USED_SIZE:       size.GenHumanReadableSize(usedBytes, decimal),
			KEY_USED_PERCENTAGE: mathutil.Round(usedPercentage, decimal),
		}
		data.Details = content
		datas = append(datas, data)
	}
	c.fillResults(datas...)
	return
}
