package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/widget"
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

		topicsList     *tview.List
		documentView   *Document
		connectionView *tview.TextView

		documents         *model.DocumentStore
		scrollToBeginning bool
	}
)

func NewCommander(ctrl CommanderController) *Commander {
	c := &Commander{
		Flex:         tview.NewFlex(),
		ctrl:         ctrl,
		topicsList:   tview.NewList(),
		documentView: NewDocumentView(),
	}

	/* Document columns */
	c.documentView.SetBorder(true)

	/* Topics column */
	c.topicsList.SetBorder(true).SetTitle("Topics")
	c.topicsList.ShowSecondaryText(false)
	c.topicsList.SetChangedFunc(func(index int, mainText, secondaryText string, short rune) {
		c.ctrl.OnTopicSelected(mainText)
	})

	/* Controls column */
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
			c.documentView.SetRenderer(renderer)
			c.Refresh()
		})

	c.connectionView = tview.NewTextView().
		SetDynamicColors(true)

	controlsFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(renderersView, 0, 1, false).
		AddItem(scrollToTopCheckbox, 1, 0, false).
		AddItem(followCheckbox, 1, 0, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(c.connectionView, 1, 0, false)

	controlsFlex.SetBorder(true).SetTitle("Controls")

	columnsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(c.topicsList, 0, 1, true).
		AddItem(c.documentView, 0, 3, false).
		AddItem(controlsFlex, 0, 1, false)

	fc := NewFocusChain(c.topicsList, c.documentView, renderersView, scrollToTopCheckbox, followCheckbox)

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
		case tcell.KeyCtrlO:
			c.ctrl.OnConnect()
		case tcell.KeyCtrlF:
			c.ctrl.OnSubscribe()
		case tcell.KeyCtrlP:
			c.ctrl.OnPublish()
		}

		return event
	})

	c.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(columnsFlex, 0, 3, true).
		AddItem(tview.NewTextView().
			SetDynamicColors(true).
			SetText("[blue](TAB):[-] Navigate  [blue](^O):[-] Open connection  [blue](^F):[-] Filters  [blue](^P):[-] Publish  [blue](^C):[-] Terminate"),
			1, 0, false)

	return c
}

func (c *Commander) SetDocumentStore(documents *model.DocumentStore) {
	c.documents = documents
}

func (c *Commander) SetConnectionStatus(s network.ConnectionStatus) {
	switch s {
	case network.StatusConnected:
		c.connectionView.SetText("[green]CONNECTED[-]")
	case network.StatusDisconnected:
		c.connectionView.SetText("[red]DISCONNECTED[-]")
	case network.StatusReconnecting:
		c.connectionView.SetText("[yellow]RECONNECTING[-]")
	}
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

	if c.topicsList.FindItems(t, subStringMatch, mustContainBoth, ignoreCase) != nil {
		return
	}

	const (
		secondaryText = ""
		shortCut      = 0
	)

	c.topicsList.AddItem(t, secondaryText, shortCut, nil)
}

func (c *Commander) setDocumentTitle(title string) {
	c.documentView.SetTitle(title)
}

func (c *Commander) setTopicsTitle(title string) {
	c.topicsList.SetTitle(title)
}

func (c *Commander) Refresh() {
	t, index := c.documents.Current()
	//c.setTopicsTitle(fmt.Sprintf("Topics %d", c.documents.Len()))
	c.setTopicsTitle(fmt.Sprintf("Topic %d/%d", c.topicsList.GetCurrentItem()+1, c.topicsList.GetItemCount()))
	if index == nil {
		c.setDocumentTitle("Document (none)")
		return
	}

	i, d := index.Current()
	c.setDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))

	c.documentView.SetDocument(d)
	c.documentView.Refresh()
	if c.scrollToBeginning {
		c.documentView.ScrollToBeginning()
	}

}
