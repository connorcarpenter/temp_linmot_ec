package safety

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// SafetyGuard provides safety validation for motion commands
type SafetyGuard struct {
	limits             *SafetyLimits
	emergencyStopActive bool
}

// SafetyLimits defines the safety limits for the system
type SafetyLimits struct {
	// Position limits (in counts)
	MinPosition float64
	MaxPosition float64
	
	// Velocity limits (in counts/s)
	MaxVelocity float64
	
	// Force limits (in counts)
	MinForce float64
	MaxForce float64
	
	// Acceleration limits (in counts/s²)
	MaxAcceleration float64
	
	// Jerk limits (in counts/s³)
	MaxJerk float64
}

// NewSafetyGuard creates a new safety guard with default limits
func NewSafetyGuard() *SafetyGuard {
	return &SafetyGuard{
		limits: &SafetyLimits{
			MinPosition:    -1000000.0, // -1000mm in counts (assuming 1000 counts/mm)
			MaxPosition:    1000000.0,  // 1000mm in counts
			MaxVelocity:    10000.0,    // 10mm/s in counts/s
			MinForce:       -1000.0,    // -10N in counts (assuming 100 counts/N)
			MaxForce:       1000.0,     // 10N in counts
			MaxAcceleration: 1000.0,    // 1mm/s² in counts/s²
			MaxJerk:        10000.0,    // 10mm/s³ in counts/s³
		},
	}
}

// NewSafetyGuardWithLimits creates a new safety guard with custom limits
func NewSafetyGuardWithLimits(limits *SafetyLimits) *SafetyGuard {
	return &SafetyGuard{
		limits: limits,
	}
}

// ValidatePosition checks if a position is within safety limits
func (sg *SafetyGuard) ValidatePosition(position float64) error {
	if position < sg.limits.MinPosition {
		return fmt.Errorf("position %f is below minimum limit %f", position, sg.limits.MinPosition)
	}
	if position > sg.limits.MaxPosition {
		return fmt.Errorf("position %f is above maximum limit %f", position, sg.limits.MaxPosition)
	}
	return nil
}

// ValidateVelocity checks if a velocity is within safety limits
func (sg *SafetyGuard) ValidateVelocity(velocity float64) error {
	if velocity < 0 {
		return fmt.Errorf("velocity %f is negative", velocity)
	}
	if velocity > sg.limits.MaxVelocity {
		return fmt.Errorf("velocity %f exceeds maximum limit %f", velocity, sg.limits.MaxVelocity)
	}
	return nil
}

// ValidateForce checks if a force is within safety limits
func (sg *SafetyGuard) ValidateForce(force float64) error {
	if force < sg.limits.MinForce {
		return fmt.Errorf("force %f is below minimum limit %f", force, sg.limits.MinForce)
	}
	if force > sg.limits.MaxForce {
		return fmt.Errorf("force %f exceeds maximum limit %f", force, sg.limits.MaxForce)
	}
	return nil
}

// ValidateAcceleration checks if an acceleration is within safety limits
func (sg *SafetyGuard) ValidateAcceleration(acceleration float64) error {
	if acceleration < 0 {
		return fmt.Errorf("acceleration %f is negative", acceleration)
	}
	if acceleration > sg.limits.MaxAcceleration {
		return fmt.Errorf("acceleration %f exceeds maximum limit %f", acceleration, sg.limits.MaxAcceleration)
	}
	return nil
}

// ValidateJerk checks if a jerk is within safety limits
func (sg *SafetyGuard) ValidateJerk(jerk float64) error {
	if jerk < 0 {
		return fmt.Errorf("jerk %f is negative", jerk)
	}
	if jerk > sg.limits.MaxJerk {
		return fmt.Errorf("jerk %f exceeds maximum limit %f", jerk, sg.limits.MaxJerk)
	}
	return nil
}

// ValidateMotionCommand validates all parameters of a motion command
func (sg *SafetyGuard) ValidateMotionCommand(command *types.Command, unitConverter *types.UnitConverter) error {
	pe := types.NewParameterExtractor()
	
	// Extract and validate position (if present)
	if position, err := pe.ExtractPosition(command.Parameters, "position"); err == nil {
		positionInCounts := unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
		if err := sg.ValidatePosition(positionInCounts.Value); err != nil {
			return fmt.Errorf("position validation failed: %w", err)
		}
	}
	
	// Extract and validate velocity (if present)
	if velocity, err := pe.ExtractVelocity(command.Parameters, "velocity"); err == nil {
		velocityInCounts := unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
		if err := sg.ValidateVelocity(velocityInCounts.Value); err != nil {
			return fmt.Errorf("velocity validation failed: %w", err)
		}
	}
	
	// Extract and validate acceleration (if present)
	if acceleration, err := pe.ExtractAcceleration(command.Parameters, "acceleration"); err == nil {
		accelerationInCounts := unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
		if err := sg.ValidateAcceleration(accelerationInCounts.Value); err != nil {
			return fmt.Errorf("acceleration validation failed: %w", err)
		}
	}
	
	// Extract and validate jerk (if present)
	if jerk, err := pe.ExtractJerk(command.Parameters, "jerk"); err == nil {
		jerkInCounts := unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
		if err := sg.ValidateJerk(jerkInCounts.Value); err != nil {
			return fmt.Errorf("jerk validation failed: %w", err)
		}
	}
	
	return nil
}

// ValidateForceCommand validates all parameters of a force command
func (sg *SafetyGuard) ValidateForceCommand(command *types.Command, unitConverter *types.UnitConverter) error {
	pe := types.NewParameterExtractor()
	
	// Extract and validate force (if present)
	if force, err := pe.ExtractForce(command.Parameters, "force"); err == nil {
		forceInCounts := unitConverter.ConvertForceValue(force, types.ForceUnitCounts)
		if err := sg.ValidateForce(forceInCounts.Value); err != nil {
			return fmt.Errorf("force validation failed: %w", err)
		}
	}
	
	return nil
}

// GetLimits returns the current safety limits
func (sg *SafetyGuard) GetLimits() *SafetyLimits {
	return sg.limits
}

// SetLimits updates the safety limits
func (sg *SafetyGuard) SetLimits(limits *SafetyLimits) error {
	if limits == nil {
		return fmt.Errorf("limits cannot be nil")
	}
	
	// Validate the limits
	if limits.MinPosition >= limits.MaxPosition {
		return fmt.Errorf("minimum position must be less than maximum position")
	}
	if limits.MaxVelocity <= 0 {
		return fmt.Errorf("maximum velocity must be positive")
	}
	if limits.MinForce >= limits.MaxForce {
		return fmt.Errorf("minimum force must be less than maximum force")
	}
	if limits.MaxAcceleration <= 0 {
		return fmt.Errorf("maximum acceleration must be positive")
	}
	if limits.MaxJerk <= 0 {
		return fmt.Errorf("maximum jerk must be positive")
	}
	
	sg.limits = limits
	return nil
}

// EmergencyStop represents an emergency stop condition
type EmergencyStop struct {
	Reason string
	Time   int64 // Unix timestamp
}

// IsEmergencyStopActive checks if an emergency stop is currently active
func (sg *SafetyGuard) IsEmergencyStopActive() bool {
	return sg.emergencyStopActive
}

// TriggerEmergencyStop triggers an emergency stop
func (sg *SafetyGuard) TriggerEmergencyStop(ctx context.Context, driveController types.DriveController, reason string) error {
	// Stop all motion immediately
	if err := driveController.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop motion during emergency stop: %w", err)
	}
	
	// Log the emergency stop
	// In a real implementation, this would log to a persistent store
	fmt.Printf("EMERGENCY STOP TRIGGERED: %s\n", reason)
	
	return nil
}