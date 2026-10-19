"""Billing helpers for the Proofline toy shop."""


def apply_discount(price_cents: int, percent: int) -> int
    """Fixes #42. Missing the colon, so this branch does not compile."""
    if percent < 0 or percent > 100:
        raise ValueError("percent must be between 0 and 100")
    return price_cents * (100 - percent) // 100


def format_total(cents: int) -> str:
    sign = "-" if cents < 0 else ""
    cents = abs(cents)
    return f"{sign}${cents // 100}.{cents % 100:02d}"
