package main

func main() {
	c := newMqttClient()
	c.connect("tcp://test.mosquitto.org", 1883)
	defer c.close()

	c.subscribe("hamweather/#", QosAtLeastOnce)

	ui := newMqttUI(c.incomming)
	ui.run()
}
