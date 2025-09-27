package constants

import (
	"machine"
	"test/internal/structures"
)

// The EventChannel will store all items received from BLE, and will be consumed on CPU idle.
var EventChannel = make(chan structures.Event, 5)

var ChipName string

var ForwardPin = machine.D3
var ForwardPWM = machine.PWM0

var BackwardPin = machine.D4
var BackwardPWM = machine.PWM1

var LightPin = machine.D0
var LightPWM = machine.PWM2

// Enums to give a better visual to Event.Operation field
const (
	SpeedOperation = iota
	LightOperation
)
