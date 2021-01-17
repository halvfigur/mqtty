package control

import (
	"github.com/halvfigur/mqtty/view"
	"github.com/rivo/tview"
)

type (
	Connector struct {
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

func (c *Connector) OnConnect(host string, port int, username, password string) {
	c.ctrl.OnConnect(host, port, username, password)
}

func (c *Connector) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *Connector) Cancel() {
	c.ctrl.Cancel()
}

func (c *Connector) OnError(err error) {
	c.ctrl.OnDisplayError(err)
}
