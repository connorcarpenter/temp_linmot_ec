# Stage LinMot CT Best Practices Guide

## Table of Contents

1. [General Principles](#general-principles)
2. [Command Table Design](#command-table-design)
3. [Error Handling](#error-handling)
4. [Performance Optimization](#performance-optimization)
5. [Safety Considerations](#safety-considerations)
6. [Testing Strategies](#testing-strategies)
7. [Code Organization](#code-organization)
8. [Documentation](#documentation)
9. [Common Pitfalls](#common-pitfalls)
10. [Troubleshooting](#troubleshooting)

## General Principles

### 1. Always Use Context

**✅ Good:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

err := manager.StartExecution(ctx, table)
```

**❌ Bad:**
```go
err := manager.StartExecution(context.Background(), table)
```

### 2. Handle Errors Properly

**✅ Good:**
```go
err := manager.AddCommand(table, command)
if err != nil {
    log.Printf("Failed to add command: %v", err)
    return fmt.Errorf("command addition failed: %w", err)
}
```

**❌ Bad:**
```go
manager.AddCommand(table, command) // Ignoring error
```

### 3. Use Appropriate Units

**✅ Good:**
```go
// Use millimeters for human-readable positions
position := types.NewPositionValue(100.0, types.PositionUnitMillimeters)

// Use counts for precise control
position := types.NewPositionValue(100000.0, types.PositionUnitCounts)
```

**❌ Bad:**
```go
// Mixing units without conversion
position := types.NewPositionValue(100.0, types.PositionUnitMillimeters)
velocity := types.NewVelocityValue(50.0, types.VelocityUnitCountsS) // Inconsistent units
```

## Command Table Design

### 1. Use Descriptive Command IDs

**✅ Good:**
```go
// Use sequential IDs with meaningful comments
types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdHome).
    WithComment("Home the drive to establish reference").
    Build(),

types.NewCommandBuilder().
    WithID(2).
    WithType(types.CmdMoveAbsolute).
    WithComment("Move to pick position").
    Build(),
```

**❌ Bad:**
```go
// Random IDs without comments
types.NewCommandBuilder().
    WithID(42).
    WithType(types.CmdMoveAbsolute).
    Build(),
```

### 2. Group Related Commands

**✅ Good:**
```go
// Group I/O operations together
commands := []*types.Command{
    // Setup phase
    types.NewCommandBuilder().WithID(1).WithType(types.CmdHome).Build(),
    types.NewCommandBuilder().WithID(2).WithType(types.CmdSetDigitalOutput).Build(),
    
    // Motion phase
    types.NewCommandBuilder().WithID(3).WithType(types.CmdMoveAbsolute).Build(),
    types.NewCommandBuilder().WithID(4).WithType(types.CmdWait).Build(),
    
    // Cleanup phase
    types.NewCommandBuilder().WithID(5).WithType(types.CmdClearDigitalOutput).Build(),
    types.NewCommandBuilder().WithID(6).WithType(types.CmdMoveAbsolute).Build(),
}
```

### 3. Use Command Builders

**✅ Good:**
```go
command := types.NewCommandBuilder().
    WithID(1).
    WithType(types.CmdMoveAbsolute).
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    WithComment("Move to position 100").
    Build()
```

**❌ Bad:**
```go
// Manual construction is error-prone
command := &types.Command{
    ID: 1,
    Type: types.CmdMoveAbsolute,
    Parameters: map[string]interface{}{
        "position": types.NewPositionValue(100.0, types.PositionUnitCounts),
        "velocity": types.NewVelocityValue(50.0, types.VelocityUnitCountsS),
    },
    Comment: "Move to position 100",
}
```

### 4. Validate Tables Before Execution

**✅ Good:**
```go
// Always validate before execution
err := manager.ValidateTable(table)
if err != nil {
    return fmt.Errorf("table validation failed: %w", err)
}

err = manager.StartExecution(ctx, table)
```

**❌ Bad:**
```go
// Skipping validation
err := manager.StartExecution(ctx, table)
```

## Error Handling

### 1. Use Wrapped Errors

**✅ Good:**
```go
err := manager.AddCommand(table, command)
if err != nil {
    return fmt.Errorf("failed to add command %d: %w", command.ID, err)
}
```

**❌ Bad:**
```go
err := manager.AddCommand(table, command)
if err != nil {
    return err // Loses context
}
```

### 2. Handle Specific Error Types

**✅ Good:**
```go
err := manager.StartExecution(ctx, table)
if err != nil {
    var preconditionErr *types.PreconditionError
    if errors.As(err, &preconditionErr) {
        log.Printf("Precondition failed: %s", preconditionErr.Message)
        // Handle precondition error
    } else {
        log.Printf("Execution failed: %v", err)
    }
}
```

### 3. Implement Retry Logic

**✅ Good:**
```go
const maxRetries = 3
const retryDelay = 1 * time.Second

for i := 0; i < maxRetries; i++ {
    err := manager.StartExecution(ctx, table)
    if err == nil {
        break
    }
    
    if i < maxRetries-1 {
        log.Printf("Attempt %d failed, retrying in %v: %v", i+1, retryDelay, err)
        time.Sleep(retryDelay)
    } else {
        return fmt.Errorf("execution failed after %d attempts: %w", maxRetries, err)
    }
}
```

## Performance Optimization

### 1. Use Appropriate Timeouts

**✅ Good:**
```go
// Set reasonable timeouts based on operation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// For quick operations
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 2. Monitor Execution Status Efficiently

**✅ Good:**
```go
// Use appropriate polling interval
ticker := time.NewTicker(100 * time.Millisecond)
defer ticker.Stop()

for {
    select {
    case <-ticker.C:
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            return nil
        }
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**❌ Bad:**
```go
// Too frequent polling
for {
    status := manager.GetExecutionStatus()
    if status.State == types.StateCompleted {
        return nil
    }
    // No delay - wastes CPU
}
```

### 3. Use Unit Conversion Efficiently

**✅ Good:**
```go
// Convert units once and reuse
positionCounts := converter.ConvertPositionValue(
    types.NewPositionValue(100.0, types.PositionUnitMillimeters),
    types.PositionUnitCounts,
)

// Use the converted value multiple times
command1 := types.NewCommandBuilder().
    WithParameter("position", positionCounts).
    Build()

command2 := types.NewCommandBuilder().
    WithParameter("position", positionCounts).
    Build()
```

## Safety Considerations

### 1. Always Set Safety Limits

**✅ Good:**
```go
limits := &stage_linmot_ct.SafetyLimits{
    MinPosition: -1000.0,
    MaxPosition: 1000.0,
    MaxVelocity: 100.0,
    MinForce:    -50.0,
    MaxForce:    50.0,
}
safetyGuard := stage_linmot_ct.NewSafetyGuardWithLimits(limits)
```

**❌ Bad:**
```go
// No safety limits
safetyGuard := stage_linmot_ct.NewSafetyGuard()
```

### 2. Implement Emergency Stop

**✅ Good:**
```go
// Handle emergency stop
if emergencyStop {
    safetyGuard.TriggerEmergencyStop("Emergency stop triggered")
    manager.StopExecution()
    return fmt.Errorf("emergency stop activated")
}
```

### 3. Validate Motion Commands

**✅ Good:**
```go
// The safety guard automatically validates motion commands
command := types.NewCommandBuilder().
    WithType(types.CmdMoveAbsolute).
    WithParameter("position", types.NewPositionValue(1500.0, types.PositionUnitCounts)). // Exceeds limit
    Build()

// This will be rejected by the safety guard
err := manager.AddCommand(table, command)
```

### 4. Use Precondition Checking

**✅ Good:**
```go
// Check drive state before motion
if driveState != types.DriveStateReady {
    return fmt.Errorf("drive not ready for motion, current state: %s", driveState)
}
```

## Testing Strategies

### 1. Use Mock Controllers

**✅ Good:**
```go
func TestMotionSequence(t *testing.T) {
    mockDrive := &MockDriveController{
        position: 0,
        velocity: 0,
    }
    
    manager := stage_linmot_ct.NewCommandTableManager(
        mockDrive, unitConverter, conditionEvaluator, safetyGuard,
    )
    
    // Test your sequence
}
```

### 2. Test Error Conditions

**✅ Good:**
```go
func TestErrorHandling(t *testing.T) {
    mockDrive := &MockDriveController{
        shouldError: true,
    }
    
    err := manager.StartExecution(ctx, table)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expected error message")
}
```

### 3. Test Edge Cases

**✅ Good:**
```go
func TestEmptyTable(t *testing.T) {
    table := manager.CreateTable("empty", "Empty", "Empty table")
    
    err := manager.StartExecution(ctx, table)
    assert.NoError(t, err)
    
    status := manager.GetExecutionStatus()
    assert.Equal(t, types.StateCompleted, status.State)
}
```

## Code Organization

### 1. Separate Concerns

**✅ Good:**
```go
// Separate file for command table creation
func CreatePickAndPlaceTable() *types.CommandTable {
    // Table creation logic
}

// Separate file for execution logic
func ExecuteTable(ctx context.Context, table *types.CommandTable) error {
    // Execution logic
}
```

### 2. Use Constants for Magic Numbers

**✅ Good:**
```go
const (
    PickPosition    = 100.0
    PlacePosition   = 200.0
    DefaultVelocity = 50.0
    WaitTime        = 1.0
)

command := types.NewCommandBuilder().
    WithParameter("position", types.NewPositionValue(PickPosition, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(DefaultVelocity, types.VelocityUnitCountsS)).
    Build()
```

**❌ Bad:**
```go
command := types.NewCommandBuilder().
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    Build()
```

### 3. Use Configuration Files

**✅ Good:**
```go
type Config struct {
    Positions struct {
        Home    float64 `yaml:"home"`
        Pick    float64 `yaml:"pick"`
        Place   float64 `yaml:"place"`
    } `yaml:"positions"`
    Velocities struct {
        Default float64 `yaml:"default"`
        Fast    float64 `yaml:"fast"`
        Slow    float64 `yaml:"slow"`
    } `yaml:"velocities"`
}

func LoadConfig(filename string) (*Config, error) {
    // Load configuration from file
}
```

## Documentation

### 1. Document Command Tables

**✅ Good:**
```go
// CreatePickAndPlaceTable creates a command table for pick and place operations.
// The table performs the following sequence:
// 1. Home the drive
// 2. Move to pick position
// 3. Wait for part detection
// 4. Move to place position
// 5. Wait for placement confirmation
// 6. Return to home position
func CreatePickAndPlaceTable() *types.CommandTable {
    // Implementation
}
```

### 2. Document Parameters

**✅ Good:**
```go
// CreateMotionCommand creates a motion command with the specified parameters.
// Parameters:
//   - position: Target position in counts
//   - velocity: Maximum velocity in counts/second
//   - acceleration: Maximum acceleration in counts/second²
//   - jerk: Maximum jerk in counts/second³
func CreateMotionCommand(position, velocity, acceleration, jerk float64) *types.Command {
    // Implementation
}
```

### 3. Use Examples

**✅ Good:**
```go
// Example usage:
//   table := CreatePickAndPlaceTable()
//   err := manager.StartExecution(ctx, table)
//   if err != nil {
//       log.Fatal(err)
//   }
func CreatePickAndPlaceTable() *types.CommandTable {
    // Implementation
}
```

## Common Pitfalls

### 1. Forgetting to Wait for Completion

**❌ Bad:**
```go
err := manager.StartExecution(ctx, table)
if err != nil {
    return err
}
// Execution might not be complete yet
return nil
```

**✅ Good:**
```go
err := manager.StartExecution(ctx, table)
if err != nil {
    return err
}

// Wait for completion
for {
    status := manager.GetExecutionStatus()
    if status.State == types.StateCompleted {
        break
    }
    if status.State == types.StateError {
        return status.Error
    }
    time.Sleep(100 * time.Millisecond)
}
```

### 2. Not Handling Context Cancellation

**❌ Bad:**
```go
for {
    status := manager.GetExecutionStatus()
    if status.State == types.StateCompleted {
        break
    }
    time.Sleep(100 * time.Millisecond)
}
```

**✅ Good:**
```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        status := manager.GetExecutionStatus()
        if status.State == types.StateCompleted {
            break
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

### 3. Mixing Units Without Conversion

**❌ Bad:**
```go
command := types.NewCommandBuilder().
    WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitMillimeters)).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    Build()
```

**✅ Good:**
```go
positionCounts := converter.ConvertPositionValue(
    types.NewPositionValue(100.0, types.PositionUnitMillimeters),
    types.PositionUnitCounts,
)

command := types.NewCommandBuilder().
    WithParameter("position", positionCounts).
    WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
    Build()
```

### 4. Not Validating Input Parameters

**❌ Bad:**
```go
func CreateCommand(position float64) *types.Command {
    return types.NewCommandBuilder().
        WithParameter("position", types.NewPositionValue(position, types.PositionUnitCounts)).
        Build()
}
```

**✅ Good:**
```go
func CreateCommand(position float64) (*types.Command, error) {
    if position < 0 {
        return nil, fmt.Errorf("position must be non-negative, got %f", position)
    }
    
    return types.NewCommandBuilder().
        WithParameter("position", types.NewPositionValue(position, types.PositionUnitCounts)).
        Build(), nil
}
```

## Troubleshooting

### 1. Common Error Messages

**"command validation failed: command type cannot be unknown"**
- Cause: Command type is not set or invalid
- Solution: Use `WithType()` method with valid command type

**"missing or invalid parameter: parameter X not found"**
- Cause: Required parameter is missing
- Solution: Add the missing parameter using `WithParameter()`

**"unit conversion failed: invalid unit"**
- Cause: Invalid unit type used
- Solution: Use valid unit constants from `types` package

**"precondition failed: drive not ready"**
- Cause: Drive is not in ready state
- Solution: Check drive state and ensure it's ready before motion

### 2. Debugging Tips

**Enable Debug Logging:**
```go
import "log"

log.SetLevel(log.DebugLevel)
```

**Check Execution Status:**
```go
status := manager.GetExecutionStatus()
fmt.Printf("State: %s, Current Command: %d, Error: %v\n", 
    status.State, status.CurrentCommand, status.Error)
```

**Validate Tables:**
```go
err := manager.ValidateTable(table)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

### 3. Performance Issues

**Slow Execution:**
- Check if commands are waiting for conditions that never become true
- Verify timeout values are appropriate
- Ensure drive controller is responding properly

**High CPU Usage:**
- Increase polling interval for status monitoring
- Use context cancellation to stop unnecessary operations
- Check for infinite loops in command sequences

**Memory Issues:**
- Avoid creating large numbers of command tables
- Use appropriate data types for parameters
- Clean up resources when done

### 4. Getting Help

1. Check the API documentation
2. Look at example code in the `examples/` directory
3. Run tests to see expected behavior
4. Check error messages for specific guidance
5. Verify your drive controller implementation

This guide should help you avoid common issues and write robust, maintainable code with the Stage LinMot CT library.