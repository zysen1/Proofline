// Package demo materializes the toy shop and writes the offline report.
package demo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/zysen1/Proofline/internal/gitx"
	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/rank"
	"github.com/zysen1/Proofline/internal/receipt"
	"github.com/zysen1/Proofline/internal/report"
	"github.com/zysen1/Proofline/internal/verify"
)

// Options points at the fixture directory and the docs output.
type Options struct {
	Toy    string
	Out    string
	Reruns int
}

type manifest struct {
	Branches []struct {
		Name  string `json:"name"`
		Title string `json:"title"`
		Rank  bool   `json:"rank"`
	} `json:"branches"`
}

// Run builds the toy git repo, verifies every preset branch, ranks the
// parallel-agent subset, and writes Markdown, HTML, and a PNG.
func Run(ctx context.Context, opt Options) error {
	if opt.Reruns <= 0 {
		opt.Reruns = 3
	}
	man, err := loadManifest(opt.Toy)
	if err != nil {
		return err
	}
	repo := filepath.Join(opt.Toy, ".repo")
	if err := materialize(opt.Toy, repo, man); err != nil {
		return err
	}
	var recs []model.Receipt
	byName := map[string]model.Receipt{}
	for _, br := range man.Branches {
		rec, err := verify.Verify(ctx, verify.Options{
			Repo: repo, Base: "main", Head: br.Name,
			Reruns: opt.Reruns, MaxMutants: 12, Timeout: 60 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", br.Name, err)
		}
		recs = append(recs, rec)
		byName[br.Name] = rec
		fmt.Fprintf(os.Stderr, "demo %-22s %s\n", br.Name, rec.Verdict)
	}
	var rankRecs []model.Receipt
	for _, br := range man.Branches {
		if br.Rank {
			rankRecs = append(rankRecs, byName[br.Name])
		}
	}
	rows := rank.Order(rankRecs)
	flake, err := flaky(ctx, opt.Reruns)
	if err != nil {
		return err
	}

	if err := assertStory(recs, rows); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(opt.Out, "receipts"), 0o755); err != nil {
		return err
	}
	for _, rec := range recs {
		path := filepath.Join(opt.Out, "receipts", rec.Head+".md")
		if err := os.WriteFile(path, []byte(receipt.Markdown(rec)), 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(opt.Out, "rank.md"), []byte(rank.Markdown(rows)), 0o644); err != nil {
		return err
	}
	htmlPath := filepath.Join(opt.Out, "demo-report.html")
	page := report.Page{
		Title:    "Proofline — don't trust the write-up",
		Subtitle: "One real fix and five slop branches. Only the branch whose new tests fail before the change, pass after it, and kill mutants is Verified.",
		Receipts: recs,
		Rank:     rows,
		Extra:    &flake,
		Foot:     "Offline toy shop. Badges come from pytest, not from the pull request text. Replaying merged public pull requests is not part of this demo.",
	}
	if err := report.WriteHTML(htmlPath, page); err != nil {
		return err
	}
	lines := []string{"Proofline demo"}
	for _, rec := range recs {
		lines = append(lines, fmt.Sprintf("%-22s %s", rec.Head, rec.Verdict))
	}
	lines = append(lines, "")
	for _, row := range rows {
		lines = append(lines, fmt.Sprintf("#%d %-22s %s score %d", row.Rank, row.Receipt.Head, row.Receipt.Verdict, row.Receipt.Score))
	}
	lines = append(lines, "", fmt.Sprintf("flake %s", flake.Verdict))
	png := filepath.Join(opt.Out, "demo.png")
	if err := report.Screenshot(htmlPath, png, lines); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", htmlPath)
	return nil
}

func assertStory(recs []model.Receipt, rows []rank.Row) error {
	verified := 0
	for _, rec := range recs {
		if rec.Head == "real-fix" {
			if rec.Verdict != model.VerdictVerified {
				return fmt.Errorf("real-fix verdict = %s, want Verified", rec.Verdict)
			}
		} else if rec.Verdict != model.VerdictUnverified {
			return fmt.Errorf("%s verdict = %s, want Unverified", rec.Head, rec.Verdict)
		}
		if rec.Verdict == model.VerdictVerified {
			verified++
		}
	}
	if verified != 1 {
		return fmt.Errorf("verified count = %d, want 1", verified)
	}
	if len(rows) == 0 || rows[0].Receipt.Head != "real-fix" {
		return fmt.Errorf("rank did not put real-fix first")
	}
	if rows[len(rows)-1].Receipt.Head != "slop-tautology" {
		return fmt.Errorf("rank last = %s, want slop-tautology", rows[len(rows)-1].Receipt.Head)
	}
	return nil
}

func loadManifest(toy string) (manifest, error) {
	b, err := os.ReadFile(filepath.Join(toy, "manifest.json"))
	if err != nil {
		return manifest{}, err
	}
	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return manifest{}, err
	}
	return m, nil
}

func materialize(toy, dest string, man manifest) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	if err := gitx.InitRepo(dest); err != nil {
		return err
	}
	if err := copyTree(filepath.Join(toy, "branches", "main"), dest); err != nil {
		return err
	}
	if err := gitx.CommitAll(dest, "buggy discount: apply_discount is not implemented"); err != nil {
		return err
	}
	for _, br := range man.Branches {
		if err := gitx.CheckoutNew(dest, br.Name); err != nil {
			return err
		}
		if err := gitx.CleanWorktree(dest); err != nil {
			return err
		}
		if err := copyTree(filepath.Join(toy, "branches", br.Name), dest); err != nil {
			return err
		}
		msg := br.Title + "\n"
		if err := gitx.CommitAll(dest, msg); err != nil {
			return err
		}
	}
	return gitx.Checkout(dest, "main")
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
}

func flaky(ctx context.Context, reruns int) (model.Receipt, error) {
	dir, err := os.MkdirTemp("", "proofline-flake-")
	if err != nil {
		return model.Receipt{}, err
	}
	defer os.RemoveAll(dir)
	if err := gitx.InitRepo(dir); err != nil {
		return model.Receipt{}, err
	}
	base := "def value():\n    return 1\n"
	if err := os.WriteFile(filepath.Join(dir, "mod.py"), []byte(base), 0o644); err != nil {
		return model.Receipt{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "test_mod.py"), []byte("def test_placeholder():\n    assert True\n"), 0o644); err != nil {
		return model.Receipt{}, err
	}
	if err := gitx.CommitAll(dir, "base"); err != nil {
		return model.Receipt{}, err
	}
	if err := gitx.CheckoutNew(dir, "flaky-counter"); err != nil {
		return model.Receipt{}, err
	}
	head := "def value():\n    # touched so the branch has a source diff\n    return 1\n"
	test := `import os
from mod import value

def test_value_alternates():
    value()
    path = "runs.txt"
    n = 0
    if os.path.exists(path):
        raw = open(path).read().strip() or "0"
        n = int(raw)
    open(path, "w").write(str(n + 1))
    assert n % 2 == 0
`
	if err := os.WriteFile(filepath.Join(dir, "mod.py"), []byte(head), 0o644); err != nil {
		return model.Receipt{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "test_mod.py"), []byte(test), 0o644); err != nil {
		return model.Receipt{}, err
	}
	body := "`value` in `mod.py`\n\nRepro:\npytest test_mod.py::test_value_alternates\n"
	if err := os.MkdirAll(filepath.Join(dir, ".proofline"), 0o755); err != nil {
		return model.Receipt{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, ".proofline", "pr.md"), []byte(body), 0o644); err != nil {
		return model.Receipt{}, err
	}
	if err := gitx.CommitAll(dir, "flaky counter"); err != nil {
		return model.Receipt{}, err
	}
	rec, err := verify.Verify(ctx, verify.Options{
		Repo: dir, Base: "main", Head: "flaky-counter",
		Reruns: reruns, MaxMutants: 4, Timeout: 60 * time.Second,
	})
	if err != nil {
		return model.Receipt{}, err
	}
	ch, ok := rec.Find("flake")
	if !ok || ch.Status != model.StatusFail {
		return rec, fmt.Errorf("flake fixture status = %s, want fail", ch.Status)
	}
	return rec, nil
}
