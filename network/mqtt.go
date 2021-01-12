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

	incomingChanCap = 1024
	requestsChanCap = 32

	connectTimeout     = 5 * time.Second
	publishTimeout     = 5 * time.Second
	subscribeTimeout   = 5 * time.Second
	unsubscribeTimeout = 5 * time.Second
)

var errConnected = errors.New("already connected")
var errNotConnected = errors.New("not connected")

type (
	Message struct {
		topic    string
		qos      Qos
		id       uint16
		payload  []byte
		retained bool
	}

	connectRequest struct {
		host         string
		port         int
		username     string
		password     string
		onCompletion func(error)
	}

	publishRequest struct {
		topic        string
		qos          Qos
		id           uint16
		payload      []byte
		retained     bool
		onCompletion func(error)
	}

	subscribeRequest struct {
		topic        string
		qos          Qos
		onCompletion func(error)
	}

	unsubscribeRequest struct {
		topic        string
		onCompletion func(error)
	}

	MqttClient struct {
		c        mqtt.Client
		incoming chan *Message
		requests chan interface{}
		done     chan struct{}
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
		incoming: make(chan *Message, incomingChanCap),
		requests: make(chan interface{}, requestsChanCap),
		done:     make(chan struct{}),
	}

	go c.processRequests()

	return c
}

func (c *MqttClient) processRequests() {
	for {
		select {
		case req := <-c.requests:

			switch r := req.(type) {
			case *connectRequest:
				c.handleConnectRequest(r)
			case *publishRequest:
				c.handlePublishRequest(r)
			case *subscribeRequest:
				c.handleSubscribeRequest(r)
			case *unsubscribeRequest:
				c.handleUnsubscribeRequest(r)
			}

		case <-c.done:
			return
		}
	}
}

func (c *MqttClient) handleConnectRequest(r *connectRequest) {
	if c.c != nil && c.c.IsConnected() {
		if r.onCompletion != nil {
			r.onCompletion(errConnected)
		}

		return
	}

	var username string
	var password string

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", r.host, r.port))
	opts.SetClientID("mqtty")

	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	opts.SetDefaultPublishHandler(c.onMessageArrived)
	opts.OnConnect = c.onConnect
	opts.OnConnectionLost = c.onConnectionLost
	c.c = mqtt.NewClient(opts)

	t := c.c.Connect()

	t.WaitTimeout(connectTimeout)

	if r.onCompletion != nil {
		r.onCompletion(t.Error())
	}
}

func (c *MqttClient) handlePublishRequest(r *publishRequest) {
	if !c.c.IsConnected() {
		if r.onCompletion != nil {
			r.onCompletion(errNotConnected)
		}

		return
	}

	t := c.c.Publish(r.topic, byte(r.qos), r.retained, r.payload)
	t.WaitTimeout(publishTimeout)

	if r.onCompletion != nil {
		r.onCompletion(t.Error())
	}
}

func (c *MqttClient) handleSubscribeRequest(r *subscribeRequest) {
	if !c.c.IsConnected() {
		if r.onCompletion != nil {
			r.onCompletion(errNotConnected)
		}
	}

	t := c.c.Subscribe(r.topic, byte(r.qos), nil)
	t.WaitTimeout(subscribeTimeout)

	if r.onCompletion != nil {
		r.onCompletion(t.Error())
	}
}

func (c *MqttClient) handleUnsubscribeRequest(r *unsubscribeRequest) {
	if !c.c.IsConnected() {
		if r.onCompletion != nil {
			r.onCompletion(errNotConnected)
		}
	}

	t := c.c.Unsubscribe(r.topic)
	t.WaitTimeout(unsubscribeTimeout)

	if r.onCompletion != nil {
		r.onCompletion(t.Error())
	}
}

func (c *MqttClient) Incoming() <-chan *Message {
	return c.incoming
}

func (c *MqttClient) Close() error {
	c.c.Disconnect(1000)
	close(c.done)
	close(c.incoming)
	close(c.requests)
	return nil
}

func (c *MqttClient) Connect(host string, port int, username, password string, onCompletion func(error)) {
	c.requests <- &connectRequest{
		host:         host,
		port:         port,
		username:     username,
		password:     password,
		onCompletion: onCompletion,
	}
}

func (c *MqttClient) Subscribe(topic string, qos Qos, onCompletion func(error)) {
	c.requests <- &subscribeRequest{
		topic:        topic,
		qos:          qos,
		onCompletion: onCompletion,
	}
}

func (c *MqttClient) Unsubscribe(topic string, onCompletion func(error)) {
	c.requests <- &unsubscribeRequest{
		topic:        topic,
		onCompletion: onCompletion,
	}
}

func (c *MqttClient) Publish(topic string, qos Qos, retained bool, message []byte, onCompletion func(error)) {
	c.requests <- &publishRequest{
		topic:        topic,
		qos:          qos,
		payload:      message,
		retained:     retained,
		onCompletion: onCompletion,
	}
}

func (c *MqttClient) onConnect(client mqtt.Client) {
}

func (c *MqttClient) onConnectionLost(client mqtt.Client, err error) {
}

func (c *MqttClient) onMessageArrived(client mqtt.Client, msg mqtt.Message) {
	c.incoming <- &Message{
		topic:    msg.Topic(),
		qos:      Qos(msg.Qos()),
		id:       msg.MessageID(),
		payload:  msg.Payload(),
		retained: msg.Retained(),
	}
}
