# stage_linmot_drive

CPython-backed Go module that provides a **1:1 interface to LinMot’s vendor Python control script** (as provided by support). This is the lowest layer in our stack: it speaks EtherCAT (via the vendor code) and exposes a stable Go API for higher layers.

- **Goal:** expose the vendor script’s capabilities with minimal translation, so we can upgrade or diff against the Python source easily.
- **Scope:** lifecycle, motion, I/O, diagnostics, and status calls; **no business logic** and no “command-table equivalents” here.
- **Used by:** `stage_linmot_ct` (our high-level motion layer consumed by the app).

---

## Design

- **Binding style:** CPython embedded in Go (cgo). We load the vendor Python module and call its functions directly.
- **API philosophy:** names, parameters, and result shapes mirror the Python script **as-is**; errors are mapped into a small, typed Go error set.
- **EtherCAT reality:** explicitly EtherCAT-based; the module expects an interface name (e.g., `eth0`) or equivalent handle required by the vendor code.
- **Threading model:** a dedicated OS thread owns the Python interpreter (GIL). Calls are serialized through a bounded work queue to keep semantics predictable.

---

## Public surface (categories)

TBD!

---

## Error model

TBD!

Each error includes a machine-readable code and a short message; original vendor text is preserved for logs.

---

## Concurrency & timing

TBD!

---

## Configuration inputs (from the application)

TBD!

---

## Repository layout

TBD!

---

## Build prerequisites

TBD!

> Exact Python version will match the vendor package we receive; we will pin it and document here.

---

## Testing approach

- **Unit tests:** mock the interpreter boundary; verify argument mapping, error translation, ect.
- **Integration tests:** run against the real vendor Python (no hardware), asserting lifecycles and basic non-motion calls.
- **Hardware smoke (run manually from higher layers):** TBD

---

## Versioning & compatibility

TBD!

---

## Out of scope

- Command-table equivalents (recipes/verbs).
- Unit conversions (counts↔mm/mm_s).
- gRPC surfaces or authorization.

These live in `stage_linmot_ct` and the application module.