package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
)

// The person Type (more like an object)
type MainInfo struct {
    Version   string `json:"version,omitempty"`
    Info      string `json:"info,omitempty"`
}
type RBDInfo struct {
    PoolName  string `json:"poolname,omitempty"`
    ImageName []string `json:"imagename,omitempty"`
}
type genRBDRequest struct {
    Name      string
}
type delRBDRequest struct {
    Name      string
}

var maininfo []MainInfo
var rbdinfo []RBDInfo

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
    conn := newConn()
    defer conn.Shutdown()
    ioctx := openIOCtx(conn, DefaultPoolName)
    rbdinfo := RBDInfo{PoolName: DefaultPoolName, ImageName: listImages(ioctx)}
    json.NewEncoder(w).Encode(rbdinfo)
}

// Display a single data
/*
func GetPerson(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, item := range people {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    json.NewEncoder(w).Encode(&Person{})
}
*/

// Generate New Image for POST Method
func genRBDImage(w http.ResponseWriter, r *http.Request) {
  var imageSize uint64
  //var imageFeature uint64
  var rbd genRBDRequest

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, DefaultPoolName)
  imageSize = 4 * 1024 * 1024 * 1024
  createImage(ioctx, rbd.Name, imageSize, 1)
  fmt.Println("Image ",rbd.Name," Successfully Created")
}

// Delete Image for DELETE Method
func delRBDImage(w http.ResponseWriter, r *http.Request) {
  var rbd delRBDRequest

	_ = json.NewDecoder(r.Body).Decode(&rbd)

  conn := newConn()
  defer conn.Shutdown()

  ioctx := openIOCtx(conn, DefaultPoolName)
  deleteImage(ioctx, rbd.Name)
  fmt.Println("Image ",rbd.Name," Successfully Deleted")
}

// main function to boot up everything
func main() {

    //Restfull Address
    var addr string
    addr = "10.10.2.104:8080"

    //various routes
    router := mux.NewRouter()
    router.HandleFunc("/", getMain).Methods("GET")
    router.HandleFunc("/rbd", getRBD).Methods("GET")
    router.HandleFunc("/rbd", genRBDImage).Methods("POST")
    router.HandleFunc("/rbd", delRBDImage).Methods("DELETE")
    log.Fatal(http.ListenAndServe(addr, router))
}
