package setup

import (
	"github.com/soypat/cyw43439"
	"log/slog"
	"machine"
	"time"
	"chip/internal/configuration"
)

func init() {

	// Configure watchdog and sanity
    machine.Watchdog.Configure(machine.WatchdogConfig{
		TimeoutMillis: machine.WatchdogMaxTimeout, // Approx 8s. See https://tinygo.org/docs/reference/microcontrollers/machine/pico 
	})

	machine.Watchdog.Start()

	var err error

	//--- Init RPi-Pico configuration
	Logger = slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	machine.Watchdog.Update()

	dev := cyw43439.NewPicoWDevice()
	cfg := cyw43439.DefaultWifiConfig()

	if configuration.Debug {
		cfg.Logger = Logger
	}

	machine.Watchdog.Update()

	for err = dev.Init(cfg) ; err != nil ; {
		Logger.Error("Error on init: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	machine.Watchdog.Update()

	//--- Setup WiFi

	wiFiConfig := WiFiSetupConfig{
		Hostname: configuration.TrainName,
		RequestedIP: configuration.TrainIP,
		WiFiSSID: configuration.WiFiSSID,
		WiFiPassword: configuration.WiFiPassword,
		TCPPorts: 4,
		UDPPorts: 10,
	}


	for PortStack, err = setupWiFi(dev, wiFiConfig) ; err != nil ; {
		Logger.Error("Error on WiFi login: " + err.Error() + ". Will try again in 5 seconds")
		time.Sleep(time.Second * 5)
	}

	machine.Watchdog.Update()

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
	}

	machine.Watchdog.Update()

	// Speed
	SpeedPWM.Configure(machine.PWMConfig{
		Period: 1e9/500,
	})
	
	SpeedChannel, err = SpeedPWM.Channel(SpeedPin)
	if err != nil {
		Logger.Error("Error on SpeedChannel init: " + err.Error())
	}

	machine.Watchdog.Update()

	// Led
	LedPWM.Configure(machine.PWMConfig{
		Period: 1e9/500,
	})
	
	LedChannel, err = LedPWM.Channel(LedPin)
	if err != nil {
		Logger.Error("Error on LedChannel init: " + err.Error())
	}

	machine.Watchdog.Update()

	go sanity()

}