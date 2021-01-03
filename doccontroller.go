package main

import (
	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

type (
	DocumentController struct {
		view  *view.DocumentView
		model *model.Document
	}
)

func NewDocumentController() *DocumentController {
	return &DocumentController{
		view:  view.NewDocumentView(),
		model: model.NewDocument(),
	}
}

func (c *DocumentController) SetDocument(d *data.Document) {
	c.model.SetDocument(d)
	c.view.SetDocument(c.model)
	c.view.Refresh()
	c.view.ScrollToBeginning()
}

func (c *DocumentController) SetRenderer(r view.Renderer) {
	c.view.SetRenderer(r)
	c.view.Refresh()
	c.view.ScrollToBeginning()
}

func (c *DocumentController) View() *view.DocumentView {
	return c.view
}
