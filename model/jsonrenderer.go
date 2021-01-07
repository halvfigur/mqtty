package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonRenderer struct{}

func NewJsonRenderer() *JsonRenderer {
	return new(JsonRenderer)
}

func (r *JsonRenderer) Name() string {
	return "JSON"
}

func (r *JsonRenderer) Render(data []byte) ([]byte, bool) {
	if data == nil {
		return nil, false
	}

	b := new(bytes.Buffer)

	if err := json.Indent(b, data, "    ", ""); err != nil {
		b := new(bytes.Buffer)
		b.WriteString(fmt.Sprintf("[red]Document is not valid JSON:[-] %s\n\n", err.Error()))
		b.Write(data)
		return b.Bytes(), true
	}

	return b.Bytes(), false
}
