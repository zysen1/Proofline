package scope

import "testing"

func TestUnrelatedAndDependency(t *testing.T) {
	ch := Evaluate([]File{
		{Path: "mod.py", Kind: "source", Added: 8, Deleted: 1, HeadSource: "def discount(n):\n    return n\n", BaseSource: "def discount(n):\n    raise NotImplementedError\n"},
		{Path: "test_mod.py", Kind: "test", Added: 10, Deleted: 0, HeadSource: "from mod import discount\n"},
		{Path: "extras/noise.py", Kind: "source", Added: 80, Deleted: 0, HeadSource: "def noise(v):\n    return v\n"},
		{Path: "requirements.txt", Kind: "manifest", Added: 1, Deleted: 0, Patch: "+requests==2.32.3\n", HeadSource: "requests==2.32.3\n"},
	}, "Fixes #101\n\n`discount` in `mod.py`\n", []string{"from mod import discount\n"})
	if ch.Status != "fail" {
		t.Fatalf("status %s: %s", ch.Status, ch.Summary)
	}
}

func TestTightFixPasses(t *testing.T) {
	ch := Evaluate([]File{
		{Path: "mod.py", Kind: "source", Added: 6, Deleted: 1, HeadSource: "def discount(n):\n    return n * 2 + 1\n", BaseSource: "def discount(n):\n    raise NotImplementedError\n"},
		{Path: "test_mod.py", Kind: "test", Added: 8, Deleted: 0, HeadSource: "from mod import discount\n"},
		{Path: ".proofline/pr.md", Kind: "meta", Added: 4, Deleted: 0},
	}, "`discount` in `mod.py`\n", []string{"from mod import discount\n"})
	if ch.Status != "pass" {
		t.Fatalf("status %s: %s", ch.Status, ch.Summary)
	}
}
