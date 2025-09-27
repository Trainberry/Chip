# Trainberry - Chip

**This code was made for [Seeed Studio nRF52840 / Seeed XIAO BLE](https://wiki.seeedstudio.com/XIAO_BLE/). It may (will)
not work on other boards.** In theory, you should be able to make it work on any TinyGo supported board which have PWM
and BLE support; but you'll have to change constants in `internal/constants/constants.go`.

This project contains the code to deploy on train chip.

## v2 - Breaking change!

The first version of the project was using a Raspberry Pi Zero W with a WiFi stack. This version was great but had some
drawbacks (bad WiFi support, huge hardware fingerprint, 200ms latency for each action, and so on. See my talk at
Volcamp2025 for more details!).

This new v2.0 version uses Bluetooth BLE as communication support. It brings a very low latency (<1ms write, <20ms
read), low energy consumption (the superconductor is not even required anymore), and a much, much simpler code. The
counterpart is that the stack is less known as a classic WiFi + TCP/IP HTTP stack, so you might need to discover a bit
Bluetooth BLE.

## Flash

Use the TinyGo CLI to flash the card:

`tinygo flash -monitor -ldflags="-X 'main.chipName=YOUR_TRAIN_NAME' -X -main.reverse=IS_REVERSE" -target=xiao-ble`

Replace `YOUR_TRAIN_NAME` by the model of your train (eg. `BB92001`). This name **must be unique on your circuit**. If
you have multiple instance of the same model, add `_1`, `_2` and so on at the end; the frontend will hide it.

Replace `IS_REVERSE` by `true` or `false` to invert train direction (notably if you made a mistake while soldering... :D)

You're ready to go!

### Start with debug

To start debug mode (which attach serial UART to show logs on your computer), add `main.debug=true` flag:

`tinygo flash -monitor -ldflags="-X 'main.chipName=YOUR_TRAIN_NAME' -X 'main.debug=true'" -target=xiao-ble`

## BLE Services & Characteristics

The chip exposes 1 service containing 2 characteristics in a GATT server.

* The service have UUID `97980000-0000-0000-0000-000000000000`
* The characteristics are the following:
    * Ping Characteristic
        * Identified by UUID `00000000-0000-0000-0000-000000000000`
        * Always returns `0`
    * Speed Characteristic
        * Identified by UUID `10000000-0000-0000-0000-000000000000`
        * Expected payload is a signed integer (in a byte-array style) between `-100` and `100`
            * If value is below 0, it will make the train go backward (forward if above 0)
            * `0` will stop the train
    * Light Characteristic
        * Identified by UUID `20000000-0000-0000-0000-000000000000`
        * Expected payload is `on` or `off` in a byte-array style
            * `off` will turn off the light
            * `on` will turn on the light

# License

This project is licensed under [GNU Affero General Public License v3.0](https://choosealicense.com/licenses/agpl-3.0/).
