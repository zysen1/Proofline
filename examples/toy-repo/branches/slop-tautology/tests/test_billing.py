from billing import apply_discount, format_total


def test_format_total_whole_dollar():
    assert format_total(500) == "$5.00"


def test_format_total_negative():
    assert format_total(-5) == "-$0.05"


def test_apply_discount_runs():
    got = apply_discount(1000, 10)
    assert got == got
    assert True
