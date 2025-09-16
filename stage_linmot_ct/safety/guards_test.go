package safety

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewSafetyGuard(t *testing.T) {
	guard := NewSafetyGuard()
	
	if guard == nil {
		t.Fatal("Expected non-nil safety guard")
	}
	
	if guard.limits == nil {
		t.Fatal("Expected non-nil limits")
	}
	
	// Check default limits
	limits := guard.GetLimits()
	if limits.MinPosition != -1000000.0 {
		t.Errorf("Expected MinPosition -1000000.0, got %f", limits.MinPosition)
	}
	if limits.MaxPosition != 1000000.0 {
		t.Errorf("Expected MaxPosition 1000000.0, got %f", limits.MaxPosition)
	}
	if limits.MaxVelocity != 10000.0 {
		t.Errorf("Expected MaxVelocity 10000.0, got %f", limits.MaxVelocity)
	}
}

func TestNewSafetyGuardWithLimits(t *testing.T) {
	customLimits := &SafetyLimits{
		MinPosition:    -500000.0,
		MaxPosition:    500000.0,
		MaxVelocity:    5000.0,
		MinForce:       -500.0,
		MaxForce:       500.0,
		MaxAcceleration: 500.0,
		MaxJerk:        5000.0,
	}
	
	guard := NewSafetyGuardWithLimits(customLimits)
	
	if guard == nil {
		t.Fatal("Expected non-nil safety guard")
	}
	
	limits := guard.GetLimits()
	if limits.MinPosition != -500000.0 {
		t.Errorf("Expected MinPosition -500000.0, got %f", limits.MinPosition)
	}
}

func TestSafetyGuard_ValidatePosition(t *testing.T) {
	guard := NewSafetyGuard()
	
	tests := []struct {
		name      string
		position  float64
		wantError bool
	}{
		{
			name:      "Valid position within limits",
			position:  0.0,
			wantError: false,
		},
		{
			name:      "Valid position at minimum limit",
			position:  -1000000.0,
			wantError: false,
		},
		{
			name:      "Valid position at maximum limit",
			position:  1000000.0,
			wantError: false,
		},
		{
			name:      "Position below minimum limit",
			position:  -1000001.0,
			wantError: true,
		},
		{
			name:      "Position above maximum limit",
			position:  1000001.0,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.ValidatePosition(tt.position)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePosition() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSafetyGuard_ValidateVelocity(t *testing.T) {
	guard := NewSafetyGuard()
	
	tests := []struct {
		name      string
		velocity  float64
		wantError bool
	}{
		{
			name:      "Valid velocity within limits",
			velocity:  5000.0,
			wantError: false,
		},
		{
			name:      "Valid velocity at maximum limit",
			velocity:  10000.0,
			wantError: false,
		},
		{
			name:      "Zero velocity",
			velocity:  0.0,
			wantError: false,
		},
		{
			name:      "Negative velocity",
			velocity:  -1000.0,
			wantError: true,
		},
		{
			name:      "Velocity above maximum limit",
			velocity:  10001.0,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.ValidateVelocity(tt.velocity)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateVelocity() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSafetyGuard_ValidateForce(t *testing.T) {
	guard := NewSafetyGuard()
	
	tests := []struct {
		name      string
		force     float64
		wantError bool
	}{
		{
			name:      "Valid force within limits",
			force:     0.0,
			wantError: false,
		},
		{
			name:      "Valid force at minimum limit",
			force:     -1000.0,
			wantError: false,
		},
		{
			name:      "Valid force at maximum limit",
			force:     1000.0,
			wantError: false,
		},
		{
			name:      "Force below minimum limit",
			force:     -1001.0,
			wantError: true,
		},
		{
			name:      "Force above maximum limit",
			force:     1001.0,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.ValidateForce(tt.force)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateForce() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSafetyGuard_ValidateAcceleration(t *testing.T) {
	guard := NewSafetyGuard()
	
	tests := []struct {
		name         string
		acceleration float64
		wantError    bool
	}{
		{
			name:         "Valid acceleration within limits",
			acceleration: 500.0,
			wantError:    false,
		},
		{
			name:         "Valid acceleration at maximum limit",
			acceleration: 1000.0,
			wantError:    false,
		},
		{
			name:         "Zero acceleration",
			acceleration: 0.0,
			wantError:    false,
		},
		{
			name:         "Negative acceleration",
			acceleration: -100.0,
			wantError:    true,
		},
		{
			name:         "Acceleration above maximum limit",
			acceleration: 1001.0,
			wantError:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.ValidateAcceleration(tt.acceleration)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAcceleration() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSafetyGuard_ValidateJerk(t *testing.T) {
	guard := NewSafetyGuard()
	
	tests := []struct {
		name      string
		jerk      float64
		wantError bool
	}{
		{
			name:      "Valid jerk within limits",
			jerk:      5000.0,
			wantError: false,
		},
		{
			name:      "Valid jerk at maximum limit",
			jerk:      10000.0,
			wantError: false,
		},
		{
			name:      "Zero jerk",
			jerk:      0.0,
			wantError: false,
		},
		{
			name:      "Negative jerk",
			jerk:      -1000.0,
			wantError: true,
		},
		{
			name:      "Jerk above maximum limit",
			jerk:      10001.0,
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guard.ValidateJerk(tt.jerk)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateJerk() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSafetyGuard_SetLimits(t *testing.T) {
	guard := NewSafetyGuard()
	
	// Test valid limits
	validLimits := &SafetyLimits{
		MinPosition:    -500000.0,
		MaxPosition:    500000.0,
		MaxVelocity:    5000.0,
		MinForce:       -500.0,
		MaxForce:       500.0,
		MaxAcceleration: 500.0,
		MaxJerk:        5000.0,
	}
	
	err := guard.SetLimits(validLimits)
	if err != nil {
		t.Errorf("SetLimits() error = %v", err)
	}
	
	// Test invalid limits
	invalidLimits := &SafetyLimits{
		MinPosition:    500000.0,  // Min > Max
		MaxPosition:    -500000.0,
		MaxVelocity:    5000.0,
		MinForce:       -500.0,
		MaxForce:       500.0,
		MaxAcceleration: 500.0,
		MaxJerk:        5000.0,
	}
	
	err = guard.SetLimits(invalidLimits)
	if err == nil {
		t.Error("Expected error for invalid limits, got nil")
	}
	
	// Test nil limits
	err = guard.SetLimits(nil)
	if err == nil {
		t.Error("Expected error for nil limits, got nil")
	}
}

func TestSafetyGuard_IsEmergencyStopActive(t *testing.T) {
	guard := NewSafetyGuard()
	
	// This should return false by default
	if guard.IsEmergencyStopActive() {
		t.Error("Expected emergency stop to be inactive by default")
	}
}

func TestSafetyGuard_TriggerEmergencyStop(t *testing.T) {
	guard := NewSafetyGuard()
	
	// Create a mock drive controller
	mockDrive := &MockDriveController{}
	
	ctx := context.Background()
	err := guard.TriggerEmergencyStop(ctx, mockDrive, "Test emergency stop")
	
	if err != nil {
		t.Errorf("TriggerEmergencyStop() error = %v", err)
	}
}

// MockDriveController for testing
type MockDriveController struct {
	driveState     types.DriveState
	motionComplete bool
	position       float64
	velocity       float64
	force          float64
}

func (mdc *MockDriveController) Stop(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) GetPosition(ctx context.Context) (float64, error) {
	return mdc.position, nil
}

func (mdc *MockDriveController) GetVelocity(ctx context.Context) (float64, error) {
	return mdc.velocity, nil
}

func (mdc *MockDriveController) GetForce(ctx context.Context) (float64, error) {
	return mdc.force, nil
}

func (mdc *MockDriveController) GetDriveState(ctx context.Context) (types.DriveState, error) {
	return mdc.driveState, nil
}

func (mdc *MockDriveController) IsMotionComplete(ctx context.Context) (bool, error) {
	return mdc.motionComplete, nil
}

// Implement other required methods as no-ops for testing
func (mdc *MockDriveController) MoveAbsolute(ctx context.Context, position float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) MoveRelative(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) MoveIncremental(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) Jog(ctx context.Context, velocity float64) error {
	return nil
}

func (mdc *MockDriveController) Wait(ctx context.Context, duration time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitPosition(ctx context.Context, position float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) SetDigitalOutput(ctx context.Context, output int, value bool) error {
	return nil
}

func (mdc *MockDriveController) ClearDigitalOutput(ctx context.Context, output int) error {
	return nil
}

func (mdc *MockDriveController) SetAnalogOutput(ctx context.Context, output int, value float64) error {
	return nil
}

func (mdc *MockDriveController) WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) Home(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) Reset(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SaveConfiguration(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) LoadConfiguration(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) ForceControlOn(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) ForceControlOff(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SetForce(ctx context.Context, force float64) error {
	return nil
}

func (mdc *MockDriveController) StartOscilloscope(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) StopOscilloscope(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SaveData(ctx context.Context, filename string) error {
	return nil
}

func (mdc *MockDriveController) GetDigitalInput(ctx context.Context, input int) (bool, error) {
	return false, nil
}

func (mdc *MockDriveController) GetAnalogInput(ctx context.Context, input int) (float64, error) {
	return 0.0, nil
}

func (mdc *MockDriveController) GetDigitalOutput(ctx context.Context, output int) (bool, error) {
	return false, nil
}

func (mdc *MockDriveController) GetAnalogOutput(ctx context.Context, output int) (float64, error) {
	return 0.0, nil
}