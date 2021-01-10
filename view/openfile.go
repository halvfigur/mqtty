package view

import (
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	OpenFileController interface {
		OnFileSelected(filename string)
		OnOpenCancelled()
	}

	OpenFile struct {
		*tview.Flex
		browser *widget.FileBrowser
	}
)

func NewOpenFile(root string) *OpenFile {

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

	return &OpenFile{
		Flex:    flex,
		browser: b,
	}
}

func (v *OpenFile) SetOnFileSelected(handle func(filename string)) *OpenFile {
	v.browser.SetOnFileSelected(handle)
	return v
}

func (v *OpenFile) SetOnError(handle func(err error)) *OpenFile {
	v.browser.SetOnError(handle)
	return v
}
