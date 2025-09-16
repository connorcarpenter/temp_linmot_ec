package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)


func TestNewForceCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	if executor == nil {
		t.Fatal("Expected non-nil executor")
	}

	if executor.driveController != driveController {
		t.Error("DriveController not set correctly")
	}

	if executor.unitConverter != unitConverter {
		t.Error("UnitConverter not set correctly")
	}
}

func TestForceCommandExecutor_ExecuteForceControlOn(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdForceControlOn).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestForceCommandExecutor_ExecuteForceControlOff(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdForceControlOff).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestForceCommandExecutor_ExecuteSetForce(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetForce).
		WithParameter("force", types.NewForceValue(10.0, types.ForceUnitN)).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify the force was set
	actualForce, _ := driveController.GetForce(context.Background())
	expectedForce := 10.0 * 100.0 // Convert N to counts (100 counts/N scaling)
	if actualForce != expectedForce {
		t.Errorf("Expected force %f, got %f", expectedForce, actualForce)
	}
}

func TestForceCommandExecutor_ValidateForceCommand(t *testing.T) {
	executor := NewForceCommandExecutor(NewMockDriveController(), types.NewUnitConverter(), nil)

	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid force control on command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdForceControlOn).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid force control off command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdForceControlOff).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid set force command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetForce).
				WithParameter("force", types.NewForceValue(5.0, types.ForceUnitN)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing force parameter for set force",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetForce).
				Build(),
			wantErr: true,
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
			err := executor.ValidateForceCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateForceCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestForceCommandExecutor_GetForceCommandInfo(t *testing.T) {
	executor := NewForceCommandExecutor(NewMockDriveController(), types.NewUnitConverter(), nil)

	tests := []struct {
		commandType types.CommandType
		wantDesc    string
		wantParams  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdForceControlOn,
			wantDesc:    "Enable force control mode",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdForceControlOff,
			wantDesc:    "Disable force control mode",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdSetForce,
			wantDesc:    "Set force setpoint",
			wantParams:  []string{"force"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveAbsolute, // Unsupported
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			description, parameters, err := executor.GetForceCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForceCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if description != tt.wantDesc {
					t.Errorf("GetForceCommandInfo() description = %v, want %v", description, tt.wantDesc)
				}

				if len(parameters) != len(tt.wantParams) {
					t.Errorf("GetForceCommandInfo() parameters length = %v, want %v", len(parameters), len(tt.wantParams))
				}
			}
		})
	}
}

func TestForceCommandExecutor_ExecuteForceControlOn_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdForceControlOn).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestForceCommandExecutor_ExecuteSetForce_MissingParameter(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewForceCommandExecutor(driveController, unitConverter, nil)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetForce).
		// Missing force parameter
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "missing or invalid force parameter: parameter force not found" {
		t.Errorf("Expected missing force parameter error, got: %v", err)
	}
}