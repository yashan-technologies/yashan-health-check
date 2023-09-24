package std

import (
	"path"

	"yhc/defs/runtimedef"
	"yhc/utils/stdutil"
)

const (
	_fname = "console.out"
)

var _redirecter *stdutil.Redirecter

func InitRedirecter() error {
	redirecter, err := stdutil.NewRedirecter(genOutput())
	if err != nil {
		return err
	}
	_redirecter = redirecter
	return nil
}

func GetRedirecter() *stdutil.Redirecter {
	return _redirecter
}

func WriteToFile(str string) {
	stdutil.Write(str, _redirecter.GetFileWriter())
}

func WriteToFileAndStdout(str string) {
	stdutil.WriteToStdout(str, _redirecter.GetFileWriter())
}

func genOutput() string {
	return path.Join(runtimedef.GetLogPath(), _fname)
}
