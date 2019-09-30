package main

import (
	"fmt"

	"github.com/tskippervold/control4-homekit/utils"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/tskippervold/control4-homekit/control4"
	"github.com/tskippervold/control4-homekit/models"
)

func main() {
	utils.ConfigureReportingEngine()

	c4Devices, err := models.LoadControl4Devices("./_c4devices.json")
	if err != nil {
		utils.ReportSync(err)
		panic(err)
	}

	var accessories []*accessory.Accessory
	for _, d := range c4Devices {
		var acc *accessory.Accessory
		var err error

		switch d.Service {
		case control4.Dimmer, control4.Light:
			fmt.Println("Dimmer or light")
			acc, err = d.SetupLight()

		case control4.MotionSensor:
			fmt.Println("Motion sensor")
			acc, err = d.SetupMotionSensor()

		case control4.Thermostat:
			fmt.Println("Thermostat")
			acc, err = d.SetupThermostat()

		default:
			err = fmt.Errorf("Unsupported Control4 device: %+v", d)
		}

		accessories = append(accessories, acc)
	}

	info := accessory.Info{
		Name:         "Control4",
		Manufacturer: "Theodor Tomander Skippervold",
	}
	bridge := accessory.New(info, accessory.TypeBridge)

	config := hc.Config{
		Pin:         "40601014",
		StoragePath: "_bridge",
	}
	transport, err := hc.NewIPTransport(config, bridge, accessories...)
	if err != nil {
		panic(err)
	}

	hc.OnTermination(func() {
		<-transport.Stop()
	})

	transport.Start()
}
