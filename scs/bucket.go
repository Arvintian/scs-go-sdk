package scs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Arvintian/scs-go-sdk/pkg/client"
)

func (b *Bucket) String() string {
	return b.Name
}

/*********************************************
GetObject和MultipartUpload没实现数据完整性校验
**********************************************/

// Head 获取object meta
func (b *Bucket) Head(key string) (ObjectMeta, error) {
	var m ObjectMeta
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	req := &client.Request{
		Method: "HEAD",
		Bucket: b.Name,
		Path:   fmt.Sprintf("/%s", key),
		Params: params,
	}
	header, _, err := b.c.Query(req)
	if err != nil {
		return m, err
	}
	//fmt.Println(header)
	m.ContentType = header.Get("Content-Type")
	length, err := strconv.ParseInt(header.Get("Content-Length"), 10, 0)
	if err != nil {
		length, err = strconv.ParseInt(header.Get("X-Filesize"), 10, 0)
		if err != nil {
			return m, err
		}
	}
	m.ContentLength = length
	m.ETag = header.Get("ETag")
	m.LastModified = header.Get("Last-Modified")
	m.XAmzMeta = make(map[string]string)
	for k, v := range header {
		if strings.HasPrefix(k, "X-Amz-Meta-") {
			if len(v) > 0 {
				m.XAmzMeta[k] = v[0]
			} else {
				m.XAmzMeta[k] = ""
			}
		} else {
			if strings.HasPrefix(k, "x-amz-meta-") {
				if len(v) > 0 {
					m.XAmzMeta[k] = v[0]
				} else {
					m.XAmzMeta[k] = ""
				}
			}
		}
	}
	return m, nil
}

// Get 获取object
func (b *Bucket) Get(key string, rg string) (io.ReadCloser, error) {
	var params = make(map[string][]string)
	var headers = make(http.Header)
	params["formatter"] = []string{"json"}
	if rg != "" {
		headers.Set("Range", fmt.Sprintf("bytes=%s", rg))
	}
	req := &client.Request{
		Method:  "GET",
		Bucket:  b.Name,
		Path:    fmt.Sprintf("/%s", key),
		Params:  params,
		Headers: headers,
	}
	_, data, err := b.c.Query(req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Put 上传object
func (b *Bucket) Put(key string, XAmzMeta map[string]string, data io.Reader) error {
	var params = make(map[string][]string)
	var headers = make(http.Header)
	params["formatter"] = []string{"json"}
	for k, v := range XAmzMeta {
		headers.Set(k, v)
	}
	length, err := GetReaderLen(data)
	if err != nil {
		return err
	}
	putData, md5, fd, err := calcMD5(data, length)
	if fd != nil {
		defer func() {
			fd.Close()
			os.Remove(fd.Name())
		}()
	}
	if err != nil {
		return err
	}
	//fmt.Println(md5)
	headers.Set("Content-Length", fmt.Sprint(length))
	headers.Set("Content-MD5", md5)
	req := &client.Request{
		Method:  "PUT",
		Bucket:  b.Name,
		Path:    fmt.Sprintf("/%s", key),
		Params:  params,
		Headers: headers,
		Body:    putData,
	}
	_, _, err = b.c.Query(req)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除object
func (b *Bucket) Delete(key string) error {
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	req := &client.Request{
		Method: "DELETE",
		Bucket: b.Name,
		Path:   fmt.Sprintf("/%s", key),
		Params: params,
	}
	_, _, err := b.c.Query(req)
	if err != nil {
		return err
	}
	return err
}

// List 获取obejct列表
func (b *Bucket) List(delimiter, prefix, marker string, limit int64) (ListObject, error) {
	var lo ListObject
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	if delimiter != "" {
		params["delimiter"] = []string{delimiter}
	}
	if prefix != "" {
		params["prefix"] = []string{prefix}
	}
	if marker != "" {
		params["marker"] = []string{marker}
	}
	params["max-keys"] = []string{fmt.Sprint(limit)}
	req := &client.Request{
		Method: "GET",
		Bucket: b.Name,
		Path:   "/",
		Params: params,
	}
	_, body, err := b.c.Query(req)
	if err != nil {
		return lo, err
	}
	bts, err := ioutil.ReadAll(body)
	if err != nil {
		return lo, err
	}
	//fmt.Println(string(bts))
	if err := json.Unmarshal(bts, &lo); err != nil {
		return lo, err
	}
	return lo, nil
}

// InitiateMultipartUpload 大文件分片上传
func (b *Bucket) InitiateMultipartUpload(key string, XAmzMeta map[string]string) (MultipartUpload, error) {
	var mu MultipartUpload
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	params["multipart"] = []string{""}
	var headers = make(http.Header)
	for k, v := range XAmzMeta {
		headers.Set(k, v)
	}
	req := &client.Request{
		Method:  "POST",
		Bucket:  b.Name,
		Path:    fmt.Sprintf("/%s", key),
		Params:  params,
		Headers: headers,
	}
	_, body, err := b.c.Query(req)
	if err != nil {
		return mu, err
	}
	bts, err := ioutil.ReadAll(body)
	if err != nil {
		return mu, err
	}
	if err := json.Unmarshal(bts, &mu); err != nil {
		return mu, err
	}
	return mu, nil
}

// UploadPart 上传分片
func (b *Bucket) UploadPart(key string, uploadID string, partNumber int, data io.Reader) (Part, error) {
	var p Part
	p.PartNumber = partNumber
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	params["partNumber"] = []string{fmt.Sprint(partNumber)}
	params["uploadId"] = []string{uploadID}
	var headers = make(http.Header)
	length, err := GetReaderLen(data)
	if err != nil {
		return p, err
	}
	headers.Set("Content-Length", fmt.Sprint(length))
	putData, md5, fd, err := calcMD5(data, length)
	if fd != nil {
		defer func() {
			fd.Close()
			os.Remove(fd.Name())
		}()
	}
	if err != nil {
		return p, err
	}
	headers.Set("Content-MD5", md5)
	req := &client.Request{
		Method:  "PUT",
		Bucket:  b.Name,
		Path:    fmt.Sprintf("/%s", key),
		Params:  params,
		Headers: headers,
		Body:    putData,
	}
	rspHeaders, _, err := b.c.Query(req)
	if err != nil {
		return p, err
	}
	p.Size = int(length)
	p.ETag = rspHeaders.Get("Etag")
	return p, nil
}

// CompleteMultipartUpload 完成分片上传
func (b *Bucket) CompleteMultipartUpload(key, uploadID string, parts []Part) error {
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	params["uploadId"] = []string{uploadID}
	bts, err := json.Marshal(parts)
	if err != nil {
		return err
	}
	req := &client.Request{
		Method: "POST",
		Bucket: b.Name,
		Path:   fmt.Sprintf("/%s", key),
		Params: params,
		Body:   bytes.NewBuffer(bts),
	}
	_, _, err = b.c.Query(req)
	if err != nil {
		return err
	}
	return nil
}

// ListParts 列出已经上传的所有分块
func (b *Bucket) ListParts(key, uploadID string) (ListPart, error) {
	var lp ListPart
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	params["uploadId"] = []string{uploadID}
	req := &client.Request{
		Method: "GET",
		Bucket: b.Name,
		Path:   fmt.Sprintf("/%s", key),
		Params: params,
	}
	_, body, err := b.c.Query(req)
	if err != nil {
		return lp, err
	}
	bts, err := ioutil.ReadAll(body)
	//fmt.Println(string(bts))
	if err := json.Unmarshal(bts, &lp); err != nil {
		return lp, err
	}
	return lp, nil
}
