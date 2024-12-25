# Train Chip

**This code is made for Raspberry Pi Pico W. It may not work with other chips.**

This projet contains the code that will go on the Pico located in the train.

## Quick start

### Configuration

First of all, you must set the configuration for your board : Train IP, Name, model, WiFi settings, and so on. The easiest way is to copy the `internal/configuration/configuration.go.template` to `internal/configuration/configuration.go`, and then edit values with your own.

As there is a lot of problems with ARP requests on Pico W, you must also provide by yourself the MAC address of the central server, in hex format. For example, if the MAC adress of the central server is `ab:cd:ef:01:02:03`, the `ServerMAC` should be `[6]byte{0xAB, 0xCD, 0xEF, 0x01, 0x02, 0x03}`. Easy, isn't it?

### Flash

To flash, you have to specify your target and a specific `stack-size` :

```sh
tinygo flash -target=pico -stack-size=16kb -monitor  ./cmd
```

Tested with TinyGo v0.33.0.

## General usage

On start, the chip will connect against the WiFi network and request his own IP. For performance reason, there is no "fancy" network stuff (no DHCP, no DNS, ...).

When connected, the chip will perform a `POST /register` request on the central server, with a JSON-formatted payload which looks like this:

```json
{
    "name": "BB27000",
    "model": "BB27000",
    "ip": "192.168.1.150"
}
```

The chip is now ready to accept new connections to verify state or change status.

## API

The chip exposes an API on port `80/tcp`. For performance reason, there is no authentication or anything else like this. We highly recommend to use a dedicated WiFi network for your train infrastructure (a very old 100Mbps WiFi router will be quite enough).

The whole reference is available in the `api` folder.