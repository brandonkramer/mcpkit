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
	got, err := NormalizeID(args{ID: "01"}, norm, func(a *args) *string { return &a.ID })
	if err != nil || got.ID != "N:01" {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	_, err = NormalizeID(args{ID: ""}, norm, func(a *args) *string { return &a.ID })
	if err == nil {
		t.Fatal("expected error")
	}
}
