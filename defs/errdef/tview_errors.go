package errdef

import (
	"errors"
	"fmt"
)

var (
	ErrPermission    = errors.New("some permission")
	ErrExitWithCtrlC = errors.New("exit with control c")
)

type FormItemUnFound struct {
	ItemName string
}

func NewFormItemUnFound(itemName string) *FormItemUnFound {
	return &FormItemUnFound{
		ItemName: itemName,
	}
}

func (e *FormItemUnFound) Error() string {
	return fmt.Sprintf("form item %s unfound", e.ItemName)
}

type ItemEmpty struct {
	ItemName string
}

func NewItemEmpty(itemName string) *ItemEmpty {
	return &ItemEmpty{
		ItemName: itemName,
	}
}

func (e *ItemEmpty) Error() string {
	return fmt.Sprintf("%s is empty", e.ItemName)
}
