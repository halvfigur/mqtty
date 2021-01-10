package control

import "github.com/halvfigur/mqtty/view"

const connectorLabel = "connector"

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

func (c *Connector) GetView() *view.Connector {
	return c.view
}

func (c *Connector) OnConnect(host string, port int, username, password string) error {
	if err := c.ctrl.Connect(host, port, username, password); err != nil {
		// TODO display error message
		return err
	}

	c.ctrl.Hide(connectorLabel)
	return nil
}

func (c *Connector) Stop() {
	c.ctrl.Cancel()
}
