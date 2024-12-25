package controls

import (
	"chip/internal/setup"
)

// SetSpeed defines the speed percentage for the train.
func SetSpeed(speed uint32) {
	// Sanitize if required
	if speed > 100 {
		speed = 100
	}

	setup.SpeedPWM.Set(setup.SpeedChannel, setup.SpeedPWM.Top() / 100 * speed)
	setup.CurrentSpeed = speed
}

// SetDirection defines if the train should go forward (true) or backward (false).
func SetDirection(forward bool) {
	if forward {
		setup.DirectionPWM.Set(setup.DirectionChannel, 0);
	} else {
		setup.DirectionPWM.Set(setup.DirectionChannel, setup.DirectionPWM.Top());
	}
	setup.CurrentDirection = forward
} 

// SetLedState turns on the LED if up is set to true, disable it otherwise.
func SetLedState(up bool) {
	if up {
		setup.LedPWM.Set(setup.LedChannel, setup.LedPWM.Top());
	} else {
		setup.LedPWM.Set(setup.LedChannel, 0);
	}
	setup.CurrentLedState = up
}