package safety

import (
	"context"
	"fmt"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// PreconditionChecker validates preconditions before command execution
type PreconditionChecker struct {
	driveController types.DriveController
	safetyGuard    *SafetyGuard
}

// NewPreconditionChecker creates a new precondition checker
func NewPreconditionChecker(driveController types.DriveController, safetyGuard *SafetyGuard) *PreconditionChecker {
	return &PreconditionChecker{
		driveController: driveController,
		safetyGuard:    safetyGuard,
	}
}

// PreconditionError represents a precondition validation error
type PreconditionError struct {
	Type    string
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (pe *PreconditionError) Error() string {
	return fmt.Sprintf("precondition check failed [%s]: %s", pe.Type, pe.Message)
}

// CheckDriveState validates that the drive is in a valid state for command execution
func (pc *PreconditionChecker) CheckDriveState(ctx context.Context, command *types.Command) error {
	driveState, err := pc.driveController.GetDriveState(ctx)
	if err != nil {
		return &PreconditionError{
			Type:    "drive_state",
			Message: fmt.Sprintf("failed to get drive state: %v", err),
			Details: map[string]interface{}{
				"command_type": command.Type,
				"error":        err.Error(),
			},
		}
	}
	
	// Check if drive is in a valid state for the command type
	switch command.Type {
	case types.CmdMoveAbsolute, types.CmdMoveRelative, types.CmdMoveIncremental, types.CmdJog:
		if driveState != types.DriveStateReady && driveState != types.DriveStateMoving {
			return &PreconditionError{
				Type:    "drive_state",
				Message: fmt.Sprintf("drive must be ready or moving for motion commands, current state: %s", driveState),
				Details: map[string]interface{}{
					"command_type":    command.Type,
					"current_state":   driveState,
					"required_states": []types.DriveState{types.DriveStateReady, types.DriveStateMoving},
				},
			}
		}
		
	case types.CmdStop:
		// Stop command can be executed in any state
		return nil
		
	case types.CmdForceControlOn, types.CmdForceControlOff, types.CmdSetForce:
		if driveState == types.DriveStateError {
			return &PreconditionError{
				Type:    "drive_state",
				Message: fmt.Sprintf("drive is in error state, cannot execute force control commands: %s", driveState),
				Details: map[string]interface{}{
					"command_type":  command.Type,
					"current_state": driveState,
				},
			}
		}
		
	case types.CmdHome, types.CmdReset:
		if driveState == types.DriveStateMoving {
			return &PreconditionError{
				Type:    "drive_state",
				Message: fmt.Sprintf("drive is moving, cannot execute system commands: %s", driveState),
				Details: map[string]interface{}{
					"command_type":  command.Type,
					"current_state": driveState,
				},
			}
		}
		
	default:
		// For other commands, just check that drive is not in error state
		if driveState == types.DriveStateError {
			return &PreconditionError{
				Type:    "drive_state",
				Message: fmt.Sprintf("drive is in error state: %s", driveState),
				Details: map[string]interface{}{
					"command_type":  command.Type,
					"current_state": driveState,
				},
			}
		}
	}
	
	return nil
}

// CheckMotionPreconditions validates motion-specific preconditions
func (pc *PreconditionChecker) CheckMotionPreconditions(ctx context.Context, command *types.Command) error {
	// Check if previous motion is complete
	motionComplete, err := pc.driveController.IsMotionComplete(ctx)
	if err != nil {
		return &PreconditionError{
			Type:    "motion_state",
			Message: fmt.Sprintf("failed to check motion completion: %v", err),
			Details: map[string]interface{}{
				"command_type": command.Type,
				"error":        err.Error(),
			},
		}
	}
	
	if !motionComplete {
		return &PreconditionError{
			Type:    "motion_state",
			Message: "previous motion is not complete, cannot start new motion",
			Details: map[string]interface{}{
				"command_type":    command.Type,
				"motion_complete": motionComplete,
			},
		}
	}
	
	return nil
}

// CheckSafetyPreconditions validates safety-related preconditions
func (pc *PreconditionChecker) CheckSafetyPreconditions(ctx context.Context, command *types.Command) error {
	// Check if emergency stop is active
	if pc.safetyGuard.IsEmergencyStopActive() {
		return &PreconditionError{
			Type:    "safety",
			Message: "emergency stop is active, cannot execute commands",
			Details: map[string]interface{}{
				"command_type": command.Type,
			},
		}
	}
	
	// Check current position is within limits
	currentPosition, err := pc.driveController.GetPosition(ctx)
	if err != nil {
		return &PreconditionError{
			Type:    "safety",
			Message: fmt.Sprintf("failed to get current position: %v", err),
			Details: map[string]interface{}{
				"command_type": command.Type,
				"error":        err.Error(),
			},
		}
	}
	
	if err := pc.safetyGuard.ValidatePosition(currentPosition); err != nil {
		return &PreconditionError{
			Type:    "safety",
			Message: fmt.Sprintf("current position is outside safety limits: %v", err),
			Details: map[string]interface{}{
				"command_type":    command.Type,
				"current_position": currentPosition,
				"error":           err.Error(),
			},
		}
	}
	
	return nil
}

// CheckParameterPreconditions validates parameter-specific preconditions
func (pc *PreconditionChecker) CheckParameterPreconditions(ctx context.Context, command *types.Command) error {
	pe := types.NewParameterExtractor()
	
	// Check for required parameters based on command type
	switch command.Type {
	case types.CmdMoveAbsolute, types.CmdMoveRelative, types.CmdMoveIncremental:
		// These commands require position parameter
		if _, err := pe.ExtractPosition(command.Parameters, "position"); err != nil {
			return &PreconditionError{
				Type:    "parameters",
				Message: fmt.Sprintf("missing required position parameter: %v", err),
				Details: map[string]interface{}{
					"command_type": command.Type,
					"parameter":    "position",
					"error":        err.Error(),
				},
			}
		}
		
	case types.CmdJog:
		// Jog command requires velocity parameter
		if _, err := pe.ExtractVelocity(command.Parameters, "velocity"); err != nil {
			return &PreconditionError{
				Type:    "parameters",
				Message: fmt.Sprintf("missing required velocity parameter: %v", err),
				Details: map[string]interface{}{
					"command_type": command.Type,
					"parameter":    "velocity",
					"error":        err.Error(),
				},
			}
		}
		
	case types.CmdSetForce:
		// SetForce command requires force parameter
		if _, err := pe.ExtractForce(command.Parameters, "force"); err != nil {
			return &PreconditionError{
				Type:    "parameters",
				Message: fmt.Sprintf("missing required force parameter: %v", err),
				Details: map[string]interface{}{
					"command_type": command.Type,
					"parameter":    "force",
					"error":        err.Error(),
				},
			}
		}
	}
	
	return nil
}

// CheckAllPreconditions performs all precondition checks
func (pc *PreconditionChecker) CheckAllPreconditions(ctx context.Context, command *types.Command) error {
	// Check drive state
	if err := pc.CheckDriveState(ctx, command); err != nil {
		return err
	}
	
	// Check motion preconditions for motion commands
	if isMotionCommand(command.Type) {
		if err := pc.CheckMotionPreconditions(ctx, command); err != nil {
			return err
		}
	}
	
	// Check safety preconditions
	if err := pc.CheckSafetyPreconditions(ctx, command); err != nil {
		return err
	}
	
	// Check parameter preconditions
	if err := pc.CheckParameterPreconditions(ctx, command); err != nil {
		return err
	}
	
	return nil
}

// isMotionCommand checks if a command type is a motion command
func isMotionCommand(commandType types.CommandType) bool {
	switch commandType {
	case types.CmdMoveAbsolute, types.CmdMoveRelative, types.CmdMoveIncremental, types.CmdJog, types.CmdStop:
		return true
	default:
		return false
	}
}

// ErrorRecovery provides error recovery capabilities
type ErrorRecovery struct {
	driveController types.DriveController
	maxRetries      int
	retryDelay      time.Duration
}

// NewErrorRecovery creates a new error recovery handler
func NewErrorRecovery(driveController types.DriveController) *ErrorRecovery {
	return &ErrorRecovery{
		driveController: driveController,
		maxRetries:      3,
		retryDelay:      time.Second,
	}
}

// RecoverFromError attempts to recover from an error
func (er *ErrorRecovery) RecoverFromError(ctx context.Context, command *types.Command, err error) error {
	// Check if it's a precondition error that might be recoverable
	if preErr, ok := err.(*PreconditionError); ok {
		switch preErr.Type {
		case "motion_state":
			// Try to stop any ongoing motion and retry
			if stopErr := er.driveController.Stop(ctx); stopErr != nil {
				return fmt.Errorf("failed to stop motion during recovery: %w", stopErr)
			}
			
			// Wait a bit for the motion to stop
			time.Sleep(er.retryDelay)
			
			// Check if motion is now complete
			motionComplete, checkErr := er.driveController.IsMotionComplete(ctx)
			if checkErr != nil {
				return fmt.Errorf("failed to check motion completion during recovery: %w", checkErr)
			}
			
			if motionComplete {
				return nil // Recovery successful
			}
			
			// If still not complete, return the original error
			return fmt.Errorf("motion still not complete after recovery attempt: %w", err)
			
		case "drive_state":
			// For drive state errors, we might need to reset the drive
			if command.Type == types.CmdReset {
				// If it's a reset command, we can try to execute it
				return nil
			}
		}
	}
	
	return fmt.Errorf("error not recoverable: %w", err)
}

// SetMaxRetries sets the maximum number of retry attempts
func (er *ErrorRecovery) SetMaxRetries(maxRetries int) {
	er.maxRetries = maxRetries
}

// SetRetryDelay sets the delay between retry attempts
func (er *ErrorRecovery) SetRetryDelay(delay time.Duration) {
	er.retryDelay = delay
}