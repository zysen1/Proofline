Fixes #42

`apply_discount` in `billing.py` now takes a real percentage.

Repro:
python3 -m pytest tests/test_billing.py::test_apply_discount_ten_percent
