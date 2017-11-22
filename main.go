/*
Services: selmametrics
	Author: Kjell Osse Almgren, Tetracon AB
	Date: 2017-11-02
	Description: Service to feed grafana with business metrics, just for test purpose
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
	"net/http"
	"os"
	"selmametrics/agedmetrics"
	"selmametrics/loadmetrics"
	"selmametrics/loginmetrics"
	"selmametrics/signedmetrics"
	"time"

	"selmametrics/utility"
	"selmametrics/version"
	"strconv"

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

var filename = "./timesets.json"

// init
func init() {
	// instanciate a new logger
	var log = logrus.New()
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&srv, "server", true, "run in server mode")
	flag.BoolVar(&srv, "s", true, "run in server mode (shorthand)")
	flag.BoolVar(&srv, "generate", true, "genrate data to processing 3D")
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
	if arg == "generate" {
		fmt.Printf("Start genrate data for processing 3D")
		generateProcessingMetrics()
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
	//	Read json metrics file, this is only for test
	//
	//tslists, err := loadmetrics.LoadSelmaMetrics(filename)
	//if err != nil {
	//	fmt.Printf("JSON unmarshal Error: %s\r\n", err)
	//	fmt.Printf("Check %s for JSON typing error\r\n", "./timesets.json")
	//	os.Exit(1)
	//}
	//
	//for key := range tslists {
	//	fmt.Printf("PersOrgNr: %s\r\n", tslists[key].PersOrgnr)
	//	fmt.Printf("PointInTime: %s\r\n", tslists[key].PointInTime)
	//	fmt.Printf("Stage: %s\r\n", tslists[key].Stage)
	//}
	// parse the arg
	//arg := flag.Args()[0]
	//
	// check both possible arguments
	if flag.NArg() < 1 {
		showStartup(port)
		color.Unset()
		router := mux.NewRouter()
		//
		// for ../search and ../query handler grafana want us to use method "POST"
		//
		// Package: signedmetrics.go
		router.HandleFunc("/health-check", signedmetrics.HealthCheckHandler).Methods("GET")
		router.HandleFunc("/getnumberofsigned", signedmetrics.GetNumberOfSignedHandler).Methods("GET")
		router.HandleFunc("/getnumberofsigned/search", signedmetrics.GetNumberOfSignedSearchHandler).Methods("POST")
		router.HandleFunc("/getnumberofsigned/query", signedmetrics.GetNumberOfSignedSearchHandler).Methods("POST")
		//
		// Package: loginmetrics.go
		router.HandleFunc("/health-check", loginmetrics.HealthCheck1Handler).Methods("GET")
		router.HandleFunc("/getnumberoflogin", loginmetrics.GetNumberOfLoginHandler).Methods("GET")
		router.HandleFunc("/getnumberoflogin/search", loginmetrics.GetNumberOfLoginSearchHandler).Methods("POST")
		router.HandleFunc("/getnumberoflogin/query", loginmetrics.GetNumberOfLoginSearchHandler).Methods("POST")
		//
		router.HandleFunc("/health-check", agedmetrics.HealthCheck1Handler).Methods("GET")
		router.HandleFunc("/getaverageage", agedmetrics.GetAverageAgeHandler).Methods("GET")
		router.HandleFunc("/getaverageage/search", agedmetrics.GetAverageAgeSearchHandler).Methods("POST")
		router.HandleFunc("/getaverageage/query", agedmetrics.GetAverageAgeSearchHandler).Methods("POST")
		//
		router.HandleFunc("/getaverageagea", agedmetrics.GetAverageAgeAHandler).Methods("GET")
		router.HandleFunc("/getaverageagea/search", agedmetrics.GetAverageAgeASearchHandler).Methods("POST")
		router.HandleFunc("/getaverageagea/query", agedmetrics.GetAverageAgeASearchHandler).Methods("POST")
		// To support CORS: Cross Domain Resource Sharing, AllowedOrigins should be a domain
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

// showStartup
func showStartup(port int) {

	color.Set(color.FgHiGreen)
	fmt.Printf("Selma metrics services (")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", version.SelmaMetricsVersion())
	color.Set(color.FgHiGreen)
	fmt.Printf(") Selma metrics API-services is started on server: ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", utility.GetHostname())
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

// generateProcessingMetrics
func generateProcessingMetrics() {

	layOut := "2006-01-02 15:04:05"
	//layout := "yyyy-MM-dd hh:mm:ss"
	t := time.Now()
	currentyear := t.Year()
	timesets, err := loadmetrics.LoadSelmaMetrics("./timesets.json")
	if err != nil {
		panic(err)
		fmt.Println("JSON")
		os.Exit(1)
	}
	for _, timeset := range timesets {
		//fmt.Printf("%s %s %s\r\n", timeset.PersOrgnr, timeset.PointInTime, timeset.Stage)
		sage := timeset.PersOrgnr[0:4]
		gendertemp := timeset.PersOrgnr[10:]
		sgender := gendertemp[0:1]
		gender := findGender(sgender)
		iage, err1 := strconv.Atoi(sage)
		if err1 != nil {
			panic(err1)
			os.Exit(1)
		}
		// Writeline
		//stime = timesets[key].PointInTime
		ttime, err2 := time.Parse(layOut, timeset.PointInTime)
		if err2 != nil {
			panic(err2)
			os.Exit(1)
		} else {
			_, week := ttime.ISOWeek()
			day := ttime.Weekday() - 1 // swedish week day
			hour := ttime.Hour()
			age := int64(0)
			age = int64(currentyear) - int64(iage)
			stage := timeset.Stage
			if day == -1 {
				day = 6
			}
			fmt.Printf("%d,%d,%d,%d,\"%s\",\"%s\"\r\n", week, day, hour, age, stage, gender)
		}
	}
}

// findGender
func findGender(sgender string) string {

	ret := ""
	igender, err := strconv.Atoi(sgender)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	if igender%2 == 0 {
		ret = "W"
	} else {
		ret = "M"
	}
	return ret
}
