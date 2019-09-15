package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/brutella/hc/accessory"
)

type Control4Device struct {
	HAPAccessory *accessory.Accessory

	Accessory          string `json:"accessory"`
	Name               string `json:"name"`
	Service            string `json:"service"`
	BaseURL            string `json:"base_url"`
	HasLevelControl    string `json:"has_level_control"`
	SwitchHandling     string `json:"switchHandling"`
	BrightnessHandling string `json:"brightnessHandling"`
	RefreshInterval    int64  `json:"refresh_interval"`
	Manufacturer       string `json:"manufacturer"`
	Model              string `json:"model"`
}

func LoadControl4Devices(pathToJSON string) ([]*Control4Device, error) {
	jsonFile, err := os.Open(pathToJSON)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	var devices []*Control4Device

	err = json.Unmarshal(bytes, &devices)
	return devices, err
}

func (d Control4Device) HomekitAccessory() (*accessory.Accessory, error) {

	info := accessory.Info{
		Name: d.Name,
	}

	switch strings.ToLower(d.Service) {
	case "dimmer", "light":
		ac := accessory.NewLightbulb(info)

		ac.Accessory.OnIdentify(func() {
			fmt.Println("Identify")
		})

		ac.Lightbulb.Brightness.OnValueRemoteUpdate(d.ClientSetBrightness)
		ac.Lightbulb.Brightness.OnValueRemoteGet(d.ClientGetBrightness)

		ac.Lightbulb.On.OnValueRemoteUpdate(d.ClientSetOn)
		ac.Lightbulb.On.OnValueRemoteGet(d.ClientGetOn)

		return ac.Accessory, nil

	default:
		return nil, fmt.Errorf("Unsupported Control4 device: %+v", d)
	}
}

func (d Control4Device) ID() (int, error) {
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

func (d Control4Device) SetBridgeIP(ip string) error {
	url := fmt.Sprintf("%s/SetApplianceIP/%s", d.BaseURL, ip)
	_, err := http.Get(url)
	return err
}

func (d Control4Device) StartServer() {
	id, err := d.ID()
	if err != nil {
		panic(err)
	}
	port := fmt.Sprintf(":%d", id*1+10000)

	mux := http.NewServeMux()
	mux.HandleFunc("/", d.requestHandler)

	server := &http.Server{Addr: port, Handler: mux}

	go func() {
		fmt.Printf("Starting server for: %s on port: %s\n", d.BaseURL, port)
		if err := server.ListenAndServe(); err != nil {
			// handle err
			fmt.Println(err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Server stopped on port %s\n", port)
}

type Control4ProperyType string

const (
	OnOff Control4ProperyType = "1000"
	Value Control4ProperyType = "1001"
)

func (d Control4Device) requestHandler(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")

	property := parts[1]
	value := parts[2]

	switch Control4ProperyType(property) {
	case OnOff:

	case Value:

	default:
		fmt.Printf("Unhandled c4 property type: %s, withValue: %s\n", property, value)
	}
}
