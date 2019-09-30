package control4

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tskippervold/control4-homekit/utils"
)

var outboundIP string

func SetHAPBridgeIP(deviceURL string) error {
	if outboundIP == "" {
		fmt.Println("Need to create outboundIP.")
		ip, err := findOutboundIP()
		if err != nil {
			return err
		}

		outboundIP = ip
		fmt.Printf("OutboundIP is: %s.\n", outboundIP)
	}

	url := fmt.Sprintf("%s/SetApplianceIP/%s", deviceURL, outboundIP)
	_, err := http.Get(url)
	return err
}

func StartServer(deviceID int, requestHandler func(w http.ResponseWriter, r *http.Request)) error {
	port := fmt.Sprintf(":%d", deviceID*1+10000)

	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	server := &http.Server{Addr: port, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			utils.ReportSync(err)
			panic(err)
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, os.Kill)
	signal.Notify(stop, syscall.SIGTERM)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	fmt.Printf("Server stopped on port %s\n", port)
	return nil
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func findOutboundIP() (string, error) {
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
