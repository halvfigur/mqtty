package main

func main() {
	c := newMqttClient()
	c.connect("tcp://test.mosquitto.org", 1883)

	c.subscribe("hamweather/#", QosAtLeastOnce)

	ui(c.incomming)
	/*
		for m := range c.incomming {
			log.Println(m)
		}
	*/

}
