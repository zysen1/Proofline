## Proofline receipt

![Unverified](https://img.shields.io/badge/Proofline-Unverified-d73a49)

**Verdict: Unverified**

Base `main` (`f883425683`) → head `slop-tautology` (`ff1c867e4a`).

Test adapter: `pytest`. Wall clock 2676 ms.

| Check | Result | Evidence |
|---|---|---|
| New tests fail before the fix | pass | With only the new tests applied to the base revision, the tests failed on 3/3 runs. |
| New tests pass after the fix | pass | The new tests passed on the head revision (3/3 runs). |
| Tests constrain the changed lines | fail | Killed 0/10 (0%). Operators: compare 0/1, boundary 0/3, delete 0/3, arithmetic 0/3. |
| PR description matches the diff | pass | Files, functions, issues, and repro commands in the description match the change. |
| Scope stays on the claimed fix | pass | Changed lines stay on the tested code. No unmentioned dependency, formatting sweep, or boilerplate dump. |
| Tests are stable across reruns | pass | Fail-before and pass-after each repeated 3 times with a stable result. |

Mutation kill rate: 0/10 (0%).

<details>
<summary>Commands and logs</summary>

#### New tests fail before the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.01s

```

#### New tests pass after the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```

#### Tests constrain the changed lines

- compare billing.py:11 `<` → `<=` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- boundary billing.py:11 `0` → `1` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- delete billing.py:11 delete statement → pass — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- arithmetic billing.py:13 `//` → `/` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- boundary billing.py:13 `100` → `101` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- delete billing.py:12 delete statement → pass — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- arithmetic billing.py:13 `*` → `/` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- boundary billing.py:13 `100` → `101` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- delete billing.py:13 delete statement → pass — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

- arithmetic billing.py:13 `-` → `+` — survived

```
.                                                                        [100%]
1 passed in 0.00s

```

#### PR description matches the diff


```
python3 -m pytest tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
============================= test session starts ==============================
platform linux -- Python 3.12.3, pytest-9.1.1, pluggy-1.7.0
rootdir: /tmp/proofline-wt-3280403517
configfile: pytest.ini
collected 1 item

tests/test_billing.py .                                                  [100%]

============================== 1 passed in 0.01s ===============================

```

#### Scope stays on the claimed fix

_No command log._

#### Tests are stable across reruns


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `1`.

```
F                                                                        [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2802762386/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_runs - NotImplementedError: apply_discount is not implemented
1 failed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_runs
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```

</details>

Proofline judged this from execution evidence. A model summary, when present, does not change the verdict.
