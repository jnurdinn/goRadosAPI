package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

// main function to boot up everything
func main() {
    //various routes
    router := mux.NewRouter()
    router.HandleFunc("/", getMainReq).Methods("GET")

    //CephFS Requests
    router.HandleFunc("/cephfs", getCephFSReq).Methods("GET")
    router.HandleFunc("/cephfs", genCephFSReq).Methods("POST")
    router.HandleFunc("/cephfs/{pool}", addCephFSPoolReq).Methods("POST")
    router.HandleFunc("/cephfs", delCephFSReq).Methods("DELETE")
    router.HandleFunc("/cephfs/{pool}", delCephFSPoolReq).Methods("DELETE")

    //RBD Requests
    router.HandleFunc("/rbd", getRBDReq).Methods("GET")
    router.HandleFunc("/rbd/{pool}", getRBDPoolReq).Methods("GET")
    router.HandleFunc("/rbd", genRBDImageReq).Methods("POST")
    router.HandleFunc("/rbd/{pool}", genRBDImagePoolReq).Methods("POST")
    router.HandleFunc("/rbd", delRBDImageReq).Methods("DELETE")
    router.HandleFunc("/rbd/{pool}", delRBDImagePoolReq).Methods("DELETE")

    //RGW Requests
    router.HandleFunc("/rgw", getRGWUsersReq).Methods("GET")
    router.HandleFunc("/rgw", getRGWBucketsReq).Methods("GET")
    router.HandleFunc("/rgw", genRGWUserReq).Methods("POST")
    router.HandleFunc("/rgw/s3", genRGWBucketReq).Methods("POST")
    router.HandleFunc("/rgw", editRGWUserReq).Methods("PATCH")
    router.HandleFunc("/rgw", delRGWUserReq).Methods("DELETE")
    router.HandleFunc("/rgw/s3", delRGWBucketReq).Methods("DELETE")

    log.Fatal(http.ListenAndServe(NextRadosHost, router))
}
