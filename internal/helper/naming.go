package helper

import (
	"fmt"
	"regexp"
)

var (
	resourceNamePrefix    = "hmp-job-"
	resourceNameValidator = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
)

func SetResourceNamePrefix(prefix string) error {
	// do validation inclusive of a dummy job id
	if err := validateResourceName(prefix + "123456"); err != nil {
		return err
	}
	resourceNamePrefix = prefix
	return nil
}

func ResourceName(jobID string) string {
	return resourceNamePrefix + jobID
}

// validateResourceName validates a resource name
func validateResourceName(name string) error {
	if len(name) > 50 {
		return fmt.Errorf("resource name is too long: %s", name)
	}
	if !resourceNameValidator.MatchString(name) {
		return fmt.Errorf("resource name is invalid: %s", name)
	}
	return nil
}
