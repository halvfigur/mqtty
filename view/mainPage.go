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
		OnNextDocument()
		OnPrevDocument()
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

	/* Topics list */
	p.topics.SetBorder(true).SetTitle("Topics")
	p.topics.ShowSecondaryText(false)

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
		case tcell.KeyRight:
			p.ctrl.OnNextDocument()
		case tcell.KeyLeft:
			p.ctrl.OnPrevDocument()
		default:
			return event
		}

		return nil
	})

	return p
}

func (p *MainPage) Focus(delegate func(p tview.Primitive)) {
	p.Flex.SetFullScreen(true)
	delegate(p.Flex)
}

func (p *MainPage) AddTopic(t string) {
	const (
		subStringMatch  = ""
		mustContainBoth = false
		ignoreCase      = false
	)

	if p.topics.FindItems(t, subStringMatch, mustContainBoth, ignoreCase) != nil {
		return
	}

	const (
		secondaryText = ""
		shortCut      = 0
	)

	p.topics.AddItem(t, secondaryText, shortCut, func() {
		p.ctrl.OnTopicSelected(t)
	})
}

func (p *MainPage) SetDocumentTitle(title string) {
	p.doc.SetTitle(title)
}

func (p *MainPage) SetTopicsTitle(title string) {
	p.topics.SetTitle(title)
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