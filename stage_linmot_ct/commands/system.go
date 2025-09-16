package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// SystemCommandExecutor handles system-related commands
type SystemCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (sce *SystemCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdHome:
		return sce.ExecuteHome(ctx, command)
	case types.CmdReset:
		return sce.ExecuteReset(ctx, command)
	case types.CmdSaveConfiguration:
		return sce.ExecuteSaveConfiguration(ctx, command)
	case types.CmdLoadConfiguration:
		return sce.ExecuteLoadConfiguration(ctx, command)
	default:
		return fmt.Errorf("unsupported system command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (sce *SystemCommandExecutor) Validate(command *types.Command) error {
	return sce.ValidateSystemCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (sce *SystemCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return sce.GetSystemCommandInfo(commandType)
}

// NewSystemCommandExecutor creates a new system command executor
func NewSystemCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *SystemCommandExecutor {
	return &SystemCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteHome executes a home command
func (sce *SystemCommandExecutor) ExecuteHome(ctx context.Context, command *types.Command) error {
	// Execute the home command
	return sce.driveController.Home(ctx)
}

// ExecuteReset executes a reset command
func (sce *SystemCommandExecutor) ExecuteReset(ctx context.Context, command *types.Command) error {
	// Execute the reset command
	return sce.driveController.Reset(ctx)
}

// ExecuteSaveConfiguration executes a save configuration command
func (sce *SystemCommandExecutor) ExecuteSaveConfiguration(ctx context.Context, command *types.Command) error {
	// Execute the save configuration command
	return sce.driveController.SaveConfiguration(ctx)
}

// ExecuteLoadConfiguration executes a load configuration command
func (sce *SystemCommandExecutor) ExecuteLoadConfiguration(ctx context.Context, command *types.Command) error {
	// Execute the load configuration command
	return sce.driveController.LoadConfiguration(ctx)
}

// ValidateSystemCommand validates a system command
func (sce *SystemCommandExecutor) ValidateSystemCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdHome:
		return sce.validateHome(command)
	case types.CmdReset:
		return sce.validateReset(command)
	case types.CmdSaveConfiguration:
		return sce.validateSaveConfiguration(command)
	case types.CmdLoadConfiguration:
		return sce.validateLoadConfiguration(command)
	default:
		return fmt.Errorf("unsupported system command type: %s", command.Type)
	}
}

// validateHome validates a home command
func (sce *SystemCommandExecutor) validateHome(command *types.Command) error {
	// Home commands typically don't require parameters
	return nil
}

// validateReset validates a reset command
func (sce *SystemCommandExecutor) validateReset(command *types.Command) error {
	// Reset commands typically don't require parameters
	return nil
}

// validateSaveConfiguration validates a save configuration command
func (sce *SystemCommandExecutor) validateSaveConfiguration(command *types.Command) error {
	// Save configuration commands typically don't require parameters
	return nil
}

// validateLoadConfiguration validates a load configuration command
func (sce *SystemCommandExecutor) validateLoadConfiguration(command *types.Command) error {
	// Load configuration commands typically don't require parameters
	return nil
}

// GetSystemCommandInfo returns information about system commands
func (sce *SystemCommandExecutor) GetSystemCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdHome:
		return "Home", []string{}, nil
	case types.CmdReset:
		return "Reset", []string{}, nil
	case types.CmdSaveConfiguration:
		return "SaveConfiguration", []string{}, nil
	case types.CmdLoadConfiguration:
		return "LoadConfiguration", []string{}, nil
	default:
		return "", nil, fmt.Errorf("unsupported system command type: %s", commandType)
	}
}