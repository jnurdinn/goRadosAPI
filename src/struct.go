package main

import(
  "github.com/mitchellh/goamz/s3"
  "github.com/QuentinPerez/go-radosgw/pkg/api"
  )

// Object Definitions
type Config struct {
    NextRadosHost     string
    AdminRGWHost      string
    AdminRGWUID       string
    AdminRGWAccessKey string
    AdminRGWSecretKey string
}
type MainInfo struct {
    Version           string `json:"version,omitempty"`
    Info              string `json:"info,omitempty"`
}
type Response struct {
    Success           bool   `json:"success,omitempty"`
    Message           string `json:"message,omitempty"`
}
type CephFSInfo struct {
    FSName            string     `json:"fsname,omitempty"`
    MetadataPool      string     `json:"metadatapool,omitempty"`
    DataPools         []string   `json:"datapools,omitempty"`
    MDSNode           string     `json:"mdsnode,omitempty"`
}
type genCephFSRequest struct {
    FSName            string
    MetaPool          string
    DataPool          string
}
type addCephFSPoolRequest struct {
    FSName            string
}
type delCephFSRequest struct {
    FSName            string
}
type RBDInfo struct {
    PoolName          string      `json:"poolname,omitempty"`
    Image             []ImageInfo  `json:"images,omitempty"`
}
type ImageInfo struct {
    Name              string  `json:"name,omitempty"`
    Feature           string  `json:"feature,omitempty"`
    Size              uint64  `json:"size,omitempty"`
    ObjSize           uint64  `json:"objsize,omitempty"`
    NumObjs           uint64  `json:"numobjs,omitempty"`
    Order             int     `json:"order,omitempty"`
    BlockNamePrefix   string  `json:"blocknameprefix,omitempty"`
    ParentPool        int64   `json:"parentpool,omitempty"`
    ParentName        string  `json:"parentname,omitempty"`
}
type genRBDRequest struct {
    Name              string
    Size              string
    Feature           string
}
type delRBDRequest struct {
    Name              string
}
type putRGWAccountRequest struct {
    UID               string  `json:"uid,omitempty"`
    DisplayName       string  `json:"displayname,omitempty"`
    UserCaps          string  `json:"usercaps,omitempty"`
    MaxBuckets        *int    `json:"maxbuckets,omitempty`
    IsSuspended       bool    `json:"issuspended,omitempty`
}
type delRGWAccountRequest struct {
    UID               string             `json:"uid,omitempty"`
}
type userRGWBucket    struct {
    UID               string             `json:"uid,omitempty"`
    Buckets           []radosAPI.Buckets   `json:"buckets,omitempty"`
}
type genRGWBucketRequest struct {
    Name              string
    ACL               s3.ACL
}
type delRGWBucketRequest struct {
    Name              string
}
