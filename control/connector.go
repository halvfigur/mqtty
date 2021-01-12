package control

import "github.com/halvfigur/mqtty/view"

type (
	ConnectFunc func(host string, port int, username, password string)
	Connector   struct {
		ctrl Control
		view *view.Connector
	}
)

func NewConnector(ctrl Control) *Connector {
	c := &Connector{
		ctrl: ctrl,
	}

	c.view = view.NewConnector(c)

	ctrl.Register(connectorLabel, c.view, false)
	return c
}

func (c *Connector) OnConnect(host string, port int, username, password string, onCompletion func(error)) {
	c.ctrl.OnConnect(host, port, username, password, onCompletion)
}

func (c *Connector) OnConnected() {
	c.ctrl.Cancel()
}

func (c *Connector) QueueUpdate(f func()) {
	c.ctrl.QueueUpdate(f)
}

func (c *Connector) Stop() {
	c.ctrl.Cancel()
}
