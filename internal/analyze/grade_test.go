package analyze

import (
	"sort"
	"testing"
)

func TestLetterFor_BoundaryConditions(t *testing.T) {
	cases := map[int]string{
		100: "A+",
		95:  "A+",
		94:  "A",
		90:  "A",
		89:  "A-",
		85:  "A-",
		84:  "B+",
		80:  "B+",
		79:  "B",
		75:  "B",
		74:  "B-",
		70:  "B-",
		69:  "C+",
		65:  "C+",
		64:  "C",
		60:  "C",
		59:  "C-",
		55:  "C-",
		54:  "D",
		50:  "D",
		49:  "F",
		0:   "F",
	}
	for in, want := range cases {
		if got := letterFor(in); got != want {
			t.Errorf("letterFor(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestPercentileRank_KnownDistribution(t *testing.T) {
	pop := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	cases := map[int]int{
		5:   0,   // worse than everything
		10:  0,   // equal to lowest, beats nothing
		35:  30,  // beats 10, 20, 30 → 30%
		50:  40,  // beats 10..40 → 40%
		95:  90,  // beats 10..90 → 90%
		200: 100, // beats everything
	}
	for in, want := range cases {
		if got := percentileRank(in, pop); got != want {
			t.Errorf("percentileRank(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestCalibrationScores_Invariants(t *testing.T) {
	// Must be sorted ascending; the percentile lookup is a binary
	// search.
	if !sort.IntsAreSorted(calibrationScores) {
		t.Fatal("calibrationScores must be sorted ascending")
	}
	// Must be 100 entries so percentile == count of beaten samples;
	// users learn this implicitly from the rendered text.
	if got := len(calibrationScores); got != 100 {
		t.Fatalf("calibration sample size = %d, want 100", got)
	}
	// Bounds sanity: every sample is a valid score.
	for _, v := range calibrationScores {
		if v < 0 || v > 100 {
			t.Errorf("calibration sample out of range: %d", v)
		}
	}
}

func TestGradeFor_EndToEnd(t *testing.T) {
	g := GradeFor(78)
	if g.Letter != "B" {
		t.Errorf("GradeFor(78).Letter = %q, want B", g.Letter)
	}
	if g.Sample != len(calibrationScores) {
		t.Errorf("Sample = %d", g.Sample)
	}
	// 78 should beat the bottom of the distribution but not the top;
	// concrete value depends on the curve. Just check it's in (0,100).
	if g.PercentileRank <= 0 || g.PercentileRank >= 100 {
		t.Errorf("percentile rank = %d, want strictly between 0 and 100", g.PercentileRank)
	}
}
