package view

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/widget"
)

type (
	RendererPageController interface {
		OnRendererSelected(renderer model.Renderer)
		GetView() *RendererPage
	}

	RendererPage struct {
		//*tview.DropDown
		*widget.RadioButtons

		ctrl  RendererPageController
		model []model.Renderer
	}
)

func NewRendererPage(ctrl RendererPageController) *RendererPage {
	/*
		d := tview.NewDropDown().SetLabel("[blue]Renderer:[-] ")
		d.SetBorder(true).SetTitle("Select renderer")
	*/
	r := widget.NewRadioButtons().SetHeader("[blue]Renderer[-]")

	return &RendererPage{
		//DropDown: d,
		RadioButtons: r,
		ctrl:         ctrl,
		model:        make([]model.Renderer, 0),
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
