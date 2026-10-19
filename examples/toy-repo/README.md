# Toy shop fixture

Six preset branches on top of `main`. `apply_discount` raises `NotImplementedError` on `main`. `format_total` already works.

| Branch | What it pretends | What is actually true |
|---|---|---|
| `real-fix` | Fixes #42 | New tests fail before the change, pass after it, and kill mutants |
| `slop-passes-before` | Adds coverage | The new test already passed on `main` |
| `slop-tautology` | A polished discount fix | The new test is `assert True` in disguise |
| `slop-scope` | Fixes the discount | Most of the diff is unrelated files and new dependencies |
| `slop-broken` | Fixes the discount | The function does not parse |
| `slop-false-claims` | Fixes `charge_card` in `payments.py` | The diff is the discount fix, and the repro path does not exist |

`make demo` (from the repository root) turns these snapshots into a git repo and writes the proof receipts.
