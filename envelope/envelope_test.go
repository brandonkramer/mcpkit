package envelope

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSuccessEnvelopeJSON(t *testing.T) {
	t.Parallel()
	env := Success("search").
		Summary("3 matches").
		AnswerContext("Found three files mentioning auth.").
		Subject("file", "src/auth.go").
		Data(map[string]any{"count": 3}).
		NextHints("search query=session", "read path=src/auth.go").
		Quality(QualitySignal{Name: "coverage", Score: 0.9}).
		Diagnostic(Diagnostic{Code: "truncated", Message: "preview capped", Severity: "info"}).
		Build()

	data, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, sub := range []string{
		`"ok":true`, `"tool":"search"`, `"summary":"3 matches"`,
		`"answer_context":"Found three files mentioning auth."`,
		`"kind":"file"`, `"id":"src/auth.go"`,
		`"next_actions"`, `"quality_signals"`, `"diagnostics"`,
	} {
		if !strings.Contains(text, sub) {
			t.Fatalf("json=%s want substring %q", text, sub)
		}
	}
}

func TestErrorEnvelopeSeparatesSummaryAndError(t *testing.T) {
	t.Parallel()
	env := Error("request failed", errTest("boom"))
	if env.OK || env.Summary != "request failed" || env.Error != "boom" {
		t.Fatalf("env=%+v", env)
	}
	if strings.Contains(env.AnswerContext, "boom") {
		t.Fatal("error text must not leak into answer_context")
	}
}

func TestSubjectsReplaceDomainIDs(t *testing.T) {
	t.Parallel()
	env := Success("notify").
		Subject("task", "01TASK").
		Subject("run", "01RUN").
		Subject("notice", "01NOTICE").
		Build()
	if len(env.Subjects) != 3 {
		t.Fatalf("subjects=%+v", env.Subjects)
	}
	data, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "task_id") || strings.Contains(string(data), "run_id") {
		t.Fatalf("json must not contain domain-specific id fields: %s", data)
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }
