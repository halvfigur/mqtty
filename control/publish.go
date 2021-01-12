package control

import (
	"io/ioutil"
	"os"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

type Publish struct {
	ctrl        Control
	view        *view.Publish
	fileView    *view.OpenFile
	historyCtrl *PublishHistory
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

	c.historyCtrl = NewPublishHistory(ctrl)

	c.ctrl.Register(publishLabel, c.view, false)
	c.ctrl.Register(openFileLabel, view.Center(c.fileView, 1, 1), false)

	return c
}

func (c *Publish) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *Publish) Register(label string, p tview.Primitive, visible bool) {
	c.ctrl.Register(label, p, visible)
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
	c.ctrl.OnDisplayOpenFile()
}

func (c *Publish) OnFileSelected(filename string) {
	c.readAndUpdateView(filename)
	c.ctrl.Cancel()
	c.ctrl.OnDisplayPublisher()
}

func (c *Publish) OnOpenHistory() {
	c.ctrl.OnDisplayPublishHistory()
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

func (c *Publish) OnPublish(topic string, qos network.Qos, retained bool, message []byte) {
	c.ctrl.OnPublish(topic, qos, retained, message, func(err error) {
		c.ctrl.QueueUpdate(func() {
			if err != nil {
				//TODO handle error
				return
			}
			c.historyCtrl.AddDocument(topic, data.NewDocumentBytes(message))
		})
	})
}

func (c *Publish) QueueUpdate(f func()) {
	c.ctrl.QueueUpdate(f)
}

func (c *Publish) Cancel() {
	c.ctrl.Hide(publishLabel)
}
