"""Order totals — depends on the generated client (run codegen.py first)."""

from generated_client import net_to_gross


def order_total(net_prices):
    """Sum a list of net prices and return the gross total."""
    return net_to_gross(sum(net_prices))
