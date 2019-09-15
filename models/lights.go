package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (d Control4Device) ClientSetOn(isOn bool) {
	fmt.Println("ClientSetOn")
}

func (d Control4Device) ClientSetBrightness(val int) {
	fmt.Println("ClientSetBrightness")

	url := fmt.Sprintf("%s/level/%d", d.BaseURL, val)
	http.Get(url)
}

func (d Control4Device) ClientGetOn() bool {
	fmt.Println("ClientGetOn")
	return true
}

func (d Control4Device) ClientGetBrightness() int {
	url := fmt.Sprintf("%s/brightness", d.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	value, _ := strconv.Atoi(bodyString)

	return value
}
