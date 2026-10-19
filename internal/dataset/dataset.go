// Package dataset builds the synthetic slop set and measures detection.
package dataset

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zysen1/Proofline/internal/gitx"
	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/verify"
)

const (
	ClassReal         = "real_fix"
	ClassPassesBefore = "passes_before"
	ClassTautology    = "tautology"
	ClassScope        = "scope_creep"
	ClassBroken       = "broken"
	ClassFalseClaims  = "false_claims"
)

// Classes is the synthetic taxonomy: five slop classes plus real fixes.
var Classes = []string{
	ClassReal,
	ClassPassesBefore,
	ClassTautology,
	ClassScope,
	ClassBroken,
	ClassFalseClaims,
}

var fnames = []string{
	"discount", "tax", "shipping", "tip", "late_fee", "quota", "rebate", "toll",
	"fare", "levy", "surcharge", "credit", "bonus", "penalty", "tariff", "wage",
	"stipend", "royalty", "coupon", "deposit",
}

// Spec is one generated case.
type Spec struct {
	Class   string
	Index   int
	Name    string
	K       int
	B       int
	Issue   int
	Adapter string
}

func specFor(class string, index int) Spec {
	return Spec{
		Class:   class,
		Index:   index,
		Name:    fnames[index%len(fnames)],
		K:       2 + index%5,
		B:       1 + index%7,
		Issue:   100 + index,
		Adapter: "pytest",
	}
}

// Materialize writes a git repo with branches main and pr.
func Materialize(dir string, sp Spec) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := gitx.InitRepo(dir); err != nil {
		return err
	}
	if err := writeTree(dir, sp, false); err != nil {
		return err
	}
	if err := gitx.CommitAll(dir, "base "+sp.Name); err != nil {
		return err
	}
	if err := gitx.CheckoutNew(dir, "pr"); err != nil {
		return err
	}
	if err := gitx.CleanWorktree(dir); err != nil {
		return err
	}
	if err := writeTree(dir, sp, true); err != nil {
		return err
	}
	return gitx.CommitAll(dir, fmt.Sprintf("pr %s %s\n\nFixes #%d\n", sp.Class, sp.Name, sp.Issue))
}

func writeTree(dir string, sp Spec, head bool) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	mod := pythonModule(sp, head)
	test := pythonTest(sp, head)
	body := pythonBody(sp, head)
	if err := os.WriteFile(filepath.Join(dir, "mod.py"), []byte(mod), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "test_mod.py"), []byte(test), 0o644); err != nil {
		return err
	}
	if head || sp.Class == ClassPassesBefore {
		if err := os.MkdirAll(filepath.Join(dir, ".proofline"), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, ".proofline", "pr.md"), []byte(body), 0o644); err != nil {
			return err
		}
	}
	if head && sp.Class == ClassScope {
		if err := os.MkdirAll(filepath.Join(dir, "extras"), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "extras", "noise.py"), []byte(noise(sp)), 0o644); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "requirements.txt"), []byte("requests==2.32.3\nflask==3.0.3\n"), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func pythonModule(sp Spec, head bool) string {
	correct := fmt.Sprintf(`def %s(n):
    # Fixes #%d
    if n < 0:
        raise ValueError("n must be >= 0")
    return n * %d + %d
`, sp.Name, sp.Issue, sp.K, sp.B)
	brokenImpl := fmt.Sprintf(`def %s(n):
    raise NotImplementedError("%s is not implemented")
`, sp.Name, sp.Name)
	switch sp.Class {
	case ClassPassesBefore:
		if head {
			return fmt.Sprintf(`def %s(n):
    # regression note %d
    if n < 0:
        raise ValueError("n must be >= 0")
    return n * %d + %d
`, sp.Name, sp.Index, sp.K, sp.B)
		}
		return correct
	case ClassBroken:
		if head {
			return fmt.Sprintf(`def %s(n)
    # Fixes #%d
    if n < 0:
        raise ValueError("n must be >= 0")
    return n * %d + %d
`, sp.Name, sp.Issue, sp.K, sp.B)
		}
		return brokenImpl
	default:
		if head {
			return correct
		}
		return brokenImpl
	}
}

func pythonTest(sp Spec, head bool) string {
	basic := 4*sp.K + sp.B
	strong := fmt.Sprintf(`import pytest

from mod import %s


def test_%s_basic():
    assert %s(4) == %d


def test_%s_zero():
    assert %s(0) == %d


def test_%s_rejects_negative():
    with pytest.raises(ValueError):
        %s(-1)
`, sp.Name, sp.Name, sp.Name, basic, sp.Name, sp.Name, sp.B, sp.Name, sp.Name)
	weak := fmt.Sprintf(`from mod import %s


def test_%s_runs():
    got = %s(4)
    assert got == got
    assert True
`, sp.Name, sp.Name, sp.Name)
	switch sp.Class {
	case ClassPassesBefore:
		if !head {
			return "def test_placeholder():\n    assert True\n"
		}
		return strong
	case ClassTautology:
		if !head {
			return "def test_placeholder():\n    assert True\n"
		}
		return weak
	default:
		if !head {
			return "def test_placeholder():\n    assert True\n"
		}
		return strong
	}
}

func pythonBody(sp Spec, head bool) string {
	repro := fmt.Sprintf("python3 -m pytest test_mod.py::test_%s_basic", sp.Name)
	if sp.Class == ClassTautology {
		repro = fmt.Sprintf("python3 -m pytest test_mod.py::test_%s_runs", sp.Name)
	}
	switch sp.Class {
	case ClassFalseClaims:
		return "Fixes #9999\n\n`charge` in `payments.py` is fixed.\n\nRepro:\npython3 -m pytest test_payments.py\n"
	case ClassBroken:
		return fmt.Sprintf("Fixes #%d\n\n`%s` in `mod.py` applies the formula.\n", sp.Issue, sp.Name)
	case ClassPassesBefore:
		return fmt.Sprintf("`%s` in `mod.py` already returned the formula. This adds a regression test.\n\nRepro:\n%s\n", sp.Name, repro)
	default:
		return fmt.Sprintf("Fixes #%d\n\n`%s` in `mod.py` now returns n*%d+%d and rejects negatives.\n\nRepro:\n%s\n", sp.Issue, sp.Name, sp.K, sp.B, repro)
	}
}

func noise(sp Spec) string {
	var b strings.Builder
	b.WriteString("# generated by an agent — do not edit\n")
	b.WriteString("\"\"\"Unrelated analytics dump for case " + sp.Name + ".\"\"\"\n\n")
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&b, "def noise_%s_%d(value):\n    return value + %d\n\n", sp.Name, i, i)
	}
	return b.String()
}

// CaseResult is one measured case.
type CaseResult struct {
	Class      string            `json:"class"`
	Index      int               `json:"index"`
	Verdict    model.Verdict     `json:"verdict"`
	KillRate   float64           `json:"kill_rate"`
	Detected   bool              `json:"detected"`
	Primary    string            `json:"primary"`
	Checks     map[string]string `json:"checks"`
	DurationMS int64             `json:"duration_ms"`
	Error      string            `json:"error,omitempty"`
}

// Report is the measured synthetic run.
type Report struct {
	Dataset       string                  `json:"dataset"`
	PerClass      int                     `json:"per_class"`
	Reruns        int                     `json:"reruns"`
	MaxMutants    int                     `json:"max_mutants"`
	Adapter       string                  `json:"adapter"`
	WallClockMS   int64                   `json:"wall_clock_ms"`
	MeanCaseMS    int64                   `json:"mean_case_ms"`
	Classes       map[string]ClassSummary `json:"classes"`
	Cases         []CaseResult            `json:"cases"`
	FalsePositive float64                 `json:"false_positive_rate_unverified"`
	OSSReplay     string                  `json:"oss_replay"`
	Notes         string                  `json:"notes"`
}

// ClassSummary aggregates one class.
type ClassSummary struct {
	N             int     `json:"n"`
	Verified      int     `json:"verified"`
	Partially     int     `json:"partially"`
	Unverified    int     `json:"unverified"`
	NotApplicable int     `json:"not_applicable"`
	Detected      int     `json:"detected"`
	DetectionRate float64 `json:"detection_rate"`
	PrimaryCheck  string  `json:"primary_check,omitempty"`
}

// Measure runs perClass cases of each class. OSS replay is not performed.
func Measure(ctx context.Context, perClass, reruns, maxMutants int) (Report, error) {
	if perClass <= 0 {
		perClass = 20
	}
	if reruns <= 0 {
		reruns = 2
	}
	start := time.Now()
	root, err := os.MkdirTemp("", "proofline-measure-")
	if err != nil {
		return Report{}, err
	}
	defer os.RemoveAll(root)

	rep := Report{
		Dataset:    "synthetic-slop-v1",
		PerClass:   perClass,
		Reruns:     reruns,
		MaxMutants: maxMutants,
		Adapter:    "pytest",
		Classes:    map[string]ClassSummary{},
		OSSReplay:  "WIP. This run did not replay merged public pull requests.",
		Notes:      "Each case is a generated one-module git repository. Detection means the verdict is Unverified and the class's primary check failed. The false-positive rate is the share of real_fix cases judged Unverified.",
	}
	for _, class := range Classes {
		for i := 0; i < perClass; i++ {
			sp := specFor(class, i)
			dir := filepath.Join(root, class, fmt.Sprintf("%02d", i))
			cr := CaseResult{Class: class, Index: i, Primary: primary(class), Checks: map[string]string{}}
			if err := Materialize(dir, sp); err != nil {
				return Report{}, err
			}
			rec, err := verify.Verify(ctx, verify.Options{
				Repo: dir, Base: "main", Head: "pr",
				Reruns: reruns, MaxMutants: maxMutants, Timeout: 60 * time.Second,
			})
			if err != nil {
				cr.Error = err.Error()
				rep.Cases = append(rep.Cases, cr)
				continue
			}
			cr.Verdict = rec.Verdict
			cr.KillRate = rec.KillRate
			cr.DurationMS = rec.DurationMS
			for _, c := range rec.Checks {
				cr.Checks[c.ID] = string(c.Status)
			}
			cr.Detected = detected(class, rec)
			rep.Cases = append(rep.Cases, cr)
			fmt.Fprintf(os.Stderr, "measure %s %d → %s detected=%v\n", class, i, rec.Verdict, cr.Detected)
		}
	}
	rep.WallClockMS = time.Since(start).Milliseconds()
	if len(rep.Cases) > 0 {
		rep.MeanCaseMS = rep.WallClockMS / int64(len(rep.Cases))
	}
	rep.Classes = summarize(rep.Cases)
	if real, ok := rep.Classes[ClassReal]; ok && real.N > 0 {
		rep.FalsePositive = float64(real.Unverified) / float64(real.N)
	}
	return rep, nil
}

func primary(class string) string {
	switch class {
	case ClassPassesBefore:
		return "fail_before"
	case ClassTautology:
		return "mutation"
	case ClassScope:
		return "scope"
	case ClassBroken:
		return "pass_after"
	case ClassFalseClaims:
		return "claims"
	default:
		return ""
	}
}

func detected(class string, rec model.Receipt) bool {
	if class == ClassReal {
		return rec.Verdict == model.VerdictVerified
	}
	if rec.Verdict != model.VerdictUnverified {
		return false
	}
	id := primary(class)
	c, ok := rec.Find(id)
	return ok && c.Status == model.StatusFail
}

func summarize(cases []CaseResult) map[string]ClassSummary {
	out := map[string]ClassSummary{}
	for _, class := range Classes {
		s := ClassSummary{PrimaryCheck: primary(class)}
		for _, c := range cases {
			if c.Class != class {
				continue
			}
			s.N++
			switch c.Verdict {
			case model.VerdictVerified:
				s.Verified++
			case model.VerdictPartially:
				s.Partially++
			case model.VerdictUnverified:
				s.Unverified++
			case model.VerdictNA:
				s.NotApplicable++
			}
			if class == ClassReal {
				continue
			}
			if c.Detected {
				s.Detected++
			}
		}
		if class != ClassReal && s.N > 0 {
			s.DetectionRate = float64(s.Detected) / float64(s.N)
		}
		out[class] = s
	}
	return out
}

// Write serializes a report.
func Write(path string, rep Report) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}
