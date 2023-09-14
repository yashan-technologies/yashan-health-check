package errdef

import (
	"errors"
	"fmt"
)

var (
	ErrEndLessStart        = errors.New("start time should be less than end time")
	ErrStartShouldLessCurr = errors.New("start time should be less than current time")
)

type ErrGreaterMaxDuration struct {
	MaxDuration string
}

func NewGreaterMaxDur(max string) *ErrGreaterMaxDuration {
	return &ErrGreaterMaxDuration{MaxDuration: max}
}

func (e ErrGreaterMaxDuration) Error() string {
	return fmt.Sprintf("end-start time should be less than %s, you can modify the configuration file ./config/strategy.toml 'max_duration'", e.MaxDuration)
}

type ErrLessMinDuration struct {
	MinDuration string
}

func NewLessMinDur(min string) *ErrLessMinDuration {
	return &ErrLessMinDuration{MinDuration: min}
}

func (e ErrLessMinDuration) Error() string {
	return fmt.Sprintf("end-start time should be greater than %s, you can modify the configuration file ./config/strategy.toml 'min_duration'", e.MinDuration)
}
