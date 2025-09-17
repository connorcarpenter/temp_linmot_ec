package safety

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewPreconditionChecker(t *testing.T) {
	mockDrive := &MockDriveController{}
	guard := NewSafetyGuard()
	
	checker := NewPreconditionChecker(mockDrive, guard)
	
	if checker == nil {
		t.Fatal("Expected non-nil precondition checker")
	}
	
	if checker.driveController != mockDrive {
		t.Error("Expected drive controller to be set correctly")
	}
	
	if checker.safetyGuard != guard {
		t.Error("Expected safety guard to be set correctly")
	}
}

func TestPreconditionChecker_CheckDriveState(t *testing.T) {
	tests := []struct {
		name        string
		commandType types.CommandType
		driveState  types.DriveState
		wantError   bool
	}{
		{
			name:        "Motion command with ready drive",
			commandType: types.CmdMoveAbsolute,
			driveState:  types.DriveStateReady,
			wantError:   false,
		},
		{
			name:        "Motion command with moving drive",
			commandType: types.CmdMoveAbsolute,
			driveState:  types.DriveStateMoving,
			wantError:   false,
		},
		{
			name:        "Motion command with error drive",
			commandType: types.CmdMoveAbsolute,
			driveState:  types.DriveStateError,
			wantError:   true,
		},
		{
			name:        "Stop command with any state",
			commandType: types.CmdStop,
			driveState:  types.DriveStateError,
			wantError:   false,
		},
		{
			name:        "Force control with error drive",
			commandType: types.CmdForceControlOn,
			driveState:  types.DriveStateError,
			wantError:   true,
		},
		{
			name:        "System command with moving drive",
			commandType: types.CmdHome,
			driveState:  types.DriveStateMoving,
			wantError:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDrive := &MockDriveController{
				driveState: tt.driveState,
			}
			guard := NewSafetyGuard()
			checker := NewPreconditionChecker(mockDrive, guard)
			
			command := &types.Command{
				Type: tt.commandType,
			}
			
			ctx := context.Background()
			err := checker.CheckDriveState(ctx, command)
			
			if (err != nil) != tt.wantError {
				t.Errorf("CheckDriveState() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPreconditionChecker_CheckMotionPreconditions(t *testing.T) {
	tests := []struct {
		name           string
		motionComplete bool
		wantError      bool
	}{
		{
			name:           "Motion complete",
			motionComplete: true,
			wantError:      false,
		},
		{
			name:           "Motion not complete",
			motionComplete: false,
			wantError:      true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDrive := &MockDriveController{
				motionComplete: tt.motionComplete,
			}
			guard := NewSafetyGuard()
			checker := NewPreconditionChecker(mockDrive, guard)
			
			command := &types.Command{
				Type: types.CmdMoveAbsolute,
			}
			
			ctx := context.Background()
			err := checker.CheckMotionPreconditions(ctx, command)
			
			if (err != nil) != tt.wantError {
				t.Errorf("CheckMotionPreconditions() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPreconditionChecker_CheckSafetyPreconditions(t *testing.T) {
	tests := []struct {
		name        string
		position    float64
		emergencyStop bool
		wantError   bool
	}{
		{
			name:        "Valid position, no emergency stop",
			position:    0.0,
			emergencyStop: false,
			wantError:   false,
		},
		{
			name:        "Position out of range",
			position:    1000001.0,
			emergencyStop: false,
			wantError:   true,
		},
		{
			name:        "Emergency stop active",
			position:    0.0,
			emergencyStop: true,
			wantError:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDrive := &MockDriveController{
				position: tt.position,
			}
			guard := NewSafetyGuard()
			// Override emergency stop status for testing
			guard.emergencyStopActive = tt.emergencyStop
			checker := NewPreconditionChecker(mockDrive, guard)
			
			command := &types.Command{
				Type: types.CmdMoveAbsolute,
			}
			
			ctx := context.Background()
			err := checker.CheckSafetyPreconditions(ctx, command)
			
			if (err != nil) != tt.wantError {
				t.Errorf("CheckSafetyPreconditions() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPreconditionChecker_CheckParameterPreconditions(t *testing.T) {
	tests := []struct {
		name        string
		commandType types.CommandType
		parameters  map[string]interface{}
		wantError   bool
	}{
		{
			name:        "MoveAbsolute with position parameter",
			commandType: types.CmdMoveAbsolute,
			parameters: map[string]interface{}{
				"position": &types.PositionValue{Value: 1000.0, Unit: types.PositionUnitCounts},
			},
			wantError: false,
		},
		{
			name:        "MoveAbsolute without position parameter",
			commandType: types.CmdMoveAbsolute,
			parameters:  map[string]interface{}{},
			wantError:   true,
		},
		{
			name:        "Jog with velocity parameter",
			commandType: types.CmdJog,
		parameters: map[string]interface{}{
			"velocity": &types.VelocityValue{Value: 1000.0, Unit: types.VelocityUnitCountsS},
		},
			wantError: false,
		},
		{
			name:        "Jog without velocity parameter",
			commandType: types.CmdJog,
			parameters:  map[string]interface{}{},
			wantError:   true,
		},
		{
			name:        "SetForce with force parameter",
			commandType: types.CmdSetForce,
			parameters: map[string]interface{}{
				"force": &types.ForceValue{Value: 100.0, Unit: types.ForceUnitCounts},
			},
			wantError: false,
		},
		{
			name:        "SetForce without force parameter",
			commandType: types.CmdSetForce,
			parameters:  map[string]interface{}{},
			wantError:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDrive := &MockDriveController{}
			guard := NewSafetyGuard()
			checker := NewPreconditionChecker(mockDrive, guard)
			
			command := &types.Command{
				Type:       tt.commandType,
				Parameters: tt.parameters,
			}
			
			ctx := context.Background()
			err := checker.CheckParameterPreconditions(ctx, command)
			
			if (err != nil) != tt.wantError {
				t.Errorf("CheckParameterPreconditions() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPreconditionChecker_CheckAllPreconditions(t *testing.T) {
	mockDrive := &MockDriveController{
		driveState:     types.DriveStateReady,
		motionComplete: true,
		position:       0.0,
	}
	guard := NewSafetyGuard()
	checker := NewPreconditionChecker(mockDrive, guard)
	
	command := &types.Command{
		Type: types.CmdMoveAbsolute,
		Parameters: map[string]interface{}{
			"position": &types.PositionValue{Value: 1000.0, Unit: types.PositionUnitCounts},
		},
	}
	
	ctx := context.Background()
	err := checker.CheckAllPreconditions(ctx, command)
	
	if err != nil {
		t.Errorf("CheckAllPreconditions() error = %v", err)
	}
}

func TestErrorRecovery_RecoverFromError(t *testing.T) {
	mockDrive := &MockDriveController{}
	recovery := NewErrorRecovery(mockDrive)
	
	tests := []struct {
		name    string
		command *types.Command
		err     error
		wantErr bool
	}{
		{
			name: "Recoverable motion state error",
			command: &types.Command{
				Type: types.CmdMoveAbsolute,
			},
			err: &PreconditionError{
				Type:    "motion_state",
				Message: "previous motion is not complete",
			},
			wantErr: true, // Will not be recoverable with mock drive
		},
		{
			name: "Non-recoverable error",
			command: &types.Command{
				Type: types.CmdMoveAbsolute,
			},
			err: &PreconditionError{
				Type:    "safety",
				Message: "emergency stop is active",
			},
			wantErr: true, // Should not be recoverable
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := recovery.RecoverFromError(ctx, tt.command, tt.err)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("RecoverFromError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrorRecovery_SetMaxRetries(t *testing.T) {
	mockDrive := &MockDriveController{}
	recovery := NewErrorRecovery(mockDrive)
	
	recovery.SetMaxRetries(5)
	
	if recovery.maxRetries != 5 {
		t.Errorf("Expected maxRetries to be 5, got %d", recovery.maxRetries)
	}
}

func TestErrorRecovery_SetRetryDelay(t *testing.T) {
	mockDrive := &MockDriveController{}
	recovery := NewErrorRecovery(mockDrive)
	
	delay := 2 * time.Second
	recovery.SetRetryDelay(delay)
	
	if recovery.retryDelay != delay {
		t.Errorf("Expected retryDelay to be %v, got %v", delay, recovery.retryDelay)
	}
}

func TestPreconditionError_Error(t *testing.T) {
	err := &PreconditionError{
		Type:    "drive_state",
		Message: "Drive is not ready",
		Details: map[string]interface{}{
			"command_type": "MoveAbsolute",
		},
	}
	
	expected := "precondition check failed [drive_state]: Drive is not ready"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

// MockDriveController for testing - using the one from guards_test.go