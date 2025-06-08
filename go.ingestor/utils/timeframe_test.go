package utils

import (
	"testing"

	"github.com/polygon-io/client-go/rest/models"
)

func TestParseTimeframe_ValidCases(t *testing.T) {
	tests := []struct {
		input    string
		expected PolygonTimeframe
	}{
		{"1m", PolygonTimeframe{1, models.Minute}},
		{"5h", PolygonTimeframe{5, models.Hour}},
		{"1d", PolygonTimeframe{1, models.Day}},
		{"2w", PolygonTimeframe{2, models.Week}},
		{"3M", PolygonTimeframe{3, models.Month}},
		{"1y", PolygonTimeframe{1, models.Year}},
	}

	for _, test := range tests {
		got, err := ParseTimeframe(test.input)
		if err != nil {
			t.Errorf("expected no error for input %s, got %v", test.input, err)
		}
		if got != test.expected {
			t.Errorf("for input %s, expected %+v, got %+v", test.input, test.expected, got)
		}
	}
}

func TestParseTimeframe_InvalidCases(t *testing.T) {
	invalidInputs := []string{
		"", "abc", "1", "5z", "m1", "60", "1minute", "10x", "1 h",
	}

	for _, input := range invalidInputs {
		_, err := ParseTimeframe(input)
		if err == nil {
			t.Errorf("expected error for invalid input %s, got none", input)
		}
	}
}
