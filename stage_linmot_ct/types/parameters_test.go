package types

import (
	"encoding/json"
	"testing"
)

func TestPositionUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     PositionUnit
		expected string
	}{
		{"Millimeters", PositionUnitMM, "mm"},
		{"Counts", PositionUnitCounts, "counts"},
		{"Unknown", PositionUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("PositionUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestVelocityUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     VelocityUnit
		expected string
	}{
		{"MillimetersPerSecond", VelocityUnitMMS, "mm/s"},
		{"CountsPerSecond", VelocityUnitCountsS, "counts/s"},
		{"Unknown", VelocityUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("VelocityUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAccelerationUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     AccelerationUnit
		expected string
	}{
		{"MillimetersPerSecondSquared", AccelerationUnitMMS2, "mm/s²"},
		{"CountsPerSecondSquared", AccelerationUnitCountsS2, "counts/s²"},
		{"Unknown", AccelerationUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("AccelerationUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestJerkUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     JerkUnit
		expected string
	}{
		{"MillimetersPerSecondCubed", JerkUnitMMS3, "mm/s³"},
		{"CountsPerSecondCubed", JerkUnitCountsS3, "counts/s³"},
		{"Unknown", JerkUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("JerkUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestForceUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     ForceUnit
		expected string
	}{
		{"Newtons", ForceUnitN, "N"},
		{"Counts", ForceUnitCounts, "counts"},
		{"Unknown", ForceUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("ForceUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTimeUnit_String(t *testing.T) {
	tests := []struct {
		name     string
		unit     TimeUnit
		expected string
	}{
		{"Milliseconds", TimeUnitMS, "ms"},
		{"Seconds", TimeUnitS, "s"},
		{"Unknown", TimeUnit(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.unit.String(); got != tt.expected {
				t.Errorf("TimeUnit.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewPositionValue(t *testing.T) {
	pv := NewPositionValue(100.0, PositionUnitMM)

	if pv.Value != 100.0 {
		t.Errorf("Expected value 100.0, got %f", pv.Value)
	}

	if pv.Unit != PositionUnitMM {
		t.Errorf("Expected unit PositionUnitMM, got %v", pv.Unit)
	}
}

func TestNewVelocityValue(t *testing.T) {
	vv := NewVelocityValue(50.0, VelocityUnitMMS)

	if vv.Value != 50.0 {
		t.Errorf("Expected value 50.0, got %f", vv.Value)
	}

	if vv.Unit != VelocityUnitMMS {
		t.Errorf("Expected unit VelocityUnitMMS, got %v", vv.Unit)
	}
}

func TestNewAccelerationValue(t *testing.T) {
	av := NewAccelerationValue(25.0, AccelerationUnitMMS2)

	if av.Value != 25.0 {
		t.Errorf("Expected value 25.0, got %f", av.Value)
	}

	if av.Unit != AccelerationUnitMMS2 {
		t.Errorf("Expected unit AccelerationUnitMMS2, got %v", av.Unit)
	}
}

func TestNewJerkValue(t *testing.T) {
	jv := NewJerkValue(10.0, JerkUnitMMS3)

	if jv.Value != 10.0 {
		t.Errorf("Expected value 10.0, got %f", jv.Value)
	}

	if jv.Unit != JerkUnitMMS3 {
		t.Errorf("Expected unit JerkUnitMMS3, got %v", jv.Unit)
	}
}

func TestNewForceValue(t *testing.T) {
	fv := NewForceValue(1000.0, ForceUnitN)

	if fv.Value != 1000.0 {
		t.Errorf("Expected value 1000.0, got %f", fv.Value)
	}

	if fv.Unit != ForceUnitN {
		t.Errorf("Expected unit ForceUnitN, got %v", fv.Unit)
	}
}

func TestNewTimeValue(t *testing.T) {
	tv := NewTimeValue(1000.0, TimeUnitMS)

	if tv.Value != 1000.0 {
		t.Errorf("Expected value 1000.0, got %f", tv.Value)
	}

	if tv.Unit != TimeUnitMS {
		t.Errorf("Expected unit TimeUnitMS, got %v", tv.Unit)
	}
}

func TestPositionValue_String(t *testing.T) {
	pv := NewPositionValue(123.456, PositionUnitMM)
	expected := "123.456 mm"

	if pv.String() != expected {
		t.Errorf("Expected %s, got %s", expected, pv.String())
	}
}

func TestVelocityValue_String(t *testing.T) {
	vv := NewVelocityValue(45.67, VelocityUnitMMS)
	expected := "45.670 mm/s"

	if vv.String() != expected {
		t.Errorf("Expected %s, got %s", expected, vv.String())
	}
}

func TestAccelerationValue_String(t *testing.T) {
	av := NewAccelerationValue(12.34, AccelerationUnitMMS2)
	expected := "12.340 mm/s²"

	if av.String() != expected {
		t.Errorf("Expected %s, got %s", expected, av.String())
	}
}

func TestJerkValue_String(t *testing.T) {
	jv := NewJerkValue(5.67, JerkUnitMMS3)
	expected := "5.670 mm/s³"

	if jv.String() != expected {
		t.Errorf("Expected %s, got %s", expected, jv.String())
	}
}

func TestForceValue_String(t *testing.T) {
	fv := NewForceValue(987.65, ForceUnitN)
	expected := "987.650 N"

	if fv.String() != expected {
		t.Errorf("Expected %s, got %s", expected, fv.String())
	}
}

func TestTimeValue_String(t *testing.T) {
	tv := NewTimeValue(1500.0, TimeUnitMS)
	expected := "1500.000 ms"

	if tv.String() != expected {
		t.Errorf("Expected %s, got %s", expected, tv.String())
	}
}

func TestTimeValue_Duration(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		unit     TimeUnit
		expected float64 // in milliseconds
	}{
		{"Milliseconds", 1000.0, TimeUnitMS, 1000.0},
		{"Seconds", 1.0, TimeUnitS, 1000.0},
		{"HalfSecond", 0.5, TimeUnitS, 500.0},
		{"TwoSeconds", 2.0, TimeUnitS, 2000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := NewTimeValue(tt.value, tt.unit)
			duration := tv.Duration()
		expectedDuration := tt.expected * 1e6 // Convert to nanoseconds
		actualDuration := float64(duration.Nanoseconds())
		// Allow for small floating point precision differences
		if actualDuration < expectedDuration-1 || actualDuration > expectedDuration+1 {
			t.Errorf("Expected duration %f ns, got %f ns", expectedDuration, actualDuration)
		}
		})
	}
}

func TestParameterExtractor_ExtractPosition(t *testing.T) {
	pe := NewParameterExtractor()

	tests := []struct {
		name        string
		params      map[string]interface{}
		key         string
		expected    *PositionValue
		expectError bool
	}{
		{
			name: "Valid position with unit",
			params: map[string]interface{}{
				"position": map[string]interface{}{
					"value": 100.0,
					"unit":  "mm",
				},
			},
			key: "position",
			expected: &PositionValue{
				Value: 100.0,
				Unit:  PositionUnitMM,
			},
			expectError: false,
		},
		{
			name: "Valid position with counts unit",
			params: map[string]interface{}{
				"position": map[string]interface{}{
					"value": 1000.0,
					"unit":  "counts",
				},
			},
			key: "position",
			expected: &PositionValue{
				Value: 1000.0,
				Unit:  PositionUnitCounts,
			},
			expectError: false,
		},
		{
			name: "Position as simple float",
			params: map[string]interface{}{
				"position": 100.0,
			},
			key: "position",
			expected: &PositionValue{
				Value: 100.0,
				Unit:  PositionUnitMM, // Default unit
			},
			expectError: false,
		},
		{
			name: "Missing parameter",
			params: map[string]interface{}{
				"other": 100.0,
			},
			key:         "position",
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid position value",
			params: map[string]interface{}{
				"position": "invalid",
			},
			key:         "position",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pe.ExtractPosition(tt.params, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			if result.Value != tt.expected.Value {
				t.Errorf("Expected value %f, got %f", tt.expected.Value, result.Value)
			}

			if result.Unit != tt.expected.Unit {
				t.Errorf("Expected unit %v, got %v", tt.expected.Unit, result.Unit)
			}
		})
	}
}

func TestParameterExtractor_ExtractVelocity(t *testing.T) {
	pe := NewParameterExtractor()

	params := map[string]interface{}{
		"velocity": map[string]interface{}{
			"value": 50.0,
			"unit":  "mm/s",
		},
	}

	result, err := pe.ExtractVelocity(params, "velocity")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.Value != 50.0 {
		t.Errorf("Expected value 50.0, got %f", result.Value)
	}

	if result.Unit != VelocityUnitMMS {
		t.Errorf("Expected unit VelocityUnitMMS, got %v", result.Unit)
	}
}

func TestParameterExtractor_ExtractTime(t *testing.T) {
	pe := NewParameterExtractor()

	params := map[string]interface{}{
		"timeout": map[string]interface{}{
			"value": 1000.0,
			"unit":  "ms",
		},
	}

	result, err := pe.ExtractTime(params, "timeout")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result.Value != 1000.0 {
		t.Errorf("Expected value 1000.0, got %f", result.Value)
	}

	if result.Unit != TimeUnitMS {
		t.Errorf("Expected unit TimeUnitMS, got %v", result.Unit)
	}
}

func TestParameterExtractor_ExtractInt(t *testing.T) {
	pe := NewParameterExtractor()

	tests := []struct {
		name        string
		params      map[string]interface{}
		key         string
		expected    int
		expectError bool
	}{
		{
			name: "Valid int",
			params: map[string]interface{}{
				"value": 42,
			},
			key:         "value",
			expected:    42,
			expectError: false,
		},
		{
			name: "Valid float64",
			params: map[string]interface{}{
				"value": 42.0,
			},
			key:         "value",
			expected:    42,
			expectError: false,
		},
		{
			name: "Valid string",
			params: map[string]interface{}{
				"value": "42",
			},
			key:         "value",
			expected:    42,
			expectError: false,
		},
		{
			name: "Missing parameter",
			params: map[string]interface{}{
				"other": 42,
			},
			key:         "value",
			expected:    0,
			expectError: true,
		},
		{
			name: "Invalid string",
			params: map[string]interface{}{
				"value": "invalid",
			},
			key:         "value",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pe.ExtractInt(tt.params, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestParameterExtractor_ExtractFloat(t *testing.T) {
	pe := NewParameterExtractor()

	tests := []struct {
		name        string
		params      map[string]interface{}
		key         string
		expected    float64
		expectError bool
	}{
		{
			name: "Valid float64",
			params: map[string]interface{}{
				"value": 42.5,
			},
			key:         "value",
			expected:    42.5,
			expectError: false,
		},
		{
			name: "Valid int",
			params: map[string]interface{}{
				"value": 42,
			},
			key:         "value",
			expected:    42.0,
			expectError: false,
		},
		{
			name: "Valid string",
			params: map[string]interface{}{
				"value": "42.5",
			},
			key:         "value",
			expected:    42.5,
			expectError: false,
		},
		{
			name: "Missing parameter",
			params: map[string]interface{}{
				"other": 42.5,
			},
			key:         "value",
			expected:    0.0,
			expectError: true,
		},
		{
			name: "Invalid string",
			params: map[string]interface{}{
				"value": "invalid",
			},
			key:         "value",
			expected:    0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pe.ExtractFloat(tt.params, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestParameterExtractor_ExtractBool(t *testing.T) {
	pe := NewParameterExtractor()

	tests := []struct {
		name        string
		params      map[string]interface{}
		key         string
		expected    bool
		expectError bool
	}{
		{
			name: "Valid bool true",
			params: map[string]interface{}{
				"value": true,
			},
			key:         "value",
			expected:    true,
			expectError: false,
		},
		{
			name: "Valid bool false",
			params: map[string]interface{}{
				"value": false,
			},
			key:         "value",
			expected:    false,
			expectError: false,
		},
		{
			name: "Valid string true",
			params: map[string]interface{}{
				"value": "true",
			},
			key:         "value",
			expected:    true,
			expectError: false,
		},
		{
			name: "Valid string false",
			params: map[string]interface{}{
				"value": "false",
			},
			key:         "value",
			expected:    false,
			expectError: false,
		},
		{
			name: "Valid int 1",
			params: map[string]interface{}{
				"value": 1,
			},
			key:         "value",
			expected:    true,
			expectError: false,
		},
		{
			name: "Valid int 0",
			params: map[string]interface{}{
				"value": 0,
			},
			key:         "value",
			expected:    false,
			expectError: false,
		},
		{
			name: "Valid float64 1.0",
			params: map[string]interface{}{
				"value": 1.0,
			},
			key:         "value",
			expected:    true,
			expectError: false,
		},
		{
			name: "Valid float64 0.0",
			params: map[string]interface{}{
				"value": 0.0,
			},
			key:         "value",
			expected:    false,
			expectError: false,
		},
		{
			name: "Missing parameter",
			params: map[string]interface{}{
				"other": true,
			},
			key:         "value",
			expected:    false,
			expectError: true,
		},
		{
			name: "Invalid string",
			params: map[string]interface{}{
				"value": "invalid",
			},
			key:         "value",
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pe.ExtractBool(tt.params, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestParameterExtractor_ExtractString(t *testing.T) {
	pe := NewParameterExtractor()

	tests := []struct {
		name        string
		params      map[string]interface{}
		key         string
		expected    string
		expectError bool
	}{
		{
			name: "Valid string",
			params: map[string]interface{}{
				"value": "test",
			},
			key:         "value",
			expected:    "test",
			expectError: false,
		},
		{
			name: "Valid int",
			params: map[string]interface{}{
				"value": 42,
			},
			key:         "value",
			expected:    "42",
			expectError: false,
		},
		{
			name: "Valid float64",
			params: map[string]interface{}{
				"value": 42.5,
			},
			key:         "value",
			expected:    "42.5",
			expectError: false,
		},
		{
			name: "Valid bool true",
			params: map[string]interface{}{
				"value": true,
			},
			key:         "value",
			expected:    "true",
			expectError: false,
		},
		{
			name: "Valid bool false",
			params: map[string]interface{}{
				"value": false,
			},
			key:         "value",
			expected:    "false",
			expectError: false,
		},
		{
			name: "Missing parameter",
			params: map[string]interface{}{
				"other": "test",
			},
			key:         "value",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pe.ExtractString(tt.params, tt.key)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPositionValue_JSON(t *testing.T) {
	// Test marshaling
	pv := NewPositionValue(123.45, PositionUnitMM)
	data, err := json.Marshal(pv)
	if err != nil {
		t.Errorf("Unexpected error marshaling: %v", err)
	}

	// Test unmarshaling
	var pv2 PositionValue
	err = json.Unmarshal(data, &pv2)
	if err != nil {
		t.Errorf("Unexpected error unmarshaling: %v", err)
	}

	if pv2.Value != pv.Value {
		t.Errorf("Expected value %f, got %f", pv.Value, pv2.Value)
	}

	if pv2.Unit != pv.Unit {
		t.Errorf("Expected unit %v, got %v", pv.Unit, pv2.Unit)
	}
}

func TestMotionParameters_Structure(t *testing.T) {
	mp := MotionParameters{
		Position:     NewPositionValue(100.0, PositionUnitMM),
		Velocity:     NewVelocityValue(50.0, VelocityUnitMMS),
		Acceleration: NewAccelerationValue(25.0, AccelerationUnitMMS2),
		Deceleration: NewAccelerationValue(25.0, AccelerationUnitMMS2),
		Jerk:         NewJerkValue(10.0, JerkUnitMMS3),
		Timeout:      NewTimeValue(5000.0, TimeUnitMS),
		Tolerance:    NewPositionValue(0.1, PositionUnitMM),
	}

	if mp.Position.Value != 100.0 {
		t.Errorf("Expected position 100.0, got %f", mp.Position.Value)
	}

	if mp.Velocity.Value != 50.0 {
		t.Errorf("Expected velocity 50.0, got %f", mp.Velocity.Value)
	}

	if mp.Acceleration.Value != 25.0 {
		t.Errorf("Expected acceleration 25.0, got %f", mp.Acceleration.Value)
	}

	if mp.Deceleration.Value != 25.0 {
		t.Errorf("Expected deceleration 25.0, got %f", mp.Deceleration.Value)
	}

	if mp.Jerk.Value != 10.0 {
		t.Errorf("Expected jerk 10.0, got %f", mp.Jerk.Value)
	}

	if mp.Timeout.Value != 5000.0 {
		t.Errorf("Expected timeout 5000.0, got %f", mp.Timeout.Value)
	}

	if mp.Tolerance.Value != 0.1 {
		t.Errorf("Expected tolerance 0.1, got %f", mp.Tolerance.Value)
	}
}

func TestIOParameters_Structure(t *testing.T) {
	io := IOParameters{
		OutputNumber: 1,
		InputNumber:  2,
		State:        true,
		Value:        3.14,
		Tolerance:    0.01,
		Timeout:      NewTimeValue(1000.0, TimeUnitMS),
	}

	if io.OutputNumber != 1 {
		t.Errorf("Expected output number 1, got %d", io.OutputNumber)
	}

	if io.InputNumber != 2 {
		t.Errorf("Expected input number 2, got %d", io.InputNumber)
	}

	if io.State != true {
		t.Errorf("Expected state true, got %t", io.State)
	}

	if io.Value != 3.14 {
		t.Errorf("Expected value 3.14, got %f", io.Value)
	}

	if io.Tolerance != 0.01 {
		t.Errorf("Expected tolerance 0.01, got %f", io.Tolerance)
	}

	if io.Timeout.Value != 1000.0 {
		t.Errorf("Expected timeout 1000.0, got %f", io.Timeout.Value)
	}
}

func TestLoopParameters_Structure(t *testing.T) {
	condition := Condition{
		Type:      CondDigitalInput,
		Parameter: "digital_input_1",
		Operator:  OpEqual,
		Value:     true,
	}

	lp := LoopParameters{
		CounterVariable: "loop_count",
		MaxIterations:   10,
		Condition:       condition,
		NextCommand:     5,
	}

	if lp.CounterVariable != "loop_count" {
		t.Errorf("Expected counter variable 'loop_count', got %s", lp.CounterVariable)
	}

	if lp.MaxIterations != 10 {
		t.Errorf("Expected max iterations 10, got %d", lp.MaxIterations)
	}

	if lp.Condition.Type != CondDigitalInput {
		t.Errorf("Expected condition type CondDigitalInput, got %v", lp.Condition.Type)
	}

	if lp.NextCommand != 5 {
		t.Errorf("Expected next command 5, got %d", lp.NextCommand)
	}
}

func TestJumpParameters_Structure(t *testing.T) {
	condition := Condition{
		Type:      CondPosition,
		Parameter: "position",
		Operator:  OpGreaterThan,
		Value:     NewPositionValue(100.0, PositionUnitMM),
	}

	jp := JumpParameters{
		Condition:     condition,
		TargetCommand: 10,
		FalseCommand:  20,
	}

	if jp.Condition.Type != CondPosition {
		t.Errorf("Expected condition type CondPosition, got %v", jp.Condition.Type)
	}

	if jp.TargetCommand != 10 {
		t.Errorf("Expected target command 10, got %d", jp.TargetCommand)
	}

	if jp.FalseCommand != 20 {
		t.Errorf("Expected false command 20, got %d", jp.FalseCommand)
	}
}

func TestSystemParameters_Structure(t *testing.T) {
	sp := SystemParameters{
		HomingMethod: 1,
		ResetType:    2,
		ConfigSlot:   3,
		Velocity:     NewVelocityValue(25.0, VelocityUnitMMS),
		Acceleration: NewAccelerationValue(50.0, AccelerationUnitMMS2),
		Timeout:      NewTimeValue(30000.0, TimeUnitMS),
	}

	if sp.HomingMethod != 1 {
		t.Errorf("Expected homing method 1, got %d", sp.HomingMethod)
	}

	if sp.ResetType != 2 {
		t.Errorf("Expected reset type 2, got %d", sp.ResetType)
	}

	if sp.ConfigSlot != 3 {
		t.Errorf("Expected config slot 3, got %d", sp.ConfigSlot)
	}

	if sp.Velocity.Value != 25.0 {
		t.Errorf("Expected velocity 25.0, got %f", sp.Velocity.Value)
	}

	if sp.Acceleration.Value != 50.0 {
		t.Errorf("Expected acceleration 50.0, got %f", sp.Acceleration.Value)
	}

	if sp.Timeout.Value != 30000.0 {
		t.Errorf("Expected timeout 30000.0, got %f", sp.Timeout.Value)
	}
}

func TestForceControlParameters_Structure(t *testing.T) {
	fcp := ForceControlParameters{
		ForceSetpoint:  NewForceValue(500.0, ForceUnitN),
		ForceLimit:     NewForceValue(1000.0, ForceUnitN),
		PositionLimit:  NewPositionValue(200.0, PositionUnitMM),
		TransitionTime: NewTimeValue(100.0, TimeUnitMS),
	}

	if fcp.ForceSetpoint.Value != 500.0 {
		t.Errorf("Expected force setpoint 500.0, got %f", fcp.ForceSetpoint.Value)
	}

	if fcp.ForceLimit.Value != 1000.0 {
		t.Errorf("Expected force limit 1000.0, got %f", fcp.ForceLimit.Value)
	}

	if fcp.PositionLimit.Value != 200.0 {
		t.Errorf("Expected position limit 200.0, got %f", fcp.PositionLimit.Value)
	}

	if fcp.TransitionTime.Value != 100.0 {
		t.Errorf("Expected transition time 100.0, got %f", fcp.TransitionTime.Value)
	}
}

func TestDataAcquisitionParameters_Structure(t *testing.T) {
	dap := DataAcquisitionParameters{
		SampleRate: 1000,
		NumSamples: 10000,
		Channels:   []int{1, 2, 3},
		SaveData:   true,
		Filename:   "test_data.csv",
		Format:     "CSV",
	}

	if dap.SampleRate != 1000 {
		t.Errorf("Expected sample rate 1000, got %d", dap.SampleRate)
	}

	if dap.NumSamples != 10000 {
		t.Errorf("Expected num samples 10000, got %d", dap.NumSamples)
	}

	if len(dap.Channels) != 3 {
		t.Errorf("Expected 3 channels, got %d", len(dap.Channels))
	}

	if dap.Channels[0] != 1 || dap.Channels[1] != 2 || dap.Channels[2] != 3 {
		t.Errorf("Expected channels [1, 2, 3], got %v", dap.Channels)
	}

	if dap.SaveData != true {
		t.Errorf("Expected save data true, got %t", dap.SaveData)
	}

	if dap.Filename != "test_data.csv" {
		t.Errorf("Expected filename 'test_data.csv', got %s", dap.Filename)
	}

	if dap.Format != "CSV" {
		t.Errorf("Expected format 'CSV', got %s", dap.Format)
	}
}