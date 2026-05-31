package tool

import (
	"errors"
	"testing"
)

func TestNormalizeID(t *testing.T) {
	t.Parallel()

	type args struct{ ID string }
	norm := func(s string) (string, error) {
		if s == "" {
			return "", errors.New("empty")
		}
		return "N:" + s, nil
	}

	cases := []struct {
		name    string
		id      string
		want    string
		wantErr bool
	}{
		{name: "normalizes id", id: "01", want: "N:01"},
		{name: "rejects empty", id: "", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeID(args{ID: tc.id}, norm, func(a *args) *string { return &a.ID })
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil || got.ID != tc.want {
				t.Fatalf("got=%+v err=%v", got, err)
			}
		})
	}
}
