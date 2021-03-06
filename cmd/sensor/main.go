package main

import (
	"flag"

	"mattianatali.it/tm1637-mqtt/internal/sensor"
)

func main() {
	var mqttUsername = flag.String("u", "username", "Mqtt username")
	var mqttPassword = flag.String("p", "password", "Mqtt password")
	flag.Parse()

	c := sensor.Config{
		Topic:          "/home/milano/power-consumption",
		SensorPortPath: "/dev/ttyUSB0",
		MqttBroker:     "tcp://192.168.6.117:1883",
		MqttUsername:   *mqttUsername,
		MqttPassword:   *mqttPassword,
	}

	sensor.Start(c)
}
