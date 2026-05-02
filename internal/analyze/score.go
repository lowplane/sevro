package analyze

import "github.com/optiqor/optiqor-cli/pkg/rules"

// Score is the result of `optiqor score [chart]` — a 0-100 efficiency
// score derived from the severities of detector findings, plus the
// qualitative confidence band, a letter grade, and a percentile rank
// against the baked-in calibration distribution.
//
// The numerical value is "100 minus the per-finding penalty cap";
// numerical Confidence Scores arrive in Year 2 once we have enough
// merged-PR outcomes to calibrate. For now Confidence is qualitative.
//
// Grade turns the abstract score into a screenshot-friendly social
// signal ("B+ · better than 64% of charts"); it is fully derived
// from Value and the static distribution in [GradeFor].
type Score struct {
	Workloads int              `json:"workloads_analyzed"`
	Source    string           `json:"source"`
	Value     int              `json:"score"` // 0-100
	Band      rules.Confidence `json:"confidence_band"`
	Grade     Grade            `json:"grade"`
	Penalties map[string]int   `json:"penalties"` // detector_id -> penalty points
	Findings  []rules.Finding  `json:"findings"`
}

// Penalty weights per severity. High-severity findings drag the score
// down faster than low-severity ones; the cap (100) is the maximum
// penalty any single workload can incur, so a chart with one HIGH
// finding never scores below ~50 from that finding alone.
const (
	penaltyHigh = 25
	penaltyMed  = 10
	penaltyLow  = 3
	penaltyInfo = 1
	penaltyCap  = 100
)

// Compute folds findings into a Score.
func Compute(source string, workloads int, findings []rules.Finding) Score {
	penalties := map[string]int{}
	total := 0
	for _, f := range findings {
		p := penaltyFor(f.Severity)
		penalties[f.DetectorID] += p
		total += p
	}
	if total > penaltyCap {
		total = penaltyCap
	}
	value := 100 - total
	if value < 0 {
		value = 0
	}
	return Score{
		Workloads: workloads,
		Source:    source,
		Value:     value,
		Band:      bandFor(value),
		Grade:     GradeFor(value),
		Penalties: penalties,
		Findings:  findings,
	}
}

func penaltyFor(s rules.Severity) int {
	switch s {
	case rules.SeverityHigh:
		return penaltyHigh
	case rules.SeverityMed:
		return penaltyMed
	case rules.SeverityLow:
		return penaltyLow
	default:
		return penaltyInfo
	}
}

// bandFor maps a numerical score into a qualitative confidence band.
// Stable mapping; Year-2 numerical scores derive their confidence
// differently (calibrated against measured outcomes).
func bandFor(score int) rules.Confidence {
	switch {
	case score >= 85:
		return rules.ConfidenceHigh
	case score >= 60:
		return rules.ConfidenceMed
	default:
		return rules.ConfidenceLow
	}
}
