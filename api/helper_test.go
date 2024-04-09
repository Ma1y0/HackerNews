package api

import (
	"fmt"
	"testing"
)

func TestEqIntSclice(t *testing.T) {
	testCases := []struct {
		a        []int
		b        []int
		expected bool
	}{
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, true},
		{[]int{8, 6, 56, 885}, []int{5}, false},
		{[]int{4, 5, 6}, []int{3, 6, 4656}, false},
	}

	for i, tt := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			if a := eqIntSclice(tt.a, tt.b); a != tt.expected {
				t.Errorf("Got=%v, Expected=%v", a, tt.expected)
			}
		})
	}
}
