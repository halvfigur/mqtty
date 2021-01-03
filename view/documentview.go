package view

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/model"
)

type (
	Renderer interface {
		Name() string
		Render(doc *model.Document) ([]byte, bool)
	}

	RawRenderer struct{}

	DocumentView struct {
		*tview.TextView
		renderer Renderer

		doc *model.Document
	}
)

func (r *RawRenderer) Name() string {
	return "Raw"
}

func (r *RawRenderer) Render(d *model.Document) ([]byte, bool) {
	c := d.Contents()
	if c == nil {
		return nil, false
	}

	//tv.Write(d.Contents())
	return d.Contents(), false
}

func NewDocumentView() *DocumentView {
	d := &DocumentView{
		TextView: tview.NewTextView(),
		renderer: new(RawRenderer),
	}

	d.TextView.SetTitle("Document").SetBorder(true)
	return d
}

func (v *DocumentView) SetDocument(d *model.Document) {
	v.doc = d
}

func (v *DocumentView) SetRenderer(r Renderer) {
	v.renderer = r
	v.Refresh()
}

func (v *DocumentView) Refresh() {
	v.TextView.Clear()

	if v.renderer == nil {
		v.TextView.Write(v.doc.Contents())
		return
	}

	text, colorized := v.renderer.Render(v.doc)

	v.TextView.SetDynamicColors(colorized)
	v.Write(text)
}
