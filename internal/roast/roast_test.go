package roast

import (
	"strings"
	"testing"

	"github.com/optiqor/optiqor-cli/internal/render"
	"github.com/optiqor/optiqor-cli/pkg/rules"
)

func TestApply_RewritesKnownTitles(t *testing.T) {
	in := render.Report{
		Findings: []rules.Finding{
			{DetectorID: "cpu-overprovisioned", Title: "CPU request appears overprovisioned"},
			{DetectorID: "missing-memory-limit", Title: "Memory limit not set"},
		},
	}
	out := Apply(in)
	if out.Findings[0].Title == in.Findings[0].Title {
		t.Errorf("cpu-overprovisioned title should have been rewritten")
	}
	if out.Findings[1].Title == in.Findings[1].Title {
		t.Errorf("missing-memory-limit title should have been rewritten")
	}
}

func TestApply_LeavesUnknownTitleAlone(t *testing.T) {
	in := render.Report{
		Findings: []rules.Finding{
			{DetectorID: "future-detector-not-yet-roasted", Title: "Some title"},
		},
	}
	out := Apply(in)
	if out.Findings[0].Title != "Some title" {
		t.Errorf("unknown detector title was rewritten: %q", out.Findings[0].Title)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	in := render.Report{
		Findings: []rules.Finding{
			{DetectorID: "cpu-overprovisioned", Title: "original"},
		},
	}
	original := in.Findings[0].Title
	_ = Apply(in)
	if in.Findings[0].Title != original {
		t.Errorf("input mutated: %q != %q", in.Findings[0].Title, original)
	}
}

func TestApply_PreservesMaterialFields(t *testing.T) {
	// Hard rule: only Title changes. Detail, MonthlyUSDCents,
	// Severity, Confidence, Category, Signal must round-trip.
	in := render.Report{
		Findings: []rules.Finding{{
			DetectorID:      "cpu-overprovisioned",
			Workload:        "api",
			Title:           "CPU request appears overprovisioned",
			Detail:          "Request 2 vs limit 2.5",
			MonthlyUSDCents: 12345,
			Severity:        rules.SeverityMed,
			Confidence:      rules.ConfidenceMed,
			Category:        rules.CategoryCost,
			Signal: &rules.Signal{
				Label: "CPU", Have: 2, Want: 2.5,
				HaveDisplay: "2", WantDisplay: "2.5", Note: "80% of limit",
			},
		}},
	}
	out := Apply(in).Findings[0]
	if out.Detail != "Request 2 vs limit 2.5" {
		t.Errorf("Detail rewritten")
	}
	if out.MonthlyUSDCents != 12345 {
		t.Errorf("MonthlyUSDCents changed")
	}
	if out.Severity != rules.SeverityMed {
		t.Errorf("Severity changed")
	}
	if out.Confidence != rules.ConfidenceMed {
		t.Errorf("Confidence changed")
	}
	if out.Category != rules.CategoryCost {
		t.Errorf("Category changed")
	}
	if out.Signal == nil || out.Signal.Note != "80% of limit" {
		t.Errorf("Signal lost or rewritten")
	}
}

func TestTagline_AndFooter_NotEmpty(t *testing.T) {
	if strings.TrimSpace(Tagline) == "" {
		t.Error("Tagline empty")
	}
	if strings.TrimSpace(FooterQuip) == "" {
		t.Error("FooterQuip empty")
	}
}
