## Proofline receipt

![Unverified](https://img.shields.io/badge/Proofline-Unverified-d73a49)

**Verdict: Unverified**

Base `main` (`f883425683`) → head `slop-broken` (`8dbaa21d70`).

Test adapter: `pytest`. Wall clock 1835 ms.

| Check | Result | Evidence |
|---|---|---|
| New tests fail before the fix | pass | With only the new tests applied to the base revision, the tests failed on 3/3 runs. |
| New tests pass after the fix | fail | The head revision does not build or collect tests. The new tests did not pass on the head revision (0 passing of 3 runs). |
| Tests constrain the changed lines | not_applicable | Mutation was not run because the new tests do not pass on the head revision. |
| PR description matches the diff | warn | Named files, functions, and issues appear in the diff. No repro command was declared. |
| Scope stays on the claimed fix | pass | Changed lines stay on the tested code. No unmentioned dependency, formatting sweep, or boilerplate dump. |
| Tests are stable across reruns | pass | Fail-before and pass-after each repeated 3 times with a stable result. |

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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
…[truncated]
```

#### New tests pass after the fix


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```

#### Tests constrain the changed lines

_No command log._

#### PR description matches the diff

_No command log._

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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
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
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
E   NotImplementedError: apply_discount is not implemented
/tmp/proofline-wt-2316187289/billing.py:11: NotImplementedError: apply_discount is not implemented
=========================== short test summary info ============================
FAILED tests/test_billing.py::test_apply_discount_ten_percent - NotImplementedError: apply_discount is not implemented
FAILED tests/test_billing.py::test_apply_discount
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```


```
python3 -m pytest -q --tb=line -p no:cacheprovider tests/test_billing.py::test_apply_discount_ten_percent tests/test_billing.py::test_apply_discount_zero_percent tests/test_billing.py::test_apply_discount_full tests/test_billing.py::test_apply_discount_zero_price tests/test_billing.py::test_apply_discount_rejects_high_percent tests/test_billing.py::test_apply_discount_rejects_negative_price
```

Exit code `4`.

```
ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_ten_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_full

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_zero_price

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_high_percent

ERROR: found no collectors for /tmp/proofline-wt-355006414/tests/test_billing.py::test_apply_discount_rejects_negative_price


==================================== ERRORS ====================================
____________________ ERROR collecting tests/test_billing.py ____________________
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/python.py:508: in importtestmodule
    mod = import_path(
/home/ubuntu/.local/lib/python3.12/site-packages/_pytest/pathlib.py:596: in import_path
    importlib.import_module(module_name)
/usr/lib/python3.12/importlib/__init__.py:90: in import_module
    return _bootstrap._gcd_
…[truncated]
```

</details>

Proofline judged this from execution evidence. A model summary, when present, does not change the verdict.
