// Package gitx is the small git surface Proofline needs: revisions, diffs, and worktrees.
package gitx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zysen1/Proofline/internal/diffx"
	"github.com/zysen1/Proofline/internal/model"
)

// RevParse resolves a ref to a full SHA.
func RevParse(repo, ref string) (string, error) {
	out, err := git(repo, "rev-parse", ref)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", ref, err)
	}
	return strings.TrimSpace(out), nil
}

// ChangedFiles returns the three-dot diff base...head.
func ChangedFiles(repo, base, head string) ([]model.ChangedFile, string, error) {
	spec := base + "..." + head
	nameStatus, err := git(repo, "diff", "--name-status", "-M", spec)
	if err != nil {
		return nil, "", err
	}
	numstat, err := git(repo, "diff", "--numstat", "-M", spec)
	if err != nil {
		return nil, "", err
	}
	patch, err := git(repo, "diff", "--no-color", "-U0", spec)
	if err != nil {
		return nil, "", err
	}
	lines := diffx.ParseUnified(patch)
	counts := parseNumstat(numstat)
	var files []model.ChangedFile
	for _, row := range strings.Split(strings.TrimSpace(nameStatus), "\n") {
		if strings.TrimSpace(row) == "" {
			continue
		}
		fields := strings.Split(row, "\t")
		if len(fields) < 2 {
			continue
		}
		status := fields[0]
		path := fields[len(fields)-1]
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			status = status[:1]
		}
		c := counts[path]
		files = append(files, model.ChangedFile{
			Path:     path,
			Status:   status,
			Added:    c.added,
			Deleted:  c.deleted,
			NewLines: lines[path],
			Patch:    filePatch(patch, path),
		})
	}
	return files, patch, nil
}

type count struct{ added, deleted int }

func parseNumstat(s string) map[string]count {
	out := map[string]count{}
	for _, row := range strings.Split(strings.TrimSpace(s), "\n") {
		if row == "" {
			continue
		}
		fields := strings.Split(row, "\t")
		if len(fields) < 3 {
			continue
		}
		_, path := diffx.NumstatPath(strings.Join(fields[2:], "\t"))
		a, _ := strconv.Atoi(fields[0])
		d, _ := strconv.Atoi(fields[1])
		out[path] = count{added: a, deleted: d}
	}
	return out
}

func filePatch(patch, path string) string {
	var b strings.Builder
	keep := false
	for _, line := range strings.Split(patch, "\n") {
		if strings.HasPrefix(line, "diff --git ") {
			keep = headerPath(line) == path
			if keep {
				b.WriteString(line)
				b.WriteByte('\n')
			}
			continue
		}
		if keep {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func headerPath(line string) string {
	rest := strings.TrimPrefix(line, "diff --git ")
	i := strings.LastIndex(rest, " b/")
	if i < 0 {
		return ""
	}
	return rest[i+3:]
}

// Show reads a blob. A missing path returns empty bytes and a nil error
// when the revision exists but the path does not.
func Show(repo, rev, path string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repo, "show", rev+":"+path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if strings.Contains(msg, "does not exist") || strings.Contains(msg, "exists on disk") || strings.Contains(msg, "Path") {
			return nil, nil
		}
		return nil, fmt.Errorf("git show %s:%s: %w (%s)", rev, path, err, strings.TrimSpace(msg))
	}
	return stdout.Bytes(), nil
}

// LogMessages returns commit messages reachable from head but not from base.
func LogMessages(repo, base, head string) (string, error) {
	out, err := git(repo, "log", "--format=%B%n---", base+".."+head)
	if err != nil {
		return "", err
	}
	return out, nil
}

// AddWorktree checks out sha in a detached worktree and returns a cleanup func.
func AddWorktree(repo, sha string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "proofline-wt-")
	if err != nil {
		return "", nil, err
	}
	if _, err := git(repo, "worktree", "add", "--detach", "--force", dir, sha); err != nil {
		os.RemoveAll(dir)
		return "", nil, err
	}
	cleanup := func() {
		_, _ = git(repo, "worktree", "remove", "--force", dir)
		os.RemoveAll(dir)
	}
	return dir, cleanup, nil
}

// CheckoutPaths replaces paths in a worktree with the blobs from sha.
func CheckoutPaths(dir, sha string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	args := append([]string{"-C", dir, "checkout", sha, "--"}, paths...)
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("checkout paths: %w (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func git(repo string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w (%s)", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

// InitRepo creates a repository with a local identity so commits work offline.
func InitRepo(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if _, err := git(dir, "init", "-b", "main"); err != nil {
		return err
	}
	if _, err := git(dir, "config", "user.email", "demo@proofline.local"); err != nil {
		return err
	}
	if _, err := git(dir, "config", "user.name", "Proofline Demo"); err != nil {
		return err
	}
	return nil
}

// CommitAll stages the worktree and commits. An empty tree is an error.
func CommitAll(dir, message string) error {
	if _, err := git(dir, "add", "-A"); err != nil {
		return err
	}
	cmd := exec.Command("git", "-C", dir, "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit: %w (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// CheckoutNew creates or resets a branch to the current HEAD.
func CheckoutNew(dir, branch string) error {
	_, err := git(dir, "checkout", "-B", branch)
	return err
}

// Checkout switches to an existing branch without moving it.
func Checkout(dir, branch string) error {
	_, err := git(dir, "checkout", branch)
	return err
}

// CleanWorktree removes every file except .git.
func CleanWorktree(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// Top returns the absolute repository root for dir.
func Top(dir string) (string, error) {
	out, err := git(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
