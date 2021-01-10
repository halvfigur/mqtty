package control

import (
	"io/ioutil"
	"os"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

const (
	publishPageLabel  = "publishpage"
	openFileViewLabel = "openfileview"
)

type PublishPageController struct {
	ctrl     Control
	view     *view.PublishPage
	fileView *view.OpenFileView
}

func NewPublishPageController(ctrl Control) *PublishPageController {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	c := &PublishPageController{
		ctrl: ctrl,
	}

	c.view = view.NewPublishPage(c)

	c.fileView = view.NewOpenFileView(cwd).
		SetOnFileSelected(c.OnFileSelected).
		SetOnError(c.OnError)
	c.ctrl.Register(publishPageLabel, c.view, false)
	c.ctrl.Register(openFileViewLabel, view.Center(c.fileView, 1, 1), false)

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

	c.readAndUpdateView(filename)
}

func (c *PublishPageController) OnOpenFile() {
	c.ctrl.Display(openFileViewLabel)
}

func (c *PublishPageController) OnFileSelected(filename string) {
	c.readAndUpdateView(filename)
	c.ctrl.Hide(openFileViewLabel)
	c.ctrl.Display(publishPageLabel)
}

func (c *PublishPageController) OnError(err error) {
}

func (c *PublishPageController) readAndUpdateView(filename string) {
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
