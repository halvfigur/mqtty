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
		return []byte(fmt.Sprint("[red]Document is not valid JSON:[-] ", err.Error())), true
	}

	return b.Bytes(), false
}
