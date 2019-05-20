package main

import (
  "github.com/ceph/go-ceph/rados"
  "github.com/ceph/go-ceph/rbd"
  "fmt"
  "strings"
  "strconv"
)

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

//Generalized image generate
func genImage(poolName string, imageName string, Size string, Feature string){
  var imageSize    uint64
  var imageFeature uint64

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, poolName)

  switch Feature {
    case "RbdFeatureLayering":
      imageFeature = 1
    case "RbdFeatureStripingV2":
      imageFeature = 2
    case "RbdFeatureExclusiveLock":
      imageFeature = 4
    case "RbdFeatureObjectMap":
      imageFeature = 8
    case "RbdFeatureFastDiff":
      imageFeature = 16
    case "RbdFeatureDeepFlatten":
      imageFeature = 32
    case "RbdFeatureJournaling":
      imageFeature = 64
    case "RbdFeatureDataPool":
      imageFeature = 128
    case "RbdFeaturesDefault":
      imageFeature = 61
    default:
      fmt.Println("Feature Error")
      return
  }

  sizeUnit := Size[len(Size)-1:]
  switch sizeUnit {
    case "M":
      Size  = strings.TrimSuffix(Size, "M")
      sizeVal,err := strconv.ParseUint(Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",imageName," Failed")
      }
      imageSize = sizeVal * 1024 * 1024
      createImage(ioctx, imageName, imageSize, imageFeature)
      fmt.Println("Image ",imageName," Successfully Created")
    case "G":
      Size  = strings.TrimSuffix(Size, "G")
      sizeVal,err := strconv.ParseUint(Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",imageName," Failed")
      }
      imageSize = sizeVal * 1024 * 1024 * 1024
      createImage(ioctx, imageName, imageSize, imageFeature)
      fmt.Println("Image ",imageName," Successfully Created")
    default:
      fmt.Println("Size Error")
      return
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
