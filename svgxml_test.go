package svgxml

import (
	"fmt"
	s "strings"
	"testing"
)

func TestAnyNumberToString(t *testing.T) {
	stringIntVal := "9095454"
	stringFloatVal := stringIntVal + "." + stringIntVal
	int32Val := int32(9095454)
	int64Val := int64(9095454)
	float32Val := float32(2000.9095454)
	float64Val := float64(9095454.9095454)
	badVal := "non-numeric"

	tests := []struct {
		precision       int
		expected        string
		expectedFloat32 string
	}{
		{precision: 0, expected: "9095455", expectedFloat32: "2001"},
		{precision: 1, expected: "9095454.9", expectedFloat32: "2000.9"},
		{precision: 2, expected: "9095454.91", expectedFloat32: "2000.91"},
		{precision: 3, expected: "9095454.910", expectedFloat32: "2000.910"},
		{precision: 4, expected: "9095454.9095", expectedFloat32: "2000.9095"},
		{precision: 5, expected: "9095454.90955", expectedFloat32: "2000.90955"},
		{precision: 6, expected: "9095454.909545", expectedFloat32: "2000.909546"},
		{precision: 7, expected: "9095454.9095454", expectedFloat32: "2000.9095459"},
		{precision: 8, expected: "9095454.90954540", expectedFloat32: "2000.90954590"},
		{precision: 9, expected: "9095454.909545399", expectedFloat32: "2000.909545898"},
	}

	for i, tt := range tests {
		if i == 2 {
			t.Run(fmt.Sprintf("negative_string_float_.%d", tt.precision), func(t *testing.T) {
				result := anyNumberToString("-"+stringFloatVal, tt.precision)
				expected := "-" + tt.expected
				if result != expected {
					t.Errorf("got '%s', want '%s'", result, "-"+expected)
				}
			})
			t.Run(fmt.Sprintf("invalid_string_.%d", tt.precision), func(t *testing.T) {
				result := anyNumberToString(badVal, tt.precision)
				if result != "" {
					t.Errorf("got '%s', want ''", result)
				}
			})
		}
		t.Run(fmt.Sprintf("string_float_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(stringFloatVal, tt.precision)
			if result != tt.expected {
				t.Errorf("got '%s', want '%s'", result, tt.expected)
			}
		})
		t.Run(fmt.Sprintf("float32_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(float32Val, tt.precision)
			if result != tt.expectedFloat32 {
				t.Errorf("got '%s', want '%s'", result, tt.expectedFloat32)
			}
		})
		t.Run(fmt.Sprintf("float64_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(float64Val, tt.precision)
			if result != tt.expected {
				t.Errorf("got '%s', want '%s'", result, tt.expected)
			}
		})

		// only integer-like inputs below this comment
		if dotIndex := s.Index(tt.expected, "."); dotIndex >= 0 {
			tt.expected = tt.expected[:dotIndex+1] + s.Repeat("0", len(tt.expected[dotIndex+1:]))
		} else {
			tt.expected = stringIntVal
		}
		t.Run(fmt.Sprintf("string_int_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(stringIntVal, tt.precision)
			if result != tt.expected {
				t.Errorf("got '%s', want '%s'", result, tt.expected)
			}
		})
		t.Run(fmt.Sprintf("int32_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(int32Val, tt.precision)
			if result != tt.expected {
				t.Errorf("got '%s', want '%s'", result, tt.expected)
			}
		})
		t.Run(fmt.Sprintf("int64_.%d", tt.precision), func(t *testing.T) {
			result := anyNumberToString(int64Val, tt.precision)
			if result != tt.expected {
				t.Errorf("got '%s', want '%s'", result, tt.expected)
			}
		})
	}
}

func TestGetViewBox(t *testing.T) {
	var tests = []struct {
		name     string
		svgObj   *SVG
		expected ViewBoxDef
	}{
		{"ints, float height", &SVG{ViewBox: "0 0 100 99.5"}, ViewBoxDef{0.0, 0.0, 100.0, 99.5}},
		{"ints, float width", &SVG{ViewBox: "0 0 99.5 100"}, ViewBoxDef{0.0, 0.0, 99.5, 100.0}},
		{"ints, float width, negative X", &SVG{ViewBox: "-10 0 89.5 100"}, ViewBoxDef{-10.0, 0.0, 89.5, 100.0}},
		{"ints, negative float Y", &SVG{ViewBox: "0 -10.5 99.5 89.5"}, ViewBoxDef{0.0, -10.5, 99.5, 89.5}},
		{"error: bad Y", &SVG{ViewBox: "0 -10>5 99.5 89.5"}, ViewBoxDef{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.svgObj.GetViewBox()
			if tt.name[:6] == "error:" {
				if err == nil {
					t.Errorf("expected error but got none; result=%+v", result)
				}
			} else {
				if err != nil {
					t.Errorf("got error: %s", err.Error())
				} else if result != tt.expected {
					t.Errorf("got +%v, want +%v", result, tt.expected)
				}
			}
		})
	}
}
