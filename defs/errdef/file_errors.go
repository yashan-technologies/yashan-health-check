package errdef

import (
	"errors"
	"fmt"
)

var (
	ErrPathFormat = errors.New("path format error, please check")
)

type ErrPermissionDenied struct {
	User     string
	FileName string
}

type ErrFileNotFound struct {
	FName string
}

type ErrFileParseFailed struct {
	FName string
	Err   error
}

func NewErrPermissionDenied(user string, path string) *ErrPermissionDenied {
	return &ErrPermissionDenied{
		User:     user,
		FileName: path,
	}
}

func (e *ErrPermissionDenied) Error() string {
	return fmt.Sprintf("The current user %s does not have permission to: %s", e.User, e.FileName)
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("%s is not existed", e.FName)
}

func (e *ErrFileParseFailed) Error() string {
	return fmt.Sprintf("parse %s failed: %s", e.FName, e.Err)
}
