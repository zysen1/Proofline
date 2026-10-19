package verify

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/zysen1/Proofline/internal/gitx"
	"github.com/zysen1/Proofline/internal/model"
)

func TestJSAdapters(t *testing.T) {
	if _, err := exec.LookPath("npm"); err != nil {
		t.Fatal("npm is required to exercise vitest and jest")
	}
	deps := filepath.Join(t.TempDir(), "deps")
	if err := os.MkdirAll(deps, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("npm", "install", "--silent", "--prefix", deps, "vitest@2.1.9", "jest@29.7.0")
	cmd.Env = append(os.Environ(), "npm_config_fund=false", "npm_config_audit=false")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("npm install: %v\n%s", err, out)
	}
	nodeModules := filepath.Join(deps, "node_modules")

	cases := []struct {
		adapter string
		class   string
		want    model.Verdict
		check   string
	}{
		{"vitest", "real_fix", model.VerdictVerified, ""},
		{"vitest", "tautology", model.VerdictUnverified, "mutation"},
		{"vitest", "false_claims", model.VerdictUnverified, "claims"},
		{"jest", "passes_before", model.VerdictUnverified, "fail_before"},
		{"jest", "broken", model.VerdictUnverified, "pass_after"},
		{"jest", "scope_creep", model.VerdictUnverified, "scope"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.adapter+"_"+tc.class, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), "repo")
			if err := materializeJS(dir, tc.adapter, tc.class); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(nodeModules, filepath.Join(dir, "node_modules")); err != nil {
				t.Fatal(err)
			}
			rec, err := Verify(context.Background(), Options{
				Repo: dir, Base: "main", Head: "pr",
				Reruns: 2, MaxMutants: 8, Timeout: 90 * time.Second,
			})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s %s verdict=%s adapter=%s kill=%.2f", tc.adapter, tc.class, rec.Verdict, rec.Adapter, rec.KillRate)
			for _, c := range rec.Checks {
				t.Logf("  %s %s %s", c.ID, c.Status, c.Summary)
			}
			if rec.Adapter != tc.adapter {
				t.Fatalf("adapter %s, want %s", rec.Adapter, tc.adapter)
			}
			if rec.Verdict != tc.want {
				t.Fatalf("verdict %s, want %s", rec.Verdict, tc.want)
			}
			if tc.check != "" {
				ch, ok := rec.Find(tc.check)
				if !ok || ch.Status != model.StatusFail {
					t.Fatalf("check %s status %s", tc.check, ch.Status)
				}
			}
		})
	}
}

func materializeJS(dir, adapter, class string) error {
	if err := gitx.InitRepo(dir); err != nil {
		return err
	}
	if err := writeJS(dir, adapter, class, false); err != nil {
		return err
	}
	if err := gitx.CommitAll(dir, "base"); err != nil {
		return err
	}
	if err := gitx.CheckoutNew(dir, "pr"); err != nil {
		return err
	}
	if err := gitx.CleanWorktree(dir); err != nil {
		return err
	}
	if err := writeJS(dir, adapter, class, true); err != nil {
		return err
	}
	return gitx.CommitAll(dir, "pr "+class)
}

func writeJS(dir, adapter, class string, head bool) error {
	pkg := map[string]any{
		"name":    "pricing-case",
		"private": true,
		"scripts": map[string]string{},
	}
	if adapter == "vitest" {
		pkg["type"] = "module"
		pkg["scripts"] = map[string]string{"test": "vitest run"}
	} else {
		pkg["scripts"] = map[string]string{"test": "jest --watchAll=false"}
	}
	if head && class == "scope_creep" {
		pkg["dependencies"] = map[string]string{"left-pad": "1.3.0"}
	}
	raw, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), append(raw, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "pricing.js"), []byte(jsSource(adapter, class, head)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "pricing.test.js"), []byte(jsTest(adapter, class, head)), 0o644); err != nil {
		return err
	}
	if head {
		if err := os.MkdirAll(filepath.Join(dir, ".proofline"), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, ".proofline", "pr.md"), []byte(jsBody(adapter, class)), 0o644); err != nil {
			return err
		}
	}
	if head && class == "scope_creep" {
		var noise string
		for i := 0; i < 40; i++ {
			noise += "export function noise" + itoa(i) + "(v) { return v + " + itoa(i) + "; }\n"
		}
		if adapter == "jest" {
			noise = "function noise(v) { return v; }\nmodule.exports = { noise };\n"
			for i := 0; i < 40; i++ {
				noise += "function noise" + itoa(i) + "(v) { return v + " + itoa(i) + "; }\n"
			}
		}
		return os.WriteFile(filepath.Join(dir, "noise.js"), []byte(noise), 0o644)
	}
	return nil
}

func jsSource(adapter, class string, head bool) string {
	fn := "function discount(n) {\n  throw new Error(\"not implemented\");\n}\n"
	good := "function discount(n) {\n  // Fixes #42\n  if (n < 0) {\n    throw new Error(\"n must be >= 0\");\n  }\n  return n * 2 + 1;\n}\n"
	if class == "passes_before" {
		if head {
			good = "function discount(n) {\n  // regression note\n  if (n < 0) {\n    throw new Error(\"n must be >= 0\");\n  }\n  return n * 2 + 1;\n}\n"
		}
		fn = good
		if !head {
			fn = "function discount(n) {\n  if (n < 0) {\n    throw new Error(\"n must be >= 0\");\n  }\n  return n * 2 + 1;\n}\n"
		} else {
			fn = good
		}
	} else if head && class == "broken" {
		fn = "function discount(n) {\n  // Fixes #42\n  return n *\n}\n"
	} else if head {
		fn = good
	}
	if adapter == "vitest" {
		return "export " + fn
	}
	return fn + "module.exports = { discount };\n"
}

func jsTest(adapter, class string, head bool) string {
	if !head {
		if adapter == "vitest" {
			return "import { test, expect } from \"vitest\";\ntest(\"placeholder\", () => { expect(true).toBe(true); });\n"
		}
		return "test(\"placeholder\", () => { expect(true).toBe(true); });\n"
	}
	body := ""
	switch class {
	case "tautology":
		body = "test(\"runs\", () => {\n  const got = discount(4);\n  expect(got).toBe(got);\n  expect(true).toBe(true);\n});\n"
	default:
		body = "test(\"basic\", () => { expect(discount(4)).toBe(9); });\n" +
			"test(\"zero\", () => { expect(discount(0)).toBe(1); });\n" +
			"test(\"neg\", () => { expect(() => discount(-1)).toThrow(); });\n"
	}
	if adapter == "vitest" {
		return "import { expect, test } from \"vitest\";\nimport { discount } from \"./pricing.js\";\n" + body
	}
	return "const { discount } = require(\"./pricing\");\n" + body
}

func jsBody(adapter, class string) string {
	tool := "npx vitest run pricing.test.js"
	if adapter == "jest" {
		tool = "npx jest --watchAll=false pricing.test.js"
	}
	switch class {
	case "false_claims":
		missing := "npx vitest run payments.test.js"
		if adapter == "jest" {
			missing = "npx jest --watchAll=false payments.test.js"
		}
		return "Fixes #9999\n\n`charge` in `payments.js` is fixed.\n\nRepro:\n" + missing + "\n"
	case "broken":
		return "Fixes #42\n\n`discount` in `pricing.js` applies the formula.\n"
	case "passes_before":
		return "`discount` in `pricing.js` already returned the formula.\n\nRepro:\n" + tool + "\n"
	case "tautology":
		return "Fixes #42\n\n`discount` in `pricing.js` no longer throws.\n\nRepro:\nnpx vitest run pricing.test.js -t runs\n"
	default:
		return "Fixes #42\n\n`discount` in `pricing.js` now returns n*2+1.\n\nRepro:\n" + tool + "\n"
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
