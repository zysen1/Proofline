import pytest

from billing import apply_discount, format_total


def test_format_total_whole_dollar():
    assert format_total(500) == "$5.00"


def test_format_total_negative():
    assert format_total(-5) == "-$0.05"


def test_apply_discount_ten_percent():
    assert apply_discount(1000, 10) == 900


def test_apply_discount_zero_percent():
    assert apply_discount(1000, 0) == 1000


def test_apply_discount_full():
    assert apply_discount(250, 100) == 0


def test_apply_discount_zero_price():
    assert apply_discount(0, 10) == 0


def test_apply_discount_rejects_high_percent():
    with pytest.raises(ValueError):
        apply_discount(1000, 101)


def test_apply_discount_rejects_negative_price():
    with pytest.raises(ValueError):
        apply_discount(-5, 10)
