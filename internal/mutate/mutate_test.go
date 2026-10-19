package mutate

import (
	"strings"
	"testing"

	"github.com/zysen1/Proofline/internal/lang"
)

func TestFourOperatorsPython(t *testing.T) {
	src := []byte(`def f(n):
    if n < 0:
        raise ValueError("bad")
    return n * 2 + 1
`)
	lines := map[int]bool{1: true, 2: true, 3: true, 4: true}
	got := Generate([]FileSites{{
		Path: "mod.py", Lang: lang.Python, Source: src, ChangedLines: lines,
	}}, 20)
	seen := map[string]bool{}
	for _, m := range got {
		seen[m.Operator] = true
		if string(m.Source) == string(src) {
			t.Fatalf("mutant did not change source: %s", m.Summary)
		}
	}
	for _, op := range Operators {
		if !seen[op] {
			t.Errorf("missing operator %s", op)
		}
	}
}

func TestFourOperatorsJavaScript(t *testing.T) {
	src := []byte(`function f(n) {
  if (n < 0) {
    throw new Error("bad");
  }
  return n * 2 + 1;
}
`)
	lines := map[int]bool{}
	for i := 1; i <= 6; i++ {
		lines[i] = true
	}
	got := Generate([]FileSites{{
		Path: "pricing.js", Lang: lang.JavaScript, Source: src, ChangedLines: lines,
	}}, 20)
	seen := map[string]bool{}
	for _, m := range got {
		seen[m.Operator] = true
	}
	for _, op := range Operators {
		if !seen[op] {
			t.Errorf("missing operator %s among %d mutants", op, len(got))
		}
	}
	joined := ""
	for _, m := range got {
		joined += m.Summary + "\n"
	}
	if !strings.Contains(joined, "compare") || !strings.Contains(joined, "arithmetic") {
		t.Fatalf("summaries missing operators:\n%s", joined)
	}
}
