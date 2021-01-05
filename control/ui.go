package control

import (
	"fmt"
	"log"

	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

type (
	Control interface {
		Connect(host string, port int, username, password string)

		OnSubscribe()
		Subscribe(topic string, qos network.Qos)

		Renderers() []model.Renderer
		OnRenderer()
		Renderer(renderer model.Renderer)
		Focus(p tview.Primitive)

		OnStop()
		Stop()

		Cancel()
	}

	controllers struct {
		start     *StartPageController
		main      *MainPageController
		subscribe *SubscribePageController
		renderer  *RendererPageController
	}

	MqttUI struct {
		c         *network.MqttClient
		app       *tview.Application
		pages     *tview.Pages
		ctrl      controllers
		renderers []model.Renderer
	}
)

func NewMqttUI(c *network.MqttClient) *MqttUI {
	app := tview.NewApplication()

	u := &MqttUI{
		c:   c,
		app: app,
		renderers: []model.Renderer{
			model.NewRawRenderer(),
			model.NewHexRenderer(),
			model.NewJsonRenderer(),
		},
	}

	startCtrl := NewStartPageController(u)
	startPage := view.NewStartPage(startCtrl)

	mainCtrl := NewMainPageController(u)
	mainPage := view.NewMainPage(mainCtrl)
	mainCtrl.SetView(mainPage)

	subscribeCtrl := NewSubscribePageController(u)
	subscribePage := view.NewSubscribePage(u)

	rendererCtrl := NewRendererPageController(u)
	rendererPage := view.NewRendererPage(u)

	u.pages = tview.NewPages().
		AddPage(mainPageLabel, mainPage, false, true).
		AddPage(subscribePageLabel, subscribePage, true, true).
		AddPage(rendererPageLabel, rendererPage, true, true).
		AddAndSwitchToPage(startPageLabel, startPage, true)

	u.ctrl = controllers{
		start:     startCtrl,
		main:      mainCtrl,
		subscribe: subscribeCtrl,
		renderer:  rendererCtrl,
	}

	return u
}

func (u *MqttUI) Connect(host string, port int, username, password string) {
	u.c.Connect(fmt.Sprintf("tcp://%s", host), port, username, password)
	u.c.Subscribe("hamweather/#", network.QosAtLeastOnce)
	u.pages.SwitchToPage(mainPageLabel)
}

func (u *MqttUI) OnSubscribe() {
	u.pages.ShowPage(subscribePageLabel)
}

func (u *MqttUI) Subscribe(topic string, qos network.Qos) {
	u.pages.HidePage(subscribePageLabel)
	u.c.Subscribe(topic, qos)
}

func (u *MqttUI) Renderers() []model.Renderer {
	return u.renderers
}

func (u *MqttUI) OnRenderer() {
	u.pages.ShowPage(rendererPageLabel)
}

func (u *MqttUI) Renderer(renderer model.Renderer) {
	u.pages.HidePage(rendererPageLabel)
	u.ctrl.main.SetRenderer(renderer)
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
	u.app.SetRoot(u.pages, true)

	go func() {
		/* This goroutine will exit when the incomming channel is closed */
		for m := range u.c.Incomming() {
			u.app.QueueUpdateDraw(func() {
				doc := data.NewDocumentBytes(m.Payload())
				u.ctrl.main.AddDocument(m.Topic(), doc)
			})
		}
	}()

	if err := u.app.Run(); err != nil {
		log.Fatal(err)
	}
}
