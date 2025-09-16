package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// PositionUnit represents the unit for position values
type PositionUnit int

const (
	PositionUnitMM     PositionUnit = iota // millimeters
	PositionUnitCounts                     // encoder counts
)

// String returns the string representation of the position unit
func (pu PositionUnit) String() string {
	switch pu {
	case PositionUnitMM:
		return "mm"
	case PositionUnitCounts:
		return "counts"
	default:
		return "unknown"
	}
}

// VelocityUnit represents the unit for velocity values
type VelocityUnit int

const (
	VelocityUnitMMS     VelocityUnit = iota // millimeters per second
	VelocityUnitCountsS                     // counts per second
)

// String returns the string representation of the velocity unit
func (vu VelocityUnit) String() string {
	switch vu {
	case VelocityUnitMMS:
		return "mm/s"
	case VelocityUnitCountsS:
		return "counts/s"
	default:
		return "unknown"
	}
}

// AccelerationUnit represents the unit for acceleration values
type AccelerationUnit int

const (
	AccelerationUnitMMS2     AccelerationUnit = iota // millimeters per second squared
	AccelerationUnitCountsS2                         // counts per second squared
)

// String returns the string representation of the acceleration unit
func (au AccelerationUnit) String() string {
	switch au {
	case AccelerationUnitMMS2:
		return "mm/s²"
	case AccelerationUnitCountsS2:
		return "counts/s²"
	default:
		return "unknown"
	}
}

// JerkUnit represents the unit for jerk values
type JerkUnit int

const (
	JerkUnitMMS3     JerkUnit = iota // millimeters per second cubed
	JerkUnitCountsS3                 // counts per second cubed
)

// String returns the string representation of the jerk unit
func (ju JerkUnit) String() string {
	switch ju {
	case JerkUnitMMS3:
		return "mm/s³"
	case JerkUnitCountsS3:
		return "counts/s³"
	default:
		return "unknown"
	}
}

// ForceUnit represents the unit for force values
type ForceUnit int

const (
	ForceUnitN        ForceUnit = iota // Newtons
	ForceUnitCounts                    // force counts
)

// String returns the string representation of the force unit
func (fu ForceUnit) String() string {
	switch fu {
	case ForceUnitN:
		return "N"
	case ForceUnitCounts:
		return "counts"
	default:
		return "unknown"
	}
}

// TimeUnit represents the unit for time values
type TimeUnit int

const (
	TimeUnitMS TimeUnit = iota // milliseconds
	TimeUnitS                  // seconds
)

// String returns the string representation of the time unit
func (tu TimeUnit) String() string {
	switch tu {
	case TimeUnitMS:
		return "ms"
	case TimeUnitS:
		return "s"
	default:
		return "unknown"
	}
}

// PositionValue represents a position value with unit
type PositionValue struct {
	Value float64      `json:"value" yaml:"value"`
	Unit  PositionUnit `json:"unit" yaml:"unit"`
}

// NewPositionValue creates a new position value
func NewPositionValue(value float64, unit PositionUnit) *PositionValue {
	return &PositionValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the position value
func (pv *PositionValue) String() string {
	return fmt.Sprintf("%.3f %s", pv.Value, pv.Unit.String())
}

// VelocityValue represents a velocity value with unit
type VelocityValue struct {
	Value float64      `json:"value" yaml:"value"`
	Unit  VelocityUnit `json:"unit" yaml:"unit"`
}

// NewVelocityValue creates a new velocity value
func NewVelocityValue(value float64, unit VelocityUnit) *VelocityValue {
	return &VelocityValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the velocity value
func (vv *VelocityValue) String() string {
	return fmt.Sprintf("%.3f %s", vv.Value, vv.Unit.String())
}

// AccelerationValue represents an acceleration value with unit
type AccelerationValue struct {
	Value float64            `json:"value" yaml:"value"`
	Unit  AccelerationUnit   `json:"unit" yaml:"unit"`
}

// NewAccelerationValue creates a new acceleration value
func NewAccelerationValue(value float64, unit AccelerationUnit) *AccelerationValue {
	return &AccelerationValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the acceleration value
func (av *AccelerationValue) String() string {
	return fmt.Sprintf("%.3f %s", av.Value, av.Unit.String())
}

// JerkValue represents a jerk value with unit
type JerkValue struct {
	Value float64  `json:"value" yaml:"value"`
	Unit  JerkUnit `json:"unit" yaml:"unit"`
}

// NewJerkValue creates a new jerk value
func NewJerkValue(value float64, unit JerkUnit) *JerkValue {
	return &JerkValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the jerk value
func (jv *JerkValue) String() string {
	return fmt.Sprintf("%.3f %s", jv.Value, jv.Unit.String())
}

// ForceValue represents a force value with unit
type ForceValue struct {
	Value float64   `json:"value" yaml:"value"`
	Unit  ForceUnit `json:"unit" yaml:"unit"`
}

// NewForceValue creates a new force value
func NewForceValue(value float64, unit ForceUnit) *ForceValue {
	return &ForceValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the force value
func (fv *ForceValue) String() string {
	return fmt.Sprintf("%.3f %s", fv.Value, fv.Unit.String())
}

// TimeValue represents a time value with unit
type TimeValue struct {
	Value float64  `json:"value" yaml:"value"`
	Unit  TimeUnit `json:"unit" yaml:"unit"`
}

// NewTimeValue creates a new time value
func NewTimeValue(value float64, unit TimeUnit) *TimeValue {
	return &TimeValue{
		Value: value,
		Unit:  unit,
	}
}

// String returns the string representation of the time value
func (tv *TimeValue) String() string {
	return fmt.Sprintf("%.3f %s", tv.Value, tv.Unit.String())
}

// Duration returns the time value as a Go duration
func (tv *TimeValue) Duration() time.Duration {
	switch tv.Unit {
	case TimeUnitMS:
		return time.Duration(tv.Value) * time.Millisecond
	case TimeUnitS:
		return time.Duration(tv.Value) * time.Second
	default:
		return 0
	}
}

// MotionParameters represents parameters for motion commands
type MotionParameters struct {
	Position     *PositionValue     `json:"position,omitempty" yaml:"position,omitempty"`
	Velocity     *VelocityValue     `json:"velocity,omitempty" yaml:"velocity,omitempty"`
	Acceleration *AccelerationValue `json:"acceleration,omitempty" yaml:"acceleration,omitempty"`
	Deceleration *AccelerationValue `json:"deceleration,omitempty" yaml:"deceleration,omitempty"`
	Jerk         *JerkValue         `json:"jerk,omitempty" yaml:"jerk,omitempty"`
	Timeout      *TimeValue         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Tolerance    *PositionValue     `json:"tolerance,omitempty" yaml:"tolerance,omitempty"`
}

// IOParameters represents parameters for I/O commands
type IOParameters struct {
	OutputNumber int       `json:"output_number,omitempty" yaml:"output_number,omitempty"`
	InputNumber  int       `json:"input_number,omitempty" yaml:"input_number,omitempty"`
	State        bool      `json:"state,omitempty" yaml:"state,omitempty"`
	Value        float64   `json:"value,omitempty" yaml:"value,omitempty"`
	Tolerance    float64   `json:"tolerance,omitempty" yaml:"tolerance,omitempty"`
	Timeout      *TimeValue `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// LoopParameters represents parameters for loop commands
type LoopParameters struct {
	CounterVariable string    `json:"counter_variable" yaml:"counter_variable"`
	MaxIterations   int       `json:"max_iterations" yaml:"max_iterations"`
	Condition       Condition `json:"condition,omitempty" yaml:"condition,omitempty"`
	NextCommand     int       `json:"next_command" yaml:"next_command"`
}

// JumpParameters represents parameters for jump commands
type JumpParameters struct {
	Condition      Condition `json:"condition,omitempty" yaml:"condition,omitempty"`
	TargetCommand  int       `json:"target_command" yaml:"target_command"`
	FalseCommand   int       `json:"false_command,omitempty" yaml:"false_command,omitempty"`
}

// SystemParameters represents parameters for system commands
type SystemParameters struct {
	HomingMethod   int         `json:"homing_method,omitempty" yaml:"homing_method,omitempty"`
	ResetType      int         `json:"reset_type,omitempty" yaml:"reset_type,omitempty"`
	ConfigSlot     int         `json:"config_slot,omitempty" yaml:"config_slot,omitempty"`
	Velocity       *VelocityValue     `json:"velocity,omitempty" yaml:"velocity,omitempty"`
	Acceleration   *AccelerationValue `json:"acceleration,omitempty" yaml:"acceleration,omitempty"`
	Timeout        *TimeValue         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// ForceControlParameters represents parameters for force control commands
type ForceControlParameters struct {
	ForceSetpoint  *ForceValue `json:"force_setpoint,omitempty" yaml:"force_setpoint,omitempty"`
	ForceLimit     *ForceValue `json:"force_limit,omitempty" yaml:"force_limit,omitempty"`
	PositionLimit  *PositionValue `json:"position_limit,omitempty" yaml:"position_limit,omitempty"`
	TransitionTime *TimeValue     `json:"transition_time,omitempty" yaml:"transition_time,omitempty"`
}

// DataAcquisitionParameters represents parameters for data acquisition commands
type DataAcquisitionParameters struct {
	SampleRate     int      `json:"sample_rate,omitempty" yaml:"sample_rate,omitempty"`
	NumSamples     int      `json:"num_samples,omitempty" yaml:"num_samples,omitempty"`
	Channels       []int    `json:"channels,omitempty" yaml:"channels,omitempty"`
	SaveData       bool     `json:"save_data,omitempty" yaml:"save_data,omitempty"`
	Filename       string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	Format         string   `json:"format,omitempty" yaml:"format,omitempty"`
}

// ParameterExtractor provides utilities for extracting parameters from command maps
type ParameterExtractor struct{}

// NewParameterExtractor creates a new parameter extractor
func NewParameterExtractor() *ParameterExtractor {
	return &ParameterExtractor{}
}

// ExtractPosition extracts a position value from parameters
func (pe *ParameterExtractor) ExtractPosition(params map[string]interface{}, key string) (*PositionValue, error) {
	value, ok := params[key]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", key)
	}
	
	posMap, ok := value.(map[string]interface{})
	if !ok {
		// Try to extract as simple float
		if floatVal, ok := value.(float64); ok {
			return NewPositionValue(floatVal, PositionUnitMM), nil
		}
		return nil, fmt.Errorf("parameter %s is not a position value", key)
	}
	
	val, ok := posMap["value"].(float64)
	if !ok {
		return nil, fmt.Errorf("position value is not a number")
	}
	
	unitStr, ok := posMap["unit"].(string)
	if !ok {
		unitStr = "mm" // default
	}
	
	var unit PositionUnit
	switch unitStr {
	case "mm":
		unit = PositionUnitMM
	case "counts":
		unit = PositionUnitCounts
	default:
		return nil, fmt.Errorf("unknown position unit: %s", unitStr)
	}
	
	return NewPositionValue(val, unit), nil
}

// ExtractVelocity extracts a velocity value from parameters
func (pe *ParameterExtractor) ExtractVelocity(params map[string]interface{}, key string) (*VelocityValue, error) {
	value, ok := params[key]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", key)
	}
	
	velMap, ok := value.(map[string]interface{})
	if !ok {
		// Try to extract as simple float
		if floatVal, ok := value.(float64); ok {
			return NewVelocityValue(floatVal, VelocityUnitMMS), nil
		}
		return nil, fmt.Errorf("parameter %s is not a velocity value", key)
	}
	
	val, ok := velMap["value"].(float64)
	if !ok {
		return nil, fmt.Errorf("velocity value is not a number")
	}
	
	unitStr, ok := velMap["unit"].(string)
	if !ok {
		unitStr = "mm/s" // default
	}
	
	var unit VelocityUnit
	switch unitStr {
	case "mm/s":
		unit = VelocityUnitMMS
	case "counts/s":
		unit = VelocityUnitCountsS
	default:
		return nil, fmt.Errorf("unknown velocity unit: %s", unitStr)
	}
	
	return NewVelocityValue(val, unit), nil
}

// ExtractTime extracts a time value from parameters
func (pe *ParameterExtractor) ExtractTime(params map[string]interface{}, key string) (*TimeValue, error) {
	value, ok := params[key]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", key)
	}
	
	timeMap, ok := value.(map[string]interface{})
	if !ok {
		// Try to extract as simple float
		if floatVal, ok := value.(float64); ok {
			return NewTimeValue(floatVal, TimeUnitMS), nil
		}
		return nil, fmt.Errorf("parameter %s is not a time value", key)
	}
	
	val, ok := timeMap["value"].(float64)
	if !ok {
		return nil, fmt.Errorf("time value is not a number")
	}
	
	unitStr, ok := timeMap["unit"].(string)
	if !ok {
		unitStr = "ms" // default
	}
	
	var unit TimeUnit
	switch unitStr {
	case "ms":
		unit = TimeUnitMS
	case "s":
		unit = TimeUnitS
	default:
		return nil, fmt.Errorf("unknown time unit: %s", unitStr)
	}
	
	return NewTimeValue(val, unit), nil
}

// ExtractInt extracts an integer value from parameters
func (pe *ParameterExtractor) ExtractInt(params map[string]interface{}, key string) (int, error) {
	value, ok := params[key]
	if !ok {
		return 0, fmt.Errorf("parameter %s not found", key)
	}
	
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("parameter %s is not an integer", key)
	}
}

// ExtractFloat extracts a float value from parameters
func (pe *ParameterExtractor) ExtractFloat(params map[string]interface{}, key string) (float64, error) {
	value, ok := params[key]
	if !ok {
		return 0, fmt.Errorf("parameter %s not found", key)
	}
	
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("parameter %s is not a number", key)
	}
}

// ExtractBool extracts a boolean value from parameters
func (pe *ParameterExtractor) ExtractBool(params map[string]interface{}, key string) (bool, error) {
	value, ok := params[key]
	if !ok {
		return false, fmt.Errorf("parameter %s not found", key)
	}
	
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		return false, fmt.Errorf("parameter %s is not a boolean", key)
	}
}

// ExtractString extracts a string value from parameters
func (pe *ParameterExtractor) ExtractString(params map[string]interface{}, key string) (string, error) {
	value, ok := params[key]
	if !ok {
		return "", fmt.Errorf("parameter %s not found", key)
	}
	
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("parameter %s is not a string", key)
	}
}

// MarshalJSON implements json.Marshaler for PositionValue
func (pv *PositionValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"value": pv.Value,
		"unit":  pv.Unit.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler for PositionValue
func (pv *PositionValue) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	
	if value, ok := m["value"].(float64); ok {
		pv.Value = value
	}
	
	if unitStr, ok := m["unit"].(string); ok {
		switch unitStr {
		case "mm":
			pv.Unit = PositionUnitMM
		case "counts":
			pv.Unit = PositionUnitCounts
		default:
			return fmt.Errorf("unknown position unit: %s", unitStr)
		}
	}
	
	return nil
}