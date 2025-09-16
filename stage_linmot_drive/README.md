# stage_linmot_drive

CPython-backed Go module that provides a **1:1 interface to LinMot's vendor Python control library** (as provided by support). Find it in the `python_source` top-level directory in this repo. This is the lowest layer in our stack: it speaks EtherCAT (via the vendor code) and exposes a stable Go API for higher layers.

- **Goal:** expose the vendor script's capabilities with minimal translation, so we can upgrade or diff against the Python source easily.
- **Scope:** lifecycle, motion, I/O, diagnostics, and status calls; **no business logic** and no "command-table equivalents" here.
- **Used by:** `stage_linmot_ct` (our high-level motion layer consumed by the app).

---

## Design

- **Binding style:** CPython embedded in Go using the `go-python/gpython` library. We load the vendor Python module and call its functions directly.
- **API philosophy:** names, parameters, and result shapes mirror the Python script **as-is**; errors are mapped into a small, typed Go error set.
- **EtherCAT reality:** explicitly EtherCAT-based; the module expects an interface name (e.g., `eth0`) or equivalent handle required by the vendor code.
- **Threading model:** Python interpreter calls are serialized through Go's mutex system to keep semantics predictable.

---

## Public surface (categories)

### Core Drive Interface
- `NewDrive(config DriveConfig) (*Drive, error)` - Create a new drive instance
- `Start() error` - Start EtherCAT communication
- `Stop() error` - Stop EtherCAT communication
- `IsActive() bool` - Check if drive is active
- `Close() error` - Clean up resources

### Motor Control
- `SwitchOnMotor(drive interface{}) error` - Switch on motor(s)
- `SwitchOffMotor(drive interface{}) error` - Switch off motor(s)
- `HomeMotor(drive interface{}) error` - Home motor(s)
- `MoveToPosition(drive int, position, maxVelocity, acceleration, deceleration float64, jerk int) (int, error)` - Move to absolute position
- `MoveByOffset(drive int, offset, maxVelocity, acceleration, deceleration float64, jerk int) (int, error)` - Move by relative offset

### Status and Monitoring
- `GetStatus() (map[int]DriveStatus, error)` - Get status of all drives
- `WaitForMotionFinished(drive interface{}, countNibble interface{}, timeout time.Duration) (bool, error)` - Wait for motion completion
- `WaitForTargetPosition(drive interface{}, countNibble interface{}, timeout time.Duration) (bool, error)` - Wait for target position

### Configuration
- `ReadConfig(drive int, header ConfigHeader, upid string) (int, error)` - Read configuration parameter
- `WriteConfig(drive int, header ConfigHeader, upid string, value int) error` - Write configuration parameter

### Data Acquisition
- `SetOscilloscopeRecording(enable bool) error` - Enable/disable oscilloscope recording
- `SaveOscilloscopeData(filename string) error` - Save oscilloscope data to files
- `GetErrorMessages() ([]string, error)` - Get error messages
- `GetInfoMessages() ([]string, error)` - Get info messages

### Drive Information
- `GetDriveInfo() map[int]DriveInfo` - Get information about available drives

---

## Error model

All errors are wrapped in Go's standard error interface with descriptive messages. Python exceptions are converted to Go errors with context about the operation that failed. Each error includes:

- Machine-readable error type
- Descriptive message
- Context about the operation that failed
- Original Python error information (when available)

---

## Concurrency & timing

- All operations are thread-safe through Go's mutex system
- EtherCAT communication runs in a separate Python process
- Go API calls are serialized to prevent race conditions
- Timeout support for all blocking operations
- Context-based cancellation support

---

## Configuration inputs (from the application)

```go
type DriveConfig struct {
    AdapterID           string  `json:"adapter_id"`           // EtherCAT adapter (e.g., "eth0")
    NumDevices          int     `json:"num_devices"`          // Number of connected drives
    CycleTime           float64 `json:"cycle_time"`           // Communication cycle time (seconds)
    NumMonitoring       int     `json:"num_monitoring"`       // Number of monitoring channels
    NumParameter        int     `json:"num_parameter"`        // Number of parameter channels
    ActivateLMDriveData bool    `json:"activate_lm_drive_data"` // Enable advanced data handling
    MpLogging           int     `json:"mp_logging"`           // Logging level (0-50)
}
```

---

## Repository layout

```
stage_linmot_drive/
├── go.mod                    # Go module definition
├── drive.go                  # Main drive interface
├── drive_test.go            # Unit tests
├── python_port/             # CPython bridge implementation
│   ├── ethercat_comm.go     # EtherCAT communication wrapper
│   ├── drive_data.go        # Drive data structures
│   ├── motion_control.go    # Motion control functions
│   ├── housekeeping.go      # Basic motor operations
│   └── configuration.go     # Configuration management
└── example/
    └── main.go              # Example usage
```

---

## Build prerequisites

- Go 1.22 or later
- Python 3.12 or later
- LinMot Python library (from `python_source/` directory)
- Required Python packages:
  - `pysoem` - EtherCAT communication
  - `multiprocessing` - Process management
  - `readerwriterlock` - Thread synchronization

### Installation

1. Install Python dependencies:
```bash
pip install pysoem readerwriterlock
```

2. Copy LinMot Python files to Python path or current directory

3. Build Go module:
```bash
go build ./...
```

---

## Testing approach

- **Unit tests:** Test data structures, utility functions, and API contracts
- **Integration tests:** Test against Python library (requires Python environment)
- **Hardware tests:** Manual testing with actual LinMot hardware

Run tests:
```bash
go test ./...
```

---

## Example usage

```go
package main

import (
    "log"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive"
)

func main() {
    // Configure the drive
    config := stage_linmot_drive.DriveConfig{
        AdapterID:           "eth0",
        NumDevices:          1,
        CycleTime:           0.050,
        NumMonitoring:       4,
        NumParameter:        4,
        ActivateLMDriveData: false,
        MpLogging:           20,
    }

    // Create and start drive
    drive, err := stage_linmot_drive.NewDrive(config)
    if err != nil {
        log.Fatal(err)
    }
    defer drive.Close()

    if err := drive.Start(); err != nil {
        log.Fatal(err)
    }
    defer drive.Stop()

    // Switch on and home motor
    if err := drive.SwitchOnMotor(1); err != nil {
        log.Fatal(err)
    }

    if err := drive.HomeMotor(1); err != nil {
        log.Fatal(err)
    }

    // Move to position
    countNibble, err := drive.MoveToPosition(1, 5.0, 0.01, 0.1, 0.1, 10000)
    if err != nil {
        log.Fatal(err)
    }

    // Wait for completion
    finished, err := drive.WaitForMotionFinished(1, countNibble, 30*time.Second)
    if err != nil {
        log.Fatal(err)
    }

    if finished {
        log.Println("Motion completed successfully")
    }

    // Switch off motor
    if err := drive.SwitchOffMotor(1); err != nil {
        log.Fatal(err)
    }
}
```

---

## Versioning & compatibility

- Compatible with LinMot Python library version 0.82e
- Go module version follows semantic versioning
- Python 3.12+ required
- EtherCAT master support via pysoem

---

## Out of scope

- Command-table equivalents (recipes/verbs) - handled by `stage_linmot_ct`
- Unit conversions (counts↔mm/mm_s) - handled by `stage_linmot_ct`
- gRPC surfaces or authorization - handled by `stage_linmot_app`