package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"
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

func downup(s string) (down float64, up float64) {
	pcs := strings.Split(s, ">");
	pcs = strings.Split(pcs[1], "\u00A0dB")
	const (
		FLOAT_64_SIZE=64 
	)
	var err error
	down, err = strconv.ParseFloat(pcs[0], FLOAT_64_SIZE)
	if (err != nil) {
		down = 0
		up = 0
		return
	}
	pcs = strings.Split(pcs[1], ")\u00A0")
	up, err = strconv.ParseFloat(pcs[1], FLOAT_64_SIZE)
	if (err != nil) {
		up = 0
	}
	return
}

func numbers(s string) (cnmd float64, cnmu float64, cad float64, cau float64, copd float64, copu float64) {
	lns := strings.Split(s, "\n")
	cnmd, cnmu = downup(lns[1])
	cad, cau = downup(lns[5])
	copd, copu = downup(lns[9])
	return
}

func main() {
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

	currentNoiseMarginDown, currentNoiseMarginUp, currentAttenuationDown, currentAttenuationUp, currentOutputPowerDown, currentOutputPowerUp := numbers(s)
	fmt.Printf(" (%f,%f) (%f,%f) (%f,%f)\n",
		currentNoiseMarginDown, currentNoiseMarginUp,
		currentAttenuationDown, currentAttenuationUp,
		currentOutputPowerDown, currentOutputPowerUp)
}
