# stage\_linmot\_app — System Specification (High‑Level)

## 1) Purpose & scope

* **Goal:** A single Go service (**stage\_linmot\_app**) drives a LinMot drive via a dedicated control cable (**EtherCAT**) and exposes a **gRPC** interface on the LAN for operators/automation.
* **Scope:** Architecture, contracts, configuration model, module/package structure, and dependency rules. No implementation plan here.
* **Hard split:**

  * **Control link** — real‑time fieldbus on the built‑in Ethernet port (**EtherCAT, not IP**).
  * **LAN link** — normal IP network on a USB‑Ethernet adapter (gRPC, admin, updates).

---

## 2) Physical & network topology (reference)

```
Ubuntu dev laptop ─┐
Windows laptop ────┤   (LAN / IP)
Ethernet switch ───┤────────── USB‑Ethernet (control computer) (IP)
                    └────────── Built‑in Ethernet (control cable; EtherCAT) ── LinMot drive ── machine
```

**Notes**

* Built‑in port: direct to the drive; link up, EtherCAT, power‑saving off.
* USB‑Ethernet: IP LAN for gRPC/admin only.

---

## 3) Operating model

* **Single process** on the control computer (Ubuntu/balenaOS), containerized within Docker/BalenaOS.
* **Core loops:**
  * **Control loop** — manages enable/home/motion/state over the control link.
  * **gRPC server** — serves client calls over the LAN.
* **Isolation:** control NIC never carries IP; gRPC stays on the LAN NIC.

---

## 4) gRPC API (public surface)

Proto definitions live in **stage\_linmot\_proto** (versioned; shared with clients).

**Transport**

* Bind on **LAN** adapter only

---

## 5) Configuration model

* Files are validated on load; produce immutable snapshots per command.
* **`defaults.yaml`** — interfaces, control cycle basics.

Conventions:

* Field names are descriptive & unit‑explicit (e.g., `position_mm`, `velocity_mm_s`, `timeout_ms`).
* Error names are clear nouns: `not_homed`, `fault_active`, `soft_limit_exceeded`, `target_timeout`, `interlock_violation`, `mapping_not_found`.

---

## 6) Workspace & modules (four modules, one container)

```
go.work
  use ./stage_linmot_app
      ./stage_linmot_proto
      ./stage_linmot_ct
      ./stage_linmot_drive
```

**Application module (service)**

```
stage_linmot_app/
  cmd/stage_linmot_app/
  internal/
    api_grpc/            # gRPC server + public contract (validation, errors); binds to LAN
    config/              # load/validate defaults.yaml
    observe/             # logs, metrics, audit
```

**Proto module (shared)**

```
stage_linmot_proto/
  proto/                 # gRPC services/messages
  gen/                   # generated code for consumption by app/clients
```

**Command Table module (library)**

```
stage_linmot_ct/
  public/                # EXPORTED: stable API the app uses
  internal/              # TBD
```

**Drive module (explicit 1:1 GoLang/CPython port from vendor's Python scripts)**

```
stage_linmot_drive/
  python_port/           # EXPORTED: CPython/cgo bridge exposing a 1:1 API mapped from LinMot's Python script
                         # Notes: EtherCAT explicit; mirrors vendor naming; no business logic
  internal/              # TBD
```

---

## 7) Dependency direction (downward only; no cycles)

```
stage_linmot_app
  ├─ stage_linmot_proto
  └─ stage_linmot_ct

stage_linmot_ct
  └─ stage_linmot_drive
```

Rules:

* `stage_linmot_app` depends on `stage_linmot_ct` and `stage_linmot_proto`; it does **not** import any `stage_linmot_ct/internal/...`.
* `stage_linmot_ct` reaches the **EtherCAT‑specific** drive via `stage_linmot_drive`.
* gRPC handlers call the **ct** API, not lower layers.

---

## 13) Open questions

* Exact LinMot model/ESI to finalize scaling and the signal map.
* Role boundaries for method‑level authorization.
* Status streaming cadence and retention expected by HMI/clients.
