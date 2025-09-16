package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// DataCommandExecutor handles data acquisition related commands
type DataCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (dce *DataCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdStartOscilloscope:
		return dce.ExecuteStartOscilloscope(ctx, command)
	case types.CmdStopOscilloscope:
		return dce.ExecuteStopOscilloscope(ctx, command)
	case types.CmdSaveData:
		return dce.ExecuteSaveData(ctx, command)
	default:
		return fmt.Errorf("unsupported data acquisition command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (dce *DataCommandExecutor) Validate(command *types.Command) error {
	return dce.ValidateDataCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (dce *DataCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return dce.GetDataCommandInfo(commandType)
}

// NewDataCommandExecutor creates a new data acquisition command executor
func NewDataCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *DataCommandExecutor {
	return &DataCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteStartOscilloscope executes a start oscilloscope command
func (dce *DataCommandExecutor) ExecuteStartOscilloscope(ctx context.Context, command *types.Command) error {
	// Start oscilloscope typically doesn't require parameters
	// Just start data acquisition on the drive
	return dce.driveController.StartOscilloscope(ctx)
}

// ExecuteStopOscilloscope executes a stop oscilloscope command
func (dce *DataCommandExecutor) ExecuteStopOscilloscope(ctx context.Context, command *types.Command) error {
	// Stop oscilloscope typically doesn't require parameters
	// Just stop data acquisition on the drive
	return dce.driveController.StopOscilloscope(ctx)
}

// ExecuteSaveData executes a save data command
func (dce *DataCommandExecutor) ExecuteSaveData(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract filename parameter
	filename, err := pe.ExtractString(command.Parameters, "filename")
	if err != nil {
		return fmt.Errorf("missing or invalid filename parameter: %w", err)
	}
	
	// Execute the command
	return dce.driveController.SaveData(ctx, filename)
}

// ValidateDataCommand validates a data acquisition command
func (dce *DataCommandExecutor) ValidateDataCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdStartOscilloscope, types.CmdStopOscilloscope:
		// These commands typically don't require parameters
		return nil
	case types.CmdSaveData:
		pe := types.NewParameterExtractor()
		
		// Check for required filename parameter
		_, err := pe.ExtractString(command.Parameters, "filename")
		if err != nil {
			return fmt.Errorf("missing or invalid filename parameter: %w", err)
		}
		
		return nil
	default:
		return fmt.Errorf("unsupported data acquisition command type: %s", command.Type)
	}
}

// GetDataCommandInfo returns command information for data acquisition commands
func (dce *DataCommandExecutor) GetDataCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdStartOscilloscope:
		return "Start data acquisition (oscilloscope)", []string{}, nil
	case types.CmdStopOscilloscope:
		return "Stop data acquisition (oscilloscope)", []string{}, nil
	case types.CmdSaveData:
		return "Save acquired data to file", []string{"filename"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported data acquisition command type: %s", commandType)
	}
}