package main

import (
	"github.com/halvfigur/mqtty/control"
	"github.com/halvfigur/mqtty/network"
)

func main() {
	c := network.NewMqttClient()
	//c.connect("tcp://test.mosquitto.org", 1883)
	//defer c.close()

	//c.subscribe("hamweather/#", QosAtLeastOnce)
	//c.subscribe("#", QosAtLeastOnce)

	ui := control.NewMqttUI(c)
	ui.Start()
}
