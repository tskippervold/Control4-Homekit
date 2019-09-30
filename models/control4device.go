package models

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/tskippervold/control4-homekit/customAccessory"

	"github.com/brutella/hc/accessory"

	"github.com/tskippervold/control4-homekit/control4"
)

type Control4Device struct {
	Accessory          string        `json:"accessory"`
	Name               string        `json:"name"`
	Service            control4.Kind `json:"service"`
	BaseURL            string        `json:"base_url"`
	HasLevelControl    string        `json:"has_level_control"`
	SwitchHandling     string        `json:"switchHandling"`
	BrightnessHandling string        `json:"brightnessHandling"`
	RefreshInterval    int64         `json:"refresh_interval"`
	Manufacturer       string        `json:"manufacturer"`
	Model              string        `json:"model"`

	DeviceID int
}

func LoadControl4Devices(pathToJSON string) ([]Control4Device, error) {
	jsonFile, err := os.Open(pathToJSON)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var devices []Control4Device
	err = json.Unmarshal(bytes, &devices)

	for i, d := range devices {
		id, _ := d.deviceID() // TODO: Handle error
		d.DeviceID = id
		devices[i] = d
	}

	return devices, err
}

func (d *Control4Device) deviceID() (int, error) {
	url, err := url.Parse(d.BaseURL)
	if err != nil {
		return -1, err
	}

	parts := strings.Split(url.Path, "/")
	i, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1, err
	}

	return i, nil
}

func (d *Control4Device) hapAccessoryInfo() accessory.Info {
	return accessory.Info{
		Name:         d.Name,
		Manufacturer: d.Manufacturer,
		Model:        d.Model,
	}
}

func (d *Control4Device) SetupLight() (*accessory.Accessory, error) {
	light := Control4Light{*d, accessory.NewLightbulb(d.hapAccessoryInfo())}

	light.HAPDevice.OnIdentify(light.ClientIdentify)

	light.HAPDevice.Lightbulb.Brightness.OnValueRemoteUpdate(light.ClientSetBrightness)
	light.HAPDevice.Lightbulb.Brightness.OnValueRemoteGet(light.ClientGetBrightness)

	light.HAPDevice.Lightbulb.On.OnValueRemoteUpdate(light.ClientSetOn)
	light.HAPDevice.Lightbulb.On.OnValueRemoteGet(light.ClientGetOn)

	err := control4.SetHAPBridgeIP(d.BaseURL)
	if err != nil {
		return nil, err
	}

	go control4.StartServer(light.DeviceID, light.RemoteUpdatedHandler)

	return light.HAPDevice.Accessory, nil
}

func (d *Control4Device) SetupMotionSensor() (*accessory.Accessory, error) {
	sensor := Control4MotionSensor{*d, customAccessory.NewMotionSensor(d.hapAccessoryInfo())}
	sensor.HAPDevice.IsInverted = true

	sensor.HAPDevice.Accessory.OnIdentify(sensor.ClientIdentify)
	sensor.HAPDevice.MotionSensor.MotionDetected.OnValueRemoteGet(sensor.ClientGetMotion)

	err := control4.SetHAPBridgeIP(d.BaseURL)
	if err != nil {
		return nil, err
	}

	go control4.StartServer(sensor.DeviceID, sensor.RemoteUpdatedHandler)

	return sensor.HAPDevice.Accessory, nil
}

func (d *Control4Device) SetupThermostat() (*accessory.Accessory, error) {
	return nil, nil
}
