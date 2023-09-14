package jsonutil

import "encoding/json"

func ToJSONString(any interface{}) string {
	bytes, _ := json.MarshalIndent(any, "", "    ")
	return string(bytes)
}
