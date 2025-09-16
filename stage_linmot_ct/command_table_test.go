package stage_linmot_ct

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// MockDriveController for testing
type MockDriveController struct {
	position       float64
	velocity       float64
	force          float64
	digitalInputs  map[int]bool
	analogInputs   map[int]float64
	digitalOutputs map[int]bool
	analogOutputs  map[int]float64
	driveState     types.DriveState
	motionComplete bool
	error          error
}

func NewMockDriveController() *MockDriveController {
	return &MockDriveController{
		digitalInputs:  make(map[int]bool),
		analogInputs:   make(map[int]float64),
		digitalOutputs: make(map[int]bool),
		analogOutputs:  make(map[int]float64),
		driveState:     types.DriveStateReady,
		motionComplete: true,
	}
}

func (mdc *MockDriveController) SetError(err error) {
	mdc.error = err
}

// Implement DriveController interface
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
	return nil
}

func (mdc *MockDriveController) WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error {
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
	mdc.driveState = types.DriveStateReady
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

func (mdc *MockDriveController) GetDigitalOutput(ctx context.Context, output int) (bool, error) {
	if mdc.error != nil {
		return false, mdc.error
	}
	return mdc.digitalOutputs[output], nil
}

func (mdc *MockDriveController) GetAnalogOutput(ctx context.Context, output int) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.analogOutputs[output], nil
}

func (mdc *MockDriveController) SetDigitalOutput(ctx context.Context, output int, value bool) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.digitalOutputs[output] = value
	return nil
}

func (mdc *MockDriveController) ClearDigitalOutput(ctx context.Context, output int) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.digitalOutputs[output] = false
	return nil
}

func (mdc *MockDriveController) SetAnalogOutput(ctx context.Context, output int, value float64) error {
	if mdc.error != nil {
		return mdc.error
	}
	mdc.analogOutputs[output] = value
	return nil
}

func (mdc *MockDriveController) WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	return nil
}

func TestNewCommandTableManager(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	if manager == nil {
		t.Fatal("Expected non-nil CommandTableManager")
	}

	if manager.driveController != driveController {
		t.Error("Expected driveController to be set")
	}

	if manager.unitConverter != unitConverter {
		t.Error("Expected unitConverter to be set")
	}

	if manager.conditionEvaluator != conditionEvaluator {
		t.Error("Expected conditionEvaluator to be set")
	}

	if manager.safetyGuard != safetyGuard {
		t.Error("Expected safetyGuard to be set")
	}
}

func TestCommandTableManager_CreateTable(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())

	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	if table == nil {
		t.Fatal("Expected non-nil CommandTable")
	}

	if table.ID != "test-table" {
		t.Errorf("Expected ID 'test-table', got '%s'", table.ID)
	}

	if table.Name != "Test Table" {
		t.Errorf("Expected Name 'Test Table', got '%s'", table.Name)
	}

	if table.Description != "A test command table" {
		t.Errorf("Expected Description 'A test command table', got '%s'", table.Description)
	}
}

func TestCommandTableManager_AddCommand(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	err := manager.AddCommand(table, command)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(table.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(table.Commands))
	}

	if table.Commands[0].ID != 1 {
		t.Errorf("Expected command ID 1, got %d", table.Commands[0].ID)
	}
}

func TestCommandTableManager_AddCommand_DuplicateID(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command1 := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	command2 := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveRelative).
		WithParameter("distance", types.NewPositionValue(50.0, types.PositionUnitCounts)).
		Build()

	err := manager.AddCommand(table, command1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = manager.AddCommand(table, command2)
	if err == nil {
		t.Error("Expected error for duplicate command ID")
	}
}

func TestCommandTableManager_RemoveCommand(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	err := manager.RemoveCommand(table, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(table.Commands) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(table.Commands))
	}
}

func TestCommandTableManager_RemoveCommand_NotFound(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	err := manager.RemoveCommand(table, 999)
	if err == nil {
		t.Error("Expected error for non-existent command ID")
	}
}

func TestCommandTableManager_UpdateCommand(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	updatedCommand := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(200.0, types.PositionUnitCounts)).
		Build()

	err := manager.UpdateCommand(table, 1, updatedCommand)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(table.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(table.Commands))
	}

	// Check that the position was updated
	position, err := types.NewParameterExtractor().ExtractPosition(table.Commands[0].Parameters, "position")
	if err != nil {
		t.Fatalf("Expected no error extracting position, got %v", err)
	}

	if position.Value != 200.0 {
		t.Errorf("Expected position 200.0, got %f", position.Value)
	}
}

func TestCommandTableManager_GetCommand(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	retrievedCommand, err := manager.GetCommand(table, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedCommand.ID != 1 {
		t.Errorf("Expected command ID 1, got %d", retrievedCommand.ID)
	}
}

func TestCommandTableManager_GetCommand_NotFound(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	_, err := manager.GetCommand(table, 999)
	if err == nil {
		t.Error("Expected error for non-existent command ID")
	}
}

func TestCommandTableManager_ValidateTable(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	err := manager.ValidateTable(table)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommandTableManager_ValidateTable_EmptyTable(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	err := manager.ValidateTable(table)
	if err == nil {
		t.Error("Expected error for empty table")
	}
}

func TestCommandTableManager_StartExecution(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	status := manager.GetExecutionStatus()
	if status.State != types.StateRunning {
		t.Errorf("Expected state Running, got %s", status.State.String())
	}
}

func TestCommandTableManager_StartExecution_AlreadyRunning(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	manager.StartExecution(ctx, table)

	// Try to start again
	err := manager.StartExecution(ctx, table)
	if err == nil {
		t.Error("Expected error for already running execution")
	}
}

func TestCommandTableManager_PauseExecution(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	manager.StartExecution(ctx, table)

	err := manager.PauseExecution()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	status := manager.GetExecutionStatus()
	if status.State != types.StatePaused {
		t.Errorf("Expected state Paused, got %s", status.State.String())
	}
}

func TestCommandTableManager_ResumeExecution(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	manager.StartExecution(ctx, table)
	manager.PauseExecution()

	err := manager.ResumeExecution()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	status := manager.GetExecutionStatus()
	if status.State != types.StateRunning {
		t.Errorf("Expected state Running, got %s", status.State.String())
	}
}

func TestCommandTableManager_StopExecution(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	manager.StartExecution(ctx, table)

	err := manager.StopExecution()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	status := manager.GetExecutionStatus()
	if status.State != types.StateStopped {
		t.Errorf("Expected state Stopped, got %s", status.State.String())
	}
}

func TestCommandTableManager_GetCurrentCommand(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())
	table := manager.CreateTable("test-table", "Test Table", "A test command table")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()
	manager.StartExecution(ctx, table)

	currentCommand := manager.GetCurrentCommand()
	if currentCommand == nil {
		t.Error("Expected non-nil current command")
	}

	if currentCommand.ID != 1 {
		t.Errorf("Expected current command ID 1, got %d", currentCommand.ID)
	}
}

func TestCommandTableManager_GetVariables(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())

	variables := manager.GetVariables()
	if variables == nil {
		t.Error("Expected non-nil variables map")
	}

	if len(variables) != 0 {
		t.Errorf("Expected empty variables map, got %d variables", len(variables))
	}
}

func TestCommandTableManager_SetVariable(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())

	manager.SetVariable("testVar", 42.0)

	variables := manager.GetVariables()
	if len(variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variables))
	}

	if variables["testVar"] != 42.0 {
		t.Errorf("Expected variable value 42.0, got %v", variables["testVar"])
	}
}

func TestCommandTableManager_GetExecutionHistory(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())

	history := manager.GetExecutionHistory()
	if history == nil {
		t.Error("Expected non-nil execution history")
	}

	if len(history) != 0 {
		t.Errorf("Expected empty execution history, got %d results", len(history))
	}
}

func TestCommandTableManager_ClearHistory(t *testing.T) {
	manager := NewCommandTableManager(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard())

	// This should not panic
	manager.ClearHistory()
}

func TestNewCommandTableBuilder(t *testing.T) {
	builder := NewCommandTableBuilder()
	if builder == nil {
		t.Fatal("Expected non-nil CommandTableBuilder")
	}
}

func TestCommandTableBuilder_WithID(t *testing.T) {
	builder := NewCommandTableBuilder()
	builder = builder.WithID("test-id")

	if builder.id != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", builder.id)
	}
}

func TestCommandTableBuilder_WithName(t *testing.T) {
	builder := NewCommandTableBuilder()
	builder = builder.WithName("Test Name")

	if builder.name != "Test Name" {
		t.Errorf("Expected Name 'Test Name', got '%s'", builder.name)
	}
}

func TestCommandTableBuilder_WithDescription(t *testing.T) {
	builder := NewCommandTableBuilder()
	builder = builder.WithDescription("Test Description")

	if builder.description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got '%s'", builder.description)
	}
}

func TestCommandTableBuilder_WithCommand(t *testing.T) {
	builder := NewCommandTableBuilder()
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		Build()

	builder = builder.WithCommand(command)

	if len(builder.commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(builder.commands))
	}
}

func TestCommandTableBuilder_WithVariable(t *testing.T) {
	builder := NewCommandTableBuilder()
	builder = builder.WithVariable("testVar", 42.0)

	if len(builder.variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(builder.variables))
	}

	if builder.variables["testVar"] != 42.0 {
		t.Errorf("Expected variable value 42.0, got %v", builder.variables["testVar"])
	}
}

func TestCommandTableBuilder_Build(t *testing.T) {
	builder := NewCommandTableBuilder().
		WithID("test-id").
		WithName("Test Name").
		WithDescription("Test Description")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		Build()

	builder = builder.WithCommand(command)

	table := builder.Build()

	if table == nil {
		t.Fatal("Expected non-nil CommandTable")
	}

	if table.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", table.ID)
	}

	if table.Name != "Test Name" {
		t.Errorf("Expected Name 'Test Name', got '%s'", table.Name)
	}

	if table.Description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got '%s'", table.Description)
	}

	if len(table.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(table.Commands))
	}
}

func TestNewInMemoryCommandTableRepository(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()
	if repo == nil {
		t.Fatal("Expected non-nil InMemoryCommandTableRepository")
	}
}

func TestInMemoryCommandTableRepository_Save(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()
	table := &types.CommandTable{
		ID:   "test-table",
		Name: "Test Table",
	}

	err := repo.Save(table)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	exists := repo.Exists("test-table")
	if !exists {
		t.Error("Expected table to exist")
	}
}

func TestInMemoryCommandTableRepository_Load(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()
	table := &types.CommandTable{
		ID:   "test-table",
		Name: "Test Table",
	}

	repo.Save(table)

	loadedTable, err := repo.Load("test-table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if loadedTable.ID != "test-table" {
		t.Errorf("Expected ID 'test-table', got '%s'", loadedTable.ID)
	}
}

func TestInMemoryCommandTableRepository_Load_NotFound(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()

	_, err := repo.Load("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

func TestInMemoryCommandTableRepository_Delete(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()
	table := &types.CommandTable{
		ID:   "test-table",
		Name: "Test Table",
	}

	repo.Save(table)
	repo.Delete("test-table")

	exists := repo.Exists("test-table")
	if exists {
		t.Error("Expected table to not exist after deletion")
	}
}

func TestInMemoryCommandTableRepository_List(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()
	table1 := &types.CommandTable{ID: "table1", Name: "Table 1"}
	table2 := &types.CommandTable{ID: "table2", Name: "Table 2"}

	repo.Save(table1)
	repo.Save(table2)

	tables := repo.List()
	if len(tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(tables))
	}
}

func TestInMemoryCommandTableRepository_Exists(t *testing.T) {
	repo := NewInMemoryCommandTableRepository()

	exists := repo.Exists("non-existent")
	if exists {
		t.Error("Expected non-existent table to not exist")
	}

	table := &types.CommandTable{ID: "test-table", Name: "Test Table"}
	repo.Save(table)

	exists = repo.Exists("test-table")
	if !exists {
		t.Error("Expected existing table to exist")
	}
}

func TestNewCommandTableService(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()
	repository := NewInMemoryCommandTableRepository()

	service := NewCommandTableService(driveController, unitConverter, conditionEvaluator, safetyGuard, repository)

	if service == nil {
		t.Fatal("Expected non-nil CommandTableService")
	}
}

func TestCommandTableService_CreateTable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	table, err := service.CreateTable("test-table", "Test Table", "A test command table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if table == nil {
		t.Fatal("Expected non-nil CommandTable")
	}

	if table.ID != "test-table" {
		t.Errorf("Expected ID 'test-table', got '%s'", table.ID)
	}
}

func TestCommandTableService_GetTable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Create a table first
	service.CreateTable("test-table", "Test Table", "A test command table")

	table, err := service.GetTable("test-table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if table.ID != "test-table" {
		t.Errorf("Expected ID 'test-table', got '%s'", table.ID)
	}
}

func TestCommandTableService_GetTable_NotFound(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	_, err := service.GetTable("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

func TestCommandTableService_UpdateTable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Create a table first
	service.CreateTable("test-table", "Test Table", "A test command table")

	updatedTable := &types.CommandTable{
		ID:          "test-table",
		Name:        "Updated Test Table",
		Description: "An updated test command table",
	}

	err := service.UpdateTable(updatedTable)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the update
	table, err := service.GetTable("test-table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if table.Name != "Updated Test Table" {
		t.Errorf("Expected Name 'Updated Test Table', got '%s'", table.Name)
	}
}

func TestCommandTableService_DeleteTable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Create a table first
	service.CreateTable("test-table", "Test Table", "A test command table")

	err := service.DeleteTable("test-table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = service.GetTable("test-table")
	if err == nil {
		t.Error("Expected error for deleted table")
	}
}

func TestCommandTableService_ListTables(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Create some tables
	service.CreateTable("table1", "Table 1", "First table")
	service.CreateTable("table2", "Table 2", "Second table")

	tables := service.ListTables()
	if len(tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(tables))
	}
}

func TestCommandTableService_ExecuteTable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Create a table with a command
	table := &types.CommandTable{
		ID:   "test-table",
		Name: "Test Table",
		Commands: []types.Command{
			*types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
				Build(),
		},
	}

	service.repository.Save(table)

	ctx := context.Background()
	err := service.ExecuteTable(ctx, "test-table")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommandTableService_ExecuteTable_NotFound(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	ctx := context.Background()
	err := service.ExecuteTable(ctx, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

func TestCommandTableService_GetExecutionStatus(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	status := service.GetExecutionStatus()
	if status == nil {
		t.Error("Expected non-nil execution status")
	}
}

func TestCommandTableService_ControlExecution(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	// Test pause
	err := service.Pause()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test resume
	err = service.Resume()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test stop
	err = service.Stop()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommandTableService_Start(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	ctx := context.Background()
	err := service.Start(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCommandTableService_GetStatus(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	status := service.GetStatus()
	if status == nil {
		t.Error("Expected non-nil status")
	}
}

func TestCommandTableService_GetCurrentCommand(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	command := service.GetCurrentCommand()
	// This might be nil if no execution is running, which is fine
	_ = command
}

func TestCommandTableService_GetVariables(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	variables := service.GetVariables()
	if variables == nil {
		t.Error("Expected non-nil variables map")
	}
}

func TestCommandTableService_SetVariable(t *testing.T) {
	service := NewCommandTableService(NewMockDriveController(), types.NewUnitConverter(), types.NewDefaultConditionEvaluator(), NewSafetyGuard(), NewInMemoryCommandTableRepository())

	service.SetVariable("testVar", 42.0)

	variables := service.GetVariables()
	if len(variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variables))
	}

	if variables["testVar"] != 42.0 {
		t.Errorf("Expected variable value 42.0, got %v", variables["testVar"])
	}
}