package utils

import (
	"strings"
	"testing"
)

// AssertEqualsString asserts when the expected and actual strings are qual
func AssertEqualsString(t *testing.T, expected, actual string) {
	if strings.TrimSpace(expected) != strings.TrimSpace(actual) {
		t.Helper()
		t.Fatalf("Expected %s, got %s", expected, actual)
	}
}

// AssertPrefixString asserts when the actual string contains the prefix expectedPrefix
func AssertPrefixString(t *testing.T, expectedPrefix, actual string) {
	if !strings.HasPrefix(actual, expectedPrefix) {
		t.Helper()
		t.Fatalf("Expected %s, got %s", expectedPrefix, actual)
	}
}

// AssertContainsString asserts when the actual string contains the expectedSubString
func AssertContainsString(t *testing.T, expectedSubstring, actual string) {
	if !strings.Contains(actual, expectedSubstring) {
		t.Helper()
		t.Fatalf("Expected %s, got %s", expectedSubstring, actual)
	}
}

// AssertEqualsInt asserts when expected and actual integer valus are equal
func AssertEqualsInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Helper()
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertNilError asserts when the error is nil
func AssertNilError(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("Expected no error, but got %s", err.Error())
	}
}

// AssertNotNilError asserts when the error is not nil
func AssertNotNilError(t *testing.T, err error) {
	if err == nil {
		t.Helper()
		t.Fatalf("Expected error, but it was nil.")
	}
}

// AssertStringNotEmpty asserts when the string is not empty
func AssertStringNotEmpty(t *testing.T, message, str string) {
	str = strings.TrimSpace(str)
	if str != "" {
		return
	}

	if message != "" {
		t.Helper()
		t.Fatalf("%s: expected not empty string.", message)
	} else {
		t.Helper()
		t.Fatalf("Expected not empty string.")
	}
}

// AssertNotNil asserts when the objectType is not nil
func AssertNotNil(t *testing.T, obj interface{}) {
	if obj == nil {
		t.Helper()
		t.Fatalf("Expected object %v to be not null", obj)
	}
}

// AssertTrue asserts when the message is true
func AssertTrue(t *testing.T, message string, expression bool) {
	if !expression {
		t.Helper()
		t.Fatalf("%s: expected to be true, but it is false.", message)
	}
}

// AssertNil asserts when the objectType is nil
func AssertNil(t *testing.T, obj interface{}) {
	if obj != nil {
		t.Helper()
		t.Fatalf("Expected object %v to be null", obj)
	}
}
