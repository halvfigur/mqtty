package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/rivo/tview"
)

type (
	MainPageController interface {
		OnTopicSelected(t string)
		OnRendererSelected(r Renderer)
		OnChangeFocus(p tview.Primitive)
	}

	MainPage struct {
		*tview.Flex
		ctrl MainPageController

		topics    *tview.List
		doc       *DocumentView
		renderers *tview.List
	}
)

func NewMainPage(ctrl MainPageController) *MainPage {
	p := &MainPage{
		Flex:      tview.NewFlex(),
		topics:    tview.NewList(),
		doc:       NewDocumentView(),
		ctrl:      ctrl,
		renderers: tview.NewList(),
	}

	debugView := tview.NewTextView()

	/* Topic list */
	p.topics.SetBorder(true).SetTitle("Topics")
	p.topics.ShowSecondaryText(false)
	p.topics.AddItem("iotea/ingestion/events", "", 0, nil)
	p.topics.AddItem("iotea/discovery", "", 0, nil)

	/* Renderers list */
	renderers := []Renderer{
		new(RawRenderer),
		NewHexRenderer(),
	}

	p.renderers.SetBorder(true).SetTitle("Renderers")
	p.renderers.ShowSecondaryText(false)
	for _, r := range renderers {
		p.renderers.AddItem(r.Name(), "", 0,
			func(r Renderer) func() {
				return func() {
					p.ctrl.OnRendererSelected(r)
					debugView.SetText(fmt.Sprint("Renderer: ", r.Name()))
				}
			}(r))
	}

	p.Flex.SetDirection(tview.FlexColumn)
	p.Flex.AddItem(p.topics, 0, 1, true)
	p.Flex.AddItem(p.doc, 0, 3, false)
	p.Flex.AddItem(p.renderers, 0, 1, false)

	/*
		debug := tview.NewFlex()
		debug.SetDirection(tview.FlexRow)
		debug.AddItem(flow, 0, 1, true)
		debug.AddItem(debugView, 1, 0, true)
	*/
	fc := NewFocusChain(p.topics, p.doc, p.renderers)

	p.Flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			p.ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			p.ctrl.OnChangeFocus(fc.Prev())
		default:
			return event
		}

		return nil
	})

	return p
}

func (p *MainPage) SetDocument(d *model.Document) {
	p.doc.SetDocument(d)
	p.doc.Refresh()
	p.doc.ScrollToBeginning()
}

func (p *MainPage) SetRenderer(r Renderer) {
	p.doc.SetRenderer(r)
	p.doc.Refresh()
	p.doc.ScrollToBeginning()
}
