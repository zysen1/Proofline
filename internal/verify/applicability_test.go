package verify

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/zysen1/Proofline/internal/gitx"
	"github.com/zysen1/Proofline/internal/model"
)

func TestDocsOnlyIsNotAFailure(t *testing.T) {
	dir := t.TempDir()
	if err := gitx.InitRepo(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Shop\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "base"); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CheckoutNew(dir, "docs"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Shop\n\nMore words.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "docs only"); err != nil {
		t.Fatal(err)
	}
	rec, err := Verify(context.Background(), Options{Repo: dir, Base: "main", Head: "docs", Reruns: 1})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != model.VerdictNA {
		t.Fatalf("verdict %s, want Not applicable (%s)", rec.Verdict, rec.NotApplicableReason)
	}
}

func TestBenchOnlyIsNotAFailure(t *testing.T) {
	dir := t.TempDir()
	if err := gitx.InitRepo(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bench", "bench_add.py"), []byte("def bench_add():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "base"); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CheckoutNew(dir, "perf"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bench", "bench_add.py"), []byte("def bench_add():\n    return 1 + 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "perf only"); err != nil {
		t.Fatal(err)
	}
	rec, err := Verify(context.Background(), Options{Repo: dir, Base: "main", Head: "perf", Reruns: 1, Timeout: 30 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != model.VerdictNA {
		t.Fatalf("verdict %s (%s)", rec.Verdict, rec.NotApplicableReason)
	}
}

func TestSummaryDoesNotChangeVerdict(t *testing.T) {
	dir := t.TempDir()
	if err := gitx.InitRepo(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# A\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "base"); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CheckoutNew(dir, "docs"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# B\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitx.CommitAll(dir, "docs"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PROOFLINE_LLM_CMD", "printf 'Verified\\n'")
	rec, err := Verify(context.Background(), Options{
		Repo: dir, Base: "main", Head: "docs", Reruns: 1, LLMSummary: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rec.Verdict != model.VerdictNA {
		t.Fatalf("verdict changed to %s", rec.Verdict)
	}
	if rec.SummarySource != "PROOFLINE_LLM_CMD" {
		t.Fatalf("summary source %q", rec.SummarySource)
	}
	if rec.Summary != "Verified" {
		t.Fatalf("summary %q", rec.Summary)
	}
}
