package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zysen1/Proofline/internal/dataset"
	"github.com/zysen1/Proofline/internal/demo"
	"github.com/zysen1/Proofline/internal/model"
	"github.com/zysen1/Proofline/internal/rank"
	"github.com/zysen1/Proofline/internal/receipt"
	"github.com/zysen1/Proofline/internal/verify"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var code int
	var err error
	switch os.Args[1] {
	case "verify":
		code, err = cmdVerify(os.Args[2:])
	case "rank":
		code, err = cmdRank(os.Args[2:])
	case "measure":
		code, err = cmdMeasure(os.Args[2:])
	case "demo":
		code, err = cmdDemo(os.Args[2:])
	case "help", "-h", "--help":
		usage()
		return
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "proofline:", err)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Proofline %s — proof receipts for AI-generated changes

  proofline verify --base main --head pr-branch
  proofline rank --base main branch1 branch2 ...
  proofline measure --per-class 20 --out docs/measured-results.json
  proofline demo --toy examples/toy-repo --out docs

verify exits 0 for Verified or Not applicable, 2 for Partially, 1 for Unverified.
`, model.Version)
}

func cmdVerify(args []string) (int, error) {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	repo := fs.String("repo", ".", "git repository")
	base := fs.String("base", "main", "base ref")
	head := fs.String("head", "HEAD", "head ref")
	bodyFile := fs.String("pr-body", "", "pull request description file (default: .proofline/pr.md on the head ref)")
	reruns := fs.Int("reruns", verify.DefaultReruns, "reruns of fail-before and pass-after")
	maxMutants := fs.Int("max-mutants", verify.DefaultMutants, "cap on mutation sites")
	format := fs.String("format", "markdown", "markdown or json")
	output := fs.String("output", "", "write the receipt to this path as well as stdout")
	llm := fs.Bool("llm-summary", false, "add a prose summary; does not affect the verdict")
	timeout := fs.Duration("timeout", 120*time.Second, "per-command timeout")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	body := ""
	if *bodyFile != "" {
		b, err := os.ReadFile(*bodyFile)
		if err != nil {
			return 1, err
		}
		body = string(b)
	}
	rec, err := verify.Verify(context.Background(), verify.Options{
		Repo: *repo, Base: *base, Head: *head, PRBody: body,
		Reruns: *reruns, MaxMutants: *maxMutants, Timeout: *timeout, LLMSummary: *llm,
	})
	if err != nil {
		return 1, err
	}
	text, err := render(*format, rec)
	if err != nil {
		return 1, err
	}
	if *output != "" {
		if err := os.WriteFile(*output, []byte(text), 0o644); err != nil {
			return 1, err
		}
	}
	fmt.Print(text)
	if !endsNL(text) {
		fmt.Println()
	}
	return model.ExitCode(rec.Verdict), nil
}

func cmdRank(args []string) (int, error) {
	fs := flag.NewFlagSet("rank", flag.ContinueOnError)
	repo := fs.String("repo", ".", "git repository")
	base := fs.String("base", "main", "base ref")
	reruns := fs.Int("reruns", verify.DefaultReruns, "reruns of fail-before and pass-after")
	maxMutants := fs.Int("max-mutants", verify.DefaultMutants, "cap on mutation sites")
	format := fs.String("format", "markdown", "markdown or json")
	output := fs.String("output", "", "write the leaderboard to this path")
	timeout := fs.Duration("timeout", 120*time.Second, "per-command timeout")
	llm := fs.Bool("llm-summary", false, "add a prose summary to each receipt; does not affect scores")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	branches := fs.Args()
	if len(branches) < 2 {
		return 2, fmt.Errorf("rank needs at least two branches")
	}
	var recs []model.Receipt
	worst := 0
	for _, br := range branches {
		rec, err := verify.Verify(context.Background(), verify.Options{
			Repo: *repo, Base: *base, Head: br,
			Reruns: *reruns, MaxMutants: *maxMutants, Timeout: *timeout, LLMSummary: *llm,
		})
		if err != nil {
			return 1, fmt.Errorf("%s: %w", br, err)
		}
		recs = append(recs, rec)
		if c := model.ExitCode(rec.Verdict); c > worst {
			worst = c
		}
	}
	rows := rank.Order(recs)
	var text string
	switch *format {
	case "json":
		b, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return 1, err
		}
		text = string(b) + "\n"
	default:
		text = rank.Markdown(rows)
	}
	if *output != "" {
		if err := os.WriteFile(*output, []byte(text), 0o644); err != nil {
			return 1, err
		}
	}
	fmt.Print(text)
	return worst, nil
}

func cmdMeasure(args []string) (int, error) {
	fs := flag.NewFlagSet("measure", flag.ContinueOnError)
	per := fs.Int("per-class", 20, "cases per class, including real fixes")
	reruns := fs.Int("reruns", 2, "reruns inside each case")
	maxMutants := fs.Int("max-mutants", 8, "mutation cap per case")
	out := fs.String("out", "docs/measured-results.json", "report path")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	rep, err := dataset.Measure(context.Background(), *per, *reruns, *maxMutants)
	if err != nil {
		return 1, err
	}
	if err := dataset.Write(*out, rep); err != nil {
		return 1, err
	}
	fmt.Printf("dataset %s  adapter %s  wall %d ms  mean %d ms/case\n", rep.Dataset, rep.Adapter, rep.WallClockMS, rep.MeanCaseMS)
	fmt.Printf("false-positive rate on real_fix (Unverified): %.3f\n", rep.FalsePositive)
	fmt.Printf("oss replay: %s\n", rep.OSSReplay)
	for _, class := range dataset.Classes {
		s := rep.Classes[class]
		if class == dataset.ClassReal {
			fmt.Printf("  %-16s n=%d verified=%d partially=%d unverified=%d na=%d\n", class, s.N, s.Verified, s.Partially, s.Unverified, s.NotApplicable)
			continue
		}
		fmt.Printf("  %-16s detection %.3f (%d/%d) unverified=%d\n", class, s.DetectionRate, s.Detected, s.N, s.Unverified)
	}
	fmt.Printf("wrote %s\n", *out)
	return 0, nil
}

func cmdDemo(args []string) (int, error) {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	toy := fs.String("toy", "examples/toy-repo", "toy fixture directory")
	out := fs.String("out", "docs", "directory for the report, receipts, and PNG")
	reruns := fs.Int("reruns", 3, "reruns")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	return 0, demo.Run(context.Background(), demo.Options{Toy: *toy, Out: *out, Reruns: *reruns})
}

func render(format string, rec model.Receipt) (string, error) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(rec, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b) + "\n", nil
	case "markdown", "":
		return receipt.Markdown(rec), nil
	default:
		return "", fmt.Errorf("unknown format %q", format)
	}
}

func endsNL(s string) bool {
	return s == "" || s[len(s)-1] == '\n'
}
