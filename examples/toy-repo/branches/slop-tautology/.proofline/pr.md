Fixes #42

`apply_discount` in `billing.py` no longer raises. The write-up is in `README.md`.

Repro:
python3 -m pytest tests/test_billing.py::test_apply_discount_runs
