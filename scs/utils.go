package scs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// GetReaderLen 获取reader length
func GetReaderLen(reader io.Reader) (int64, error) {
	var contentLength int64
	var err error
	switch v := reader.(type) {
	case *bytes.Buffer:
		contentLength = int64(v.Len())
	case *bytes.Reader:
		contentLength = int64(v.Len())
	case *strings.Reader:
		contentLength = int64(v.Len())
	case *os.File:
		fInfo, fError := v.Stat()
		if fError != nil {
			err = fmt.Errorf("can't get reader content length,%s", fError.Error())
		} else {
			contentLength = fInfo.Size()
		}
	case *io.LimitedReader:
		contentLength = int64(v.N)
	default:
		err = fmt.Errorf("can't get reader content length,unkown reader type")
	}
	return contentLength, err
}
