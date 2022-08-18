package internal

import "testing"

func TestJoinNotEmpty(t *testing.T) {

	tests := []struct {
		name  string
		s     []string
		wantR string
	}{
		{
			name:  "single arg",
			s:     []string{"ABCD"},
			wantR: "ABCD",
		},
		{
			name:  "one empty arg",
			s:     []string{"ABCD", "", "1234"},
			wantR: "ABCD, 1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := JoinNotEmpty(tt.s...); gotR != tt.wantR {
				t.Errorf("JoinNotEmpty() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
