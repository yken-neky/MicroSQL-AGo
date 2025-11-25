package logging

import "strings"

// RedactAuthHeader masks bearer tokens and similar Authorization header values.
// Example: "Bearer eyJhbGci..." -> "Bearer <redacted>"
func RedactAuthHeader(h string) string {
    if h == "" {
        return h
    }
    // common pattern: "Bearer <token>"
    parts := strings.SplitN(h, " ", 2)
    if len(parts) == 2 {
        scheme := parts[0]
        token := parts[1]
        return scheme + " " + RedactToken(token)
    }
    // unknown shape -> mask most of the value
    return RedactToken(h)
}

// RedactToken returns a truncated token that keeps first/last 4 chars.
func RedactToken(token string) string {
    if token == "" {
        return ""
    }
    if len(token) <= 8 {
        return "<redacted>"
    }
    return token[:4] + "..." + token[len(token)-4:]
}

// RedactValue generic redaction helper for values like passwords or secrets.
func RedactValue(v string) string {
    if v == "" {
        return v
    }
    return "<redacted>"
}
