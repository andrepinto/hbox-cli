package utils

import (
	"errors"
	"regexp"
)

var (
	FuncSlugPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)
)

func ValidateFunctionSlug(slug string) error {
	if !FuncSlugPattern.MatchString(slug) {
		return errors.New("Invalid Function name. Must start with at least one letter, and only include alphanumeric characters, underscores, and hyphens. (^[A-Za-z][A-Za-z0-9_-]*$)")
	}

	return nil
}
