package scs

//ACLPrivate 私有acl
const ACLPrivate = "private"

//ACLPublicRead 公开读
const ACLPublicRead = "public-read"

//ACLPublicReadWrite 公开读写
const ACLPublicReadWrite = "public-read-write"

//ACLAuthenticatedRead 授权读
const ACLAuthenticatedRead = "authenticated-read"

//TempFilePrefix 上传临时文件前缀
const TempFilePrefix = "scs-temp-"

//MD5Threshold Memory footprint threshold for each MD5 computation (16MB is the default), in byte. When the data is more than that, temp file is used.
const MD5Threshold = 16 * 1024 * 1024
