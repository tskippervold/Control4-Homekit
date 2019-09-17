package control4

type Kind string

const (
	Dimmer       Kind = "Dimmer"
	Light        Kind = "Light"
	MotionSensor Kind = "Motion"
	Thermostat   Kind = "Thermostat"
)
