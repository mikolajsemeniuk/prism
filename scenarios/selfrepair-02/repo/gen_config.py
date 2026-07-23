"""Setup step 1/2: writes config.json (the tax rate codegen.py needs).

Stdlib only; no network. Must run before codegen.py.
"""

import json


def main():
    with open("config.json", "w") as f:
        json.dump({"tax_rate": 0.23}, f)
    print("config.json written")


if __name__ == "__main__":
    main()
