package main

import (
	"log"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/view"
)

type (
	UI interface {
		OnIncomming(m *mqttMessage)
	}

	controllers struct {
		main *MainPageController
	}

	mqttUI struct {
		incomming <-chan *mqttMessage
		app       *tview.Application
		ctrl      controllers
	}
)

func newMqttUI(incomming <-chan *mqttMessage) *mqttUI {
	app := tview.NewApplication()
	mainCtrl := NewMainPageController(app)
	mainPage := view.NewMainPage(mainCtrl)
	mainCtrl.SetView(mainPage)

	u := &mqttUI{
		incomming: incomming,
		app:       app,
		ctrl: controllers{
			main: mainCtrl,
		},
	}

	go u.run()

	return u
}

func (u *mqttUI) run() {
	u.app.SetRoot(u.ctrl.main.view, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range u.incomming {
			u.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.payload)
				//u.ctrl.main.SetDocument(doc)
				u.ctrl.main.AddDocument(m.topic, doc)
			})
		}
	}()

	if err := u.app.Run(); err != nil {
		log.Fatal(err)
	}
}
