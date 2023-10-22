package stdutil

import (
	"fmt"
	"io"
	"os"

	"yhc/defs/runtimedef"
	"yhc/utils/userutil"
)

const (
	_newline = "\n"
)

type Redirecter struct {
	fileWriter *os.File
	fName      string
}

func NewRedirecter(fn string) (*Redirecter, error) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return nil, err
	}
	return &Redirecter{
		fileWriter: f,
		fName:      fn,
	}, nil
}

// RedirectStd writes stdout and stderr to a file and retains the original output.
func (r *Redirecter) RedirectStd() (finalize func()) {
	out := os.Stdout
	mw := io.MultiWriter(out, r.fileWriter)
	reader, writer, err := os.Pipe()
	if err != nil {
		return
	}
	os.Stdout = writer
	os.Stderr = writer

	done := make(chan bool)
	go func() {
		// io.Copy() returns when the reader gets EOF
		_, _ = io.Copy(mw, reader)
		done <- true
	}()

	finalize = func() {
		// use w.Close() to make the reader get EOF, which causes io.Copy() return
		// so that "done <- true" will be executed
		_ = writer.Close()
		// wait for goroutine done
		<-done
		// wait for flushing the file system's in-memory copy of recently written data to disk
		_ = r.fileWriter.Sync()
		// close file
		_ = r.fileWriter.Close()
		// change owner
		owner := runtimedef.GetExecuteableOwner()
		if userutil.IsCurrentUserRoot() && (owner.Uid != 0 || owner.Gid != 0) {
			_ = os.Chown(r.fName, owner.Uid, owner.Gid)
		}
	}
	return
}

func (r *Redirecter) GetFileWriter() *os.File {
	return r.fileWriter
}

// ReadFromStdin scans input from stdin, and uses extraWriters to record the value.
func ReadFromStdin(extraWriters ...*os.File) (str string) {
	fmt.Scanln(&str)
	var ioWriters []io.Writer
	for _, w := range extraWriters {
		ioWriters = append(ioWriters, w)
	}
	mutiWriter := io.MultiWriter(ioWriters...)
	_, _ = mutiWriter.Write([]byte(str + _newline))
	return str
}

// Write uses extraWriters to record the value.
func Write(str string, extraWriters ...*os.File) {
	var ioWriters []io.Writer
	for _, w := range extraWriters {
		ioWriters = append(ioWriters, w)
	}
	mutiWriter := io.MultiWriter(ioWriters...)
	_, _ = mutiWriter.Write([]byte(str))
}

// WriteToStdout writes str to stdout, and uses extraWriters to record the value.
func WriteToStdout(str string, extraWriters ...*os.File) {
	extraWriters = append(extraWriters, os.Stdout)
	Write(str, extraWriters...)
}
