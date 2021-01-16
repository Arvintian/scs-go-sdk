package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/Arvintian/scs-go-sdk/scs"
)

// {
//     "accesskey": "accesskey",
//     "secretkey": "secretkey",
//     "bucket": "demo",
//     "endpoint": "https://sinacloud.net"
// }

func main() {
	bts, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("load config error")
	}
	config := make(map[string]string)
	err = json.Unmarshal(bts, &config)
	if err != nil {
		log.Fatal("load config error")
	}
	s, err := scs.NewSCS(config["accesskey"], config["secretkey"], config["endpoint"])
	if err != nil {
		log.Fatal(err)
	}
	bucket := config["bucket"]
	fmt.Println("test get bucket meta========")
	fmt.Println(s.GetBucketMeta(bucket))

	fmt.Println("test list get bucket ============")
	fmt.Println(s.ListBuckets())
	b, _ := s.GetBucket(bucket)
	// fmt.Println(s.PutBucket("test"))
	// fmt.Println(s.GetBucketACL("test"))
	// fmt.Println(s.PutBucketACL("test"))

	fmt.Println("test put==========")
	data := bytes.NewBufferString("sssssbb")
	fmt.Println(b.Put("testkey", map[string]string{
		"x-amz-meta-test": "foo",
	}, data))

	fmt.Println("test head object================")
	fmt.Println(b.Head("testkey"))

	fmt.Println("test get object==========")
	rdata, _ := b.Get("testkey", 0, 0)
	readPrint(rdata)
	rdata, _ = b.Get("testkey", 5, 7)
	readPrint(rdata)
	rdata, _ = b.Get("testkey", 5, 10)
	readPrint(rdata)

	fmt.Println("test delete object==============")
	fmt.Println(b.Delete("testkey"))
	fmt.Println("test head object================")
	fmt.Println(b.Head("testkey"))

	fmt.Println("test list object=============")
	fmt.Println(b.List("", "", "", 2))

	// fmt.Println("test create bucket==========")
	// fmt.Println(s.PutBucket("test.create", scs.ACLPrivate))

	// fmt.Println("test delete bucket===========")
	// fmt.Println(s.DeleteBucket("test.create"))

	fmt.Println("test MultipartUpload======")
	mp, err := b.InitiateMultipartUpload("testmukey", map[string]string{
		"x-amx-meta-mu": "foo",
	})
	if err == nil {
		mps := make([]scs.Part, 0)
		for i := 1; i <= 3; i++ {
			p, _ := b.UploadPart("testmukey", mp.UploadID, i, bytes.NewBufferString("mmmmm\n"))
			mps = append(mps, p)
		}
		lps, err := b.ListParts(mp.Key, mp.UploadID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(mps)
		fmt.Println(lps)
		b.CompleteMultipartUpload(mp.Key, mp.UploadID, mps)
		rdata, _ := b.Get("testmukey", 0, 0)
		readPrint(rdata)
	} else {
		fmt.Println(err)
	}
}

func readPrint(data io.ReadCloser) {
	bts, _ := ioutil.ReadAll(data)
	fmt.Println(string(bts))
}
