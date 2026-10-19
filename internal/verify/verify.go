// Package verify runs the five deterministic checks and builds a receipt.
package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/zysen1/Proofline/internal/claims"
	"github.com/zysen1/Proofline/internal/gitx"
	"github.com/zysen1/Proofline/internal/lang"
	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/mutate"
	"github.com/zysen1/Proofline/internal/runner"
	"github.com/zysen1/Proofline/internal/scope"
)

const (
	KillPass       = 0.50
	KillWarn       = 0.25
	DefaultReruns  = 3
	DefaultMutants = 12
)

// Options controls one verification.
type Options struct {
	Repo       string
	Base       string
	Head       string
	PRBody     string
	Reruns     int
	MaxMutants int
	Parallel   int
	Timeout    time.Duration
	LLMSummary bool
}

// Verify compares head against base and returns a proof receipt.
func Verify(ctx context.Context, opt Options) (model.Receipt, error) {
	start := time.Now()
	if opt.Reruns <= 0 {
		opt.Reruns = DefaultReruns
	}
	if opt.MaxMutants <= 0 {
		opt.MaxMutants = DefaultMutants
	}
	if opt.Parallel <= 0 {
		opt.Parallel = 4
	}
	if opt.Timeout <= 0 {
		opt.Timeout = 120 * time.Second
	}
	repo := opt.Repo
	if repo == "" {
		repo = "."
	}
	abs, err := filepath.Abs(repo)
	if err != nil {
		return model.Receipt{}, err
	}
	top, err := gitx.Top(abs)
	if err != nil {
		return model.Receipt{}, err
	}
	baseSHA, err := gitx.RevParse(top, opt.Base)
	if err != nil {
		return model.Receipt{}, err
	}
	headSHA, err := gitx.RevParse(top, opt.Head)
	if err != nil {
		return model.Receipt{}, err
	}
	files, patch, err := gitx.ChangedFiles(top, opt.Base, opt.Head)
	if err != nil {
		return model.Receipt{}, err
	}
	commits, err := gitx.LogMessages(top, opt.Base, opt.Head)
	if err != nil {
		return model.Receipt{}, err
	}
	body := opt.PRBody
	if strings.TrimSpace(body) == "" {
		b, err := gitx.Show(top, headSHA, ".proofline/pr.md")
		if err != nil {
			return model.Receipt{}, err
		}
		body = string(b)
	}

	headSrc := map[string][]byte{}
	baseSrc := map[string][]byte{}
	for i := range files {
		h, err := gitx.Show(top, headSHA, files[i].Path)
		if err != nil {
			return model.Receipt{}, err
		}
		b, err := gitx.Show(top, baseSHA, files[i].Path)
		if err != nil {
			return model.Receipt{}, err
		}
		headSrc[files[i].Path] = h
		baseSrc[files[i].Path] = b
		files[i].Kind = lang.Classify(files[i].Path, h)
	}

	rec := model.Receipt{
		Tool:    "proofline",
		Version: model.Version,
		Repo:    top,
		Base:    opt.Base,
		Head:    opt.Head,
		BaseSHA: baseSHA,
		HeadSHA: headSHA,
	}

	if reason, na := notApplicable(files, body); na {
		rec.NotApplicableReason = reason
		rec.Verdict = model.VerdictNA
		rec.Badge = model.BadgeMarkdown(rec.Verdict)
		rec.Checks = []model.Check{{
			ID:      "applicability",
			Name:    "Applicability",
			Status:  model.StatusNA,
			Summary: reason + " Proofline does not treat this as a failed proof.",
		}}
		rec.DurationMS = time.Since(start).Milliseconds()
		rec.Score = model.EvidenceScore(rec)
		applySummary(&rec, opt.LLMSummary)
		return rec, nil
	}

	var selectors []lang.TestSelector
	var testBodies []string
	var testPaths []string
	for _, f := range files {
		if f.Kind != lang.KindTest || f.Status == "D" {
			continue
		}
		testPaths = append(testPaths, f.Path)
		sels := lang.ChangedTests(f.Path, baseSrc[f.Path], headSrc[f.Path])
		selectors = append(selectors, sels...)
		testBodies = append(testBodies, string(headSrc[f.Path]))
	}
	if len(selectors) == 0 {
		rec.NotApplicableReason = "No new or changed behavioral tests."
		rec.Verdict = model.VerdictNA
		rec.Badge = model.BadgeMarkdown(rec.Verdict)
		rec.Checks = []model.Check{{
			ID:      "applicability",
			Name:    "Applicability",
			Status:  model.StatusNA,
			Summary: "No new or changed behavioral tests. Proofline does not treat this as a failed proof.",
		}}
		rec.DurationMS = time.Since(start).Milliseconds()
		rec.Score = model.EvidenceScore(rec)
		applySummary(&rec, opt.LLMSummary)
		return rec, nil
	}

	headWT, cleanupH, err := gitx.AddWorktree(top, headSHA)
	if err != nil {
		return model.Receipt{}, err
	}
	defer cleanupH()
	baseWT, cleanupB, err := gitx.AddWorktree(top, baseSHA)
	if err != nil {
		return model.Receipt{}, err
	}
	defer cleanupB()
	if err := gitx.CheckoutPaths(baseWT, headSHA, testPaths); err != nil {
		return model.Receipt{}, err
	}

	linkNodeModules(top, headWT)
	linkNodeModules(top, baseWT)

	fbRuns, err := runN(ctx, baseWT, selectors, opt.Reruns, opt.Timeout)
	if err != nil {
		return model.Receipt{}, err
	}
	paRuns, err := runN(ctx, headWT, selectors, opt.Reruns, opt.Timeout)
	if err != nil {
		return model.Receipt{}, err
	}
	adapter := ""
	if len(paRuns) > 0 {
		adapter = paRuns[0].Adapter
	}
	rec.Adapter = adapter

	fb := checkFailBefore(fbRuns)
	pa := checkPassAfter(paRuns)
	fl := checkFlake(fbRuns, paRuns, opt.Reruns)

	var mutCheck model.Check
	if !anyPassed(paRuns) {
		mutCheck = model.Check{
			ID:      "mutation",
			Name:    "Tests constrain the changed lines",
			Status:  model.StatusNA,
			Summary: "Mutation was not run because the new tests do not pass on the head revision.",
		}
	} else {
		mutCheck = runMutation(ctx, headWT, files, headSrc, selectors, testBodies, opt, &rec)
	}

	var changedPaths []string
	for _, f := range files {
		changedPaths = append(changedPaths, f.Path)
	}
	claimsCheck := claims.Evaluate(claims.Input{
		Body:     body,
		DiffText: patch,
		Commits:  commits,
		Changed:  changedPaths,
		RunRepro: func(command string) (int, string, error) {
			return runner.Shell(ctx, headWT, command, opt.Timeout)
		},
	})

	var scopeFiles []scope.File
	for _, f := range files {
		scopeFiles = append(scopeFiles, scope.File{
			Path:       f.Path,
			Kind:       f.Kind,
			Added:      f.Added,
			Deleted:    f.Deleted,
			BaseSource: string(baseSrc[f.Path]),
			HeadSource: string(headSrc[f.Path]),
			Patch:      f.Patch,
		})
	}
	scopeCheck := scope.Evaluate(scopeFiles, body, testBodies)

	rec.Checks = []model.Check{fb, pa, mutCheck, claimsCheck, scopeCheck, fl}
	rec.Verdict = model.Decide(false, rec.Checks)
	rec.Badge = model.BadgeMarkdown(rec.Verdict)
	rec.DurationMS = time.Since(start).Milliseconds()
	rec.Score = model.EvidenceScore(rec)
	applySummary(&rec, opt.LLMSummary)
	return rec, nil
}

type timedRun struct {
	runner.Result
	Adapter string
}

func runN(ctx context.Context, dir string, sels []lang.TestSelector, n int, timeout time.Duration) ([]timedRun, error) {
	out := make([]timedRun, 0, n)
	for i := 0; i < n; i++ {
		res, adapter, err := runner.Run(ctx, dir, sels, timeout)
		if err != nil {
			return nil, err
		}
		out = append(out, timedRun{Result: res, Adapter: adapter})
	}
	return out, nil
}

func linkNodeModules(repo, wt string) {
	src := filepath.Join(repo, "node_modules")
	if _, err := os.Stat(src); err != nil {
		return
	}
	dst := filepath.Join(wt, "node_modules")
	if _, err := os.Lstat(dst); err == nil {
		return
	}
	_ = os.Symlink(src, dst)
}

func anyPassed(runs []timedRun) bool {
	for _, r := range runs {
		if r.Passed {
			return true
		}
	}
	return false
}

func allFailed(runs []timedRun) bool {
	if len(runs) == 0 {
		return false
	}
	for _, r := range runs {
		if r.Passed {
			return false
		}
	}
	return true
}

func allPassed(runs []timedRun) bool {
	if len(runs) == 0 {
		return false
	}
	for _, r := range runs {
		if !r.Passed {
			return false
		}
	}
	return true
}

func checkFailBefore(runs []timedRun) model.Check {
	ev := evidenceFrom(runs)
	dur := sumDur(runs)
	if allFailed(runs) {
		return model.Check{
			ID: "fail_before", Name: "New tests fail before the fix", Status: model.StatusPass,
			Summary:  fmt.Sprintf("With only the new tests applied to the base revision, the tests failed on %d/%d runs.", len(runs), len(runs)),
			Evidence: ev, DurationMS: dur,
		}
	}
	if allPassed(runs) {
		return model.Check{
			ID: "fail_before", Name: "New tests fail before the fix", Status: model.StatusFail,
			Summary:  fmt.Sprintf("The new tests already passed on the base revision (%d/%d runs). They do not require this change.", len(runs), len(runs)),
			Evidence: ev, DurationMS: dur,
		}
	}
	return model.Check{
		ID: "fail_before", Name: "New tests fail before the fix", Status: model.StatusFail,
		Summary:  "Fail-before runs disagreed, so the tests do not reliably fail without the fix.",
		Evidence: ev, DurationMS: dur,
	}
}

func checkPassAfter(runs []timedRun) model.Check {
	ev := evidenceFrom(runs)
	dur := sumDur(runs)
	if allPassed(runs) {
		return model.Check{
			ID: "pass_after", Name: "New tests pass after the fix", Status: model.StatusPass,
			Summary:  fmt.Sprintf("The new tests passed on the head revision (%d/%d runs).", len(runs), len(runs)),
			Evidence: ev, DurationMS: dur,
		}
	}
	summary := fmt.Sprintf("The new tests did not pass on the head revision (%d passing of %d runs).", passCount(runs), len(runs))
	if looksBroken(runs) {
		summary = "The head revision does not build or collect tests. " + summary
	}
	return model.Check{
		ID: "pass_after", Name: "New tests pass after the fix", Status: model.StatusFail,
		Summary: summary, Evidence: ev, DurationMS: dur,
	}
}

func looksBroken(runs []timedRun) bool {
	for _, r := range runs {
		s := r.Output
		if strings.Contains(s, "SyntaxError") || strings.Contains(s, "IndentationError") ||
			strings.Contains(s, "Syntax error") || strings.Contains(s, "ERROR: ") ||
			strings.Contains(s, "does not provide an export") || strings.Contains(s, "Unexpected token") {
			return true
		}
	}
	return false
}

func checkFlake(fb, pa []timedRun, n int) model.Check {
	if n < 2 {
		return model.Check{
			ID: "flake", Name: "Tests are stable across reruns", Status: model.StatusNA,
			Summary: "Fewer than 2 reruns were requested, so flakiness was not judged.",
		}
	}
	fbStable := sameOutcome(fb)
	paStable := sameOutcome(pa)
	ev := append(evidenceFrom(fb), evidenceFrom(pa)...)
	if fbStable && paStable {
		return model.Check{
			ID: "flake", Name: "Tests are stable across reruns", Status: model.StatusPass,
			Summary:  fmt.Sprintf("Fail-before and pass-after each repeated %d times with a stable result.", n),
			Evidence: ev,
		}
	}
	return model.Check{
		ID: "flake", Name: "Tests are stable across reruns", Status: model.StatusFail,
		Summary:  fmt.Sprintf("Outcomes disagreed across %d reruns.", n),
		Evidence: ev,
	}
}

func sameOutcome(runs []timedRun) bool {
	if len(runs) == 0 {
		return true
	}
	first := runs[0].Passed
	for _, r := range runs[1:] {
		if r.Passed != first {
			return false
		}
	}
	return true
}

func passCount(runs []timedRun) int {
	n := 0
	for _, r := range runs {
		if r.Passed {
			n++
		}
	}
	return n
}

func sumDur(runs []timedRun) int64 {
	var d time.Duration
	for _, r := range runs {
		d += r.Duration
	}
	return d.Milliseconds()
}

func evidenceFrom(runs []timedRun) []model.Evidence {
	var ev []model.Evidence
	for _, r := range runs {
		code := r.ExitCode
		ev = append(ev, model.Evidence{
			Command:  r.Command,
			ExitCode: &code,
			Log:      r.Output,
		})
	}
	return ev
}

func runMutation(ctx context.Context, headWT string, files []model.ChangedFile, headSrc map[string][]byte, sels []lang.TestSelector, testBodies []string, opt Options, rec *model.Receipt) model.Check {
	var sites []mutate.FileSites
	for _, f := range files {
		if f.Kind != lang.KindSource {
			continue
		}
		// Unrelated files are scored by the scope check. Mutation asks whether
		// the new tests constrain the source they actually touch.
		if !referencedByTests(f.Path, testBodies) {
			continue
		}
		id, ok := lang.Language(f.Path)
		if !ok {
			continue
		}
		lines := map[int]bool{}
		for _, ln := range f.NewLines {
			lines[ln] = true
		}
		sites = append(sites, mutate.FileSites{
			Path: f.Path, Lang: id, Source: headSrc[f.Path], ChangedLines: lines,
		})
	}
	mutants := mutate.Generate(sites, opt.MaxMutants)
	if len(mutants) == 0 {
		return model.Check{
			ID: "mutation", Name: "Tests constrain the changed lines", Status: model.StatusNA,
			Summary: "No compare, boundary, delete, or arithmetic sites were found on the changed source lines.",
		}
	}
	type outcome struct {
		m      mutate.Mutant
		killed bool
		log    string
	}
	results := make([]outcome, len(mutants))
	sem := make(chan struct{}, opt.Parallel)
	var wg sync.WaitGroup
	for i, m := range mutants {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, m mutate.Mutant) {
			defer wg.Done()
			defer func() { <-sem }()
			killed, log := evalMutant(ctx, headWT, m, sels, opt.Timeout)
			results[i] = outcome{m: m, killed: killed, log: log}
		}(i, m)
	}
	wg.Wait()

	byOp := map[string]int{}
	killedOp := map[string]int{}
	killed := 0
	var ev []model.Evidence
	for _, r := range results {
		byOp[r.m.Operator]++
		status := "survived"
		if r.killed {
			killed++
			killedOp[r.m.Operator]++
			status = "killed"
		}
		ev = append(ev, model.Evidence{Detail: r.m.Summary + " — " + status, Log: r.log})
	}
	rec.Mutants = len(mutants)
	rec.Killed = killed
	rec.MutantsByOp = byOp
	rec.KilledByOp = killedOp
	rec.KillRate = float64(killed) / float64(len(mutants))

	var parts []string
	for _, op := range mutate.Operators {
		if byOp[op] == 0 {
			parts = append(parts, op+" 0 sites")
			continue
		}
		parts = append(parts, fmt.Sprintf("%s %d/%d", op, killedOp[op], byOp[op]))
	}
	summary := fmt.Sprintf("Killed %s. Operators: %s.", model.FormatRate(killed, len(mutants)), strings.Join(parts, ", "))
	st := model.StatusFail
	switch {
	case rec.KillRate >= KillPass:
		st = model.StatusPass
	case rec.KillRate >= KillWarn:
		st = model.StatusWarn
	}
	if len(ev) > 16 {
		ev = ev[:16]
	}
	return model.Check{
		ID: "mutation", Name: "Tests constrain the changed lines", Status: st,
		Summary: summary, Evidence: ev,
	}
}

func referencedByTests(path string, tests []string) bool {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if len(base) < 2 {
		return false
	}
	for _, t := range tests {
		if strings.Contains(t, base) {
			return true
		}
	}
	return false
}

func evalMutant(ctx context.Context, headWT string, m mutate.Mutant, sels []lang.TestSelector, timeout time.Duration) (bool, string) {
	dir, err := copyTree(headWT)
	if err != nil {
		return false, err.Error()
	}
	defer os.RemoveAll(dir)
	dest := filepath.Join(dir, filepath.FromSlash(m.File))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return false, err.Error()
	}
	if err := os.WriteFile(dest, m.Source, 0o644); err != nil {
		return false, err.Error()
	}
	res, _, err := runner.Run(ctx, dir, sels, timeout)
	if err != nil {
		return false, err.Error()
	}
	// A mutant is killed when the new tests stop passing.
	return !res.Passed, res.Output
}

func copyTree(src string) (string, error) {
	dst, err := os.MkdirTemp("", "proofline-mut-")
	if err != nil {
		return "", err
	}
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if info.Name() == ".git" || info.Name() == "__pycache__" || info.Name() == ".pytest_cache" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		os.RemoveAll(dst)
		return "", err
	}
	return dst, nil
}

func notApplicable(files []model.ChangedFile, body string) (string, bool) {
	if len(files) == 0 {
		return "The diff is empty.", true
	}
	lower := strings.ToLower(body)
	if strings.Contains(lower, "proofline: not-applicable") {
		return "The description marks this change proofline: not-applicable.", true
	}
	var sawTest, sawSource, sawDoc, sawBench, sawOther bool
	for _, f := range files {
		switch f.Kind {
		case lang.KindTest:
			sawTest = true
		case lang.KindSource:
			sawSource = true
		case lang.KindDoc:
			sawDoc = true
		case lang.KindBench:
			sawBench = true
		case lang.KindMeta, lang.KindConfig, lang.KindManifest:
		default:
			sawOther = true
		}
	}
	if sawTest || sawSource || sawOther {
		return "", false
	}
	if sawBench {
		return "Performance or benchmark-only change.", true
	}
	if sawDoc {
		return "Documentation-only change.", true
	}
	return "", false
}

func applySummary(rec *model.Receipt, llm bool) {
	if !llm {
		return
	}
	rec.Summary = localSummary(*rec)
	rec.SummarySource = "local"
	cmd := os.Getenv("PROOFLINE_LLM_CMD")
	if cmd == "" {
		return
	}
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = strings.NewReader(rec.Summary)
	out, err := c.Output()
	if err != nil {
		rec.Summary = rec.Summary + "\n\nLocal summary kept; PROOFLINE_LLM_CMD failed: " + err.Error()
		return
	}
	text := strings.TrimSpace(string(out))
	if text != "" {
		rec.Summary = text
		rec.SummarySource = "PROOFLINE_LLM_CMD"
	}
}

func localSummary(r model.Receipt) string {
	var bits []string
	for _, c := range r.Checks {
		bits = append(bits, c.ID+"="+string(c.Status))
	}
	return fmt.Sprintf("Verdict %s from deterministic checks (%s). This paragraph is a local rendering of those checks and does not change the verdict.", r.Verdict, strings.Join(bits, ", "))
}
