package processutil_test

import (
	"encoding/json"
	"testing"

	"yhc/utils/processutil"
)

func TestIsRunning(t *testing.T) {
	p := processutil.NewProcess(1)
	_, ok := p.IsRunning()
	if !ok {
		t.Fail()
	}
}

func TestListProcess(t *testing.T) {
	processes, err := processutil.ListProcess()
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(processes, "", "  ")
	t.Log(len(processes))
	t.Log(string(data))
}
