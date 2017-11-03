//
package loadmetrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//
type Timeset struct {
	PersOrgnr   string `json:"persorgnr"`
	PointInTime string `json:"pointintime"`
	Stage       string `json:"stage"`
}

//
func LoadSelmaMetrics(file string) ([]Timeset, error) {

	var timesets []Timeset
	//timesetFile, err := os.Open(file)
	timesetFile, err := ioutil.ReadFile(file)
	//defer timesetFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err1 := json.Unmarshal(timesetFile, &timesets)
	return timesets, err1
}
