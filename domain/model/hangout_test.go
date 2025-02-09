package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMinute_ValidDuration_ReturnsNoError(t *testing.T) {
	tc := []struct {
		name  string
		input int
		want  Minutes
	}{
		{name: "zero minutes", input: 0, want: Minutes(0)},
		{name: "one minute", input: 1, want: Minutes(1)},
		{name: "large number", input: 10000, want: Minutes(10000)},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			min, err := NewMinute(tt.input)
			assert.NoError(t, err, "valid duration should not return an error")
			assert.Equal(t, tt.want, min, "should return correct Minutes value")
		})
	}
}

func Test_NewMinute_NegativeDuration_ReturnsError(t *testing.T) {
	tc := []struct {
		name  string
		input int
	}{
		{name: "negative one", input: -1},
		{name: "large negative", input: -10000},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMinute(tt.input)
			assert.Error(t, err, "negative duration should return error")
		})
	}
}
