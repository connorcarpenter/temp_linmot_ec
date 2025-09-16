# stage_linmot_ct

High-level motion layer that replaces LinMot **Command Tables** with clear, blocking **verbs**. This module composes vendor-level operations from `stage_linmot_drive` into safe, unit-explicit commands the application can call and reason about.

- **Goal:** provide a stable, human-readable motion API (the “CT-analog” surface).
- **Scope:** commands (verbs), behavioral guards, completion rules, error translation, and status shaping.
- **Not included:** vendor bindings, gRPC, or config file loading.

---

## Why this module exists

Command Tables offer named, parameterized macros inside the drive. When controlling directly, you still need that structure: small, predictable **primitives** and a handful of **macro-level verbs** with strict preconditions and deterministic completion. This module is that layer.

---

## Public surface (package `drive`)

All operations are **blocking** and accept a **deadline** (`context.Context`). Inputs and status fields are **unit-explicit**.

### Commands (verbs)

- TBD!

### Status (indicative fields)

- TBD!

### Error taxonomy

- TBD!

Errors include a machine-readable code and short message; underlying vendor text is preserved in logs upstream.

---

## Design contracts

- **Preconditions first:** verbs refuse motion unless required conditions hold (enabled ∧ !fault ∧ homed ∧ within limits ∧ interlocks satisfied).
- **Deterministic completion:** verbs end by (a) reaching a drive-bit condition, (b) meeting a numeric predicate (`|pos_actual − target| ≤ tolerance_mm`), or (c) timing out.
- **Idempotency:** safe re-entry where sensible (e.g., `Home` when already homed returns quickly).
- **Units at the edge:** callers pass mm / mm_s / ms; internal scaling uses the resolved map from the app.
- **Observability hooks:** each verb accepts a logger/metrics handle from the app; no global logging.

---

## Dependencies

- **Downstream:** `stage_linmot_drive` (vendor CPython 1:1 binding).  
- **Upstream:** the application (e.g., `stage_linmot_app`) for configuration snapshots and telemetry handles.  
- **Wire:** none (no networking here).

This module does **not** import the app; the app constructs and passes a validated configuration snapshot.

---

## Configuration inputs (provided by the app)

> The app owns file I/O and validation; this module consumes an immutable snapshot.

---

## Concurrency & timing guarantees

- One verb is **active at a time** per controller instance.  
- Deadlines are honored; on timeout the verb exits with `target_timeout` and leaves the drive in a consistent, known state (documented per verb).  
- Long-running verbs poll vendor status at a reasonable cadence; jitter tolerance is bounded by configuration.

---

## Package layout

    stage_linmot_ct/
      public/                # exported: public API (verbs, status, errors); no vendor calls
      internal/              # TBD

---

## Testing stance (module level)

- Table-driven unit tests for verbs (success, timeout, fault paths).  
- Behavior tests for transitions and status shaping.  
- Adapter tests proving `port_python` satisfies `port_interface`.

---

## Notes

- Keep the public surface **small and boring**. Prefer composing verbs to adding bespoke one-offs.  
- Only the adapter touches vendor calls; everywhere else speaks in **names** and **units**.
