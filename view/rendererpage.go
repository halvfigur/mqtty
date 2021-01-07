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

		onSelected func(renderer model.Renderer)
		model      []model.Renderer
	}
)

func NewRendererPage() *RendererPage {
	/*
		d := tview.NewDropDown().SetLabel("[blue]Renderer:[-] ")
		d.SetBorder(true).SetTitle("Select renderer")
	*/
	r := widget.NewRadioButtons().SetLabel("[blue]Renderer[-]")

	return &RendererPage{
		//DropDown: d,
		RadioButtons: r,
		model:        make([]model.Renderer, 0),
	}
}

func (p *RendererPage) SetSelectedFunc(handler func(renderer model.Renderer)) *RendererPage {
	p.onSelected = handler
	return p
}

func (p *RendererPage) SetRenderers(renderers []model.Renderer) *RendererPage {
	p.model = renderers

	options := make([]string, len(renderers))
	for i, r := range renderers {
		options[i] = r.Name()
	}
	p.SetOptions(options, func(text string, index int) {
		if p.onSelected != nil {
			p.onSelected(p.model[index])
		}
	})

	p.SetCurrentOption(0)

	return p
}
