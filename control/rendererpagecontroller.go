package control

import "github.com/halvfigur/mqtty/model"

const rendererPageLabel = "rendererpage"

type (
	RendererPageController struct {
		ctrl Control
	}
)

func NewRendererPageController(ctrl Control) *RendererPageController {
	return &RendererPageController{
		ctrl: ctrl,
	}
}

func (c *RendererPageController) Renderers() []model.Renderer {
	return c.ctrl.Renderers()
}

func (c *RendererPageController) OnRenderer() {
	c.ctrl.OnRenderer()
}

func (c *RendererPageController) Renderer(renderer model.Renderer) {
	c.ctrl.Renderer(renderer)
}

func (c *RendererPageController) Cancel() {
	c.ctrl.Cancel()
}
