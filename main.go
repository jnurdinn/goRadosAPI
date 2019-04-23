package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "strings"
    "strconv"
)

// Object Definitions
type MainInfo struct {
    Version   string `json:"version,omitempty"`
    Info      string `json:"info,omitempty"`
}
type Response struct {
    Success   bool   `json:"success,omitempty"`
    Message   string `json:"message,omitempty"`
}
type RBDInfo struct {
    PoolName  string      `json:"poolname,omitempty"`
    Image    []ImageInfo  `json:"images,omitempty"`
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
    Name      string
    Size      string
    Feature   string
}
type delRBDRequest struct {
    Name      string
}

var maininfo []MainInfo
var rbdinfo []RBDInfo
var response []Response

const (
    DefaultPoolName         = "rbd"
)

//Display Main Page Info
func getMain(w http.ResponseWriter, r *http.Request){
    maininfo := MainInfo{Version: "1.0", Info: "NextRados Seems to Work Properly"}
    json.NewEncoder(w).Encode(maininfo)
}

//Display RBD Pool & Image Info
func getRBD(w http.ResponseWriter, r *http.Request){
    var imageInfos []ImageInfo

    conn := newConn()
    defer conn.Shutdown()
    ioctx  := openIOCtx(conn, DefaultPoolName)
    imageList := listImages(ioctx)

    for _,image := range imageList {
      imageInfos = append(imageInfos, getImageInfo(ioctx,image))
    }

    rbdinfo := RBDInfo{PoolName: DefaultPoolName, Image: imageInfos}
    json.NewEncoder(w).Encode(rbdinfo)
}

//Display RBD Pool & Image Info
func getRBDPool(w http.ResponseWriter, r *http.Request){
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
func genRBDImage(w http.ResponseWriter, r *http.Request) {
  var rbd          genRBDRequest
  var imageSize    uint64
  var imageFeature uint64
  //var response Response

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, DefaultPoolName)

  switch rbd.Feature {
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

  sizeUnit := rbd.Size[len(rbd.Size)-1:]
  switch sizeUnit {
    case "M":
      rbd.Size  = strings.TrimSuffix(rbd.Size, "M")
      sizeVal,err := strconv.ParseUint(rbd.Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",rbd.Name," Failed",rbd.Size)
      }
      imageSize = sizeVal * 1024 * 1024
      createImage(ioctx, rbd.Name, imageSize, imageFeature)
      fmt.Println("Image ",rbd.Name," Successfully Created")
    case "G":
      rbd.Size  = strings.TrimSuffix(rbd.Size, "G")
      sizeVal,err := strconv.ParseUint(rbd.Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",rbd.Name," Failed",rbd.Size)
      }
      imageSize = sizeVal * 1024 * 1024
      createImage(ioctx, rbd.Name, imageSize, imageFeature)
      fmt.Println("Image ",rbd.Name," Successfully Created")
    default:
      fmt.Println("Size Error")
      return
      //response = Response{Success:false, Message: "Image Creation Fail"}
  }
  //response = Response{Success:true, Message: "Image Successfully Created"}
  //json.NewEncoder(w).Encode(response)
}

// Generate New Image for POST Method in Particular Pool
func genRBDImagePool(w http.ResponseWriter, r *http.Request) {
  var rbd          genRBDRequest
  var imageSize    uint64
  var imageFeature uint64
  params := mux.Vars(r)
  //var response Response

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, params["pool"])

  switch rbd.Feature {
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

  sizeUnit := rbd.Size[len(rbd.Size)-1:]
  switch sizeUnit {
    case "M":
      rbd.Size  = strings.TrimSuffix(rbd.Size, "M")
      sizeVal,err := strconv.ParseUint(rbd.Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",rbd.Name," Failed",rbd.Size)
      }
      imageSize = sizeVal * 1024 * 1024
      createImage(ioctx, rbd.Name, imageSize, imageFeature)
      fmt.Println("Image ",rbd.Name," Successfully Created")
    case "G":
      rbd.Size  = strings.TrimSuffix(rbd.Size, "G")
      sizeVal,err := strconv.ParseUint(rbd.Size, 10, 64)
      if err != nil {
        fmt.Println("Image ",rbd.Name," Failed",rbd.Size)
      }
      imageSize = sizeVal * 1024 * 1024
      createImage(ioctx, rbd.Name, imageSize, imageFeature)
      fmt.Println("Image ",rbd.Name," Successfully Created")
    default:
      fmt.Println("Size Error")
      return
      //response = Response{Success:false, Message: "Image Creation Fail"}
  }
  //response = Response{Success:true, Message: "Image Successfully Created"}
  //json.NewEncoder(w).Encode(response)
}

// Delete Image for DELETE Method
func delRBDImage(w http.ResponseWriter, r *http.Request) {
  var rbd delRBDRequest
  //var response Response

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, DefaultPoolName)
  deleteImage(ioctx, rbd.Name)
  fmt.Println("Image ",rbd.Name," Successfully Deleted")
  //response = Response{Success:true, Message: "Image Successfully Deleted"}
  //json.NewEncoder(w).Encode(response)
}

// Delete Image for DELETE Method
func delRBDImagePool(w http.ResponseWriter, r *http.Request) {
  var rbd delRBDRequest
  params := mux.Vars(r)
  //var response Response

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, params["pool"])
  deleteImage(ioctx, rbd.Name)
  fmt.Println("Image ",rbd.Name," Successfully Deleted")
  //response = Response{Success:true, Message: "Image Successfully Deleted"}
  //json.NewEncoder(w).Encode(response)
}

// main function to boot up everything
func main() {

    //Restfull Address
    var addr string
    addr = "10.10.2.103:8080"

    //various routes
    router := mux.NewRouter()
    router.HandleFunc("/", getMain).Methods("GET")
    router.HandleFunc("/rbd", getRBD).Methods("GET")
    router.HandleFunc("/rbd/{pool}", getRBDPool).Methods("GET")
    router.HandleFunc("/rbd", genRBDImage).Methods("POST")
    router.HandleFunc("/rbd/{pool}", genRBDImagePool).Methods("POST")
    router.HandleFunc("/rbd", delRBDImage).Methods("DELETE")
    router.HandleFunc("/rbd/{pool}", delRBDImagePool).Methods("DELETE")
    log.Fatal(http.ListenAndServe(addr, router))
}
