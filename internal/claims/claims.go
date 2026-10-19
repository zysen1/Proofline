// Package claims checks that a PR description matches the diff and that
// declared repro commands actually run.
package claims

import (
	"regexp"
	"strings"

	"github.com/zysen1/Proofline/internal/model"
)

// Input is everything the claims check can see. RunRepro may be nil when
// the caller only wants the static half (tests).
type Input struct {
	Body     string
	DiffText string
	Commits  string
	Changed  []string // repo-relative paths
	RunRepro func(command string) (exit int, output string, err error)
}

var (
	fileRe = regexp.MustCompile(`(?:^|[^\w./-])((?:[\w.-]+/)*[\w.-]+\.(?:py|pyi|js|jsx|mjs|cjs|ts|tsx|md|json|toml|txt|yml|yaml))`)
	btRe   = regexp.MustCompile("`([A-Za-z_][A-Za-z0-9_]*)`")
	fnRe   = regexp.MustCompile(`(?i)\b(?:function|def)\s+([A-Za-z_][A-Za-z0-9_]*)`)
	issRe  = regexp.MustCompile(`#(\d+)`)
	fence  = regexp.MustCompile("(?s)```(?:bash|sh|shell|console)?[ \\t]*\\n(.*?)```")
)

var allowPrefixes = []string{
	"pytest",
	"python -m pytest",
	"python3 -m pytest",
	"npx vitest",
	"npx jest",
	"npm test",
	"npm run test",
}

// Evaluate returns the claims check.
func Evaluate(in Input) model.Check {
	body := strings.TrimSpace(in.Body)
	files := unique(filesIn(body))
	funcs := unique(funcsIn(body))
	issues := unique(issuesIn(body))
	repros := unique(reprosIn(body))

	if len(files) == 0 && len(funcs) == 0 && len(issues) == 0 && len(repros) == 0 {
		return model.Check{
			ID:      "claims",
			Name:    "PR description matches the diff",
			Status:  model.StatusWarn,
			Summary: "The description has no checkable file, function, issue, or repro claim.",
		}
	}

	var missing []string
	var evidence []model.Evidence
	haystack := in.DiffText + "\n" + in.Commits

	for _, f := range files {
		if !fileInDiff(f, in.Changed) {
			missing = append(missing, "file "+f)
		}
	}
	for _, fn := range funcs {
		if !strings.Contains(haystack, fn) {
			missing = append(missing, "function "+fn)
		}
	}
	for _, iss := range issues {
		if !strings.Contains(haystack, "#"+iss) {
			missing = append(missing, "issue #"+iss)
		}
	}

	reproFail := 0
	for _, cmd := range repros {
		if !allowed(cmd) {
			reproFail++
			missing = append(missing, "repro not allowlisted: "+cmd)
			evidence = append(evidence, model.Evidence{Command: cmd, Detail: "rejected by the test-command allowlist"})
			continue
		}
		if in.RunRepro == nil {
			continue
		}
		code, out, err := in.RunRepro(cmd)
		ev := model.Evidence{Command: cmd, ExitCode: &code, Log: out}
		if err != nil {
			ev.Detail = err.Error()
			reproFail++
			missing = append(missing, "repro error: "+cmd)
		} else if code != 0 {
			reproFail++
			missing = append(missing, "repro failed: "+cmd)
		}
		evidence = append(evidence, ev)
	}

	warns := 0
	if len(repros) == 0 && (len(files)+len(funcs)+len(issues) > 0) {
		warns++
	}

	if len(missing) > 0 {
		return model.Check{
			ID:       "claims",
			Name:     "PR description matches the diff",
			Status:   model.StatusFail,
			Summary:  "Claims not backed by the diff or a passing repro: " + strings.Join(missing, "; ") + ".",
			Evidence: evidence,
		}
	}
	if warns > 0 {
		return model.Check{
			ID:       "claims",
			Name:     "PR description matches the diff",
			Status:   model.StatusWarn,
			Summary:  "Named files, functions, and issues appear in the diff. No repro command was declared.",
			Evidence: evidence,
		}
	}
	return model.Check{
		ID:       "claims",
		Name:     "PR description matches the diff",
		Status:   model.StatusPass,
		Summary:  "Files, functions, issues, and repro commands in the description match the change.",
		Evidence: evidence,
	}
}

func filesIn(body string) []string {
	var out []string
	for _, m := range fileRe.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	return out
}

func funcsIn(body string) []string {
	var out []string
	for _, m := range btRe.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	for _, m := range fnRe.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	return out
}

func issuesIn(body string) []string {
	var out []string
	for _, m := range issRe.FindAllStringSubmatch(body, -1) {
		out = append(out, m[1])
	}
	return out
}

func reprosIn(body string) []string {
	var out []string
	for _, m := range fence.FindAllStringSubmatch(body, -1) {
		for _, line := range strings.Split(m[1], "\n") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "$"))
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if looksLikeRepro(line) {
				out = append(out, line)
			}
		}
	}
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if looksLikeRepro(line) {
			out = append(out, line)
		}
	}
	return out
}

func looksLikeRepro(line string) bool {
	for _, p := range allowPrefixes {
		if line == p || strings.HasPrefix(line, p+" ") {
			return true
		}
	}
	return false
}

func allowed(cmd string) bool {
	if strings.ContainsAny(cmd, ";|&`$<>") {
		return false
	}
	return looksLikeRepro(strings.TrimSpace(cmd))
}

func fileInDiff(mentioned string, changed []string) bool {
	mentioned = strings.TrimPrefix(mentioned, "./")
	for _, c := range changed {
		c = strings.TrimPrefix(c, "./")
		if c == mentioned || strings.HasSuffix(c, "/"+mentioned) || strings.HasSuffix(mentioned, "/"+c) {
			return true
		}
	}
	return false
}

func unique(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
