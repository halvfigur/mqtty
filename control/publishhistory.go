package control

import (
	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
	"github.com/rivo/tview"
)

const publishHistoryLabel = "publishhistory"

type PublishHistory struct {
	ctrl      Control
	view      *view.PublishHistory
	documents *model.DocumentStore
}

func NewPublishHistory(ctrl Control) *PublishHistory {
	c := &PublishHistory{
		ctrl:      ctrl,
		documents: model.NewDocumentStore(),
	}

	c.view = view.NewPublishHistory(c)
	c.view.SetDocumentStore(c.documents)

	ctrl.Register(publishHistoryLabel, c.view, false)

	return c
}

func (h *PublishHistory) OnTopicSelected(t string) {
	h.documents.SetCurrent(t)
	h.view.Refresh()
}

func (h *PublishHistory) OnDocumentSelected(topic string, doc *data.Document) {
}

func (h *PublishHistory) OnChangeFocus(p tview.Primitive) {
	h.ctrl.Focus(p)
}

func (h *PublishHistory) OnNextDocument() {
	h.documents.Next()
	h.view.Refresh()
}

func (h *PublishHistory) OnPrevDocument() {
	h.documents.Prev()
	h.view.Refresh()
}

func (h *PublishHistory) Cancel() {
	h.ctrl.Cancel()
}

func (h *PublishHistory) AddDocument(t string, d *data.Document) {
	h.documents.Store(t, d)
	h.view.Refresh()
}
