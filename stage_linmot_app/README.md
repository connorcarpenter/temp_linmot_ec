# stage_linmot_app

A Go service that drives a LinMot stage over a dedicated **EtherCAT** control link and exposes a **gRPC** API on the LAN. This repository contains only the application/service layer (gRPC, configuration, telemetry). Motion logic and vendor bindings live in sibling modules.

---

## Why this exists

- **Control link (EtherCAT):** direct cable from the host’s built-in Ethernet port to the LinMot drive (no IP).
- **LAN link (IP):** USB/Ethernet (or another NIC) for gRPC, admin, and updates.
- **Process model:** one process (containerized in Docker/BalenaOS) with a control loop (in lower libs) and a gRPC server (here). The control NIC never carries IP.

---

## Repository layout

    stage_linmot_app/
      cmd/stage_linmot_app/       # process entrypoint (main)
      internal/
        api_grpc/                 # gRPC server + public contract (validation)
        config/                   # load/validate: defaults.yaml
        observe/                  # logs, metrics, audit

**External modules consumed by this app:**
- `stage_linmot_proto/` — versioned `.proto` definitions and generated stubs.
- `stage_linmot_ct/` — high-level motion API that replaces command tables
- `stage_linmot_drive/` — CPython/cgo 1:1 port of the vendor Python (invoked by the command-table layer).

> The app must not import any `internal/...` packages from those modules.

---

## Configuration

Configuration is validated on startup and snapshotted per command.

    config/
      defaults.yaml     # interfaces, control cycle basics

**Conventions**
- Field names include units, e.g., `position_mm`, `velocity_mm_s`, `timeout_ms`.
- Error names are clear nouns: `not_homed`, `fault_active`, `soft_limit_exceeded`, `target_timeout`, `interlock_violation`, `mapping_not_found`.

---

## gRPC API (summary)

TBD!

See `stage_linmot_proto/proto/` for message schemas. All fields use explicit units in their names.

Bind only on the LAN adapter.

---

## Build

    go build ./cmd/stage_linmot_app

Prereqs: Go ≥ 1.22. If you modify protos, regenerate stubs in `stage_linmot_proto` before building.

---

## Run (development)

    ./stage_linmot_app

The server will load and validate configuration, start the gRPC endpoint on the LAN interface, and refuse motion until lower layers report a safe, enabled state.

---

## Telemetry

- **Logs:** structured, TBD
- **Metrics:** control-loop latency/jitter (as surfaced by lower layers), command durations, fault counters, stream health.
- **Health:** liveness/ready, version, and a redacted config snapshot.

---

## Notes

- The application layer is intentionally thin. It delegates motion to the command-table layer, which calls the vendor-compatible drive binding.
