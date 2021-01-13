package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	PublishHistoryController interface {
		OnTopicSelected(t string)
		OnDocumentSelected(topic string, doc *data.Document)
		OnChangeFocus(p tview.Primitive)
		OnNextDocument()
		OnPrevDocument()

		Cancel()
	}

	PublishHistory struct {
		*tview.Flex
		ctrl      PublishHistoryController
		documents *model.DocumentStore

		topicList    *tview.List
		documentView *Document
	}
)

func NewPublishHistory(ctrl PublishHistoryController, documents *model.DocumentStore) *PublishHistory {

	topicList := tview.NewList().
		ShowSecondaryText(false)

	documentView := NewDocumentView()

	loadButton := tview.NewButton("Load").
		SetSelectedFunc(func() {
			t, index := documents.Current()
			if index == nil {
				return
			}

			_, doc := index.Current()
			ctrl.OnDocumentSelected(t, doc)
		})

	cancelButton := tview.NewButton("Cancel").
		SetSelectedFunc(func() {
			ctrl.Cancel()
		})

	buttonFlex := Space(tview.FlexColumn, loadButton, cancelButton)

	fc := NewFocusChain(topicList, documentView, loadButton, cancelButton)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topicList, 0, 1, true).
		AddItem(widget.NewDivider().SetLabel("Document"), 1, 0, false).
		AddItem(documentView, 0, 3, false).
		AddItem(widget.NewDivider(), 1, 0, false).
		AddItem(buttonFlex, 1, 0, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			ctrl.OnChangeFocus(fc.Prev())
		case tcell.KeyRight:
			ctrl.OnNextDocument()
		case tcell.KeyLeft:
			ctrl.OnPrevDocument()
		}
		return event
	}).SetBorder(true).SetTitle("Topics")

	return &PublishHistory{
		Flex:         Center(flex, 500, 200),
		ctrl:         ctrl,
		documents:    documents,
		topicList:    topicList,
		documentView: documentView,
	}
}

func (h *PublishHistory) refreshTopics() {
	h.SetTitle(fmt.Sprintf("Topics %d", h.documents.Len()))

	h.topicList.Clear()
	for _, t := range h.documents.Topics() {
		h.topicList.AddItem(t, "", 0, func() {
			h.ctrl.OnTopicSelected(t)
		})
	}
}

func (h *PublishHistory) refreshDocument() {
	t, index := h.documents.Current()
	if index == nil {
		h.documentView.SetTitle("Document (none)")
		return
	}

	i, d := index.Current()
	h.documentView.SetTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))

	h.documentView.SetDocument(d)
	h.documentView.Refresh()
}

func (h *PublishHistory) Refresh() {
	h.refreshTopics()
	h.refreshDocument()
}
