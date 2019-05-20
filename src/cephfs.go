package main

import (
  "os/exec"
  "fmt"
  "strings"
)

//List all CephFS
func listFS() CephFSInfo { //CephFSInfo
  var fsInfo CephFSInfo
  cmd := exec.Command("ceph", "fs", "ls")
  stdout, err := cmd.Output()
  if err != nil {
      fmt.Println(err.Error())
  }
  out := strings.Split(string(stdout),",")
  if (out[0][:7] == "No file") {
    fsInfo.FSName = out[0]
    fsInfo.MetadataPool = ""
    fsInfo.DataPools = append(fsInfo.DataPools, "")
    fsInfo.MDSNode = ""
  } else {
    fsInfo.FSName = out[0][6:]
    fsInfo.MetadataPool = out[1][16:]
    fsInfo.DataPools = strings.Split(out[2][14:len(out[2])-3]," ")
    /*cmd = exec.Command("ceph", "fs", "get", fsInfo.FSName)
    stdout, err = cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
    }
    out = strings.Split(string(stdout),"\n")
    fsInfo.MDSNode = strings.Split(out[26],"'")[1]*/
  }
  return fsInfo
}

//Generate new CephFS
func genFS(fsName string, metaPool string, dataPool string) string {
  cmd := exec.Command("ceph", "fs", "new", fsName, metaPool, dataPool)
  stdout, err := cmd.Output()
  if err != nil {
      fmt.Println(err.Error())
  }
  return string(stdout)
}

//Delete CephFS
func delFS(fsName string) string {
  cmd := exec.Command("ceph", "fs", "rm", fsName, "--yes-i-really-mean-it")
  stdout, err := cmd.Output()
  if err != nil {
      fmt.Println(err.Error())
  }
  fmt.Println(string(stdout))
  return string(stdout)
}


//Add data pool to CephFS
func addFSPool(fsName string, dataPool string) string {
  cmd := exec.Command("ceph", "fs", "add_data_pool", fsName, dataPool)
  stdout, err := cmd.Output()
  if err != nil {
      fmt.Println(err.Error())
  }
  return string(stdout)
}

//Remove data pool from CephFS
func delFSPool(fsName string, dataPool string) string {
  cmd := exec.Command("ceph", "fs", "rm_data_pool", fsName, dataPool)
  stdout, err := cmd.Output()
  if err != nil {
      fmt.Println(err.Error())
  }
  return string(stdout)
}
