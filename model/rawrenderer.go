package model

type (
	Renderer interface {
		Name() string
		Render(data []byte) ([]byte, bool)
	}

	RawRenderer struct{}
)

func NewRawRenderer() *RawRenderer {
	return new(RawRenderer)
}

func (r *RawRenderer) Name() string {
	return "Raw"
}

func (r *RawRenderer) Render(data []byte) ([]byte, bool) {
	printable := make([]byte, len(data))

	for i, b := range data {
		printable[i] = toPrintable(b)
	}

	return printable, false
}
