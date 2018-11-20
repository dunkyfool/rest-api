package main

import (
    //"encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
)

// our main function
func main() {
    router := mux.NewRouter()
    //router.HandleFunc("/search/{jid}", GetScanID).Methods("GET")
    router.HandleFunc("/status/{scanid}", GetStatus).Methods("GET")
    router.HandleFunc("/scan/{scanid}", TriggerScan).Methods("GET")
    log.Fatal(http.ListenAndServe(":5000", router))
}

//func GetScanID(w http.ResponseWriter, r *http.Request) {}
func GetStatus(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    fmt.Println(params["scanid"])
}

func TriggerScan(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r)
}
