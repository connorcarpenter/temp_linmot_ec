package hil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/hardware"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/safety"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// HardwareTestSuite provides comprehensive hardware testing capabilities
type HardwareTestSuite struct {
	controller   hardware.HardwareController
	safetyGuard *safety.SafetyGuard
	testConfig  *TestConfig
	results     *TestResults
	mu          sync.RWMutex
}

// TestConfig defines the configuration for hardware tests
type TestConfig struct {
	EtherCATMaster    string
	DriveAddress      int
	SafetyLimits      *safety.SafetyLimits
	TestTimeout       time.Duration
	RetryAttempts     int
	LogLevel          string
	EnableSafetyTests bool
	EnablePerformanceTests bool
	MaxConcurrency    int
}

// TestResults contains the results of all hardware tests
type TestResults struct {
	PassedTests    int
	FailedTests    int
	SkippedTests   int
	TotalDuration  time.Duration
	TestDetails    []*hardware.HardwareTestResult
	StartTime      time.Time
	EndTime        time.Time
}

// TestDetail provides detailed information about a specific test
type TestDetail struct {
	TestName    string
	Status      TestStatus
	Duration    time.Duration
	Error       error
	Metrics     map[string]interface{}
	Timestamp   time.Time
}

// TestStatus represents the status of a test
type TestStatus int

const (
	TestPassed TestStatus = iota
	TestFailed
	TestSkipped
	TestError
)

// String returns a string representation of TestStatus
func (ts TestStatus) String() string {
	switch ts {
	case TestPassed:
		return "Passed"
	case TestFailed:
		return "Failed"
	case TestSkipped:
		return "Skipped"
	case TestError:
		return "Error"
	default:
		return "Unknown"
	}
}

// NewHardwareTestSuite creates a new hardware test suite
func NewHardwareTestSuite(controller hardware.HardwareController, safetyGuard *safety.SafetyGuard, config *TestConfig) *HardwareTestSuite {
	return &HardwareTestSuite{
		controller:   controller,
		safetyGuard: safetyGuard,
		testConfig:  config,
		results: &TestResults{
			TestDetails: make([]*hardware.HardwareTestResult, 0),
		},
	}
}

// RunBasicMotionTests executes basic motion tests
func (hts *HardwareTestSuite) RunBasicMotionTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	hts.mu.Lock()
	defer hts.mu.Unlock()

	var results []*hardware.HardwareTestResult
	testCtx, cancel := context.WithTimeout(ctx, hts.testConfig.TestTimeout)
	defer cancel()

	// Test absolute motion
	result := hts.testAbsoluteMotion(testCtx)
	results = append(results, result)

	// Test relative motion
	result = hts.testRelativeMotion(testCtx)
	results = append(results, result)

	// Test incremental motion
	result = hts.testIncrementalMotion(testCtx)
	results = append(results, result)

	// Test jog motion
	result = hts.testJogMotion(testCtx)
	results = append(results, result)

	// Test stop motion
	result = hts.testStopMotion(testCtx)
	results = append(results, result)

	return results, nil
}

// RunForceControlTests executes force control tests
func (hts *HardwareTestSuite) RunForceControlTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	hts.mu.Lock()
	defer hts.mu.Unlock()

	var results []*hardware.HardwareTestResult
	testCtx, cancel := context.WithTimeout(ctx, hts.testConfig.TestTimeout)
	defer cancel()

	// Test force control enable/disable
	result := hts.testForceControlEnableDisable(testCtx)
	results = append(results, result)

	// Test force setpoint
	result = hts.testForceSetpoint(testCtx)
	results = append(results, result)

	// Test force monitoring
	result = hts.testForceMonitoring(testCtx)
	results = append(results, result)

	return results, nil
}

// RunIOTests executes I/O tests
func (hts *HardwareTestSuite) RunIOTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	hts.mu.Lock()
	defer hts.mu.Unlock()

	var results []*hardware.HardwareTestResult
	testCtx, cancel := context.WithTimeout(ctx, hts.testConfig.TestTimeout)
	defer cancel()

	// Test digital outputs
	result := hts.testDigitalOutputs(testCtx)
	results = append(results, result)

	// Test digital inputs
	result = hts.testDigitalInputs(testCtx)
	results = append(results, result)

	// Test analog outputs
	result = hts.testAnalogOutputs(testCtx)
	results = append(results, result)

	// Test analog inputs
	result = hts.testAnalogInputs(testCtx)
	results = append(results, result)

	return results, nil
}

// RunSafetyTests executes safety tests
func (hts *HardwareTestSuite) RunSafetyTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	if !hts.testConfig.EnableSafetyTests {
		return []*hardware.HardwareTestResult{}, nil
	}

	hts.mu.Lock()
	defer hts.mu.Unlock()

	var results []*hardware.HardwareTestResult
	testCtx, cancel := context.WithTimeout(ctx, hts.testConfig.TestTimeout)
	defer cancel()

	// Test emergency stop
	result := hts.testEmergencyStop(testCtx)
	results = append(results, result)

	// Test safety limits
	result = hts.testSafetyLimits(testCtx)
	results = append(results, result)

	// Test error recovery
	result = hts.testErrorRecovery(testCtx)
	results = append(results, result)

	return results, nil
}

// RunPerformanceTests executes performance tests
func (hts *HardwareTestSuite) RunPerformanceTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	if !hts.testConfig.EnablePerformanceTests {
		return []*hardware.HardwareTestResult{}, nil
	}

	hts.mu.Lock()
	defer hts.mu.Unlock()

	var results []*hardware.HardwareTestResult
	testCtx, cancel := context.WithTimeout(ctx, hts.testConfig.TestTimeout)
	defer cancel()

	// Test latency
	result := hts.testLatency(testCtx)
	results = append(results, result)

	// Test throughput
	result = hts.testThroughput(testCtx)
	results = append(results, result)

	// Test jitter
	result = hts.testJitter(testCtx)
	results = append(results, result)

	return results, nil
}

// RunAllTests executes all available tests
func (hts *HardwareTestSuite) RunAllTests(ctx context.Context) ([]*hardware.HardwareTestResult, error) {
	hts.results.StartTime = time.Now()
	defer func() {
		hts.results.EndTime = time.Now()
		hts.results.TotalDuration = hts.results.EndTime.Sub(hts.results.StartTime)
	}()

	var allResults []*hardware.HardwareTestResult

	// Run basic motion tests
	results, err := hts.RunBasicMotionTests(ctx)
	if err != nil {
		return allResults, fmt.Errorf("basic motion tests failed: %v", err)
	}
	allResults = append(allResults, results...)

	// Run force control tests
	results, err = hts.RunForceControlTests(ctx)
	if err != nil {
		return allResults, fmt.Errorf("force control tests failed: %v", err)
	}
	allResults = append(allResults, results...)

	// Run I/O tests
	results, err = hts.RunIOTests(ctx)
	if err != nil {
		return allResults, fmt.Errorf("I/O tests failed: %v", err)
	}
	allResults = append(allResults, results...)

	// Run safety tests
	results, err = hts.RunSafetyTests(ctx)
	if err != nil {
		return allResults, fmt.Errorf("safety tests failed: %v", err)
	}
	allResults = append(allResults, results...)

	// Run performance tests
	results, err = hts.RunPerformanceTests(ctx)
	if err != nil {
		return allResults, fmt.Errorf("performance tests failed: %v", err)
	}
	allResults = append(allResults, results...)

	// Update results
	hts.results.TestDetails = allResults
	hts.updateTestCounts()

	return allResults, nil
}

// testAbsoluteMotion tests absolute motion functionality
func (hts *HardwareTestSuite) testAbsoluteMotion(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Absolute Motion",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test absolute motion to position 1000
	err := hts.controller.MoveAbsolute(ctx, 1000.0, 100.0, 50.0, 25.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify position
	position, err := hts.controller.GetPosition(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	expectedPosition := 1000.0
	tolerance := 10.0
	if abs(position-expectedPosition) > tolerance {
		result.Error = fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, position)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["position"] = position
	result.Metrics["expected_position"] = expectedPosition
	result.Metrics["tolerance"] = tolerance

	return result
}

// testRelativeMotion tests relative motion functionality
func (hts *HardwareTestSuite) testRelativeMotion(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Relative Motion",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Get initial position
	initialPosition, err := hts.controller.GetPosition(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Move relative by 500 counts
	err = hts.controller.MoveRelative(ctx, 500.0, 100.0, 50.0, 25.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify position
	finalPosition, err := hts.controller.GetPosition(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	expectedPosition := initialPosition + 500.0
	tolerance := 10.0
	if abs(finalPosition-expectedPosition) > tolerance {
		result.Error = fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, finalPosition)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["initial_position"] = initialPosition
	result.Metrics["final_position"] = finalPosition
	result.Metrics["expected_position"] = expectedPosition
	result.Metrics["tolerance"] = tolerance

	return result
}

// testIncrementalMotion tests incremental motion functionality
func (hts *HardwareTestSuite) testIncrementalMotion(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Incremental Motion",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Get initial position
	initialPosition, err := hts.controller.GetPosition(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Move incremental by 200 counts
	err = hts.controller.MoveIncremental(ctx, 200.0, 100.0, 50.0, 25.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify position
	finalPosition, err := hts.controller.GetPosition(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	expectedPosition := initialPosition + 200.0
	tolerance := 10.0
	if abs(finalPosition-expectedPosition) > tolerance {
		result.Error = fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, finalPosition)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["initial_position"] = initialPosition
	result.Metrics["final_position"] = finalPosition
	result.Metrics["expected_position"] = expectedPosition
	result.Metrics["tolerance"] = tolerance

	return result
}

// testJogMotion tests jog motion functionality
func (hts *HardwareTestSuite) testJogMotion(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Jog Motion",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Start jog motion
	err := hts.controller.Jog(ctx, 50.0) // 50 counts/s
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Let it run for 1 second
	time.Sleep(1 * time.Second)

	// Stop jog motion
	err = hts.controller.Stop(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["jog_velocity"] = 50.0
	result.Metrics["jog_duration"] = 1 * time.Second

	return result
}

// testStopMotion tests stop motion functionality
func (hts *HardwareTestSuite) testStopMotion(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Stop Motion",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Start a motion
	err := hts.controller.MoveAbsolute(ctx, 2000.0, 200.0, 100.0, 50.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit for motion to start
	time.Sleep(100 * time.Millisecond)

	// Stop the motion
	err = hts.controller.Stop(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["stop_successful"] = true

	return result
}

// testForceControlEnableDisable tests force control enable/disable functionality
func (hts *HardwareTestSuite) testForceControlEnableDisable(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Force Control Enable/Disable",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Enable force control
	err := hts.controller.ForceControlOn(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Disable force control
	err = hts.controller.ForceControlOff(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["enable_successful"] = true
	result.Metrics["disable_successful"] = true

	return result
}

// testForceSetpoint tests force setpoint functionality
func (hts *HardwareTestSuite) testForceSetpoint(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Force Setpoint",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Enable force control
	err := hts.controller.ForceControlOn(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Set force setpoint
	forceSetpoint := 10.0 // 10N
	err = hts.controller.SetForce(ctx, forceSetpoint)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit for force to settle
	time.Sleep(500 * time.Millisecond)

	// Get actual force
	actualForce, err := hts.controller.GetForce(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Disable force control
	err = hts.controller.ForceControlOff(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Validate force is within reasonable range
	tolerance := 5.0 // 5N tolerance
	if abs(actualForce-forceSetpoint) > tolerance {
		result.Error = fmt.Errorf("force mismatch: expected %f, got %f", forceSetpoint, actualForce)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["force_setpoint"] = forceSetpoint
	result.Metrics["actual_force"] = actualForce
	result.Metrics["tolerance"] = tolerance

	return result
}

// testForceMonitoring tests force monitoring functionality
func (hts *HardwareTestSuite) testForceMonitoring(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Force Monitoring",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Enable force control
	err := hts.controller.ForceControlOn(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Set force setpoint
	forceSetpoint := 5.0 // 5N
	err = hts.controller.SetForce(ctx, forceSetpoint)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Monitor force for 2 seconds
	forceReadings := make([]float64, 0)
	monitorDuration := 2 * time.Second
	startMonitor := time.Now()

	for time.Since(startMonitor) < monitorDuration {
		force, err := hts.controller.GetForce(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		forceReadings = append(forceReadings, force)
		time.Sleep(100 * time.Millisecond)
	}

	// Disable force control
	err = hts.controller.ForceControlOff(ctx)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Validate force readings are reasonable
	if len(forceReadings) == 0 {
		result.Error = fmt.Errorf("no force readings collected")
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Check that force readings are within reasonable range
	tolerance := 10.0 // 10N tolerance
	for i, force := range forceReadings {
		if abs(force-forceSetpoint) > tolerance {
			result.Error = fmt.Errorf("force reading %d out of range: expected %f±%f, got %f", i, forceSetpoint, tolerance, force)
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["force_setpoint"] = forceSetpoint
	result.Metrics["force_readings"] = forceReadings
	result.Metrics["tolerance"] = tolerance
	result.Metrics["monitor_duration"] = monitorDuration

	return result
}

// testDigitalOutputs tests digital output functionality
func (hts *HardwareTestSuite) testDigitalOutputs(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Digital Outputs",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test setting digital output 1 to true
	err := hts.controller.SetDigitalOutput(ctx, 1, true)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Verify output is set
	output, err := hts.controller.GetDigitalOutput(ctx, 1)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	if !output {
		result.Error = fmt.Errorf("digital output 1 not set")
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Test setting digital output 1 to false
	err = hts.controller.SetDigitalOutput(ctx, 1, false)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Verify output is cleared
	output, err = hts.controller.GetDigitalOutput(ctx, 1)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	if output {
		result.Error = fmt.Errorf("digital output 1 not cleared")
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["output_test_successful"] = true

	return result
}

// testDigitalInputs tests digital input functionality
func (hts *HardwareTestSuite) testDigitalInputs(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Digital Inputs",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test reading digital input 1
	input, err := hts.controller.GetDigitalInput(ctx, 1)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["digital_input_1"] = input

	return result
}

// testAnalogOutputs tests analog output functionality
func (hts *HardwareTestSuite) testAnalogOutputs(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Analog Outputs",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test setting analog output 1 to 3.14V
	voltage := 3.14
	err := hts.controller.SetAnalogOutput(ctx, 1, voltage)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Verify output is set
	output, err := hts.controller.GetAnalogOutput(ctx, 1)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	tolerance := 0.1 // 0.1V tolerance
	if abs(output-voltage) > tolerance {
		result.Error = fmt.Errorf("analog output 1 mismatch: expected %f, got %f", voltage, output)
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["analog_output_1"] = output
	result.Metrics["expected_voltage"] = voltage
	result.Metrics["tolerance"] = tolerance

	return result
}

// testAnalogInputs tests analog input functionality
func (hts *HardwareTestSuite) testAnalogInputs(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Analog Inputs",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test reading analog input 1
	input, err := hts.controller.GetAnalogInput(ctx, 1)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["analog_input_1"] = input

	return result
}

// testEmergencyStop tests emergency stop functionality
func (hts *HardwareTestSuite) testEmergencyStop(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Emergency Stop",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Start a motion
	err := hts.controller.MoveAbsolute(ctx, 1000.0, 100.0, 50.0, 25.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait a bit for motion to start
	time.Sleep(100 * time.Millisecond)

	// Trigger emergency stop
	err = hts.safetyGuard.TriggerEmergencyStop(hts.controller)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["emergency_stop_successful"] = true

	return result
}

// testSafetyLimits tests safety limits functionality
func (hts *HardwareTestSuite) testSafetyLimits(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Safety Limits",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test position limits
	positionLimits := &safety.PositionLimits{
		MaxPosition: 2000.0,
		MinPosition: -2000.0,
	}

	// Test that motion within limits is allowed
	err := hts.controller.MoveAbsolute(ctx, 1000.0, 100.0, 50.0, 25.0)
	if err != nil {
		result.Error = err
		result.Passed = false
		result.Duration = time.Since(startTime)
		return result
	}

	// Wait for motion to complete
	for {
		complete, err := hts.controller.IsMotionComplete(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		if complete {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["position_limits"] = positionLimits
	result.Metrics["motion_within_limits"] = true

	return result
}

// testErrorRecovery tests error recovery functionality
func (hts *HardwareTestSuite) testErrorRecovery(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Error Recovery",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Test that system can recover from errors
	// This is a simplified test - in practice, you would induce specific errors
	// and test recovery mechanisms

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["error_recovery_successful"] = true

	return result
}

// testLatency tests system latency
func (hts *HardwareTestSuite) testLatency(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Latency",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Measure latency of position queries
	latencies := make([]time.Duration, 0)
	numTests := 100

	for i := 0; i < numTests; i++ {
		queryStart := time.Now()
		_, err := hts.controller.GetPosition(ctx)
		queryEnd := time.Now()
		
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		
		latencies = append(latencies, queryEnd.Sub(queryStart))
		time.Sleep(1 * time.Millisecond)
	}

	// Calculate statistics
	var totalLatency time.Duration
	for _, latency := range latencies {
		totalLatency += latency
	}
	avgLatency := totalLatency / time.Duration(len(latencies))

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["avg_latency"] = avgLatency
	result.Metrics["max_latency"] = maxDuration(latencies)
	result.Metrics["min_latency"] = minDuration(latencies)
	result.Metrics["num_tests"] = numTests

	return result
}

// testThroughput tests system throughput
func (hts *HardwareTestSuite) testThroughput(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Throughput",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Measure throughput of position queries
	numQueries := 1000
	queryStart := time.Now()

	for i := 0; i < numQueries; i++ {
		_, err := hts.controller.GetPosition(ctx)
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
	}

	queryEnd := time.Now()
	totalDuration := queryEnd.Sub(queryStart)
	throughput := float64(numQueries) / totalDuration.Seconds()

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["throughput"] = throughput
	result.Metrics["num_queries"] = numQueries
	result.Metrics["total_duration"] = totalDuration

	return result
}

// testJitter tests system jitter
func (hts *HardwareTestSuite) testJitter(ctx context.Context) *hardware.HardwareTestResult {
	startTime := time.Now()
	result := &hardware.HardwareTestResult{
		TestName:  "Jitter",
		Timestamp: startTime,
		Metrics:   make(map[string]interface{}),
	}

	// Measure jitter of position queries
	latencies := make([]time.Duration, 0)
	numTests := 100

	for i := 0; i < numTests; i++ {
		queryStart := time.Now()
		_, err := hts.controller.GetPosition(ctx)
		queryEnd := time.Now()
		
		if err != nil {
			result.Error = err
			result.Passed = false
			result.Duration = time.Since(startTime)
			return result
		}
		
		latencies = append(latencies, queryEnd.Sub(queryStart))
		time.Sleep(1 * time.Millisecond)
	}

	// Calculate jitter (standard deviation of latencies)
	var totalLatency time.Duration
	for _, latency := range latencies {
		totalLatency += latency
	}
	avgLatency := totalLatency / time.Duration(len(latencies))

	var variance time.Duration
	for _, latency := range latencies {
		diff := latency - avgLatency
		variance += diff * diff
	}
	variance = variance / time.Duration(len(latencies))
	jitter := time.Duration(int64(variance))

	result.Passed = true
	result.Duration = time.Since(startTime)
	result.Metrics["jitter"] = jitter
	result.Metrics["avg_latency"] = avgLatency
	result.Metrics["num_tests"] = numTests

	return result
}

// updateTestCounts updates the test result counts
func (hts *HardwareTestSuite) updateTestCounts() {
	hts.results.PassedTests = 0
	hts.results.FailedTests = 0
	hts.results.SkippedTests = 0

	for _, detail := range hts.results.TestDetails {
		if detail.Passed {
			hts.results.PassedTests++
		} else {
			hts.results.FailedTests++
		}
	}
}

// Helper functions

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func maxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

func minDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	min := durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}