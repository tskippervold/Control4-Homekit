# Homekit bridge for Control4
A program written in Golang for enabling HomeKit control of your Control4 system.

It supports control of:
- Lights and dimmers
- Switches
- Thermostats
- Motion sensors

This program is bi-directional, it receives commands from Control4 when something changes (ie. a dimmer has changed brightness from a Control4 remote), and updates Control4 when it receives a command from Homekit.

## Running & Building
- `go run main.go` for just compiling and running the source code.
- Run `make compile` will compile executables for OSX and Linux.

### Dependencies
Go modules is used. `go.mod` lists dependencies.
- Run `go mod tidy` to ensure you have the neccessary modules installed.

## Running "in production"

### Second device, ie. Raspberry Pi.
_wip_

### Directly on Control4 device.
_wip_. This is "hacky" and requires you to tamper with your Control4 controller using `root` priveleges.

Transfer the compiled Linux executable to the Control4 controller using `scp`:
- `scp ./bin/linux-hap-bridge root@<CONTROLLER_IP>:/hap/`. Control4 root user password: `t0talc0ntr0l4!`

_wip: Make sure the program is started when the C4 controller boots after director._