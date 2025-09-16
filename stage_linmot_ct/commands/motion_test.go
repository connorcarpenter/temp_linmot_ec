package commands

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
	digitalOutputs map[int]bool
	analogOutputs  map[int]float64
	driveState    types.DriveState
	motionComplete bool
	error         error
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
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
	time.Sleep(10 * time.Millisecond)
	return nil
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
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error {
	if mdc.error != nil {
		return mdc.error
	}
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

func TestNewMotionCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	if executor == nil {
		t.Fatal("NewMotionCommandExecutor returned nil")
	}
	
	if executor.driveController != driveController {
		t.Error("DriveController not set correctly")
	}
	
	if executor.unitConverter != unitConverter {
		t.Error("UnitConverter not set correctly")
	}
}

func TestMotionCommandExecutor_ExecuteMoveAbsolute(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Create a move absolute command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := executor.ExecuteMoveAbsolute(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteMoveAbsolute failed: %v", err)
	}
	
	// Verify the drive controller was called with correct values
	// Position should be converted from mm to counts (10.0 * 1000 = 10000)
	expectedPosition := 10000.0
	if driveController.position != expectedPosition {
		t.Errorf("Expected position %f, got %f", expectedPosition, driveController.position)
	}
}

func TestMotionCommandExecutor_ExecuteMoveRelative(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Set initial position
	driveController.SetPosition(5000.0) // 5mm in counts
	
	// Create a move relative command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveRelative).
		WithParameter("distance", types.NewPositionValue(5.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(2.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(5.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(50.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := executor.ExecuteMoveRelative(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteMoveRelative failed: %v", err)
	}
	
	// Verify the drive controller was called with correct values
	// Distance should be converted from mm to counts (5.0 * 1000 = 5000)
	expectedPosition := 5000.0 + 5000.0 // Initial + distance
	if driveController.position != expectedPosition {
		t.Errorf("Expected position %f, got %f", expectedPosition, driveController.position)
	}
}

func TestMotionCommandExecutor_ExecuteMoveIncremental(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Set initial position
	driveController.SetPosition(3000.0) // 3mm in counts
	
	// Create a move incremental command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveIncremental).
		WithParameter("distance", types.NewPositionValue(2.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(1.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(2.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(20.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := executor.ExecuteMoveIncremental(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteMoveIncremental failed: %v", err)
	}
	
	// Verify the drive controller was called with correct values
	// Distance should be converted from mm to counts (2.0 * 1000 = 2000)
	expectedPosition := 3000.0 + 2000.0 // Initial + distance
	if driveController.position != expectedPosition {
		t.Errorf("Expected position %f, got %f", expectedPosition, driveController.position)
	}
}

func TestMotionCommandExecutor_ExecuteJog(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Create a jog command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdJog).
		WithParameter("velocity", types.NewVelocityValue(3.0, types.VelocityUnitMMS)).
		Build()
	
	// Execute the command
	err := executor.ExecuteJog(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteJog failed: %v", err)
	}
	
	// Verify the drive controller was called with correct values
	// Velocity should be converted from mm/s to counts/s (3.0 * 1000 = 3000)
	expectedVelocity := 3000.0
	if driveController.velocity != expectedVelocity {
		t.Errorf("Expected velocity %f, got %f", expectedVelocity, driveController.velocity)
	}
}

func TestMotionCommandExecutor_ExecuteStop(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Set initial velocity
	driveController.SetVelocity(1000.0)
	
	// Create a stop command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdStop).
		Build()
	
	// Execute the command
	err := executor.ExecuteStop(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteStop failed: %v", err)
	}
	
	// Verify the drive controller was called
	if driveController.velocity != 0 {
		t.Errorf("Expected velocity 0, got %f", driveController.velocity)
	}
}

func TestMotionCommandExecutor_ValidateMotionCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid move absolute command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
				WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing position parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
				WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
				Build(),
			wantErr: true,
		},
		{
			name: "Invalid position parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", "invalid").
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
				WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid jog command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJog).
				WithParameter("velocity", types.NewVelocityValue(3.0, types.VelocityUnitMMS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing velocity parameter for jog",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJog).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid stop command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdStop).
				Build(),
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.ValidateMotionCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMotionCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMotionCommandExecutor_GetMotionCommandInfo(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	tests := []struct {
		commandType types.CommandType
		description string
		parameters  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdMoveAbsolute,
			description: "Move to absolute position",
			parameters:  []string{"position", "velocity", "acceleration", "jerk"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveRelative,
			description: "Move by relative distance",
			parameters:  []string{"distance", "velocity", "acceleration", "jerk"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveIncremental,
			description: "Move by fixed increment",
			parameters:  []string{"distance", "velocity", "acceleration", "jerk"},
			wantErr:     false,
		},
		{
			commandType: types.CmdJog,
			description: "Continuous motion",
			parameters:  []string{"velocity"},
			wantErr:     false,
		},
		{
			commandType: types.CmdStop,
			description: "Stop motion",
			parameters:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdWait, // Invalid for motion executor
			description: "",
			parameters:  nil,
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			description, parameters, err := executor.GetMotionCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMotionCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if description != tt.description {
					t.Errorf("GetMotionCommandInfo() description = %v, want %v", description, tt.description)
				}
				
				if len(parameters) != len(tt.parameters) {
					t.Errorf("GetMotionCommandInfo() parameters length = %v, want %v", len(parameters), len(tt.parameters))
				}
				
				for i, param := range parameters {
					if param != tt.parameters[i] {
						t.Errorf("GetMotionCommandInfo() parameters[%d] = %v, want %v", i, param, tt.parameters[i])
					}
				}
			}
		})
	}
}

func TestMotionCommandExecutor_ExecuteMoveAbsolute_Error(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Set drive controller to return error
	driveController.SetError(errors.New("drive error"))
	
	// Create a move absolute command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := executor.ExecuteMoveAbsolute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got %v", err)
	}
}

func TestMotionCommandExecutor_ExecuteMoveAbsolute_MissingParameter(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Create a move absolute command with missing position parameter
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := executor.ExecuteMoveAbsolute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if err.Error() != "missing position parameter: parameter position not found" {
		t.Errorf("Expected missing position parameter error, got %v", err)
	}
}

// Additional methods needed for I/O command testing
func (mdc *MockDriveController) GetDigitalOutput(output int) bool {
	return mdc.digitalOutputs[output]
}

func (mdc *MockDriveController) GetAnalogOutput(output int) float64 {
	return mdc.analogOutputs[output]
}
