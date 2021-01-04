package main

import "github.com/rivo/tview"

type (
	connectFunc         func(host string, port int, username, password string)
	StartPageController struct {
		app     *tview.Application
		connect connectFunc
	}
)

func NewStartPageController(a *tview.Application, connect connectFunc) *StartPageController {
	return &StartPageController{
		app:     a,
		connect: connect,
	}
}

func (c *StartPageController) OnConnect(host string, port int, username, password string) {
	c.connect(host, port, username, password)
}

func (c *StartPageController) OnExit() {
	c.app.Stop()
}
