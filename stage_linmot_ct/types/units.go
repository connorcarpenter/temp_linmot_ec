package types

import (
	"fmt"
	"math"
)

// UnitConverter provides unit conversion utilities
type UnitConverter struct {
	positionScalingFactor float64 // counts per mm
	forceScalingFactor    float64 // counts per Newton
}

// NewUnitConverter creates a new unit converter with default scaling factors
func NewUnitConverter() *UnitConverter {
	return &UnitConverter{
		positionScalingFactor: 1000.0, // Default: 1000 counts per mm
		forceScalingFactor:    100.0,  // Default: 100 counts per Newton
	}
}

// NewUnitConverterWithFactors creates a new unit converter with custom scaling factors
func NewUnitConverterWithFactors(posFactor, forceFactor float64) *UnitConverter {
	return &UnitConverter{
		positionScalingFactor: posFactor,
		forceScalingFactor:    forceFactor,
	}
}

// ConvertPosition converts a position value between units
func (uc *UnitConverter) ConvertPosition(value float64, from, to PositionUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == PositionUnitMM && to == PositionUnitCounts:
		return value * uc.positionScalingFactor
	case from == PositionUnitCounts && to == PositionUnitMM:
		return value / uc.positionScalingFactor
	default:
		return value
	}
}

// ConvertVelocity converts a velocity value between units
func (uc *UnitConverter) ConvertVelocity(value float64, from, to VelocityUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == VelocityUnitMMS && to == VelocityUnitCountsS:
		return value * uc.positionScalingFactor
	case from == VelocityUnitCountsS && to == VelocityUnitMMS:
		return value / uc.positionScalingFactor
	default:
		return value
	}
}

// ConvertAcceleration converts an acceleration value between units
func (uc *UnitConverter) ConvertAcceleration(value float64, from, to AccelerationUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == AccelerationUnitMMS2 && to == AccelerationUnitCountsS2:
		return value * uc.positionScalingFactor
	case from == AccelerationUnitCountsS2 && to == AccelerationUnitMMS2:
		return value / uc.positionScalingFactor
	default:
		return value
	}
}

// ConvertJerk converts a jerk value between units
func (uc *UnitConverter) ConvertJerk(value float64, from, to JerkUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == JerkUnitMMS3 && to == JerkUnitCountsS3:
		return value * uc.positionScalingFactor
	case from == JerkUnitCountsS3 && to == JerkUnitMMS3:
		return value / uc.positionScalingFactor
	default:
		return value
	}
}

// ConvertForce converts a force value between units
func (uc *UnitConverter) ConvertForce(value float64, from, to ForceUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == ForceUnitN && to == ForceUnitCounts:
		return value * uc.forceScalingFactor
	case from == ForceUnitCounts && to == ForceUnitN:
		return value / uc.forceScalingFactor
	default:
		return value
	}
}

// ConvertTime converts a time value between units
func (uc *UnitConverter) ConvertTime(value float64, from, to TimeUnit) float64 {
	if from == to {
		return value
	}
	
	switch {
	case from == TimeUnitMS && to == TimeUnitS:
		return value / 1000.0
	case from == TimeUnitS && to == TimeUnitMS:
		return value * 1000.0
	default:
		return value
	}
}

// ConvertPositionValue converts a PositionValue to different units
func (uc *UnitConverter) ConvertPositionValue(pv *PositionValue, to PositionUnit) *PositionValue {
	if pv.Unit == to {
		return pv
	}
	
	convertedValue := uc.ConvertPosition(pv.Value, pv.Unit, to)
	return &PositionValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// ConvertVelocityValue converts a VelocityValue to different units
func (uc *UnitConverter) ConvertVelocityValue(vv *VelocityValue, to VelocityUnit) *VelocityValue {
	if vv.Unit == to {
		return vv
	}
	
	convertedValue := uc.ConvertVelocity(vv.Value, vv.Unit, to)
	return &VelocityValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// ConvertAccelerationValue converts an AccelerationValue to different units
func (uc *UnitConverter) ConvertAccelerationValue(av *AccelerationValue, to AccelerationUnit) *AccelerationValue {
	if av.Unit == to {
		return av
	}
	
	convertedValue := uc.ConvertAcceleration(av.Value, av.Unit, to)
	return &AccelerationValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// ConvertJerkValue converts a JerkValue to different units
func (uc *UnitConverter) ConvertJerkValue(jv *JerkValue, to JerkUnit) *JerkValue {
	if jv.Unit == to {
		return jv
	}
	
	convertedValue := uc.ConvertJerk(jv.Value, jv.Unit, to)
	return &JerkValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// ConvertForceValue converts a ForceValue to different units
func (uc *UnitConverter) ConvertForceValue(fv *ForceValue, to ForceUnit) *ForceValue {
	if fv.Unit == to {
		return fv
	}
	
	convertedValue := uc.ConvertForce(fv.Value, fv.Unit, to)
	return &ForceValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// ConvertTimeValue converts a TimeValue to different units
func (uc *UnitConverter) ConvertTimeValue(tv *TimeValue, to TimeUnit) *TimeValue {
	if tv.Unit == to {
		return tv
	}
	
	convertedValue := uc.ConvertTime(tv.Value, tv.Unit, to)
	return &TimeValue{
		Value: convertedValue,
		Unit:  to,
	}
}

// SetPositionScalingFactor sets the position scaling factor
func (uc *UnitConverter) SetPositionScalingFactor(factor float64) error {
	if factor <= 0 {
		return fmt.Errorf("position scaling factor must be positive, got %f", factor)
	}
	uc.positionScalingFactor = factor
	return nil
}

// SetForceScalingFactor sets the force scaling factor
func (uc *UnitConverter) SetForceScalingFactor(factor float64) error {
	if factor <= 0 {
		return fmt.Errorf("force scaling factor must be positive, got %f", factor)
	}
	uc.forceScalingFactor = factor
	return nil
}

// GetPositionScalingFactor returns the position scaling factor
func (uc *UnitConverter) GetPositionScalingFactor() float64 {
	return uc.positionScalingFactor
}

// GetForceScalingFactor returns the force scaling factor
func (uc *UnitConverter) GetForceScalingFactor() float64 {
	return uc.forceScalingFactor
}

// ValidateScalingFactors validates that the scaling factors are reasonable
func (uc *UnitConverter) ValidateScalingFactors() error {
	if uc.positionScalingFactor <= 0 {
		return fmt.Errorf("position scaling factor must be positive, got %f", uc.positionScalingFactor)
	}
	if uc.forceScalingFactor <= 0 {
		return fmt.Errorf("force scaling factor must be positive, got %f", uc.forceScalingFactor)
	}
	if math.IsInf(uc.positionScalingFactor, 0) || math.IsNaN(uc.positionScalingFactor) {
		return fmt.Errorf("position scaling factor is invalid: %f", uc.positionScalingFactor)
	}
	if math.IsInf(uc.forceScalingFactor, 0) || math.IsNaN(uc.forceScalingFactor) {
		return fmt.Errorf("force scaling factor is invalid: %f", uc.forceScalingFactor)
	}
	return nil
}

// UnitConversionError represents an error in unit conversion
type UnitConversionError struct {
	FromUnit string
	ToUnit   string
	Value    float64
	Message  string
}

// Error returns the error message
func (uce *UnitConversionError) Error() string {
	return fmt.Sprintf("unit conversion error from %s to %s for value %f: %s", 
		uce.FromUnit, uce.ToUnit, uce.Value, uce.Message)
}

// NewUnitConversionError creates a new unit conversion error
func NewUnitConversionError(fromUnit, toUnit string, value float64, message string) *UnitConversionError {
	return &UnitConversionError{
		FromUnit: fromUnit,
		ToUnit:   toUnit,
		Value:    value,
		Message:  message,
	}
}

// UnitSystem represents a complete unit system configuration
type UnitSystem struct {
	PositionUnit     PositionUnit
	VelocityUnit     VelocityUnit
	AccelerationUnit AccelerationUnit
	JerkUnit         JerkUnit
	ForceUnit        ForceUnit
	TimeUnit         TimeUnit
}

// DefaultUnitSystem returns the default unit system (metric)
func DefaultUnitSystem() *UnitSystem {
	return &UnitSystem{
		PositionUnit:     PositionUnitMM,
		VelocityUnit:     VelocityUnitMMS,
		AccelerationUnit: AccelerationUnitMMS2,
		JerkUnit:         JerkUnitMMS3,
		ForceUnit:        ForceUnitN,
		TimeUnit:         TimeUnitMS,
	}
}

// ImperialUnitSystem returns an imperial unit system
func ImperialUnitSystem() *UnitSystem {
	// Note: This would need additional unit types for imperial units
	// For now, we'll use the default system
	return DefaultUnitSystem()
}

// ConvertToUnitSystem converts all values in a command to the specified unit system
func (uc *UnitConverter) ConvertToUnitSystem(cmd *Command, targetSystem *UnitSystem) error {
	// This is a placeholder for converting command parameters to a target unit system
	// The actual implementation would depend on the specific command type and parameters
	
	switch cmd.Type {
	case CmdMoveAbsolute, CmdMoveRelative, CmdMoveIncremental:
		return uc.convertMotionCommandToUnitSystem(cmd, targetSystem)
	case CmdJog:
		return uc.convertJogCommandToUnitSystem(cmd, targetSystem)
	case CmdWaitPosition, CmdWaitVelocity, CmdWaitForce:
		return uc.convertWaitCommandToUnitSystem(cmd, targetSystem)
	default:
		// No conversion needed for other command types
		return nil
	}
}

// convertMotionCommandToUnitSystem converts motion command parameters to target unit system
func (uc *UnitConverter) convertMotionCommandToUnitSystem(cmd *Command, targetSystem *UnitSystem) error {
	// Extract and convert position parameters
	if pos, ok := cmd.Parameters["position"]; ok {
		if posValue, ok := pos.(*PositionValue); ok {
			converted := uc.ConvertPositionValue(posValue, targetSystem.PositionUnit)
			cmd.Parameters["position"] = converted
		}
	}
	
	// Extract and convert velocity parameters
	if vel, ok := cmd.Parameters["velocity"]; ok {
		if velValue, ok := vel.(*VelocityValue); ok {
			converted := uc.ConvertVelocityValue(velValue, targetSystem.VelocityUnit)
			cmd.Parameters["velocity"] = converted
		}
	}
	
	// Extract and convert acceleration parameters
	if accel, ok := cmd.Parameters["acceleration"]; ok {
		if accelValue, ok := accel.(*AccelerationValue); ok {
			converted := uc.ConvertAccelerationValue(accelValue, targetSystem.AccelerationUnit)
			cmd.Parameters["acceleration"] = converted
		}
	}
	
	// Extract and convert jerk parameters
	if jerk, ok := cmd.Parameters["jerk"]; ok {
		if jerkValue, ok := jerk.(*JerkValue); ok {
			converted := uc.ConvertJerkValue(jerkValue, targetSystem.JerkUnit)
			cmd.Parameters["jerk"] = converted
		}
	}
	
	return nil
}

// convertJogCommandToUnitSystem converts jog command parameters to target unit system
func (uc *UnitConverter) convertJogCommandToUnitSystem(cmd *Command, targetSystem *UnitSystem) error {
	// Extract and convert velocity parameters
	if vel, ok := cmd.Parameters["velocity"]; ok {
		if velValue, ok := vel.(*VelocityValue); ok {
			converted := uc.ConvertVelocityValue(velValue, targetSystem.VelocityUnit)
			cmd.Parameters["velocity"] = converted
		}
	}
	
	// Extract and convert acceleration parameters
	if accel, ok := cmd.Parameters["acceleration"]; ok {
		if accelValue, ok := accel.(*AccelerationValue); ok {
			converted := uc.ConvertAccelerationValue(accelValue, targetSystem.AccelerationUnit)
			cmd.Parameters["acceleration"] = converted
		}
	}
	
	return nil
}

// convertWaitCommandToUnitSystem converts wait command parameters to target unit system
func (uc *UnitConverter) convertWaitCommandToUnitSystem(cmd *Command, targetSystem *UnitSystem) error {
	// Extract and convert position parameters for wait position
	if pos, ok := cmd.Parameters["position"]; ok {
		if posValue, ok := pos.(*PositionValue); ok {
			converted := uc.ConvertPositionValue(posValue, targetSystem.PositionUnit)
			cmd.Parameters["position"] = converted
		}
	}
	
	// Extract and convert velocity parameters for wait velocity
	if vel, ok := cmd.Parameters["velocity"]; ok {
		if velValue, ok := vel.(*VelocityValue); ok {
			converted := uc.ConvertVelocityValue(velValue, targetSystem.VelocityUnit)
			cmd.Parameters["velocity"] = converted
		}
	}
	
	// Extract and convert force parameters for wait force
	if force, ok := cmd.Parameters["force"]; ok {
		if forceValue, ok := force.(*ForceValue); ok {
			converted := uc.ConvertForceValue(forceValue, targetSystem.ForceUnit)
			cmd.Parameters["force"] = converted
		}
	}
	
	// Extract and convert timeout parameters
	if timeout, ok := cmd.Parameters["timeout"]; ok {
		if timeoutValue, ok := timeout.(*TimeValue); ok {
			converted := uc.ConvertTimeValue(timeoutValue, targetSystem.TimeUnit)
			cmd.Parameters["timeout"] = converted
		}
	}
	
	return nil
}

// UnitConverterFactory creates unit converters for different drive types
type UnitConverterFactory struct {
	driveTypeScalingFactors map[string]map[string]float64
}

// NewUnitConverterFactory creates a new unit converter factory
func NewUnitConverterFactory() *UnitConverterFactory {
	return &UnitConverterFactory{
		driveTypeScalingFactors: make(map[string]map[string]float64),
	}
}

// RegisterDriveType registers scaling factors for a specific drive type
func (ucf *UnitConverterFactory) RegisterDriveType(driveType string, factors map[string]float64) {
	ucf.driveTypeScalingFactors[driveType] = factors
}

// CreateUnitConverter creates a unit converter for a specific drive type
func (ucf *UnitConverterFactory) CreateUnitConverter(driveType string) (*UnitConverter, error) {
	factors, ok := ucf.driveTypeScalingFactors[driveType]
	if !ok {
		return NewUnitConverter(), nil // Use default factors
	}
	
	posFactor, ok := factors["position"]
	if !ok {
		posFactor = 1000.0 // Default
	}
	
	forceFactor, ok := factors["force"]
	if !ok {
		forceFactor = 100.0 // Default
	}
	
	return NewUnitConverterWithFactors(posFactor, forceFactor), nil
}

// GetSupportedDriveTypes returns the list of supported drive types
func (ucf *UnitConverterFactory) GetSupportedDriveTypes() []string {
	types := make([]string, 0, len(ucf.driveTypeScalingFactors))
	for driveType := range ucf.driveTypeScalingFactors {
		types = append(types, driveType)
	}
	return types
}