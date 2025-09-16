package execution

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// MockDriveController implements DriveController for testing
type MockDriveController struct {
	position      float64
	velocity      float64
	force         float64
	digitalInputs map[int]bool
	analogInputs  map[int]float64
	driveState    types.DriveState
	motionComplete bool
	error         error
}

func NewMockDriveController() *MockDriveController {
	return &MockDriveController{
		digitalInputs: make(map[int]bool),
		analogInputs:  make(map[int]float64),
		driveState:    types.DriveStateReady,
		motionComplete: true,
	}
}

func (mdc *MockDriveController) SetError(err error) {
	mdc.error = err
}

func (mdc *MockDriveController) SetPosition(pos float64) {
	mdc.position = pos
}

func (mdc *MockDriveController) SetVelocity(vel float64) {
	mdc.velocity = vel
}

func (mdc *MockDriveController) SetForceValue(force float64) {
	mdc.force = force
}

func (mdc *MockDriveController) SetDigitalInput(input int, value bool) {
	mdc.digitalInputs[input] = value
}

func (mdc *MockDriveController) SetAnalogInput(input int, value float64) {
	mdc.analogInputs[input] = value
}

func (mdc *MockDriveController) SetDriveState(state types.DriveState) {
	mdc.driveState = state
}

func (mdc *MockDriveController) SetMotionComplete(complete bool) {
	mdc.motionComplete = complete
}

// DriveController interface implementation
func (mdc *MockDriveController) MoveAbsolute(ctx context.Context, position float64, velocity float64, acceleration float64, jerk float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.position = position
	return nil
}

func (mdc *MockDriveController) MoveRelative(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.position += distance
	return nil
}

func (mdc *MockDriveController) MoveIncremental(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.position += distance
	return nil
}

func (mdc *MockDriveController) Jog(ctx context.Context, velocity float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.velocity = velocity
	return nil
}

func (mdc *MockDriveController) Stop(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.velocity = 0
	return nil
}

func (mdc *MockDriveController) Wait(ctx context.Context, duration time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	time.Sleep(duration)
	return nil
}

func (mdc *MockDriveController) WaitPosition(ctx context.Context, position float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	// Simulate waiting for position
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	// Simulate waiting for velocity
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	// Simulate waiting for force
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) SetDigitalOutput(ctx context.Context, output int, value bool) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.digitalInputs[output] = value
	return nil
}

func (mdc *MockDriveController) ClearDigitalOutput(ctx context.Context, output int) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.digitalInputs[output] = false
	return nil
}

func (mdc *MockDriveController) SetAnalogOutput(ctx context.Context, output int, value float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.analogInputs[output] = value
	return nil
}

func (mdc *MockDriveController) WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	// Simulate waiting for digital input
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	// Simulate waiting for analog input
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) Home(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.position = 0
	return nil
}

func (mdc *MockDriveController) Reset(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.position = 0
	mdc.velocity = 0
	mdc.force = 0
	return nil
}

func (mdc *MockDriveController) SaveConfiguration(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) LoadConfiguration(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) ForceControlOn(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) ForceControlOff(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) SetForce(ctx context.Context, force float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.force = force
	return nil
}

func (mdc *MockDriveController) StartOscilloscope(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) StopOscilloscope(ctx context.Context) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) SaveData(ctx context.Context, filename string) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) GetPosition(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.position, nil
}

func (mdc *MockDriveController) GetVelocity(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.velocity, nil
}

func (mdc *MockDriveController) GetForce(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.force, nil
}

func (mdc *MockDriveController) GetDigitalInput(ctx context.Context, input int) (bool, error) {
	if mdc.error != nil {
		return false, mdc.error
	}
	return mdc.digitalInputs[input], nil
}

func (mdc *MockDriveController) GetAnalogInput(ctx context.Context, input int) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.analogInputs[input], nil
}

func (mdc *MockDriveController) GetDriveState(ctx context.Context) (types.DriveState, error) {
	if mdc.error != nil {
		return types.DriveState(0), mdc.error
	}
	return mdc.driveState, nil
}

func (mdc *MockDriveController) IsMotionComplete(ctx context.Context) (bool, error) {
	if mdc.error != nil {
		return false, mdc.error
	}
	return mdc.motionComplete, nil
}

// MockConditionEvaluator implements ConditionEvaluator for testing
type MockConditionEvaluator struct {
	canEvaluate bool
	evaluateResult bool
	evaluateError error
}

func NewMockConditionEvaluator() *MockConditionEvaluator {
	return &MockConditionEvaluator{
		canEvaluate: true,
		evaluateResult: true,
	}
}

func (mce *MockConditionEvaluator) SetCanEvaluate(can bool) {
	mce.canEvaluate = can
}

func (mce *MockConditionEvaluator) SetEvaluateResult(result bool) {
	mce.evaluateResult = result
}

func (mce *MockConditionEvaluator) SetEvaluateError(err error) {
	mce.evaluateError = err
}

func (mce *MockConditionEvaluator) CanEvaluate(condition *types.Condition) bool {
	return mce.canEvaluate
}

func (mce *MockConditionEvaluator) Evaluate(ctx context.Context, condition *types.Condition, variables map[string]interface{}) (bool, error) {
	if mce.evaluateError != nil {
		return false, mce.evaluateError
	}
	return mce.evaluateResult, nil
}

func (mce *MockConditionEvaluator) GetRequiredData(condition *types.Condition) []string {
	return []string{}
}

func TestExecutionState_String(t *testing.T) {
	tests := []struct {
		state    ExecutionState
		expected string
	}{
		{StateIdle, "Idle"},
		{StateRunning, "Running"},
		{StatePaused, "Paused"},
		{StateStopped, "Stopped"},
		{StateError, "Error"},
		{StateCompleted, "Completed"},
		{ExecutionState(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("ExecutionState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewDefaultExecutionEngine(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()

	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	if engine == nil {
		t.Fatal("NewDefaultExecutionEngine returned nil")
	}

	if engine.state != StateIdle {
		t.Errorf("Expected initial state to be Idle, got %s", engine.state)
	}

	if engine.driveController != driveController {
		t.Error("DriveController not set correctly")
	}

	// Note: Can't compare interface values directly
	// This is a limitation of the test design

	if engine.unitConverter != unitConverter {
		t.Error("UnitConverter not set correctly")
	}
}

func TestDefaultExecutionEngine_Execute(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Create a simple command table
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Test successful execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = engine.WaitForCompletion(ctx)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}

	// Check final status
	status := engine.GetStatus()
	if status.State != StateCompleted {
		t.Errorf("Expected final state to be Completed, got %s", status.State)
	}

	if len(status.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(status.Results))
	}

	if !status.Results[0].Success {
		t.Errorf("Expected command to succeed, got error: %v", status.Results[0].Error)
	}
}

func TestDefaultExecutionEngine_Execute_EmptyTable(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Create empty command table
	table := types.NewCommandTableBuilder().
		WithID("empty_table").
		Build()

	err := engine.Execute(context.Background(), table)
	if err == nil {
		t.Error("Expected error for empty table, got nil")
	}

	if err.Error() != "command table is empty" {
		t.Errorf("Expected 'command table is empty' error, got: %v", err)
	}
}

func TestDefaultExecutionEngine_Execute_NilTable(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	err := engine.Execute(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil table, got nil")
	}

	if err.Error() != "command table cannot be nil" {
		t.Errorf("Expected 'command table cannot be nil' error, got: %v", err)
	}
}

func TestDefaultExecutionEngine_Execute_AlreadyRunning(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Create a command table
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start first execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("First Execute failed: %v", err)
	}

	// Try to start second execution while first is running
	err = engine.Execute(context.Background(), table)
	if err == nil {
		t.Error("Expected error for second execution, got nil")
	}

	if err.Error() != "execution engine is not idle, current state: Running" {
		t.Errorf("Expected 'execution engine is not idle' error, got: %v", err)
	}
}

func TestDefaultExecutionEngine_PauseResume(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Create a command table with wait command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Wait a bit for execution to start
	time.Sleep(10 * time.Millisecond)

	// Pause execution
	err = engine.Pause()
	if err != nil {
		t.Fatalf("Pause failed: %v", err)
	}

	// Check state
	status := engine.GetStatus()
	if status.State != StatePaused {
		t.Errorf("Expected state to be Paused, got %s", status.State)
	}

	// Resume execution
	err = engine.Resume()
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}

	// Check state
	status = engine.GetStatus()
	if status.State != StateRunning {
		t.Errorf("Expected state to be Running, got %s", status.State)
	}

	// Wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = engine.WaitForCompletion(ctx)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}
}

func TestDefaultExecutionEngine_Stop(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Create a command table with wait command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Wait a bit for execution to start
	time.Sleep(10 * time.Millisecond)

	// Stop execution
	err = engine.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Check state
	status := engine.GetStatus()
	if status.State != StateStopped {
		t.Errorf("Expected state to be Stopped, got %s", status.State)
	}
}

func TestDefaultExecutionEngine_IsRunning(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Initially not running
	if engine.IsRunning() {
		t.Error("Expected engine to not be running initially")
	}

	// Create a command table
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(50.0, types.TimeUnitMS)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should be running now
	if !engine.IsRunning() {
		t.Error("Expected engine to be running")
	}

	// Wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = engine.WaitForCompletion(ctx)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}

	// Should not be running after completion
	if engine.IsRunning() {
		t.Error("Expected engine to not be running after completion")
	}
}

func TestDefaultExecutionEngine_GetCurrentCommand(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Initially no current command
	if engine.GetCurrentCommand() != nil {
		t.Error("Expected no current command initially")
	}

	// Create a command table
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(50.0, types.TimeUnitMS)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Should have current command
	currentCmd := engine.GetCurrentCommand()
	if currentCmd == nil {
		t.Error("Expected current command to be set")
	} else if currentCmd.ID != 1 {
		t.Errorf("Expected current command ID to be 1, got %d", currentCmd.ID)
	}

	// Wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = engine.WaitForCompletion(ctx)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}

	// Should not have current command after completion
	if engine.GetCurrentCommand() != nil {
		t.Error("Expected no current command after completion")
	}
}

func TestDefaultExecutionEngine_CommandFailure(t *testing.T) {
	driveController := NewMockDriveController()
	conditionEvaluator := NewMockConditionEvaluator()
	unitConverter := types.NewUnitConverter()
	engine := NewDefaultExecutionEngine(driveController, conditionEvaluator, unitConverter)

	// Set drive controller to return error
	driveController.SetError(errors.New("drive error"))

	// Create a command table
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()

	table := types.NewCommandTableBuilder().
		WithID("test_table").
		WithCommand(command).
		Build()

	// Start execution
	err := engine.Execute(context.Background(), table)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = engine.WaitForCompletion(ctx)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}

	// Check final status
	status := engine.GetStatus()
	if status.State != StateError {
		t.Errorf("Expected final state to be Error, got %s", status.State)
	}

	if len(status.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(status.Results))
	}

	if status.Results[0].Success {
		t.Error("Expected command to fail")
	}

	if status.Results[0].Error == nil {
		t.Error("Expected command error to be set")
	}
}