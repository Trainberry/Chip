package controls

import (
	"machine"
	"strconv"
	"test/internal/constants"
)

var speedForwardChannel uint8
var speedBackwardChannel uint8
var lightChannel uint8

func init() {
	// Initialize PWM with the maximum speed
	_ = constants.ForwardPWM.Configure(machine.PWMConfig{Period: 1e9 / 500})
	_ = constants.BackwardPWM.Configure(machine.PWMConfig{Period: 1e9 / 500})
	_ = constants.LightPWM.Configure(machine.PWMConfig{Period: 1e9 / 500})

	// Create and store PWM channel to communicate with them later on
	speedForwardChannel, _ = constants.ForwardPWM.Channel(constants.ForwardPin)
	speedBackwardChannel, _ = constants.BackwardPWM.Channel(constants.BackwardPin)
	lightChannel, _ = constants.LightPWM.Channel(constants.LightPin)
}

// SetSpeed sets the speed of the motor. It takes an integer comprised between -100 (full-backward) and 100 (full-forward).
func SetSpeed(rawData []byte, reverse bool) {
	speed, err := strconv.Atoi(string(rawData))
	if err != nil || speed < -100 || speed > 100 {
		return
	}

	if reverse {
		speed = -speed
	}

	unsignedSpeed := uint32(speed)

	// ALWAYS set the unused channel to 0 BEFORE setting the other one to the desired state.

	if speed == 0 {
		constants.ForwardPWM.Set(speedForwardChannel, 0)
		constants.BackwardPWM.Set(speedBackwardChannel, 0)
		return
	}

	if speed < 0 {
		constants.ForwardPWM.Set(speedForwardChannel, 0)
		constants.BackwardPWM.Set(speedBackwardChannel, constants.BackwardPWM.Top()/100*unsignedSpeed)
		return
	}

	constants.BackwardPWM.Set(speedBackwardChannel, 0)
	constants.ForwardPWM.Set(speedForwardChannel, constants.ForwardPWM.Top()/100*unsignedSpeed)
}

// SetLight turns on/off the light. The expected payload is []byte("on") to light it up, []byte("off") otherwise.
func SetLight(rawData []byte) {
	if string(rawData) == "off" {
		constants.LightPWM.Set(lightChannel, 0)
		return
	}
	constants.LightPWM.Set(lightChannel, constants.LightPWM.Top())
}
