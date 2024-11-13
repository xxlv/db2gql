package main

import "testing"

func TestAsCamStyle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "aA",
			expected: "Aa",
		},
		{
			name:     "multiple words",
			input:    "ax_aa",
			expected: "AxAa",
		},
		{
			name:     "a_b",
			input:    "a_b",
			expected: "Ab",
		},
		{
			name:     "aBc",
			input:    "aBc",
			expected: "Abc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := asCamStyle(tc.input)
			if actual != tc.expected {
				t.Errorf("asCamStyle(%q) = %q, expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}
