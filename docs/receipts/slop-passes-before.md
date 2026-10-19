## Proofline receipt

![Unverified](https://img.shields.io/badge/Proofline-Unverified-d73a49)

**Verdict: Unverified**

Base `main` (`f883425683`) → head `slop-passes-before` (`8b4cca76a5`).

Test adapter: `pytest`. Wall clock 1893 ms.

| Check | Result | Evidence |
|---|---|---|
| New tests fail before the fix | fail | The new tests already passed on the base revision (3/3 runs). They do not require this change. |
| New tests pass after the fix | pass | The new tests passed on the head revision (3/3 runs). |
| Tests constrain the changed lines | not_applicable | No compare, boundary, delete, or arithmetic sites were found on the changed source lines. |
| PR description matches the diff | pass | Files, functions, issues, and repro commands in the description match the change. |
| Scope stays on the claimed fix | pass | Changed lines stay on the tested code. No unmentioned dependency, formatting sweep, or boilerplate dump. |
| Tests are stable across reruns | pass | Fail-before and pass-after each repeated 3 times with a stable result. |

<details>
<summary>Commands and logs</summary>

#### New tests fail before the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```

#### New tests pass after the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```

#### Tests constrain the changed lines

_No command log._

#### PR description matches the diff


```
python3 -m pytest tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
============================= test session starts ==============================
platform linux -- Python 3.12.3, pytest-9.1.1, pluggy-1.7.0
rootdir: /tmp/proofline-wt-1305999788
configfile: pytest.ini
collected 1 item

tests/test_billing.py .                                                  [100%]

============================== 1 passed in 0.01s ===============================

```

#### Scope stays on the claimed fix

_No command log._

#### Tests are stable across reruns


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_format_total_cents
```

Exit code `0`.

```
.                                                                        [100%]
1 passed in 0.00s

```

</details>

Proofline judged this from execution evidence. A model summary, when present, does not change the verdict.
