package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"
	"errors"
	"time"
)

func stuff(b []byte) string {
	lns := strings.Split(string(b), "Current Noise Margin")
	if len(lns) < 2 {
		fmt.Printf("%s\n", lns[1])
		return ""
	}
	lns = strings.Split(lns[1], "DSLAM Vendor Information")
	return lns[0]
}

// Holds a pair of magnitudes: one for downstream, one for upstream
type magnitude struct {
	down float64
	up float64
}

func downup(s string) magnitude {
	mag := magnitude{0.0, 0.0}
	pcs := strings.Split(s, ">");
	pcs = strings.Split(pcs[1], "\u00A0dB")
	const (
		FLOAT_64_SIZE=64 
	)
	var err error
	mag.down, err = strconv.ParseFloat(pcs[0], FLOAT_64_SIZE)
	if err != nil {
		err = errors.New("Error parsing downstream number")
		return mag
	}
	pcs = strings.Split(pcs[1], ")\u00A0")
	mag.up, err = strconv.ParseFloat(pcs[1], FLOAT_64_SIZE)
	if err != nil {
		err = errors.New("Error parsing upstream number")
	}
	return mag
}

type currents struct {
	noiseMargin magnitude
	attenuation magnitude
	outputPower magnitude
}

func numbers(s string) (curr currents) {
	lns := strings.Split(s, "\n")
	curr.noiseMargin = downup(lns[1])
	curr.attenuation = downup(lns[5])
	curr.outputPower = downup(lns[9])
	return
}

func pollrouter() {
	res, err := http.Get("http://192.168.1.254/xslt?PAGE=B02&THISPAGE=B01&NEXTPAGE=B02")
	if err != nil {
		log.Fatal(err)
	}
	page, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := stuff(page)
	current := numbers(s)
	fmt.Printf("%s\t(%6.3f,%6.3f)\t(%6.3f,%6.3f)\t(%6.3f,%6.3f)\n", time.Now(),
		current.noiseMargin.down, current.noiseMargin.up,
		current.attenuation.down, current.attenuation.up,
		current.outputPower.down, current.outputPower.up)
}

func main() {

	for {
		pollrouter()
		time.Sleep(time.Minute)
	}
}

