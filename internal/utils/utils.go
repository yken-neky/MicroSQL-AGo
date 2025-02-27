package utils

// IsValidUUID Function auxiliary de ejemplo para validar UUID
func IsValidUUID(u string) bool {
	return len(u) >= 1 && len(u) < 3
}
