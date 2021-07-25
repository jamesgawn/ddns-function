package main

import (
	"os"
	"strconv"
	"testing"
)

func TestObtainVersion(t *testing.T) {
	result := obtainVersion()
	EqualsString(t, "0.0.0", result)
}

func TestAuthenticate(t *testing.T) {
	userErr := os.Setenv("username", "testuser")
	if userErr != nil {
		t.Error(userErr)
	}
	passErr := os.Setenv("password", "testing")
	if passErr != nil {
		t.Error(passErr)
	}
	result := authenticate("Basic dGVzdHVzZXI6dGVzdGluZw==")
	EqualsBool(t, true, result)
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

func EqualsBool(t *testing.T, expected bool, actual bool) {
	if expected != actual {
		t.Errorf(
			"Expected %s, but got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual),
		)
	}
}
