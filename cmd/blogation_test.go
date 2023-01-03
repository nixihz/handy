package cmd

import "testing"

func Test_transMermaid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"hehe", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transMermaid(tt.input)
		})
	}
}
