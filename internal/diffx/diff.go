// Package diffx parses unified git diffs into changed line numbers.
package diffx

import (
	"strconv"
	"strings"
)

// ParseUnified returns 1-indexed added or replaced lines in the new file,
// keyed by the b/ path.
func ParseUnified(diff string) map[string][]int {
	out := map[string][]int{}
	var file string
	newLine := 0
	inHunk := false
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			file = pathFromGitHeader(line)
			inHunk = false
			if file != "" {
				if _, ok := out[file]; !ok {
					out[file] = nil
				}
			}
		case strings.HasPrefix(line, "+++ "):
			p := strings.TrimPrefix(line, "+++ ")
			if p == "/dev/null" {
				file = ""
				continue
			}
			p = strings.TrimPrefix(p, "b/")
			if p != "" {
				file = p
			}
		case strings.HasPrefix(line, "@@"):
			start, ok := newStart(line)
			if !ok {
				inHunk = false
				continue
			}
			newLine = start
			inHunk = true
		case !inHunk || file == "":
			continue
		case strings.HasPrefix(line, "+"):
			if strings.HasPrefix(line, "+++") {
				continue
			}
			out[file] = append(out[file], newLine)
			newLine++
		case strings.HasPrefix(line, "-"):
			// removed line, new-file cursor stays
		case strings.HasPrefix(line, "\\"):
			// "\ No newline at end of file"
		default:
			newLine++
		}
	}
	return out
}

func pathFromGitHeader(line string) string {
	// diff --git a/foo b/foo  or  diff --git "a/foo bar" "b/foo bar"
	rest := strings.TrimPrefix(line, "diff --git ")
	b := strings.LastIndex(rest, " b/")
	if b < 0 {
		return ""
	}
	return rest[b+3:]
}

func newStart(hunk string) (int, bool) {
	// @@ -l,c +l,c @@
	plus := strings.Index(hunk, "+")
	if plus < 0 {
		return 0, false
	}
	rest := hunk[plus+1:]
	end := strings.IndexAny(rest, " ,")
	if end < 0 {
		return 0, false
	}
	n, err := strconv.Atoi(rest[:end])
	if err != nil {
		return 0, false
	}
	return n, true
}

// NumstatPath normalizes a numstat path that may contain a rename.
func NumstatPath(field string) (statusHint, path string) {
	if strings.Contains(field, " => ") {
		// "old => new" or "{old => new}/file"
		if strings.Contains(field, "{") && strings.Contains(field, "}") {
			// keep it simple: take the text after the last " => " and strip braces
			parts := strings.Split(field, " => ")
			right := strings.TrimSuffix(parts[len(parts)-1], "}")
			left := parts[0]
			prefix := ""
			if i := strings.Index(left, "{"); i >= 0 {
				prefix = left[:i]
			}
			return "R", prefix + right
		}
		parts := strings.Split(field, " => ")
		return "R", strings.TrimSpace(parts[len(parts)-1])
	}
	return "", field
}
