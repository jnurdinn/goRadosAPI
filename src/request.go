package main

import (
  "encoding/json"
  "github.com/gorilla/mux"
  "fmt"
  "net/http"
  "github.com/QuentinPerez/go-radosgw/pkg/api"
)

const (
    DefaultPoolRBD         = "rbd"
)

var (
  config      Config
)

/* CEPHFS AVAILABLE REQUESTS LIST  */

// Display Main Page Info for GET Method
func getMainReq(w http.ResponseWriter, r *http.Request){
    maininfo := MainInfo{Version: "1.0", Info: "NextRados Seems to Work Properly"}
    json.NewEncoder(w).Encode(maininfo)
}

//Display CephFS Info for GET Method
func getCephFSReq(w http.ResponseWriter, r *http.Request){
    var listfs CephFSInfo
    listfs = listFS()
    cephfsinfo := CephFSInfo{FSName: listfs.FSName, MetadataPool: listfs.MetadataPool, DataPools: listfs.DataPools, MDSNode: listfs.MDSNode}
    json.NewEncoder(w).Encode(cephfsinfo)
}

// Generate new CephFS for POST Method
func genCephFSReq(w http.ResponseWriter, r *http.Request){
  var cephfs  genCephFSRequest
  _ = json.NewDecoder(r.Body).Decode(&cephfs)
  genFS(cephfs.FSName,cephfs.MetaPool,cephfs.DataPool)
}

// Add data pool to CephFS for POST Method
func addCephFSPoolReq(w http.ResponseWriter, r *http.Request){
  var cephfs  addCephFSPoolRequest
  params := mux.Vars(r)
  _ = json.NewDecoder(r.Body).Decode(&cephfs)
  addFSPool(cephfs.FSName, params["pool"])
}

// Delete CephFS for DELETE Method
func delCephFSReq(w http.ResponseWriter, r *http.Request) {
  var cephfs delCephFSRequest
	_ = json.NewDecoder(r.Body).Decode(&cephfs)
  delFS(cephfs.FSName)
}

// Delete Datapool from CephFS for POST METHOD
func delCephFSPoolReq(w http.ResponseWriter, r *http.Request) {
  var cephfs delCephFSRequest
  params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&cephfs)
  delFSPool(cephfs.FSName, params["pool"])
}

/* RBD AVAILABLE REQUESTS LIST  */

// Display RBD Pool & Image Info for GET Method
func getRBDReq(w http.ResponseWriter, r *http.Request){
    var imageInfos []ImageInfo

    conn := newConn()
    defer conn.Shutdown()
    ioctx  := openIOCtx(conn, DefaultPoolRBD)
    imageList := listImages(ioctx)

    for _,image := range imageList {
      imageInfos = append(imageInfos, getImageInfo(ioctx,image))
    }

    rbdinfo := RBDInfo{PoolName: DefaultPoolRBD, Image: imageInfos}
    json.NewEncoder(w).Encode(rbdinfo)
}

// Display RBD Pool & Image Info for GET Method
func getRBDPoolReq(w http.ResponseWriter, r *http.Request){
    var imageInfos []ImageInfo
    params := mux.Vars(r)

    conn := newConn()
    defer conn.Shutdown()
    ioctx  := openIOCtx(conn, params["pool"])
    imageList := listImages(ioctx)

    for _,image := range imageList {
      imageInfos = append(imageInfos, getImageInfo(ioctx,image))
    }

    rbdinfo := RBDInfo{PoolName: params["pool"], Image: imageInfos}
    json.NewEncoder(w).Encode(rbdinfo)
}

// Generate New Image for POST Method
func genRBDImageReq(w http.ResponseWriter, r *http.Request) {
  var rbd          genRBDRequest
	_ = json.NewDecoder(r.Body).Decode(&rbd)
  genImage(DefaultPoolRBD,rbd.Name,rbd.Size,rbd.Feature)
}

// Generate New Image for POST Method in Particular Pool
func genRBDImagePoolReq(w http.ResponseWriter, r *http.Request) {
  var rbd          genRBDRequest
  params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&rbd)
  fmt.Println(params["pool"])
  genImage(params["pool"],rbd.Name,rbd.Size,rbd.Feature)
}

// Delete Image for DELETE Method
func delRBDImageReq(w http.ResponseWriter, r *http.Request) {
  var rbd delRBDRequest

	_ = json.NewDecoder(r.Body).Decode(&rbd)
  conn := newConn()
  defer conn.Shutdown()
  ioctx := openIOCtx(conn, DefaultPoolRBD)
  deleteImage(ioctx, rbd.Name)
  fmt.Println("Image ",rbd.Name," Successfully Deleted")
}

// Delete Image for DELETE Method
func delRBDImagePoolReq(w http.ResponseWriter, r *http.Request) {
  var rbd delRBDRequest
  params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, params["pool"])
  deleteImage(ioctx, rbd.Name)
  fmt.Println("Image ",rbd.Name," Successfully Deleted")
}

/* RGW AVAILABLE REQUESTS LIST  */
// List All RGW Users for GET Method
func getRGWUsersReq(w http.ResponseWriter, r *http.Request){
  var rgwUsers []*radosAPI.User
  adminAPI := initAdminAPI(conf.AdminRGWHost, conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  listUsers := listRGWUsers(adminAPI)
  for _, uids := range listUsers {
    rgwUsers = append(rgwUsers, getRGWUser(adminAPI,uids))
  }
  json.NewEncoder(w).Encode(rgwUsers)
}

// Create New RGW User for POST Method
func genRGWUserReq(w http.ResponseWriter, r *http.Request){
  var rgw      putRGWAccountRequest
	_ = json.NewDecoder(r.Body).Decode(&rgw)
  adminAPI := initAdminAPI(conf.AdminRGWHost, conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  fmt.Println(adminAPI, rgw.UID , rgw.DisplayName, rgw.UserCaps, rgw.MaxBuckets , rgw.IsSuspended)
  genRGWUser(adminAPI, rgw.UID , rgw.DisplayName, rgw.UserCaps, rgw.MaxBuckets , rgw.IsSuspended)
  fmt.Println("New Account Successfully Created")
}

// Edit RGW User for PATCH Method
func editRGWUserReq(w http.ResponseWriter, r *http.Request){
  var rgw      putRGWAccountRequest
	_ = json.NewDecoder(r.Body).Decode(&rgw)
  adminAPI := initAdminAPI(conf.AdminRGWHost, conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  editRGWUser(adminAPI, rgw.UID , rgw.DisplayName, rgw.UserCaps, rgw.MaxBuckets , rgw.IsSuspended)
  fmt.Println("Account Successfully Updated")
}

// Delete RGW User for DELETE Method
func delRGWUserReq(w http.ResponseWriter, r *http.Request){
  var rgw      delRGWAccountRequest
	_ = json.NewDecoder(r.Body).Decode(&rgw)
  adminAPI := initAdminAPI(conf.AdminRGWHost, conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  delRGWUser(adminAPI, rgw.UID)
  fmt.Println("Account Successfully Deleted")
}

// List All Available Buckets Corresponding to User for GET Method
func getRGWBucketReq(w http.ResponseWriter, r *http.Request){
  var (
    rgwUserBuckets      []userRGWBucket
    bucketInfos         []radosAPI.Buckets
  )
  adminAPI := initAdminAPI(conf.AdminRGWHost, conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  userBuckets := getRGWBuckets(adminAPI, conf.AdminRGWUID)
  for _, bucket := range userBuckets {
    bucketInfos = append(bucketInfos, getRGWBucket(adminAPI,bucket.Name))
  }
  rgwUserBuckets = append(rgwUserBuckets, userRGWBucket{UID: conf.AdminRGWUID, Buckets : bucketInfos})
  json.NewEncoder(w).Encode(rgwUserBuckets)
}

// Create New RGW S3 Bucket for POST Method
func genRGWBucketReq(w http.ResponseWriter, r *http.Request){
  var rgw      genRGWBucketRequest
	_ = json.NewDecoder(r.Body).Decode(&rgw)
  s3API := initS3API(conf.AdminRGWHost,conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  genS3Bucket(s3API,rgw.Name, rgw.ACL)
  fmt.Println("New Bucket Successfully Created")
}

// Delete RGW S3 Bucket for DELETE Method
func delRGWBucketReq(w http.ResponseWriter, r *http.Request){
  var rgw      genRGWBucketRequest
	_ = json.NewDecoder(r.Body).Decode(&rgw)
  s3API := initS3API(conf.AdminRGWHost,conf.AdminRGWAccessKey, conf.AdminRGWSecretKey)
  delS3Bucket(s3API,rgw.Name)
  fmt.Println("Bucket Successfully Deleted")
}
