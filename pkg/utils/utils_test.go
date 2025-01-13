package utils_test

import (
	"testing"

	"github.com/dim2k2006/correlateapp-be/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Positive numbers", 2, 3, 5},
		{"Negative numbers", -2, -3, -5},
		{"Mixed numbers", -2, 3, 1},
		{"Zero", 0, 5, 5},
	}

	for _, tt := range tests {
		sum := utils.Add(tt.a, tt.b)
		assert.Equal(t, tt.expected, sum, "Add(%d, %d) should be %d", tt.a, tt.b, tt.expected)
	}
}
