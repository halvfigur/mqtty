package view

import (
	"github.com/rivo/tview"
)

type ModalControl interface {
	Cancel()
}

func NewErrorModal(ctrl ModalControl) *tview.Modal {
	m := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(index int, label string) {
			ctrl.Cancel()
		})

	m.SetTitle("Error").SetBorder(true)

	return m
}

func NewWaitModal() *tview.Modal {
	m := tview.NewModal()

	m.SetTitle("...").SetBorder(true)

	return m
}
