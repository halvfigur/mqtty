package view

import (
	"github.com/halvfigur/mqtty/model"
	"github.com/halvfigur/mqtty/widget"
)

type (
	RendererSelectionController interface {
		OnRendererSelected(renderer model.Renderer)
		GetView() *RendererSelector
	}

	RendererSelector struct {
		*widget.RadioButtons

		onSelected func(renderer model.Renderer)
		model      []model.Renderer
	}
)

func NewDocumentRenderer() *RendererSelector {
	s := widget.NewRadioButtons().SetLabel("[blue]Renderer[-]")

	return &RendererSelector{
		RadioButtons: s,
		model:        make([]model.Renderer, 0),
	}
}

func (s *RendererSelector) SetSelectedFunc(handler func(renderer model.Renderer)) *RendererSelector {
	s.onSelected = handler
	return s
}

func (s *RendererSelector) SetRenderers(renderers []model.Renderer) *RendererSelector {
	s.model = renderers

	options := make([]string, len(renderers))
	for i, s := range renderers {
		options[i] = s.Name()
	}
	s.SetOptions(options, func(text string, index int) {
		if s.onSelected != nil {
			s.onSelected(s.model[index])
		}
	})

	s.SetCurrentOption(0)

	return s
}
