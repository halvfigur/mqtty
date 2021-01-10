package widget

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Divider struct {
	*tview.Box
	label string
}

func NewDivider() *Divider {
	return &Divider{
		Box: tview.NewBox(),
	}
}

func (d *Divider) SetLabel(label string) *Divider {
	d.label = label
	return d
}

func (d *Divider) Draw(screen tcell.Screen) {
	d.Box.DrawForSubclass(screen, d)
	x, y, width, _ := d.GetInnerRect()

	lineLen := (width - len(d.label)) / 2
	if lineLen <= 0 {
		return
	}

	tail := (width + len(d.label)) % 2
	s1 := strings.Repeat(string(tview.Borders.Horizontal), lineLen+tail)
	s2 := strings.Repeat(string(tview.Borders.Horizontal), lineLen)

	line := fmt.Sprintf("%s%s%s", s1, d.label, s2)
	tview.Print(screen, line, x, y, width, tview.AlignLeft, tcell.ColorWhite)
}
