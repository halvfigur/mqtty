package model

import (
	"bytes"
	"io"
	"testing"

	"github.com/halvfigur/mqtty/data"
)

func TestRawRenderer(t *testing.T) {
	tests := []struct {
		name   string
		input  io.Reader
		expect string
	}{
		{
			"simple",
			bytes.NewBufferString("Hello world"),
			"48 65 6c 6c   6f 20 77 6f\n72 6c 64 ",
		},
	}

	for _, test := range tests {
		interp := new(HexRenderer)

		t.Run(test.name, func(t *testing.T) {
			d, err := data.NewDocument(test.input)
			if err != nil {
				t.Fatal("Failed to read from bytes.Buffer")
			}

			have := interp.Render(d)
			if have != test.expect {
				t.Errorf("Render() = [%s]; want [%s]", have, test.expect)
			}
		})
	}

}
