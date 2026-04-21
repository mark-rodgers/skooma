// Package validators contains validation functions for user input and configuration.
package validators

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"

	"github.com/skooma-cli/skooma/internal/sanitize"
	"github.com/skooma-cli/skooma/internal/types"
)

var ValidatorRegistry = map[string]func(string) types.ValidatorFunc{
	"not_empty":       NotEmpty,
	"no_spaces":       NoSpaces,
	"no_underscores":  NoUnderscores,
	"valid_url":       ValidURL,
	"rfc5322_address": RFC5322Address,
}

// NotEmpty checks that a string is not empty or whitespace-only.
func NotEmpty(label string) types.ValidatorFunc {
	return func(str string) error {
		if strings.TrimSpace(str) == "" {
			return errors.New(strings.ToLower(label) + " can't be empty")
		}
		return nil
	}
}

// NoSpaces checks that a string contains no spaces.
func NoSpaces(label string) types.ValidatorFunc {
	return func(str string) error {
		if strings.Contains(str, " ") {
			return errors.New(strings.ToLower(label) + " can't contain spaces")
		}
		return nil
	}
}

// NoUnderscores checks that a string contains no underscores.
func NoUnderscores(label string) types.ValidatorFunc {
	return func(str string) error {
		if strings.Contains(str, "_") {
			return errors.New(strings.ToLower(label) + " can't contain underscores")
		}
		return nil
	}
}

// ValidURL checks that a string is a valid URL when prefixed with https://.
// Handles input that may or may not include an http/https prefix.
func ValidURL(label string) types.ValidatorFunc {
	return func(str string) error {
		cleaned := sanitize.StripHTTPPrefix(str)
		u, err := url.ParseRequestURI("https://" + cleaned)
		if err != nil || u.Host == "" || !strings.Contains(u.Host, ".") {
			return errors.New(strings.ToLower(label) + " must be a valid URL (e.g., github.com/user/repo)")
		}
		parts := strings.SplitN(u.Host, ".", 2)
		if parts[0] == "" || parts[1] == "" {
			return errors.New(strings.ToLower(label) + " must be a valid URL (e.g., github.com/user/repo)")
		}
		return nil
	}
}

// RFC5322Address validates the "Name <email>" format using net/mail.
func RFC5322Address(label string) types.ValidatorFunc {
	return func(str string) error {
		addr, err := mail.ParseAddress(str)
		if err != nil || addr.Name == "" {
			return errors.New(strings.ToLower(label) + " must be in format: Name <email@example.com>")
		}
		return nil
	}
}

// ResolveValidators builds a composed ValidatorFunc from the validator names
// defined on a TemplateConfigVariable, threading the variable's prompt as the label.
func ResolveValidators(v types.TemplateConfigVariable) (types.ValidatorFunc, error) {
	var fns []types.ValidatorFunc
	for _, name := range v.Validators {
		factory, ok := ValidatorRegistry[name]
		if !ok {
			return nil, fmt.Errorf("unknown validator %q for variable %q", name, v.Name)
		}
		fns = append(fns, factory(v.Prompt))
	}
	if !v.Required {
		return AllowEmpty(fns...), nil
	}
	return All(fns...), nil
}

// AllowEmpty wraps one or more validators to allow empty strings. If the input
// is empty or whitespace-only, it skips validation and returns nil. Otherwise,
// it runs the validators in order.
func AllowEmpty(validators ...types.ValidatorFunc) types.ValidatorFunc {
	return func(str string) error {
		if strings.TrimSpace(str) == "" {
			return nil
		}
		return All(validators...)(str)
	}
}

// All composes multiple validators into one, running them in order and
// stopping at the first error.
func All(validators ...types.ValidatorFunc) types.ValidatorFunc {
	return func(str string) error {
		for _, v := range validators {
			if err := v(str); err != nil {
				return err
			}
		}
		return nil
	}
}
