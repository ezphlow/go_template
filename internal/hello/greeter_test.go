package hello

import "testing"

func TestGreet(t *testing.T) {
	var s Service
	if got := s.Greet("").Message; got != "hello, world" {
		t.Errorf("empty name = %q, want hello, world", got)
	}
	if got := s.Greet("ez").Message; got != "hello, ez" {
		t.Errorf("named = %q, want hello, ez", got)
	}
}
