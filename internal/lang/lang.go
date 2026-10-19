// Package lang classifies source and test files with tree-sitter.
package lang

import (
	"context"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	js "github.com/smacker/go-tree-sitter/javascript"
	py "github.com/smacker/go-tree-sitter/python"
	tsx "github.com/smacker/go-tree-sitter/typescript/tsx"
	ts "github.com/smacker/go-tree-sitter/typescript/typescript"
)

type ID string

const (
	Python     ID = "python"
	JavaScript ID = "javascript"
	TypeScript ID = "typescript"
	TSX        ID = "tsx"
)

const (
	KindSource   = "source"
	KindTest     = "test"
	KindDoc      = "doc"
	KindConfig   = "config"
	KindManifest = "manifest"
	KindBench    = "bench"
	KindMeta     = "meta"
	KindOther    = "other"
)

func Language(path string) (ID, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".py", ".pyi":
		return Python, true
	case ".js", ".mjs", ".cjs", ".jsx":
		return JavaScript, true
	case ".ts":
		return TypeScript, true
	case ".tsx":
		return TSX, true
	default:
		return "", false
	}
}

func grammar(id ID) *sitter.Language {
	switch id {
	case Python:
		return py.GetLanguage()
	case JavaScript:
		return js.GetLanguage()
	case TypeScript:
		return ts.GetLanguage()
	case TSX:
		return tsx.GetLanguage()
	default:
		return nil
	}
}

// Parse returns the root node. The node keeps the tree alive.
func Parse(id ID, src []byte) (*sitter.Node, error) {
	g := grammar(id)
	if g == nil {
		return nil, nil
	}
	return sitter.ParseCtx(context.Background(), src, g)
}

// Classify labels a changed path. Tree-sitter is used for files whose names
// do not already say "test": a file that contains test functions is not mutated.
func Classify(path string, src []byte) string {
	p := filepath.ToSlash(path)
	base := filepath.Base(p)
	if strings.HasPrefix(p, ".proofline/") || p == ".proofline" {
		return KindMeta
	}
	if isBench(p) {
		return KindBench
	}
	if isDoc(base) {
		return KindDoc
	}
	if isManifest(base) {
		return KindManifest
	}
	if isTestPath(p) {
		return KindTest
	}
	if id, ok := Language(p); ok {
		if len(src) > 0 && hasTests(id, src) {
			return KindTest
		}
		return KindSource
	}
	if isConfig(base) {
		return KindConfig
	}
	return KindOther
}

func isDoc(base string) bool {
	switch strings.ToLower(filepath.Ext(base)) {
	case ".md", ".rst", ".adoc", ".txt", ".markdown":
		return true
	default:
		return false
	}
}

func isConfig(base string) bool {
	switch strings.ToLower(filepath.Ext(base)) {
	case ".yml", ".yaml", ".json", ".toml", ".ini", ".cfg", ".xml":
		return true
	default:
		return false
	}
}

func isManifest(base string) bool {
	switch strings.ToLower(base) {
	case "requirements.txt", "requirements-dev.txt", "package.json", "package-lock.json",
		"pnpm-lock.yaml", "yarn.lock", "go.mod", "go.sum", "pyproject.toml", "pipfile", "pipfile.lock":
		return true
	default:
		return false
	}
}

func isBench(p string) bool {
	base := filepath.Base(p)
	lower := strings.ToLower(base)
	if strings.HasPrefix(lower, "bench_") || strings.HasSuffix(lower, "_bench.py") || strings.HasSuffix(lower, ".bench.js") || strings.HasSuffix(lower, ".bench.ts") {
		return true
	}
	for _, part := range strings.Split(p, "/") {
		switch strings.ToLower(part) {
		case "bench", "benches", "benchmark", "benchmarks", "perf":
			return true
		}
	}
	return false
}

func isTestPath(p string) bool {
	base := filepath.Base(p)
	lower := strings.ToLower(base)
	if strings.HasPrefix(lower, "test_") && strings.HasSuffix(lower, ".py") {
		return true
	}
	if strings.HasSuffix(lower, "_test.py") {
		return true
	}
	if strings.HasSuffix(lower, ".test.js") || strings.HasSuffix(lower, ".test.jsx") ||
		strings.HasSuffix(lower, ".test.ts") || strings.HasSuffix(lower, ".test.tsx") ||
		strings.HasSuffix(lower, ".test.mjs") || strings.HasSuffix(lower, ".test.cjs") ||
		strings.HasSuffix(lower, ".spec.js") || strings.HasSuffix(lower, ".spec.jsx") ||
		strings.HasSuffix(lower, ".spec.ts") || strings.HasSuffix(lower, ".spec.tsx") ||
		strings.HasSuffix(lower, ".spec.mjs") {
		return true
	}
	dir := filepath.ToSlash(filepath.Dir(p))
	if (dir == "tests" || strings.HasPrefix(dir, "tests/") || strings.Contains(dir, "/tests/") ||
		dir == "test" || strings.HasPrefix(dir, "test/") || strings.Contains(dir, "/test/") ||
		dir == "__tests__" || strings.Contains(dir, "/__tests__")) &&
		strings.HasSuffix(lower, ".py") && lower != "__init__.py" {
		return true
	}
	return false
}

func hasTests(id ID, src []byte) bool {
	root, err := Parse(id, src)
	if err != nil || root == nil {
		return false
	}
	found := false
	walk(root, func(n *sitter.Node) {
		if found {
			return
		}
		if isTestNode(id, n, src) {
			found = true
		}
	})
	return found
}

func isTestNode(id ID, n *sitter.Node, src []byte) bool {
	if n == nil {
		return false
	}
	switch id {
	case Python:
		if n.Type() == "function_definition" {
			name := childName(n, src)
			if strings.HasPrefix(name, "test_") {
				return true
			}
		}
		if n.Type() == "class_definition" {
			name := childName(n, src)
			if strings.HasPrefix(name, "Test") {
				return true
			}
		}
	default:
		if n.Type() == "call_expression" {
			fn := callName(n, src)
			switch fn {
			case "test", "it", "describe":
				return true
			}
		}
	}
	return false
}

func childName(n *sitter.Node, src []byte) string {
	name := n.ChildByFieldName("name")
	if name == nil {
		return ""
	}
	return name.Content(src)
}

func callName(n *sitter.Node, src []byte) string {
	fn := n.ChildByFieldName("function")
	if fn == nil && n.ChildCount() > 0 {
		fn = n.Child(0)
	}
	if fn == nil {
		return ""
	}
	if fn.Type() == "identifier" || fn.Type() == "property_identifier" {
		return fn.Content(src)
	}
	// test.each( or describe.skip(
	if fn.Type() == "member_expression" {
		obj := fn.ChildByFieldName("object")
		if obj != nil {
			return obj.Content(src)
		}
	}
	return ""
}

func walk(n *sitter.Node, fn func(*sitter.Node)) {
	if n == nil || n.IsNull() {
		return
	}
	fn(n)
	for i := 0; i < int(n.NamedChildCount()); i++ {
		walk(n.NamedChild(i), fn)
	}
	// Anonymous children carry operators. Callers that need them walk Child().
}

// WalkAll visits every node, including anonymous operators.
func WalkAll(n *sitter.Node, fn func(*sitter.Node)) {
	if n == nil || n.IsNull() {
		return
	}
	fn(n)
	for i := 0; i < int(n.ChildCount()); i++ {
		WalkAll(n.Child(i), fn)
	}
}

// TestSelector is a runner argument for one changed test.
type TestSelector struct {
	Name string // pytest node id or jest/vitest test title
	File string
}

// ChangedTests returns test selectors whose bodies differ from base.
// When parsing fails, the whole file is the selector (empty Name).
func ChangedTests(path string, baseSrc, headSrc []byte) []TestSelector {
	id, ok := Language(path)
	if !ok {
		return []TestSelector{{File: path}}
	}
	head, err := Parse(id, headSrc)
	if err != nil || head == nil {
		return []TestSelector{{File: path}}
	}
	var base *sitter.Node
	if len(baseSrc) > 0 {
		base, _ = Parse(id, baseSrc)
	}
	baseBodies := map[string]string{}
	if base != nil {
		for _, t := range collectTests(id, base, baseSrc, path) {
			baseBodies[t.Name] = t.norm
		}
	}
	var out []TestSelector
	for _, t := range collectTests(id, head, headSrc, path) {
		prev, ok := baseBodies[t.Name]
		if !ok || prev != t.norm {
			out = append(out, TestSelector{Name: t.Name, File: path})
		}
	}
	if len(out) == 0 {
		// The file changed but we could not see a behavioral test diff.
		if isTestPath(path) && string(baseSrc) != string(headSrc) {
			return []TestSelector{{File: path}}
		}
	}
	return out
}

type collected struct {
	Name string
	norm string
}

func collectTests(id ID, root *sitter.Node, src []byte, path string) []collected {
	var out []collected
	if id == Python {
		WalkAll(root, func(n *sitter.Node) {
			if n.Type() != "function_definition" {
				return
			}
			name := childName(n, src)
			if !strings.HasPrefix(name, "test_") {
				return
			}
			nodeID := name
			if p := n.Parent(); p != nil && p.Type() == "decorated_definition" {
				// compare the decorated form so a new decorator counts
			}
			class := enclosingClass(n, src)
			qual := name
			if class != "" {
				qual = class + "::" + name
			}
			bodyNode := n
			if p := n.Parent(); p != nil && p.Type() == "decorated_definition" {
				bodyNode = p
			}
			out = append(out, collected{
				Name: filepath.ToSlash(path) + "::" + qual,
				norm: normSpace(bodyNode.Content(src)),
			})
			_ = nodeID
		})
		return out
	}
	WalkAll(root, func(n *sitter.Node) {
		if n.Type() != "call_expression" {
			return
		}
		fn := callName(n, src)
		if fn != "test" && fn != "it" {
			return
		}
		title := firstStringArg(n, src)
		if title == "" {
			title = fn
		}
		out = append(out, collected{Name: title, norm: normSpace(n.Content(src))})
	})
	return out
}

func enclosingClass(n *sitter.Node, src []byte) string {
	for p := n.Parent(); p != nil && !p.IsNull(); p = p.Parent() {
		if p.Type() == "class_definition" {
			return childName(p, src)
		}
	}
	return ""
}

func firstStringArg(n *sitter.Node, src []byte) string {
	args := n.ChildByFieldName("arguments")
	if args == nil {
		return ""
	}
	var title string
	WalkAll(args, func(c *sitter.Node) {
		if title != "" {
			return
		}
		if c.Type() == "string" || c.Type() == "string_fragment" {
			title = trimQuotes(c.Content(src))
		}
	})
	return title
}

func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'`")
	return s
}

func normSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
