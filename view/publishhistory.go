package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
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
	}
)

func NewPublishHistory(ctrl PublishHistoryController) *PublishHistory {
	h := &PublishHistory{
		ctrl: ctrl,
	}

	topicList := tview.NewList().
		ShowSecondaryText(false).
		SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			ctrl.OnTopicSelected(mainText)
		})
	topicList.SetTitle("Topic").
		SetBorder(true)

	documentView := NewDocumentView()
	documentView.SetTitle("Document").
		SetBorder(true)

	loadButton := tview.NewButton("Load").
		SetSelectedFunc(func() {
			t, index := h.documents.Current()
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

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(loadButton, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(cancelButton, 0, 1, false).
		AddItem(nil, 0, 1, false)

	fc := NewFocusChain(topicList, documentView, loadButton, cancelButton)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topicList, 0, 1, true).
		AddItem(documentView, 0, 3, false).
		AddItem(buttonFlex, 1, 0, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			ctrl.OnChangeFocus(fc.Next())
		case tcell.KeyBacktab:
			ctrl.OnChangeFocus(fc.Prev())
		}
		return event
	})

	h.Flex = Center(flex, 6, 3)

	return h
}

func (h *PublishHistory) SetDocumentStore(documents *model.DocumentStore) {
	h.documents = documents
}

func (h *PublishHistory) Refresh() {
}
