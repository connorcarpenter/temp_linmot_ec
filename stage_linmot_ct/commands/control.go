package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// ControlCommandExecutor handles control-related commands
type ControlCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (cce *ControlCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdWait:
		return cce.ExecuteWait(ctx, command)
	case types.CmdWaitPosition:
		return cce.ExecuteWaitPosition(ctx, command)
	case types.CmdWaitVelocity:
		return cce.ExecuteWaitVelocity(ctx, command)
	case types.CmdWaitForce:
		return cce.ExecuteWaitForce(ctx, command)
	default:
		return fmt.Errorf("unsupported control command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (cce *ControlCommandExecutor) Validate(command *types.Command) error {
	return cce.ValidateControlCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (cce *ControlCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return cce.GetControlCommandInfo(commandType)
}

// NewControlCommandExecutor creates a new control command executor
func NewControlCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *ControlCommandExecutor {
	return &ControlCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteWait executes a wait command
func (cce *ControlCommandExecutor) ExecuteWait(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	duration, err := extractor.ExtractTime(command.Parameters, "duration")
	if err != nil {
		return fmt.Errorf("missing duration parameter: %w", err)
	}
	
	// Execute the command
	return cce.driveController.Wait(ctx, duration.Duration())
}

// ExecuteWaitPosition executes a wait position command
func (cce *ControlCommandExecutor) ExecuteWaitPosition(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	position, err := extractor.ExtractPosition(command.Parameters, "position")
	if err != nil {
		return fmt.Errorf("missing position parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractPosition(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	posValue := cce.unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
	tolValue := cce.unitConverter.ConvertPositionValue(tolerance, types.PositionUnitCounts)
	
	// Execute the command
	return cce.driveController.WaitPosition(ctx, posValue.Value, tolValue.Value, timeout.Duration())
}

// ExecuteWaitVelocity executes a wait velocity command
func (cce *ControlCommandExecutor) ExecuteWaitVelocity(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractVelocity(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	velValue := cce.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	tolValue := cce.unitConverter.ConvertVelocityValue(tolerance, types.VelocityUnitCountsS)
	
	// Execute the command
	return cce.driveController.WaitVelocity(ctx, velValue.Value, tolValue.Value, timeout.Duration())
}

// ExecuteWaitForce executes a wait force command
func (cce *ControlCommandExecutor) ExecuteWaitForce(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	force, err := extractor.ExtractForce(command.Parameters, "force")
	if err != nil {
		return fmt.Errorf("missing force parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractForce(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	forceValue := cce.unitConverter.ConvertForceValue(force, types.ForceUnitCounts)
	tolValue := cce.unitConverter.ConvertForceValue(tolerance, types.ForceUnitCounts)
	
	// Execute the command
	return cce.driveController.WaitForce(ctx, forceValue.Value, tolValue.Value, timeout.Duration())
}

// ValidateControlCommand validates control command parameters
func (cce *ControlCommandExecutor) ValidateControlCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdWait:
		return cce.validateWaitParameters(command)
	case types.CmdWaitPosition:
		return cce.validateWaitPositionParameters(command)
	case types.CmdWaitVelocity:
		return cce.validateWaitVelocityParameters(command)
	case types.CmdWaitForce:
		return cce.validateWaitForceParameters(command)
	default:
		return fmt.Errorf("unsupported control command type: %s", command.Type)
	}
}

// validateWaitParameters validates parameters for wait command
func (cce *ControlCommandExecutor) validateWaitParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	if _, exists := command.Parameters["duration"]; !exists {
		return fmt.Errorf("missing required parameter: duration")
	}
	
	// Validate duration parameter
	_, err := extractor.ExtractTime(command.Parameters, "duration")
	if err != nil {
		return fmt.Errorf("invalid duration parameter: %w", err)
	}
	
	return nil
}

// validateWaitPositionParameters validates parameters for wait position command
func (cce *ControlCommandExecutor) validateWaitPositionParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	requiredParams := []string{"position", "tolerance", "timeout"}
	for _, param := range requiredParams {
		if _, exists := command.Parameters[param]; !exists {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	
	// Validate position parameter
	_, err := extractor.ExtractPosition(command.Parameters, "position")
	if err != nil {
		return fmt.Errorf("invalid position parameter: %w", err)
	}
	
	// Validate tolerance parameter
	_, err = extractor.ExtractPosition(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("invalid tolerance parameter: %w", err)
	}
	
	// Validate timeout parameter
	_, err = extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("invalid timeout parameter: %w", err)
	}
	
	return nil
}

// validateWaitVelocityParameters validates parameters for wait velocity command
func (cce *ControlCommandExecutor) validateWaitVelocityParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	requiredParams := []string{"velocity", "tolerance", "timeout"}
	for _, param := range requiredParams {
		if _, exists := command.Parameters[param]; !exists {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	
	// Validate velocity parameter
	_, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("invalid velocity parameter: %w", err)
	}
	
	// Validate tolerance parameter
	_, err = extractor.ExtractVelocity(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("invalid tolerance parameter: %w", err)
	}
	
	// Validate timeout parameter
	_, err = extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("invalid timeout parameter: %w", err)
	}
	
	return nil
}

// validateWaitForceParameters validates parameters for wait force command
func (cce *ControlCommandExecutor) validateWaitForceParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	requiredParams := []string{"force", "tolerance", "timeout"}
	for _, param := range requiredParams {
		if _, exists := command.Parameters[param]; !exists {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	
	// Validate force parameter
	_, err := extractor.ExtractForce(command.Parameters, "force")
	if err != nil {
		return fmt.Errorf("invalid force parameter: %w", err)
	}
	
	// Validate tolerance parameter
	_, err = extractor.ExtractForce(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("invalid tolerance parameter: %w", err)
	}
	
	// Validate timeout parameter
	_, err = extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("invalid timeout parameter: %w", err)
	}
	
	return nil
}

// GetControlCommandInfo returns information about control commands
func (cce *ControlCommandExecutor) GetControlCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdWait:
		return "Wait for specified duration", []string{"duration"}, nil
	case types.CmdWaitPosition:
		return "Wait for position with tolerance", []string{"position", "tolerance", "timeout"}, nil
	case types.CmdWaitVelocity:
		return "Wait for velocity with tolerance", []string{"velocity", "tolerance", "timeout"}, nil
	case types.CmdWaitForce:
		return "Wait for force with tolerance", []string{"force", "tolerance", "timeout"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported control command type: %s", commandType)
	}
}

// WaitForCondition waits for a condition to be met with timeout
func (cce *ControlCommandExecutor) WaitForCondition(ctx context.Context, condition func() (bool, error), timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	ticker := time.NewTicker(10 * time.Millisecond) // Check every 10ms
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			met, err := condition()
			if err != nil {
				return err
			}
			if met {
				return nil
			}
		}
	}
}