# Stage LinMot CT Troubleshooting Guide

## Table of Contents

1. [Common Issues](#common-issues)
2. [Error Messages](#error-messages)
3. [Performance Problems](#performance-problems)
4. [Safety Issues](#safety-issues)
5. [Debugging Techniques](#debugging-techniques)
6. [FAQ](#faq)
7. [Getting Help](#getting-help)

## Common Issues

### 1. Command Table Execution Fails

**Symptoms:**
- `StartExecution` returns an error
- Execution stops immediately
- Status shows `StateError`

**Possible Causes:**
- Invalid command table
- Missing required parameters
- Drive controller not responding
- Safety limits violated

**Solutions:**
```go
// 1. Validate the table before execution
err := manager.ValidateTable(table)
if err != nil {
    log.Printf("Table validation failed: %v", err)
    return err
}

// 2. Check drive controller status
state, err := driveController.GetDriveState()
if err != nil {
    log.Printf("Failed to get drive state: %v", err)
    return err
}
log.Printf("Drive state: %s", state)

// 3. Check safety limits
if safetyGuard.IsEmergencyStopActive() {
    log.Println("Emergency stop is active")
    return fmt.Errorf("emergency stop active")
}

// 4. Start execution with proper error handling
err = manager.StartExecution(ctx, table)
if err != nil {
    log.Printf("Execution failed: %v", err)
    return err
}
```

### 2. Commands Not Executing in Order

**Symptoms:**
- Commands execute out of sequence
- Some commands are skipped
- Loop commands not working properly

**Possible Causes:**
- Invalid command IDs
- Missing loop end commands
- Jump commands with invalid targets

**Solutions:**
```go
// 1. Use sequential command IDs
for i, cmd := range commands {
    cmd.ID = i + 1
}

// 2. Ensure loop commands are properly paired
loopStart := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdLoopStart).
    WithParameter("count", 3).
    Build()

loopEnd := types.NewCommandBuilder().
    WithID(5).
    WithType(types.CmdLoopEnd).
    Build()

// 3. Validate jump targets exist
jumpCmd := types.NewCommandBuilder().
    WithID(3).
    WithType(types.CmdJump).
    WithParameter("target_id", 1). // Make sure command 1 exists
    Build()
```

### 3. Motion Commands Not Working

**Symptoms:**
- Position not changing
- Velocity commands ignored
- Motion stops unexpectedly

**Possible Causes:**
- Drive not in ready state
- Safety limits exceeded
- Invalid position/velocity values
- Drive controller error

**Solutions:**
```go
// 1. Check drive state
state, err := driveController.GetDriveState()
if err != nil {
    return fmt.Errorf("failed to get drive state: %w", err)
}

if state != types.DriveStateReady {
    return fmt.Errorf("drive not ready, current state: %s", state)
}

// 2. Check safety limits
position, err := driveController.GetPosition()
if err != nil {
    return fmt.Errorf("failed to get position: %w", err)
}

if position < safetyLimits.MinPosition || position > safetyLimits.MaxPosition {
    return fmt.Errorf("position %.2f outside safety limits [%.2f, %.2f]", 
        position, safetyLimits.MinPosition, safetyLimits.MaxPosition)
}

// 3. Validate motion parameters
if velocity <= 0 {
    return fmt.Errorf("velocity must be positive, got %.2f", velocity)
}

if acceleration <= 0 {
    return fmt.Errorf("acceleration must be positive, got %.2f", acceleration)
}
```

### 4. I/O Commands Not Working

**Symptoms:**
- Digital outputs not changing
- Analog outputs not setting
- Wait commands timing out

**Possible Causes:**
- Invalid I/O channel numbers
- Drive controller not implementing I/O methods
- Timeout values too short

**Solutions:**
```go
// 1. Check I/O channel numbers
if output < 1 || output > 8 {
    return fmt.Errorf("invalid output channel %d, must be 1-8", output)
}

// 2. Verify drive controller implements I/O methods
if driveController == nil {
    return fmt.Errorf("drive controller is nil")
}

// 3. Use appropriate timeout values
timeout := 5 * time.Second // Adjust based on your system
err := driveController.WaitDigitalInput(ctx, input, value, timeout)
if err != nil {
    log.Printf("WaitDigitalInput failed: %v", err)
}
```

## Error Messages

### 1. "command validation failed: command type cannot be unknown"

**Cause:** Command type is not set or invalid.

**Solution:**
```go
// Make sure to set the command type
command := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdMoveAbsolute). // This is required
    WithParameter("position", position).
    Build()
```

### 2. "missing or invalid parameter: parameter X not found"

**Cause:** Required parameter is missing from command.

**Solution:**
```go
// Add the missing parameter
command := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdMoveAbsolute).
    WithParameter("position", position). // Required parameter
    WithParameter("velocity", velocity). // Optional but recommended
    Build()
```

### 3. "unit conversion failed: invalid unit"

**Cause:** Invalid unit type used in parameter.

**Solution:**
```go
// Use valid unit constants
position := types.NewPositionValue(100.0, types.PositionUnitCounts) // ✅ Correct
position := types.NewPositionValue(100.0, "invalid") // ❌ Wrong
```

### 4. "precondition failed: drive not ready"

**Cause:** Drive is not in ready state for motion.

**Solution:**
```go
// Check drive state before motion
state, err := driveController.GetDriveState()
if err != nil {
    return fmt.Errorf("failed to get drive state: %w", err)
}

if state != types.DriveStateReady {
    // Wait for drive to be ready or handle the state
    log.Printf("Drive not ready, current state: %s", state)
    return fmt.Errorf("drive not ready")
}
```

### 5. "safety limit violated: position X exceeds limit Y"

**Cause:** Motion command would violate safety limits.

**Solution:**
```go
// Set appropriate safety limits
limits := &stage_linmot_ct.SafetyLimits{
    MinPosition: -1000.0,
    MaxPosition: 1000.0,
    MaxVelocity: 100.0,
}
safetyGuard := stage_linmot_ct.NewSafetyGuardWithLimits(limits)

// Or adjust the command parameters
position := types.NewPositionValue(500.0, types.PositionUnitCounts) // Within limits
```

## Performance Problems

### 1. Slow Execution

**Symptoms:**
- Commands take too long to execute
- Overall sequence is slow
- Timeouts occurring

**Solutions:**
```go
// 1. Check command timeouts
command := types.NewCommandBuilder().
    WithType(types.CmdWaitDigitalInput).
    WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitSeconds)). // Reasonable timeout
    Build()

// 2. Optimize polling frequency
ticker := time.NewTicker(100 * time.Millisecond) // Not too frequent
defer ticker.Stop()

// 3. Use appropriate velocity values
velocity := types.NewVelocityValue(50.0, types.VelocityUnitCountsS) // Not too slow
```

### 2. High CPU Usage

**Symptoms:**
- CPU usage is high
- System becomes unresponsive
- Battery drains quickly

**Solutions:**
```go
// 1. Increase polling interval
ticker := time.NewTicker(500 * time.Millisecond) // Less frequent polling
defer ticker.Stop()

// 2. Use context cancellation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// 3. Avoid busy waiting
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-ticker.C:
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            return nil
        }
    }
}
```

### 3. Memory Issues

**Symptoms:**
- Memory usage increases over time
- Out of memory errors
- System becomes slow

**Solutions:**
```go
// 1. Clean up resources
defer func() {
    if monitor != nil {
        monitor.Stop()
    }
}()

// 2. Use appropriate data types
// Use float64 instead of float32 for better precision
position := types.NewPositionValue(100.0, types.PositionUnitCounts)

// 3. Avoid creating too many command tables
// Reuse tables when possible
```

## Safety Issues

### 1. Emergency Stop Not Working

**Symptoms:**
- Emergency stop button not responding
- Motion continues after emergency stop
- Safety system not activated

**Solutions:**
```go
// 1. Check emergency stop status
if safetyGuard.IsEmergencyStopActive() {
    log.Println("Emergency stop is active")
    return fmt.Errorf("emergency stop active")
}

// 2. Implement emergency stop handling
func handleEmergencyStop() {
    safetyGuard.TriggerEmergencyStop("Emergency stop triggered")
    manager.StopExecution()
    log.Println("Emergency stop activated")
}

// 3. Monitor emergency stop in main loop
for {
    if emergencyStopPressed {
        handleEmergencyStop()
        break
    }
    // ... other code
}
```

### 2. Safety Limits Not Enforced

**Symptoms:**
- Motion exceeds safety limits
- No error when limits are violated
- Safety system not working

**Solutions:**
```go
// 1. Set safety limits
limits := &stage_linmot_ct.SafetyLimits{
    MinPosition: -1000.0,
    MaxPosition: 1000.0,
    MaxVelocity: 100.0,
    MinForce:    -50.0,
    MaxForce:    50.0,
}
safetyGuard := stage_linmot_ct.NewSafetyGuardWithLimits(limits)

// 2. Enable safety validation
manager := stage_linmot_ct.NewCommandTableManager(
    driveController, unitConverter, conditionEvaluator, safetyGuard,
)

// 3. Check limits before motion
err := safetyGuard.ValidateMotionCommand(command)
if err != nil {
    return fmt.Errorf("motion command violates safety limits: %w", err)
}
```

## Debugging Techniques

### 1. Enable Debug Logging

```go
import "log"

// Set debug level
log.SetLevel(log.DebugLevel)

// Add debug statements
log.Debugf("Executing command %d: %s", command.ID, command.Type.String())
log.Debugf("Current position: %.2f", position)
log.Debugf("Drive state: %s", state)
```

### 2. Check Execution Status

```go
// Monitor execution status
for {
    status := manager.GetExecutionStatus()
    log.Printf("State: %s, Current Command: %d, Error: %v", 
        status.State, status.CurrentCommand, status.Error)
    
    if status.State == types.StateError {
        log.Printf("Execution error: %v", status.Error)
        break
    }
    
    time.Sleep(100 * time.Millisecond)
}
```

### 3. Validate Command Tables

```go
// Validate table before execution
err := manager.ValidateTable(table)
if err != nil {
    log.Printf("Table validation failed: %v", err)
    return err
}

// Check individual commands
for _, cmd := range table.Commands {
    err := cmd.Validate()
    if err != nil {
        log.Printf("Command %d validation failed: %v", cmd.ID, err)
    }
}
```

### 4. Test Drive Controller

```go
// Test drive controller methods
position, err := driveController.GetPosition()
if err != nil {
    log.Printf("GetPosition failed: %v", err)
}

state, err := driveController.GetDriveState()
if err != nil {
    log.Printf("GetDriveState failed: %v", err)
}

motionComplete, err := driveController.IsMotionComplete()
if err != nil {
    log.Printf("IsMotionComplete failed: %v", err)
}
```

## FAQ

### Q: Why is my command table not executing?

A: Check the following:
1. Validate the table before execution
2. Ensure the drive controller is responding
3. Check that all required parameters are present
4. Verify the drive is in ready state
5. Check for safety limit violations

### Q: How do I handle errors during execution?

A: Use proper error handling:
```go
err := manager.StartExecution(ctx, table)
if err != nil {
    log.Printf("Execution failed: %v", err)
    return err
}

// Monitor execution status
for {
    status := manager.GetExecutionStatus()
    if status.State == types.StateError {
        log.Printf("Execution error: %v", status.Error)
        break
    }
    time.Sleep(100 * time.Millisecond)
}
```

### Q: Why are my I/O commands not working?

A: Check the following:
1. Verify I/O channel numbers are valid (1-8)
2. Ensure drive controller implements I/O methods
3. Check timeout values are appropriate
4. Verify digital/analog values are correct

### Q: How do I implement emergency stop?

A: Use the safety guard:
```go
// Trigger emergency stop
safetyGuard.TriggerEmergencyStop("Emergency stop triggered")

// Check if emergency stop is active
if safetyGuard.IsEmergencyStopActive() {
    log.Println("Emergency stop is active")
    return fmt.Errorf("emergency stop active")
}
```

### Q: Why is my motion command being rejected?

A: Check the following:
1. Drive state is ready
2. Safety limits are not violated
3. Position/velocity values are valid
4. Required parameters are present

### Q: How do I debug unit conversion issues?

A: Use the unit converter directly:
```go
// Test unit conversion
positionMM := types.NewPositionValue(100.0, types.PositionUnitMillimeters)
positionCounts := converter.ConvertPositionValue(positionMM, types.PositionUnitCounts)
log.Printf("100mm = %.2f counts", positionCounts.Value)
```

## Getting Help

### 1. Check Documentation

- Read the API documentation
- Check the best practices guide
- Look at example code

### 2. Run Tests

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestMotionSequence

# Run with verbose output
go test -v ./...
```

### 3. Enable Debug Logging

```go
import "log"

log.SetLevel(log.DebugLevel)
```

### 4. Check Error Messages

- Read error messages carefully
- Check the error context
- Look for specific error types

### 5. Verify Configuration

- Check safety limits
- Verify unit conversion factors
- Ensure drive controller is properly implemented

### 6. Contact Support

If you're still having issues:
1. Collect error messages and logs
2. Provide a minimal reproduction case
3. Include your system configuration
4. Describe the expected vs actual behavior

This troubleshooting guide should help you resolve most common issues with the Stage LinMot CT library.