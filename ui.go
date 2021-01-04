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
				u.ctrl.main.SetDocument(doc)
			})
		}
	}()

	/*
		u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			//debugView.SetText(event.Name())

			switch event.Key() {
			case tcell.KeyTab:
				u.app.SetFocus(fc.Next())
			case tcell.KeyBacktab:
				u.app.SetFocus(fc.Prev())
			default:
				return event
			}

			return nil
		})
	*/

	if err := u.app.Run(); err != nil {
		log.Fatal(err)
	}
}

/*
func (u *mqttUI) setup() {
	// Document data
	d := data.NewDocumentEmpty()

	debugView := tview.NewTextView()

	docCtrl := NewDocumentController()
	docCtrl.SetDocument(d)

	topicList := tview.NewList()
	topicList.SetBorder(true).SetTitle("Topics")
	topicList.ShowSecondaryText(false)
	topicList.AddItem("iotea/ingestion/events", "", 0, nil)
	topicList.AddItem("iotea/discovery", "", 0, nil)

	renderers := []view.Renderer{
		new(view.RawRenderer),
		view.NewHexRenderer(),
	}

	renderersList := tview.NewList()
	renderersList.SetBorder(true).SetTitle("Renderers")
	renderersList.ShowSecondaryText(false)
	for _, r := range renderers {
		renderersList.AddItem(r.Name(), "", 0,
			func(r view.Renderer) func() {
				return func() {
					docCtrl.SetRenderer(r)
					debugView.SetText(fmt.Sprint("Renderer: ", r.Name()))
				}
			}(r))
	}

	flow := tview.NewFlex()
	flow.AddItem(topicList, 0, 1, true)
	flow.AddItem(docCtrl.View(), 0, 3, false)
	flow.AddItem(renderersList, 0, 1, false)

	debug := tview.NewFlex()
	debug.SetDirection(tview.FlexRow)
	debug.AddItem(flow, 0, 1, true)
	debug.AddItem(debugView, 1, 0, true)

	fc := NewFocusChain(topicList, docCtrl.View(), renderersList)
}
*/
