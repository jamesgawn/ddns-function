package ddnsfunction

import (
	"testing"
)

func TestObtainVersion(t *testing.T) {
	result := ObtainVersion()
	EqualsString(t, "0.0.0", result)
}

func EqualsString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf(
			"Expected %s, but got %s",
			expected,
			actual,
		)
	}
}
