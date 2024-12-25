package setup

import (
	"github.com/soypat/seqs/stacks"
	"machine"
	"log/slog"
)


//--- Constants
var DirectionPWM = machine.PWM1
const DirectionPin = machine.GPIO18

var SpeedPWM = machine.PWM1
const SpeedPin = machine.GPIO19

var LedPWM = machine.PWM5
const LedPin = machine.GPIO26

//--- Variables acessible from anywhere
var SpeedChannel uint8
var DirectionChannel uint8
var LedChannel uint8

var PortStack *stacks.PortStack

//--- Store state
var CurrentDirection bool
var CurrentLedState bool
var CurrentSpeed uint32

//--- Shared Logger
var Logger *slog.Logger