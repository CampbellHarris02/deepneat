package genetics

import (
	"errors"
	"io"
)

const errAlwaysText = "always error"

var errAlways = errors.New(errAlwaysText)

type ErrorWriter int

func (e ErrorWriter) Write(_ []byte) (int, error) {
	return 0, errAlways
}

type ErrorReader int

func (e ErrorReader) Read(_ []byte) (int, error) {
	return 0, errAlways
}

func NewErrorReader(n int) io.Reader {
	return ErrorReader(n)
}

func NewErrorWriter(n int) io.Writer {
	return ErrorWriter(n)
}
