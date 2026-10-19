// Package runner executes the project's own test tool: pytest, vitest, or jest.
package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/zysen1/Proofline/internal/lang"
)

// Result is one test process.
type Result struct {
	Command  string
	ExitCode int
	Output   string
	Passed   bool
	Duration time.Duration
	Err      error
}

// Run executes every selector. Python and JS selectors are split across adapters.
// The returned name is "pytest", "vitest", "jest", or a join when both run.
func Run(ctx context.Context, repo string, selectors []lang.TestSelector, timeout time.Duration) (Result, string, error) {
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	var py, js []lang.TestSelector
	for _, s := range selectors {
		id, _ := lang.Language(s.File)
		if id == lang.Python {
			py = append(py, s)
		} else {
			js = append(js, s)
		}
	}
	var results []Result
	var names []string
	if len(py) > 0 {
		if err := pytestAvailable(); err != nil {
			return Result{}, "", err
		}
		results = append(results, runPytest(ctx, repo, py, timeout))
		names = append(names, "pytest")
	}
	if len(js) > 0 {
		kind, err := detectJS(repo)
		if err != nil {
			return Result{}, "", err
		}
		results = append(results, runJS(ctx, repo, kind, js, timeout))
		names = append(names, kind)
	}
	if len(results) == 0 {
		return Result{}, "", errors.New("no test selectors to run")
	}
	return merge(results), strings.Join(names, "+"), nil
}

// Shell runs an allowlisted repro command in dir.
func Shell(ctx context.Context, dir, command string, timeout time.Duration) (int, string, error) {
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	r := run(ctx, dir, timeout, "sh", "-c", command)
	if r.Err != nil && r.ExitCode == -1 {
		return -1, r.Output, r.Err
	}
	return r.ExitCode, r.Output, nil
}

func merge(results []Result) Result {
	out := results[0]
	var cmds []string
	var logs []string
	passed := true
	for _, r := range results {
		cmds = append(cmds, r.Command)
		logs = append(logs, r.Output)
		out.Duration += r.Duration
		if !r.Passed {
			passed = false
			out.ExitCode = r.ExitCode
			out.Err = r.Err
		}
	}
	out.Command = strings.Join(cmds, " && ")
	out.Output = strings.Join(logs, "\n")
	out.Passed = passed
	if passed {
		out.ExitCode = 0
		out.Err = nil
	}
	return out
}

func runPytest(ctx context.Context, repo string, sels []lang.TestSelector, timeout time.Duration) Result {
	args := []string{"-m", "pytest", "-q", "--tb=line", "-p", "no:cacheprovider"}
	args = append(args, pytestArgs(sels)...)
	r := run(ctx, repo, timeout, "python3", args...)
	r.Command = "python3 " + strings.Join(args, " ")
	return r
}

func pytestArgs(sels []lang.TestSelector) []string {
	var args []string
	fileOnly := false
	for _, s := range sels {
		if s.Name == "" {
			fileOnly = true
			break
		}
	}
	if fileOnly {
		seen := map[string]bool{}
		for _, s := range sels {
			if !seen[s.File] {
				args = append(args, filepath.ToSlash(s.File))
				seen[s.File] = true
			}
		}
		return args
	}
	for _, s := range sels {
		args = append(args, s.Name)
	}
	return args
}

func runJS(ctx context.Context, repo, kind string, sels []lang.TestSelector, timeout time.Duration) Result {
	files := uniqueFiles(sels)
	bin, baseArgs := jsCommand(repo, kind)
	args := append(baseArgs, files...)
	r := run(ctx, repo, timeout, bin, args...)
	r.Command = bin + " " + strings.Join(args, " ")
	return r
}

func jsCommand(repo, kind string) (string, []string) {
	switch kind {
	case "vitest":
		if bin := localBin(repo, "vitest"); bin != "" {
			return bin, []string{"run", "--reporter=dot"}
		}
		return "npx", []string{"--yes", "vitest", "run", "--reporter=dot"}
	default:
		if bin := localBin(repo, "jest"); bin != "" {
			return bin, []string{"--watchAll=false", "--ci", "--colors=false", "--runInBand"}
		}
		return "npx", []string{"jest", "--watchAll=false", "--ci", "--colors=false", "--runInBand"}
	}
}

func localBin(repo, name string) string {
	p := filepath.Join(repo, "node_modules", ".bin", name)
	if st, err := os.Stat(p); err == nil && !st.IsDir() {
		return p
	}
	return ""
}

func uniqueFiles(sels []lang.TestSelector) []string {
	var out []string
	seen := map[string]bool{}
	for _, s := range sels {
		if !seen[s.File] {
			out = append(out, filepath.ToSlash(s.File))
			seen[s.File] = true
		}
	}
	return out
}

func detectJS(repo string) (string, error) {
	raw, err := os.ReadFile(filepath.Join(repo, "package.json"))
	if err != nil {
		return "", fmt.Errorf("package.json is required for vitest or jest: %w", err)
	}
	var doc struct {
		Scripts         map[string]string `json:"scripts"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return "", fmt.Errorf("parse package.json: %w", err)
	}
	script := doc.Scripts["test"]
	has := func(name string) bool {
		_, a := doc.Dependencies[name]
		_, b := doc.DevDependencies[name]
		return a || b || strings.Contains(script, name)
	}
	if strings.Contains(script, "vitest") || (has("vitest") && !strings.Contains(script, "jest")) {
		return "vitest", nil
	}
	if strings.Contains(script, "jest") || has("jest") {
		return "jest", nil
	}
	if localBin(repo, "vitest") != "" {
		return "vitest", nil
	}
	if localBin(repo, "jest") != "" {
		return "jest", nil
	}
	return "", errors.New("neither vitest nor jest is declared in package.json")
}

func pytestAvailable() error {
	cmd := exec.Command("python3", "-c", "import pytest")
	if err := cmd.Run(); err != nil {
		return errors.New("pytest is not importable (python3 -c 'import pytest'). Install it with pip install pytest")
	}
	return nil
}

func run(ctx context.Context, dir string, timeout time.Duration, name string, args ...string) Result {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"PYTHONPATH="+dir,
		"PYTHONDONTWRITEBYTECODE=1",
		"CI=1",
	)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	start := time.Now()
	err := cmd.Run()
	res := Result{
		Command:  name + " " + strings.Join(args, " "),
		Output:   truncate(buf.String(), 4000),
		Duration: time.Since(start),
	}
	if cctx.Err() == context.DeadlineExceeded {
		res.ExitCode = -1
		res.Err = fmt.Errorf("%s timed out after %s", name, timeout)
		res.Output = truncate(res.Output+"\n[timed out]", 4000)
		return res
	}
	if err == nil {
		res.Passed = true
		return res
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
		return res
	}
	res.ExitCode = -1
	res.Err = err
	return res
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…[truncated]"
}
