package present

import "testing"

func TestRedactText(t *testing.T) {
	t.Parallel()
	prev, cleared := RedactText("hello", 10)
	if prev.Preview != "hello" || prev.Truncated || cleared != "hello" {
		t.Fatalf("preview=%+v cleared=%q", prev, cleared)
	}
	prev, cleared = RedactText("0123456789abcdef", 10)
	if prev.Preview != "0123456789" || !prev.Truncated || cleared != "" {
		t.Fatalf("preview=%+v cleared=%q", prev, cleared)
	}
}
