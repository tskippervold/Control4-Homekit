package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tskippervold/control4-homekit/control4"

	"github.com/brutella/hc/accessory"
)

type Control4Light struct {
	Control4Device
	HAPDevice *accessory.Lightbulb
}

func (d Control4Light) ClientIdentify() {
	fmt.Printf("\nIdentify:\n%+v\n", d)
}

func (d Control4Light) ClientSetOn(isOn bool) {
	if isOn {
		brightness := d.ClientGetBrightness()
		if brightness > 0 {
			d.ClientSetBrightness(brightness)
		} else {
			d.ClientSetBrightness(100)
		}
	} else {
		d.ClientSetBrightness(0)
	}
}

func (d Control4Light) ClientSetBrightness(val int) {
	url := fmt.Sprintf("%s/level/%d", d.BaseURL, val)
	http.Get(url)
}

func (d Control4Light) ClientGetOn() bool {
	brightness := d.ClientGetBrightness()
	return brightness > 0
}

func (d Control4Light) ClientGetBrightness() int {
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

func (d Control4Light) RemoteUpdatedHandler(w http.ResponseWriter, r *http.Request) {
	prop, val := control4.PropertyFrom(r)

	switch prop {
	case control4.Power:
		isOn := val == "1"
		d.HAPDevice.Lightbulb.On.SetValue(isOn)

	case control4.Brightness:
		brightness, _ := strconv.Atoi(val)
		d.HAPDevice.Lightbulb.Brightness.SetValue(brightness)

	}

	w.WriteHeader(http.StatusOK)
}
