# Proofline rank

Branches are ordered by execution evidence, not by how polished the diff looks.

| Rank | Branch | Verdict | Score | Kill rate | Why it is here |
|---|---|---|---:|---:|---|
| 1 | `real-fix` | Verified | 148 | 9/10 (90%) | Fail-before held, tests passed, mutants died, claims and scope matched. |
| 2 | `slop-false-claims` | Unverified | 118 | 9/10 (90%) | PR description matches the diff: Claims not backed by the diff or a passing repro: file payments.py; file tests/test_payments.py; repro failed: python3 -m pytest tests/test_payment… |
| 3 | `slop-scope` | Unverified | 118 | 9/10 (90%) | Scope stays on the claimed fix: unrelated lines 164/206 (80%). |
| 4 | `slop-tautology` | Unverified | 100 | 0/10 (0%) | Tests constrain the changed lines: Killed 0/10 (0%). Operators: compare 0/1, boundary 0/3, delete 0/3, arithmetic 0/3. |

Score: +20 per passing check, +8 per warning, +20×kill rate, +10 when the verdict is Verified. Not-applicable branches score 0. The summary text is not an input.
