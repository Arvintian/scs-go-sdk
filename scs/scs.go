package scs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Arvintian/scs-go-sdk/pkg/client"
)

//NewSCS 获取scs实例
func NewSCS(accesskey, secretkey, endpoint string) (*SCS, error) {
	return &SCS{
		c: client.NewClient(accesskey, secretkey, endpoint),
	}, nil
}

//ListBuckets 获取所有buckets
func (s *SCS) ListBuckets() ([]Bucket, error) {
	var bs BucketList
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	req := &client.Request{
		Method: "GET",
		Path:   "/",
		Params: params,
	}
	_, rc, err := s.c.Query(req)
	defer rc.Close()
	if err != nil {
		return bs.Buckets, err
	}
	bts, err := ioutil.ReadAll(rc)
	if err != nil {
		return bs.Buckets, err
	}
	if err := json.Unmarshal(bts, &bs); err != nil {
		return bs.Buckets, err
	}
	result := make([]Bucket, len(bs.Buckets))
	for i, b := range bs.Buckets {
		b.c = s.c
		result[i] = b
	}
	return result, nil
}

//GetBucket 获取bucket实例
func (s *SCS) GetBucket(name string) (Bucket, error) {
	var b Bucket
	bl, err := s.ListBuckets()
	if err != nil {
		return b, err
	}
	for _, v := range bl {
		if v.Name == name {
			return v, nil
		}
	}
	return b, fmt.Errorf("not found bucket %s", name)
}

//GetBucketMeta 获取bucket meta
func (s *SCS) GetBucketMeta(name string) (BucketMeta, error) {
	var meta BucketMeta
	var params = make(map[string][]string)
	params["meta"] = []string{""}
	params["formatter"] = []string{"json"}
	req := &client.Request{
		Method: "GET",
		Bucket: name,
		Path:   "/",
		Params: params,
	}
	_, rc, err := s.c.Query(req)
	defer rc.Close()
	if err != nil {
		return meta, err
	}
	bts, err := ioutil.ReadAll(rc)
	if err != nil {
		return meta, err
	}
	if err := json.Unmarshal(bts, &meta); err != nil {
		return meta, err
	}
	return meta, nil
}

//PutBucket 创建bucket
func (s *SCS) PutBucket(name string, acl string) error {
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	var headers = make(http.Header)
	if acl != ACLPrivate && acl != ACLPublicRead && acl != ACLPublicReadWrite && acl != ACLAuthenticatedRead {
		return errors.New("acl error")
	}
	headers.Set("x-amz-acl", acl)
	req := &client.Request{
		Method:  "PUT",
		Bucket:  name,
		Path:    "/",
		Params:  params,
		Headers: headers,
	}
	_, body, err := s.c.Query(req)
	defer body.Close()
	if err != nil {
		return err
	}
	return nil
}

//DeleteBucket 删除bucket
func (s *SCS) DeleteBucket(name string) error {
	var params = make(map[string][]string)
	params["formatter"] = []string{"json"}
	req := &client.Request{
		Method: "DELETE",
		Bucket: name,
		Path:   "/",
		Params: params,
	}
	_, body, err := s.c.Query(req)
	defer body.Close()
	if err != nil {
		return err
	}
	return nil
}

//GetBucketACL 获取bucket acl
func (s *SCS) GetBucketACL(name string) error {
	return errors.New("no implement")
}

//PutBucketACL 指定Bucket设置ACL规则
func (s *SCS) PutBucketACL(name string) error {
	return errors.New("no implement")
}
