package xmath

import "testing"

func TestMinInt64(t *testing.T) {
	if got := MinInt64(2, 1); got != 1 {
		t.Fatalf("MinInt64() = %d", got)
	}
}

func TestMinMaxInt(t *testing.T) {
	if got := MinInt(2, 1); got != 1 {
		t.Fatalf("MinInt() = %d", got)
	}
	if got := MaxInt(2, 1); got != 2 {
		t.Fatalf("MaxInt() = %d", got)
	}
	if got := MaxInt64(2, 1); got != 2 {
		t.Fatalf("MaxInt64() = %d", got)
	}
}

func TestAbs(t *testing.T) {
	if got := AbsInt(-3); got != 3 {
		t.Fatalf("AbsInt() = %d", got)
	}
	if got := AbsInt64(-3); got != 3 {
		t.Fatalf("AbsInt64() = %d", got)
	}
}

func TestClamp(t *testing.T) {
	if got := ClampInt(10, 1, 5); got != 5 {
		t.Fatalf("ClampInt() = %d", got)
	}
	if got := ClampInt(0, 5, 1); got != 1 {
		t.Fatalf("ClampInt(reverse) = %d", got)
	}
	if got := ClampInt64(3, 1, 5); got != 3 {
		t.Fatalf("ClampInt64() = %d", got)
	}
}

func TestInRange(t *testing.T) {
	if !InRangeInt(3, 5, 1) {
		t.Fatal("InRangeInt(reverse) = false")
	}
	if InRangeInt64(6, 1, 5) {
		t.Fatal("InRangeInt64() = true")
	}
}

func TestFloatRangeHelpers(t *testing.T) {
	if got := MinFloat64(2.5, 1.5); got != 1.5 {
		t.Fatalf("MinFloat64() = %f", got)
	}
	if got := MaxFloat64(2.5, 1.5); got != 2.5 {
		t.Fatalf("MaxFloat64() = %f", got)
	}
	if got := ClampFloat64(10.5, 1.5, 5.5); got != 5.5 {
		t.Fatalf("ClampFloat64() = %f", got)
	}
	if !InRangeFloat64(3.5, 5.5, 1.5) {
		t.Fatal("InRangeFloat64(reverse) = false")
	}
}

func TestRoundFloat64(t *testing.T) {
	if got := RoundFloat64(12.345, 2); got != 12.35 {
		t.Fatalf("RoundFloat64() = %f", got)
	}
	if got := RoundFloat64(126, -1); got != 130 {
		t.Fatalf("RoundFloat64(negative precision) = %f", got)
	}
}

func TestFloorCeilTruncFloat64(t *testing.T) {
	if got := FloorFloat64(12.349, 2); got != 12.34 {
		t.Fatalf("FloorFloat64() = %f", got)
	}
	if got := CeilFloat64(12.341, 2); got != 12.35 {
		t.Fatalf("CeilFloat64() = %f", got)
	}
	if got := TruncFloat64(12.349, 2); got != 12.34 {
		t.Fatalf("TruncFloat64() = %f", got)
	}
}

func TestSum(t *testing.T) {
	if got := SumInt([]int{1, 2, 3}); got != 6 {
		t.Fatalf("SumInt() = %d", got)
	}
	if got := SumInt64([]int64{1, 2, 3}); got != 6 {
		t.Fatalf("SumInt64() = %d", got)
	}
	if got := SumFloat64([]float64{1.5, 2.5}); got != 4 {
		t.Fatalf("SumFloat64() = %f", got)
	}
}

func TestAvg(t *testing.T) {
	if _, ok := AvgInt(nil); ok {
		t.Fatal("AvgInt(nil) ok = true")
	}
	if got, ok := AvgInt([]int{1, 2, 3}); !ok || got != 2 {
		t.Fatalf("AvgInt() = %f, %v", got, ok)
	}
	if got, ok := AvgInt64([]int64{1, 2, 3}); !ok || got != 2 {
		t.Fatalf("AvgInt64() = %f, %v", got, ok)
	}
	if got, ok := AvgFloat64([]float64{1.5, 2.5}); !ok || got != 2 {
		t.Fatalf("AvgFloat64() = %f, %v", got, ok)
	}
}

func TestPercentRatio(t *testing.T) {
	if _, ok := Percent(1, 0); ok {
		t.Fatal("Percent(div zero) ok = true")
	}
	if got, ok := Percent(2, 4); !ok || got != 50 {
		t.Fatalf("Percent() = %f, %v", got, ok)
	}
	if got, ok := Ratio(2, 4); !ok || got != 0.5 {
		t.Fatalf("Ratio() = %f, %v", got, ok)
	}
}
