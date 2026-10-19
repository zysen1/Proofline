// Package receipt renders a proof receipt as a pull-request comment.
package receipt

import (
	"fmt"
	"strings"

	"github.com/zysen1/Proofline/internal/model"
)

// Markdown is the comment body. Logs sit inside details blocks.
func Markdown(r model.Receipt) string {
	var b strings.Builder
	b.WriteString("## Proofline receipt\n\n")
	b.WriteString(r.Badge + "\n\n")
	b.WriteString(fmt.Sprintf("**Verdict: %s**\n\n", r.Verdict))
	if r.NotApplicableReason != "" {
		b.WriteString(r.NotApplicableReason + " This is not a negative judgment.\n\n")
	}
	fmt.Fprintf(&b, "Base `%s` (`%s`) → head `%s` (`%s`).\n\n", r.Base, short(r.BaseSHA), r.Head, short(r.HeadSHA))
	if r.Adapter != "" {
		fmt.Fprintf(&b, "Test adapter: `%s`. ", r.Adapter)
	}
	fmt.Fprintf(&b, "Wall clock %d ms.\n\n", r.DurationMS)
	b.WriteString("| Check | Result | Evidence |\n|---|---|---|\n")
	for _, c := range r.Checks {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", c.Name, c.Status, escape(c.Summary))
	}
	b.WriteString("\n")
	if r.Mutants > 0 {
		fmt.Fprintf(&b, "Mutation kill rate: %s.\n\n", model.FormatRate(r.Killed, r.Mutants))
	}
	if r.Summary != "" {
		fmt.Fprintf(&b, "### Summary (%s)\n\n%s\n\n", r.SummarySource, r.Summary)
	}
	b.WriteString("<details>\n<summary>Commands and logs</summary>\n\n")
	for _, c := range r.Checks {
		fmt.Fprintf(&b, "#### %s\n\n", c.Name)
		if len(c.Evidence) == 0 {
			b.WriteString("_No command log._\n\n")
			continue
		}
		for _, ev := range c.Evidence {
			if ev.Detail != "" {
				fmt.Fprintf(&b, "- %s\n", escape(ev.Detail))
			}
			if ev.Command != "" {
				fmt.Fprintf(&b, "\n```\n%s\n```\n", ev.Command)
			}
			if ev.ExitCode != nil {
				fmt.Fprintf(&b, "\nExit code `%d`.\n", *ev.ExitCode)
			}
			if ev.Log != "" {
				fmt.Fprintf(&b, "\n```\n%s\n```\n", clip(ev.Log, 1200))
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("</details>\n\n")
	b.WriteString("Proofline judged this from execution evidence. A model summary, when present, does not change the verdict.\n")
	return b.String()
}

func short(sha string) string {
	if len(sha) > 10 {
		return sha[:10]
	}
	return sha
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…[truncated]"
}
