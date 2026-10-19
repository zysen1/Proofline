package dataset

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/verify"
)

func TestOneOfEachClass(t *testing.T) {
	for _, class := range Classes {
		class := class
		t.Run(class, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "repo")
			sp := specFor(class, 0)
			if err := Materialize(dir, sp); err != nil {
				t.Fatal(err)
			}
			rec, err := verify.Verify(context.Background(), verify.Options{
				Repo: dir, Base: "main", Head: "pr",
				Reruns: 2, MaxMutants: 8, Timeout: 60 * time.Second,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s verdict=%s kill=%.2f", class, rec.Verdict, rec.KillRate)
			for _, c := range rec.Checks {
				t.Logf("  %s %s %s", c.ID, c.Status, c.Summary)
			}
			if class == ClassReal {
				if rec.Verdict != model.VerdictVerified {
					t.Fatalf("real fix verdict %s", rec.Verdict)
				}
				return
			}
			if !detected(class, rec) {
				t.Fatalf("class %s not detected (verdict %s)", class, rec.Verdict)
			}
		})
	}
}
