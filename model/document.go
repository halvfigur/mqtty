package model

import (
	"github.com/halvfigur/mqtty/data"
)

type (
	Document struct {
		doc *data.Document
	}
)

func NewDocument() *Document {
	return &Document{}
}

func (d *Document) SetDocument(doc *data.Document) {
	d.doc = doc
}

func (d *Document) Contents() []byte {
	return d.doc.Contents()
}
