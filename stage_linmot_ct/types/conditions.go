package types

import (
	"context"
	"fmt"
	"time"
)

// ConditionType represents the type of condition
type ConditionType int

const (
	CondDigitalInput    ConditionType = iota // Digital input state
	CondAnalogInput                          // Analog input value
	CondPosition                             // Position condition
	CondVelocity                             // Velocity condition
	CondForce                                // Force condition
	CondTimer                                // Timer condition
	CondVariable                             // Variable condition
	CondError                                // Error condition
	CondDriveState                           // Drive state condition
	CondMotionComplete                       // Motion completion condition
)

// String returns the string representation of the condition type
func (ct ConditionType) String() string {
	switch ct {
	case CondDigitalInput:
		return "digital_input"
	case CondAnalogInput:
		return "analog_input"
	case CondPosition:
		return "position"
	case CondVelocity:
		return "velocity"
	case CondForce:
		return "force"
	case CondTimer:
		return "timer"
	case CondVariable:
		return "variable"
	case CondError:
		return "error"
	case CondDriveState:
		return "drive_state"
	case CondMotionComplete:
		return "motion_complete"
	default:
		return "unknown"
	}
}

// ComparisonOperator represents the comparison operator for conditions
type ComparisonOperator int

const (
	OpEqual              ComparisonOperator = iota // ==
	OpNotEqual                                     // !=
	OpGreaterThan                                  // >
	OpLessThan                                     // <
	OpGreaterThanOrEqual                           // >=
	OpLessThanOrEqual                              // <=
	OpAnd                                          // &&
	OpOr                                           // ||
	OpNot                                          // !
	OpIn                                           // in
	OpNotIn                                        // not in
)

// String returns the string representation of the comparison operator
func (co ComparisonOperator) String() string {
	switch co {
	case OpEqual:
		return "=="
	case OpNotEqual:
		return "!="
	case OpGreaterThan:
		return ">"
	case OpLessThan:
		return "<"
	case OpGreaterThanOrEqual:
		return ">="
	case OpLessThanOrEqual:
		return "<="
	case OpAnd:
		return "&&"
	case OpOr:
		return "||"
	case OpNot:
		return "!"
	case OpIn:
		return "in"
	case OpNotIn:
		return "not in"
	default:
		return "unknown"
	}
}

// Condition represents a condition that can be evaluated
type Condition struct {
	Type        ConditionType      `json:"type" yaml:"type"`
	Parameter   string             `json:"parameter" yaml:"parameter"`
	Operator    ComparisonOperator `json:"operator" yaml:"operator"`
	Value       interface{}        `json:"value" yaml:"value"`
	Timeout     *TimeValue         `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
}

// ConditionEvaluator defines the interface for evaluating conditions
type ConditionEvaluator interface {
	Evaluate(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error)
	CanEvaluate(condition *Condition) bool
	GetRequiredData(condition *Condition) []string
}

// ConditionBuilder provides a fluent interface for building conditions
type ConditionBuilder struct {
	condition *Condition
}

// NewConditionBuilder creates a new condition builder
func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{
		condition: &Condition{},
	}
}

// WithType sets the condition type
func (cb *ConditionBuilder) WithType(condType ConditionType) *ConditionBuilder {
	cb.condition.Type = condType
	return cb
}

// WithParameter sets the parameter name
func (cb *ConditionBuilder) WithParameter(param string) *ConditionBuilder {
	cb.condition.Parameter = param
	return cb
}

// WithOperator sets the comparison operator
func (cb *ConditionBuilder) WithOperator(op ComparisonOperator) *ConditionBuilder {
	cb.condition.Operator = op
	return cb
}

// WithValue sets the comparison value
func (cb *ConditionBuilder) WithValue(value interface{}) *ConditionBuilder {
	cb.condition.Value = value
	return cb
}

// WithTimeout sets the timeout for the condition
func (cb *ConditionBuilder) WithTimeout(timeout *TimeValue) *ConditionBuilder {
	cb.condition.Timeout = timeout
	return cb
}

// WithDescription sets the condition description
func (cb *ConditionBuilder) WithDescription(desc string) *ConditionBuilder {
	cb.condition.Description = desc
	return cb
}

// Build returns the constructed condition
func (cb *ConditionBuilder) Build() *Condition {
	return cb.condition
}

// DigitalInputCondition creates a digital input condition
func DigitalInputCondition(inputNum int, expectedState bool, timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondDigitalInput,
		Parameter: fmt.Sprintf("digital_input_%d", inputNum),
		Operator:  OpEqual,
		Value:     expectedState,
		Timeout:   timeout,
	}
}

// AnalogInputCondition creates an analog input condition
func AnalogInputCondition(inputNum int, operator ComparisonOperator, value float64, tolerance float64, timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondAnalogInput,
		Parameter: fmt.Sprintf("analog_input_%d", inputNum),
		Operator:  operator,
		Value:     value,
		Timeout:   timeout,
	}
}

// PositionCondition creates a position condition
func PositionCondition(operator ComparisonOperator, position *PositionValue, timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondPosition,
		Parameter: "position",
		Operator:  operator,
		Value:     position,
		Timeout:   timeout,
	}
}

// VelocityCondition creates a velocity condition
func VelocityCondition(operator ComparisonOperator, velocity *VelocityValue, timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondVelocity,
		Parameter: "velocity",
		Operator:  operator,
		Value:     velocity,
		Timeout:   timeout,
	}
}

// ForceCondition creates a force condition
func ForceCondition(operator ComparisonOperator, force *ForceValue, timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondForce,
		Parameter: "force",
		Operator:  operator,
		Value:     force,
		Timeout:   timeout,
	}
}

// TimerCondition creates a timer condition
func TimerCondition(operator ComparisonOperator, duration *TimeValue) *Condition {
	return &Condition{
		Type:      CondTimer,
		Parameter: "timer",
		Operator:  operator,
		Value:     duration,
	}
}

// VariableCondition creates a variable condition
func VariableCondition(varName string, operator ComparisonOperator, value interface{}) *Condition {
	return &Condition{
		Type:      CondVariable,
		Parameter: varName,
		Operator:  operator,
		Value:     value,
	}
}

// DriveStateCondition creates a drive state condition
func DriveStateCondition(operator ComparisonOperator, state DriveState) *Condition {
	return &Condition{
		Type:      CondDriveState,
		Parameter: "drive_state",
		Operator:  operator,
		Value:     state,
	}
}

// MotionCompleteCondition creates a motion completion condition
func MotionCompleteCondition(timeout *TimeValue) *Condition {
	return &Condition{
		Type:      CondMotionComplete,
		Parameter: "motion_complete",
		Operator:  OpEqual,
		Value:     true,
		Timeout:   timeout,
	}
}

// DriveState represents the state of a drive
type DriveState int

const (
	DriveStateDisabled    DriveState = iota
	DriveStateEnabled
	DriveStateHoming
	DriveStateMoving
	DriveStateHolding
	DriveStateError
	DriveStateFault
	DriveStateReady
)

// String returns the string representation of the drive state
func (ds DriveState) String() string {
	switch ds {
	case DriveStateDisabled:
		return "disabled"
	case DriveStateEnabled:
		return "enabled"
	case DriveStateHoming:
		return "homing"
	case DriveStateMoving:
		return "moving"
	case DriveStateHolding:
		return "holding"
	case DriveStateError:
		return "error"
	case DriveStateFault:
		return "fault"
	case DriveStateReady:
		return "ready"
	default:
		return "unknown"
	}
}

// ConditionGroup represents a group of conditions with logical operators
type ConditionGroup struct {
	Conditions []Condition        `json:"conditions" yaml:"conditions"`
	Operator   ComparisonOperator `json:"operator" yaml:"operator"`
}

// Evaluate evaluates the condition group
func (cg *ConditionGroup) Evaluate(ctx context.Context, evaluator ConditionEvaluator, vars map[string]interface{}) (bool, error) {
	if len(cg.Conditions) == 0 {
		return true, nil
	}
	
	switch cg.Operator {
	case OpAnd:
		for _, condition := range cg.Conditions {
			result, err := evaluator.Evaluate(ctx, &condition, vars)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
		return true, nil
		
	case OpOr:
		for _, condition := range cg.Conditions {
			result, err := evaluator.Evaluate(ctx, &condition, vars)
			if err != nil {
				return false, err
			}
			if result {
				return true, nil
			}
		}
		return false, nil
		
	default:
		return false, fmt.Errorf("unsupported group operator: %s", cg.Operator.String())
	}
}

// ConditionValidator defines the interface for validating conditions
type ConditionValidator interface {
	ValidateCondition(condition *Condition) error
	ValidateConditionGroup(group *ConditionGroup) error
	CheckDependencies(condition *Condition) error
}

// DefaultConditionEvaluator provides a default implementation of ConditionEvaluator
type DefaultConditionEvaluator struct {
	driveStatusProvider DriveStatusProvider
}

// DriveStatusProvider defines the interface for providing drive status
type DriveStatusProvider interface {
	GetDriveStatus(driveID int) (*DriveStatus, error)
	GetDigitalInput(driveID, inputNum int) (bool, error)
	GetAnalogInput(driveID, inputNum int) (float64, error)
	GetPosition(driveID int) (*PositionValue, error)
	GetVelocity(driveID int) (*VelocityValue, error)
	GetForce(driveID int) (*ForceValue, error)
}

// DriveStatus represents the status of a drive
type DriveStatus struct {
	DriveID        int
	State          DriveState
	Position       *PositionValue
	Velocity       *VelocityValue
	Force          *ForceValue
	DigitalInputs  []bool
	AnalogInputs   []float64
	DigitalOutputs []bool
	AnalogOutputs  []float64
	Error          error
	Timestamp      time.Time
}

// NewDefaultConditionEvaluator creates a new default condition evaluator
func NewDefaultConditionEvaluator(provider DriveStatusProvider) *DefaultConditionEvaluator {
	return &DefaultConditionEvaluator{
		driveStatusProvider: provider,
	}
}

// Evaluate evaluates a condition
func (dce *DefaultConditionEvaluator) Evaluate(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	switch condition.Type {
	case CondDigitalInput:
		return dce.evaluateDigitalInput(ctx, condition, vars)
	case CondAnalogInput:
		return dce.evaluateAnalogInput(ctx, condition, vars)
	case CondPosition:
		return dce.evaluatePosition(ctx, condition, vars)
	case CondVelocity:
		return dce.evaluateVelocity(ctx, condition, vars)
	case CondForce:
		return dce.evaluateForce(ctx, condition, vars)
	case CondTimer:
		return dce.evaluateTimer(ctx, condition, vars)
	case CondVariable:
		return dce.evaluateVariable(ctx, condition, vars)
	case CondDriveState:
		return dce.evaluateDriveState(ctx, condition, vars)
	case CondMotionComplete:
		return dce.evaluateMotionComplete(ctx, condition, vars)
	default:
		return false, fmt.Errorf("unknown condition type: %s", condition.Type.String())
	}
}

// CanEvaluate checks if the condition can be evaluated
func (dce *DefaultConditionEvaluator) CanEvaluate(condition *Condition) bool {
	return condition.Type != CondError // Error conditions are handled separately
}

// GetRequiredData returns the data required for the condition
func (dce *DefaultConditionEvaluator) GetRequiredData(condition *Condition) []string {
	switch condition.Type {
	case CondDigitalInput, CondAnalogInput, CondPosition, CondVelocity, CondForce, CondDriveState, CondMotionComplete:
		return []string{"drive_status"}
	case CondTimer:
		return []string{"timer"}
	case CondVariable:
		return []string{"variables"}
	default:
		return []string{}
	}
}

// evaluateDigitalInput evaluates a digital input condition
func (dce *DefaultConditionEvaluator) evaluateDigitalInput(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	// Extract drive ID and input number from parameter
	driveID := 1 // Default drive ID, should be extracted from context or vars
	inputNum := 1 // Default input number, should be extracted from parameter
	
	actualState, err := dce.driveStatusProvider.GetDigitalInput(driveID, inputNum)
	if err != nil {
		return false, err
	}
	
	expectedState, ok := condition.Value.(bool)
	if !ok {
		return false, fmt.Errorf("expected boolean value for digital input condition")
	}
	
	return actualState == expectedState, nil
}

// evaluateAnalogInput evaluates an analog input condition
func (dce *DefaultConditionEvaluator) evaluateAnalogInput(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	// Extract drive ID and input number from parameter
	driveID := 1 // Default drive ID, should be extracted from context or vars
	inputNum := 1 // Default input number, should be extracted from parameter
	
	actualValue, err := dce.driveStatusProvider.GetAnalogInput(driveID, inputNum)
	if err != nil {
		return false, err
	}
	
	expectedValue, ok := condition.Value.(float64)
	if !ok {
		return false, fmt.Errorf("expected float64 value for analog input condition")
	}
	
	return dce.compareValues(actualValue, condition.Operator, expectedValue), nil
}

// evaluatePosition evaluates a position condition
func (dce *DefaultConditionEvaluator) evaluatePosition(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	driveID := 1 // Default drive ID, should be extracted from context or vars
	
	actualPosition, err := dce.driveStatusProvider.GetPosition(driveID)
	if err != nil {
		return false, err
	}
	
	expectedPosition, ok := condition.Value.(*PositionValue)
	if !ok {
		return false, fmt.Errorf("expected PositionValue for position condition")
	}
	
	// Convert to same units for comparison
	actualValue := actualPosition.Value
	expectedValue := expectedPosition.Value
	
	// TODO: Add unit conversion logic here
	
	return dce.compareValues(actualValue, condition.Operator, expectedValue), nil
}

// evaluateVelocity evaluates a velocity condition
func (dce *DefaultConditionEvaluator) evaluateVelocity(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	driveID := 1 // Default drive ID, should be extracted from context or vars
	
	actualVelocity, err := dce.driveStatusProvider.GetVelocity(driveID)
	if err != nil {
		return false, err
	}
	
	expectedVelocity, ok := condition.Value.(*VelocityValue)
	if !ok {
		return false, fmt.Errorf("expected VelocityValue for velocity condition")
	}
	
	// Convert to same units for comparison
	actualValue := actualVelocity.Value
	expectedValue := expectedVelocity.Value
	
	// TODO: Add unit conversion logic here
	
	return dce.compareValues(actualValue, condition.Operator, expectedValue), nil
}

// evaluateForce evaluates a force condition
func (dce *DefaultConditionEvaluator) evaluateForce(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	driveID := 1 // Default drive ID, should be extracted from context or vars
	
	actualForce, err := dce.driveStatusProvider.GetForce(driveID)
	if err != nil {
		return false, err
	}
	
	expectedForce, ok := condition.Value.(*ForceValue)
	if !ok {
		return false, fmt.Errorf("expected ForceValue for force condition")
	}
	
	// Convert to same units for comparison
	actualValue := actualForce.Value
	expectedValue := expectedForce.Value
	
	// TODO: Add unit conversion logic here
	
	return dce.compareValues(actualValue, condition.Operator, expectedValue), nil
}

// evaluateTimer evaluates a timer condition
func (dce *DefaultConditionEvaluator) evaluateTimer(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	// Timer conditions are typically handled by the execution engine
	// This is a placeholder implementation
	return false, fmt.Errorf("timer conditions not implemented")
}

// evaluateVariable evaluates a variable condition
func (dce *DefaultConditionEvaluator) evaluateVariable(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	actualValue, ok := vars[condition.Parameter]
	if !ok {
		return false, fmt.Errorf("variable %s not found", condition.Parameter)
	}
	
	return dce.compareValues(actualValue, condition.Operator, condition.Value), nil
}

// evaluateDriveState evaluates a drive state condition
func (dce *DefaultConditionEvaluator) evaluateDriveState(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	driveID := 1 // Default drive ID, should be extracted from context or vars
	
	status, err := dce.driveStatusProvider.GetDriveStatus(driveID)
	if err != nil {
		return false, err
	}
	
	expectedState, ok := condition.Value.(DriveState)
	if !ok {
		return false, fmt.Errorf("expected DriveState for drive state condition")
	}
	
	return status.State == expectedState, nil
}

// evaluateMotionComplete evaluates a motion completion condition
func (dce *DefaultConditionEvaluator) evaluateMotionComplete(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	driveID := 1 // Default drive ID, should be extracted from context or vars
	
	status, err := dce.driveStatusProvider.GetDriveStatus(driveID)
	if err != nil {
		return false, err
	}
	
	// Motion is complete when the drive is not moving and not in error state
	return status.State != DriveStateMoving && status.State != DriveStateError, nil
}

// compareValues compares two values using the specified operator
func (dce *DefaultConditionEvaluator) compareValues(actual, expected interface{}) bool {
	switch dce.condition.Operator {
	case OpEqual:
		return actual == expected
	case OpNotEqual:
		return actual != expected
	case OpGreaterThan:
		return dce.compareNumeric(actual, expected) > 0
	case OpLessThan:
		return dce.compareNumeric(actual, expected) < 0
	case OpGreaterThanOrEqual:
		return dce.compareNumeric(actual, expected) >= 0
	case OpLessThanOrEqual:
		return dce.compareNumeric(actual, expected) <= 0
	default:
		return false
	}
}

// compareNumeric compares two numeric values
func (dce *DefaultConditionEvaluator) compareNumeric(actual, expected interface{}) int {
	// Convert to float64 for comparison
	actualFloat, ok1 := dce.toFloat64(actual)
	expectedFloat, ok2 := dce.toFloat64(expected)
	
	if !ok1 || !ok2 {
		return 0
	}
	
	if actualFloat > expectedFloat {
		return 1
	} else if actualFloat < expectedFloat {
		return -1
	}
	return 0
}

// toFloat64 converts a value to float64
func (dce *DefaultConditionEvaluator) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	default:
		return 0, false
	}
}