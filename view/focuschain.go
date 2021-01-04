package view

import (
	"github.com/rivo/tview"
)

type (
	FocusChain struct {
		items []tview.Primitive
		index int
	}
)

func NewFocusChain(items ...tview.Primitive) *FocusChain {
	return &FocusChain{
		items: items,
		index: 0,
	}
}

func (c *FocusChain) Next() tview.Primitive {
	c.index = (c.index + 1) % len(c.items)
	return c.items[c.index]
}

func (c *FocusChain) Prev() tview.Primitive {
	c.index = (len(c.items) + c.index - 1) % len(c.items)
	return c.items[c.index]
}
