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

// ListObject type
type ListObject struct {
	Delimiter              string         `json:"Delimiter"`
	Prefix                 string         `json:"Prefix"`
	CommonPrefixes         []CommonPrefix `json:"CommonPrefixes"`
	Marker                 string         `json:"Marker"`
	ContentsQuantity       int64          `json:"ContentsQuantity"`
	CommonPrefixesQuantity int64          `json:"CommonPrefixesQuantity"`
	NextMarker             string         `json:"NextMarker"`
	IsTruncated            bool           `json:"IsTruncated"`
	Contents               []Object       `json:"Contents"`
}

// CommonPrefix type
type CommonPrefix struct {
	Prefix string `json:"Prefix"`
}

// Object type
type Object struct {
	SHA1         string `json:"SHA1"`
	Name         string `json:"Name"`
	LastModified string `json:"Last-Modified"`
	Owner        string `json:"Owner"`
	MD5          string `json:"MD5"`
	ContentType  string `json:"ContentType"`
	Size         int64  `json:"Size"`
}
