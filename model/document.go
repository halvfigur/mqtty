package model

import (
	"github.com/halvfigur/mqtty/data"
)

type (
	Renderer interface {
		Name() string
		Render(doc *data.Document) string
	}

	RawRenderer struct{}

	Document struct {
		doc      *data.Document
		renderer Renderer
	}
)

func (r *RawRenderer) Name() string {
	return "Raw"
}

func (r *RawRenderer) Render(d *data.Document) string {
	return string(d.Contents())
}

func NewDocument() *Document {
	return &Document{
		renderer: new(RawRenderer),
	}
}

func (d *Document) SetDocument(doc *data.Document) {
	d.doc = doc
}

func (d *Document) SetRenderer(r Renderer) {
	d.renderer = r
}

func (d *Document) Contents() []byte {
	return []byte(d.renderer.Render(d.doc))
}
