package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/rivo/tview"
)

type (
	MainPageController interface {
		OnTopicSelected(t string)
		OnChangeFocus(p tview.Primitive)
		OnNextDocument()
		OnPrevDocument()
		OnSubscribe()
		OnRendererSelected(renderer model.Renderer)
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

	/* Topics list */
	p.topics.SetBorder(true).SetTitle("Topics")
	p.topics.ShowSecondaryText(false)

	p.Flex.SetDirection(tview.FlexRow)
	p.Flex.AddItem(p.topics, 0, 1, true)
	p.Flex.AddItem(p.doc, 0, 3, false)
	p.Flex.AddItem(tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue](TAB):[-] navigate  [blue](F):[-] filters [blue](R):[-] renderer"),
		1, 0, false)

	fc := NewFocusChain(p.topics, p.doc)

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
		}

		switch event.Rune() {
		case 'f', 'F':
			p.ctrl.OnSubscribe()
		}

		return event
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
}

func (p *MainPage) Refresh() {
	p.doc.Refresh()
	p.doc.ScrollToBeginning()
}
