package check

import "yhc/internal/modules/yhc/check/define"

func (c *YHCChecker) GetYasdbVMSwapRate() (err error) {
	data, err := c.querySingleRow(define.METRIC_YASDB_VM_SWAP_RATE)
	defer c.fillResult(data)
	return
}
