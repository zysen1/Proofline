// Package model holds the proof receipt types shared by the CLI and the checks.
package model

import "fmt"

const Version = "0.1.0"

type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
	StatusWarn Status = "warn"
	StatusNA   Status = "not_applicable"
)

type Verdict string

const (
	VerdictVerified   Verdict = "Verified"
	VerdictPartially  Verdict = "Partially"
	VerdictUnverified Verdict = "Unverified"
	VerdictNA         Verdict = "Not applicable"
)

// Evidence is one command, log, or static observation behind a check.
type Evidence struct {
	Command  string `json:"command,omitempty"`
	ExitCode *int   `json:"exit_code,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Log      string `json:"log,omitempty"`
}

// Check is one deterministic proof step.
type Check struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Status     Status     `json:"status"`
	Summary    string     `json:"summary"`
	Evidence   []Evidence `json:"evidence,omitempty"`
	DurationMS int64      `json:"duration_ms"`
}

// Receipt is the proof Proofline attaches to a branch.
type Receipt struct {
	Tool                string         `json:"tool"`
	Version             string         `json:"version"`
	Repo                string         `json:"repo"`
	Base                string         `json:"base"`
	Head                string         `json:"head"`
	BaseSHA             string         `json:"base_sha"`
	HeadSHA             string         `json:"head_sha"`
	Adapter             string         `json:"adapter,omitempty"`
	Verdict             Verdict        `json:"verdict"`
	Badge               string         `json:"badge"`
	Checks              []Check        `json:"checks"`
	KillRate            float64        `json:"kill_rate"`
	Mutants             int            `json:"mutants"`
	Killed              int            `json:"killed"`
	MutantsByOp         map[string]int `json:"mutants_by_op,omitempty"`
	KilledByOp          map[string]int `json:"killed_by_op,omitempty"`
	DurationMS          int64          `json:"duration_ms"`
	Summary             string         `json:"summary,omitempty"`
	SummarySource       string         `json:"summary_source,omitempty"`
	NotApplicableReason string         `json:"not_applicable_reason,omitempty"`
	Score               int            `json:"score"`
}

// ChangedFile is one path in the three-dot diff.
type ChangedFile struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	Added    int    `json:"added"`
	Deleted  int    `json:"deleted"`
	Kind     string `json:"kind"`
	NewLines []int  `json:"new_lines,omitempty"`
	Patch    string `json:"-"`
}

func (r Receipt) Find(id string) (Check, bool) {
	for _, c := range r.Checks {
		if c.ID == id {
			return c, true
		}
	}
	return Check{}, false
}

// Decide turns check statuses into a verdict.
// Not applicable is decided by the caller before this function; a true na
// argument forces that verdict and does not count as a negative result.
func Decide(notApplicable bool, checks []Check) Verdict {
	if notApplicable {
		return VerdictNA
	}
	fail, warn := 0, 0
	for _, c := range checks {
		switch c.Status {
		case StatusFail:
			fail++
		case StatusWarn:
			warn++
		}
	}
	if fail > 0 {
		return VerdictUnverified
	}
	if warn > 0 {
		return VerdictPartially
	}
	return VerdictVerified
}

// EvidenceScore ranks branches by how much deterministic evidence they produced.
// Passing checks dominate. Mutation kill rate breaks ties inside a check.
// The prose summary is not an input.
func EvidenceScore(r Receipt) int {
	if r.Verdict == VerdictNA {
		return 0
	}
	score := 0
	for _, c := range r.Checks {
		switch c.Status {
		case StatusPass:
			score += 20
		case StatusWarn:
			score += 8
		}
	}
	score += int(r.KillRate * 20)
	if r.Verdict == VerdictVerified {
		score += 10
	}
	return score
}

// BadgeMarkdown is the shield used in the pull-request comment.
func BadgeMarkdown(v Verdict) string {
	switch v {
	case VerdictVerified:
		return "![Verified](https://img.shields.io/badge/Proofline-Verified-2ea44f)"
	case VerdictPartially:
		return "![Partially](https://img.shields.io/badge/Proofline-Partially-dbab09)"
	case VerdictUnverified:
		return "![Unverified](https://img.shields.io/badge/Proofline-Unverified-d73a49)"
	default:
		return "![Not applicable](https://img.shields.io/badge/Proofline-Not%20applicable-6e7781)"
	}
}

// ExitCode is the process status for a receipt.
// Verified and Not applicable are 0. Partially is 2. Unverified is 1.
func ExitCode(v Verdict) int {
	switch v {
	case VerdictVerified, VerdictNA:
		return 0
	case VerdictPartially:
		return 2
	default:
		return 1
	}
}

func (s Status) String() string { return string(s) }

func FormatRate(killed, total int) string {
	if total == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%d/%d (%.0f%%)", killed, total, 100*float64(killed)/float64(total))
}
