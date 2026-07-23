# orders-service

A small orders helper library.

## Setup

This project uses a generated API client. Before running anything, generate it:

```
python3 codegen.py
```

This writes `generated_client.py` (it is intentionally not committed). The
application and its tests import from that generated module.

## Tests

```
python3 -m unittest
```
