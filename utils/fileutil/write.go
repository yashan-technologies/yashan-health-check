package fileutil

import (
	"io/fs"
	"os"
)

const (
	DEFAULT_FILE_MODE fs.FileMode = 0644
)

func WriteFile(fname string, data []byte) error {
	return os.WriteFile(fname, data, DEFAULT_FILE_MODE)
}

func RewriteFile(str string, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}
