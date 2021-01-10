package control

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

const (
	mainPageLabel = "mainpage"
)

type (
	MainPageController struct {
		ctrl        Control
		mainView    *view.MainPage
		connectCtrl *StartPageController
		filtersCtrl *SubscriptionFiltersViewController
		publishCtrl *PublishPageController
		docModel    *model.Document
		documents   *model.DocumentStore
	}
)

func NewMainPageController(ctrl Control) *MainPageController {
	c := &MainPageController{
		ctrl:      ctrl,
		docModel:  model.NewDocument(),
		documents: model.NewDocumentStore(),
	}

	c.mainView = view.NewMainPage(c)
	c.mainView.SetDocumentStore(c.documents)

	c.connectCtrl = NewStartPageController(ctrl)
	c.filtersCtrl = NewSubscriptionFiltersViewController(ctrl)
	c.publishCtrl = NewPublishPageController(ctrl)

	ctrl.Register(mainPageLabel, c.mainView, false)

	return c
}

func (c *MainPageController) AddDocument(t string, d *data.Document) {
	c.documents.Store(t, d)
	c.mainView.AddTopic(t)
	c.mainView.Refresh()
}

func (c *MainPageController) OnTopicSelected(t string) {
	c.documents.SetCurrent(t)
	c.mainView.Refresh()
}

func (c *MainPageController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *MainPageController) OnNextDocument() {
	c.documents.Next()
	c.mainView.Refresh()
}

func (c *MainPageController) OnPrevDocument() {
	c.documents.Prev()
}

func (c *MainPageController) OnSubscribe() {
	c.ctrl.OnSubscribe()
}

func (c *MainPageController) OnPublish() {
	c.ctrl.OnPublish()
}

func (c *MainPageController) OnSetScrollToTop(enabled bool) {
}

func (c *MainPageController) OnSetFollow(enabled bool) {
	c.documents.Follow(enabled)
	c.mainView.Refresh()
}

func (c *MainPageController) OnConnect() {
	c.ctrl.OnConnect()
}

func (c *MainPageController) Cancel() {
	c.ctrl.Hide(subscriptionFiltersViewLabel)
}
