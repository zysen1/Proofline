## Proofline receipt

![Verified](https://img.shields.io/badge/Proofline-Verified-2ea44f)

**Verdict: Verified**

Base `main` (`f883425683`) → head `real-fix` (`6c2db6e03d`).

Test adapter: `pytest`. Wall clock 2833 ms.

| Check | Result | Evidence |
|---|---|---|
| New tests fail before the fix | pass | With only the new tests applied to the base revision, the tests failed on 3/3 runs. |
| New tests pass after the fix | pass | The new tests passed on the head revision (3/3 runs). |
| Tests constrain the changed lines | pass | Killed 9/10 (90%). Operators: compare 1/1, boundary 3/3, delete 3/3, arithmetic 2/3. |
| PR description matches the diff | pass | Files, functions, issues, and repro commands in the description match the change. |
| Scope stays on the claimed fix | pass | Changed lines stay on the tested code. No unmentioned dependency, formatting sweep, or boilerplate dump. |
| Tests are stable across reruns | pass | Fail-before and pass-after each repeated 3 times with a stable result. |

Mutation kill rate: 9/10 (90%).

<details>
<summary>Commands and logs</summary>

#### New tests fail before the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```

#### New tests pass after the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```

#### Tests constrain the changed lines

- compare billing.py:11 `<` → `<=` — killed

```
...F..                                                                   [100%]
=================================== FAILURES ===================================
E   ValueError: price must be >= 0
/tmp/proofline-mut-387065412/billing.py:12: ValueError: price must be >= 0
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_zero_price - ValueError: price must be >= 0
1 failed, 5 passed in 0.01s

```

- boundary billing.py:11 `0` → `1` — killed

```
...F..                                                                   [100%]
=================================== FAILURES ===================================
E   ValueError: price must be >= 0
/tmp/proofline-mut-1859692140/billing.py:12: ValueError: price must be >= 0
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_zero_price - ValueError: price must be >= 0
1 failed, 5 passed in 0.01s

```

- delete billing.py:11 delete statement → pass — killed

```
.....F                                                                   [100%]
=================================== FAILURES ===================================
E   Failed: DID NOT RAISE ValueError
/tmp/proofline-mut-653268291/tests/test_billing.py:36: Failed: DID NOT RAISE ValueError
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_rejects_negative_price - Failed: DID NOT RAISE ValueError
1 failed, 5 passed in 0.01s

```

- arithmetic billing.py:13 `//` → `/` — survived

```
......                                                                   [100%]
6 passed in 0.01s

```

- boundary billing.py:13 `100` → `101` — killed

```
FFF...                                                                   [100%]
=================================== FAILURES ===================================
E   assert 910 == 900
     +  where 910 = apply_discount(1000, 10)
/tmp/proofline-mut-310632321/tests/test_billing.py:15: assert 910 == 900
E   assert 1010 == 1000
     +  where 1010 = apply_discount(1000, 0)
/tmp/proofline-mut-310632321/tests/test_billing.py:19: assert 1010 == 1000
E   assert 2 == 0
     +  where 2 = apply_discount(250, 100)
/tmp/proofline-mut-310632321/tests/test_billing.py:23: assert 2 == 0
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - assert 910 == 900
 +  where 910 = apply_discount(1000, 10)
FAILED tests/test_billing.py::test_apply_discount_zero_percent - assert 1010 == 1000
 +  where 1010 = apply_discount(1000, 0)
FAILED tests/test_billing.py::test_apply_discount_full - assert 2 == 0
 +  where 2 = apply_discount(250, 100)
3 failed, 3 passed in 0.01s

```

- delete billing.py:12 delete statement → pass — killed

```
.....F                                                                   [100%]
=================================== FAILURES ===================================
E   Failed: DID NOT RAISE ValueError
/tmp/proofline-mut-2187686703/tests/test_billing.py:36: Failed: DID NOT RAISE ValueError
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_rejects_negative_price - Failed: DID NOT RAISE ValueError
1 failed, 5 passed in 0.01s

```

- arithmetic billing.py:13 `*` → `/` — killed

```
FFF...                                                                   [100%]
=================================== FAILURES ===================================
E   assert 0.0 == 900
     +  where 0.0 = apply_discount(1000, 10)
/tmp/proofline-mut-1459284025/tests/test_billing.py:15: assert 0.0 == 900
E   assert 0.0 == 1000
     +  where 0.0 = apply_discount(1000, 0)
/tmp/proofline-mut-1459284025/tests/test_billing.py:19: assert 0.0 == 1000
E   ZeroDivisionError: division by zero
/tmp/proofline-mut-1459284025/billing.py:13: ZeroDivisionError: division by zero
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - assert 0.0 == 900
 +  where 0.0 = apply_discount(1000, 10)
FAILED tests/test_billing.py::test_apply_discount_zero_percent - assert 0.0 == 1000
 +  where 0.0 = apply_discount(1000, 0)
FAILED tests/test_billing.py::test_apply_discount_full - ZeroDivisionError: division by zero
3 failed, 3 passed in 0.01s

```

- boundary billing.py:13 `100` → `101` — killed

```
FF....                                                                   [100%]
=================================== FAILURES ===================================
E   assert 891 == 900
     +  where 891 = apply_discount(1000, 10)
/tmp/proofline-mut-2852353256/tests/test_billing.py:15: assert 891 == 900
E   assert 990 == 1000
     +  where 990 = apply_discount(1000, 0)
/tmp/proofline-mut-2852353256/tests/test_billing.py:19: assert 990 == 1000
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - assert 891 == 900
 +  where 891 = apply_discount(1000, 10)
FAILED tests/test_billing.py::test_apply_discount_zero_percent - assert 990 == 1000
 +  where 990 = apply_discount(1000, 0)
2 failed, 4 passed in 0.01s

```

- delete billing.py:13 delete statement → pass — killed

```
FFFF..                                                                   [100%]
=================================== FAILURES ===================================
E   assert None == 900
     +  where None = apply_discount(1000, 10)
/tmp/proofline-mut-2448647554/tests/test_billing.py:15: assert None == 900
E   assert None == 1000
     +  where None = apply_discount(1000, 0)
/tmp/proofline-mut-2448647554/tests/test_billing.py:19: assert None == 1000
E   assert None == 0
     +  where None = apply_discount(250, 100)
/tmp/proofline-mut-2448647554/tests/test_billing.py:23: assert None == 0
E   assert None == 0
     +  where None = apply_discount(0, 10)
/tmp/proofline-mut-2448647554/tests/test_billing.py:27: assert None == 0
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - assert None == 900
 +  where None = apply_discount(1000, 10)
FAILED tests/test_billing.py::test_apply_discount_zero_percent - assert None == 1000
 +  where None = apply_discount(1000, 0)
FAILED tests/test_billing.py::test_apply_discount_full - assert None == 0
 +  where None = apply_discount(250, 100)
FAILED tests/test_billing
…[truncated]
```

- arithmetic billing.py:13 `-` → `+` — killed

```
F.F...                                                                   [100%]
=================================== FAILURES ===================================
E   assert 1100 == 900
     +  where 1100 = apply_discount(1000, 10)
/tmp/proofline-mut-2559857544/tests/test_billing.py:15: assert 1100 == 900
E   assert 500 == 0
     +  where 500 = apply_discount(250, 100)
/tmp/proofline-mut-2559857544/tests/test_billing.py:23: assert 500 == 0
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - assert 1100 == 900
 +  where 1100 = apply_discount(1000, 10)
FAILED tests/test_billing.py::test_apply_discount_full - assert 500 == 0
 +  where 500 = apply_discount(250, 100)
2 failed, 4 passed in 0.01s

```

#### PR description matches the diff


```
python3 -m pytest tests/test_billing.py::test_apply_discount_ten_percent
```

Exit code `0`.

```
============================= test session starts ==============================
platform linux -- Python 3.12.3, pytest-9.1.1, pluggy-1.7.0
rootdir: /tmp/proofline-wt-3422017218
configfile: pytest.ini
collected 1 item

tests/test_billing.py .                                                  [100%]

============================== 1 passed in 0.01s ===============================

```

#### Scope stays on the claimed fix

_No command log._

#### Tests are stable across reruns


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `1`.

```
FFFF.F                                                                   [100%]
=================================== FAILURES ===================================
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-607889404/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount_zero
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `0`.

```
......                                                                   [100%]
6 passed in 0.01s

```

</details>

Proofline judged this from execution evidence. A model summary, when present, does not change the verdict.
