package main

import (
	"github.com/halvfigur/mqtty/control"
	"github.com/halvfigur/mqtty/network"
)

func main() {
	c := network.NewMqttClient()
	control.NewMqttApp(c).Start()
}
