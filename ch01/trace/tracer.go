package trace

import (
	"fmt"
	"io"
)

// Tracer는 코드 전체에서 이벤트를 추적할 수 있는
// 깩체를 설명하는 인터페이스다.

type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer)Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}
//Off는 Trace에 대한 호출을 무시할 Tracer를 생성한다.
func Off() Tracer {
	return &nilTracer{}
}