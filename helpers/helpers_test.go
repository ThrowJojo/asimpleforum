package helpers

import "testing"

func TestIntInSlice(t *testing.T) {

	list := []int{1, 2, 3}

	if result := IntInSlice(1, list); !result {
		t.Error("Expected result to be true")
	}

	if result := IntInSlice(5, list); result {
		t.Error("Expected result to be false")
	}

}