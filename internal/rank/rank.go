// Package rank orders branches by deterministic evidence strength.
package rank

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/receipt"
)

// Row is one ranked branch.
type Row struct {
	Rank    int
	Receipt model.Receipt
}

// Order sorts a copy. Higher evidence score comes first.
// Equal scores break on kill rate, then on the head name.
func Order(in []model.Receipt) []Row {
	cp := append([]model.Receipt(nil), in...)
	sort.SliceStable(cp, func(i, j int) bool {
		if cp[i].Score != cp[j].Score {
			return cp[i].Score > cp[j].Score
		}
		if cp[i].KillRate != cp[j].KillRate {
			return cp[i].KillRate > cp[j].KillRate
		}
		return cp[i].Head < cp[j].Head
	})
	rows := make([]Row, len(cp))
	for i := range cp {
		rows[i] = Row{Rank: i + 1, Receipt: cp[i]}
	}
	return rows
}

// Markdown renders the leaderboard.
func Markdown(rows []Row) string {
	var b strings.Builder
	b.WriteString("# Proofline rank\n\n")
	b.WriteString("Branches are ordered by execution evidence, not by how polished the diff looks.\n\n")
	b.WriteString("| Rank | Branch | Verdict | Score | Kill rate | Why it is here |\n|---|---|---|---:|---:|---|\n")
	for _, row := range rows {
		r := row.Receipt
		fmt.Fprintf(&b, "| %d | `%s` | %s | %d | %s | %s |\n",
			row.Rank, r.Head, r.Verdict, r.Score, model.FormatRate(r.Killed, r.Mutants), escape(why(r)))
	}
	b.WriteString("\nScore: +20 per passing check, +8 per warning, +20×kill rate, +10 when the verdict is Verified. Not-applicable branches score 0. The summary text is not an input.\n")
	return b.String()
}

func why(r model.Receipt) string {
	if r.Verdict == model.VerdictNA {
		return r.NotApplicableReason
	}
	var failed []string
	for _, c := range r.Checks {
		if c.Status == model.StatusFail {
			failed = append(failed, c.Name+": "+c.Summary)
		}
	}
	if len(failed) == 0 {
		if r.Verdict == model.VerdictVerified {
			return "Fail-before held, tests passed, mutants died, claims and scope matched."
		}
		return "No failed check, but at least one warning remains."
	}
	return failed[0]
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 180 {
		return s[:180] + "…"
	}
	return s
}

// Comment bundles a rank table plus each receipt, for a single maintainer comment.
func Comment(rows []Row) string {
	var b strings.Builder
	b.WriteString(Markdown(rows))
	b.WriteString("\n")
	for _, row := range rows {
		b.WriteString("\n---\n\n")
		b.WriteString(receipt.Markdown(row.Receipt))
	}
	return b.String()
}
