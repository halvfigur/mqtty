package view

import (
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	OpenFileViewController interface {
		OnFileSelected(filename string)
		OnOpenCancelled()
	}

	OpenFileView struct {
		*tview.Flex
		browser *widget.FileBrowser
	}
)

func NewOpenFileView(root string) *OpenFileView {

	b := widget.NewFileBrowser()
	b.SetDir(root)
	openButton := tview.NewButton("Open")
	cancelButton := tview.NewButton("Cancel")

	buttonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(openButton, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(cancelButton, 0, 1, false).
		AddItem(nil, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(b, 0, 1, true).
		AddItem(buttonFlex, 1, 0, false)

	//v.Flex = center(flex, 1, 1)
	flex.SetTitle("Open file").
		SetBorder(true)

	return &OpenFileView{
		Flex:    flex,
		browser: b,
	}
}

func (v *OpenFileView) SetOnFileSelected(handle func(filename string)) *OpenFileView {
	v.browser.SetOnFileSelected(handle)
	return v
}

func (v *OpenFileView) SetOnError(handle func(err error)) *OpenFileView {
	v.browser.SetOnError(handle)
	return v
}
