package model

import (
	"github.com/halvfigur/mqtty/data"
)

type (
	Document struct {
		doc      *data.Document
		renderer Renderer
	}
)

func NewDocument() *Document {
	return &Document{
		renderer: NewRawRenderer(),
	}
}

func (d *Document) SetDocument(doc *data.Document) {
	d.doc = doc
}

func (d *Document) SetRenderer(r Renderer) {
	d.renderer = r
}

func (d *Document) Contents() ([]byte, bool) {
	return d.renderer.Render(d.doc.Contents())
}
