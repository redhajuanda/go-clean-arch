package logger

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

func BenchmarkLogStack(b *testing.B) {
	out := &bytes.Buffer{}
	log := New("", "")
	SetOutput(out)
	err := errors.Wrap(errors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.WithStack(err).Error(err)
		out.Reset()
	}
}
