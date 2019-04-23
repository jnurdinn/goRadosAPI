package main

import (
  "github.com/ceph/go-ceph/rados"
  "github.com/ceph/go-ceph/rbd"
  "fmt"
)

// Make Connection to Cluster
func newConn() *rados.Conn {
    conn, err := rados.NewConn()
    if err != nil {
      fmt.Println("error when init", err)
    }
    err = conn.ReadDefaultConfigFile()
    if err != nil {
      fmt.Println("error when Read Config", err)
    }
    err = conn.Connect()
    if err != nil {
      fmt.Println("error when connect", err)
    }
    return conn
}

//OpenIOContext
func openIOCtx(conn *rados.Conn, poolName string) *rados.IOContext {
  ioctx, err := conn.OpenIOContext(poolName)
  if err != nil {
      fmt.Println("error when openIOContext", err)
  }
  return ioctx
}

//List RBD Image
func listImages(ioctx *rados.IOContext) []string {
    imageNames, err := rbd.GetImageNames(ioctx)
    if err != nil {
        fmt.Println("error when getImagesNames", err)
    }
    return imageNames
}

//OpenIOContext
func createImage(ioctx *rados.IOContext, imageName string, imageSize uint64, feature uint64) {
  _,err := rbd.Create(ioctx,imageName,imageSize, 22, feature)
  if err != nil {
      fmt.Println("RBD create image failed: ",err)
  }
}

//Select Image to Modify
func getImage(ioctx *rados.IOContext, imageName string) *rbd.Image {
  img := rbd.GetImage(ioctx,imageName)
  if err := img.Open(); err != nil {
      fmt.Println("RBD open image  failed: ",err)
  }
  defer img.Close()
  return img
}

//Get All Image Info
func getImageInfo(ioctx *rados.IOContext, imageName string) ImageInfo {
  var imageFeature string
  var imageInfo    ImageInfo

  img := getImage(ioctx,imageName)
  _ = img.Open()
  feature, err := img.GetFeatures()
  if err != nil {
      fmt.Println("Get Image Feature failed: ",err)
  }
  switch feature {
    case 1 :
      imageFeature = "RbdFeatureLayering"
    case 2 :
      imageFeature = "RbdFeatureStripingV2"
    case 4 :
      imageFeature = "RbdFeatureExclusiveLock"
    case 8 :
      imageFeature = "RbdFeatureObjectMap"
    case 16 :
      imageFeature = "RbdFeatureFastDiff"
    case 32 :
      imageFeature = "RbdFeatureDeepFlatten"
    case 64 :
      imageFeature = "RbdFeatureJournaling"
    case 128 :
      imageFeature = "RbdFeatureDataPool"
    case 61 :
      imageFeature = "RbdFeaturesDefault"
    default:
      imageFeature = "Error"
      fmt.Println("Get Feature Error")
  }
  stat, err := img.Stat()
  if err != nil {
    fmt.Println("Get Image Info: ",err)
  }
  _ = img.Close()
  imageInfo = ImageInfo{Name : imageName, Feature : imageFeature, Size : stat.Size, ObjSize : stat.Obj_size, NumObjs : stat.Num_objs, Order : stat.Order, BlockNamePrefix : stat.Block_name_prefix, ParentPool : stat.Parent_pool, ParentName : stat.Parent_name}
  return imageInfo
}

//Delete Selected Image
func deleteImage(ioctx *rados.IOContext, imageName string) {
  img := getImage(ioctx,imageName)
  _ = img.Remove()
}
