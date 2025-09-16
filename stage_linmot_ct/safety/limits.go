package safety

import (
	"fmt"
	"math"
)

// LimitType represents the type of safety limit
type LimitType int

const (
	LimitTypePosition LimitType = iota
	LimitTypeVelocity
	LimitTypeForce
	LimitTypeAcceleration
	LimitTypeJerk
)

// LimitViolation represents a safety limit violation
type LimitViolation struct {
	Type      LimitType
	Value     float64
	Limit     float64
	IsMinimum bool
	Message   string
}

// String returns a string representation of the limit violation
func (lv *LimitViolation) String() string {
	limitType := "unknown"
	switch lv.Type {
	case LimitTypePosition:
		limitType = "position"
	case LimitTypeVelocity:
		limitType = "velocity"
	case LimitTypeForce:
		limitType = "force"
	case LimitTypeAcceleration:
		limitType = "acceleration"
	case LimitTypeJerk:
		limitType = "jerk"
	}
	
	limitKind := "maximum"
	if lv.IsMinimum {
		limitKind = "minimum"
	}
	
	return fmt.Sprintf("%s %s limit violation: value %f exceeds %s limit %f", 
		limitType, limitKind, lv.Value, limitKind, lv.Limit)
}

// Error implements the error interface
func (lv *LimitViolation) Error() string {
	return lv.String()
}

// LimitChecker provides advanced limit checking capabilities
type LimitChecker struct {
	limits *SafetyLimits
}

// NewLimitChecker creates a new limit checker
func NewLimitChecker(limits *SafetyLimits) *LimitChecker {
	return &LimitChecker{
		limits: limits,
	}
}

// CheckPositionLimits checks position against all applicable limits
func (lc *LimitChecker) CheckPositionLimits(position float64) []*LimitViolation {
	var violations []*LimitViolation
	
	if position < lc.limits.MinPosition {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypePosition,
			Value:     position,
			Limit:     lc.limits.MinPosition,
			IsMinimum: true,
			Message:   fmt.Sprintf("position %f is below minimum limit %f", position, lc.limits.MinPosition),
		})
	}
	
	if position > lc.limits.MaxPosition {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypePosition,
			Value:     position,
			Limit:     lc.limits.MaxPosition,
			IsMinimum: false,
			Message:   fmt.Sprintf("position %f exceeds maximum limit %f", position, lc.limits.MaxPosition),
		})
	}
	
	return violations
}

// CheckVelocityLimits checks velocity against all applicable limits
func (lc *LimitChecker) CheckVelocityLimits(velocity float64) []*LimitViolation {
	var violations []*LimitViolation
	
	if velocity < 0 {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeVelocity,
			Value:     velocity,
			Limit:     0,
			IsMinimum: true,
			Message:   fmt.Sprintf("velocity %f is negative", velocity),
		})
	}
	
	if velocity > lc.limits.MaxVelocity {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeVelocity,
			Value:     velocity,
			Limit:     lc.limits.MaxVelocity,
			IsMinimum: false,
			Message:   fmt.Sprintf("velocity %f exceeds maximum limit %f", velocity, lc.limits.MaxVelocity),
		})
	}
	
	return violations
}

// CheckForceLimits checks force against all applicable limits
func (lc *LimitChecker) CheckForceLimits(force float64) []*LimitViolation {
	var violations []*LimitViolation
	
	if force < lc.limits.MinForce {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeForce,
			Value:     force,
			Limit:     lc.limits.MinForce,
			IsMinimum: true,
			Message:   fmt.Sprintf("force %f is below minimum limit %f", force, lc.limits.MinForce),
		})
	}
	
	if force > lc.limits.MaxForce {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeForce,
			Value:     force,
			Limit:     lc.limits.MaxForce,
			IsMinimum: false,
			Message:   fmt.Sprintf("force %f exceeds maximum limit %f", force, lc.limits.MaxForce),
		})
	}
	
	return violations
}

// CheckAccelerationLimits checks acceleration against all applicable limits
func (lc *LimitChecker) CheckAccelerationLimits(acceleration float64) []*LimitViolation {
	var violations []*LimitViolation
	
	if acceleration < 0 {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeAcceleration,
			Value:     acceleration,
			Limit:     0,
			IsMinimum: true,
			Message:   fmt.Sprintf("acceleration %f is negative", acceleration),
		})
	}
	
	if acceleration > lc.limits.MaxAcceleration {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeAcceleration,
			Value:     acceleration,
			Limit:     lc.limits.MaxAcceleration,
			IsMinimum: false,
			Message:   fmt.Sprintf("acceleration %f exceeds maximum limit %f", acceleration, lc.limits.MaxAcceleration),
		})
	}
	
	return violations
}

// CheckJerkLimits checks jerk against all applicable limits
func (lc *LimitChecker) CheckJerkLimits(jerk float64) []*LimitViolation {
	var violations []*LimitViolation
	
	if jerk < 0 {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeJerk,
			Value:     jerk,
			Limit:     0,
			IsMinimum: true,
			Message:   fmt.Sprintf("jerk %f is negative", jerk),
		})
	}
	
	if jerk > lc.limits.MaxJerk {
		violations = append(violations, &LimitViolation{
			Type:      LimitTypeJerk,
			Value:     jerk,
			Limit:     lc.limits.MaxJerk,
			IsMinimum: false,
			Message:   fmt.Sprintf("jerk %f exceeds maximum limit %f", jerk, lc.limits.MaxJerk),
		})
	}
	
	return violations
}

// CheckAllLimits checks all parameters against their respective limits
func (lc *LimitChecker) CheckAllLimits(position, velocity, force, acceleration, jerk float64) []*LimitViolation {
	var allViolations []*LimitViolation
	
	// Check each parameter if it's not NaN (indicating it was provided)
	if !math.IsNaN(position) {
		allViolations = append(allViolations, lc.CheckPositionLimits(position)...)
	}
	if !math.IsNaN(velocity) {
		allViolations = append(allViolations, lc.CheckVelocityLimits(velocity)...)
	}
	if !math.IsNaN(force) {
		allViolations = append(allViolations, lc.CheckForceLimits(force)...)
	}
	if !math.IsNaN(acceleration) {
		allViolations = append(allViolations, lc.CheckAccelerationLimits(acceleration)...)
	}
	if !math.IsNaN(jerk) {
		allViolations = append(allViolations, lc.CheckJerkLimits(jerk)...)
	}
	
	return allViolations
}

// ValidateLimits validates that the limits themselves are reasonable
func (lc *LimitChecker) ValidateLimits() error {
	if lc.limits == nil {
		return fmt.Errorf("limits cannot be nil")
	}
	
	if lc.limits.MinPosition >= lc.limits.MaxPosition {
		return fmt.Errorf("minimum position %f must be less than maximum position %f", 
			lc.limits.MinPosition, lc.limits.MaxPosition)
	}
	
	if lc.limits.MaxVelocity <= 0 {
		return fmt.Errorf("maximum velocity %f must be positive", lc.limits.MaxVelocity)
	}
	
	if lc.limits.MinForce >= lc.limits.MaxForce {
		return fmt.Errorf("minimum force %f must be less than maximum force %f", 
			lc.limits.MinForce, lc.limits.MaxForce)
	}
	
	if lc.limits.MaxAcceleration <= 0 {
		return fmt.Errorf("maximum acceleration %f must be positive", lc.limits.MaxAcceleration)
	}
	
	if lc.limits.MaxJerk <= 0 {
		return fmt.Errorf("maximum jerk %f must be positive", lc.limits.MaxJerk)
	}
	
	// Check for reasonable ranges
	if math.Abs(lc.limits.MaxPosition) > 1e6 {
		return fmt.Errorf("maximum position %f is unreasonably large", lc.limits.MaxPosition)
	}
	
	if lc.limits.MaxVelocity > 1e6 {
		return fmt.Errorf("maximum velocity %f is unreasonably large", lc.limits.MaxVelocity)
	}
	
	return nil
}

// GetLimits returns the current limits
func (lc *LimitChecker) GetLimits() *SafetyLimits {
	return lc.limits
}

// SetLimits updates the limits
func (lc *LimitChecker) SetLimits(limits *SafetyLimits) error {
	if limits == nil {
		return fmt.Errorf("limits cannot be nil")
	}
	
	// Validate the new limits
	checker := NewLimitChecker(limits)
	if err := checker.ValidateLimits(); err != nil {
		return fmt.Errorf("invalid limits: %w", err)
	}
	
	lc.limits = limits
	return nil
}