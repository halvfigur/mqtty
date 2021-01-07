package network

import (
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Qos byte

const (
	QosAtMostOnce Qos = iota
	QosAtLeastOnce
	QosExatlyOnce
)

type (
	Message struct {
		topic    string
		qos      Qos
		id       uint16
		payload  []byte
		retained bool
	}

	MqttClient struct {
		c         mqtt.Client
		incomming chan *Message
		outgoing  chan *Message
		done      chan struct{}
	}
)

func (m *Message) Topic() string {
	return m.topic
}

func (m *Message) Qos() Qos {
	return m.qos
}

func (m *Message) Id() uint16 {
	return m.id
}

func (m *Message) Payload() []byte {
	return m.payload
}

func (m *Message) Retained() bool {
	return m.retained
}

func NewMqttClient() *MqttClient {
	c := &MqttClient{
		incomming: make(chan *Message, 32),
		outgoing:  make(chan *Message, 32),
		done:      make(chan struct{}),
	}

	go c.processPublished()

	return c
}

func (c *MqttClient) processPublished() {
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

func (c *MqttClient) Incomming() <-chan *Message {
	return c.incomming
}

func (c *MqttClient) Close() error {
	c.c.Disconnect(1000)
	close(c.done)
	close(c.incomming)
	close(c.outgoing)
	return nil
}

func (c *MqttClient) Connect(host string, port int, credentials ...string) error {
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
	if token := c.c.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MqttClient) Subscribe(topic string, qos Qos) error {
	if !c.c.IsConnected() {
		return errors.New("not connected")
	}

	t := c.c.Subscribe(topic, byte(qos), nil)
	t.Wait()
	return t.Error()
}

func (c *MqttClient) Unsubscribe(topic string) error {
	if !c.c.IsConnected() {
		return errors.New("not connected")
	}

	t := c.c.Unsubscribe(topic)
	t.Wait()
	return t.Error()
}

func (c *MqttClient) Publish(topic string, qos Qos, retained bool, message []byte) error {
	c.outgoing <- &Message{
		topic:    topic,
		qos:      qos,
		payload:  message,
		retained: retained,
	}

	return nil
}

func (c *MqttClient) onConnect(client mqtt.Client) {
}

func (c *MqttClient) onConnectionLost(client mqtt.Client, err error) {
}

func (c *MqttClient) onMessageArrived(client mqtt.Client, msg mqtt.Message) {
	c.incomming <- &Message{
		topic:    msg.Topic(),
		qos:      Qos(msg.Qos()),
		id:       msg.MessageID(),
		payload:  msg.Payload(),
		retained: msg.Retained(),
	}
}
