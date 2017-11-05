/*
Services: selmametrics
	Author: Kjell Osse Almgren, Tetracon AB
	Date: 2017-11-02
	Description: Service to feed grafana with metrics, test purpose
	Architecture:
	win32: GOOS=windows GOARCH=386 go build -v
	win64: GOOS=windows GOARCH=amd64 go build -v
	arm64: GOOS=linux GOARCH=arm64 go build -v
	arm: GOOS=linux GOARCH=arm go build -v
	exprimental: GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags '-w -s' -a -installsuffix cgo -o pingservices
	expriemntal: CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -tags pingservices -ldflags '-w'
	exprimental: GOOS=linux GOARCH=arm64 go build -a --ldflags 'extldflags "-static"' -tags pingservices -installsuffix pingservices .
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"selmametrics/loadmetrics"
	"selmametrics/utility"
	"selmametrics/version"
	"strconv"
	"time"

	"github.com/fatih/color"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Banner text
const (
	// TETRACON banner
	TETRACON = `
_________    __
|__    __|   | |
   |  |  ___ | |_   ____  ___   ___ ___  _ __ 
   |  | / _ \|  _| /  __|/ _ \ / __/ _ \| '_ \
   |  | \ __/| |_  | |  | (_| | (_| (_) | | | | 
   |__| \___| \__| |_|   \__,_|\___\___/|_| |_| 
version: %s
`
)

//
//
var (
	srv  bool
	vrsn bool
)

var (
	arg string
)

//
type Timeset struct {
	PersOrgnr   string `json:"persorgnr"`
	PointInTime string `json:"pointintime"`
	Stage       string `json:"stage"`
}

// init
func init() {
	// instanciate a new logger
	var log = logrus.New()
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&srv, "server", true, "run in server mode")
	flag.BoolVar(&srv, "s", true, "run in server mode (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(TETRACON, version.SelmaMetricsVersion()))
		flag.PrintDefaults()
	}

	flag.Parse()
	if vrsn {
		fmt.Printf("Selma Metrics Version %s\n", version.SelmaMetricsVersion())
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		arg = flag.Args()[0]
	}

	if arg == "help" {
		usageAndExit("", 0)
	}

	if arg == "version" {
		fmt.Printf("flag version %s\n", version.SelmaMetricsVersion())
		os.Exit(0)
	}
	//log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) // default

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	log.Level = logrus.DebugLevel
}

// here we go
func main() {

	port := 8000
	//
	//	Read json metrics file
	//
	tslists, err := loadmetrics.LoadSelmaMetrics("./timesets.json")
	if err != nil {
		fmt.Printf("JSON unmarshal Error: %s\r\n", err)
		fmt.Printf("Check %s for JSON typing error\r\n", "./timesets.json")
		os.Exit(1)
	}
	//
	for key := range tslists {
		fmt.Printf("PersOrgNr: %s\r\n", tslists[key].PersOrgnr)
		fmt.Printf("PointInTime: %s\r\n", tslists[key].PointInTime)
		fmt.Printf("Stage: %s\r\n", tslists[key].Stage)
	}
	// parse the arg
	//arg := flag.Args()[0]
	//
	// check both possible arguments
	if flag.NArg() < 1 {
		showStartup(port)
		color.Unset()
		router := mux.NewRouter()
		router.HandleFunc("/health-check", HealthCheckHandler).Methods("GET")
		router.HandleFunc("/getnumberofsigned", GetNumberOfSignedHandler).Methods("GET")
		router.HandleFunc("/getnumberofsigned/search", GetNumberOfSignedSearchHandler).Methods("POST")
		router.HandleFunc("/getnumberofsigned/query", GetNumberOfSignedSearchHandler).Methods("POST")

		//router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

		//
		/*
			handlers.AllowedMethods([]string{"GET, POST, OPTIONS, PUT, DELETE"}),
					handlers.AllowedOrigins([]string{"*"}),
					handlers.AllowedHeaders([]string{"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With"}
		*/
		//err := http.ListenAndServe(":"+strconv.Itoa(port), router)

		headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "X-Requested-With"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS"})

		err := http.ListenAndServe(":"+strconv.Itoa(port),
			handlers.LoggingHandler(os.Stdout, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
		if err != nil {
			fmt.Printf("ListenAndServer Error: %s", err.Error())
			logrus.Fatal(err)
		}
	}
}

// GetNumberOfSignedHandler
func GetNumberOfSignedHandler(w http.ResponseWriter, r *http.Request) {

	numberofsigned := 19
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	//io.WriteString(w, `{"NumberOfSigned": `+fmt.Sprintf("%d", numberofsigned)+`}`)
	io.WriteString(w, `[{"text": "upper_50", "value": `+fmt.Sprintf("%d", numberofsigned)+`}]`)
}

// GetNumberOfSignedSearchHandler
func GetNumberOfSignedSearchHandler(w http.ResponseWriter, r *http.Request) {

	numberofsigned := 19
	//timenumber := 1450754160000
	timenumber := time.Now().UnixNano()
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `[{"target": "upper_50", "datapoints":[ [`+fmt.Sprintf("%d", numberofsigned)+`,`+fmt.Sprintf("%d", timenumber)+`]]`+` }]`)

}

// GetNumberOfSignedQueryHandler
func GetNumberOfSignedQueryHandler(w http.ResponseWriter, r *http.Request) {

	numberofsigned := 24
	//fmt.Printf("Hostname: %s", GetHostname())
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `[{"target": "upper_50", "value": `+fmt.Sprintf("%d", numberofsigned)+`}]`)

}

// HealthCheckHandler
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
	io.WriteString(w, `{"status":`+fmt.Sprintf("%d", http.StatusOK)+`}`)
	io.WriteString(w, `{"server":`+fmt.Sprintf("%s", Utility.GetHostname())+`}`)

	fmt.Printf("Http-Status %d received\r\n", http.StatusOK)
}

// showStartup
func showStartup(port int) {

	color.Set(color.FgHiGreen)
	fmt.Printf("Selma metrics services (")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", version.SelmaMetricsVersion())
	color.Set(color.FgHiGreen)
	fmt.Printf(") Selma metrics API-services is started on server: ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", Utility.GetHostname())
	color.Set(color.FgHiGreen)
	fmt.Printf(" is listen on port ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%d", port)
	color.Set(color.FgHiGreen)
	fmt.Println(" tls")
	color.Unset()
}

// usageAndExit
func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\r\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
