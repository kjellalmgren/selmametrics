// Package: Utility.go
// Author: Tetracon AB, 2017
// Developer: Kjell Almgren
//
package utility

import "os"

//
// GetHostname
func GetHostname() string {

	hostname, err1 := os.Hostname()
	if err1 == nil {
		//log.Printf("Hostname: %s", hostname)
		//fmt.Println("Hostname: ", hostname)
	}
	return hostname
}
