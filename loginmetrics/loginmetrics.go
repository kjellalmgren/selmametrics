// Packages: loginmetrics.go
// Author: Tetracon AB, 2017
// Developer: Kjell Almgren
//
package loginmetrics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"selmametrics/loadmetrics"
	"selmametrics/utility"
	"time"
)

// Timeset
type Timeset struct {
	PersOrgnr   string `json:"persorgnr"`
	PointInTime string `json:"pointintime"`
	Stage       string `json:"stage"`
}

var filename = "./timesets.json"

// GetNumberOfLoginHandler
func GetNumberOfLoginHandler(w http.ResponseWriter, r *http.Request) {

	numberoflogin := loadMetrics(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	//io.WriteString(w, `{"NumberOfLogin": `+fmt.Sprintf("%d", numberoflogin)+`}`)
	io.WriteString(w, `[{"text": "upper_50", "value": `+fmt.Sprintf("%d", numberoflogin)+`}]`)
}

// GetNumberOfLoginSearchHandler
func GetNumberOfLoginSearchHandler(w http.ResponseWriter, r *http.Request) {

	numberoflogin := loadMetrics(filename)
	//timenumber := 1450754160000
	timenumber := time.Now().UnixNano()
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `[{"target": "upper_50", "datapoints":[ [`+fmt.Sprintf("%d", numberoflogin)+`,`+fmt.Sprintf("%d", timenumber)+`]]`+` }]`)
}

// GetNumberOfLoginQueryHandler
func GetNumberOfLoginQueryHandler(w http.ResponseWriter, r *http.Request) {

	numberoflogin := loadMetrics(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `[{"target": "upper_50", "value": `+fmt.Sprintf("%d", numberoflogin)+`}]`)
}

// HealthCheckHandler
func HealthCheck1Handler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
	io.WriteString(w, `{"status":`+fmt.Sprintf("%d", http.StatusOK)+`}`)
	io.WriteString(w, `{"server":`+fmt.Sprintf("%s", utility.GetHostname())+`}`)

	fmt.Printf("Http-Status %d received\r\n", http.StatusOK)
}

// loadMetrics()  - private function to this package
func loadMetrics(filename string) int {

	numberoflogin := 0
	tslists, err := loadmetrics.LoadSelmaMetrics(filename)
	if err != nil {
		fmt.Printf("JSON unmarshal Error: %s\r\n", err)
		fmt.Printf("Check %s for JSON typing error\r\n", "./timesets.json")
		os.Exit(1)
	}
	for key := range tslists {
		if (tslists[key].Stage) == "A" {
			numberoflogin++
		}
	}
	return numberoflogin
}
