package widget

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PrettyCheckbox struct {
	*tview.Box

	label   string
	checked bool

	onChangedFunc func(bool)
}

func NewPrettyCheckbox() *PrettyCheckbox {
	return &PrettyCheckbox{
		Box: tview.NewBox(),
	}
}

func (b *PrettyCheckbox) SetLabel(label string) *PrettyCheckbox {
	b.label = label
	return b
}

func (b *PrettyCheckbox) SetChecked(checked bool) *PrettyCheckbox {
	b.checked = checked
	return b
}

func (b *PrettyCheckbox) SetChangedFunc(f func(bool)) *PrettyCheckbox {
	b.onChangedFunc = f
	return b
}

func (b *PrettyCheckbox) IsChecked() bool {
	return b.checked
}

func (b *PrettyCheckbox) Draw(screen tcell.Screen) {
	b.Box.DrawForSubclass(screen, b)
	x, y, width, _ := b.GetInnerRect()

	var text string
	if b.HasFocus() {
		text = fmt.Sprint("[white::u]", b.label, "[-:-:-]")
	} else {
		text = b.label
	}

	if b.checked {
		text = fmt.Sprint(text, " \u25a0")
	} else {
		text = fmt.Sprint(text, " \u25a1")
	}

	if b.label != "" {
		tview.Print(screen, text, x, y, width, tview.AlignLeft, tcell.ColorWhite)
		y += 1
	}
}

// InputHandler returns the handler for this primitive.
func (b *PrettyCheckbox) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyEnter || event.Rune() == ' ' {
			b.checked = !b.checked
		}

		if b.onChangedFunc != nil {
			b.onChangedFunc(b.checked)
		}
	})
}
