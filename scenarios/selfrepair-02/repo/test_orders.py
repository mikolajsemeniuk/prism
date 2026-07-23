import unittest

from orders import order_total


class OrderTotalTest(unittest.TestCase):
    def test_single(self):
        self.assertEqual(order_total([100.0]), 123.0)

    def test_many(self):
        self.assertEqual(order_total([10.0, 20.0, 70.0]), 123.0)


if __name__ == "__main__":
    unittest.main()
