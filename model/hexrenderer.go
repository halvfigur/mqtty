package model

import (
	"strings"

	"github.com/halvfigur/mqtty/data"
)

func hex(b byte) []byte {
	table := []byte("0123456789abcdef")

	b0 := table[b>>4]
	b1 := table[b&0xf]

	return []byte{b0, b1}

}

type HexRenderer struct {
}

func (r *HexRenderer) Name() string {
	return "Hex"
}

func (r *HexRenderer) Render(doc *data.Document) string {
	c := doc.Contents()
	if c == nil {
		return ""
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}

		return b
	}

	padding := [][]byte{
		[]byte(" "),
		[]byte(" "),
		[]byte(" "),
		[]byte("   "),
		[]byte(" "),
		[]byte(" "),
		[]byte(" "),
		[]byte(""),
	}

	sb := new(strings.Builder)
	for i := 0; i < len(c); i += 8 {
		if i > 0 {
			sb.WriteString("\n")
		}

		l := min(8, len(c)-i)

		for j := 0; j < l; j++ {
			sb.Write(hex(c[i+j]))
			sb.Write(padding[j])
		}
	}

	return sb.String()
}
