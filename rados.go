package main

import (
  "github.com/ceph/go-ceph/rados"
  "github.com/ceph/go-ceph/rbd"
  "os"
  "fmt"
)

// Make Connection to Cluster
func newConn() *rados.Conn {
    conn, err := rados.NewConn()
    if err != nil {
      fmt.Println("error when init", err)
      os.Exit(1)
    }
    err = conn.ReadDefaultConfigFile()
    if err != nil {
      fmt.Println("error when Read Config", err)
      os.Exit(1)
    }
    err = conn.Connect()
    if err != nil {
      fmt.Println("error when connect", err)
      os.Exit(1)
    }
    return conn
}

//List All Pools
func listPools(conn *rados.Conn) []string {
    pools, err := conn.ListPools()
    if err != nil {
        fmt.Println("error when list pool", err)
        os.Exit(1)
    }
    return pools
}

//OpenIOContext
func openIOCtx(conn *rados.Conn, poolName string) *rados.IOContext {
  ioctx, err := conn.OpenIOContext(poolName)
  if err != nil {
      fmt.Println("error when openIOContext", err)
      os.Exit(1)
  }
  return ioctx
}

//List RBD Image
func listImages(ioctx *rados.IOContext) []string {
    imageNames, err := rbd.GetImageNames(ioctx)
    if err != nil {
        fmt.Println("error when getImagesNames", err)
        os.Exit(1)
    }
    return imageNames
}

//OpenIOContext
func createImage(ioctx *rados.IOContext, imageName string, imageSize uint64, feature uint64) {
  _,err := rbd.Create(ioctx,imageName,imageSize, 22, feature)
  if err != nil {
      fmt.Println("RBD create image failed: ",err)
      os.Exit(1)
  }
}

//Select Image to Modify
func getImage(ioctx *rados.IOContext, imageName string) *rbd.Image {
  img := rbd.GetImage(ioctx,imageName)
  if err := img.Open(); err != nil {
      fmt.Println("RBD open image  failed: ",err)
      os.Exit(1)
  }
  defer img.Close()
  return img
}

//Delete Selected Image
func deleteImage(ioctx *rados.IOContext, imageName string) {
  img := getImage(ioctx,imageName)
  _ = img.Remove()
}
