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
	"net/http"
	"os"
	"pingservices/version"
	"strconv"

	"github.com/fatih/color"
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
func init() {
	// instanciate a new logger
	var log = logrus.New()
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&srv, "server", true, "run in server mode")
	flag.BoolVar(&srv, "s", true, "run in server mode (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(TETRACON, version.PingVersion()))
		flag.PrintDefaults()
	}

	flag.Parse()
	if vrsn {
		fmt.Printf("flag version %s\n", version.PingVersion())
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		arg = flag.Args()[0]
	}

	if arg == "help" {
		usageAndExit("", 0)
	}

	if arg == "version" {
		fmt.Printf("flag version %s\n", version.PingVersion())
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
	//	Read json configuration file
	//
	// parse the arg
	//arg := flag.Args()[0]
	//
	// check both possible arguments
	if flag.NArg() < 1 {
		showStartup(port)
		color.Unset()
		router := mux.NewRouter()

		router.HandleFunc("/getnumberofsigned", ActionHandler).Methods("GET")

		//router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

		//
		err := http.ListenAndServe(":"+strconv.Itoa(port), router)
		if err != nil {
			fmt.Printf("ListenAndServer Error: %s", err.Error())
			logrus.Fatal(err)
		}
	}
}

// Action handler
func ActionHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Hostname: %s", GetHostname())
}

//
//	Get hostname of running server
//
func GetHostname() string {

	hostname, err1 := os.Hostname()
	if err1 == nil {
		//log.Printf("Hostname: %s", hostname)
		//fmt.Println("Hostname: ", hostname)
	}
	return hostname
}

// showStartup
func showStartup(port int) {
	color.Set(color.FgHiGreen)
	fmt.Printf("Selma metrics services is started on server: ")
	color.Set(color.FgHiWhite)
	fmt.Printf("%s", GetHostname())
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
