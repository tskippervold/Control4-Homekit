package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tskippervold/control4-homekit/control4"

	"github.com/tskippervold/control4-homekit/customAccessory"
)

type Control4MotionSensor struct {
	Control4Device
	HAPDevice *customAccessory.MotionSensor
}

func (d Control4MotionSensor) ClientIdentify() {
	fmt.Printf("\nIdentify:\n%+v\n", d)
}

func (d Control4MotionSensor) ClientGetMotion() bool {
	url := fmt.Sprintf("%s/contact_state", d.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	if d.HAPDevice.IsInverted {
		return bodyString != "1"
	}

	return bodyString == "1"
}

func (d Control4MotionSensor) RemoteUpdatedHandler(w http.ResponseWriter, r *http.Request) {
	prop, val := control4.PropertyFrom(r)

	switch prop {
	case control4.Power:
		var isOn bool
		if d.HAPDevice.IsInverted {
			isOn = val != "1"
		} else {
			isOn = val == "1"
		}

		d.HAPDevice.MotionSensor.MotionDetected.SetValue(isOn)
	}

	w.WriteHeader(http.StatusOK)
}
