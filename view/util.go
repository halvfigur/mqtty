package view

import (
	"github.com/rivo/tview"
)

type FocusChain struct {
	items []tview.Primitive
	index int
}

func NewFocusChain(items ...tview.Primitive) *FocusChain {
	return &FocusChain{
		items: items,
		index: 0,
	}
}

func (c *FocusChain) Add(items ...tview.Primitive) {
	c.items = append(c.items, items...)
}

func (c *FocusChain) Next() tview.Primitive {
	c.index = (c.index + 1) % len(c.items)
	return c.items[c.index]
}

func (c *FocusChain) Prev() tview.Primitive {
	c.index = (len(c.items) + c.index - 1) % len(c.items)
	return c.items[c.index]
}

func (c *FocusChain) SetFocus(index int) tview.Primitive {
	c.index = index
	return c.items[c.index]
}

func (c *FocusChain) Reset() tview.Primitive {
	c.index = 0
	return c.items[c.index]
}

func Center(p tview.Primitive, rowProportion, colProportion int) *tview.Flex {
	cols := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(p, 0, colProportion, true).
		AddItem(nil, 0, 1, false)
	rows := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(cols, 0, rowProportion, true).
		AddItem(nil, 0, 1, false)

	rows.SetBorder(true)
	return rows
}

func Space(direction int, ps ...tview.Primitive) *tview.Flex {
	flex := tview.NewFlex()

	if len(ps) == 0 {
		return flex
	}

	flex.SetDirection(direction)
	for _, p := range ps {
		flex.AddItem(tview.NewTextView(), 0, 1, false)
		flex.AddItem(p, 0, 1, false)
	}
	flex.AddItem(tview.NewTextView(), 0, 1, false)

	return flex
}
