package view

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/rivo/tview"
)

type (
	RendererPageController interface {
		OnRendererSelected(renderer model.Renderer)
	}

	RendererPage struct {
		*tview.DropDown
		ctrl  RendererPageController
		model []model.Renderer
	}
)

func NewRendererPage(ctrl RendererPageController) *RendererPage {
	d := tview.NewDropDown().SetLabel("[blue]Renderer:[-] ")
	d.SetBorder(true).SetTitle("Select renderer")

	return &RendererPage{
		DropDown: d,
		ctrl:     ctrl,
		model:    make([]model.Renderer, 0),
	}
}

func (p *RendererPage) SetRenderers(renderers []model.Renderer) {
	p.model = renderers

	options := make([]string, len(renderers))
	for i, r := range renderers {
		options[i] = r.Name()
	}
	p.SetOptions(options, func(text string, index int) {
		p.ctrl.OnRendererSelected(p.model[index])
	})

	p.SetCurrentOption(0)
}
