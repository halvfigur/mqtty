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

type Publish struct {
	ctrl     Control
	view     *view.Publish
	fileView *view.OpenFile
}

func NewPublish(ctrl Control) *Publish {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	c := &Publish{
		ctrl: ctrl,
	}

	c.view = view.NewPublish(c)

	c.fileView = view.NewOpenFile(cwd).
		SetOnFileSelected(c.OnFileSelected).
		SetOnError(c.OnError)
	c.ctrl.Register(publishPageLabel, c.view, false)
	c.ctrl.Register(openFileViewLabel, view.Center(c.fileView, 1, 1), false)

	return c
}

func (c *Publish) GetView(p tview.Primitive) *view.Publish {
	return c.view
}

func (c *Publish) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *Publish) OnLaunchEditor() {
	filename, err := c.ctrl.OnLaunchEditor()
	if err != nil {
		// Handle error
		return
	}
	defer os.Remove(filename)

	c.readAndUpdateView(filename)
}

func (c *Publish) OnOpenFile() {
	c.ctrl.Display(openFileViewLabel)
}

func (c *Publish) OnFileSelected(filename string) {
	c.readAndUpdateView(filename)
	c.ctrl.Hide(openFileViewLabel)
	c.ctrl.Display(publishPageLabel)
}

func (c *Publish) OnError(err error) {
}

func (c *Publish) readAndUpdateView(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Handle error
	}

	c.view.SetData(data)
}

func (c *Publish) Publish(topic string, qos network.Qos, retained bool, message []byte) error {
	return c.ctrl.Publish(topic, qos, retained, message)
}

func (c *Publish) Cancel() {
	c.ctrl.Hide(publishPageLabel)
}