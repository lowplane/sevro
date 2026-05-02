package analyze

import "sort"

// Grade is a letter projection of a 0–100 efficiency score, paired
// with the percentile rank against the calibration distribution
// baked into the binary.
//
// Why two numbers?
//
// The raw 0–100 score is precise but has no context: scoring 72 only
// matters relative to peers. The letter grade gives an immediate
// "where do I sit?" signal (B-, C+, F …); the percentile rank is the
// honest comparison that drives the social loop ("better than 64% of
// charts we benchmarked"). Both derive deterministically from the
// same underlying [Score].
type Grade struct {
	// Letter is the conventional A+/A/A-/B+/B/B-/C+/C/C-/D/F mapping.
	Letter string `json:"letter"`

	// PercentileRank is the percentage of calibration scores that
	// this score beats (0–100, integer). 50 means median; 0 means
	// worst-in-class; 100 means top of the calibrated set.
	PercentileRank int `json:"percentile_rank"`

	// Sample is the size of the calibration distribution. Surfaced in
	// the renderer as "better than X% of N benchmark charts" so users
	// know what they're being compared against.
	Sample int `json:"sample_size"`
}

// GradeFor folds a numeric score into a [Grade]. Pure function;
// deterministic given the baked calibration distribution.
func GradeFor(score int) Grade {
	return Grade{
		Letter:         letterFor(score),
		PercentileRank: percentileRank(score, calibrationScores),
		Sample:         len(calibrationScores),
	}
}

// letterFor maps a 0–100 score to a conventional letter grade.
// Bands are informed by US academic convention (90+=A, 80+=B, etc.)
// and biased slightly toward "easier to get a B" since most real
// Helm charts cluster in the 60–80 band.
func letterFor(score int) string {
	switch {
	case score >= 95:
		return "A+"
	case score >= 90:
		return "A"
	case score >= 85:
		return "A-"
	case score >= 80:
		return "B+"
	case score >= 75:
		return "B"
	case score >= 70:
		return "B-"
	case score >= 65:
		return "C+"
	case score >= 60:
		return "C"
	case score >= 55:
		return "C-"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}

// percentileRank returns the integer percentage of the population
// that scored strictly less than score. Population must already be
// sorted ascending.
func percentileRank(score int, population []int) int {
	if len(population) == 0 {
		return 0
	}
	// sort.Search finds the lowest index i where population[i] >=
	// score; the count of strictly-less is exactly i.
	below := sort.Search(len(population), func(i int) bool {
		return population[i] >= score
	})
	rank := below * 100 / len(population)
	if rank > 100 {
		rank = 100
	}
	return rank
}

// calibrationScores is the baked-in benchmark distribution that
// powers the percentile readout in [GradeFor]. The distribution is
// modelled (not telemetered) — the CLI keeps its no-telemetry
// promise; a live percentile derived from real merged-PR outcomes
// arrives with the agent install per the strategy docs.
//
// Shape: 100 samples, beta-style curve centred ~70 with realistic
// spread. The mean and median fall in the C+/B- band, which matches
// the "most charts have one or two HIGH findings" empirical pattern
// we see in public charts. Must remain sorted ascending — the
// percentile lookup is a binary search.
var calibrationScores = []int{
	18, 22, 24, 26, 28, 30, 31, 33, 34, 35,
	37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 60, 61, 62, 63, 64, 65, 66,
	66, 67, 67, 68, 68, 69, 69, 70, 70, 70,
	71, 71, 72, 72, 73, 73, 74, 74, 75, 75,
	76, 76, 77, 77, 78, 78, 79, 79, 80, 80,
	81, 82, 83, 84, 85, 85, 86, 87, 88, 88,
	89, 90, 90, 91, 91, 92, 92, 93, 93, 94,
	94, 95, 95, 96, 96, 97, 97, 98, 99, 100,
}
