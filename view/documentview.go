package view

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/model"
)

type (
	Renderer interface {
		Name() string
		Render(doc *model.Document, tv *tview.TextView)
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

func (r *RawRenderer) Render(d *model.Document, tv *tview.TextView) {
	c := d.Contents()
	if c == nil {
		return
	}

	tv.Write(d.Contents())
}

func NewDocumentView() *DocumentView {
	d := &DocumentView{
		TextView: tview.NewTextView(),
		renderer: new(RawRenderer),
	}

	d.TextView.SetTitle("Document").SetBorder(true)
	d.TextView.SetDynamicColors(true)
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

	v.renderer.Render(v.doc, v.TextView)
}
