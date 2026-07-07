package xrand

import (
	"strings"
	"testing"
)

func TestRangeInt(t *testing.T) {
	if got := RangeInt(5, 5); got != 5 {
		t.Fatalf("RangeInt() = %d", got)
	}
	for range 20 {
		got := RangeInt(10, 1)
		if got < 1 || got > 10 {
			t.Fatalf("RangeInt() = %d, want in [1,10]", got)
		}
	}
}

func TestRangeInt64(t *testing.T) {
	for range 20 {
		got := RangeInt64(10, 1)
		if got < 1 || got > 10 {
			t.Fatalf("RangeInt64() = %d, want in [1,10]", got)
		}
	}
}

func TestRangeFloat64(t *testing.T) {
	if got := RangeFloat64(1.5, 1.5); got != 1.5 {
		t.Fatalf("RangeFloat64() = %f", got)
	}
	for range 20 {
		got := RangeFloat64(2, 1)
		if got < 1 || got >= 2 {
			t.Fatalf("RangeFloat64() = %f, want in [1,2)", got)
		}
	}
}

func TestBool(t *testing.T) {
	_ = Bool()
}

func TestPick(t *testing.T) {
	got := Pick([]int{1, 2}, 5)
	if len(got) != 5 {
		t.Fatalf("Pick() len = %d", len(got))
	}
	if got := Pick([]int{1, 2}, 0); len(got) != 0 {
		t.Fatalf("Pick() len = %d", len(got))
	}
}

func TestPickOne(t *testing.T) {
	if _, ok := PickOne([]int{}); ok {
		t.Fatal("PickOne(empty) ok = true")
	}
	got, ok := PickOne([]int{1})
	if !ok || got != 1 {
		t.Fatalf("PickOne() = %d, %v", got, ok)
	}
}

func TestPickUnique(t *testing.T) {
	src := []int{1, 2, 3}
	got := PickUnique(src, 5)
	if len(got) != len(src) {
		t.Fatalf("PickUnique() len = %d", len(got))
	}
	if len(src) != 3 || src[0] != 1 || src[1] != 2 || src[2] != 3 {
		t.Fatalf("PickUnique() modified src: %#v", src)
	}
}

func TestShuffle(t *testing.T) {
	src := []int{1, 2, 3}
	got := Shuffle(src)
	if len(got) != len(src) {
		t.Fatalf("Shuffle() len = %d", len(got))
	}
	if len(src) != 3 || src[0] != 1 || src[1] != 2 || src[2] != 3 {
		t.Fatalf("Shuffle() modified src: %#v", src)
	}

	ShuffleInPlace(src)
	if len(src) != 3 {
		t.Fatalf("ShuffleInPlace() len = %d", len(src))
	}
}

func TestChance(t *testing.T) {
	if Chance(0) {
		t.Fatal("Chance(0) = true")
	}
	if !Chance(1) {
		t.Fatal("Chance(1) = false")
	}
	if ChancePercent(0) {
		t.Fatal("ChancePercent(0) = true")
	}
	if !ChancePercent(100) {
		t.Fatal("ChancePercent(100) = false")
	}
}

func TestWeightedIndex(t *testing.T) {
	if _, ok := WeightedIndex(nil); ok {
		t.Fatal("WeightedIndex(nil) ok = true")
	}
	if _, ok := WeightedIndex([]int{0, -1}); ok {
		t.Fatal("WeightedIndex(empty weights) ok = true")
	}
	for range 20 {
		got, ok := WeightedIndex([]int{0, 10, 0})
		if !ok || got != 1 {
			t.Fatalf("WeightedIndex() = %d, %v", got, ok)
		}
	}
}

func TestWeightedPick(t *testing.T) {
	if _, ok := WeightedPick([]string{}, []int{1}); ok {
		t.Fatal("WeightedPick(empty values) ok = true")
	}
	got, ok := WeightedPick([]string{"a", "b", "c"}, []int{0, 5, 0})
	if !ok || got != "b" {
		t.Fatalf("WeightedPick() = %q, %v", got, ok)
	}
	got, ok = WeightedPick([]string{"a"}, []int{0, 10})
	if ok || got != "" {
		t.Fatalf("WeightedPick(truncated) = %q, %v", got, ok)
	}
}

func TestPermutation(t *testing.T) {
	if got := Permutation(0); len(got) != 0 {
		t.Fatalf("Permutation(0) len = %d", len(got))
	}

	got := Permutation(5)
	if len(got) != 5 {
		t.Fatalf("Permutation() len = %d", len(got))
	}

	seen := make(map[int]bool, len(got))
	for _, value := range got {
		if value < 0 || value >= 5 {
			t.Fatalf("Permutation() contains %d", value)
		}
		if seen[value] {
			t.Fatalf("Permutation() duplicated %d", value)
		}
		seen[value] = true
	}
}

func TestCryptoRandomHelpers(t *testing.T) {
	if got, err := Bytes(8); err != nil || len(got) != 8 {
		t.Fatalf("Bytes() len = %d, err = %v", len(got), err)
	}
	if got, err := Hex(4); err != nil || len(got) != 8 {
		t.Fatalf("Hex() = %q, %v", got, err)
	}
	if got, err := Number(6); err != nil || len(got) != 6 || containsOutside(got, NumberLetters) {
		t.Fatalf("Number() = %q, %v", got, err)
	}
	if _, err := String(1, ""); err == nil {
		t.Fatal("String(empty letters) error = nil")
	}
	if _, err := Bytes(-1); err == nil {
		t.Fatal("Bytes(negative) error = nil")
	}
}

func containsOutside(value string, letters string) bool {
	for _, r := range value {
		if !strings.ContainsRune(letters, r) {
			return true
		}
	}
	return false
}
