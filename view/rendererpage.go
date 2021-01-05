package view

import (
	"github.com/gdamore/tcell/v2"
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

func NewRendererPage(ctrl RendererPageController) *tview.Flex {
	renderers := ctrl.Renderers()

	d := tview.NewDropDown().SetLabel("[blue]Qos:[-] ")
	d.SetBorder(true).SetTitle("Select renderer")
	for _, r := range renderers {
		d.AddOption(r.Name(), nil)
	}

	d.SetCurrentOption(0)
	d.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			i, _ := d.GetCurrentOption()
			ctrl.Renderer(renderers[i])
			ctrl.Cancel()
		case tcell.KeyEscape:
			ctrl.Cancel()
		}

		return event
	})

	return center(d, 1, 1)
}
