package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "crypto/tls"
    "strings"
    "time"
)

type Links struct {
        Self struct {
                Href string `json:"href"`
        } `json:"self"`
}

type Step struct {
        Name                string `json:"name"`
        Status              string `json:"status"`
}

type ScanStatus struct {
        ID                  string `json:"id"`
        Name                string `json:"name"`
        Status              string `json:"status"`
        Stages              []Step `json:"stages"`
}

type Response struct {
    ID                      int    `json:"id"`
    Executable              struct {
        Number              int    `json:"number"`
        Url                 string `json:"url"`
    } `json:"executable"`
}


// main
func main() {
    router := mux.NewRouter()
    //router.HandleFunc("/search/{jid}", GetScanID).Methods("GET")
    router.HandleFunc("/status/{scanid}", GetStatus).Methods("GET")
    router.HandleFunc("/scan/{scanid}", TriggerScan).Methods("GET")
    log.Fatal(http.ListenAndServe(":5000", router))
}

// use jid to get all scan id
/*
func GetScanID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    fmt.Println("[INFO] Project ID: " + params["jid"])
}
*/

// Get Queue & last build concised log
func GetStatus(w http.ResponseWriter, r *http.Request) {
    var scanStatus []ScanStatus

    params := mux.Vars(r)
    fmt.Println("[INFO] Scan ID: " + params["scanid"])

    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    resp, err := http.Post("https://jenkins/job/"+params["scanid"]+"/wfapi/runs",
                           "application/x-www-form-urlencoded",
                            strings.NewReader(""))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println(resp)
    err = json.NewDecoder(resp.Body).Decode(&scanStatus)
    json.NewEncoder(w).Encode(scanStatus[:5])
}

// Trigger scan
func TriggerScan(w http.ResponseWriter, r *http.Request) {
    // Trigger Scan
    var queueID  string
    var response Response
    var retry    int = 0

    params := mux.Vars(r)
    fmt.Println("[INFO] Scan ID: " + params["scanid"])

    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    resp, err := http.Post("https://jenkins/job/"+params["scanid"]+"/build",
                           "application/x-www-form-urlencoded",
                           strings.NewReader(""))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    queueID = strings.Trim(resp.Header.Get("Location")[50:], "/")
    fmt.Println(queueID)
    fmt.Println(resp.Header.Get("Location"))

    // Get JobID
    for retry < 10 {
        resp, err = http.Get("https://jenkins/queue/item/"+queueID+"/api/json")

        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
        time.Sleep(3000 * time.Millisecond)

        fmt.Println(resp.Body)

        err = json.NewDecoder(resp.Body).Decode(&response)
        if err != nil {
            panic(err)
        }
        fmt.Println(response)

        if response.Executable.Number == 0 {
            retry += 1
        } else {
            json.NewEncoder(w).Encode(response.Executable.Number)
            break
        }
    }

    if retry >= 10 {
        json.NewEncoder(w).Encode("Fail to trigger scan...")
    }
}
