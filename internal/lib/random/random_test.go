package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size 1",
			size: 1,
		},
		{
			name: "size 6",
			size: 6,
		},
		{
			name: "size 8",
			size: 8,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 25",
			size: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			time.Sleep(1 * time.Nanosecond)
			str2 := NewRandomString(tt.size)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			assert.NotEqual(t, str1, str2)
		})
	}
}
