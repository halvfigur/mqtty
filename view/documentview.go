package view

import (
	"github.com/rivo/tview"

	"github.com/halvfigur/mqtty/model"
)

type DocumentView struct {
	*tview.TextView

	doc *model.Document
}

func NewDocumentView() *DocumentView {
	d := &DocumentView{
		TextView: tview.NewTextView(),
	}

	d.TextView.SetTitle("Document").SetBorder(true)
	return d
}

func (v *DocumentView) SetDocument(d *model.Document) {
	v.doc = d
}

func (v *DocumentView) Refresh() {
	v.TextView.Clear()
	v.TextView.Write(v.doc.Contents())
}
