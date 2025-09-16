package stage_linmot_drive

import (
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive/python_port"
)

func TestDriveConfig(t *testing.T) {
	config := DriveConfig{
		AdapterID:           "eth0",
		NumDevices:          1,
		CycleTime:           0.050,
		NumMonitoring:       4,
		NumParameter:        4,
		ActivateLMDriveData: false,
		MpLogging:           20,
	}

	if config.AdapterID != "eth0" {
		t.Errorf("Expected adapter ID 'eth0', got '%s'", config.AdapterID)
	}

	if config.NumDevices != 1 {
		t.Errorf("Expected 1 device, got %d", config.NumDevices)
	}

	if config.CycleTime != 0.050 {
		t.Errorf("Expected cycle time 0.050, got %f", config.CycleTime)
	}
}

func TestMotionHeaders(t *testing.T) {
	// Test motion header constants
	headers := []python_port.MotionHeader{
		python_port.AbsoluteVAI,
		python_port.RelativeVAI,
		python_port.AbsoluteVAJI,
		python_port.RelativeVAJI,
		python_port.IncrActPosRstI,
		python_port.AbsoluteSin,
		python_port.RelativeSin,
	}

	expectedHeaders := []string{
		"Absolute_VAI",
		"Relative_VAI",
		"Absolute_VAJI",
		"Relative_VAJI",
		"Incr_Act_Pos_RstI",
		"Absolute_Sin",
		"Relative_Sin",
	}

	for i, header := range headers {
		if string(header) != expectedHeaders[i] {
			t.Errorf("Expected header '%s', got '%s'", expectedHeaders[i], string(header))
		}
	}
}

func TestForceControlHeaders(t *testing.T) {
	// Test force control header constants
	headers := []python_port.ForceControlHeader{
		python_port.VAI_GoToPosWithHigherForceCtrlLimit,
		python_port.VAI_GoToPosWithLowerForceCtrlLimit,
		python_port.VAI_IncActPosWithHigherForceCtrlLimit,
		python_port.VAI_IncActPosWithLowerForceCtrlLimit,
		python_port.VAI_GoToPosFromActPosAndResetForceControl,
		python_port.VAI_IncrementActPosAndResetForceControl,
	}

	expectedHeaders := []string{
		"VAI Go To Pos With Higher Force Ctrl Limit and Target Force",
		"VAI Go To Pos With Lower Force Ctrl Limit and Target Force",
		"VAI Inc Act Pos With Higher Force Ctrl Limit and Target Force",
		"VAI Inc Act Pos With Lower Force Ctrl Limit and Target Force",
		"VAI Go To Pos From Act Pos And Reset Force Control Set I",
		"VAI Increment Act Pos And Reset Force Control Set I",
	}

	for i, header := range headers {
		if string(header) != expectedHeaders[i] {
			t.Errorf("Expected header '%s', got '%s'", expectedHeaders[i], string(header))
		}
	}
}

func TestConfigHeaders(t *testing.T) {
	// Test configuration header constants
	headers := []python_port.ConfigHeader{
		python_port.ReadValueROM,
		python_port.ReadValueRAM,
		python_port.WriteValueROM,
		python_port.WriteValueRAM,
		python_port.WriteValueRAMAndROM,
	}

	expectedHeaders := []string{
		"Read_Value_ROM",
		"Read_Value_RAM",
		"Write_Value_ROM",
		"Write_Value_RAM",
		"Write_Value_RAM_and_ROM",
	}

	for i, header := range headers {
		if string(header) != expectedHeaders[i] {
			t.Errorf("Expected header '%s', got '%s'", expectedHeaders[i], string(header))
		}
	}
}

func TestDriveState(t *testing.T) {
	config := DriveConfig{
		AdapterID:           "eth0",
		NumDevices:          1,
		CycleTime:           0.050,
		NumMonitoring:       4,
		NumParameter:        4,
		ActivateLMDriveData: false,
		MpLogging:           20,
	}

	// Note: This test will fail without actual Python environment and hardware
	// It's included to show the expected API usage
	drive, err := NewDrive(config)
	if err != nil {
		t.Logf("Expected error without Python environment: %v", err)
		return
	}
	defer drive.Close()

	// Test initial state
	if drive.IsActive() {
		t.Error("Drive should not be active initially")
	}

	// Test starting (will fail without hardware)
	err = drive.Start()
	if err != nil {
		t.Logf("Expected error starting without hardware: %v", err)
	}

	// Test stopping
	err = drive.Stop()
	if err != nil {
		t.Logf("Error stopping drive: %v", err)
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test UnsignedToSigned16Bit
	testCases := []struct {
		input    uint16
		expected int16
	}{
		{0, 0},
		{100, 100},
		{32767, 32767},
		{32768, -32768},
		{65535, -1},
	}

	for _, tc := range testCases {
		result := python_port.UnsignedToSigned16Bit(tc.input)
		if result != tc.expected {
			t.Errorf("UnsignedToSigned16Bit(%d) = %d, expected %d", tc.input, result, tc.expected)
		}
	}

	// Test IEEE754BitsToFloat
	floatTests := []struct {
		value     uint32
		precision string
		expected  float64
	}{
		{0x3F800000, "single", 1.0}, // 1.0 in IEEE 754 single precision
		{0x40000000, "single", 2.0}, // 2.0 in IEEE 754 single precision
	}

	for _, tc := range floatTests {
		result, err := python_port.IEEE754BitsToFloat(tc.value, tc.precision)
		if err != nil {
			t.Errorf("IEEE754BitsToFloat(%d, %s) error: %v", tc.value, tc.precision, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("IEEE754BitsToFloat(%d, %s) = %f, expected %f", tc.value, tc.precision, result, tc.expected)
		}
	}
}

func TestDriveDataStructures(t *testing.T) {
	// Test DriveConfig structure
	config := DriveConfig{
		IsRotaryMotor:       false,
		PosScaleNumerator:   10000.0,
		PosScaleDenominator: 1.0,
		UnitScale:           10000.0,
		ModuloFactor:        360000,
		FcForceScale:        0.1,
		FcTorqueScale:       0.00057295779513082,
		DriveName:           "LMDrive",
		DriveType:           "0",
	}

	if config.UnitScale != 10000.0 {
		t.Errorf("Expected UnitScale 10000.0, got %f", config.UnitScale)
	}

	// Test DriveStatus structure
	status := python_port.DriveStatus{
		OperationEnabled:   true,
		SwitchOnLocked:     false,
		Homed:              true,
		MotionActive:       false,
		Jogging:            false,
		Warning:            false,
		Error:              false,
		ErrorCode:          0,
		DemandPosition:     0.0,
		ActualPosition:     0.0,
		DifferencePosition: 0.0,
		ActualCurrent:      0.0,
		NrOfRevolutions:    0,
	}

	if !status.OperationEnabled {
		t.Error("Expected OperationEnabled to be true")
	}

	if !status.Homed {
		t.Error("Expected Homed to be true")
	}

	if status.MotionActive {
		t.Error("Expected MotionActive to be false")
	}
}

func TestDriveInputsOutputs(t *testing.T) {
	// Test DriveInputs structure
	inputs := python_port.DriveInputs{
		StateVar:     0x0000,
		StatusWord:   0x0000,
		WarnWord:     0x0000,
		DemandPos:    0,
		ActualPos:    0,
		DemandCurr:   0,
		CfgStatus:    0x0000,
		CfgIndexIn:   0x0000,
		CfgValueIn:   0x00000000,
		MonChannels:  make([]int32, 4),
	}

	if len(inputs.MonChannels) != 4 {
		t.Errorf("Expected 4 monitoring channels, got %d", len(inputs.MonChannels))
	}

	// Test DriveOutputs structure
	outputs := python_port.DriveOutputs{
		ControlWord:  0x003E,
		McHeader:     0x0000,
		McParaWords:  [10]uint16{},
		CfgControl:   0x0000,
		CfgIndexOut:  0x0000,
		CfgValueOut:  0x00000000,
		ParChannels:  make([]uint16, 4),
	}

	if outputs.ControlWord != 0x003E {
		t.Errorf("Expected ControlWord 0x003E, got 0x%04X", outputs.ControlWord)
	}

	if len(outputs.ParChannels) != 4 {
		t.Errorf("Expected 4 parameter channels, got %d", len(outputs.ParChannels))
	}
}

// Benchmark tests
func BenchmarkUnsignedToSigned16Bit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		python_port.UnsignedToSigned16Bit(uint16(i % 65536))
	}
}

func BenchmarkIEEE754BitsToFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		python_port.IEEE754BitsToFloat(uint32(i%1000000), "single")
	}
}

// Example test showing how to use the drive (will fail without hardware)
func ExampleDrive() {
	config := DriveConfig{
		AdapterID:           "eth0",
		NumDevices:          1,
		CycleTime:           0.050,
		NumMonitoring:       4,
		NumParameter:        4,
		ActivateLMDriveData: false,
		MpLogging:           20,
	}

	drive, err := NewDrive(config)
	if err != nil {
		// Handle error
		return
	}
	defer drive.Close()

	// Start communication
	if err := drive.Start(); err != nil {
		// Handle error
		return
	}
	defer drive.Stop()

	// Switch on motor
	if err := drive.SwitchOnMotor(1); err != nil {
		// Handle error
		return
	}

	// Home motor
	if err := drive.HomeMotor(1); err != nil {
		// Handle error
		return
	}

	// Move to position
	countNibble, err := drive.MoveToPosition(1, 5.0, 0.01, 0.1, 0.1, 10000)
	if err != nil {
		// Handle error
		return
	}

	// Wait for motion to complete
	finished, err := drive.WaitForMotionFinished(1, countNibble, 30*time.Second)
	if err != nil {
		// Handle error
		return
	}

	if finished {
		// Motion completed successfully
	}

	// Switch off motor
	if err := drive.SwitchOffMotor(1); err != nil {
		// Handle error
		return
	}
}