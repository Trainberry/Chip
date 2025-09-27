package main

import (
	"fmt"
	"machine"
	"test/internal/bluetooth"
	"test/internal/constants"
	"test/internal/controls"
)

var chipName string
var debug string
var reverse string

func main() {
	if debug == "true" && machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200}) != nil {
		panic("failed to configure serial port")
	}

	// chipName is defined here because ldflags cannot set value in internal packages.
	constants.ChipName = fmt.Sprintf("Trainberry::%s", chipName)

	// Announce on Bluetooth and add handler on events.
	bluetooth.StartBLE()

	// Turn on the light to indicate that train is ready to operate
	controls.SetLight([]byte("on"))

	reverseBool := false
	if reverse == "true" {
		reverseBool = true
	}

	// Loop over events. As the chan is blocking if there's nothing to do, we don't have to do an awful trick with
	// time.Sleep or something similar.
	for {
		event := <-constants.EventChannel
		switch event.Operation {
		case constants.LightOperation:
			controls.SetLight(event.Data)
			break
		case constants.SpeedOperation:
			controls.SetSpeed(event.Data, reverseBool)
			break
		}
	}
}
