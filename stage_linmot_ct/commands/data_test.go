package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)


func TestNewDataCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

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

func TestDataCommandExecutor_ExecuteStartOscilloscope(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdStartOscilloscope).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestDataCommandExecutor_ExecuteStopOscilloscope(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdStopOscilloscope).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestDataCommandExecutor_ExecuteSaveData(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSaveData).
		WithParameter("filename", "test_data.csv").
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestDataCommandExecutor_ValidateDataCommand(t *testing.T) {
	executor := NewDataCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid start oscilloscope command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdStartOscilloscope).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid stop oscilloscope command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdStopOscilloscope).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid save data command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSaveData).
				WithParameter("filename", "data.csv").
				Build(),
			wantErr: false,
		},
		{
			name: "Missing filename parameter for save data",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSaveData).
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
			err := executor.ValidateDataCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataCommandExecutor_GetDataCommandInfo(t *testing.T) {
	executor := NewDataCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		commandType types.CommandType
		wantDesc    string
		wantParams  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdStartOscilloscope,
			wantDesc:    "Start data acquisition (oscilloscope)",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdStopOscilloscope,
			wantDesc:    "Stop data acquisition (oscilloscope)",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdSaveData,
			wantDesc:    "Save acquired data to file",
			wantParams:  []string{"filename"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveAbsolute, // Unsupported
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			description, parameters, err := executor.GetDataCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if description != tt.wantDesc {
					t.Errorf("GetDataCommandInfo() description = %v, want %v", description, tt.wantDesc)
				}

				if len(parameters) != len(tt.wantParams) {
					t.Errorf("GetDataCommandInfo() parameters length = %v, want %v", len(parameters), len(tt.wantParams))
				}
			}
		})
	}
}

func TestDataCommandExecutor_ExecuteStartOscilloscope_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdStartOscilloscope).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestDataCommandExecutor_ExecuteSaveData_MissingParameter(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewDataCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSaveData).
		// Missing filename parameter
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "missing or invalid filename parameter: parameter filename not found" {
		t.Errorf("Expected missing filename parameter error, got: %v", err)
	}
}