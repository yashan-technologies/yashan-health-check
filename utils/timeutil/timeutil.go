package timeutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"yhc/defs/regexpdef"
)

const (
	YEAR_INDEX = iota
	MONTH_INDEX
	DAY_INDEX
	HOUR_INDEX
	MINUTE_INDEX
	SECOND_INDEX
)

const (
	minute = "m"
	hour   = "h"
	day    = "d"
	month  = "M"

	minute_dur time.Duration = time.Minute
	hour_dur                 = time.Hour
	day_dur                  = hour_dur * 24
	month_dur                = day_dur * 30
	year_dur                 = month_dur * 12

	NOON = time.Hour * 12
)

var (
	ErrDurationInvalid = errors.New("duration invalid")
)

var (
	monthMap = map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"May": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Aug": time.August,
		"Sep": time.September,
		"Oct": time.October,
		"Nov": time.November,
		"Dec": time.December,
	}

	durMap = map[string]time.Duration{
		minute: minute_dur,
		hour:   hour_dur,
		day:    day_dur,
		month:  month_dur,
	}
)

func GetTimeDivBySepa(timeStr, sepa string) (time.Time, error) {
	dateFields := [6]int{}
	parts := strings.Split(timeStr, sepa)
	for index := range parts {
		field, err := strconv.ParseInt(parts[index], 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("prase time err: %s", err.Error())
		}
		dateFields[index] = int(field)
	}
	t := time.Date(
		dateFields[YEAR_INDEX],
		time.Month(dateFields[MONTH_INDEX]),
		dateFields[DAY_INDEX],
		dateFields[HOUR_INDEX],
		dateFields[MINUTE_INDEX],
		dateFields[SECOND_INDEX],
		0, time.Local,
	)
	return t, nil
}

func GetDuration(s string) (d time.Duration, err error) {
	if !regexpdef.RangeRegex.MatchString(s) {
		err = ErrDurationInvalid
		return
	}
	var p int64
	suffix := s[len(s)-1:]
	prefix := s[:len(s)-1]
	p, err = strconv.ParseInt(prefix, 10, 64)
	if err != nil {
		return
	}
	d = duration(p, durMap[suffix])
	return
}

func duration(base int64, unit time.Duration) time.Duration {
	return time.Duration(base * int64(unit))
}

func GetMonth(str string) (month time.Month, err error) {
	month, ok := monthMap[str]
	if !ok {
		err = fmt.Errorf("invalid month %s", str)
	}
	return
}
