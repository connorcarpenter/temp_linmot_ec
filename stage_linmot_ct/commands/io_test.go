package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)


func TestNewIOCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

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

func TestIOCommandExecutor_ExecuteSetDigitalOutput(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetDigitalOutput).
		WithParameter("output", 1).
		WithParameter("value", true).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify the command was executed
	if !driveController.GetDigitalOutput(1) {
		t.Error("Expected digital output 1 to be set to true")
	}
}

func TestIOCommandExecutor_ExecuteClearDigitalOutput(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	// Set the output first
	driveController.SetDigitalOutput(context.Background(), 1, true)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdClearDigitalOutput).
		WithParameter("output", 1).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify the command was executed
	if driveController.GetDigitalOutput(1) {
		t.Error("Expected digital output 1 to be cleared")
	}
}

func TestIOCommandExecutor_ExecuteSetAnalogOutput(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetAnalogOutput).
		WithParameter("output", 1).
		WithParameter("value", 3.14).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify the command was executed
	expected := 3.14
	if driveController.GetAnalogOutput(1) != expected {
		t.Errorf("Expected analog output 1 to be %f, got %f", expected, driveController.GetAnalogOutput(1))
	}
}

func TestIOCommandExecutor_ExecuteWaitDigitalInput(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	// Set the input to the expected value
	driveController.SetDigitalInput(1, true)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWaitDigitalInput).
		WithParameter("input", 1).
		WithParameter("value", true).
		WithParameter("timeout", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestIOCommandExecutor_ExecuteWaitAnalogInput(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	// Set the input to the expected value
	driveController.SetAnalogInput(1, 2.5)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWaitAnalogInput).
		WithParameter("input", 1).
		WithParameter("value", 2.5).
		WithParameter("tolerance", 0.1).
		WithParameter("timeout", types.NewTimeValue(100.0, types.TimeUnitMS)).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestIOCommandExecutor_ValidateIOCommand(t *testing.T) {
	executor := NewIOCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid set digital output command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetDigitalOutput).
				WithParameter("output", 1).
				WithParameter("value", true).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing output parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetDigitalOutput).
				WithParameter("value", true).
				Build(),
			wantErr: true,
		},
		{
			name: "Missing value parameter",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetDigitalOutput).
				WithParameter("output", 1).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid clear digital output command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdClearDigitalOutput).
				WithParameter("output", 1).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid set analog output command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetAnalogOutput).
				WithParameter("output", 1).
				WithParameter("value", 3.14).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid wait digital input command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitDigitalInput).
				WithParameter("input", 1).
				WithParameter("value", true).
				WithParameter("timeout", types.NewTimeValue(100.0, types.TimeUnitMS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing timeout parameter for wait digital input",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitDigitalInput).
				WithParameter("input", 1).
				WithParameter("value", true).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid wait analog input command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitAnalogInput).
				WithParameter("input", 1).
				WithParameter("value", 2.5).
				WithParameter("tolerance", 0.1).
				WithParameter("timeout", types.NewTimeValue(100.0, types.TimeUnitMS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing tolerance parameter for wait analog input",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWaitAnalogInput).
				WithParameter("input", 1).
				WithParameter("value", 2.5).
				WithParameter("timeout", types.NewTimeValue(100.0, types.TimeUnitMS)).
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
			err := executor.ValidateIOCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIOCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIOCommandExecutor_GetIOCommandInfo(t *testing.T) {
	executor := NewIOCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		commandType types.CommandType
		wantName    string
		wantParams  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdSetDigitalOutput,
			wantName:    "SetDigitalOutput",
			wantParams:  []string{"output", "value"},
			wantErr:     false,
		},
		{
			commandType: types.CmdClearDigitalOutput,
			wantName:    "ClearDigitalOutput",
			wantParams:  []string{"output"},
			wantErr:     false,
		},
		{
			commandType: types.CmdSetAnalogOutput,
			wantName:    "SetAnalogOutput",
			wantParams:  []string{"output", "value"},
			wantErr:     false,
		},
		{
			commandType: types.CmdWaitDigitalInput,
			wantName:    "WaitDigitalInput",
			wantParams:  []string{"input", "value", "timeout"},
			wantErr:     false,
		},
		{
			commandType: types.CmdWaitAnalogInput,
			wantName:    "WaitAnalogInput",
			wantParams:  []string{"input", "value", "tolerance", "timeout"},
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
			name, params, err := executor.GetIOCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIOCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if name != tt.wantName {
				t.Errorf("GetIOCommandInfo() name = %v, want %v", name, tt.wantName)
			}
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetIOCommandInfo() params length = %v, want %v", len(params), len(tt.wantParams))
			}
		})
	}
}

func TestIOCommandExecutor_ExecuteSetDigitalOutput_Error(t *testing.T) {
	driveController := NewMockDriveController()
	driveController.SetError(errors.New("drive error"))
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetDigitalOutput).
		WithParameter("output", 1).
		WithParameter("value", true).
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "drive error" {
		t.Errorf("Expected 'drive error', got: %v", err)
	}
}

func TestIOCommandExecutor_ExecuteSetDigitalOutput_MissingParameter(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewIOCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdSetDigitalOutput).
		WithParameter("output", 1).
		// Missing value parameter
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "missing or invalid value parameter: parameter value not found" {
		t.Errorf("Expected missing value parameter error, got: %v", err)
	}
}