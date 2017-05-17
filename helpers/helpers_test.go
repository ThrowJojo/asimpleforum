package helpers

import "testing"

func TestIntInSlice(t *testing.T) {

	list := []int{1, 2, 3}

	if result := IntInSlice(list, 1); !result {
		t.Error("Expected result to be true")
	}

	if result := IntInSlice(list, 5); result {
		t.Error("Expected result to be false")
	}

}