package utils_test

import (
	"testing"

	"github.com/dim2k2006/correlateapp-be/pkg/utils"
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
		t.Run(tt.name, func(t *testing.T) {
			sum := utils.Add(tt.a, tt.b)
			if sum != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, sum, tt.expected)
			}
		})
	}
}
