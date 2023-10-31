package mathutil

import (
	"fmt"
	"math"

	"git.yasdb.com/go/yasutil"
)

const (
	thousand        = 1000
	ten_thounsand   = 10 * thousand
	million         = 100 * ten_thounsand
	ten_million     = 10 * million
	hundred_million = 10 * ten_million
)

func Round(num float64, decimal int) float64 {
	pow := math.Pow(10, float64(decimal))
	return math.Round(num*pow) / pow
}

func GenHumanReadableNumber(num float64, decimal int) string {
	if num == 0 || num < thousand {
		return yasutil.FormatFloat(num, decimal)
	}
	if num < ten_thounsand {
		return fmt.Sprintf("%s千", yasutil.FormatFloat(num/thousand, decimal))
	}
	if num < million {
		return fmt.Sprintf("%s万", yasutil.FormatFloat(num/ten_thounsand, decimal))
	}
	if num < ten_million {
		return fmt.Sprintf("%s百万", yasutil.FormatFloat(num/million, decimal))
	}
	if num < hundred_million {
		return fmt.Sprintf("%s千万", yasutil.FormatFloat(num/ten_million, decimal))
	}
	return fmt.Sprintf("%s亿", yasutil.FormatFloat(num/ten_million, decimal))
}
