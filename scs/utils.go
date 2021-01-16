package scs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
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

func calcMD5(body io.Reader, contentLen int64) (reader io.Reader, b64 string, tempFile *os.File, err error) {
	if contentLen == 0 || contentLen > MD5Threshold {
		// Huge body, use temporary file
		tempFile, err = ioutil.TempFile(os.TempDir(), TempFilePrefix)
		if tempFile != nil {
			io.Copy(tempFile, body)
			tempFile.Seek(0, os.SEEK_SET)
			md5 := md5.New()
			io.Copy(md5, tempFile)
			sum := md5.Sum(nil)
			b64 = base64.StdEncoding.EncodeToString(sum[:])
			tempFile.Seek(0, os.SEEK_SET)
			reader = tempFile
		}
	} else {
		// Small body, use memory
		var buf []byte
		buf, err = ioutil.ReadAll(body)
		sum := md5.Sum(buf)
		b64 = base64.StdEncoding.EncodeToString(sum[:])
		reader = bytes.NewReader(buf)
	}
	return
}

func calcCRC32(body io.Reader, contentLen int64) (reader io.Reader, crc string, tempFile *os.File, err error) {
	if contentLen == 0 || contentLen > MD5Threshold {
		// Huge body, use temporary file
		tempFile, err = ioutil.TempFile(os.TempDir(), TempFilePrefix)
		if tempFile != nil {
			io.Copy(tempFile, body)
			tempFile.Seek(0, os.SEEK_SET)
			mycrc32 := crc32.NewIEEE()
			io.Copy(mycrc32, tempFile)
			crc = strings.ToUpper(fmt.Sprintf("%x", mycrc32.Sum32()))
			tempFile.Seek(0, os.SEEK_SET)
			reader = tempFile
		}
	} else {
		// Small body, use memory
		var buf []byte
		buf, err = ioutil.ReadAll(body)
		mycrc32 := crc32.NewIEEE()
		mycrc32.Write(buf)
		crc = strings.ToUpper(fmt.Sprintf("%x", mycrc32.Sum32()))
		reader = bytes.NewReader(buf)
	}
	return
}
