package scs

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Arvintian/scs-go-sdk/pkg/client"
)

func (b *Bucket) String() string {
	return b.Name
}

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
func (b *Bucket) Get(key string, start, end int64) (io.ReadCloser, error) {
	var params = make(map[string][]string)
	var headers = make(http.Header)
	params["formatter"] = []string{"json"}
	if !(start == 0 && end == 0) {
		headers.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
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
	headers.Set("Content-Length", fmt.Sprint(length))
	req := &client.Request{
		Method:  "PUT",
		Bucket:  b.Name,
		Path:    fmt.Sprintf("/%s", key),
		Params:  params,
		Headers: headers,
		Body:    data,
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
