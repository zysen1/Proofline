from billing import format_total


def test_format_total_whole_dollar():
    assert format_total(500) == "$5.00"


def test_format_total_negative():
    assert format_total(-5) == "-$0.05"


def test_format_total_cents():
    assert format_total(1) == "$0.01"
