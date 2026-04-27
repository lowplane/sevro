package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoad_NoFileReturnsZero(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	_ = os.Chdir(dir)

	c, err := Load("")
	if err != nil {
		t.Fatalf("expected zero config, got err: %v", err)
	}
	if !reflect.DeepEqual(c, Config{}) {
		t.Errorf("expected zero, got %+v", c)
	}
}

func TestLoad_PicksUpDotSevro(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	_ = os.Chdir(dir)

	body := `min_severity: med
detectors:
  - cpu-overprovisioned
  - missing-memory-limit
fail_on: high
no_color: true
`
	if err := os.WriteFile(filepath.Join(dir, ConfigName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.MinSeverity != "med" || c.FailOn != "high" || !c.NoColor {
		t.Errorf("config not loaded: %+v", c)
	}
	if len(c.Detectors) != 2 {
		t.Errorf("detectors lost: %v", c.Detectors)
	}
}

func TestLoad_ExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yaml")
	if err := os.WriteFile(path, []byte("min_severity: high"), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.MinSeverity != "high" {
		t.Errorf("MinSeverity = %q", c.MinSeverity)
	}
}

func TestLoad_ExplicitMissingErrors(t *testing.T) {
	if _, err := Load("/no/such/path"); err == nil {
		t.Fatal("expected error on missing explicit path")
	}
}

func TestLoad_EnvVarFallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "env-config.yaml")
	if err := os.WriteFile(path, []byte("min_severity: low"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SEVRO_CONFIG", path)
	c, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if c.MinSeverity != "low" {
		t.Errorf("MinSeverity = %q", c.MinSeverity)
	}
}

func TestValidate_BadSeverity(t *testing.T) {
	c := Config{MinSeverity: "noisy"}
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "min_severity") {
		t.Fatalf("expected min_severity validation error, got %v", err)
	}
}

func TestValidate_AcceptsAliases(t *testing.T) {
	for _, v := range []string{"low", "med", "medium", "high", "LOW", "Med", "HIGH"} {
		c := Config{MinSeverity: v}
		if err := c.Validate(); err != nil {
			t.Errorf("Validate(%q): unexpected error %v", v, err)
		}
	}
}

func TestDecode_EmptyOK(t *testing.T) {
	c, err := Decode(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, Config{}) {
		t.Errorf("empty body should yield zero, got %+v", c)
	}
}

func TestDecode_BadYAML(t *testing.T) {
	if _, err := Decode(strings.NewReader("not: valid: yaml::")); err == nil {
		t.Fatal("expected yaml error")
	}
}
