package scs

import "github.com/Arvintian/scs-go-sdk/pkg/client"

//SCS type
type SCS struct {
	c     *client.Client
	debug bool
}

//BucketMeta bucket meta信息
type BucketMeta struct {
	DeleteQuantity   int64  `json:"DeleteQuantity"`
	Capacity         int64  `json:"Capacity"`
	ACL              ACL    `json:"ACL"`
	ProjectID        int64  `json:"ProjectID"`
	DownloadQuantity int64  `json:"DownloadQuantity"`
	DownloadCapacity int64  `json:"DownloadCapacity"`
	CapacityC        int64  `json:"CapacityC"`
	QuantityC        int64  `json:"QuantityC"`
	Project          string `json:"Project"`
	UploadCapacity   int64  `json:"UploadCapacity"`
	UploadQuantity   int64  `json:"UploadQuantity"`
	LastModified     string `json:"LastModified"`
	SizeC            int64  `json:"SizeC"`
	Owner            string `json:"Owner"`
	DeleteCapacity   int64  `json:"DeleteCapacity"`
	Quantity         int64  `json:"Quantity"`
}

//ACL acl规则
type ACL map[string][]string

//ACLInfo acl信息
type ACLInfo struct {
	Owner string `json:"Owner"`
	ACL   ACL    `json:"ACL"`
}

// BucketList type
type BucketList struct {
	Owner   Owner
	Buckets []Bucket
}

// Bucket type
type Bucket struct {
	Name          string `json:"Name"`
	ConsumedBytes int64  `json:"ConsumedBytes"`
	CreationDate  string `json:"CreationDate"`
	c             *client.Client
}

// Owner type
type Owner struct {
	ID string `json:"ID"`
}

// ObjectMeta type
type ObjectMeta struct {
	ContentType   string
	ContentLength int64
	ETag          string
	LastModified  string
	XAmzMeta      map[string]string
}

// Content-Type: <object-mime-type>
// Content-Length: <object-file-bytes>
// ETag: "<文件的MD5值>"
// Last-Modified: <最后修改时间>
// X-RequestId: 00078d50-1404-0810-5947-782bcb10b128
// X-Requester: Your UserId
// x-amz-meta-foo1: <value1> #自定义meta：foo1
// x-amz-meta-foo2: <value2> #自定义meta：foo2
