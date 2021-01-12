package view

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/halvfigur/mqtty/widget"
	"github.com/rivo/tview"
)

type (
	OpenFileController interface {
		OnChangeFocus(p tview.Primitive)
		OnFileSelected(filename string)
		Cancel()
	}

	OpenFile struct {
		*tview.Flex
		browser *widget.FileBrowser
	}
)

func NewOpenFile(ctrl OpenFileController) *OpenFile {

	fileBrowser := widget.NewFileBrowser().
		SetOnFileSelected(func(filename string) {
			ctrl.OnFileSelected(filename)
		})

	openButton := tview.NewButton("Open").
		SetSelectedFunc(func() {
			file := fileBrowser.GetCurrentFile()
			f, err := os.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				return
			}

			if info.IsDir() {
				fileBrowser.SetDir(file)
				return
			}

			ctrl.OnFileSelected(file)
		})

	cancelButton := tview.NewButton("Cancel").
		SetSelectedFunc(func() {
			ctrl.Cancel()
		})

		/*
			buttonFlex := tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(nil, 0, 1, false).
				AddItem(openButton, 0, 1, false).
				AddItem(nil, 0, 1, false).
				AddItem(cancelButton, 0, 1, false).
				AddItem(nil, 0, 1, false)
		*/
	buttonFlex := Space(tview.FlexColumn, openButton, cancelButton)

	fc := NewFocusChain(fileBrowser, openButton, cancelButton)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(fileBrowser, 0, 1, true).
		AddItem(buttonFlex, 1, 0, false)

	flex.SetTitle("Open file").
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyTab:
				ctrl.OnChangeFocus(fc.Next())
			case tcell.KeyBacktab:
				ctrl.OnChangeFocus(fc.Prev())
			}

			return event
		})

	return &OpenFile{
		Flex:    flex,
		browser: fileBrowser,
	}
}

func (v *OpenFile) SetDir(dir string) *OpenFile {
	v.browser.SetDir(dir)
	return v
}

func (v *OpenFile) SetOnFileSelected(handle func(filename string)) *OpenFile {
	v.browser.SetOnFileSelected(handle)
	return v
}

func (v *OpenFile) SetOnError(handle func(err error)) *OpenFile {
	v.browser.SetOnError(handle)
	return v
}
