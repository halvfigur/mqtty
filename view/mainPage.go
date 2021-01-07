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
		OnChangeFocus(p tview.Primitive)
		OnNextDocument()
		OnPrevDocument()
		OnSubscribe()
		OnSetFollow(enabled bool)
	}

	MainPage struct {
		*tview.Flex
		ctrl MainPageController

		topics            *tview.List
		docView           *DocumentView
		documents         *model.DocumentStore
		scrollToBeginning bool
	}
)

func NewMainPage(ctrl MainPageController) *MainPage {
	p := &MainPage{
		Flex:    tview.NewFlex(),
		topics:  tview.NewList(),
		docView: NewDocumentView(),
		ctrl:    ctrl,
	}

	/* Topics list */
	p.topics.SetBorder(true).SetTitle("Topics")
	p.topics.ShowSecondaryText(false)
	p.topics.SetChangedFunc(func(index int, mainText, secondaryText string, short rune) {
		p.ctrl.OnTopicSelected(mainText)
	})

	scrollToTopCheckbox := tview.NewCheckbox().
		SetLabel("Scroll to top: ").
		SetChangedFunc(func(checked bool) { p.scrollToBeginning = checked })

	followCheckbox := tview.NewCheckbox().
		SetLabel("Follow: ").
		SetChangedFunc(ctrl.OnSetFollow)

	renderersView := NewRendererPage().
		SetRenderers([]model.Renderer{
			model.NewRawRenderer(),
			model.NewHexRenderer(),
			model.NewJsonRenderer(),
		}).
		SetSelectedFunc(func(renderer model.Renderer) {
			p.docView.SetRenderer(renderer)
			p.Refresh()
		})
	controlsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	controlsFlex.SetBorder(true).SetTitle("Controls")
	controlsFlex.AddItem(renderersView, 0, 1, false).
		AddItem(scrollToTopCheckbox, 0, 1, false).
		AddItem(followCheckbox, 0, 1, false)

	columnsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(p.topics, 0, 1, true).
		AddItem(p.docView, 0, 3, false).
		AddItem(controlsFlex, 0, 1, false)
	fc := NewFocusChain(p.topics, p.docView, renderersView, scrollToTopCheckbox, followCheckbox)

	columnsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	p.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(columnsFlex, 0, 3, true).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[blue](TAB):[-] navigate  [blue](F):[-] filters [blue](R):[-] renderer"),
			1, 0, false)

	return p
}

func (p *MainPage) SetDocumentStore(documents *model.DocumentStore) {
	p.documents = documents
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

	p.topics.AddItem(t, secondaryText, shortCut, nil)
	/*
		p.topics.AddItem(t, secondaryText, shortCut, func() {
			p.ctrl.OnTopicSelected(t)
		})
	*/
}

func (p *MainPage) SetDocumentTitle(title string) {
	p.docView.SetTitle(title)
}

func (p *MainPage) SetTopicsTitle(title string) {
	p.topics.SetTitle(title)
}

func (p *MainPage) Refresh() {
	t, index := p.documents.Current()
	if index == nil {
		return
	}

	i, d := index.Current()

	p.docView.SetDocument(d)
	p.docView.Refresh()
	if p.scrollToBeginning {
		p.docView.ScrollToBeginning()
	}

	p.SetTopicsTitle(fmt.Sprintf("Topics %d", p.documents.Len()))
	p.SetDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))
}
