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
)

type (
	AppController interface {
		OnConnect()
		OnSubscribe()
		OnPublish()
		OnLaunchEditor() (string, error)
		Stop()
	}

	MqttController interface {
		Connect(host string, port int, username, password string) error
		Subscribe(topic string, qos network.Qos) error
		Unsubscribe(topic string) error
		Publish(topic string, qos network.Qos, retained bool, message []byte) error
	}

	ViewController interface {
		Register(pageLabel string, p tview.Primitive, visible bool)
		Display(pageLabel string)
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
		c     *network.MqttClient
		app   *tview.Application
		pages *tview.Pages
		main  *CommanderController
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

	u.main = NewMainPageController(u)

	return u
}

func (a *MqttApp) Connect(host string, port int, username, password string) error {
	if err := a.c.Connect(fmt.Sprintf("tcp://%s", host), port, username, password); err != nil {
		return err
	}
	a.Display(commanderLabel)
	return nil
}

func (a *MqttApp) Subscribe(filter string, qos network.Qos) error {
	return a.c.Subscribe(filter, qos)
}

func (a *MqttApp) Unsubscribe(filter string) error {
	return a.c.Unsubscribe(filter)
}

func (a *MqttApp) Publish(topic string, qos network.Qos, retained bool, message []byte) error {
	return a.c.Publish(topic, qos, retained, message)
}

func (a *MqttApp) Register(pageLabel string, p tview.Primitive, visible bool) {
	a.pages.AddPage(pageLabel, p, true, visible)
}

func (a *MqttApp) Display(pageLabel string) {
	a.pages.SendToFront(pageLabel)
	a.pages.ShowPage(pageLabel)
}

func (a *MqttApp) OnConnect() {
	a.Display(connectorLabel)
}

func (a *MqttApp) OnSubscribe() {
	a.Display(filtersLabel)
}

func (a *MqttApp) OnPublish() {
	a.Display(publishPageLabel)
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

func (a *MqttApp) Start() {
	a.Display(connectorLabel)

	a.app.SetRoot(a.pages, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range a.c.Incomming() {
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
