package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewEmail_ValidEmail_ReturnsNoError(t *testing.T) {
	email, err := NewEmail("test@example.com")
	assert.NoError(t, err, "valid email should not return an error")
	assert.Equal(t, Email("test@example.com"), email, "email should be correctly parsed")
}

func Test_NewEmail_InvalidEmail_ReturnsError(t *testing.T) {
	tc := []struct {
		name  string
		input string
	}{
		{name: "empty email", input: ""},
		{name: "missing domain", input: "test@"},
		{name: "missing username", input: "@example.com"},
		{name: "missing @ symbol", input: "testexample.com"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewEmail(tt.input)
			assert.Error(t, err, "invalid email should return error")
		})
	}
}
