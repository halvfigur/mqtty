package control

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

const (
	commanderLabel      = "commander"
	connectorLabel      = "connector"
	filtersLabel        = "filters"
	publishLabel        = "publish"
	publishHistoryLabel = "publishhistory"
	openFileLabel       = "openfile"
	errorLabel          = "error"
)

type (
	AppController interface {
		OnLaunchEditor() (string, error)
		QueueUpdate(func())
		QueueUpdateDraw(func())
		Stop()
	}

	MqttController interface {
		OnConnect(host string, port int, username, password string, onCompletion func(error))
		OnSubscribe(topic string, qos network.Qos, onCompletion func(error))
		OnUnsubscribe(topic string, onCompletion func(error))
		OnPublish(topic string, qos network.Qos, retained bool, message []byte, onCompletion func(error))
	}

	ViewController interface {
		OnDisplayConnector()
		OnDisplayCommander()
		OnDisplaySubscriber()
		OnDisplayPublisher()
		OnDisplayPublishHistory()
		OnDisplayOpenFile()
		OnDisplayError(err error)

		Register(pageLabel string, p tview.Primitive, visible bool)
		Hide(pageLabel string)
		Focus(p tview.Primitive)
		Cancel()
	}

	Control interface {
		AppController
		MqttController
		ViewController
	}

	MqttApp struct {
		c          *network.MqttClient
		app        *tview.Application
		pages      *tview.Pages
		errorModal *tview.Modal
		main       *CommanderController
	}
)

func NewMqttApp(c *network.MqttClient) *MqttApp {
	tview.Styles.TitleColor = tcell.ColorBlue
	app := tview.NewApplication()

	u := &MqttApp{
		c:     c,
		app:   app,
		pages: tview.NewPages(),
	}

	u.errorModal = view.NewErrorModal(u)

	u.Register(errorLabel, u.errorModal, false)

	u.main = NewCommanderController(u)

	return u
}

func (a *MqttApp) OnConnect(host string, port int, username, password string, onCompletion func(error)) {
	a.c.Connect(fmt.Sprintf("tcp://%s", host), port, username, password, onCompletion)
}

func (a *MqttApp) OnSubscribe(filter string, qos network.Qos, onCompletion func(error)) {
	a.c.Subscribe(filter, qos, onCompletion)
}

func (a *MqttApp) OnUnsubscribe(filter string, onCompletion func(error)) {
	a.c.Unsubscribe(filter, onCompletion)
}

func (a *MqttApp) OnPublish(topic string, qos network.Qos, retained bool, message []byte, onCompletion func(error)) {
	a.c.Publish(topic, qos, retained, message, onCompletion)
}

func (a *MqttApp) Register(pageLabel string, p tview.Primitive, visible bool) {
	a.pages.AddPage(pageLabel, p, true, visible)
}

func (a *MqttApp) display(pageLabel string) {
	a.pages.SendToFront(pageLabel)
	a.pages.ShowPage(pageLabel)
}

func (a *MqttApp) OnDisplayConnector() {
	a.display(connectorLabel)
}

func (a *MqttApp) OnDisplayCommander() {
	//a.display(commanderLabel)
	a.pages.SwitchToPage(commanderLabel)
}

func (a *MqttApp) OnDisplaySubscriber() {
	a.display(filtersLabel)
}

func (a *MqttApp) OnDisplayPublisher() {
	a.display(publishLabel)
}

func (a *MqttApp) OnDisplayPublishHistory() {
	a.display(publishHistoryLabel)
}

func (a *MqttApp) OnDisplayOpenFile() {
	a.display(openFileLabel)
}

func (a *MqttApp) OnDisplayError(err error) {
	a.errorModal.SetText(err.Error())
	a.display(errorLabel)
}

func (a *MqttApp) Hide(pageLabel string) {
	a.pages.HidePage(pageLabel)
}

func (a *MqttApp) Focus(p tview.Primitive) {
	a.app.SetFocus(p)
}

func (a *MqttApp) OnLaunchEditor() (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mqtty-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	a.app.Suspend(func() {
		cmd := exec.Command("/usr/bin/nvim", tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	})

	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func (a *MqttApp) OnStop() {
	a.Stop()
}

func (a *MqttApp) Stop() {
	a.app.Stop()
}

func (a *MqttApp) Cancel() {
	/*
		name, _ := a.pages.GetFrontPage()
		a.pages.HidePage(name)
	*/
	name, _ := a.pages.GetFrontPage()
	a.pages.HidePage(name)
	a.pages.SendToBack(name)
}

func (a *MqttApp) QueueUpdate(f func()) {
	a.app.QueueUpdate(f)
}

func (a *MqttApp) QueueUpdateDraw(f func()) {
	a.app.QueueUpdateDraw(f)
}

func (a *MqttApp) Start() {

	a.OnDisplayCommander()

	a.app.SetRoot(a.pages, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range a.c.Incoming() {
			a.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.Payload())
				a.main.AddDocument(m.Topic(), doc)
			})
		}
	}()

	if err := a.app.Run(); err != nil {
		log.Fatal(err)
	}
}
