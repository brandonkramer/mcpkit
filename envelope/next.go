package envelope

import (
	"fmt"
	"strings"
)

// Hint builds a string continuation hint such as `tool arg=value`.
func Hint(tool string, parts ...string) string {
	if tool == "" {
		return strings.Join(parts, " ")
	}
	if len(parts) == 0 {
		return tool
	}
	return tool + " " + strings.Join(parts, " ")
}

// HintsToActions converts string hints into structured next actions.
func HintsToActions(hints ...string) []NextAction {
	if len(hints) == 0 {
		return nil
	}
	out := make([]NextAction, 0, len(hints))
	for _, hint := range hints {
		if hint == "" {
			continue
		}
		out = append(out, NextActionFromHint(hint))
	}
	return out
}

// NextActionFromHint parses a simple `tool key=value` hint into a NextAction.
func NextActionFromHint(hint string) NextAction {
	hint = strings.TrimSpace(hint)
	if hint == "" {
		return NextAction{}
	}
	fields := strings.Fields(hint)
	action := NextAction{Hint: hint, Tool: fields[0]}
	if len(fields) == 1 {
		return action
	}
	args := make(map[string]string)
	for _, field := range fields[1:] {
		key, val, ok := strings.Cut(field, "=")
		if !ok {
			continue
		}
		args[key] = strings.Trim(val, "\"")
	}
	if len(args) > 0 {
		action.Args = args
	}
	return action
}

// HintParam formats a key=value hint argument using Go quoting rules.
func HintParam(key, value string) string {
	if key == "" {
		return value
	}
	return fmt.Sprintf("%s=%q", key, value)
}

// FilterSubjects returns subjects with non-empty kind and id.
func FilterSubjects(subjects ...Subject) []Subject {
	if len(subjects) == 0 {
		return nil
	}
	out := make([]Subject, 0, len(subjects))
	for _, subject := range subjects {
		if subject.Kind == "" || subject.ID == "" {
			continue
		}
		out = append(out, subject)
	}
	return out
}
