package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewControlCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	if executor == nil {
		t.Fatal("NewControlCommandExecutor returned nil")
	}
	
	if executor.driveController != driveController {
		t.Error("DriveController not set correctly")
	}
	
	if executor.unitConverter != unitConverter {
		t.Error("UnitConverter not set correctly")
	}
}

func TestControlCommandExecutor_ExecuteWait(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Create a wait command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()
	
	// Execute the command
	start := time.Now()
	err := executor.ExecuteWait(context.Background(), command)
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("ExecuteWait failed: %v", err)
	}
	
	// Verify the duration is approximately correct (allow some tolerance)
	expectedDuration := 100 * time.Millisecond
	if duration < expectedDuration-10*time.Millisecond || duration > expectedDuration+10*time.Millisecond {
		t.Errorf("Expected duration around %v, got %v", expectedDuration, duration)
	}
}

func TestControlCommandExecutor_ExecuteWaitPosition(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Create a wait position command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWaitPosition).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("tolerance", types.NewPositionValue(0.1, types.PositionUnitMM)).
		WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitS)).
		Build()
	
	// Execute the command
	err := executor.ExecuteWaitPosition(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteWaitPosition failed: %v", err)
	}
	
	// The mock drive controller should have been called
	// We can't easily verify the exact values without exposing them,
	// but we can verify no error occurred
}

func TestControlCommandExecutor_ExecuteWaitVelocity(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Create a wait velocity command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWaitVelocity).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("tolerance", types.NewVelocityValue(0.1, types.VelocityUnitMMS)).
		WithParameter("timeout", types.NewTimeValue(3.0, types.TimeUnitS)).
		Build()
	
	// Execute the command
	err := executor.ExecuteWaitVelocity(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteWaitVelocity failed: %v", err)
	}
}

func TestControlCommandExecutor_ExecuteWaitForce(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Create a wait force command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWaitForce).
		WithParameter("force", types.NewForceValue(10.0, types.ForceUnitN)).
		WithParameter("tolerance", types.NewForceValue(0.5, types.ForceUnitN)).
		WithParameter("timeout", types.NewTimeValue(2.0, types.TimeUnitS)).
		Build()
	
	// Execute the command
	err := executor.ExecuteWaitForce(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteWaitForce failed: %v", err)
	}
}

func TestControlCommandExecutor_ValidateControlCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid wait command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWait).
				WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing duration parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWait).
				Build(),
			wantErr: true,
		},
		{
			name: "Invalid duration parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWait).
				WithParameter("duration", "invalid").
				Build(),
			wantErr: true,
		},
		{
			name: "Valid wait position command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitPosition).
				WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
				WithParameter("tolerance", types.NewPositionValue(0.1, types.PositionUnitMM)).
				WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing position parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitPosition).
				WithParameter("tolerance", types.NewPositionValue(0.1, types.PositionUnitMM)).
				WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitS)).
				Build(),
			wantErr: true,
		},
		{
			name: "Missing tolerance parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitPosition).
				WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
				WithParameter("timeout", types.NewTimeValue(5.0, types.TimeUnitS)).
				Build(),
			wantErr: true,
		},
		{
			name: "Missing timeout parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitPosition).
				WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
				WithParameter("tolerance", types.NewPositionValue(0.1, types.PositionUnitMM)).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid wait velocity command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitVelocity).
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("tolerance", types.NewVelocityValue(0.1, types.VelocityUnitMMS)).
				WithParameter("timeout", types.NewTimeValue(3.0, types.TimeUnitS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid wait force command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitForce).
				WithParameter("force", types.NewForceValue(10.0, types.ForceUnitN)).
				WithParameter("tolerance", types.NewForceValue(0.5, types.ForceUnitN)).
				WithParameter("timeout", types.NewTimeValue(2.0, types.TimeUnitS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Unsupported command type",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				Build(),
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.ValidateControlCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateControlCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestControlCommandExecutor_GetControlCommandInfo(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	tests := []struct {
		commandType types.CommandType
		description string
		parameters  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdWait,
			description: "Wait for specified duration",
			parameters:  []string{"duration"},
			wantErr:     false,
		},
		{
			commandType: types.CmdWaitPosition,
			description: "Wait for position with tolerance",
			parameters:  []string{"position", "tolerance", "timeout"},
			wantErr:     false,
		},
		{
			commandType: types.CmdWaitVelocity,
			description: "Wait for velocity with tolerance",
			parameters:  []string{"velocity", "tolerance", "timeout"},
			wantErr:     false,
		},
		{
			commandType: types.CmdWaitForce,
			description: "Wait for force with tolerance",
			parameters:  []string{"force", "tolerance", "timeout"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveAbsolute, // Invalid for control executor
			description: "",
			parameters:  nil,
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			description, parameters, err := executor.GetControlCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetControlCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if description != tt.description {
					t.Errorf("GetControlCommandInfo() description = %v, want %v", description, tt.description)
				}
				
				if len(parameters) != len(tt.parameters) {
					t.Errorf("GetControlCommandInfo() parameters length = %v, want %v", len(parameters), len(tt.parameters))
				}
				
				for i, param := range parameters {
					if param != tt.parameters[i] {
						t.Errorf("GetControlCommandInfo() parameters[%d] = %v, want %v", i, param, tt.parameters[i])
					}
				}
			}
		})
	}
}

func TestControlCommandExecutor_ExecuteWait_Error(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Set drive controller to return error
	driveController.SetError(errors.New("drive error"))
	
	// Create a wait command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()
	
	// Execute the command
	err := executor.ExecuteWait(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got %v", err)
	}
}

func TestControlCommandExecutor_ExecuteWait_MissingParameter(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Create a wait command with missing duration parameter
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		Build()
	
	// Execute the command
	err := executor.ExecuteWait(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if err.Error() != "missing duration parameter: parameter duration not found" {
		t.Errorf("Expected missing duration parameter error, got %v", err)
	}
}

func TestControlCommandExecutor_WaitForCondition(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Test successful condition
	conditionMet := false
	condition := func() (bool, error) {
		conditionMet = true
		return true, nil
	}
	
	err := executor.WaitForCondition(context.Background(), condition, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("WaitForCondition failed: %v", err)
	}
	
	if !conditionMet {
		t.Error("Condition was not called")
	}
}

func TestControlCommandExecutor_WaitForCondition_Timeout(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Test timeout condition
	condition := func() (bool, error) {
		return false, nil // Never met
	}
	
	err := executor.WaitForCondition(context.Background(), condition, 50*time.Millisecond)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
	
	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestControlCommandExecutor_WaitForCondition_Error(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewControlCommandExecutor(driveController, unitConverter)
	
	// Test error condition
	expectedErr := errors.New("condition error")
	condition := func() (bool, error) {
		return false, expectedErr
	}
	
	err := executor.WaitForCondition(context.Background(), condition, 100*time.Millisecond)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if err != expectedErr {
		t.Errorf("Expected %v, got %v", expectedErr, err)
	}
}