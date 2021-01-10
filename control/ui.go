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

	MqttUI struct {
		c     *network.MqttClient
		app   *tview.Application
		pages *tview.Pages
		main  *MainPageController
	}
)

func NewMqttUI(c *network.MqttClient) *MqttUI {
	tview.Styles.TitleColor = tcell.ColorBlue
	app := tview.NewApplication()

	u := &MqttUI{
		c:     c,
		app:   app,
		pages: tview.NewPages(),
	}

	u.main = NewMainPageController(u)

	return u
}

func (u *MqttUI) Connect(host string, port int, username, password string) error {
	if err := u.c.Connect(fmt.Sprintf("tcp://%s", host), port, username, password); err != nil {
		return err
	}
	u.Display(mainPageLabel)
	return nil
}

func (u *MqttUI) Subscribe(filter string, qos network.Qos) error {
	return u.c.Subscribe(filter, qos)
}

func (u *MqttUI) Unsubscribe(filter string) error {
	return u.c.Unsubscribe(filter)
}

func (u *MqttUI) Publish(topic string, qos network.Qos, retained bool, message []byte) error {
	return u.c.Publish(topic, qos, retained, message)
}

func (u *MqttUI) Register(pageLabel string, p tview.Primitive, visible bool) {
	u.pages.AddPage(pageLabel, p, true, visible)
}

func (u *MqttUI) Display(pageLabel string) {
	u.pages.ShowPage(pageLabel)
}

func (u *MqttUI) OnConnect() {
	u.Display(startPageLabel)
}

func (u *MqttUI) OnSubscribe() {
	u.Display(subscriptionFiltersViewLabel)
}

func (u *MqttUI) OnPublish() {
	// TODO
	//u.Display(subscriptionFiltersViewLabel)
}

func (u *MqttUI) Hide(pageLabel string) {
	u.pages.HidePage(pageLabel)
}

func (u *MqttUI) Focus(p tview.Primitive) {
	u.app.SetFocus(p)
}

func (u *MqttUI) OnLaunchEditor() (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mqtty-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	u.app.Suspend(func() {
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

func (u *MqttUI) OnStop() {
	u.Stop()
}

func (u *MqttUI) Stop() {
	u.app.Stop()
}

func (u *MqttUI) Cancel() {
	name, _ := u.pages.GetFrontPage()
	u.pages.HidePage(name)
}

func (u *MqttUI) Start() {
	//u.Display(mainPageLabel)
	//u.Display(startPageLabel)
	u.Display(publishPageLabel)

	u.app.SetRoot(u.pages, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range u.c.Incomming() {
			u.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.Payload())
				u.main.AddDocument(m.Topic(), doc)
			})
		}
	}()

	if err := u.app.Run(); err != nil {
		log.Fatal(err)
	}
}
