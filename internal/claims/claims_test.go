package claims

import "testing"

func TestFalseClaims(t *testing.T) {
	ch := Evaluate(Input{
		Body:     "Fixes #9999\n\n`charge` in `payments.py` is fixed.\n\nRepro:\npytest test_payments.py\n",
		DiffText: "diff --git a/mod.py b/mod.py\n+# Fixes #101\n+def discount\n",
		Changed:  []string{"mod.py", "test_mod.py"},
		RunRepro: func(command string) (int, string, error) {
			return 1, "not found", nil
		},
	})
	if ch.Status != "fail" {
		t.Fatalf("status %s: %s", ch.Status, ch.Summary)
	}
}

func TestMatchingClaims(t *testing.T) {
	ch := Evaluate(Input{
		Body:     "Fixes #101\n\n`discount` in `mod.py` now works.\n\nRepro:\npytest test_mod.py::test_discount_basic\n",
		DiffText: "diff --git a/mod.py b/mod.py\n+# Fixes #101\n+def discount(n):\n",
		Commits:  "Fixes #101",
		Changed:  []string{"mod.py", "test_mod.py"},
		RunRepro: func(command string) (int, string, error) {
			if command != "pytest test_mod.py::test_discount_basic" {
				t.Fatalf("command %q", command)
			}
			return 0, "ok", nil
		},
	})
	if ch.Status != "pass" {
		t.Fatalf("status %s: %s", ch.Status, ch.Summary)
	}
}
