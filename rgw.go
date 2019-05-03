package main

import (
  "github.com/QuentinPerez/go-radosgw/pkg/api"
  "github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/s3"
  "fmt"
)

//QUENTINPEREZ API, CONSISTS OF ADMIN OPS

//Init RGW Admin API
func initAdminAPI(addr string, accessKey string, secretKey string) *radosAPI.API {
  api, err := radosAPI.New(addr, accessKey, secretKey)
  if err != nil {
    fmt.Println("Error", err)
  }
  fmt.Println("Admin API Initiated", api)
  return api
}

//List RGW Users List
func listRGWUsers(api *radosAPI.API) []string {
  listUID, err := api.GetUIDs()
  if err != nil {
    fmt.Println("UID List Error", err)
  }
  return listUID
}

//Get Specific RGW User Info
func getRGWUser(api *radosAPI.API, uid string) *radosAPI.User  {
  getUser, err := api.GetUser(uid)
  if err != nil {
    fmt.Println("Get User Error", err)
  }
  return getUser
}

//Get Bucket Info from Specific UID
func getRGWBucket(api *radosAPI.API, uid string) radosAPI.Buckets  {
  var conf radosAPI.BucketConfig
  conf.UID = uid
  getBuckets, err := api.GetBucket(conf)
  if err != nil {
    fmt.Println("Get User Error", err)
  }
  return getBuckets
}


//Create New RGW User
func genRGWUser(api *radosAPI.API, uid string, displayName string, userCaps string, maxBuckets *int, isSuspended bool){
  user, err := api.CreateUser(radosAPI.UserConfig{
    UID:         uid,
    DisplayName: displayName,
    UserCaps:    userCaps,
    MaxBuckets:  maxBuckets,
    Suspended:   isSuspended,
  })
  if err != nil {
    fmt.Println("User Creation Error", err)
  }
  fmt.Println("User Creation Success : ", user)
}

//Edit RGW User
func editRGWUser(api *radosAPI.API, uid string, displayName string, userCaps string, maxBuckets *int, isSuspended bool){
  user, err := api.UpdateUser(radosAPI.UserConfig{
    UID:         uid,
    DisplayName: displayName,
    UserCaps:    userCaps,
    MaxBuckets:  maxBuckets,
    Suspended:   isSuspended,
  })
  if err != nil {
    fmt.Println("User Update Error", err)
  }
  fmt.Println("User Update Success : ", user)
}

//Delete RGW User
func delRGWUser(api *radosAPI.API, uid string){
  err := api.RemoveUser(radosAPI.UserConfig{
    UID: uid,
  })
  if err != nil {
    fmt.Println("User Deletion Error", err)
  }
  fmt.Println("User Deletion Success")
}

//GOAMZ AWS API FOR S3 STORAGE, CONSISTS OF S3 API OPS

//Init RGW S3 API
func initS3API(addr string, accessKey string, secretKey string) *s3.S3 {
  auth := aws.Auth{
      AccessKey: accessKey,
      SecretKey: secretKey,
  }
  conn := s3.New(auth, aws.Region{"","",addr,"",true,true,"","","","","","","","",})
  fmt.Println("S3 API Initiated")
  return conn
}

//Generate New S3 Bucket
func genS3Bucket(conn *s3.S3,bucketName string, acl s3.ACL){
  bucket := conn.Bucket(bucketName)
  res := bucket.PutBucket(acl)
  fmt.Println(res)
}

//Delete S3 Bucket
func delS3Bucket(conn *s3.S3,bucketName string){
  bucket := conn.Bucket(bucketName)
  res := bucket.DelBucket()
  fmt.Println(res)
}
