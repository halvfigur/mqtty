package control

import (
	"fmt"

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

	c.mainView = view.NewMainPage(c, NewRendererPageController(c))

	c.connectCtrl = NewStartPageController(ctrl)
	c.filtersCtrl = NewSubscriptionFiltersViewController(ctrl)

	ctrl.Register(mainPageLabel, c.mainView, false)
	ctrl.Register(subscriptionFiltersViewLabel, c.filtersCtrl.GetView(), false)
	ctrl.Register(startPageLabel, c.connectCtrl.GetView(), true)

	return c
}

func (c *MainPageController) SetDocument(d *data.Document) {
	c.docModel.SetDocument(d)
	if c.mainView != nil {
		c.mainView.SetDocument(c.docModel)
	}
}

func (c *MainPageController) AddDocument(t string, d *data.Document) {
	c.documents.Store(t, d)

	if c.mainView != nil {
		c.mainView.AddTopic(t)
		c.updateDocumentView()
	}
}

func (c *MainPageController) OnTopicSelected(t string) {
	c.documents.SetCurrent(t)
	c.updateDocumentView()
}

func (c *MainPageController) OnChangeFocus(p tview.Primitive) {
	c.ctrl.Focus(p)
}

func (c *MainPageController) OnNextDocument() {
	_, index := c.documents.Current()
	index.Next()
	c.updateDocumentView()
}

func (c *MainPageController) OnPrevDocument() {
	_, index := c.documents.Current()
	index.Prev()
	c.updateDocumentView()
}

func (c *MainPageController) OnSubscribe() {
	c.ctrl.Display(subscriptionFiltersViewLabel)
}

func (c *MainPageController) OnRendererSelected(renderer model.Renderer) {
	c.docModel.SetRenderer(renderer)
	c.updateDocumentView()
}

func (c *MainPageController) updateDocumentView() {

	t, index := c.documents.Current()
	if index == nil {
		return
	}

	i, d := index.Current()

	c.docModel.SetDocument(d)

	c.mainView.SetDocument(c.docModel)
	c.mainView.SetTopicsTitle(fmt.Sprintf("Topics %d", c.documents.Len()))
	c.mainView.SetDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))
}

func (c *MainPageController) Cancel() {
	c.ctrl.Hide(subscriptionFiltersViewLabel)
}
