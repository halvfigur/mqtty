package view

import (
	"github.com/rivo/tview"
)

type ModalActionPair struct {
	Name string
	Func func()
}

func NewActionModal(text string, actions []ModalActionPair) *tview.Modal {
	buttons := make([]string, len(actions))
	for i, a := range actions {
		buttons[i] = a.Name
	}

	return tview.NewModal().SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(index int, label string) {
			actions[index].Func()
		})
}

type ModalPage struct {
	*tview.Flex
}

func NewModalPage() *ModalPage {
	return &ModalPage{
		tview.NewFlex(),
	}
}

func (p *ModalPage) SetModal(m *tview.Modal) {
	p.Clear()
	p.AddItem(m, 0, 0, true)
}
