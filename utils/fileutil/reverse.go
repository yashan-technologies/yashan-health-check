package fileutil

import (
	"bytes"
	"io"
	"os"
)

// Reversefile manages a file opened for reading line-by-line in reverse.
type Reversefile struct {
	fh      *os.File
	offset  int64
	bufsize int
	lines   [][]byte
	i       int
}

// NewReverseFile returns Reversefile handle to be read in reverse line-by-line.
func NewReverseFile(file string) (*Reversefile, error) {
	var err error
	rf := new(Reversefile)
	if rf.fh, err = os.Open(file); err != nil {
		return nil, err
	}
	fi, _ := rf.fh.Stat()
	rf.offset = fi.Size()
	rf.bufsize = 4096
	return rf, nil
}

// Close file that was opened.
func (rf *Reversefile) Close() {
	rf.fh.Close()
}

// ReadLine returns the  next previous line, beginning with the last line in the file.
// When the beginning of the file is reached: "", io.EOF is returned.
func (rf *Reversefile) ReadLine() (string, error) {
	if rf.i > 0 {
		rf.i--
		return string(rf.lines[rf.i+1]), nil
	}
	if rf.i < 0 {
		return "", io.EOF
	}
	if rf.offset == 0 {
		rf.i-- // use as flag to send EOF on next call
		return string(rf.lines[0]), nil
	}

	// get more from file - back up from end-of-file
	rf.offset -= int64(rf.bufsize)
	if rf.offset < 0 {
		rf.bufsize += int(rf.offset) // rf.offset is negative
		rf.offset = 0
	}
	_, err := rf.fh.Seek(rf.offset, 0)
	if err != nil {
		return "", err
	}

	// size buffer
	buf := make([]byte, rf.bufsize)
	if n, err := rf.fh.Read(buf); err != nil && err != io.EOF {
		return "", err
	} else if n != rf.bufsize { // shouldn't happen
		buf = buf[:n]
	}

	// get the lines in the buffer, append what was carried over
	if len(rf.lines) > 0 {
		buf = append(buf, rf.lines[0]...)
	}
	// clean up terminating NL so we don't get empty slice from Split
	if buf[len(buf)-1] == '\n' {
		buf = buf[:len(buf)-1]
	}
	rf.lines = bytes.Split(buf, []byte("\n"))
	rf.i = len(rf.lines) - 1

	return rf.ReadLine() // now read the next line back
}

// Tail returns the last N lines of a file.  If the file has
// fewer lines than 'n', the whole file will be returned.
func Tail(file string, n int) ([]string, error) {
	lines := make([]string, n)

	rf, err := NewReverseFile(file)
	if err != nil {
		return nil, err
	}
	defer rf.Close()

	// Save the lines in reverse order.
	var i int
	for i = n - 1; i >= 0; i-- {
		lines[i], err = rf.ReadLine()
		if err == io.EOF {
			lines = lines[i+1:] // back it out
			break
		} else if err != nil {
			return lines[i:], err
		}
	}
	return lines, nil
}
