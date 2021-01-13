package control

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/network"
	"github.com/halvfigur/mqtty/view"
)

type CommanderController struct {
	ctrl        Control
	mainView    *view.Commander
	connectCtrl *Connector
	filtersCtrl *Filters
	publishCtrl *Publish
	documents   *model.DocumentStore
}

func NewCommanderController(ctrl Control) *CommanderController {
	c := &CommanderController{
		ctrl:      ctrl,
		documents: model.NewDocumentStore(),
	}

	c.mainView = view.NewCommander(c)
	c.mainView.SetDocumentStore(c.documents)
	c.mainView.SetConnectionStatus(network.StatusDisconnected)

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

func (c *CommanderController) SetConnectionStatus(s network.ConnectionStatus) {
	c.mainView.SetConnectionStatus(s)
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
	c.ctrl.OnDisplaySubscriber()
}

func (c *CommanderController) OnPublish() {
	c.ctrl.OnDisplayPublisher()
}

func (c *CommanderController) OnSetScrollToTop(enabled bool) {
}

func (c *CommanderController) OnSetFollow(enabled bool) {
	c.documents.Follow(enabled)
	c.mainView.Refresh()
}

func (c *CommanderController) OnConnect() {
	c.ctrl.OnDisplayConnector()
}
