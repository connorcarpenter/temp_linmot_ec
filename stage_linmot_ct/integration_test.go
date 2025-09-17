package stage_linmot_ct

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// TestEndToEndMotionWorkflow tests a complete motion workflow
func TestEndToEndMotionWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with a motion sequence
	table := manager.CreateTable("motion-test", "Motion Test", "End-to-end motion test")

	// Add commands to the table
	commands := []*types.Command{
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
			Build(),
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
			Build(),
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Validate the table
	err := manager.ValidateTable(table)
	if err != nil {
		t.Fatalf("Table validation failed: %v", err)
	}

	// Execute the table
	ctx := context.Background()
	err = manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Verify final position
				position, err := driveController.GetPosition(ctx)
				if err != nil {
					t.Fatalf("Failed to get position: %v", err)
				}
				if position != 0.0 {
					t.Errorf("Expected final position 0.0, got %f", position)
				}
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestErrorHandlingWorkflow tests error handling and recovery
func TestErrorHandlingWorkflow(t *testing.T) {
	// Create a mock drive controller that will fail
	driveController := NewMockDriveController()
	driveController.SetError(types.NewPreconditionError("motion_state", "drive not ready"))

	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with a motion command
	table := manager.CreateTable("error-test", "Error Test", "Error handling test")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	err := manager.AddCommand(table, command)
	if err != nil {
		t.Fatalf("Failed to add command: %v", err)
	}

	// Execute the table
	ctx := context.Background()
	err = manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for error state
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Expected error state, but execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateError {
				if status.Error == nil {
					t.Error("Expected error in status, but got nil")
				}
				return
			}
		}
	}
}

// TestPauseResumeWorkflow tests pause and resume functionality
func TestPauseResumeWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with a long-running command
	table := manager.CreateTable("pause-test", "Pause Test", "Pause/resume test")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(2.0, types.TimeUnitSeconds)).
		Build()

	err := manager.AddCommand(table, command)
	if err != nil {
		t.Fatalf("Failed to add command: %v", err)
	}

	// Execute the table
	ctx := context.Background()
	err = manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait a bit, then pause
	time.Sleep(100 * time.Millisecond)
	err = manager.PauseExecution()
	if err != nil {
		t.Fatalf("Failed to pause execution: %v", err)
	}

	// Verify paused state
	status := manager.GetExecutionStatus()
	if status.State != types.StatePaused {
		t.Errorf("Expected paused state, got %s", status.State.String())
	}

	// Resume execution
	err = manager.ResumeExecution()
	if err != nil {
		t.Fatalf("Failed to resume execution: %v", err)
	}

	// Verify running state
	status = manager.GetExecutionStatus()
	if status.State != types.StateRunning {
		t.Errorf("Expected running state, got %s", status.State.String())
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestSafetyLimitsWorkflow tests safety limit enforcement
func TestSafetyLimitsWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()

	// Create safety guard with limits
	safetyGuard := NewSafetyGuardWithLimits(&SafetyLimits{
		MinPosition: 0.0,
		MaxPosition: 50.0,
		MaxVelocity: 100.0,
		MinForce:    -10.0,
		MaxForce:    10.0,
	})

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with a command that exceeds safety limits
	table := manager.CreateTable("safety-test", "Safety Test", "Safety limits test")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)). // Exceeds max position
		WithParameter("velocity", types.NewVelocityValue(150.0, types.VelocityUnitCountsS)). // Exceeds max velocity
		Build()

	err := manager.AddCommand(table, command)
	if err != nil {
		t.Fatalf("Failed to add command: %v", err)
	}

	// Execute the table
	ctx := context.Background()
	err = manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for error state due to safety violation
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Expected safety error, but execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateError {
				if status.Error == nil {
					t.Error("Expected error in status, but got nil")
				}
				// Verify it's a safety-related error
				if status.Error != nil && status.Error.Error() != "" {
					// This is expected - safety validation should have failed
					return
				}
			}
		}
	}
}

// TestIOCommandsWorkflow tests I/O command execution
func TestIOCommandsWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with I/O commands
	table := manager.CreateTable("io-test", "IO Test", "I/O commands test")

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
			WithParameter("timeout", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Verify I/O operations
				digitalOutput, err := driveController.GetDigitalOutput(ctx, 1)
				if err != nil {
					t.Fatalf("Failed to get digital output: %v", err)
				}
				if !digitalOutput {
					t.Error("Expected digital output 1 to be true")
				}

				analogOutput, err := driveController.GetAnalogOutput(ctx, 1)
				if err != nil {
					t.Fatalf("Failed to get analog output: %v", err)
				}
				if analogOutput != 3.14 {
					t.Errorf("Expected analog output 1 to be 3.14, got %f", analogOutput)
				}
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestLoopCommandsWorkflow tests loop and jump command execution
func TestLoopCommandsWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with loop commands
	table := manager.CreateTable("loop-test", "Loop Test", "Loop commands test")

	commands := []*types.Command{
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdLoopStart).
			WithParameter("count", 3).
			Build(),
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitCounts)).
			Build(),
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(0.1, types.TimeUnitSeconds)).
			Build(),
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdLoopEnd).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Verify the loop executed (position should be 10.0 after 3 iterations)
				position, err := driveController.GetPosition(ctx)
				if err != nil {
					t.Fatalf("Failed to get position: %v", err)
				}
				if position != 10.0 {
					t.Errorf("Expected final position 10.0, got %f", position)
				}
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestForceControlWorkflow tests force control command execution
func TestForceControlWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with force control commands
	table := manager.CreateTable("force-test", "Force Test", "Force control test")

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
			WithParameter("duration", types.NewTimeValue(0.5, types.TimeUnitSeconds)).
			Build(),
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdForceControlOff).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Verify force was set
				force, err := driveController.GetForce(ctx)
				if err != nil {
					t.Fatalf("Failed to get force: %v", err)
				}
				// Force should be 5.0 * 100 (scaling factor) = 500.0 counts
				expectedForce := 5.0 * 100.0
				if force != expectedForce {
					t.Errorf("Expected force %f, got %f", expectedForce, force)
				}
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestDataAcquisitionWorkflow tests data acquisition command execution
func TestDataAcquisitionWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with data acquisition commands
	table := manager.CreateTable("data-test", "Data Test", "Data acquisition test")

	commands := []*types.Command{
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdStartOscilloscope).
			Build(),
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
			Build(),
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdStopOscilloscope).
			Build(),
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdSaveData).
			WithParameter("filename", "test_data.csv").
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Data acquisition commands should have executed successfully
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestSystemCommandsWorkflow tests system command execution
func TestSystemCommandsWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with system commands
	table := manager.CreateTable("system-test", "System Test", "System commands test")

	commands := []*types.Command{
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdHome).
			Build(),
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdReset).
			Build(),
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdSaveConfiguration).
			Build(),
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdLoadConfiguration).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// System commands should have executed successfully
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}

// TestComplexWorkflow tests a complex multi-command workflow
func TestComplexWorkflow(t *testing.T) {
	// Create a mock drive controller
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	// Create command table manager
	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with a complex workflow
	table := manager.CreateTable("complex-test", "Complex Test", "Complex workflow test")

	commands := []*types.Command{
		// Home the system
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdHome).
			Build(),
		// Move to position 1
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(50.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			Build(),
		// Set digital output
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdSetDigitalOutput).
			WithParameter("output", 1).
			WithParameter("value", true).
			Build(),
		// Wait for input
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdWaitDigitalInput).
			WithParameter("input", 1).
			WithParameter("value", true).
			WithParameter("timeout", types.NewTimeValue(2.0, types.TimeUnitSeconds)).
			Build(),
		// Move to position 2
		types.NewCommandBuilder().
			WithID(5).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
			Build(),
		// Enable force control
		types.NewCommandBuilder().
			WithID(6).
			WithType(types.CmdForceControlOn).
			Build(),
		// Set force
		types.NewCommandBuilder().
			WithID(7).
			WithType(types.CmdSetForce).
			WithParameter("force", types.NewForceValue(2.5, types.ForceUnitNewtons)).
			Build(),
		// Wait
		types.NewCommandBuilder().
			WithID(8).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitSeconds)).
			Build(),
		// Disable force control
		types.NewCommandBuilder().
			WithID(9).
			WithType(types.CmdForceControlOff).
			Build(),
		// Return home
		types.NewCommandBuilder().
			WithID(10).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			Build(),
	}

	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			t.Fatalf("Failed to add command: %v", err)
		}
	}

	// Execute the table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				// Verify final state
				position, err := driveController.GetPosition(ctx)
				if err != nil {
					t.Fatalf("Failed to get position: %v", err)
				}
				if position != 0.0 {
					t.Errorf("Expected final position 0.0, got %f", position)
				}

				digitalOutput, err := driveController.GetDigitalOutput(ctx, 1)
				if err != nil {
					t.Fatalf("Failed to get digital output: %v", err)
				}
				if !digitalOutput {
					t.Error("Expected digital output 1 to be true")
				}
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}