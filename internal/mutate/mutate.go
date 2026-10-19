// Package mutate applies four operators to changed source lines:
// comparison replacement, boundary tweak, statement deletion, and arithmetic/logical replacement.
package mutate

import (
	"fmt"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/zysen1/Proofline/internal/lang"
)

const (
	OpCompare    = "compare"
	OpBoundary   = "boundary"
	OpDelete     = "delete"
	OpArithmetic = "arithmetic"
)

// Operators is the fixed set. Order is the round-robin order.
var Operators = []string{OpCompare, OpBoundary, OpDelete, OpArithmetic}

// Mutant is one edited copy of a source file.
type Mutant struct {
	Operator string
	File     string
	Line     int
	Summary  string
	Source   []byte
}

// FileSites describes one changed source file.
type FileSites struct {
	Path         string
	Lang         lang.ID
	Source       []byte
	ChangedLines map[int]bool // 1-indexed
}

// Generate returns up to max mutants, spread across the four operators.
func Generate(files []FileSites, max int) []Mutant {
	if max <= 0 {
		max = 12
	}
	buckets := map[string][]Mutant{}
	for _, f := range files {
		if len(f.Source) == 0 || len(f.ChangedLines) == 0 {
			continue
		}
		root, err := lang.Parse(f.Lang, f.Source)
		if err != nil || root == nil {
			continue
		}
		for _, m := range sites(f, root) {
			buckets[m.Operator] = append(buckets[m.Operator], m)
		}
	}
	var out []Mutant
	for len(out) < max {
		progress := false
		for _, op := range Operators {
			b := buckets[op]
			if len(b) == 0 {
				continue
			}
			out = append(out, b[0])
			buckets[op] = b[1:]
			progress = true
			if len(out) >= max {
				break
			}
		}
		if !progress {
			break
		}
	}
	return out
}

func sites(f FileSites, root *sitter.Node) []Mutant {
	var out []Mutant
	seen := map[string]bool{}
	lang.WalkAll(root, func(n *sitter.Node) {
		if n == nil || n.IsNull() {
			return
		}
		line := int(n.StartPoint().Row) + 1
		if !f.ChangedLines[line] {
			return
		}
		switch n.Type() {
		case "comparison_operator", "boolean_operator", "binary_operator", "binary_expression":
			for i := 0; i < int(n.ChildCount()); i++ {
				c := n.Child(i)
				if c == nil || c.IsNamed() {
					continue
				}
				op := c.Content(f.Source)
				if repl, ok := compareSwap[op]; ok {
					add(&out, seen, f, OpCompare, c, fmt.Sprintf("`%s` → `%s`", op, repl), []byte(repl))
				} else if repl, ok := arithSwap[op]; ok {
					add(&out, seen, f, OpArithmetic, c, fmt.Sprintf("`%s` → `%s`", op, repl), []byte(repl))
				}
			}
		case "integer", "number":
			text := n.Content(f.Source)
			if !digitsOnly(text) {
				return
			}
			repl := bump(text)
			if repl == text {
				return
			}
			add(&out, seen, f, OpBoundary, n, fmt.Sprintf("`%s` → `%s`", text, repl), []byte(repl))
		}
		if isDeletable(n.Type()) && spanLines(n) <= 4 && !insideImport(n) {
			// The node's start is the first token, so the line's existing indent stays put.
			repl := passToken(f.Lang)
			add(&out, seen, f, OpDelete, n, "delete statement → "+strings.TrimSpace(repl), []byte(repl))
		}
	})
	return out
}

var compareSwap = map[string]string{
	"<":   "<=",
	"<=":  "<",
	">":   ">=",
	">=":  ">",
	"==":  "!=",
	"!=":  "==",
	"===": "!==",
	"!==": "===",
}

var arithSwap = map[string]string{
	"+":   "-",
	"-":   "+",
	"*":   "/",
	"/":   "*",
	"//":  "/",
	"%":   "+",
	"and": "or",
	"or":  "and",
	"&&":  "||",
	"||":  "&&",
}

func digitsOnly(s string) bool {
	if s == "" || len(s) > 6 {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func bump(text string) string {
	n, err := strconv.Atoi(text)
	if err != nil {
		return text
	}
	if n == 0 {
		return "1"
	}
	return strconv.Itoa(n + 1)
}

func add(out *[]Mutant, seen map[string]bool, f FileSites, op string, n *sitter.Node, summary string, repl []byte) {
	start := int(n.StartByte())
	end := int(n.EndByte())
	key := fmt.Sprintf("%s:%d:%d:%s", f.Path, start, end, op)
	if seen[key] {
		return
	}
	seen[key] = true
	src := splice(f.Source, start, end, repl)
	if string(src) == string(f.Source) {
		return
	}
	*out = append(*out, Mutant{
		Operator: op,
		File:     f.Path,
		Line:     int(n.StartPoint().Row) + 1,
		Summary:  fmt.Sprintf("%s %s:%d %s", op, f.Path, int(n.StartPoint().Row)+1, summary),
		Source:   src,
	})
}

func splice(src []byte, start, end int, repl []byte) []byte {
	if start < 0 || end > len(src) || start > end {
		return src
	}
	out := make([]byte, 0, len(src)-(end-start)+len(repl))
	out = append(out, src[:start]...)
	out = append(out, repl...)
	out = append(out, src[end:]...)
	return out
}

func isDeletable(t string) bool {
	switch t {
	case "expression_statement", "return_statement", "raise_statement", "assert_statement",
		"throw_statement", "lexical_declaration", "variable_declaration",
		"if_statement", "for_statement", "while_statement", "with_statement":
		return true
	default:
		return false
	}
}

func spanLines(n *sitter.Node) int {
	return int(n.EndPoint().Row) - int(n.StartPoint().Row) + 1
}

func insideImport(n *sitter.Node) bool {
	for p := n; p != nil && !p.IsNull(); p = p.Parent() {
		switch p.Type() {
		case "import_statement", "import_from_statement", "import_declaration":
			return true
		}
	}
	return false
}

func passToken(id lang.ID) string {
	if id == lang.Python {
		return "pass"
	}
	return "void 0;"
}
