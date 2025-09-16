package safety

import (
	"math"
	"testing"
)

func TestNewLimitChecker(t *testing.T) {
	limits := &SafetyLimits{
		MinPosition:    -1000000.0,
		MaxPosition:    1000000.0,
		MaxVelocity:    10000.0,
		MinForce:       -1000.0,
		MaxForce:       1000.0,
		MaxAcceleration: 1000.0,
		MaxJerk:        10000.0,
	}
	
	checker := NewLimitChecker(limits)
	
	if checker == nil {
		t.Fatal("Expected non-nil limit checker")
	}
	
	if checker.limits != limits {
		t.Error("Expected limits to be set correctly")
	}
}

func TestLimitChecker_CheckPositionLimits(t *testing.T) {
	limits := &SafetyLimits{
		MinPosition: -1000000.0,
		MaxPosition: 1000000.0,
	}
	checker := NewLimitChecker(limits)
	
	tests := []struct {
		name        string
		position    float64
		wantViolations int
	}{
		{
			name:        "Valid position within limits",
			position:    0.0,
			wantViolations: 0,
		},
		{
			name:        "Valid position at minimum limit",
			position:    -1000000.0,
			wantViolations: 0,
		},
		{
			name:        "Valid position at maximum limit",
			position:    1000000.0,
			wantViolations: 0,
		},
		{
			name:        "Position below minimum limit",
			position:    -1000001.0,
			wantViolations: 1,
		},
		{
			name:        "Position above maximum limit",
			position:    1000001.0,
			wantViolations: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := checker.CheckPositionLimits(tt.position)
			if len(violations) != tt.wantViolations {
				t.Errorf("CheckPositionLimits() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestLimitChecker_CheckVelocityLimits(t *testing.T) {
	limits := &SafetyLimits{
		MaxVelocity: 10000.0,
	}
	checker := NewLimitChecker(limits)
	
	tests := []struct {
		name        string
		velocity    float64
		wantViolations int
	}{
		{
			name:        "Valid velocity within limits",
			velocity:    5000.0,
			wantViolations: 0,
		},
		{
			name:        "Valid velocity at maximum limit",
			velocity:    10000.0,
			wantViolations: 0,
		},
		{
			name:        "Zero velocity",
			velocity:    0.0,
			wantViolations: 0,
		},
		{
			name:        "Negative velocity",
			velocity:    -1000.0,
			wantViolations: 1,
		},
		{
			name:        "Velocity above maximum limit",
			velocity:    10001.0,
			wantViolations: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := checker.CheckVelocityLimits(tt.velocity)
			if len(violations) != tt.wantViolations {
				t.Errorf("CheckVelocityLimits() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestLimitChecker_CheckForceLimits(t *testing.T) {
	limits := &SafetyLimits{
		MinForce: -1000.0,
		MaxForce: 1000.0,
	}
	checker := NewLimitChecker(limits)
	
	tests := []struct {
		name        string
		force       float64
		wantViolations int
	}{
		{
			name:        "Valid force within limits",
			force:       0.0,
			wantViolations: 0,
		},
		{
			name:        "Valid force at minimum limit",
			force:       -1000.0,
			wantViolations: 0,
		},
		{
			name:        "Valid force at maximum limit",
			force:       1000.0,
			wantViolations: 0,
		},
		{
			name:        "Force below minimum limit",
			force:       -1001.0,
			wantViolations: 1,
		},
		{
			name:        "Force above maximum limit",
			force:       1001.0,
			wantViolations: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := checker.CheckForceLimits(tt.force)
			if len(violations) != tt.wantViolations {
				t.Errorf("CheckForceLimits() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestLimitChecker_CheckAllLimits(t *testing.T) {
	limits := &SafetyLimits{
		MinPosition:    -1000000.0,
		MaxPosition:    1000000.0,
		MaxVelocity:    10000.0,
		MinForce:       -1000.0,
		MaxForce:       1000.0,
		MaxAcceleration: 1000.0,
		MaxJerk:        10000.0,
	}
	checker := NewLimitChecker(limits)
	
	tests := []struct {
		name        string
		position    float64
		velocity    float64
		force       float64
		acceleration float64
		jerk        float64
		wantViolations int
	}{
		{
			name:        "All valid parameters",
			position:    0.0,
			velocity:    5000.0,
			force:       0.0,
			acceleration: 500.0,
			jerk:        5000.0,
			wantViolations: 0,
		},
		{
			name:        "Position out of range",
			position:    1000001.0,
			velocity:    5000.0,
			force:       0.0,
			acceleration: 500.0,
			jerk:        5000.0,
			wantViolations: 1,
		},
		{
			name:        "Multiple violations",
			position:    1000001.0,
			velocity:    10001.0,
			force:       1001.0,
			acceleration: 1001.0,
			jerk:        10001.0,
			wantViolations: 5,
		},
		{
			name:        "NaN values should be ignored",
			position:    math.NaN(),
			velocity:    math.NaN(),
			force:       math.NaN(),
			acceleration: math.NaN(),
			jerk:        math.NaN(),
			wantViolations: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := checker.CheckAllLimits(tt.position, tt.velocity, tt.force, tt.acceleration, tt.jerk)
			if len(violations) != tt.wantViolations {
				t.Errorf("CheckAllLimits() violations = %d, want %d", len(violations), tt.wantViolations)
			}
		})
	}
}

func TestLimitChecker_ValidateLimits(t *testing.T) {
	checker := NewLimitChecker(&SafetyLimits{})
	
	tests := []struct {
		name    string
		limits  *SafetyLimits
		wantErr bool
	}{
		{
			name: "Valid limits",
			limits: &SafetyLimits{
				MinPosition:    -1000000.0,
				MaxPosition:    1000000.0,
				MaxVelocity:    10000.0,
				MinForce:       -1000.0,
				MaxForce:       1000.0,
				MaxAcceleration: 1000.0,
				MaxJerk:        10000.0,
			},
			wantErr: false,
		},
		{
			name: "Invalid position limits",
			limits: &SafetyLimits{
				MinPosition:    1000000.0,  // Min > Max
				MaxPosition:    -1000000.0,
				MaxVelocity:    10000.0,
				MinForce:       -1000.0,
				MaxForce:       1000.0,
				MaxAcceleration: 1000.0,
				MaxJerk:        10000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid velocity limits",
			limits: &SafetyLimits{
				MinPosition:    -1000000.0,
				MaxPosition:    1000000.0,
				MaxVelocity:    0.0,  // Max velocity <= 0
				MinForce:       -1000.0,
				MaxForce:       1000.0,
				MaxAcceleration: 1000.0,
				MaxJerk:        10000.0,
			},
			wantErr: true,
		},
		{
			name: "Invalid force limits",
			limits: &SafetyLimits{
				MinPosition:    -1000000.0,
				MaxPosition:    1000000.0,
				MaxVelocity:    10000.0,
				MinForce:       1000.0,  // Min > Max
				MaxForce:       -1000.0,
				MaxAcceleration: 1000.0,
				MaxJerk:        10000.0,
			},
			wantErr: true,
		},
		{
			name: "Unreasonably large position",
			limits: &SafetyLimits{
				MinPosition:    -1000000.0,
				MaxPosition:    2e6,  // Too large
				MaxVelocity:    10000.0,
				MinForce:       -1000.0,
				MaxForce:       1000.0,
				MaxAcceleration: 1000.0,
				MaxJerk:        10000.0,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker.limits = tt.limits
			err := checker.ValidateLimits()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLimits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLimitChecker_SetLimits(t *testing.T) {
	checker := NewLimitChecker(&SafetyLimits{})
	
	validLimits := &SafetyLimits{
		MinPosition:    -500000.0,
		MaxPosition:    500000.0,
		MaxVelocity:    5000.0,
		MinForce:       -500.0,
		MaxForce:       500.0,
		MaxAcceleration: 500.0,
		MaxJerk:        5000.0,
	}
	
	err := checker.SetLimits(validLimits)
	if err != nil {
		t.Errorf("SetLimits() error = %v", err)
	}
	
	if checker.limits != validLimits {
		t.Error("Expected limits to be set correctly")
	}
	
	// Test with invalid limits
	invalidLimits := &SafetyLimits{
		MinPosition:    500000.0,  // Min > Max
		MaxPosition:    -500000.0,
		MaxVelocity:    5000.0,
		MinForce:       -500.0,
		MaxForce:       500.0,
		MaxAcceleration: 500.0,
		MaxJerk:        5000.0,
	}
	
	err = checker.SetLimits(invalidLimits)
	if err == nil {
		t.Error("Expected error for invalid limits, got nil")
	}
}

func TestLimitViolation_String(t *testing.T) {
	violation := &LimitViolation{
		Type:      LimitTypePosition,
		Value:     1000001.0,
		Limit:     1000000.0,
		IsMinimum: false,
		Message:   "position 1000001.0 exceeds maximum limit 1000000.0",
	}
	
	expected := "position maximum limit violation: value 1000001.000000 exceeds maximum limit 1000000.000000"
	if violation.String() != expected {
		t.Errorf("String() = %v, want %v", violation.String(), expected)
	}
}

func TestLimitViolation_Error(t *testing.T) {
	violation := &LimitViolation{
		Type:      LimitTypeVelocity,
		Value:     -1000.0,
		Limit:     0.0,
		IsMinimum: true,
		Message:   "velocity -1000.0 is negative",
	}
	
	expected := "velocity minimum limit violation: value -1000.000000 exceeds minimum limit 0.000000"
	if violation.Error() != expected {
		t.Errorf("Error() = %v, want %v", violation.Error(), expected)
	}
}