package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
)

func main() {
	buf := bytes.NewBufferString(`
	{
		"alpha": {
			"beta": [0, 1, 2, 3},
			"gamma": "hello world",
			"epsilon": null
		}
	}
	`)

	// Document data
	d, err := data.NewDocument(buf)

	if err != nil {
		log.Fatal("Failed to read document")
	}
	app := tview.NewApplication()
	debugView := tview.NewTextView()

	docCtrl := NewDocumentController()
	docCtrl.SetDocument(d)

	/* Topic list */
	topicList := tview.NewList()
	topicList.SetBorder(true).SetTitle("Topics")
	topicList.ShowSecondaryText(false)
	topicList.AddItem("iotea/ingestion/events", "", 0, nil)
	topicList.AddItem("iotea/discovery", "", 0, nil)

	/* Renderers list */
	renderers := []model.Renderer{
		new(model.RawRenderer),
		new(model.HexRenderer),
	}

	renderersList := tview.NewList()
	renderersList.SetBorder(true).SetTitle("Renderers")
	renderersList.ShowSecondaryText(false)
	for _, r := range renderers {
		renderersList.AddItem(r.Name(), "", 0,
			func(r model.Renderer) func() {
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
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		debugView.SetText(event.Name())

		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(fc.Next())
		case tcell.KeyBacktab:
			app.SetFocus(fc.Prev())
		default:
			return event
		}

		return nil
	})

	app.SetRoot(debug, true)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
