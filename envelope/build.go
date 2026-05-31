package envelope

//
// ────────────────────────────────────────
// envelope builder.
//

import "fmt"

const defaultFailureSummary = "request failed"

// Builder constructs an Envelope fluently.
type Builder struct {
	env Envelope
}

// Success starts a successful envelope for tool or method name.
func Success(name string) *Builder {
	return &Builder{env: Envelope{OK: true, Tool: name, Method: name}}
}

// Error builds a failure envelope.
func Error(summary string, err error) Envelope {
	env := Envelope{OK: false, Summary: summary}
	if summary == "" {
		env.Summary = defaultFailureSummary
	}
	if err != nil {
		env.Error = err.Error()
		if env.Summary == "" {
			env.Summary = env.Error
		}
	} else if env.Summary != "" {
		env.Error = env.Summary
	}
	return env
}

// Errorf builds a failure envelope from a formatted message.
func Errorf(format string, args ...any) Envelope {
	msg := fmt.Sprintf(format, args...)
	return Envelope{OK: false, Summary: defaultFailureSummary, Error: msg}
}

// Tool sets the tool name.
func (b *Builder) Tool(name string) *Builder {
	b.env.Tool = name
	return b
}

// Method sets the method name.
func (b *Builder) Method(name string) *Builder {
	b.env.Method = name
	return b
}

// Summary sets the short user-meaningful summary.
func (b *Builder) Summary(summary string) *Builder {
	b.env.Summary = summary
	return b
}

// AnswerContext sets concise LLM-facing synthesis context.
func (b *Builder) AnswerContext(context string) *Builder {
	b.env.AnswerContext = context
	return b
}

// Subject appends one subject reference.
func (b *Builder) Subject(kind, id string) *Builder {
	if kind == "" || id == "" {
		return b
	}
	b.env.Subjects = append(b.env.Subjects, Subject{Kind: kind, ID: id})
	return b
}

// SubjectRole appends one subject with a role.
func (b *Builder) SubjectRole(kind, id, role string) *Builder {
	if kind == "" || id == "" {
		return b
	}
	b.env.Subjects = append(b.env.Subjects, Subject{Kind: kind, ID: id, Role: role})
	return b
}

// Subjects replaces all subjects.
func (b *Builder) Subjects(subjects ...Subject) *Builder {
	b.env.Subjects = append([]Subject(nil), subjects...)
	return b
}

// Data sets structured payload data.
func (b *Builder) Data(data any) *Builder {
	b.env.Data = data
	return b
}

// NextActions sets structured continuations.
func (b *Builder) NextActions(actions ...NextAction) *Builder {
	b.env.NextActions = append([]NextAction(nil), actions...)
	return b
}

// NextHints sets string continuation hints and derives structured next actions.
func (b *Builder) NextHints(hints ...string) *Builder {
	b.env.Next = append([]string(nil), hints...)
	b.env.NextActions = HintsToActions(hints...)
	return b
}

// Quality appends quality signals.
func (b *Builder) Quality(signals ...QualitySignal) *Builder {
	b.env.Quality = append(b.env.Quality, signals...)
	return b
}

// Diagnostic appends diagnostics.
func (b *Builder) Diagnostic(items ...Diagnostic) *Builder {
	b.env.Diagnostics = append(b.env.Diagnostics, items...)
	return b
}

// Build returns the constructed envelope.
func (b *Builder) Build() Envelope {
	return b.env
}
