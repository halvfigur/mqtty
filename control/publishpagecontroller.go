package control

import (
	"io/ioutil"
	"os"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

const publishPageLabel = "publishpage"

type PublishPageController struct {
	ctrl Control
	view *view.PublishPage
}

func NewPublishPageController(ctrl Control) *PublishPageController {
	c := &PublishPageController{
		ctrl: ctrl,
	}

	c.view = view.NewPublishPage(c)
	c.ctrl.Register(publishPageLabel, c.view, false)

	return c
}

func (c *PublishPageController) GetView(p tview.Primitive) *view.PublishPage {
	return c.view
}

func (c *PublishPageController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *PublishPageController) OnLaunchEditor() {
	filename, err := c.ctrl.OnLaunchEditor()
	if err != nil {
		// Handle error
		return
	}
	defer os.Remove(filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Handle error
	}

	c.view.SetData(data)
}

func (c *PublishPageController) Publish(topic string, qos network.Qos, retained bool, message []byte) error {
	return c.ctrl.Publish(topic, qos, retained, message)
}

func (c *PublishPageController) Cancel() {
	c.ctrl.Hide(publishPageLabel)
}
