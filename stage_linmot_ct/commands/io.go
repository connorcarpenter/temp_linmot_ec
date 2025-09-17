package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// IOCommandExecutor handles I/O-related commands
type IOCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (ioe *IOCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdSetDigitalOutput:
		return ioe.ExecuteSetDigitalOutput(ctx, command)
	case types.CmdClearDigitalOutput:
		return ioe.ExecuteClearDigitalOutput(ctx, command)
	case types.CmdSetAnalogOutput:
		return ioe.ExecuteSetAnalogOutput(ctx, command)
	case types.CmdWaitDigitalInput:
		return ioe.ExecuteWaitDigitalInput(ctx, command)
	case types.CmdWaitAnalogInput:
		return ioe.ExecuteWaitAnalogInput(ctx, command)
	default:
		return fmt.Errorf("unsupported I/O command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (ioe *IOCommandExecutor) Validate(command *types.Command) error {
	return ioe.ValidateIOCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (ioe *IOCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return ioe.GetIOCommandInfo(commandType)
}

// NewIOCommandExecutor creates a new I/O command executor
func NewIOCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *IOCommandExecutor {
	return &IOCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteSetDigitalOutput executes a set digital output command
func (ioe *IOCommandExecutor) ExecuteSetDigitalOutput(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract output number
	output, err := pe.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing or invalid output parameter: %w", err)
	}
	
	// Extract value
	value, err := pe.ExtractBool(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing or invalid value parameter: %w", err)
	}
	
	// Execute the command
	return ioe.driveController.SetDigitalOutput(ctx, output, value)
}

// ExecuteClearDigitalOutput executes a clear digital output command
func (ioe *IOCommandExecutor) ExecuteClearDigitalOutput(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract output number
	output, err := pe.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing or invalid output parameter: %w", err)
	}
	
	// Execute the command
	return ioe.driveController.ClearDigitalOutput(ctx, output)
}

// ExecuteSetAnalogOutput executes a set analog output command
func (ioe *IOCommandExecutor) ExecuteSetAnalogOutput(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract output number
	output, err := pe.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing or invalid output parameter: %w", err)
	}
	
	// Extract value
	value, err := pe.ExtractFloat(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing or invalid value parameter: %w", err)
	}
	
	// Execute the command
	return ioe.driveController.SetAnalogOutput(ctx, output, value)
}

// ExecuteWaitDigitalInput executes a wait digital input command
func (ioe *IOCommandExecutor) ExecuteWaitDigitalInput(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract input number
	input, err := pe.ExtractInt(command.Parameters, "input")
	if err != nil {
		return fmt.Errorf("missing or invalid input parameter: %w", err)
	}
	
	// Extract expected value
	value, err := pe.ExtractBool(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing or invalid value parameter: %w", err)
	}
	
	// Extract timeout
	timeout, err := pe.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing or invalid timeout parameter: %w", err)
	}
	
	// Execute the command
	return ioe.driveController.WaitDigitalInput(ctx, input, value, timeout.Duration())
}

// ExecuteWaitAnalogInput executes a wait analog input command
func (ioe *IOCommandExecutor) ExecuteWaitAnalogInput(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract input number
	input, err := pe.ExtractInt(command.Parameters, "input")
	if err != nil {
		return fmt.Errorf("missing or invalid input parameter: %w", err)
	}
	
	// Extract expected value
	value, err := pe.ExtractFloat(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing or invalid value parameter: %w", err)
	}
	
	// Extract tolerance
	tolerance, err := pe.ExtractFloat(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing or invalid tolerance parameter: %w", err)
	}
	
	// Extract timeout
	timeout, err := pe.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing or invalid timeout parameter: %w", err)
	}
	
	// Execute the command
	return ioe.driveController.WaitAnalogInput(ctx, input, value, tolerance, timeout.Duration())
}

// ValidateIOCommand validates an I/O command
func (ioe *IOCommandExecutor) ValidateIOCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdSetDigitalOutput:
		return ioe.validateSetDigitalOutput(command)
	case types.CmdClearDigitalOutput:
		return ioe.validateClearDigitalOutput(command)
	case types.CmdSetAnalogOutput:
		return ioe.validateSetAnalogOutput(command)
	case types.CmdWaitDigitalInput:
		return ioe.validateWaitDigitalInput(command)
	case types.CmdWaitAnalogInput:
		return ioe.validateWaitAnalogInput(command)
	default:
		return fmt.Errorf("unsupported I/O command type: %s", command.Type)
	}
}

// validateSetDigitalOutput validates a set digital output command
func (ioe *IOCommandExecutor) validateSetDigitalOutput(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "output"); err != nil {
		return fmt.Errorf("missing output parameter")
	}
	
	if _, err := pe.ExtractBool(command.Parameters, "value"); err != nil {
		return fmt.Errorf("missing value parameter")
	}
	
	return nil
}

// validateClearDigitalOutput validates a clear digital output command
func (ioe *IOCommandExecutor) validateClearDigitalOutput(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "output"); err != nil {
		return fmt.Errorf("missing output parameter")
	}
	
	return nil
}

// validateSetAnalogOutput validates a set analog output command
func (ioe *IOCommandExecutor) validateSetAnalogOutput(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "output"); err != nil {
		return fmt.Errorf("missing output parameter")
	}
	
	if _, err := pe.ExtractFloat(command.Parameters, "value"); err != nil {
		return fmt.Errorf("missing value parameter")
	}
	
	return nil
}

// validateWaitDigitalInput validates a wait digital input command
func (ioe *IOCommandExecutor) validateWaitDigitalInput(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "input"); err != nil {
		return fmt.Errorf("missing input parameter")
	}
	
	if _, err := pe.ExtractBool(command.Parameters, "value"); err != nil {
		return fmt.Errorf("missing value parameter")
	}
	
	if _, err := pe.ExtractTime(command.Parameters, "timeout"); err != nil {
		return fmt.Errorf("missing timeout parameter")
	}
	
	return nil
}

// validateWaitAnalogInput validates a wait analog input command
func (ioe *IOCommandExecutor) validateWaitAnalogInput(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "input"); err != nil {
		return fmt.Errorf("missing input parameter")
	}
	
	if _, err := pe.ExtractFloat(command.Parameters, "value"); err != nil {
		return fmt.Errorf("missing value parameter")
	}
	
	if _, err := pe.ExtractFloat(command.Parameters, "tolerance"); err != nil {
		return fmt.Errorf("missing tolerance parameter")
	}
	
	if _, err := pe.ExtractTime(command.Parameters, "timeout"); err != nil {
		return fmt.Errorf("missing timeout parameter")
	}
	
	return nil
}

// GetIOCommandInfo returns information about I/O commands
func (ioe *IOCommandExecutor) GetIOCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdSetDigitalOutput:
		return "SetDigitalOutput", []string{"output", "value"}, nil
	case types.CmdClearDigitalOutput:
		return "ClearDigitalOutput", []string{"output"}, nil
	case types.CmdSetAnalogOutput:
		return "SetAnalogOutput", []string{"output", "value"}, nil
	case types.CmdWaitDigitalInput:
		return "WaitDigitalInput", []string{"input", "value", "timeout"}, nil
	case types.CmdWaitAnalogInput:
		return "WaitAnalogInput", []string{"input", "value", "tolerance", "timeout"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported I/O command type: %s", commandType)
	}
}