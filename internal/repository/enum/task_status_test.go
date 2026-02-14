package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		status   TaskStatus
		expected string
	}{
		{Created, "Created"},
		{Started, "Started"},
		{Done, "Done"},
		{Failed, "Failed"},
		{Delayed, "Delayed"},
		{Canceled, "Canceled"},
		{TaskStatus(99), ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}
