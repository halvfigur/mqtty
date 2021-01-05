package view

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/rivo/tview"
)

type (
	RendererPageController interface {
		Renderers() []model.Renderer
		OnRenderer()
		Renderer(renderer model.Renderer)
		Cancel()
	}
)

func NewRendererPage(ctrl RendererPageController) *tview.Form {
	renderers := ctrl.Renderers()

	options := make([]string, len(renderers))
	for i, r := range renderers {
		options[i] = r.Name()
	}

	var renderer model.Renderer
	return tview.NewForm().
		AddDropDown("Renderer", options, 0, func(opt string, index int) {
			renderer = renderers[index]
		}).
		AddButton("OK", func() {
			ctrl.Renderer(renderer)
		}).
		AddButton("Cancel", func() {
			ctrl.Cancel()
		})
}
