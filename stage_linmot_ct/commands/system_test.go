package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewSystemCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

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

func TestSystemCommandExecutor_ExecuteHome(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdHome).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteReset(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdReset).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteSaveConfiguration(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSaveConfiguration).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteLoadConfiguration(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoadConfiguration).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestSystemCommandExecutor_ValidateSystemCommand(t *testing.T) {
	executor := NewSystemCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid home command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdHome).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid reset command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdReset).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid save configuration command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSaveConfiguration).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid load configuration command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdLoadConfiguration).
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
			err := executor.ValidateSystemCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSystemCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemCommandExecutor_GetSystemCommandInfo(t *testing.T) {
	executor := NewSystemCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		commandType types.CommandType
		wantName    string
		wantParams  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdHome,
			wantName:    "Home",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdReset,
			wantName:    "Reset",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdSaveConfiguration,
			wantName:    "SaveConfiguration",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdLoadConfiguration,
			wantName:    "LoadConfiguration",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveAbsolute, // Unsupported
			wantName:    "",
			wantParams:  nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			name, params, err := executor.GetSystemCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if name != tt.wantName {
				t.Errorf("GetSystemCommandInfo() name = %v, want %v", name, tt.wantName)
			}
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetSystemCommandInfo() params length = %v, want %v", len(params), len(tt.wantParams))
			}
		})
	}
}

func TestSystemCommandExecutor_ExecuteHome_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdHome).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteReset_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdReset).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteSaveConfiguration_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSaveConfiguration).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestSystemCommandExecutor_ExecuteLoadConfiguration_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewSystemCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoadConfiguration).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}