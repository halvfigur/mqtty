package view

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
)

type (
	DocumentView struct {
		*tview.TextView

		doc *model.Document
	}
)

func NewDocumentView() *DocumentView {
	d := &DocumentView{
		TextView: tview.NewTextView(),
		doc:      model.NewDocument(),
	}

	d.TextView.SetTitle("Document").SetBorder(true)
	return d
}

func (v *DocumentView) SetRenderer(r model.Renderer) {
	v.doc.SetRenderer(r)
}

func (v *DocumentView) SetDocument(d *data.Document) {
	v.doc.SetDocument(d)
}

func (v *DocumentView) Refresh() {
	v.TextView.Clear()

	text, colorized := v.doc.Contents()

	v.TextView.SetDynamicColors(colorized)
	v.Write(text)
}
