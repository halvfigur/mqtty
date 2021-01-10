package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/rivo/tview"
)

type (
	CommanderController interface {
		OnTopicSelected(t string)
		OnChangeFocus(p tview.Primitive)
		OnNextDocument()
		OnPrevDocument()
		OnConnect()
		OnSubscribe()
		OnPublish()
		OnSetFollow(enabled bool)
	}

	Commander struct {
		*tview.Flex
		ctrl CommanderController

		topics            *tview.List
		docView           *Document
		documents         *model.DocumentStore
		scrollToBeginning bool
	}
)

func NewMainPage(ctrl CommanderController) *Commander {
	c := &Commander{
		Flex:    tview.NewFlex(),
		topics:  tview.NewList(),
		docView: NewDocumentView(),
		ctrl:    ctrl,
	}

	/* Topics list */
	c.topics.SetBorder(true).SetTitle("Topics")
	c.topics.ShowSecondaryText(false)
	c.topics.SetChangedFunc(func(index int, mainText, secondaryText string, short rune) {
		c.ctrl.OnTopicSelected(mainText)
	})

	scrollToTopCheckbox := tview.NewCheckbox().
		SetLabel("Scroll to top: ").
		SetChangedFunc(func(checked bool) { c.scrollToBeginning = checked })

	followCheckbox := tview.NewCheckbox().
		SetLabel("Follow: ").
		SetChangedFunc(ctrl.OnSetFollow)

	renderersView := NewDocumentRenderer().
		SetRenderers([]model.Renderer{
			model.NewRawRenderer(),
			model.NewHexRenderer(),
			model.NewJsonRenderer(),
		}).
		SetSelectedFunc(func(renderer model.Renderer) {
			c.docView.SetRenderer(renderer)
			c.Refresh()
		})
	controlsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	controlsFlex.SetBorder(true).SetTitle("Controls")
	controlsFlex.AddItem(renderersView, 0, 1, false).
		AddItem(scrollToTopCheckbox, 0, 1, false).
		AddItem(followCheckbox, 0, 1, false)

	columnsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(c.topics, 0, 1, true).
		AddItem(c.docView, 0, 3, false).
		AddItem(controlsFlex, 0, 1, false)
	fc := NewFocusChain(c.topics, c.docView, renderersView, scrollToTopCheckbox, followCheckbox)

	columnsFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			c.ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			c.ctrl.OnChangeFocus(fc.Prev())
		case tcell.KeyRight:
			c.ctrl.OnNextDocument()
		case tcell.KeyLeft:
			c.ctrl.OnPrevDocument()
		}

		switch event.Rune() {
		case 'f', 'F':
			c.ctrl.OnSubscribe()
		case 'p', 'P':
			c.ctrl.OnPublish()
		}

		return event
	})

	c.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(columnsFlex, 0, 3, true).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[blue](TAB):[-] navigate  [blue](F):[-] filters  [blue](P):[-] publish"),
			1, 0, false)

	return c
}

func (c *Commander) SetDocumentStore(documents *model.DocumentStore) {
	c.documents = documents
}

func (c *Commander) Focus(delegate func(p tview.Primitive)) {
	c.Flex.SetFullScreen(true)
	delegate(c.Flex)
}

func (c *Commander) AddTopic(t string) {
	const (
		subStringMatch  = ""
		mustContainBoth = false
		ignoreCase      = false
	)

	if c.topics.FindItems(t, subStringMatch, mustContainBoth, ignoreCase) != nil {
		return
	}

	const (
		secondaryText = ""
		shortCut      = 0
	)

	c.topics.AddItem(t, secondaryText, shortCut, nil)
}

func (c *Commander) setDocumentTitle(title string) {
	c.docView.SetTitle(title)
}

func (c *Commander) setTopicsTitle(title string) {
	c.topics.SetTitle(title)
}

func (c *Commander) Refresh() {

	t, index := c.documents.Current()
	c.setTopicsTitle(fmt.Sprintf("Topics %d", c.documents.Len()))
	if index == nil {
		c.setDocumentTitle("Document (none)")
		return
	}

	i, d := index.Current()
	c.setDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))

	c.docView.SetDocument(d)
	c.docView.Refresh()
	if c.scrollToBeginning {
		c.docView.ScrollToBeginning()
	}

}
