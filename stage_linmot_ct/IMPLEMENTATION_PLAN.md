# stage_linmot_ct Implementation Plan

## Overview

The `stage_linmot_ct` module provides a high-level motion control layer that replaces LinMot Command Tables with clear, type-safe Go verbs. This module sits between `stage_linmot_app` (gRPC service) and `stage_linmot_drive` (CPython bridge), providing unit conversion, safety guards, and status shaping.

## Architecture Design

### Module Structure
```
stage_linmot_ct/
├── go.mod
├── README.md
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

## Core Components

### 1. Command Table System

#### Command Table Structure
```go
type CommandTable struct {
    ID          string
    Name        string
    Description string
    Commands    []Command
    Variables   map[string]interface{}
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Command struct {
    ID          int
    Type        CommandType
    Parameters  map[string]interface{}
    Conditions  []Condition
    NextCommand int
    LineNumber  int
    Comment     string
}
```

#### Command Types
```go
type CommandType int

const (
    // Motion Commands
    CmdMoveAbsolute    CommandType = iota // MA
    CmdMoveRelative                       // MR
    CmdMoveIncremental                    // MI
    CmdJog                                // JO
    CmdStop                               // ST
    
    // Control Commands
    CmdWait                               // WA
    CmdWaitPosition                       // WP
    CmdWaitVelocity                       // WV
    CmdWaitForce                          // WF
    
    // I/O Commands
    CmdSetDigitalOutput                   // DO
    CmdClearDigitalOutput                 // CO
    CmdSetAnalogOutput                    // AO
    CmdWaitDigitalInput                   // DI
    CmdWaitAnalogInput                    // AI
    
    // Loop Commands
    CmdLoopStart                          // LS
    CmdLoopEnd                            // LE
    CmdLoopBreak                          // LB
    
    // Jump Commands
    CmdJump                               // JP
    CmdJumpIfTrue                         // JT
    CmdJumpIfFalse                        // JF
    
    // System Commands
    CmdHome                               // HO
    CmdReset                              // RE
    CmdSaveConfiguration                  // SC
    CmdLoadConfiguration                  // LC
    
    // Force Control Commands
    CmdForceControlOn                     // FC
    CmdForceControlOff                    // FO
    CmdSetForce                           // SF
    
    // Data Acquisition Commands
    CmdStartOscilloscope                  // SO
    CmdStopOscilloscope                   // SP
    CmdSaveData                           // SD
)
```

### 2. Parameter Structures

#### Motion Parameters
```go
type MotionParameters struct {
    Position     *PositionValue
    Velocity     *VelocityValue
    Acceleration *AccelerationValue
    Deceleration *AccelerationValue
    Jerk         *JerkValue
    Timeout      *TimeValue
    Tolerance    *PositionValue
}

type PositionValue struct {
    Value float64
    Unit  PositionUnit
}

type VelocityValue struct {
    Value float64
    Unit  VelocityUnit
}

type AccelerationValue struct {
    Value float64
    Unit  AccelerationUnit
}

type JerkValue struct {
    Value float64
    Unit  JerkUnit
}

type TimeValue struct {
    Value float64
    Unit  TimeUnit
}
```

#### I/O Parameters
```go
type IOParameters struct {
    OutputNumber int
    InputNumber  int
    State        bool
    Value        float64
    Tolerance    float64
    Timeout      *TimeValue
}
```

#### Loop Parameters
```go
type LoopParameters struct {
    CounterVariable string
    MaxIterations   int
    Condition       Condition
    NextCommand     int
}
```

### 3. Condition System

#### Condition Types
```go
type ConditionType int

const (
    CondDigitalInput    ConditionType = iota
    CondAnalogInput
    CondPosition
    CondVelocity
    CondForce
    CondTimer
    CondVariable
    CondError
)

type Condition struct {
    Type        ConditionType
    Parameter   string
    Operator    ComparisonOperator
    Value       interface{}
    Timeout     *TimeValue
}

type ComparisonOperator int

const (
    OpEqual              ComparisonOperator = iota
    OpNotEqual
    OpGreaterThan
    OpLessThan
    OpGreaterThanOrEqual
    OpLessThanOrEqual
    OpAnd
    OpOr
    OpNot
)
```

### 4. Execution Engine

#### Execution State
```go
type ExecutionState int

const (
    StateIdle        ExecutionState = iota
    StateRunning
    StatePaused
    StateStopped
    StateError
    StateCompleted
)

type ExecutionContext struct {
    CommandTable    *CommandTable
    CurrentCommand  int
    State           ExecutionState
    Variables       map[string]interface{}
    StartTime       time.Time
    LastUpdateTime  time.Time
    Error           error
    CancelFunc      context.CancelFunc
}
```

#### Execution Engine Interface
```go
type ExecutionEngine interface {
    Start(ctx context.Context, table *CommandTable) error
    Pause() error
    Resume() error
    Stop() error
    GetStatus() ExecutionStatus
    GetCurrentCommand() *Command
    GetVariables() map[string]interface{}
    SetVariable(name string, value interface{}) error
}
```

### 5. Unit Conversion System

#### Unit Types
```go
type PositionUnit int
const (
    PositionUnitMM    PositionUnit = iota
    PositionUnitCounts
)

type VelocityUnit int
const (
    VelocityUnitMMS    VelocityUnit = iota
    VelocityUnitCountsS
)

type AccelerationUnit int
const (
    AccelerationUnitMMS2    AccelerationUnit = iota
    AccelerationUnitCountsS2
)

type ForceUnit int
const (
    ForceUnitN        ForceUnit = iota
    ForceUnitCounts
)

type TimeUnit int
const (
    TimeUnitMS    TimeUnit = iota
    TimeUnitS
)
```

#### Conversion Functions
```go
type UnitConverter interface {
    ConvertPosition(value float64, from, to PositionUnit) float64
    ConvertVelocity(value float64, from, to VelocityUnit) float64
    ConvertAcceleration(value float64, from, to AccelerationUnit) float64
    ConvertForce(value float64, from, to ForceUnit) float64
    ConvertTime(value float64, from, to TimeUnit) float64
}
```

### 6. Safety System

#### Safety Guards
```go
type SafetyGuard interface {
    CheckPreconditions(ctx context.Context, cmd *Command) error
    ValidateParameters(cmd *Command) error
    CheckLimits(cmd *Command) error
    HandleError(err error) error
}

type PositionLimits struct {
    MinPosition float64
    MaxPosition float64
    Unit        PositionUnit
}

type ForceLimits struct {
    MaxForce float64
    Unit     ForceUnit
}

type VelocityLimits struct {
    MaxVelocity float64
    Unit        VelocityUnit
}
```

### 7. Status Management

#### Status Structure
```go
type Status struct {
    DriveID          int
    Position         PositionValue
    Velocity         VelocityValue
    Force            ForceValue
    State            DriveState
    Error            error
    DigitalInputs    []bool
    AnalogInputs     []float64
    DigitalOutputs   []bool
    AnalogOutputs    []float64
    Timestamp        time.Time
}

type DriveState int
const (
    DriveStateDisabled    DriveState = iota
    DriveStateEnabled
    DriveStateHoming
    DriveStateMoving
    DriveStateHolding
    DriveStateError
)
```

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1-2)

#### 1.1 Project Setup
- [ ] Initialize Go module
- [ ] Set up project structure
- [ ] Create basic type definitions
- [ ] Implement unit conversion system

#### 1.2 Command Table Foundation
- [ ] Implement CommandTable structure
- [ ] Create command type definitions
- [ ] Implement parameter structures
- [ ] Add condition system

#### 1.3 Basic Execution Engine
- [ ] Create execution context
- [ ] Implement basic command execution
- [ ] Add state management
- [ ] Implement command validation

### Phase 2: Motion Commands (Week 3-4)

#### 2.1 Motion Command Implementation
- [ ] MoveAbsolute command
- [ ] MoveRelative command
- [ ] MoveIncremental command
- [ ] Jog command
- [ ] Stop command

#### 2.2 Control Commands
- [ ] Wait command
- [ ] WaitPosition command
- [ ] WaitVelocity command
- [ ] WaitForce command

#### 2.3 Integration with stage_linmot_drive
- [ ] Connect to drive interface
- [ ] Implement command translation
- [ ] Add error handling
- [ ] Test basic motion

### Phase 3: Advanced Commands (Week 5-6)

#### 3.1 I/O Commands
- [ ] Digital output commands
- [ ] Analog output commands
- [ ] Digital input commands
- [ ] Analog input commands

#### 3.2 Loop and Jump Commands
- [ ] Loop start/end commands
- [ ] Loop break command
- [ ] Jump commands
- [ ] Conditional execution

#### 3.3 System Commands
- [ ] Home command
- [ ] Reset command
- [ ] Configuration commands

### Phase 4: Force Control (Week 7-8)

#### 4.1 Force Control Commands
- [ ] Force control on/off
- [ ] Set force command
- [ ] Force monitoring

#### 4.2 Data Acquisition
- [ ] Oscilloscope commands
- [ ] Data saving commands
- [ ] Real-time data streaming

### Phase 5: Safety and Validation (Week 9-10)

#### 5.1 Safety System
- [ ] Position limits
- [ ] Force limits
- [ ] Velocity limits
- [ ] Emergency stop

#### 5.2 Precondition Checking
- [ ] Drive state validation
- [ ] Parameter validation
- [ ] Error recovery

#### 5.3 Status Management
- [ ] Status monitoring
- [ ] Error translation
- [ ] Status shaping

### Phase 6: Testing and Documentation (Week 11-12)

#### 6.1 Unit Testing
- [ ] Command table tests
- [ ] Execution engine tests
- [ ] Safety system tests
- [ ] Unit conversion tests

#### 6.2 Integration Testing
- [ ] End-to-end motion tests
- [ ] Error handling tests
- [ ] Performance tests

#### 6.3 Documentation
- [ ] API documentation
- [ ] Usage examples
- [ ] Best practices guide
- [ ] Troubleshooting guide

## API Design

### Public Interface

#### Command Table Management
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

#### Execution Control
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

#### Status Monitoring
```go
// Get drive status
func (sm *StatusMonitor) GetStatus(driveID int) (*Status, error)

// Stream drive status
func (sm *StatusMonitor) StreamStatus(ctx context.Context, driveID int) (<-chan *Status, error)

// Get execution status
func (sm *StatusMonitor) GetExecutionStatus() ExecutionStatus
```

### Configuration

#### Command Table Configuration
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
```

#### Safety Configuration
```yaml
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
  emergency_stop:
    enabled: true
    timeout: 1000
```

## Error Handling

### Error Types
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

### Error Recovery
- Automatic retry for transient errors
- Graceful degradation for non-critical errors
- Emergency stop for safety violations
- Error logging and reporting

## Performance Considerations

### Optimization Strategies
- Command pre-validation
- Parameter caching
- Efficient unit conversion
- Minimal memory allocation
- Concurrent execution where possible

### Resource Management
- Connection pooling
- Memory management
- Garbage collection optimization
- Resource cleanup

## Testing Strategy

### Unit Tests
- Command table validation
- Unit conversion accuracy
- Safety guard functionality
- Error handling

### Integration Tests
- End-to-end motion sequences
- Error recovery scenarios
- Performance under load
- Real hardware testing

### Test Data
- Synthetic command tables
- Edge case parameters
- Error injection scenarios
- Performance benchmarks

## Dependencies

### Internal Dependencies
- `stage_linmot_drive`: CPython bridge for drive communication
- `stage_linmot_proto`: Protocol definitions (future)

### External Dependencies
- `context`: Context management
- `time`: Time handling
- `sync`: Synchronization primitives
- `errors`: Error handling

## Future Enhancements

### Planned Features
- Command table persistence
- Visual command table editor
- Real-time monitoring dashboard
- Advanced debugging tools
- Performance analytics

### Extensibility
- Plugin system for custom commands
- Custom condition types
- User-defined unit systems
- Custom safety guards

## Success Criteria

### Functional Requirements
- [ ] All LinMot Command Table commands implemented
- [ ] Unit conversion system working
- [ ] Safety guards functional
- [ ] Error handling comprehensive
- [ ] Performance meets requirements

### Quality Requirements
- [ ] 90%+ test coverage
- [ ] Documentation complete
- [ ] Code review passed
- [ ] Performance benchmarks met
- [ ] Security audit passed

### Integration Requirements
- [ ] Works with stage_linmot_drive
- [ ] Compatible with stage_linmot_app
- [ ] gRPC interface ready
- [ ] Configuration system integrated
- [ ] Monitoring system connected

This implementation plan provides a comprehensive roadmap for building the `stage_linmot_ct` module with Command Table-like abstractions that will serve as a robust foundation for the LinMot control system.