package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/safety"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// MotionCommandExecutor handles motion-related commands
type MotionCommandExecutor struct {
	driveController types.DriveController
	unitConverter   *types.UnitConverter
	safetyGuard     *safety.SafetyGuard
}

// Execute implements CommandExecutor interface
func (mce *MotionCommandExecutor) Execute(ctx context.Context, command *types.Command) error {
	// Validate safety limits before execution
	if mce.safetyGuard != nil {
		if err := mce.safetyGuard.ValidateMotionCommand(command, mce.unitConverter); err != nil {
			return fmt.Errorf("safety validation failed: %w", err)
		}
	}
	
	switch command.Type {
	case types.CmdMoveAbsolute:
		return mce.ExecuteMoveAbsolute(ctx, command)
	case types.CmdMoveRelative:
		return mce.ExecuteMoveRelative(ctx, command)
	case types.CmdMoveIncremental:
		return mce.ExecuteMoveIncremental(ctx, command)
	case types.CmdJog:
		return mce.ExecuteJog(ctx, command)
	case types.CmdStop:
		return mce.ExecuteStop(ctx, command)
	default:
		return fmt.Errorf("unsupported motion command type: %s", command.Type)
	}
}

// Validate implements CommandExecutor interface
func (mce *MotionCommandExecutor) Validate(command *types.Command) error {
	return mce.ValidateMotionCommand(command)
}

// GetCommandInfo implements CommandExecutor interface
func (mce *MotionCommandExecutor) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	return mce.GetMotionCommandInfo(commandType)
}

// NewMotionCommandExecutor creates a new motion command executor
func NewMotionCommandExecutor(driveController types.DriveController, unitConverter *types.UnitConverter, safetyGuard *safety.SafetyGuard) *MotionCommandExecutor {
	return &MotionCommandExecutor{
		driveController: driveController,
		unitConverter:   unitConverter,
		safetyGuard:     safetyGuard,
	}
}

// ExecuteMoveAbsolute executes a move absolute command
func (mce *MotionCommandExecutor) ExecuteMoveAbsolute(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	position, err := extractor.ExtractPosition(command.Parameters, "position")
	if err != nil {
		return fmt.Errorf("missing position parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	posValue := mce.unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
	velValue := mce.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := mce.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := mce.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	// Execute the command
	return mce.driveController.MoveAbsolute(ctx, posValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

// ExecuteMoveRelative executes a move relative command
func (mce *MotionCommandExecutor) ExecuteMoveRelative(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	distance, err := extractor.ExtractPosition(command.Parameters, "distance")
	if err != nil {
		return fmt.Errorf("missing distance parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	distValue := mce.unitConverter.ConvertPositionValue(distance, types.PositionUnitCounts)
	velValue := mce.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := mce.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := mce.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	// Execute the command
	return mce.driveController.MoveRelative(ctx, distValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

// ExecuteMoveIncremental executes a move incremental command
func (mce *MotionCommandExecutor) ExecuteMoveIncremental(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	distance, err := extractor.ExtractPosition(command.Parameters, "distance")
	if err != nil {
		return fmt.Errorf("missing distance parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	distValue := mce.unitConverter.ConvertPositionValue(distance, types.PositionUnitCounts)
	velValue := mce.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := mce.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := mce.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	// Execute the command
	return mce.driveController.MoveIncremental(ctx, distValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

// ExecuteJog executes a jog command
func (mce *MotionCommandExecutor) ExecuteJog(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Extract required parameters
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	// Convert to drive units (counts)
	velValue := mce.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	
	// Execute the command
	return mce.driveController.Jog(ctx, velValue.Value)
}

// ExecuteStop executes a stop command
func (mce *MotionCommandExecutor) ExecuteStop(ctx context.Context, command *types.Command) error {
	// Stop command has no parameters
	return mce.driveController.Stop(ctx)
}

// ValidateMotionCommand validates motion command parameters
func (mce *MotionCommandExecutor) ValidateMotionCommand(command *types.Command) error {
	switch command.Type {
	case types.CmdMoveAbsolute, types.CmdMoveRelative, types.CmdMoveIncremental:
		return mce.validateMotionParameters(command)
	case types.CmdJog:
		return mce.validateJogParameters(command)
	case types.CmdStop:
		return nil // No parameters to validate
	default:
		return fmt.Errorf("unsupported motion command type: %s", command.Type)
	}
}

// validateMotionParameters validates parameters for motion commands
func (mce *MotionCommandExecutor) validateMotionParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	requiredParams := []string{"position", "velocity", "acceleration", "jerk"}
	if command.Type == types.CmdMoveRelative || command.Type == types.CmdMoveIncremental {
		requiredParams[0] = "distance"
	}
	
	for _, param := range requiredParams {
		if _, exists := command.Parameters[param]; !exists {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	
	// Validate position/distance parameter
	var posParam string
	if command.Type == types.CmdMoveAbsolute {
		posParam = "position"
	} else {
		posParam = "distance"
	}
	
	_, err := extractor.ExtractPosition(command.Parameters, posParam)
	if err != nil {
		return fmt.Errorf("invalid %s parameter: %w", posParam, err)
	}
	
	// Validate velocity parameter
	_, err = extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("invalid velocity parameter: %w", err)
	}
	
	// Validate acceleration parameter
	_, err = extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("invalid acceleration parameter: %w", err)
	}
	
	// Validate jerk parameter
	_, err = extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("invalid jerk parameter: %w", err)
	}
	
	return nil
}

// validateJogParameters validates parameters for jog command
func (mce *MotionCommandExecutor) validateJogParameters(command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	// Check required parameters
	if _, exists := command.Parameters["velocity"]; !exists {
		return fmt.Errorf("missing required parameter: velocity")
	}
	
	// Validate velocity parameter
	_, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("invalid velocity parameter: %w", err)
	}
	
	return nil
}

// GetMotionCommandInfo returns information about motion commands
func (mce *MotionCommandExecutor) GetMotionCommandInfo(commandType types.CommandType) (string, []string, error) {
	switch commandType {
	case types.CmdMoveAbsolute:
		return "Move to absolute position", []string{"position", "velocity", "acceleration", "jerk"}, nil
	case types.CmdMoveRelative:
		return "Move by relative distance", []string{"distance", "velocity", "acceleration", "jerk"}, nil
	case types.CmdMoveIncremental:
		return "Move by fixed increment", []string{"distance", "velocity", "acceleration", "jerk"}, nil
	case types.CmdJog:
		return "Continuous motion", []string{"velocity"}, nil
	case types.CmdStop:
		return "Stop motion", []string{}, nil
	default:
		return "", nil, fmt.Errorf("unsupported motion command type: %s", commandType)
	}
}