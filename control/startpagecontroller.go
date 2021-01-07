package control

import "github.com/halvfigur/mqtty/view"

const startPageLabel = "startpage"

type (
	ConnectFunc         func(host string, port int, username, password string)
	StartPageController struct {
		ctrl Control
		view *view.StartPage
	}
)

func NewStartPageController(ctrl Control) *StartPageController {
	c := &StartPageController{
		ctrl: ctrl,
	}

	c.view = view.NewStartPage(c)

	ctrl.Register(startPageLabel, c.view, false)
	return c
}

func (c *StartPageController) GetView() *view.StartPage {
	return c.view
}

func (c *StartPageController) OnConnect(host string, port int, username, password string) error {
	if err := c.ctrl.Connect(host, port, username, password); err != nil {
		// TODO display error message
		return err
	}

	c.ctrl.Hide(startPageLabel)
	return nil
}

func (c *StartPageController) Stop() {
	c.ctrl.Cancel()
}
