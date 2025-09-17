package types

import (
	"math"
	"testing"
)

func TestUnitConverter_NewUnitConverter(t *testing.T) {
	uc := NewUnitConverter()

	if uc.positionScalingFactor != 1000.0 {
		t.Errorf("Expected default position scaling factor 1000.0, got %f", uc.positionScalingFactor)
	}

	if uc.forceScalingFactor != 100.0 {
		t.Errorf("Expected default force scaling factor 100.0, got %f", uc.forceScalingFactor)
	}
}

func TestUnitConverter_NewUnitConverterWithFactors(t *testing.T) {
	uc := NewUnitConverterWithFactors(2000.0, 200.0)

	if uc.positionScalingFactor != 2000.0 {
		t.Errorf("Expected position scaling factor 2000.0, got %f", uc.positionScalingFactor)
	}

	if uc.forceScalingFactor != 200.0 {
		t.Errorf("Expected force scaling factor 200.0, got %f", uc.forceScalingFactor)
	}
}

func TestUnitConverter_ConvertPosition(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     PositionUnit
		to       PositionUnit
		expected float64
	}{
		{"Same unit MM", 100.0, PositionUnitMM, PositionUnitMM, 100.0},
		{"Same unit Counts", 1000.0, PositionUnitCounts, PositionUnitCounts, 1000.0},
		{"MM to Counts", 100.0, PositionUnitMM, PositionUnitCounts, 100000.0},
		{"Counts to MM", 100000.0, PositionUnitCounts, PositionUnitMM, 100.0},
		{"Zero value", 0.0, PositionUnitMM, PositionUnitCounts, 0.0},
		{"Negative value", -50.0, PositionUnitMM, PositionUnitCounts, -50000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertPosition(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertVelocity(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     VelocityUnit
		to       VelocityUnit
		expected float64
	}{
		{"Same unit MMS", 50.0, VelocityUnitMMS, VelocityUnitMMS, 50.0},
		{"Same unit CountsS", 50000.0, VelocityUnitCountsS, VelocityUnitCountsS, 50000.0},
		{"MMS to CountsS", 50.0, VelocityUnitMMS, VelocityUnitCountsS, 50000.0},
		{"CountsS to MMS", 50000.0, VelocityUnitCountsS, VelocityUnitMMS, 50.0},
		{"Zero value", 0.0, VelocityUnitMMS, VelocityUnitCountsS, 0.0},
		{"Negative value", -25.0, VelocityUnitMMS, VelocityUnitCountsS, -25000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertVelocity(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertAcceleration(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     AccelerationUnit
		to       AccelerationUnit
		expected float64
	}{
		{"Same unit MMS2", 25.0, AccelerationUnitMMS2, AccelerationUnitMMS2, 25.0},
		{"Same unit CountsS2", 25000.0, AccelerationUnitCountsS2, AccelerationUnitCountsS2, 25000.0},
		{"MMS2 to CountsS2", 25.0, AccelerationUnitMMS2, AccelerationUnitCountsS2, 25000.0},
		{"CountsS2 to MMS2", 25000.0, AccelerationUnitCountsS2, AccelerationUnitMMS2, 25.0},
		{"Zero value", 0.0, AccelerationUnitMMS2, AccelerationUnitCountsS2, 0.0},
		{"Negative value", -12.5, AccelerationUnitMMS2, AccelerationUnitCountsS2, -12500.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertAcceleration(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertJerk(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     JerkUnit
		to       JerkUnit
		expected float64
	}{
		{"Same unit MMS3", 10.0, JerkUnitMMS3, JerkUnitMMS3, 10.0},
		{"Same unit CountsS3", 10000000.0, JerkUnitCountsS3, JerkUnitCountsS3, 10000000.0},
		{"MMS3 to CountsS3", 10.0, JerkUnitMMS3, JerkUnitCountsS3, 10000000000.0},
		{"CountsS3 to MMS3", 10000000000.0, JerkUnitCountsS3, JerkUnitMMS3, 10.0},
		{"Zero value", 0.0, JerkUnitMMS3, JerkUnitCountsS3, 0.0},
		{"Negative value", -5.0, JerkUnitMMS3, JerkUnitCountsS3, -5000000000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertJerk(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertForce(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     ForceUnit
		to       ForceUnit
		expected float64
	}{
		{"Same unit N", 1000.0, ForceUnitN, ForceUnitN, 1000.0},
		{"Same unit Counts", 100000.0, ForceUnitCounts, ForceUnitCounts, 100000.0},
		{"N to Counts", 1000.0, ForceUnitN, ForceUnitCounts, 100000.0},
		{"Counts to N", 100000.0, ForceUnitCounts, ForceUnitN, 1000.0},
		{"Zero value", 0.0, ForceUnitN, ForceUnitCounts, 0.0},
		{"Negative value", -500.0, ForceUnitN, ForceUnitCounts, -50000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertForce(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertTime(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name     string
		value    float64
		from     TimeUnit
		to       TimeUnit
		expected float64
	}{
		{"Same unit MS", 1000.0, TimeUnitMS, TimeUnitMS, 1000.0},
		{"Same unit S", 1.0, TimeUnitS, TimeUnitS, 1.0},
		{"MS to S", 1000.0, TimeUnitMS, TimeUnitS, 1.0},
		{"S to MS", 1.0, TimeUnitS, TimeUnitMS, 1000.0},
		{"Zero value", 0.0, TimeUnitMS, TimeUnitS, 0.0},
		{"Negative value", -500.0, TimeUnitMS, TimeUnitS, -0.5},
		{"Fractional seconds", 0.5, TimeUnitS, TimeUnitMS, 500.0},
		{"Large milliseconds", 5000.0, TimeUnitMS, TimeUnitS, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.ConvertTime(tt.value, tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitConverter_ConvertPositionValue(t *testing.T) {
	uc := NewUnitConverter()

	pv := NewPositionValue(100.0, PositionUnitMM)
	converted := uc.ConvertPositionValue(pv, PositionUnitCounts)

	if converted.Value != 100000.0 {
		t.Errorf("Expected value 100000.0, got %f", converted.Value)
	}

	if converted.Unit != PositionUnitCounts {
		t.Errorf("Expected unit PositionUnitCounts, got %v", converted.Unit)
	}

	// Test same unit conversion
	sameUnit := uc.ConvertPositionValue(pv, PositionUnitMM)
	if sameUnit.Value != pv.Value {
		t.Errorf("Expected same value %f, got %f", pv.Value, sameUnit.Value)
	}

	if sameUnit.Unit != pv.Unit {
		t.Errorf("Expected same unit %v, got %v", pv.Unit, sameUnit.Unit)
	}
}

func TestUnitConverter_ConvertVelocityValue(t *testing.T) {
	uc := NewUnitConverter()

	vv := NewVelocityValue(50.0, VelocityUnitMMS)
	converted := uc.ConvertVelocityValue(vv, VelocityUnitCountsS)

	if converted.Value != 50000.0 {
		t.Errorf("Expected value 50000.0, got %f", converted.Value)
	}

	if converted.Unit != VelocityUnitCountsS {
		t.Errorf("Expected unit VelocityUnitCountsS, got %v", converted.Unit)
	}
}

func TestUnitConverter_ConvertAccelerationValue(t *testing.T) {
	uc := NewUnitConverter()

	av := NewAccelerationValue(25.0, AccelerationUnitMMS2)
	converted := uc.ConvertAccelerationValue(av, AccelerationUnitCountsS2)

	if converted.Value != 25000.0 {
		t.Errorf("Expected value 25000.0, got %f", converted.Value)
	}

	if converted.Unit != AccelerationUnitCountsS2 {
		t.Errorf("Expected unit AccelerationUnitCountsS2, got %v", converted.Unit)
	}
}

func TestUnitConverter_ConvertJerkValue(t *testing.T) {
	uc := NewUnitConverter()

	jv := NewJerkValue(10.0, JerkUnitMMS3)
	converted := uc.ConvertJerkValue(jv, JerkUnitCountsS3)

	if converted.Value != 10000000000.0 {
		t.Errorf("Expected value 10000000000.0, got %f", converted.Value)
	}

	if converted.Unit != JerkUnitCountsS3 {
		t.Errorf("Expected unit JerkUnitCountsS3, got %v", converted.Unit)
	}
}

func TestUnitConverter_ConvertForceValue(t *testing.T) {
	uc := NewUnitConverter()

	fv := NewForceValue(1000.0, ForceUnitN)
	converted := uc.ConvertForceValue(fv, ForceUnitCounts)

	if converted.Value != 100000.0 {
		t.Errorf("Expected value 100000.0, got %f", converted.Value)
	}

	if converted.Unit != ForceUnitCounts {
		t.Errorf("Expected unit ForceUnitCounts, got %v", converted.Unit)
	}
}

func TestUnitConverter_ConvertTimeValue(t *testing.T) {
	uc := NewUnitConverter()

	tv := NewTimeValue(1000.0, TimeUnitMS)
	converted := uc.ConvertTimeValue(tv, TimeUnitS)

	if converted.Value != 1.0 {
		t.Errorf("Expected value 1.0, got %f", converted.Value)
	}

	if converted.Unit != TimeUnitS {
		t.Errorf("Expected unit TimeUnitS, got %v", converted.Unit)
	}
}

func TestUnitConverter_SetPositionScalingFactor(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name        string
		factor      float64
		expectError bool
	}{
		{"Valid positive factor", 2000.0, false},
		{"Valid small factor", 0.1, false},
		{"Zero factor", 0.0, true},
		{"Negative factor", -100.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.SetPositionScalingFactor(tt.factor)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if uc.positionScalingFactor != tt.factor {
					t.Errorf("Expected scaling factor %f, got %f", tt.factor, uc.positionScalingFactor)
				}
			}
		})
	}
}

func TestUnitConverter_SetForceScalingFactor(t *testing.T) {
	uc := NewUnitConverter()

	tests := []struct {
		name        string
		factor      float64
		expectError bool
	}{
		{"Valid positive factor", 200.0, false},
		{"Valid small factor", 0.5, false},
		{"Zero factor", 0.0, true},
		{"Negative factor", -50.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.SetForceScalingFactor(tt.factor)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if uc.forceScalingFactor != tt.factor {
					t.Errorf("Expected scaling factor %f, got %f", tt.factor, uc.forceScalingFactor)
				}
			}
		})
	}
}

func TestUnitConverter_GetPositionScalingFactor(t *testing.T) {
	uc := NewUnitConverter()
	uc.positionScalingFactor = 1500.0

	if uc.GetPositionScalingFactor() != 1500.0 {
		t.Errorf("Expected 1500.0, got %f", uc.GetPositionScalingFactor())
	}
}

func TestUnitConverter_GetForceScalingFactor(t *testing.T) {
	uc := NewUnitConverter()
	uc.forceScalingFactor = 150.0

	if uc.GetForceScalingFactor() != 150.0 {
		t.Errorf("Expected 150.0, got %f", uc.GetForceScalingFactor())
	}
}

func TestUnitConverter_ValidateScalingFactors(t *testing.T) {
	tests := []struct {
		name        string
		posFactor   float64
		forceFactor float64
		expectError bool
	}{
		{"Valid factors", 1000.0, 100.0, false},
		{"Zero position factor", 0.0, 100.0, true},
		{"Zero force factor", 1000.0, 0.0, true},
		{"Negative position factor", -1000.0, 100.0, true},
		{"Negative force factor", 1000.0, -100.0, true},
		{"Inf position factor", math.Inf(1), 100.0, true},
		{"NaN position factor", math.NaN(), 100.0, true},
		{"Inf force factor", 1000.0, math.Inf(1), true},
		{"NaN force factor", 1000.0, math.NaN(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUnitConverterWithFactors(tt.posFactor, tt.forceFactor)
			err := uc.ValidateScalingFactors()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUnitConversionError(t *testing.T) {
	err := NewUnitConversionError("mm", "counts", 100.0, "test error")
	
	if err.FromUnit != "mm" {
		t.Errorf("Expected FromUnit 'mm', got %s", err.FromUnit)
	}
	
	if err.ToUnit != "counts" {
		t.Errorf("Expected ToUnit 'counts', got %s", err.ToUnit)
	}
	
	if err.Value != 100.0 {
		t.Errorf("Expected Value 100.0, got %f", err.Value)
	}
	
	if err.Message != "test error" {
		t.Errorf("Expected Message 'test error', got %s", err.Message)
	}
	
	expectedErrorMsg := "unit conversion error from mm to counts for value 100.000000: test error"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestUnitSystem_DefaultUnitSystem(t *testing.T) {
	us := DefaultUnitSystem()

	if us.PositionUnit != PositionUnitMM {
		t.Errorf("Expected PositionUnitMM, got %v", us.PositionUnit)
	}

	if us.VelocityUnit != VelocityUnitMMS {
		t.Errorf("Expected VelocityUnitMMS, got %v", us.VelocityUnit)
	}

	if us.AccelerationUnit != AccelerationUnitMMS2 {
		t.Errorf("Expected AccelerationUnitMMS2, got %v", us.AccelerationUnit)
	}

	if us.JerkUnit != JerkUnitMMS3 {
		t.Errorf("Expected JerkUnitMMS3, got %v", us.JerkUnit)
	}

	if us.ForceUnit != ForceUnitN {
		t.Errorf("Expected ForceUnitN, got %v", us.ForceUnit)
	}

	if us.TimeUnit != TimeUnitMS {
		t.Errorf("Expected TimeUnitMS, got %v", us.TimeUnit)
	}
}

func TestUnitSystem_ImperialUnitSystem(t *testing.T) {
	us := ImperialUnitSystem()

	// Note: Currently returns DefaultUnitSystem, but this tests the function exists
	if us == nil {
		t.Error("Expected non-nil unit system")
	}
}

func TestUnitConverterFactory_NewUnitConverterFactory(t *testing.T) {
	factory := NewUnitConverterFactory()

	if factory.driveTypeScalingFactors == nil {
		t.Error("Expected non-nil driveTypeScalingFactors map")
	}

	if len(factory.driveTypeScalingFactors) != 0 {
		t.Errorf("Expected empty driveTypeScalingFactors map, got %d entries", len(factory.driveTypeScalingFactors))
	}
}

func TestUnitConverterFactory_RegisterDriveType(t *testing.T) {
	factory := NewUnitConverterFactory()

	factors := map[string]float64{
		"position": 2000.0,
		"force":    200.0,
	}

	factory.RegisterDriveType("C1250-EC", factors)

	if len(factory.driveTypeScalingFactors) != 1 {
		t.Errorf("Expected 1 drive type, got %d", len(factory.driveTypeScalingFactors))
	}

	registeredFactors, exists := factory.driveTypeScalingFactors["C1250-EC"]
	if !exists {
		t.Error("Expected drive type 'C1250-EC' to be registered")
	}

	if registeredFactors["position"] != 2000.0 {
		t.Errorf("Expected position factor 2000.0, got %f", registeredFactors["position"])
	}

	if registeredFactors["force"] != 200.0 {
		t.Errorf("Expected force factor 200.0, got %f", registeredFactors["force"])
	}
}

func TestUnitConverterFactory_CreateUnitConverter(t *testing.T) {
	factory := NewUnitConverterFactory()

	// Test with registered drive type
	factors := map[string]float64{
		"position": 2000.0,
		"force":    200.0,
	}
	factory.RegisterDriveType("C1250-EC", factors)

	converter, err := factory.CreateUnitConverter("C1250-EC")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if converter.positionScalingFactor != 2000.0 {
		t.Errorf("Expected position scaling factor 2000.0, got %f", converter.positionScalingFactor)
	}

	if converter.forceScalingFactor != 200.0 {
		t.Errorf("Expected force scaling factor 200.0, got %f", converter.forceScalingFactor)
	}

	// Test with unregistered drive type (should use defaults)
	converter2, err := factory.CreateUnitConverter("Unknown")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if converter2.positionScalingFactor != 1000.0 {
		t.Errorf("Expected default position scaling factor 1000.0, got %f", converter2.positionScalingFactor)
	}

	if converter2.forceScalingFactor != 100.0 {
		t.Errorf("Expected default force scaling factor 100.0, got %f", converter2.forceScalingFactor)
	}
}

func TestUnitConverterFactory_GetSupportedDriveTypes(t *testing.T) {
	factory := NewUnitConverterFactory()

	// Initially empty
	types := factory.GetSupportedDriveTypes()
	if len(types) != 0 {
		t.Errorf("Expected empty drive types, got %d", len(types))
	}

	// Register some drive types
	factory.RegisterDriveType("C1250-EC", map[string]float64{"position": 2000.0, "force": 200.0})
	factory.RegisterDriveType("C1200-EC", map[string]float64{"position": 1500.0, "force": 150.0})

	types = factory.GetSupportedDriveTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 drive types, got %d", len(types))
	}

	// Check that both types are present
	foundC1250 := false
	foundC1200 := false
	for _, driveType := range types {
		if driveType == "C1250-EC" {
			foundC1250 = true
		}
		if driveType == "C1200-EC" {
			foundC1200 = true
		}
	}

	if !foundC1250 {
		t.Error("Expected 'C1250-EC' to be in supported drive types")
	}

	if !foundC1200 {
		t.Error("Expected 'C1200-EC' to be in supported drive types")
	}
}

func TestUnitConverter_ConvertToUnitSystem(t *testing.T) {
	uc := NewUnitConverter()

	// Create a test command
	cmd := &Command{
		Type: CmdMoveAbsolute,
		Parameters: map[string]interface{}{
			"position": map[string]interface{}{
				"value": 100.0,
				"unit":  "mm",
			},
			"velocity": map[string]interface{}{
				"value": 50.0,
				"unit":  "mm/s",
			},
		},
	}

	// Test conversion to counts unit system
	targetSystem := &UnitSystem{
		PositionUnit:     PositionUnitCounts,
		VelocityUnit:     VelocityUnitCountsS,
		AccelerationUnit: AccelerationUnitCountsS2,
		JerkUnit:         JerkUnitCountsS3,
		ForceUnit:        ForceUnitCounts,
		TimeUnit:         TimeUnitMS,
	}

	err := uc.ConvertToUnitSystem(cmd, targetSystem)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: The actual conversion logic would need to be implemented
	// This test verifies the function exists and can be called
}

func TestUnitConverter_ConvertMotionCommandToUnitSystem(t *testing.T) {
	uc := NewUnitConverter()

	cmd := &Command{
		Type: CmdMoveAbsolute,
		Parameters: map[string]interface{}{
			"position": map[string]interface{}{
				"value": 100.0,
				"unit":  "mm",
			},
		},
	}

	targetSystem := &UnitSystem{
		PositionUnit: PositionUnitCounts,
	}

	err := uc.convertMotionCommandToUnitSystem(cmd, targetSystem)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: The actual conversion logic would need to be implemented
	// This test verifies the function exists and can be called
}

func TestUnitConverter_ConvertJogCommandToUnitSystem(t *testing.T) {
	uc := NewUnitConverter()

	cmd := &Command{
		Type: CmdJog,
		Parameters: map[string]interface{}{
			"velocity": map[string]interface{}{
				"value": 50.0,
				"unit":  "mm/s",
			},
		},
	}

	targetSystem := &UnitSystem{
		VelocityUnit: VelocityUnitCountsS,
	}

	err := uc.convertJogCommandToUnitSystem(cmd, targetSystem)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: The actual conversion logic would need to be implemented
	// This test verifies the function exists and can be called
}

func TestUnitConverter_ConvertWaitCommandToUnitSystem(t *testing.T) {
	uc := NewUnitConverter()

	cmd := &Command{
		Type: CmdWaitPosition,
		Parameters: map[string]interface{}{
			"position": map[string]interface{}{
				"value": 100.0,
				"unit":  "mm",
			},
			"timeout": map[string]interface{}{
				"value": 5000.0,
				"unit":  "ms",
			},
		},
	}

	targetSystem := &UnitSystem{
		PositionUnit: PositionUnitCounts,
		TimeUnit:     TimeUnitS,
	}

	err := uc.convertWaitCommandToUnitSystem(cmd, targetSystem)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Note: The actual conversion logic would need to be implemented
	// This test verifies the function exists and can be called
}