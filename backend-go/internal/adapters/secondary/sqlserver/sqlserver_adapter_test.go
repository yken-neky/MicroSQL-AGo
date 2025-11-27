package sqlserver

import (
	"testing"
)

func TestConvertResultToBool_variousTypes(t *testing.T) {
	tests := []struct {
		in   interface{}
		want bool
		ok   bool
	}{
		{true, true, true},
		{int64(1), true, true},
		{int64(0), false, true},
		{int(2), true, true},
		{float64(0.0), false, true},
		{[]byte("TRUE"), true, true},
		{[]byte("FALSE"), false, true},
		{"TRUE", true, true},
		{"FALSE", false, true},
		{"1", true, true},
		{"0", false, true},
		{"maybe", false, false},
	}

	for _, tc := range tests {
		got, err := convertResultToBool(tc.in)
		if tc.ok {
			if err != nil {
				t.Fatalf("expected ok for %v but got err: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("expected %v -> %v, got %v", tc.in, tc.want, got)
			}
		} else {
			if err == nil {
				t.Fatalf("expected error for %v, got nil", tc.in)
			}
		}
	}
}
