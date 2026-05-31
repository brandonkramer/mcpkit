package envelope

import "testing"

func TestNextActionFromHint(t *testing.T) {
	t.Parallel()
	action := NextActionFromHint(`runs.get id="01RUN"`)
	if action.Tool != "runs.get" {
		t.Fatalf("tool=%q", action.Tool)
	}
	if action.Args["id"] != "01RUN" {
		t.Fatalf("args=%v", action.Args)
	}
	if action.Hint == "" {
		t.Fatal("hint required")
	}
}

func TestHintsToActions(t *testing.T) {
	t.Parallel()
	actions := HintsToActions("tasks.list", "runs.poll id=01 since_seq=7")
	if len(actions) != 2 {
		t.Fatalf("actions=%+v", actions)
	}
	if actions[1].Args["since_seq"] != "7" {
		t.Fatalf("args=%v", actions[1].Args)
	}
}

func TestHintParam(t *testing.T) {
	t.Parallel()
	if got := HintParam("id", "01RUN"); got != `id="01RUN"` {
		t.Fatalf("HintParam=%q", got)
	}
}

func TestFilterSubjects(t *testing.T) {
	t.Parallel()
	subjects := FilterSubjects(
		Subject{Kind: "task", ID: "01"},
		Subject{Kind: "run", ID: ""},
		Subject{Kind: "", ID: "x"},
	)
	if len(subjects) != 1 || subjects[0].Kind != "task" {
		t.Fatalf("subjects=%+v", subjects)
	}
}
