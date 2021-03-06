package widget

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// RadioButtons implements a simple primitive for radio button selections.
type RadioButtons struct {
	*tview.Box
	label          string
	options        []string
	currentOption  int
	onSelectedFunc func(text string, index int)
}

// NewRadioButtons returns a new radio button primitive.
func NewRadioButtons() *RadioButtons {
	return &RadioButtons{
		Box: tview.NewBox(),
	}
}

func (r *RadioButtons) SetLabel(label string) *RadioButtons {
	r.label = label
	return r
}

func (r *RadioButtons) SetOptions(options []string, handler func(text string, index int)) *RadioButtons {
	r.options = options
	r.onSelectedFunc = handler

	return r
}

func (r *RadioButtons) SetCurrentOption(index int) *RadioButtons {
	r.currentOption = index
	return r
}

// Draw draws this primitive onto the screen.
func (r *RadioButtons) Draw(screen tcell.Screen) {
	r.Box.DrawForSubclass(screen, r)
	x, y, width, height := r.GetInnerRect()

	if r.label != "" {
		tview.Print(screen, r.label, x, y, width, tview.AlignLeft, tcell.ColorWhite)
		y += 1
	}

	for index, option := range r.options {
		if index >= height {
			break
		}
		radioButton := "\u25ef" // Unchecked.
		if index == r.currentOption {
			radioButton = "\u25c9" // Checked.
		}

		colorTag := "[white]"
		if r.Box.HasFocus() && r.currentOption == index {
			colorTag = "[white::u]"
		}
		line := fmt.Sprintf(`%s %s%s`, radioButton, colorTag, option)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorWhite)
	}
}

// InputHandler returns the handler for this primitive.
func (r *RadioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = 0
			}
		case tcell.KeyDown:
			r.currentOption++
			if r.currentOption >= len(r.options) {
				r.currentOption = len(r.options) - 1
			}
		}

		if r.onSelectedFunc != nil {
			r.onSelectedFunc(r.options[r.currentOption], r.currentOption)
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (r *RadioButtons) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		_, rectY, _, _ := r.GetInnerRect()
		if !r.InRect(x, y) {
			return false, nil
		}

		if action == tview.MouseLeftClick {
			setFocus(r)
			index := y - rectY
			if index >= 0 && index < len(r.options) {
				r.currentOption = index
				consumed = true
			}
		}

		return
	})
}
