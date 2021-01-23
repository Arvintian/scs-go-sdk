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

	// fmt.Println("test list get bucket ============")
	// fmt.Println(s.ListBuckets())
	b, _ := s.GetBucket(bucket)
	// // fmt.Println(s.PutBucket("test"))
	// // fmt.Println(s.GetBucketACL("test"))
	// // fmt.Println(s.PutBucketACL("test"))

	fmt.Println("test put==========")
	data := bytes.NewBufferString("asssssbb")
	err = b.Put("testkey", map[string]string{
		"x-amz-meta-test": "foo",
	}, data)
	if err != nil {
		fmt.Println("put error", err)
	}

	// fmt.Println("test head object================")
	// fmt.Println(b.Head("testkey"))

	fmt.Println("test get object==========")
	rdata, err := b.Get("testkey", "0-0")
	if err != nil {
		fmt.Println(err)
	}
	readPrint(rdata)

	// fmt.Println("test delete object==============")
	// fmt.Println(b.Delete("testkey"))
	// fmt.Println("test head object================")
	// fmt.Println(b.Head("testkey"))

	// fmt.Println("test list object=============")
	// fmt.Println(b.List("", "", "", 2))

	// // fmt.Println("test create bucket==========")
	// // fmt.Println(s.PutBucket("test.create", scs.ACLPrivate))

	// // fmt.Println("test delete bucket===========")
	// // fmt.Println(s.DeleteBucket("test.create"))

	// fmt.Println("test MultipartUpload======")
	// mp, err := b.InitiateMultipartUpload("testmukey", map[string]string{
	// 	"x-amx-meta-mu": "foo",
	// })
	// if err == nil {
	// 	mps := make([]scs.Part, 0)
	// 	for i := 1; i <= 3; i++ {
	// 		p, err := b.UploadPart("testmukey", mp.UploadID, i, bytes.NewBufferString("ammmmm"))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		mps = append(mps, p)
	// 	}
	// 	lps, err := b.ListParts(mp.Key, mp.UploadID)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(mps)
	// 	fmt.Println(lps)
	// 	b.CompleteMultipartUpload(mp.Key, mp.UploadID, mps)
	// 	rdata, _ := b.Get("testmukey", "0-")
	// 	readPrint(rdata)
	// } else {
	// 	fmt.Println(err)
	// }

	get := func(key string, off, limit int64) (io.ReadCloser, error) {
		if off > 0 || limit > 0 {
			var r string
			if limit > 0 {
				r = fmt.Sprintf("%d-%d", off, off+limit-1)
			} else {
				r = fmt.Sprintf("%d-", off)
			}
			fmt.Println(key, r)
			return b.Get(key, r)
		}
		fmt.Println(key, "")
		return b.Get(key, "")
	}
	//fmt.Println(get)
	for {
		size := 5 << 20
		for i := 0; i < 10; i++ {
			r, err := get("dir1/CentOS-8.3.2011-aarch64-dvd1.iso", int64(i*size), int64(size))
			if err != nil {
				fmt.Println("eee", err)
			}
			data := make([]byte, size)
			_, err = io.ReadFull(r, data)
			if err != nil {
				fmt.Println("sss", err)
			}
			fmt.Println("bbb", i, len(data))
			// bts, err = ioutil.ReadAll(r)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// fmt.Println(len(bts))
		}
	}
}

func readPrint(data io.ReadCloser) {
	bts, _ := ioutil.ReadAll(data)
	fmt.Println(string(bts))
}
