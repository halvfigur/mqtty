package view

import (
	"github.com/rivo/tview"
)

type ErrorModalControl interface {
	Cancel()
}

func NewErrorModal(ctrl ErrorModalControl) *tview.Modal {
	m := tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(index int, label string) {
			ctrl.Cancel()
		})

	m.SetTitle("Error").SetBorder(true)

	return m
}
