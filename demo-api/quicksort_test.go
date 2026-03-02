package main

import (
	"reflect"
	"testing"
)

func TestSortArrayReturnsSortedList(t *testing.T) {
	input := []int{5, 3, 8, 1, 2, 7}
	got := sortArray(input)
	want := []int{1, 2, 3, 5, 7, 8}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected sorted array %v, got %v", want, got)
	}
}

func TestSortArrayDoesNotMutateInput(t *testing.T) {
	input := []int{4, 1, 4, 2}
	original := append([]int(nil), input...)

	_ = sortArray(input)

	if !reflect.DeepEqual(input, original) {
		t.Fatalf("sortArray mutated input: got %v, want %v", input, original)
	}
}

func TestSortArrayHandlesEmptyAndSingleItem(t *testing.T) {
	empty := []int{}
	if got := sortArray(empty); len(got) != 0 {
		t.Fatalf("expected empty sorted list, got %v", got)
	}

	single := []int{42}
	if got := sortArray(single); !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("expected single-item sorted list, got %v", got)
	}
}
