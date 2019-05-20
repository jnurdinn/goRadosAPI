package main

import (
  "github.com/ceph/go-ceph/rados"
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
