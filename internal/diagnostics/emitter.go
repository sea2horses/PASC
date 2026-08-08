package diagnostics

import (
	"fmt"
	"pseint-compiled/internal/models"
)

type DiagnosticEmitter struct {
	diagnostics []Diagnostic
}

func (d *DiagnosticEmitter) diagnosis(span models.Span, level DiagnosticLevel, format string, args ...any) {
	Dbg("Filing diagnosis of type: ", level, ". span: ", span)
	d.diagnostics = append(d.diagnostics, Diagnostic{
		Span:  span,
		Msg:   fmt.Sprintf(format, args...),
		Level: level,
	})
}

func NewDiagnosticEmitter() DiagnosticEmitter {
	/* Create new Diagnostic Emitter */
	return DiagnosticEmitter{
		diagnostics: []Diagnostic{},
	}
}

func (d *DiagnosticEmitter) Diagnostics() []Diagnostic {
	return d.diagnostics
}

func (d *DiagnosticEmitter) Report(span models.Span, format string, args ...any) {
	d.diagnosis(span, ERROR, format, args...)
}

func (d *DiagnosticEmitter) Warn(span models.Span, format string, args ...any) {
	d.diagnosis(span, WARNING, format, args...)
}

func (d *DiagnosticEmitter) Info(span models.Span, format string, args ...any) {
	d.diagnosis(span, INFO, format, args...)
}
