package models

import "fmt"

type DiagnosticLevel uint8

const (
	ERROR DiagnosticLevel = iota
	WARNING
	INFO
)

type Diagnostic struct {
	Span  Span
	Msg   string
	Level DiagnosticLevel
}

type Diagnostics struct {
	Src  []rune
	list []Diagnostic
}

func (d *Diagnostics) Errorf(s Span, level DiagnosticLevel, f string, a ...any) {
	d.list = append(d.list, Diagnostic{
		Span: s,
		Msg:  fmt.Sprintf(f, a...),
		Level: level,
	})
}
