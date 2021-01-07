package control

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/view"
)

const rendererPageLabel = "rendererpage"

type (
	RendererPageController struct {
		ctrl      *MainPageController
		view      *view.RendererPage
		renderers []model.Renderer
	}
)

func NewRendererPageController(ctrl *MainPageController) *RendererPageController {
	c := &RendererPageController{
		ctrl: ctrl,
		renderers: []model.Renderer{
			model.NewRawRenderer(),
			model.NewHexRenderer(),
			model.NewJsonRenderer(),
		},
	}
	c.view = view.NewRendererPage(c)
	c.view.SetRenderers(c.renderers)

	return c
}

func (c *RendererPageController) GetView() *view.RendererPage {
	return c.view
}

func (c *RendererPageController) OnRendererSelected(renderer model.Renderer) {
	c.ctrl.OnRendererSelected(renderer)
}
