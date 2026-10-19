"""Billing helpers for the Proofline toy shop."""


def apply_discount(price_cents: int, percent: int) -> int:
    """Return price_cents after a percentage discount.

    Not implemented on main. 10% off 1000 cents is 900.
    """
    if percent < 0 or percent > 100:
        raise ValueError("percent must be between 0 and 100")
    raise NotImplementedError("apply_discount is not implemented")


def format_total(cents: int) -> str:
    # regression note for the receipt printer
    sign = "-" if cents < 0 else ""
    cents = abs(cents)
    return f"{sign}${cents // 100}.{cents % 100:02d}"
