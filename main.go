package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/tskippervold/control4-homekit/models"
)

func main() {
	outboundIP, err := getOutboundIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("OutboundIP is: %s\n", outboundIP)

	c4Devices, err := models.LoadControl4Devices("_c4devices.json")
	if err != nil {
		panic(err)
	}

	var homekitDevices []*accessory.Accessory
	for _, d := range c4Devices {
		ipOut, err := getOutboundIP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go d.StartServer()
		if err := d.SetBridgeIP(ipOut); err != nil {
			fmt.Println(err)
			continue
		}

		acc, err := d.HomekitAccessory()
		if err != nil {
			fmt.Println(err)
			continue
		}
		homekitDevices = append(homekitDevices, acc)
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
	transport, err := hc.NewIPTransport(config, bridge, homekitDevices...)
	if err != nil {
		panic(err)
	}

	hc.OnTermination(func() {
		<-transport.Stop()
	})

	transport.Start()
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func getOutboundIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
