package control

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

const (
	commanderLabel = "commander"
)

type (
	CommanderController struct {
		ctrl        Control
		mainView    *view.Commander
		connectCtrl *Connector
		filtersCtrl *Filters
		publishCtrl *Publish
		docModel    *model.Document
		documents   *model.DocumentStore
	}
)

func NewMainPageController(ctrl Control) *CommanderController {
	c := &CommanderController{
		ctrl:      ctrl,
		docModel:  model.NewDocument(),
		documents: model.NewDocumentStore(),
	}

	c.mainView = view.NewMainPage(c)
	c.mainView.SetDocumentStore(c.documents)

	c.connectCtrl = NewConnector(ctrl)
	c.filtersCtrl = NewFilters(ctrl)
	c.publishCtrl = NewPublish(ctrl)

	ctrl.Register(commanderLabel, c.mainView, false)

	return c
}

func (c *CommanderController) AddDocument(t string, d *data.Document) {
	c.documents.Store(t, d)
	c.mainView.AddTopic(t)
	c.mainView.Refresh()
}

func (c *CommanderController) OnTopicSelected(t string) {
	c.documents.SetCurrent(t)
	c.mainView.Refresh()
}

func (c *CommanderController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *CommanderController) OnNextDocument() {
	c.documents.Next()
	c.mainView.Refresh()
}

func (c *CommanderController) OnPrevDocument() {
	c.documents.Prev()
}

func (c *CommanderController) OnSubscribe() {
	c.ctrl.OnSubscribe()
}

func (c *CommanderController) OnPublish() {
	c.ctrl.OnPublish()
}

func (c *CommanderController) OnSetScrollToTop(enabled bool) {
}

func (c *CommanderController) OnSetFollow(enabled bool) {
	c.documents.Follow(enabled)
	c.mainView.Refresh()
}

func (c *CommanderController) OnConnect() {
	c.ctrl.OnConnect()
}

func (c *CommanderController) Cancel() {
	c.ctrl.Hide(filtersLabel)
}
