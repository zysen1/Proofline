package report

import (
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/rank"
)

// Page is the offline demo report.
type Page struct {
	Title    string
	Subtitle string
	Receipts []model.Receipt
	Rank     []rank.Row
	Extra    *model.Receipt
	Foot     string
}

// WriteHTML writes a self-contained report.
func WriteHTML(path string, page Page) error {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\">")
	b.WriteString("<title>" + html.EscapeString(page.Title) + "</title>")
	b.WriteString(`<style>
body{margin:0;background:#f4f0e6;color:#1c1915;font:15px/1.45 "Iowan Old Style",Palatino,"Palatino Linotype",serif}
main{max-width:1080px;margin:0 auto;padding:28px 32px 48px}
h1{font-size:34px;line-height:1.1;margin:0 0 6px}
.sub{margin:0 0 22px;color:#4a453c}
table{width:100%;border-collapse:collapse;background:#fffdf8;margin:8px 0 22px}
th,td{border-bottom:1px solid #e4dccb;text-align:left;padding:8px 10px;vertical-align:top}
th{font:12px/1.3 ui-sans-serif,system-ui,sans-serif;letter-spacing:.04em;text-transform:uppercase;color:#6b6458}
code{font:12px/1.3 ui-monospace,SFMono-Regular,Menlo,monospace}
.pill{display:inline-block;padding:2px 8px;border-radius:999px;font:12px/1.4 ui-sans-serif,system-ui,sans-serif;color:#fff}
.Verified{background:#1f7a4d}.Unverified{background:#a33b32}.Partially{background:#8a6a12}
.na{background:#5c6570}
.why{color:#3f3a33}
h2{font-size:20px;margin:22px 0 8px}
footer{margin-top:12px;color:#5c564c;font:13px/1.4 ui-sans-serif,system-ui,sans-serif}
</style></head><body><main>`)
	fmt.Fprintf(&b, "<h1>%s</h1><p class=\"sub\">%s</p>", html.EscapeString(page.Title), html.EscapeString(page.Subtitle))
	b.WriteString("<h2>Six pull requests, one green badge</h2>")
	b.WriteString("<table><thead><tr><th>Branch</th><th>Verdict</th><th>What the evidence says</th></tr></thead><tbody>")
	for _, r := range page.Receipts {
		fmt.Fprintf(&b, "<tr><td><code>%s</code></td><td>%s</td><td class=\"why\">%s</td></tr>",
			html.EscapeString(r.Head), pill(r.Verdict), html.EscapeString(oneLine(r)))
	}
	b.WriteString("</tbody></table>")
	b.WriteString("<h2>Same bug, four agent branches</h2>")
	b.WriteString("<table><thead><tr><th>Rank</th><th>Branch</th><th>Verdict</th><th>Score</th><th>Kill rate</th></tr></thead><tbody>")
	for _, row := range page.Rank {
		fmt.Fprintf(&b, "<tr><td>%d</td><td><code>%s</code></td><td>%s</td><td>%d</td><td>%s</td></tr>",
			row.Rank, html.EscapeString(row.Receipt.Head), pill(row.Receipt.Verdict), row.Receipt.Score,
			html.EscapeString(model.FormatRate(row.Receipt.Killed, row.Receipt.Mutants)))
	}
	b.WriteString("</tbody></table>")
	if page.Extra != nil {
		b.WriteString("<h2>Flakiness, separately</h2><table><tbody><tr><td><code>flaky-counter</code></td><td>")
		b.WriteString(pill(page.Extra.Verdict))
		b.WriteString("</td><td class=\"why\">")
		b.WriteString(html.EscapeString(oneLine(*page.Extra)))
		b.WriteString("</td></tr></tbody></table>")
	}
	fmt.Fprintf(&b, "<footer>%s</footer></main></body></html>", html.EscapeString(page.Foot))
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func pill(v model.Verdict) string {
	cls := string(v)
	if v == model.VerdictNA {
		cls = "na"
	}
	return fmt.Sprintf("<span class=\"pill %s\">%s</span>", cls, html.EscapeString(string(v)))
}

func oneLine(r model.Receipt) string {
	if r.NotApplicableReason != "" {
		return r.NotApplicableReason
	}
	for _, c := range r.Checks {
		if c.Status == model.StatusFail {
			return c.Summary
		}
	}
	if r.Verdict == model.VerdictVerified {
		return "New tests failed before the fix, passed after, and killed the mutants."
	}
	return string(r.Verdict)
}
