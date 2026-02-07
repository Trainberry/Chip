package bluetooth

import (
	"test/internal/constants"
	"test/internal/structures"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// StartBLE enables Bluetooth, defines services and characteristics and advertise bluetooth to make it visible.
func StartBLE() {
	var err error

	// Activate Bluetooth
	if err = adapter.Enable(); err != nil {
		panic(err)
	}

	// Create our service and characteristics to control our train
	// See documentation for full reference (service & characteristics IDs, expected payloads and so on)
	serviceUUID := bluetooth.NewUUID([16]byte{0x97, 0x98, 0x00, 0x00})

	pingCharacteristicUUID := bluetooth.NewUUID([16]byte{0x00, 0x00, 0x00, 0x00})
	speedCharacteristicUUID := bluetooth.NewUUID([16]byte{0x10, 0x00, 0x00, 0x00})
	lightCharacteristicUUID := bluetooth.NewUUID([16]byte{0x20, 0x00, 0x00, 0x00})

	// Set flags on our characteristics: they are readable, writable and writable without acknowledge
	basicFlags := bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission

	err = adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  pingCharacteristicUUID,
				Flags: basicFlags,
				Value: []byte{0x0},
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					constants.EventChannel <- structures.Event{
						Data:      value,
						Operation: constants.UpdateWatchdog,
					}
				},
			},
			{
				UUID:  speedCharacteristicUUID,
				Flags: basicFlags,
				Value: constants.SpeedState,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					constants.EventChannel <- structures.Event{
						Data:      value,
						Operation: constants.SpeedOperation,
					}
					constants.SpeedState = value
				},
			},
			{
				UUID:  lightCharacteristicUUID,
				Flags: basicFlags,
				Value: constants.LightState,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					constants.EventChannel <- structures.Event{
						Data:      value,
						Operation: constants.LightOperation,
					}
					constants.LightState = value
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}

	// Retrieve & configure the default advertisement with the desired name
	adv := adapter.DefaultAdvertisement()

	if err = adv.Configure(bluetooth.AdvertisementOptions{LocalName: constants.ChipName}); err != nil {
		panic(err)
	}

	// Start advertising
	if err = adv.Start(); err != nil {
		panic(err)
	}
}
