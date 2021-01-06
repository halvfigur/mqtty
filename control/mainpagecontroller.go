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
		ctrl      Control
		view      *view.MainPage
		model     *model.Document
		renderer  model.Renderer
		documents *model.DocumentStore
	}
)

func NewMainPageController(ctrl Control) *MainPageController {
	return &MainPageController{
		ctrl:      ctrl,
		model:     model.NewDocument(),
		documents: model.NewDocumentStore(),
		renderer:  ctrl.Renderers()[0],
	}
}

func (c *MainPageController) SetView(v *view.MainPage) {
	c.view = v
}

func (c *MainPageController) SetDocument(d *data.Document) {
	c.model.SetDocument(d)
	if c.view != nil {
		c.view.SetDocument(c.model)
	}
}

func (c *MainPageController) AddDocument(t string, d *data.Document) {
	c.documents.Store(t, d)

	if c.view != nil {
		c.view.AddTopic(t)
		c.updateDocumentView()
	}
}

func (c *MainPageController) OnTopicSelected(t string) {
	c.documents.SetCurrent(t)
	c.updateDocumentView()
}

func (c *MainPageController) SetRenderer(r model.Renderer) {
	c.renderer = r
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
	c.ctrl.OnSubscribe()
}

func (c *MainPageController) OnRenderer() {
	c.ctrl.OnRenderer()
}

func (c *MainPageController) updateDocumentView() {
	c.model.SetRenderer(c.renderer)

	t, index := c.documents.Current()
	if index == nil {
		return
	}

	i, d := index.Current()

	c.model.SetDocument(d)

	c.view.SetDocument(c.model)
	c.view.SetTopicsTitle(fmt.Sprintf("Topics %d", c.documents.Len()))
	c.view.SetDocumentTitle(fmt.Sprintf("%s (%d/%d)", t, i+1, index.Len()))
}
