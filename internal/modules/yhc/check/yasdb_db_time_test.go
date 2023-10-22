package check

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"

	"yhc/commons/constants"

	"git.yasdb.com/go/yaslog"
	"github.com/stretchr/testify/assert"
)

var m = []map[string]string{
	{"DB_TIMES": "4644180", "SNAP_TIME": "2023-10-15 18:49:03"},
	{"DB_TIMES": "4651365", "SNAP_TIME": "2023-10-15 19:49:04"},
	{"DB_TIMES": "4658609", "SNAP_TIME": "2023-10-15 20:49:05"},
	{"DB_TIMES": "4", "SNAP_TIME": "2023-10-16 03:45:52"},
	{"DB_TIMES": "3249", "SNAP_TIME": "2023-10-16 10:42:39"},
	{"DB_TIMES": "4059", "SNAP_TIME": "2023-10-16 11:42:40"},
	{"DB_TIMES": "4097", "SNAP_TIME": "2023-10-16 12:42:42"},
	{"DB_TIMES": "4123", "SNAP_TIME": "2023-10-16 13:42:43"},
	{"DB_TIMES": "4155", "SNAP_TIME": "2023-10-16 14:42:44"},
	{"DB_TIMES": "4187", "SNAP_TIME": "2023-10-16 15:42:45"},
}

func TestGetDBTimes(t *testing.T) {
	checker := NewYHCChecker(nil, nil)
	res, err := checker.getDBTimes(yaslog.NewDefaultConsoleLogger(), m)
	if err != nil {
		log.Fatal(err)
	}
	for _, dbTimes := range res {
		tmp := map[string]map[string]interface{}{}
		bytes, _ := json.Marshal(dbTimes)
		_ = json.Unmarshal(bytes, &tmp)
		dbTimeStr := tmp[KEY_DB_TIME_MS][KEY_DB_TIMES].(string)
		dbTime, _ := strconv.ParseFloat(dbTimeStr, constants.BIT_SIZE_64)
		assert.Greater(t, dbTime, float64(0))
	}
}
