# Stage LinMot CT API Documentation

## Overview

The Stage LinMot CT (Command Table) library provides a high-level Go API for controlling LinMot C1250-EC servo drives over EtherCAT. It replaces the traditional LinMot-Talk Command Tables with a more flexible and type-safe Go implementation.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Core Concepts](#core-concepts)
3. [Command Types](#command-types)
4. [Command Table Management](#command-table-management)
5. [Execution Engine](#execution-engine)
6. [Safety System](#safety-system)
7. [Status Monitoring](#status-monitoring)
8. [Unit Conversion](#unit-conversion)
9. [Error Handling](#error-handling)
10. [Best Practices](#best-practices)
11. [Examples](#examples)

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
    // Create a drive controller (replace with actual implementation)
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    // Create command table manager
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, 
        unitConverter, 
        conditionEvaluator, 
        safetyGuard,
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

## Core Concepts

### Command Table

A Command Table is a sequence of commands that define a motion or operation sequence. It's similar to a LinMot-Talk Command Table but implemented in Go with type safety and better error handling.

### Command

A Command represents a single operation to be performed by the drive. Each command has:
- **ID**: Unique identifier within the table
- **Type**: The type of operation (e.g., MoveAbsolute, Wait, SetDigitalOutput)
- **Parameters**: Key-value pairs defining command behavior
- **Conditions**: Optional conditions that must be met before execution

### Execution Engine

The Execution Engine manages the execution of command tables, providing:
- Sequential command execution
- Pause/Resume functionality
- Error handling and recovery
- Status monitoring

## Command Types

### Motion Commands

#### MoveAbsolute
Move to an absolute position.

```go
cmd := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdMoveAbsolute).
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    WithParameter("acceleration", types.NewAccelerationValue(100.0, types.AccelerationUnitCountsS2)).
    WithParameter("jerk", types.NewJerkValue(1000.0, types.JerkUnitCountsS3)).
    Build()
```

**Parameters:**
- `position` (PositionValue): Target position
- `velocity` (VelocityValue): Maximum velocity (optional)
- `acceleration` (AccelerationValue): Maximum acceleration (optional)
- `jerk` (JerkValue): Maximum jerk (optional)

#### MoveRelative
Move by a relative distance.

```go
cmd := types.NewCommandBuilder().
    WithID(2).
    WithType(types.CmdMoveRelative).
    WithParameter("distance", types.NewPositionValue(50.0, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
    Build()
```

#### MoveIncremental
Move by an incremental distance (same as MoveRelative).

#### Jog
Continuous motion at specified velocity.

```go
cmd := types.NewCommandBuilder().
    WithID(3).
    WithType(types.CmdJog).
    WithParameter("velocity", types.NewVelocityValue(10.0, types.VelocityUnitCountsS)).
    Build()
```

#### Stop
Stop all motion.

```go
cmd := types.NewCommandBuilder().
    WithID(4).
    WithType(types.CmdStop).
    Build()
```

### Control Commands

#### Wait
Wait for a specified duration.

```go
cmd := types.NewCommandBuilder().
    WithID(5).
    WithType(types.CmdWait).
    WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
    Build()
```

#### WaitPosition
Wait until position is reached.

```go
cmd := types.NewCommandBuilder().
    WithID(6).
    WithType(types.CmdWaitPosition).
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
    WithParameter("tolerance", 0.1).
    WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitSeconds)).
    Build()
```

#### WaitVelocity
Wait until velocity is reached.

```go
cmd := types.NewCommandBuilder().
    WithID(7).
    WithType(types.CmdWaitVelocity).
    WithParameter("velocity", types.NewVelocityValue(0.0, types.VelocityUnitCountsS)).
    WithParameter("tolerance", 0.1).
    WithParameter("timeout", types.NewTimeValue(2.0, types.TimeUnitSeconds)).
    Build()
```

#### WaitForce
Wait until force is reached.

```go
cmd := types.NewCommandBuilder().
    WithID(8).
    WithType(types.CmdWaitForce).
    WithParameter("force", types.NewForceValue(5.0, types.ForceUnitNewtons)).
    WithParameter("tolerance", 0.1).
    WithParameter("timeout", types.NewTimeValue(3.0, types.TimeUnitSeconds)).
    Build()
```

### I/O Commands

#### SetDigitalOutput
Set a digital output.

```go
cmd := types.NewCommandBuilder().
    WithID(9).
    WithType(types.CmdSetDigitalOutput).
    WithParameter("output", 1).
    WithParameter("value", true).
    Build()
```

#### ClearDigitalOutput
Clear a digital output.

```go
cmd := types.NewCommandBuilder().
    WithID(10).
    WithType(types.CmdClearDigitalOutput).
    WithParameter("output", 1).
    Build()
```

#### SetAnalogOutput
Set an analog output.

```go
cmd := types.NewCommandBuilder().
    WithID(11).
    WithType(types.CmdSetAnalogOutput).
    WithParameter("output", 1).
    WithParameter("value", 3.14).
    Build()
```

#### WaitDigitalInput
Wait for a digital input.

```go
cmd := types.NewCommandBuilder().
    WithID(12).
    WithType(types.CmdWaitDigitalInput).
    WithParameter("input", 1).
    WithParameter("value", true).
    WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitSeconds)).
    Build()
```

#### WaitAnalogInput
Wait for an analog input.

```go
cmd := types.NewCommandBuilder().
    WithID(13).
    WithType(types.CmdWaitAnalogInput).
    WithParameter("input", 1).
    WithParameter("value", 2.5).
    WithParameter("tolerance", 0.1).
    WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitSeconds)).
    Build()
```

### Loop and Jump Commands

#### LoopStart
Start a loop.

```go
cmd := types.NewCommandBuilder().
    WithID(14).
    WithType(types.CmdLoopStart).
    WithParameter("count", 5).
    Build()
```

#### LoopEnd
End a loop.

```go
cmd := types.NewCommandBuilder().
    WithID(15).
    WithType(types.CmdLoopEnd).
    Build()
```

#### LoopBreak
Break out of a loop.

```go
cmd := types.NewCommandBuilder().
    WithID(16).
    WithType(types.CmdLoopBreak).
    Build()
```

#### Jump
Jump to a command.

```go
cmd := types.NewCommandBuilder().
    WithID(17).
    WithType(types.CmdJump).
    WithParameter("target_id", 5).
    Build()
```

#### JumpIfTrue
Jump to a command if condition is true.

```go
cmd := types.NewCommandBuilder().
    WithID(18).
    WithType(types.CmdJumpIfTrue).
    WithParameter("target_id", 10).
    WithParameter("condition", true).
    Build()
```

#### JumpIfFalse
Jump to a command if condition is false.

```go
cmd := types.NewCommandBuilder().
    WithID(19).
    WithType(types.CmdJumpIfFalse).
    WithParameter("target_id", 15).
    WithParameter("condition", false).
    Build()
```

### System Commands

#### Home
Home the drive.

```go
cmd := types.NewCommandBuilder().
    WithID(20).
    WithType(types.CmdHome).
    Build()
```

#### Reset
Reset the drive.

```go
cmd := types.NewCommandBuilder().
    WithID(21).
    WithType(types.CmdReset).
    Build()
```

#### SaveConfiguration
Save drive configuration.

```go
cmd := types.NewCommandBuilder().
    WithID(22).
    WithType(types.CmdSaveConfiguration).
    Build()
```

#### LoadConfiguration
Load drive configuration.

```go
cmd := types.NewCommandBuilder().
    WithID(23).
    WithType(types.CmdLoadConfiguration).
    Build()
```

### Force Control Commands

#### ForceControlOn
Enable force control.

```go
cmd := types.NewCommandBuilder().
    WithID(24).
    WithType(types.CmdForceControlOn).
    Build()
```

#### ForceControlOff
Disable force control.

```go
cmd := types.NewCommandBuilder().
    WithID(25).
    WithType(types.CmdForceControlOff).
    Build()
```

#### SetForce
Set force setpoint.

```go
cmd := types.NewCommandBuilder().
    WithID(26).
    WithType(types.CmdSetForce).
    WithParameter("force", types.NewForceValue(10.0, types.ForceUnitNewtons)).
    Build()
```

### Data Acquisition Commands

#### StartOscilloscope
Start oscilloscope data acquisition.

```go
cmd := types.NewCommandBuilder().
    WithID(27).
    WithType(types.CmdStartOscilloscope).
    Build()
```

#### StopOscilloscope
Stop oscilloscope data acquisition.

```go
cmd := types.NewCommandBuilder().
    WithID(28).
    WithType(types.CmdStopOscilloscope).
    Build()
```

#### SaveData
Save acquired data.

```go
cmd := types.NewCommandBuilder().
    WithID(29).
    WithType(types.CmdSaveData).
    WithParameter("filename", "data.csv").
    Build()
```

## Command Table Management

### Creating Command Tables

```go
// Create a new command table
table := manager.CreateTable("my-table", "My Table", "Description of my table")

// Or use the builder pattern
table := types.NewCommandTableBuilder().
    WithID("my-table").
    WithName("My Table").
    WithDescription("Description of my table").
    WithCommand(command1).
    WithCommand(command2).
    Build()
```

### Adding Commands

```go
// Add a single command
err := manager.AddCommand(table, command)

// Add multiple commands
for _, cmd := range commands {
    err := manager.AddCommand(table, cmd)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Removing Commands

```go
// Remove a command by ID
err := manager.RemoveCommand(table, commandID)
```

### Updating Commands

```go
// Update an existing command
err := manager.UpdateCommand(table, commandID, updatedCommand)
```

### Validating Tables

```go
// Validate a command table
err := manager.ValidateTable(table)
if err != nil {
    log.Fatal("Table validation failed:", err)
}
```

## Execution Engine

### Starting Execution

```go
ctx := context.Background()
err := manager.StartExecution(ctx, table)
if err != nil {
    log.Fatal("Failed to start execution:", err)
}
```

### Controlling Execution

```go
// Pause execution
err := manager.PauseExecution()

// Resume execution
err := manager.ResumeExecution()

// Stop execution
err := manager.StopExecution()
```

### Monitoring Execution

```go
// Get execution status
status := manager.GetExecutionStatus()
fmt.Printf("State: %s\n", status.State)
fmt.Printf("Current Command: %d\n", status.CurrentCommand)
fmt.Printf("Error: %v\n", status.Error)

// Get current command
currentCmd := manager.GetCurrentCommand()
if currentCmd != nil {
    fmt.Printf("Current command type: %s\n", currentCmd.Type.String())
}
```

### Waiting for Completion

```go
// Wait for completion
for {
    status := manager.GetExecutionStatus()
    switch status.State {
    case types.StateCompleted:
        fmt.Println("Execution completed successfully")
        return
    case types.StateError:
        log.Fatal("Execution failed:", status.Error)
    case types.StateStopped:
        fmt.Println("Execution stopped")
        return
    }
    time.Sleep(100 * time.Millisecond)
}
```

## Safety System

### Creating Safety Guards

```go
// Create safety guard with no limits (default)
safetyGuard := stage_linmot_ct.NewSafetyGuard()

// Create safety guard with specific limits
limits := &stage_linmot_ct.SafetyLimits{
    MinPosition: 0.0,
    MaxPosition: 1000.0,
    MaxVelocity: 100.0,
    MinForce:    -50.0,
    MaxForce:    50.0,
}
safetyGuard := stage_linmot_ct.NewSafetyGuardWithLimits(limits)
```

### Setting Safety Limits

```go
// Set position limits
safetyGuard.SetLimits(&stage_linmot_ct.SafetyLimits{
    MinPosition: -500.0,
    MaxPosition: 500.0,
})

// Set velocity limits
safetyGuard.SetLimits(&stage_linmot_ct.SafetyLimits{
    MaxVelocity: 200.0,
})

// Set force limits
safetyGuard.SetLimits(&stage_linmot_ct.SafetyLimits{
    MinForce: -100.0,
    MaxForce: 100.0,
})
```

### Emergency Stop

```go
// Trigger emergency stop
safetyGuard.TriggerEmergencyStop("Emergency stop triggered")

// Check if emergency stop is active
if safetyGuard.IsEmergencyStopActive() {
    fmt.Println("Emergency stop is active")
}
```

## Status Monitoring

### Creating Status Monitor

```go
// Create status monitor
monitor := stage_linmot_ct.NewStatusMonitor(driveController, 100*time.Millisecond)

// Start monitoring
monitor.Start()
defer monitor.Stop()
```

### Getting Status Information

```go
// Get current status
status := monitor.GetStatus()
fmt.Printf("Position: %f\n", status.Position)
fmt.Printf("Velocity: %f\n", status.Velocity)
fmt.Printf("Force: %f\n", status.Force)
fmt.Printf("Drive State: %s\n", status.DriveState)
fmt.Printf("Motion Complete: %t\n", status.MotionComplete)

// Get individual values
position, err := monitor.GetPosition()
velocity, err := monitor.GetVelocity()
force, err := monitor.GetForce()
driveState, err := monitor.GetDriveState()
motionComplete, err := monitor.IsMotionComplete()

// Check if system is healthy
isHealthy := monitor.IsHealthy()
```

### Error Translation

```go
// Create error translator
translator := stage_linmot_ct.NewErrorTranslator()

// Translate an error
err := driveController.MoveAbsolute(ctx, 100.0, 50.0, 100.0, 1000.0)
if err != nil {
    translatedErr := translator.TranslateError(err)
    fmt.Printf("Translated error: %s\n", translatedErr.Error())
}
```

## Unit Conversion

### Creating Unit Converter

```go
// Create unit converter with default scaling factors
converter := types.NewUnitConverter()

// Create unit converter with custom scaling factors
converter := types.NewUnitConverterWithFactors(1000.0, 100.0) // position, force
```

### Converting Values

```go
// Convert position
positionMM := types.NewPositionValue(100.0, types.PositionUnitMillimeters)
positionCounts := converter.ConvertPositionValue(positionMM, types.PositionUnitCounts)

// Convert velocity
velocityMMS := types.NewVelocityValue(50.0, types.VelocityUnitMillimetersPerSecond)
velocityCountsS := converter.ConvertVelocityValue(velocityMMS, types.VelocityUnitCountsS)

// Convert force
forceN := types.NewForceValue(10.0, types.ForceUnitNewtons)
forceCounts := converter.ConvertForceValue(forceN, types.ForceUnitCounts)

// Convert time
timeS := types.NewTimeValue(1.5, types.TimeUnitSeconds)
timeMS := converter.ConvertTimeValue(timeS, types.TimeUnitMilliseconds)
```

### Setting Scaling Factors

```go
// Set position scaling factor (counts per mm)
converter.SetPositionScalingFactor(1000.0)

// Set force scaling factor (counts per Newton)
converter.SetForceScalingFactor(100.0)

// Get scaling factors
posFactor := converter.GetPositionScalingFactor()
forceFactor := converter.GetForceScalingFactor()
```

## Error Handling

### Error Types

The library defines several error types:

```go
// PreconditionError - when preconditions are not met
err := types.NewPreconditionError("motion_state", "drive not ready")

// UnitConversionError - when unit conversion fails
err := types.NewUnitConversionError("invalid unit", "PositionUnit", "Invalid")

// LimitViolation - when safety limits are violated
err := stage_linmot_ct.NewLimitViolation("position", 1500.0, 1000.0, 0.0)
```

### Error Recovery

```go
// Create error recovery system
recovery := stage_linmot_ct.NewErrorRecovery(driveController)

// Set retry parameters
recovery.SetMaxRetries(3)
recovery.SetRetryDelay(1 * time.Second)

// Attempt recovery
err := driveController.MoveAbsolute(ctx, 100.0, 50.0, 100.0, 1000.0)
if err != nil {
    recovered, recoveryErr := recovery.RecoverFromError(command, err)
    if !recovered {
        log.Fatal("Recovery failed:", recoveryErr)
    }
}
```

## Best Practices

### 1. Always Use Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := manager.StartExecution(ctx, table)
```

### 2. Handle Errors Properly

```go
err := manager.AddCommand(table, command)
if err != nil {
    log.Printf("Failed to add command: %v", err)
    return err
}
```

### 3. Use Appropriate Units

```go
// Use millimeters for human-readable positions
position := types.NewPositionValue(100.0, types.PositionUnitMillimeters)

// Use counts for precise control
position := types.NewPositionValue(100000.0, types.PositionUnitCounts)
```

### 4. Implement Safety Limits

```go
// Always set appropriate safety limits
limits := &stage_linmot_ct.SafetyLimits{
    MinPosition: -1000.0,
    MaxPosition: 1000.0,
    MaxVelocity: 100.0,
    MinForce:    -50.0,
    MaxForce:    50.0,
}
safetyGuard := stage_linmot_ct.NewSafetyGuardWithLimits(limits)
```

### 5. Monitor Execution Status

```go
// Always monitor execution status
go func() {
    for {
        status := manager.GetExecutionStatus()
        switch status.State {
        case types.StateCompleted:
            fmt.Println("Execution completed")
            return
        case types.StateError:
            log.Printf("Execution error: %v", status.Error)
            return
        case types.StateStopped:
            fmt.Println("Execution stopped")
            return
        }
        time.Sleep(100 * time.Millisecond)
    }
}()
```

### 6. Use Command Builders

```go
// Use command builders for type safety
command := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdMoveAbsolute).
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    WithComment("Move to position 100").
    Build()
```

### 7. Validate Tables Before Execution

```go
// Always validate tables before execution
err := manager.ValidateTable(table)
if err != nil {
    log.Fatal("Table validation failed:", err)
}

err = manager.StartExecution(ctx, table)
```

## Examples

### Example 1: Simple Motion Sequence

```go
func simpleMotionSequence() {
    // Create components
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Create table
    table := manager.CreateTable("simple", "Simple Motion", "Basic motion sequence")
    
    // Add commands
    commands := []*types.Command{
        types.NewCommandBuilder().
            WithID(1).
            WithType(types.CmdHome).
            Build(),
        types.NewCommandBuilder().
            WithID(2).
            WithType(types.CmdMoveAbsolute).
            WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
            WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
            Build(),
        types.NewCommandBuilder().
            WithID(3).
            WithType(types.CmdWait).
            WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(4).
            WithType(types.CmdMoveAbsolute).
            WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
            WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
            Build(),
    }
    
    for _, cmd := range commands {
        manager.AddCommand(table, cmd)
    }
    
    // Execute
    ctx := context.Background()
    manager.StartExecution(ctx, table)
    
    // Wait for completion
    for {
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        if status.State == types.StateError {
            log.Fatal("Execution failed:", status.Error)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

### Example 2: I/O Control Sequence

```go
func ioControlSequence() {
    // Create components
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Create table
    table := manager.CreateTable("io-control", "IO Control", "I/O control sequence")
    
    // Add commands
    commands := []*types.Command{
        types.NewCommandBuilder().
            WithID(1).
            WithType(types.CmdSetDigitalOutput).
            WithParameter("output", 1).
            WithParameter("value", true).
            Build(),
        types.NewCommandBuilder().
            WithID(2).
            WithType(types.CmdSetAnalogOutput).
            WithParameter("output", 1).
            WithParameter("value", 3.14).
            Build(),
        types.NewCommandBuilder().
            WithID(3).
            WithType(types.CmdWaitDigitalInput).
            WithParameter("input", 1).
            WithParameter("value", true).
            WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(4).
            WithType(types.CmdClearDigitalOutput).
            WithParameter("output", 1).
            Build(),
    }
    
    for _, cmd := range commands {
        manager.AddCommand(table, cmd)
    }
    
    // Execute
    ctx := context.Background()
    manager.StartExecution(ctx, table)
    
    // Wait for completion
    for {
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        if status.State == types.StateError {
            log.Fatal("Execution failed:", status.Error)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

### Example 3: Force Control Sequence

```go
func forceControlSequence() {
    // Create components
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Create table
    table := manager.CreateTable("force-control", "Force Control", "Force control sequence")
    
    // Add commands
    commands := []*types.Command{
        types.NewCommandBuilder().
            WithID(1).
            WithType(types.CmdForceControlOn).
            Build(),
        types.NewCommandBuilder().
            WithID(2).
            WithType(types.CmdSetForce).
            WithParameter("force", types.NewForceValue(5.0, types.ForceUnitNewtons)).
            Build(),
        types.NewCommandBuilder().
            WithID(3).
            WithType(types.CmdWait).
            WithParameter("duration", types.NewTimeValue(2.0, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(4).
            WithType(types.CmdSetForce).
            WithParameter("force", types.NewForceValue(10.0, types.ForceUnitNewtons)).
            Build(),
        types.NewCommandBuilder().
            WithID(5).
            WithType(types.CmdWait).
            WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(6).
            WithType(types.CmdForceControlOff).
            Build(),
    }
    
    for _, cmd := range commands {
        manager.AddCommand(table, cmd)
    }
    
    // Execute
    ctx := context.Background()
    manager.StartExecution(ctx, table)
    
    // Wait for completion
    for {
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        if status.State == types.StateError {
            log.Fatal("Execution failed:", status.Error)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

### Example 4: Loop and Jump Sequence

```go
func loopJumpSequence() {
    // Create components
    driveController := &MyDriveController{}
    unitConverter := types.NewUnitConverter()
    conditionEvaluator := types.NewDefaultConditionEvaluator()
    safetyGuard := stage_linmot_ct.NewSafetyGuard()
    
    manager := stage_linmot_ct.NewCommandTableManager(
        driveController, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Create table
    table := manager.CreateTable("loop-jump", "Loop Jump", "Loop and jump sequence")
    
    // Add commands
    commands := []*types.Command{
        types.NewCommandBuilder().
            WithID(1).
            WithType(types.CmdLoopStart).
            WithParameter("count", 3).
            Build(),
        types.NewCommandBuilder().
            WithID(2).
            WithType(types.CmdMoveAbsolute).
            WithParameter("position", types.NewPositionValue(50.0, types.PositionUnitCounts)).
            Build(),
        types.NewCommandBuilder().
            WithID(3).
            WithType(types.CmdWait).
            WithParameter("duration", types.NewTimeValue(0.5, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(4).
            WithType(types.CmdMoveAbsolute).
            WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
            Build(),
        types.NewCommandBuilder().
            WithID(5).
            WithType(types.CmdWait).
            WithParameter("duration", types.NewTimeValue(0.5, types.TimeUnitSeconds)).
            Build(),
        types.NewCommandBuilder().
            WithID(6).
            WithType(types.CmdLoopEnd).
            Build(),
    }
    
    for _, cmd := range commands {
        manager.AddCommand(table, cmd)
    }
    
    // Execute
    ctx := context.Background()
    manager.StartExecution(ctx, table)
    
    // Wait for completion
    for {
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        if status.State == types.StateError {
            log.Fatal("Execution failed:", status.Error)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

This documentation provides a comprehensive guide to using the Stage LinMot CT library. For more examples and advanced usage, refer to the test files in the repository.