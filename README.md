# han-mqtt
This is a very simple application to forward readings from electricity meters to MQTT. The app connects to a serial interface, assumes readouts according to IEC 62056-21 and sends the parsed message to an MQTT broker. The idea is to consume these messages in home automation setups such as Home Assistant or Node-RED.

The app was written for the HAN port in mind but can just as well be used with other interfaces outputting data according to IEC 62056-21 such as IR interfaces.

---
**NOTE**

This app does not request data from the source but assumes that data will be sent on a regular interval. For HAN / P1, data is sent every 5-10 seconds if PIN 2 (Data request) is high

---

## Running with npm
`npm install -g han-mqtt`

`HAN_MQTT_CONFIG=./config.yaml han-mqtt`

## Running with docker
`sudo docker run -e HAN_MQTT_CONFIG=/config/config.yaml -v /home/pa/han-mqtt-config/:/config --privileged pakerfeldt/han-mqtt:latest`
This assumes a directory called `/home/pa/han-mqtt-config/` on the host machine which contains your `config.yaml`. Obviously change the host path to wherever you choose to store the config file. `--privileged` is needed for the container to gain access of the physical serial interface.

## Running with docker-compose
Perhaps the most convenient of choices. Ensure you have docker and docker-compose installed. Then create a `docker-compose.yaml` containing:
```yaml
version: '3.2'
services:
  hanmqtt:
    image: pakerfeldt/han-mqtt:latest
    environment:
      - NODE_ENV=production
      - HAN_MQTT_CONFIG=/config/config.yaml
    volumes:
      - /home/pa/han-mqtt-config:/config
    privileged: true
    restart: unless-stopped
```
Now you can build and run the image. 

`sudo docker-compose build hanmqtt`

`sudo docker-compose up -d`

`privileged: true` is needed for the container to gain access of the physical serial interface.

## HAN / P1 port
HAN or P1 is the name of the RJ12 port found on Swedish electricity meters and is meant to be used for local applications such as consumption metrics or load balancers. This interface uses 115200 baud rate with 8 databits, 1 stopbit and parity none. However, signal is inverted thus cannot be connected directly to a serial device unless first taken care of.

This readme won't go into details on how to accomplish this. There's already guides on how to do that. In short, you'll need a transistor and two resistors. Here's a complete BOM from electrokit if you choose to go that route:

| Item                                |
|-------------------------------------|
| [BC547C TO-92 NPN 45V 100mA](https://www.electrokit.com/produkt/bc547c/)
| [Motstånd kolfilm 0.25W 10kohm (10k)](https://www.electrokit.com/produkt/motstand-kolfilm-0-25w-10kohm-10k/)
| [Motstånd kolfilm 0.25W 4.7kohm (4k7)](https://www.electrokit.com/produkt/motstand-kolfilm-0-25w-4-7kohm-4k7/)
| [Experimentkort för FB01](https://www.electrokit.com/produkt/experimentkort-for-fb01/)
| [Apparatlåda FB01 grå 60x65x25mm](https://www.electrokit.com/produkt/apparatlada-fb01-gra-60x65x25mm/)
| [USB-serielladapter PL2303](https://www.electrokit.com/produkt/usb-serielladapter-pl2303/)

You might very well find smaller circuit boards but this was the one I went for.
Here's the schematics which can also be found, together with more description in Swedish, on https://www.akehedman.se/wordpress/?cat=39

<img src="https://user-images.githubusercontent.com/195860/155875137-820ef95f-fce9-412d-8723-b0e575b98b13.png" width="640">
Additional information on the HAN port on Swedish electricity meters can be found here: https://hanporten.se/svenska/porten/

