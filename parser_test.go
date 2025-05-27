package main

import (
	"testing"
)

func TestParseObisLine(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "(0*k.:4.8.0(016800.7-000.0:2.7.0(0000*kW)2(0",
			expected: `{"id":"016800.7-000.0:2.7.0","value":0,"unit":"kW"}`,
		},
		{
			input:    "1-.0(002.9*A)",
			expected: `{"id":"1-.0","value":2.9,"unit":"A"}`,
		},
		{
			input:    "1-1-0:2.8.0(017520.933*kWh)",
			expected: `{"id":"1-1-0:2.8.0","value":17520.933,"unit":"kWh"}`,
		},
		{
			input:    "0:21.7.1-0:22.7.0(000*kW)",
			expected: `{"id":"0:21.7.1-0:22.7.0","value":0,"unit":"kW"}`,
		},
		{
			input:    ".0(0*kW)",
			expected: `{"id":".0","value":0,"unit":"kW"}`,
		},
		{
			input:    "1--0:51.7.0(2.9*A)",
			expected: `{"id":"1--0:51.7.0","value":2.9,"unit":"A"}`,
		},
		{
			input:    "1-0:50:71.7.0(001*A)",
			expected: `{"id":"1-0:50:71.7.0","value":1,"unit":"A"}`,
		},
		{
			input:    "1-0:1.8.0(12345.67*kWh)",
			expected: `{"id":"1-0:1.8.0","value":12345.67,"unit":"kWh"}`,
		},
		{
			input:    "1-0:2.8.0(ABC*kWh)", // invalid number
			expected: ``,
		},
		{
			input:    "bad line with no match",
			expected: ``,
		},
	}

	for _, test := range tests {
		msg := ParseObisLine(test.input)
		if test.expected == "" {
			if msg != nil {
				t.Errorf("Expected nil, got %+v for input %q", msg, test.input)
			}
			continue
		}
		if msg == nil {
			t.Errorf("Expected non-nil result for input %q", test.input)
			continue
		}
		result := MustJSON(msg)
		if result != test.expected {
			t.Errorf("Mismatch:\nInput:    %q\nExpected: %s\nGot:      %s", test.input, test.expected, result)
		}
	}
}
