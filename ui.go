package main

import (
	"fmt"
	"log"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/view"
)

const (
	startPageLabel = "startpage"
	mainPageLabel  = "mainpage"
)

type (
	UI interface {
		OnIncomming(m *mqttMessage)
	}

	controllers struct {
		start *StartPageController
		main  *MainPageController
	}

	mqttUI struct {
		c     *mqttClient
		app   *tview.Application
		pages *tview.Pages
		ctrl  controllers
	}
)

func newMqttUI(c *mqttClient) *mqttUI {
	app := tview.NewApplication()

	u := &mqttUI{
		c:   c,
		app: app,
	}
	startCtrl := NewStartPageController(app, u.connect)
	startPage := view.NewStartPage(startCtrl)

	mainCtrl := NewMainPageController(app)
	mainPage := view.NewMainPage(mainCtrl)
	mainCtrl.SetView(mainPage)

	u.pages = tview.NewPages().
		AddPage(mainPageLabel, mainPage, false, true).
		AddAndSwitchToPage(startPageLabel, startPage, true)

	u.ctrl = controllers{
		start: startCtrl,
		main:  mainCtrl,
	}

	go u.run()

	return u
}

func (u *mqttUI) connect(host string, port int, username, password string) {
	u.c.connect(fmt.Sprintf("tcp://%s", host), port, username, password)
	u.c.subscribe("hamweather/#", QosAtLeastOnce)
	u.pages.SwitchToPage(mainPageLabel)
}

func (u *mqttUI) run() {
	u.app.SetRoot(u.pages, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range u.c.incomming {
			u.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.payload)
				u.ctrl.main.AddDocument(m.topic, doc)
			})
		}
	}()

	if err := u.app.Run(); err != nil {
		log.Fatal(err)
	}
}
