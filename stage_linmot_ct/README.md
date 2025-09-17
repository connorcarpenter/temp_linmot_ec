# Stage LinMot CT (Command Table)

A comprehensive Go library that replaces LinMot Command Tables with modern, type-safe Go APIs. This module provides high-level motion control, unit conversion, safety guards, and status monitoring for LinMot C1250-EC servo drives over EtherCAT.

## Overview

The `stage_linmot_ct` module provides a complete replacement for LinMot-Talk Command Tables, offering:

- **Type-Safe Go APIs**: Modern Go interfaces with compile-time type checking
- **Comprehensive Command Support**: All LinMot Command Table commands implemented
- **Unit Conversion**: Automatic conversion between different unit systems
- **Safety Guards**: Built-in safety validation and limit checking
- **Status Monitoring**: Real-time drive status and execution monitoring
- **Error Recovery**: Robust error handling and recovery mechanisms
- **Test-Driven Development**: Comprehensive test coverage and examples

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func main() {
    // Create components
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    // Create command table manager
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Create a command table
    table := manager.CreateTable("my-table", "My Table", "A simple motion sequence")
    
    // Add commands
    moveCmd := types.NewCommandBuilder().
        WithID(1).
        WithType(types.CmdMoveAbsolute).
        WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
        WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
        Build()
    
    manager.AddCommand(table, moveCmd)
    
    // Execute the table
    ctx := context.Background()
    err := manager.StartExecution(ctx, table)
    if err != nil {
        log.Fatal(err)
    }
    
    // Wait for completion
    for {
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        if status.State == types.StateError {
            log.Fatal("Execution failed:", status.Error)
        }
    }
}
```

## Documentation

- **[API Documentation](API_DOCUMENTATION.md)**: Complete API reference with examples
- **[Best Practices](BEST_PRACTICES.md)**: Guidelines for writing robust code
- **[Troubleshooting](TROUBLESHOOTING.md)**: Common issues and solutions
- **[Command Table Reference](LINMOT_COMMAND_TABLE_REFERENCE.md)**: LinMot Command Table documentation
- **[Implementation Plan](IMPLEMENTATION_PLAN.md)**: Development roadmap and phases

## Examples

- **[Simple Motion](examples/simple_motion/main.go)**: Basic motion sequence
- **[I/O Control](examples/io_control/main.go)**: Digital and analog I/O operations
- **[Force Control](examples/force_control/main.go)**: Force control operations
- **[Loop and Jump](examples/loop_jump/main.go)**: Loop and conditional execution

## Architecture

This module implements a comprehensive Command Table system that mirrors the functionality of LinMot-Talk Command Tables while providing modern Go interfaces and type safety.

### Key Components

- **Command Table System**: Complete implementation of LinMot Command Table functionality
- **Execution Engine**: Robust command execution with state management
- **Unit Conversion**: Automatic conversion between different unit systems
- **Safety Guards**: Comprehensive safety checks and validation
- **Condition System**: Flexible condition evaluation for conditional execution
- **Status Monitoring**: Real-time status monitoring and reporting

### Command Types

The module supports all LinMot Command Table command types:

#### Motion Commands
- **MA (Move Absolute)**: Move to absolute position
- **MR (Move Relative)**: Move by relative distance
- **MI (Move Incremental)**: Move by fixed increment
- **JO (Jog)**: Continuous motion
- **ST (Stop)**: Stop motion immediately

#### Control Commands
- **WA (Wait)**: Wait for specified time
- **WP (Wait Position)**: Wait for position condition
- **WV (Wait Velocity)**: Wait for velocity condition
- **WF (Wait Force)**: Wait for force condition

#### I/O Commands
- **DO (Set Digital Output)**: Set digital output state
- **CO (Clear Digital Output)**: Clear digital output state
- **AO (Set Analog Output)**: Set analog output value
- **DI (Wait Digital Input)**: Wait for digital input condition
- **AI (Wait Analog Input)**: Wait for analog input condition

#### Loop Commands
- **LS (Loop Start)**: Start loop
- **LE (Loop End)**: End loop
- **LB (Loop Break)**: Break loop

#### Jump Commands
- **JP (Jump)**: Unconditional jump
- **JT (Jump If True)**: Jump if condition true
- **JF (Jump If False)**: Jump if condition false

#### System Commands
- **HO (Home)**: Home motor
- **RE (Reset)**: Reset drive
- **SC (Save Configuration)**: Save configuration
- **LC (Load Configuration)**: Load configuration

#### Force Control Commands
- **FC (Force Control On)**: Enable force control
- **FO (Force Control Off)**: Disable force control
- **SF (Set Force)**: Set force setpoint

#### Data Acquisition Commands
- **SO (Start Oscilloscope)**: Start data acquisition
- **SP (Stop Oscilloscope)**: Stop data acquisition
- **SD (Save Data)**: Save acquired data

## Public API

### Command Table Management

```go
// Create new command table
func NewCommandTable(id, name string) *CommandTable

// Add command to table
func (ct *CommandTable) AddCommand(cmd *Command) error

// Remove command from table
func (ct *CommandTable) RemoveCommand(id int) error

// Update command in table
func (ct *CommandTable) UpdateCommand(id int, cmd *Command) error

// Validate command table
func (ct *CommandTable) Validate() error
```

### Execution Control

```go
// Start execution
func (ee *ExecutionEngine) Start(ctx context.Context, table *CommandTable) error

// Pause execution
func (ee *ExecutionEngine) Pause() error

// Resume execution
func (ee *ExecutionEngine) Resume() error

// Stop execution
func (ee *ExecutionEngine) Stop() error

// Get execution status
func (ee *ExecutionEngine) GetStatus() ExecutionStatus
```

### Status Monitoring

```go
// Get drive status
func (sm *StatusMonitor) GetStatus(driveID int) (*Status, error)

// Stream drive status
func (sm *StatusMonitor) StreamStatus(ctx context.Context, driveID int) (<-chan *Status, error)

// Get execution status
func (sm *StatusMonitor) GetExecutionStatus() ExecutionStatus
```

## Design Contracts

- **Thread Safety**: All operations are thread-safe
- **Context Support**: All operations support context cancellation
- **Error Handling**: Comprehensive error types and messages
- **Unit Conversion**: Automatic conversion between mm and counts
- **Safety First**: All operations include safety checks
- **Type Safety**: Strong typing for all parameters and values

## Dependencies

- `stage_linmot_drive`: CPython bridge for drive communication
- `context`: Context management
- `time`: Time handling
- `sync`: Synchronization primitives
- `encoding/json`: JSON serialization
- `github.com/pkg/errors`: Enhanced error handling

## Configuration

Configuration is provided through:

- **YAML Files**: `defaults.yaml` for default settings
- **Environment Variables**: Runtime configuration
- **gRPC Calls**: Dynamic configuration updates
- **Command Table Parameters**: Per-command configuration

### Example Configuration

```yaml
command_tables:
  - id: "basic_motion"
    name: "Basic Motion Sequence"
    description: "Simple move and return sequence"
    commands:
      - type: "MA"
        parameters:
          position: 100.0
          unit: "mm"
          velocity: 50.0
          acceleration: 100.0
        next_command: 2
      - type: "WA"
        parameters:
          time: 1000
          unit: "ms"
        next_command: 3
      - type: "MA"
        parameters:
          position: 0.0
          unit: "mm"
          velocity: 50.0
          acceleration: 100.0
        next_command: 0

safety:
  position_limits:
    min: -1000.0
    max: 1000.0
    unit: "mm"
  force_limits:
    max: 1000.0
    unit: "N"
  velocity_limits:
    max: 100.0
    unit: "mm/s"
```

## Concurrency

- **Thread-safe**: All operations are thread-safe with proper mutex protection
- **Context-aware**: Operations respect context cancellation
- **Non-blocking**: Status operations are non-blocking
- **Concurrent execution**: Multiple command tables can be executed concurrently

## Package Layout

```
stage_linmot_ct/
├── go.mod
├── README.md
├── IMPLEMENTATION_PLAN.md
├── LINMOT_COMMAND_TABLE_REFERENCE.md
├── command_table.go          # Command table management
├── commands/                 # Command implementations
│   ├── motion.go            # Motion commands (MA, MR, MI, JO, ST)
│   ├── control.go           # Control commands (WA, WP, WV, WF)
│   ├── io.go                # I/O commands (DO, CO, AO, DI, AI)
│   ├── loop.go              # Loop commands (LS, LE, LB)
│   ├── jump.go              # Jump commands (JP, JT, JF)
│   ├── system.go            # System commands (HO, RE, SC, LC)
│   ├── force.go             # Force control commands (FC, FO, SF)
│   └── data.go              # Data acquisition commands (SO, SP, SD)
├── types/                   # Type definitions
│   ├── commands.go          # Command type definitions
│   ├── parameters.go        # Parameter structures
│   ├── conditions.go        # Condition definitions
│   └── units.go             # Unit conversion utilities
├── execution/               # Execution engine
│   ├── engine.go           # Main execution engine
│   ├── state.go            # Execution state management
│   ├── scheduler.go        # Command scheduling
│   └── validator.go        # Command validation
├── safety/                  # Safety and validation
│   ├── guards.go           # Safety guards
│   ├── limits.go           # Position/force limits
│   └── preconditions.go    # Precondition checking
├── conversion/              # Unit conversion
│   ├── position.go         # Position unit conversion
│   ├── velocity.go         # Velocity unit conversion
│   ├── force.go            # Force unit conversion
│   └── time.go             # Time unit conversion
├── status/                  # Status management
│   ├── monitor.go          # Status monitoring
│   ├── shaping.go          # Status shaping
│   └── errors.go           # Error translation
├── examples/                # Usage examples
│   ├── basic_motion.go     # Basic motion examples
│   ├── complex_sequence.go # Complex sequence examples
│   └── error_handling.go   # Error handling examples
└── tests/                   # Test files
    ├── command_table_test.go
    ├── execution_test.go
    ├── safety_test.go
    └── conversion_test.go
```

## Usage Examples

### Basic Command Table Creation

```go
// Create command table manager
manager := NewCommandTableManager(executionEngine, unitConverter, validator)

// Create a new command table
table := manager.CreateTable("basic_motion", "Basic Motion", "Simple move sequence")

// Add commands
cmd1 := NewCommandBuilder().
    WithID(1).
    WithType(CmdMoveAbsolute).
    WithParameter("position", NewPositionValue(100.0, PositionUnitMM)).
    WithParameter("velocity", NewVelocityValue(50.0, VelocityUnitMMS)).
    WithNextCommand(2).
    Build()

cmd2 := NewCommandBuilder().
    WithID(2).
    WithType(CmdWait).
    WithParameter("time", NewTimeValue(1000.0, TimeUnitMS)).
    WithNextCommand(3).
    Build()

cmd3 := NewCommandBuilder().
    WithID(3).
    WithType(CmdMoveAbsolute).
    WithParameter("position", NewPositionValue(0.0, PositionUnitMM)).
    WithParameter("velocity", NewVelocityValue(50.0, VelocityUnitMMS)).
    WithNextCommand(0).
    Build()

// Add commands to table
manager.AddCommand(table, cmd1)
manager.AddCommand(table, cmd2)
manager.AddCommand(table, cmd3)

// Execute the table
ctx := context.Background()
err := manager.StartExecution(ctx, table)
```

### Conditional Execution

```go
// Create a command with conditions
cmd := NewCommandBuilder().
    WithID(1).
    WithType(CmdMoveAbsolute).
    WithParameter("position", NewPositionValue(100.0, PositionUnitMM)).
    WithCondition(DigitalInputCondition(1, true, nil)).
    WithNextCommand(2).
    Build()
```

### Loop Execution

```go
// Create a loop
loopStart := NewCommandBuilder().
    WithID(1).
    WithType(CmdLoopStart).
    WithParameter("counter_variable", "loop_count").
    WithParameter("max_iterations", 10).
    WithNextCommand(2).
    Build()

loopEnd := NewCommandBuilder().
    WithID(3).
    WithType(CmdLoopEnd).
    WithParameter("counter_variable", "loop_count").
    WithNextCommand(1).
    Build()
```

## Testing

- **Unit Tests**: Comprehensive unit test coverage for all components
- **Integration Tests**: Tests with real hardware via `stage_linmot_drive`
- **Mock Tests**: Tests with mocked drive interface
- **Performance Tests**: Latency and throughput tests
- **Safety Tests**: Validation of safety guards and limits

## Error Handling

The module provides comprehensive error handling with typed errors:

```go
type CommandTableError struct {
    Type        ErrorType
    CommandID   int
    Message     string
    Cause       error
}

type ErrorType int
const (
    ErrInvalidCommand    ErrorType = iota
    ErrInvalidParameter
    ErrTimeout
    ErrSafetyViolation
    ErrDriveFault
    ErrCommunicationError
    ErrExecutionError
)
```

## Performance Considerations

- **Command pre-validation**: Commands are validated before execution
- **Parameter caching**: Frequently used parameters are cached
- **Efficient unit conversion**: Optimized conversion algorithms
- **Minimal memory allocation**: Reduced garbage collection pressure
- **Concurrent execution**: Multiple command tables can run simultaneously

## Future Enhancements

- **Command Table Persistence**: Save/load command tables to/from files
- **Visual Command Table Editor**: GUI for creating and editing command tables
- **Real-time Monitoring Dashboard**: Web-based monitoring interface
- **Advanced Debugging Tools**: Step-by-step execution and variable inspection
- **Performance Analytics**: Detailed performance metrics and optimization suggestions

## Out of Scope

- **Low-level EtherCAT**: Handled by `stage_linmot_drive`
- **gRPC Interface**: Handled by `stage_linmot_app`
- **Configuration Persistence**: Handled by `stage_linmot_app`
- **User Interface**: Handled by external applications