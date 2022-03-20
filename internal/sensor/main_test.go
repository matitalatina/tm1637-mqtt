package sensor

import (
	"reflect"
	"testing"
)

func TestSplitToDigits(t *testing.T) {
	assertDigits([]int{2, 3, 4, 5}, 2345, t)
	assertDigits([]int{1, 0, 0, 0}, 1000, t)
	assertDigits([]int{0, 0, 0, 0}, 0, t)
	assertDigits([]int{0, 0, 0, 0}, -100, t)
	assertDigits([]int{9, 9, 9, 9}, 10000, t)
}

func assertDigits(expected []int, number int, t *testing.T) {
	if result := splitToDigits(number); !reflect.DeepEqual(result, expected) {
		t.Errorf("SplitToDigits should be %d but got %d", expected, result)
	}
}
