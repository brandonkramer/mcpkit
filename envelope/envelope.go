package envelope

//
// ────────────────────────────────────────
// result types.
//

// Envelope is the agent-facing tool result payload.
type Envelope struct {
	OK            bool            `json:"ok"`
	Tool          string          `json:"tool,omitempty"`
	Method        string          `json:"method,omitempty"`
	Summary       string          `json:"summary,omitempty"`
	AnswerContext string          `json:"answer_context,omitempty"`
	Subjects      []Subject       `json:"subjects,omitempty"`
	Data          any             `json:"data,omitempty"`
	NextActions   []NextAction    `json:"next_actions,omitempty"`
	Next          []string        `json:"next,omitempty"`
	Quality       []QualitySignal `json:"quality_signals,omitempty"`
	Diagnostics   []Diagnostic    `json:"diagnostics,omitempty"`
	Error         string          `json:"error,omitempty"`
}

// Subject identifies a domain entity referenced by the tool result.
type Subject struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Role string `json:"role,omitempty"`
}

// NextAction describes a capability or continuation the caller can invoke next.
type NextAction struct {
	Tool string            `json:"tool,omitempty"`
	Hint string            `json:"hint,omitempty"`
	Args map[string]string `json:"args,omitempty"`
}

// QualitySignal carries optional confidence or completeness metadata.
type QualitySignal struct {
	Name  string  `json:"name"`
	Value any     `json:"value,omitempty"`
	Score float64 `json:"score,omitempty"`
	Note  string  `json:"note,omitempty"`
}

// Diagnostic carries implementation detail kept separate from answer content.
type Diagnostic struct {
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`
	Detail   any    `json:"detail,omitempty"`
}
