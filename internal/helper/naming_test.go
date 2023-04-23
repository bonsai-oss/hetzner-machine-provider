package helper

import (
	"testing"
)

func TestSetResourceNamePrefix(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		prefix     string
		shouldFail bool
	}{
		{"ValidPrefix", "test-job-", false},
		{"InvalidPrefixTooLong", "test-job-12345678901234567890123456789012345678901234", true},
		{"InvalidPrefixWithInvalidChars", "test$job-", true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := SetResourceNamePrefix(testCase.prefix)
			if testCase.shouldFail && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !testCase.shouldFail && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}

func TestResourceName(t *testing.T) {
	jobID := "123456"
	expected := "test-job-123456"
	resourceName := ResourceName(jobID)

	if resourceName != expected {
		t.Errorf("Expected: %s, got: %s", expected, resourceName)
	}
}

func TestValidateResourceName(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{"ValidResourceName", "test-job-123456", false},
		{"InvalidResourceNameTooLong", "test-job-123456789012345678901234567890123456789012345678901", true},
		{"InvalidResourceNameWithInvalidChars", "test-job-12@3456", true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateResourceName(testCase.input)
			if testCase.shouldFail && err == nil {
				t.Errorf("Expected error, but got nil")
			}
			if !testCase.shouldFail && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
		})
	}
}
