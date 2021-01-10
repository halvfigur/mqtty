package view

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/data"
	"github.com/halvfigur/mqtty/model"
)

type (
	Document struct {
		*tview.TextView

		doc *model.Document
	}
)

func NewDocumentView() *Document {
	d := &Document{
		TextView: tview.NewTextView(),
		doc:      model.NewDocument(),
	}

	d.TextView.SetTitle("Document").SetBorder(true)
	return d
}

func (v *Document) SetRenderer(r model.Renderer) {
	v.doc.SetRenderer(r)
}

func (v *Document) SetDocument(d *data.Document) {
	v.doc.SetDocument(d)
}

func (v *Document) Refresh() {
	v.TextView.Clear()

	text, colorized := v.doc.Contents()

	v.TextView.SetDynamicColors(colorized)
	v.Write(text)
}
