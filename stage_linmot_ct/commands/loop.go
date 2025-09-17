package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// LoopCommandExecutor handles loop and jump commands
type LoopCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
}

// Execute implements CommandExecutor interface
func (lce *LoopCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	switch command.Type {
	case types.CmdLoopStart:
		return lce.ExecuteLoopStart(ctx, command)
	case types.CmdLoopEnd:
		return lce.ExecuteLoopEnd(ctx, command)
	case types.CmdLoopBreak:
		return lce.ExecuteLoopBreak(ctx, command)
	case types.CmdJump:
		return lce.ExecuteJump(ctx, command)
	case types.CmdJumpIfTrue:
		return lce.ExecuteJumpIfTrue(ctx, command)
	case types.CmdJumpIfFalse:
		return lce.ExecuteJumpIfFalse(ctx, command)
	default:
		return fmt.Errorf("unsupported loop/jump command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (lce *LoopCommandExecutor) Validate(command *types.Command) error {
	return lce.ValidateLoopCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (lce *LoopCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return lce.GetLoopCommandInfo(commandType)
}

// NewLoopCommandExecutor creates a new loop command executor
func NewLoopCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter) *LoopCommandExecutor {
	return &LoopCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
	}
}

// ExecuteLoopStart executes a loop start command
func (lce *LoopCommandExecutor) ExecuteLoopStart(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract loop count
	count, err := pe.ExtractInt(command.Parameters, "count")
	if err != nil {
		return fmt.Errorf("missing or invalid count parameter: %w", err)
	}
	
	// Validate count
	if count <= 0 {
		return fmt.Errorf("loop count must be positive, got %d", count)
	}
	
	// TODO: Implement loop start logic
	// This would typically involve setting up loop state in the execution engine
	// For now, we'll just validate the parameters
	
	return nil
}

// ExecuteLoopEnd executes a loop end command
func (lce *LoopCommandExecutor) ExecuteLoopEnd(ctx context.Context, command *types.Command) error {
	// TODO: Implement loop end logic
	// This would typically involve checking if we should continue the loop
	// or exit the loop based on the current loop state
	
	return nil
}

// ExecuteLoopBreak executes a loop break command
func (lce *LoopCommandExecutor) ExecuteLoopBreak(ctx context.Context, command *types.Command) error {
	// TODO: Implement loop break logic
	// This would typically involve breaking out of the current loop
	// and jumping to the end of the loop
	
	return nil
}

// ExecuteJump executes a jump command
func (lce *LoopCommandExecutor) ExecuteJump(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract target command ID
	targetID, err := pe.ExtractInt(command.Parameters, "target_id")
	if err != nil {
		return fmt.Errorf("missing or invalid target_id parameter: %w", err)
	}
	
	// Validate target ID
	if targetID <= 0 {
		return fmt.Errorf("target command ID must be positive, got %d", targetID)
	}
	
	// TODO: Implement jump logic
	// This would typically involve setting the current command index
	// to the target command ID in the execution engine
	
	return nil
}

// ExecuteJumpIfTrue executes a conditional jump if true command
func (lce *LoopCommandExecutor) ExecuteJumpIfTrue(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract target command ID
	targetID, err := pe.ExtractInt(command.Parameters, "target_id")
	if err != nil {
		return fmt.Errorf("missing or invalid target_id parameter: %w", err)
	}
	
	// Extract condition
	_, err = pe.ExtractBool(command.Parameters, "condition")
	if err != nil {
		return fmt.Errorf("missing or invalid condition parameter: %w", err)
	}
	
	// Validate target ID
	if targetID <= 0 {
		return fmt.Errorf("target command ID must be positive, got %d", targetID)
	}
	
	// TODO: Implement conditional jump logic
	// This would typically involve evaluating the condition
	// and jumping to the target if the condition is true
	
	return nil
}

// ExecuteJumpIfFalse executes a conditional jump if false command
func (lce *LoopCommandExecutor) ExecuteJumpIfFalse(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Extract target command ID
	targetID, err := pe.ExtractInt(command.Parameters, "target_id")
	if err != nil {
		return fmt.Errorf("missing or invalid target_id parameter: %w", err)
	}
	
	// Extract condition
	_, err = pe.ExtractBool(command.Parameters, "condition")
	if err != nil {
		return fmt.Errorf("missing or invalid condition parameter: %w", err)
	}
	
	// Validate target ID
	if targetID <= 0 {
		return fmt.Errorf("target command ID must be positive, got %d", targetID)
	}
	
	// TODO: Implement conditional jump logic
	// This would typically involve evaluating the condition
	// and jumping to the target if the condition is false
	
	return nil
}

// ValidateLoopCommand validates a loop/jump command
func (lce *LoopCommandExecutor) ValidateLoopCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdLoopStart:
		return lce.validateLoopStart(command)
	case types.CmdLoopEnd:
		return lce.validateLoopEnd(command)
	case types.CmdLoopBreak:
		return lce.validateLoopBreak(command)
	case types.CmdJump:
		return lce.validateJump(command)
	case types.CmdJumpIfTrue:
		return lce.validateJumpIfTrue(command)
	case types.CmdJumpIfFalse:
		return lce.validateJumpIfFalse(command)
	default:
		return fmt.Errorf("unsupported loop/jump command type: %s", command.Type)
	}
}

// validateLoopStart validates a loop start command
func (lce *LoopCommandExecutor) validateLoopStart(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "count"); err != nil {
		return fmt.Errorf("missing count parameter")
	}
	
	return nil
}

// validateLoopEnd validates a loop end command
func (lce *LoopCommandExecutor) validateLoopEnd(command *types.Command) error {
	// Loop end commands typically don't require parameters
	return nil
}

// validateLoopBreak validates a loop break command
func (lce *LoopCommandExecutor) validateLoopBreak(command *types.Command) error {
	// Loop break commands typically don't require parameters
	return nil
}

// validateJump validates a jump command
func (lce *LoopCommandExecutor) validateJump(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "target_id"); err != nil {
		return fmt.Errorf("missing target_id parameter")
	}
	
	return nil
}

// validateJumpIfTrue validates a jump if true command
func (lce *LoopCommandExecutor) validateJumpIfTrue(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "target_id"); err != nil {
		return fmt.Errorf("missing target_id parameter")
	}
	
	if _, err := pe.ExtractBool(command.Parameters, "condition"); err != nil {
		return fmt.Errorf("missing condition parameter")
	}
	
	return nil
}

// validateJumpIfFalse validates a jump if false command
func (lce *LoopCommandExecutor) validateJumpIfFalse(command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check required parameters
	if _, err := pe.ExtractInt(command.Parameters, "target_id"); err != nil {
		return fmt.Errorf("missing target_id parameter")
	}
	
	if _, err := pe.ExtractBool(command.Parameters, "condition"); err != nil {
		return fmt.Errorf("missing condition parameter")
	}
	
	return nil
}

// GetLoopCommandInfo returns information about loop/jump commands
func (lce *LoopCommandExecutor) GetLoopCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdLoopStart:
		return "LoopStart", []string{"count"}, nil
	case types.CmdLoopEnd:
		return "LoopEnd", []string{}, nil
	case types.CmdLoopBreak:
		return "LoopBreak", []string{}, nil
	case types.CmdJump:
		return "Jump", []string{"target_id"}, nil
	case types.CmdJumpIfTrue:
		return "JumpIfTrue", []string{"target_id", "condition"}, nil
	case types.CmdJumpIfFalse:
		return "JumpIfFalse", []string{"target_id", "condition"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported loop/jump command type: %s", commandType)
	}
}