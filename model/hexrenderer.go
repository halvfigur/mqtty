package model

import (
	"fmt"
	"strings"
)

func hex(b byte) []byte {
	table := []byte("0123456789abcdef")

	b0 := table[b>>4]
	b1 := table[b&0xf]

	return []byte{b0, b1}

}

func toPrintable(b byte) byte {
	if b >= ' ' && b <= '~' {
		return b
	}

	return '.'
}

type HexRenderer struct {
	sb *strings.Builder
}

func NewHexRenderer() *HexRenderer {
	return &HexRenderer{
		sb: new(strings.Builder),
	}
}

func (r *HexRenderer) Name() string {
	return "Hex"
}

func (r *HexRenderer) Render(data []byte) ([]byte, bool) {
	if data == nil {
		return nil, false
	}

	r.sb.Reset()

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
		[]byte("   "),
		[]byte(" "),
		[]byte(" "),
		[]byte(" "),
		[]byte("   "),
		[]byte(" "),
		[]byte(" "),
		[]byte(" "),
		[]byte("  "),
	}

	r.sb.Write([]byte("[blue]0000>[-] "))
	for i := 0; i < len(data); i += len(padding) {
		if i > 0 {
			r.sb.Write([]byte("\n"))
			r.sb.Write([]byte(fmt.Sprintf("[blue]%04x>[-] ", i)))
		}

		l := min(len(padding), len(data)-i)

		for j := 0; j < l; j++ {
			r.sb.Write(hex(data[i+j]))
			r.sb.Write(padding[j])
		}

		for j := 0; j < l; j++ {
			r.sb.WriteByte(toPrintable(data[i+j]))
			r.sb.Write(padding[j][0 : len(padding[j])-1])
		}
	}

	return []byte(r.sb.String()), true
}
