package customAccessory

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type MotionSensor struct {
	*accessory.Accessory

	MotionSensor *service.MotionSensor

	IsInverted bool
}

func NewMotionSensor(info accessory.Info) *MotionSensor {
	acc := MotionSensor{IsInverted: false}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)
	acc.MotionSensor = service.NewMotionSensor()

	acc.AddService(acc.MotionSensor.Service)

	return &acc
}
