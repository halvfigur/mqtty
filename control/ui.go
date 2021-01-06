package control

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/network"
)

type (
	AppController interface {
		Stop()
	}

	MqttController interface {
		Connect(host string, port int, username, password string) error
		Subscribe(topic string, qos network.Qos) error
		Unsubscribe(topic string) error
	}

	ViewController interface {
		Register(pageLabel string, p tview.Primitive, visible bool)
		Display(pageLabel string)
		Hide(pageLabel string)
		Focus(p tview.Primitive)
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

func (u *MqttUI) Register(pageLabel string, p tview.Primitive, visible bool) {
	u.pages.AddPage(pageLabel, p, true, visible)
}

func (u *MqttUI) Display(pageLabel string) {
	u.pages.ShowPage(pageLabel)
}

func (u *MqttUI) Hide(pageLabel string) {
	u.pages.HidePage(pageLabel)
}

func (u *MqttUI) Focus(p tview.Primitive) {
	u.app.SetFocus(p)
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
	u.pages.SendToFront(startPageLabel)

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
