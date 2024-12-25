package setup

import (
	"github.com/soypat/cyw43439"
	"log/slog"
	"machine"
	"chip/internal/configuration"
)

func init() {
	//--- Init RPi-Pico configuration
	Logger = slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	dev := cyw43439.NewPicoWDevice()
	cfg := cyw43439.DefaultWifiConfig()
	cfg.Logger = Logger
	err := dev.Init(cfg)

	if err != nil {
		Logger.Error("Error on init: " + err.Error())
		panic(err)
	}

	//--- Setup WiFi
	PortStack, err = setupWiFi(dev, WiFiSetupConfig{
		Hostname: configuration.TrainName,
		RequestedIP: configuration.TrainIP,
		WiFiSSID: configuration.WiFiSSID,
		WiFiPassword: configuration.WiFiPassword,
		TCPPorts: 4,
		UDPPorts: 10,
	})
	if err != nil {
		Logger.Error("Error on WiFi login: " + err.Error())
		panic(err)
	}

	//--- Set built-in LED on to have a live notification
	dev.GPIOSet(0, true)

	//--- Configure our PWM signals

	// Direction
	DirectionPWM.Configure(machine.PWMConfig{
		Period: 1e9/500,
	})

	DirectionChannel, err = DirectionPWM.Channel(DirectionPin)
	if err != nil {
		Logger.Error("Error on DirectionChannel init: " + err.Error())
		return
	}

	// Speed
	SpeedPWM.Configure(machine.PWMConfig{
		Period: 1e9/500,
	})
	
	SpeedChannel, err = SpeedPWM.Channel(SpeedPin)
	if err != nil {
		Logger.Error("Error on SpeedChannel init: " + err.Error())
		return
	}

	// Led
	LedPWM.Configure(machine.PWMConfig{
		Period: 1e9/500,
	})
	
	LedChannel, err = LedPWM.Channel(LedPin)
	if err != nil {
		Logger.Error("Error on LedChannel init: " + err.Error())
		return
	}

}