// Packages: loginmetrics.go
// Author: Tetracon AB, 2017
// Developer: Kjell Almgren
//
package agedmetrics

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"selmametrics/loadmetrics"
	"selmametrics/utility"
	"strconv"
	"time"
)

// Timeset
type Timeset struct {
	PersOrgnr   string `json:"persorgnr"`
	PointInTime string `json:"pointintime"`
	Stage       string `json:"stage"`
}

var filename = "./timesets.json"

// GetAverageAgeHandler
func GetAverageAgeHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetrics(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	//io.WriteString(w, `{"NumberOfLogin": `+fmt.Sprintf("%d", numberoflogin)+`}`)
	io.WriteString(w, `[{"text": "upper_50", "value": `+fmt.Sprintf("%d", averageage)+`}]`)
}

// GetAverageAgeSearchHandler
func GetAverageAgeSearchHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetrics(filename)
	//timenumber := 1450754160000
	timenumber := time.Now().UnixNano()
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `[{"target": "upper_50", "datapoints":[ [`+fmt.Sprintf("%d", averageage)+`,`+fmt.Sprintf("%d", timenumber)+`]]`+` }]`)
}

// GetAverageAgeQueryHandler
func GetAverageAgeQueryHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetrics(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `[{"target": "upper_50", "value": `+fmt.Sprintf("%d", averageage)+`}]`)
}

// GetAverageAgeAHandler
func GetAverageAgeAHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetricsA(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	//io.WriteString(w, `{"NumberOfLogin": `+fmt.Sprintf("%d", numberoflogin)+`}`)
	io.WriteString(w, `[{"text": "upper_50", "value": `+fmt.Sprintf("%d", averageage)+`}]`)
}

// GetAverageAgeAQueryHandler
func GetAverageAgeAQueryHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetricsA(filename)
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `[{"target": "upper_50", "value": `+fmt.Sprintf("%d", averageage)+`}]`)
}

// GetAverageAgeASearchHandler
func GetAverageAgeASearchHandler(w http.ResponseWriter, r *http.Request) {

	averageage := loadMetricsA(filename)
	//timenumber := 1450754160000
	timenumber := time.Now().UnixNano()
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `[{"target": "upper_50", "datapoints":[ [`+fmt.Sprintf("%d", averageage)+`,`+fmt.Sprintf("%d", timenumber)+`]]`+` }]`)
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
func loadMetrics(filename string) int64 {

	averageage := int64(0)
	tslists, err := loadmetrics.LoadSelmaMetrics(filename)
	if err != nil {
		fmt.Printf("JSON unmarshal Error: %s\r\n", err)
		fmt.Printf("Check %s for JSON typing error\r\n", "./timesets.json")
		os.Exit(1)
	}
	age := int64(0)
	count := 0
	for key := range tslists {
		if (tslists[key].Stage) == "S" {
			sage := tslists[key].PersOrgnr[0:4]
			//fmt.Printf("%s\r\n", sage)
			iage, err := strconv.Atoi(sage)
			if err != nil {
				panic(err)
				os.Exit(1)
			}
			//fage := fmt.Sprintf("%d", iage)
			age += int64(2017) - int64(iage)
			count++
			fmt.Printf("2017-%d=%d\r\n", iage, age)
		}
	}
	averageage = age / int64(count)
	fmt.Printf("\r\n average age: %d\r\n", averageage)
	return averageage
}

// loadMetricsA()  - private function to this package
func loadMetricsA(filename string) int64 {

	averageage := int64(0)
	tslists, err := loadmetrics.LoadSelmaMetrics(filename)
	//tslists1 := removeDuplicates1(tslists)
	if err != nil {
		fmt.Printf("JSON unmarshal Error: %s\r\n", err)
		fmt.Printf("Check %s for JSON typing error\r\n", "./timesets.json")
		os.Exit(1)
	}
	age := int64(0)
	count := 0
	for key := range tslists {
		if (tslists[key].Stage) == "A" {
			sage := tslists[key].PersOrgnr[0:4]
			//fmt.Printf("%s\r\n", sage)
			iage, err := strconv.Atoi(sage)
			if err != nil {
				panic(err)
				os.Exit(1)
			}
			//fage := fmt.Sprintf("%d", iage)
			age += int64(2017) - int64(iage)
			count++
			fmt.Printf("2017-%d=%d\r\n", iage, age)
		}
	}
	averageage = age / int64(count)
	fmt.Printf("\r\n average age: %d\r\n", averageage)
	return averageage
}

// removeDuplicates1
func removeDuplicates1(lists []struct {
	PersOrgnr   string `json:"persorgnr"`
	PointInTime string `json:"pointintime"`
	Stage       string `json:"stage"`
}) int {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	// Copy to PersOrgnr to a string array
	i := 0
	persons := []string{}
	for _, w := range lists {
		persons = append(persons, w.PersOrgnr)
		i++
		//fmt.Printf("( %d : (%s)\r\n", i, w.PersOrgnr)
	}
	// set all map items to true
	for person := range persons {
		encountered[persons[person]] = true
	}
	// prepare result, number of unique customers
	result := []string{}
	j := 0
	for key, _ := range encountered {
		result = append(result, key)
		j++
	}
	return j
}
