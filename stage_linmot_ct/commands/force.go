package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// ForceCommandExecutor handles force control related commands
type ForceCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (fce *ForceCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdForceControlOn:
		return fce.ExecuteForceControlOn(ctx, command)
	case types.CmdForceControlOff:
		return fce.ExecuteForceControlOff(ctx, command)
	case types.CmdSetForce:
		return fce.ExecuteSetForce(ctx, command)
	default:
		return fmt.Errorf("unsupported force control command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (fce *ForceCommandExecutor) Validate(command *types.Command) error {
	return fce.ValidateForceCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (fce *ForceCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return fce.GetForceCommandInfo(commandType)
}

// NewForceCommandExecutor creates a new force control command executor
func NewForceCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *ForceCommandExecutor {
	return &ForceCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteForceControlOn executes a force control on command
func (fce *ForceCommandExecutor) ExecuteForceControlOn(ctx context.Context, command *types.Command) error {
	// Force control on typically doesn't require parameters
	// Just enable force control mode on the drive
	return fce.driveController.ForceControlOn(ctx)
}

// ExecuteForceControlOff executes a force control off command
func (fce *ForceCommandExecutor) ExecuteForceControlOff(ctx context.Context, command *types.Command) error {
	// Force control off typically doesn't require parameters
	// Just disable force control mode on the drive
	return fce.driveController.ForceControlOff(ctx)
}

// ExecuteSetForce executes a set force command
func (fce *ForceCommandExecutor) ExecuteSetForce(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract force value
	forceValue, err := pe.ExtractForce(command.Parameters, "force")
	if err != nil {
		return fmt.Errorf("missing or invalid force parameter: %w", err)
	}
	
	// Convert force to counts if needed
	forceInCounts := fce.unitConverter.ConvertForceValue(forceValue, types.ForceUnitCounts)
	
	// Execute the command
	return fce.driveController.SetForce(ctx, forceInCounts.Value)
}

// ValidateForceCommand validates a force control command
func (fce *ForceCommandExecutor) ValidateForceCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdForceControlOn, types.CmdForceControlOff:
		// These commands typically don't require parameters
		return nil
	case types.CmdSetForce:
		pe := types.NewParameterExtractor()
		
		// Check for required force parameter
		_, err := pe.ExtractForce(command.Parameters, "force")
		if err != nil {
			return fmt.Errorf("missing or invalid force parameter: %w", err)
		}
		
		return nil
	default:
		return fmt.Errorf("unsupported force control command type: %s", command.Type)
	}
}

// GetForceCommandInfo returns command information for force control commands
func (fce *ForceCommandExecutor) GetForceCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdForceControlOn:
		return "Enable force control mode", []string{}, nil
	case types.CmdForceControlOff:
		return "Disable force control mode", []string{}, nil
	case types.CmdSetForce:
		return "Set force setpoint", []string{"force"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported force control command type: %s", commandType)
	}
}