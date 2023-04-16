package sql

import (
	"runtime"
)

const maxDepth = 32

func trace(skip int) StackTrace {
	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+2, pcs[:])
	st := make([]Frame, 0, n)
	for i := 0; i < n; i++ {
		st = append(st, Frame(pcs[i]))
	}
	return st
}

type StackTrace []Frame

type Frame uintptr

func (f Frame) Func() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	return fn.Name()
}

func (f Frame) File() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return ""
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}
