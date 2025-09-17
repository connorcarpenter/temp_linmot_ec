package types

import (
	"context"
	"testing"
)

func TestConditionType_String(t *testing.T) {
	tests := []struct {
		name     string
		condType ConditionType
		expected string
	}{
		{"DigitalInput", CondDigitalInput, "digital_input"},
		{"AnalogInput", CondAnalogInput, "analog_input"},
		{"Position", CondPosition, "position"},
		{"Velocity", CondVelocity, "velocity"},
		{"Force", CondForce, "force"},
		{"Timer", CondTimer, "timer"},
		{"Variable", CondVariable, "variable"},
		{"Error", CondError, "error"},
		{"DriveState", CondDriveState, "drive_state"},
		{"MotionComplete", CondMotionComplete, "motion_complete"},
		{"Unknown", ConditionType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condType.String(); got != tt.expected {
				t.Errorf("ConditionType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestComparisonOperator_String(t *testing.T) {
	tests := []struct {
		name     string
		operator ComparisonOperator
		expected string
	}{
		{"Equal", OpEqual, "=="},
		{"NotEqual", OpNotEqual, "!="},
		{"GreaterThan", OpGreaterThan, ">"},
		{"LessThan", OpLessThan, "<"},
		{"GreaterThanOrEqual", OpGreaterThanOrEqual, ">="},
		{"LessThanOrEqual", OpLessThanOrEqual, "<="},
		{"And", OpAnd, "&&"},
		{"Or", OpOr, "||"},
		{"Not", OpNot, "!"},
		{"In", OpIn, "in"},
		{"NotIn", OpNotIn, "not in"},
		{"Unknown", ComparisonOperator(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.String(); got != tt.expected {
				t.Errorf("ComparisonOperator.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConditionBuilder(t *testing.T) {
	timeout := NewTimeValue(5000.0, TimeUnitMS)
	
	condition := NewConditionBuilder().
		WithType(CondDigitalInput).
		WithParameter("digital_input_1").
		WithOperator(OpEqual).
		WithValue(true).
		WithTimeout(timeout).
		WithDescription("Check if digital input 1 is high").
		Build()

	if condition.Type != CondDigitalInput {
		t.Errorf("Expected type CondDigitalInput, got %v", condition.Type)
	}

	if condition.Parameter != "digital_input_1" {
		t.Errorf("Expected parameter 'digital_input_1', got %s", condition.Parameter)
	}

	if condition.Operator != OpEqual {
		t.Errorf("Expected operator OpEqual, got %v", condition.Operator)
	}

	if condition.Value != true {
		t.Errorf("Expected value true, got %v", condition.Value)
	}

	if condition.Timeout == nil {
		t.Error("Expected non-nil timeout")
	} else if condition.Timeout.Value != 5000.0 {
		t.Errorf("Expected timeout value 5000.0, got %f", condition.Timeout.Value)
	}

	if condition.Description != "Check if digital input 1 is high" {
		t.Errorf("Expected description 'Check if digital input 1 is high', got %s", condition.Description)
	}
}

func TestDigitalInputCondition(t *testing.T) {
	timeout := NewTimeValue(1000.0, TimeUnitMS)
	condition := DigitalInputCondition(1, true, timeout)

	if condition.Type != CondDigitalInput {
		t.Errorf("Expected type CondDigitalInput, got %v", condition.Type)
	}

	if condition.Parameter != "digital_input_1" {
		t.Errorf("Expected parameter 'digital_input_1', got %s", condition.Parameter)
	}

	if condition.Operator != OpEqual {
		t.Errorf("Expected operator OpEqual, got %v", condition.Operator)
	}

	if condition.Value != true {
		t.Errorf("Expected value true, got %v", condition.Value)
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestAnalogInputCondition(t *testing.T) {
	timeout := NewTimeValue(2000.0, TimeUnitMS)
	condition := AnalogInputCondition(2, OpGreaterThan, 3.14, 0.01, timeout)

	if condition.Type != CondAnalogInput {
		t.Errorf("Expected type CondAnalogInput, got %v", condition.Type)
	}

	if condition.Parameter != "analog_input_2" {
		t.Errorf("Expected parameter 'analog_input_2', got %s", condition.Parameter)
	}

	if condition.Operator != OpGreaterThan {
		t.Errorf("Expected operator OpGreaterThan, got %v", condition.Operator)
	}

	if condition.Value != 3.14 {
		t.Errorf("Expected value 3.14, got %v", condition.Value)
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestPositionCondition(t *testing.T) {
	position := NewPositionValue(100.0, PositionUnitMM)
	timeout := NewTimeValue(5000.0, TimeUnitMS)
	condition := PositionCondition(OpLessThan, position, timeout)

	if condition.Type != CondPosition {
		t.Errorf("Expected type CondPosition, got %v", condition.Type)
	}

	if condition.Parameter != "position" {
		t.Errorf("Expected parameter 'position', got %s", condition.Parameter)
	}

	if condition.Operator != OpLessThan {
		t.Errorf("Expected operator OpLessThan, got %v", condition.Operator)
	}

	if condition.Value != position {
		t.Errorf("Expected value to match provided position")
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestVelocityCondition(t *testing.T) {
	velocity := NewVelocityValue(50.0, VelocityUnitMMS)
	timeout := NewTimeValue(3000.0, TimeUnitMS)
	condition := VelocityCondition(OpGreaterThanOrEqual, velocity, timeout)

	if condition.Type != CondVelocity {
		t.Errorf("Expected type CondVelocity, got %v", condition.Type)
	}

	if condition.Parameter != "velocity" {
		t.Errorf("Expected parameter 'velocity', got %s", condition.Parameter)
	}

	if condition.Operator != OpGreaterThanOrEqual {
		t.Errorf("Expected operator OpGreaterThanOrEqual, got %v", condition.Operator)
	}

	if condition.Value != velocity {
		t.Errorf("Expected value to match provided velocity")
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestForceCondition(t *testing.T) {
	force := NewForceValue(1000.0, ForceUnitN)
	timeout := NewTimeValue(10000.0, TimeUnitMS)
	condition := ForceCondition(OpEqual, force, timeout)

	if condition.Type != CondForce {
		t.Errorf("Expected type CondForce, got %v", condition.Type)
	}

	if condition.Parameter != "force" {
		t.Errorf("Expected parameter 'force', got %s", condition.Parameter)
	}

	if condition.Operator != OpEqual {
		t.Errorf("Expected operator OpEqual, got %v", condition.Operator)
	}

	if condition.Value != force {
		t.Errorf("Expected value to match provided force")
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestTimerCondition(t *testing.T) {
	duration := NewTimeValue(5000.0, TimeUnitMS)
	condition := TimerCondition(OpGreaterThan, duration)

	if condition.Type != CondTimer {
		t.Errorf("Expected type CondTimer, got %v", condition.Type)
	}

	if condition.Parameter != "timer" {
		t.Errorf("Expected parameter 'timer', got %s", condition.Parameter)
	}

	if condition.Operator != OpGreaterThan {
		t.Errorf("Expected operator OpGreaterThan, got %v", condition.Operator)
	}

	if condition.Value != duration {
		t.Errorf("Expected value to match provided duration")
	}
}

func TestVariableCondition(t *testing.T) {
	condition := VariableCondition("test_var", OpNotEqual, "test_value")

	if condition.Type != CondVariable {
		t.Errorf("Expected type CondVariable, got %v", condition.Type)
	}

	if condition.Parameter != "test_var" {
		t.Errorf("Expected parameter 'test_var', got %s", condition.Parameter)
	}

	if condition.Operator != OpNotEqual {
		t.Errorf("Expected operator OpNotEqual, got %v", condition.Operator)
	}

	if condition.Value != "test_value" {
		t.Errorf("Expected value 'test_value', got %v", condition.Value)
	}
}

func TestDriveStateCondition(t *testing.T) {
	condition := DriveStateCondition(OpEqual, DriveStateEnabled)

	if condition.Type != CondDriveState {
		t.Errorf("Expected type CondDriveState, got %v", condition.Type)
	}

	if condition.Parameter != "drive_state" {
		t.Errorf("Expected parameter 'drive_state', got %s", condition.Parameter)
	}

	if condition.Operator != OpEqual {
		t.Errorf("Expected operator OpEqual, got %v", condition.Operator)
	}

	if condition.Value != DriveStateEnabled {
		t.Errorf("Expected value DriveStateEnabled, got %v", condition.Value)
	}
}

func TestMotionCompleteCondition(t *testing.T) {
	timeout := NewTimeValue(30000.0, TimeUnitMS)
	condition := MotionCompleteCondition(timeout)

	if condition.Type != CondMotionComplete {
		t.Errorf("Expected type CondMotionComplete, got %v", condition.Type)
	}

	if condition.Parameter != "motion_complete" {
		t.Errorf("Expected parameter 'motion_complete', got %s", condition.Parameter)
	}

	if condition.Operator != OpEqual {
		t.Errorf("Expected operator OpEqual, got %v", condition.Operator)
	}

	if condition.Value != true {
		t.Errorf("Expected value true, got %v", condition.Value)
	}

	if condition.Timeout != timeout {
		t.Errorf("Expected timeout to match provided timeout")
	}
}

func TestDriveState_String(t *testing.T) {
	tests := []struct {
		name     string
		state    DriveState
		expected string
	}{
		{"Disabled", DriveStateDisabled, "disabled"},
		{"Enabled", DriveStateEnabled, "enabled"},
		{"Homing", DriveStateHoming, "homing"},
		{"Moving", DriveStateMoving, "moving"},
		{"Holding", DriveStateHolding, "holding"},
		{"Error", DriveStateError, "error"},
		{"Fault", DriveStateFault, "fault"},
		{"Ready", DriveStateReady, "ready"},
		{"Unknown", DriveState(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("DriveState.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConditionGroup_AndEvaluation(t *testing.T) {
	// Mock condition evaluator
	evaluator := &MockConditionEvaluator{
		results: map[string]bool{
			"condition1": true,
			"condition2": true,
			"condition3": false,
		},
	}

	group := ConditionGroup{
		Conditions: []Condition{
			{Type: CondVariable, Parameter: "condition1", Operator: OpEqual, Value: true},
			{Type: CondVariable, Parameter: "condition2", Operator: OpEqual, Value: true},
			{Type: CondVariable, Parameter: "condition3", Operator: OpEqual, Value: true},
		},
		Operator: OpAnd,
	}

	ctx := context.Background()
	vars := map[string]interface{}{}

	result, err := group.Evaluate(ctx, evaluator, vars)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should be false because condition3 is false
	if result != false {
		t.Errorf("Expected false, got %t", result)
	}
}

func TestConditionGroup_OrEvaluation(t *testing.T) {
	// Mock condition evaluator
	evaluator := &MockConditionEvaluator{
		results: map[string]bool{
			"condition1": false,
			"condition2": false,
			"condition3": true,
		},
	}

	group := ConditionGroup{
		Conditions: []Condition{
			{Type: CondVariable, Parameter: "condition1", Operator: OpEqual, Value: true},
			{Type: CondVariable, Parameter: "condition2", Operator: OpEqual, Value: true},
			{Type: CondVariable, Parameter: "condition3", Operator: OpEqual, Value: true},
		},
		Operator: OpOr,
	}

	ctx := context.Background()
	vars := map[string]interface{}{}

	result, err := group.Evaluate(ctx, evaluator, vars)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should be true because condition3 is true
	if result != true {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestConditionGroup_EmptyConditions(t *testing.T) {
	evaluator := &MockConditionEvaluator{}

	group := ConditionGroup{
		Conditions: []Condition{},
		Operator:   OpAnd,
	}

	ctx := context.Background()
	vars := map[string]interface{}{}

	result, err := group.Evaluate(ctx, evaluator, vars)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Empty conditions should evaluate to true
	if result != true {
		t.Errorf("Expected true for empty conditions, got %t", result)
	}
}

func TestConditionGroup_UnsupportedOperator(t *testing.T) {
	evaluator := &MockConditionEvaluator{}

	group := ConditionGroup{
		Conditions: []Condition{
			{Type: CondVariable, Parameter: "condition1", Operator: OpEqual, Value: true},
		},
		Operator: OpNot, // Unsupported for groups
	}

	ctx := context.Background()
	vars := map[string]interface{}{}

	_, err := group.Evaluate(ctx, evaluator, vars)
	if err == nil {
		t.Error("Expected error for unsupported operator")
	}
}

func TestDefaultConditionEvaluator_CanEvaluate(t *testing.T) {
	evaluator := &TestConditionEvaluator{}

	tests := []struct {
		name     string
		condition *Condition
		expected bool
	}{
		{
			name: "Digital input condition",
			condition: &Condition{
				Type: CondDigitalInput,
			},
			expected: true,
		},
		{
			name: "Error condition",
			condition: &Condition{
				Type: CondError,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.CanEvaluate(tt.condition)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestDefaultConditionEvaluator_GetRequiredData(t *testing.T) {
	evaluator := &TestConditionEvaluator{}

	tests := []struct {
		name     string
		condition *Condition
		expected []string
	}{
		{
			name: "Digital input condition",
			condition: &Condition{
				Type: CondDigitalInput,
			},
			expected: []string{"drive_status"},
		},
		{
			name: "Variable condition",
			condition: &Condition{
				Type: CondVariable,
			},
			expected: []string{"variables"},
		},
		{
			name: "Timer condition",
			condition: &Condition{
				Type: CondTimer,
			},
			expected: []string{"timer"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluator.GetRequiredData(tt.condition)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d required data items, got %d", len(tt.expected), len(result))
			}
		})
	}
}

func TestDefaultConditionEvaluator_CompareValues(t *testing.T) {
	evaluator := &TestConditionEvaluator{}

	tests := []struct {
		name     string
		actual   interface{}
		operator ComparisonOperator
		expected interface{}
		result   bool
	}{
		{"Equal int", 42, OpEqual, 42, true},
		{"Not equal int", 42, OpNotEqual, 43, true},
		{"Greater than int", 43, OpGreaterThan, 42, true},
		{"Less than int", 41, OpLessThan, 42, true},
		{"Greater than or equal int", 42, OpGreaterThanOrEqual, 42, true},
		{"Less than or equal int", 42, OpLessThanOrEqual, 42, true},
		{"Equal float", 42.0, OpEqual, 42.0, true},
		{"Not equal float", 42.0, OpNotEqual, 43.0, true},
		{"Greater than float", 43.0, OpGreaterThan, 42.0, true},
		{"Less than float", 41.0, OpLessThan, 42.0, true},
		{"Equal string", "test", OpEqual, "test", true},
		{"Not equal string", "test", OpNotEqual, "other", true},
		{"Equal bool", true, OpEqual, true, true},
		{"Not equal bool", true, OpNotEqual, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a condition with the operator
			condition := &Condition{
				Operator: tt.operator,
			}
			evaluator.condition = condition

			result := evaluator.compareValues(tt.actual, tt.expected)
			if result != tt.result {
				t.Errorf("Expected %t, got %t", tt.result, result)
			}
		})
	}
}

func TestDefaultConditionEvaluator_ToFloat64(t *testing.T) {
	evaluator := &TestConditionEvaluator{}

	tests := []struct {
		name     string
		value    interface{}
		expected float64
		valid    bool
	}{
		{"Float64", 42.5, 42.5, true},
		{"Int", 42, 42.0, true},
		{"Int64", int64(42), 42.0, true},
		{"Float32", float32(42.5), 42.5, true},
		{"String", "invalid", 0.0, false},
		{"Bool", true, 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, valid := evaluator.toFloat64(tt.value)
			if valid != tt.valid {
				t.Errorf("Expected valid %t, got %t", tt.valid, valid)
			}
			if valid && result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

// MockConditionEvaluator for testing
type MockConditionEvaluator struct {
	results map[string]bool
}

func (mce *MockConditionEvaluator) Evaluate(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	if result, ok := mce.results[condition.Parameter]; ok {
		return result, nil
	}
	return false, nil
}

func (mce *MockConditionEvaluator) CanEvaluate(condition *Condition) bool {
	return true
}

func (mce *MockConditionEvaluator) GetRequiredData(condition *Condition) []string {
	return []string{"test_data"}
}

// TestConditionEvaluator for testing
type TestConditionEvaluator struct {
	condition *Condition
}

func (tce *TestConditionEvaluator) Evaluate(ctx context.Context, condition *Condition, vars map[string]interface{}) (bool, error) {
	return false, nil
}

func (tce *TestConditionEvaluator) CanEvaluate(condition *Condition) bool {
	return condition.Type != CondError
}

func (tce *TestConditionEvaluator) GetRequiredData(condition *Condition) []string {
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

func (tce *TestConditionEvaluator) compareValues(actual, expected interface{}) bool {
	switch tce.condition.Operator {
	case OpEqual:
		return actual == expected
	case OpNotEqual:
		return actual != expected
	case OpGreaterThan:
		return tce.compareNumeric(actual, expected) > 0
	case OpLessThan:
		return tce.compareNumeric(actual, expected) < 0
	case OpGreaterThanOrEqual:
		return tce.compareNumeric(actual, expected) >= 0
	case OpLessThanOrEqual:
		return tce.compareNumeric(actual, expected) <= 0
	default:
		return false
	}
}

func (tce *TestConditionEvaluator) compareNumeric(actual, expected interface{}) int {
	actualFloat, ok1 := tce.toFloat64(actual)
	expectedFloat, ok2 := tce.toFloat64(expected)

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

func (tce *TestConditionEvaluator) toFloat64(value interface{}) (float64, bool) {
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