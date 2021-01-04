package main

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

type (
	MainPageController struct {
		app   *tview.Application
		view  *view.MainPage
		model *model.Document
	}
)

func NewMainPageController(a *tview.Application) *MainPageController {
	return &MainPageController{
		app:   a,
		model: model.NewDocument(),
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

func (c *MainPageController) OnTopicSelected(t string) {
}

func (c *MainPageController) OnRendererSelected(r view.Renderer) {
	if c.view != nil {
		c.view.SetRenderer(r)
	}
}

func (c *MainPageController) OnChangeFocus(p tview.Primitive) {
	c.app.SetFocus(p)
}
