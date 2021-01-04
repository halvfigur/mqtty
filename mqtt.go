package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Qos byte

const (
	QosAtMostOnce Qos = iota
	QosAtLeastOnce
	QosExatlyOnce
)

type (
	mqttMessage struct {
		topic    string
		qos      Qos
		id       uint16
		payload  []byte
		retained bool
	}

	mqttClient struct {
		c         mqtt.Client
		incomming chan *mqttMessage
		outgoing  chan *mqttMessage
		done      chan struct{}
	}
)

func newMqttClient() *mqttClient {
	c := &mqttClient{
		incomming: make(chan *mqttMessage, 32),
		outgoing:  make(chan *mqttMessage, 32),
		done:      make(chan struct{}),
	}

	go c.processPublished()

	return c
}

func (c *mqttClient) processPublished() {
	for {
		select {
		case m := <-c.outgoing:
			if !c.c.IsConnected() {
				// Log error
				break
			}

			c.c.Publish(m.topic, byte(m.qos), m.retained, m.payload)
		case <-c.done:
			return
		}
	}
}

func (c *mqttClient) close() error {
	close(c.done)
	close(c.incomming)
	close(c.outgoing)
	return nil
}

func (c *mqttClient) connect(host string, port int, credentials ...string) error {
	var username string
	var password string

	if len(credentials) > 0 {
		if len(credentials) != 2 {
			panic("invalid credentials")
			username = credentials[0]
			password = credentials[1]
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", host, port))
	opts.SetClientID("mqtty")

	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	opts.SetDefaultPublishHandler(c.onMessageArrived)
	opts.OnConnect = c.onConnect
	opts.OnConnectionLost = c.onConnectionLost
	c.c = mqtt.NewClient(opts)
	if token := c.c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *mqttClient) subscribe(topic string, qos Qos) error {
	t := c.c.Subscribe(topic, byte(qos), nil)
	t.Wait()
	return t.Error()
}

func (c *mqttClient) unsubscribe(topic string) error {
	t := c.c.Unsubscribe(topic)
	t.Wait()
	return t.Error()
}

func (c *mqttClient) publish(topic string, qos Qos, retained bool, message []byte) error {
	c.outgoing <- &mqttMessage{
		topic:    topic,
		qos:      qos,
		payload:  message,
		retained: retained,
	}

	return nil
}

func (c *mqttClient) onConnect(client mqtt.Client) {
}

func (c *mqttClient) onConnectionLost(client mqtt.Client, err error) {
}

func (c *mqttClient) onMessageArrived(client mqtt.Client, msg mqtt.Message) {
	c.incomming <- &mqttMessage{
		topic:    msg.Topic(),
		qos:      Qos(msg.Qos()),
		id:       msg.MessageID(),
		payload:  msg.Payload(),
		retained: msg.Retained(),
	}
}
