package widget

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	FileBrowser struct {
		*tview.List
		root string

		onFileSelected func(filename string)
		onError        func(err error)
	}
)

func NewFileBrowser() *FileBrowser {
	b := &FileBrowser{
		List: tview.NewList().
			ShowSecondaryText(false).
			SetWrapAround(false),
	}

	return b
}

func (b *FileBrowser) SetOnFileSelected(handle func(filename string)) *FileBrowser {
	b.onFileSelected = handle
	return b
}

func (b *FileBrowser) SetOnError(handle func(err error)) *FileBrowser {
	b.onError = handle
	return b
}

func (b *FileBrowser) fileInfo(path string) (os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Stat()
}

func (b *FileBrowser) isDir(path string) (bool, error) {
	info, err := b.fileInfo(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func (b *FileBrowser) SetDir(dir string) error {
	b.Clear()

	b.root = dir

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	parent := filepath.Dir(dir)

	text := fmt.Sprint("[red]", b.root, "[-]")
	b.AddItem(text, "", 0, nil)

	if dir != "/" {
		b.AddItem("[green]..[-]", parent, 0, nil)
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		_ = path

		itemPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			text := fmt.Sprint("[green]", file.Name(), "[-]")
			b.AddItem(text, itemPath, 0, nil)
		} else {
			b.AddItem(file.Name(), itemPath, 0, nil)
		}
	}

	b.SetCurrentItem(1)
	return nil
}

func (b *FileBrowser) GetCurrentFile() string {
	_, path := b.GetItemText(b.GetCurrentItem())
	return path
}

func (b *FileBrowser) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		curr := b.GetCurrentItem()

		switch event.Key() {
		case tcell.KeyUp:
			// The first item is the name of the current directory
			if curr > 1 {
				b.SetCurrentItem(curr - 1)
			}
		case tcell.KeyDown:
			if curr < b.GetItemCount()-1 {
				b.SetCurrentItem(curr + 1)
			}
		case tcell.KeyRight, tcell.KeyEnter:
			_, path := b.GetItemText(b.GetCurrentItem())
			info, err := b.fileInfo(path)
			if err != nil {
				if b.onError != nil {
					b.onError(err)
				}
				return
			}

			if info.IsDir() {
				b.SetDir(path)
				return
			}

			if info.Mode().IsRegular() && b.onFileSelected != nil {
				b.onFileSelected(path)
			}
		case tcell.KeyLeft:
			parent := filepath.Dir(b.root)
			b.SetDir(parent)
		}
	})
}
