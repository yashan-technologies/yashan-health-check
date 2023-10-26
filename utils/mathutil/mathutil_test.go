package mathutil_test

import (
	"fmt"
	"testing"

	"yhc/utils/mathutil"
)

func TestRound(t *testing.T) {
	res := mathutil.Round(3.141592, 2)
	fmt.Println(res)
}

func TestNumber(t *testing.T) {
	res := mathutil.GenHumanReadableNumber(10000000, 2)
	fmt.Println(res)
}
