package logging

import "testing"

func TestRedactAuthHeader(t *testing.T) {
    in := "Bearer abcdef1234567890XYZ"
    got := RedactAuthHeader(in)
    if got == in {
        t.Fatalf("expected header to be redacted, got same value")
    }
    if len(got) == 0 {
        t.Fatalf("redacted value empty")
    }
}

func TestRedactTokenShort(t *testing.T) {
    if RedactToken("abcd") != "<redacted>" {
        t.Fatalf("expected short token to be fully redacted")
    }
}
