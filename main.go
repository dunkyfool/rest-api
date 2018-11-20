package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
)

type Status struct {
    ScanID      string  `json:"scanid,omitempty"`
    QueueStatus string  `json:"queueStatus,omitempty"`
    ConcisedLog string  `json:"concisedLog,omitempty"`
}

// main
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/search/{jid}", GetScanID).Methods("GET")
    router.HandleFunc("/status/{scanid}", GetStatus).Methods("GET")
    router.HandleFunc("/scan/{scanid}", TriggerScan).Methods("GET")
    log.Fatal(http.ListenAndServe(":5000", router))
}

// use jid to get all scan id
func GetScanID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    fmt.Println("[INFO] Project ID: " + params["jid"])
}

// Get Queue & last build concised log
func GetStatus(w http.ResponseWriter, r *http.Request) {
    var currentStatus Status
    currentStatus = Status{ScanID: "1001", QueueStatus: "1 job(s)",  ConcisedLog: "WTF"}

    params := mux.Vars(r)
    fmt.Println("[INFO] Scan ID: " + params["scanid"])
    json.NewEncoder(w).Encode(currentStatus)
}

// Trigger scan
func TriggerScan(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    fmt.Println("[INFO] Scan ID: " + params["scanid"])
}
