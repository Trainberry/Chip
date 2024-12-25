package main

import (
	"chip/internal/controls"
	"chip/internal/http/client"
	"chip/internal/http/server"
)


func main() {
	// Set default state : enable LED, speed 0 and forward direction
	controls.SetSpeed(0)
	controls.SetDirection(true)
	controls.SetLedState(true)
	
	// Register to central server, let REST API start during this step...
	go client.RegisterTrain()

	// Set-up web server
	server.StartServer()
}


