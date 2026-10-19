`format_total` in `billing.py` already formatted cents. This adds a regression test.

Repro:
python3 -m pytest tests/test_billing.py::test_format_total_cents
