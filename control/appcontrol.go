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
	waitLabel           = "wait"
)

type (
	Config struct {
		Server string
		Port   int
		Topics []string
	}

	AppController interface {
		OnLaunchEditor() (string, error)
		QueueUpdate(func())
		QueueUpdateDraw(func())
		Stop()
	}

	MqttController interface {
		OnConnect(server string, port int, username, password string)
		OnSubscribe(topic string, qos network.Qos)
		OnUnsubscribe(topic string)
		OnPublish(topic string, qos network.Qos, retained bool, message []byte)
	}

	ViewController interface {
		OnDisplayConnector()
		OnDisplayCommander()
		OnDisplaySubscriber()
		OnDisplayPublisher()
		OnDisplayPublishHistory()
		OnDisplayOpenFile()
		OnDisplayError(err error)
		OnDisplayWait(msg string)

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
		conf       Config
		app        *tview.Application
		pages      *tview.Pages
		errorModal *tview.Modal
		waitModal  *tview.Modal
		main       *CommanderController
	}
)

func NewMqttApp(c *network.MqttClient, conf Config) *MqttApp {
	tview.Styles.TitleColor = tcell.ColorBlue
	app := tview.NewApplication()

	u := &MqttApp{
		c:     c,
		conf:  conf,
		app:   app,
		pages: tview.NewPages(),
	}

	u.errorModal = view.NewErrorModal(u)
	u.waitModal = view.NewWaitModal()

	u.Register(errorLabel, u.errorModal, false)
	u.Register(waitLabel, u.waitModal, false)

	u.main = NewCommanderController(u)

	u.c.SetConnectionStatusFunc(func(s network.ConnectionStatus) {
		u.app.QueueUpdateDraw(func() {
			u.main.SetConnectionStatus(s)
		})
	})

	return u
}

func (a *MqttApp) OnConnect(server string, port int, username, password string) {
	a.OnDisplayWait(fmt.Sprint("Connecting to ", a.conf.Server))
	a.c.Connect(server, port, username, password, func(err error) {
		a.QueueUpdateDraw(func() {
			a.OnDisplayCommander()

			if err != nil {
				a.main.SetConnectionStatus(network.StatusDisconnected)
				a.OnDisplayError(err)
			}

			a.main.SetConnectionStatus(network.StatusConnected)
		})
	})
}

func (a *MqttApp) OnSubscribe(filter string, qos network.Qos) {
	a.c.Subscribe(filter, qos, func(err error) {
		a.OnDisplayWait(fmt.Sprint("Subscribing to ", filter))
		a.app.QueueUpdateDraw(func() {
			// Hide wait modal
			a.Cancel()
			if err != nil {
				a.OnDisplayError(err)
				return
			}

			a.main.AddFilter(filter, qos)
		})
	})
}

func (a *MqttApp) OnUnsubscribe(filter string) {
	a.c.Unsubscribe(filter, func(err error) {
		a.app.QueueUpdateDraw(func() {
			if err != nil {
				a.OnDisplayError(err)
				return
			}

			a.main.RemoveFilter(filter)
		})
	})
}

func (a *MqttApp) OnPublish(topic string, qos network.Qos, retained bool, message []byte) {
	a.c.Publish(topic, qos, retained, message, func(err error) {
		a.QueueUpdateDraw(func() {
			if err != nil {
				a.OnDisplayError(err)
				return
			}
			a.main.AddPublishedDocument(topic, data.NewDocumentBytes(message))
		})
	})
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
	a.display(commanderLabel)
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

func (a *MqttApp) OnDisplayWait(msg string) {
	a.waitModal.SetText(msg)
	a.display(waitLabel)
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
	name, _ := a.pages.GetFrontPage()
	a.pages.HidePage(name)
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

	if a.conf.Server != "" {
		a.c.Connect(a.conf.Server, a.conf.Port, "", "", func(err error) {
			if err != nil {
				a.OnDisplayError(err)
				return
			}

			for _, f := range a.conf.Topics {
				a.c.Subscribe(f, network.QosAtMostOnce, func(err error) {
					a.QueueUpdateDraw(func() {
						if err != nil {
							a.OnDisplayError(err)
							return
						}

						a.main.AddFilter(f, network.QosAtMostOnce)
					})
				})
			}
		})
	}

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range a.c.Incoming() {
			a.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.Payload())
				a.main.AddDocument(m.Topic(), doc)
			})
		}
	}()

	err := a.app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
