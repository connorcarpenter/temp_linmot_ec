package stage_linmot_ct

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// BenchmarkCommandExecution benchmarks command execution performance
func BenchmarkCommandExecution(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a simple command table
	table := manager.CreateTable("benchmark", "Benchmark", "Performance benchmark")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.StartExecution(ctx, table)
		// Wait for completion
		for {
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted || status.State == types.StateError {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// BenchmarkCommandTableCreation benchmarks command table creation
func BenchmarkCommandTableCreation(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table := manager.CreateTable("benchmark", "Benchmark", "Performance benchmark")
		_ = table
	}
}

// BenchmarkCommandValidation benchmarks command validation
func BenchmarkCommandValidation(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	table := manager.CreateTable("benchmark", "Benchmark", "Performance benchmark")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
		Build()

	manager.AddCommand(table, command)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.ValidateTable(table)
	}
}

// BenchmarkUnitConversion benchmarks unit conversion performance
func BenchmarkUnitConversion(b *testing.B) {
	converter := types.NewUnitConverter()

	position := types.NewPositionValue(100.0, types.PositionUnitMillimeters)
	velocity := types.NewVelocityValue(50.0, types.VelocityUnitMillimetersPerSecond)
	force := types.NewForceValue(10.0, types.ForceUnitNewtons)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.ConvertPositionValue(position, types.PositionUnitCounts)
		converter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
		converter.ConvertForceValue(force, types.ForceUnitCounts)
	}
}

// BenchmarkSafetyValidation benchmarks safety validation performance
func BenchmarkSafetyValidation(b *testing.B) {
	safetyGuard := NewSafetyGuardWithLimits(&SafetyLimits{
		MinPosition: 0.0,
		MaxPosition: 1000.0,
		MaxVelocity: 100.0,
		MinForce:    -50.0,
		MaxForce:    50.0,
	})

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(500.0, types.PositionUnitCounts)).
		WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
		Build()

	unitConverter := types.NewUnitConverter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		safetyGuard.ValidateMotionCommand(command, unitConverter)
	}
}

// BenchmarkConditionEvaluation benchmarks condition evaluation performance
func BenchmarkConditionEvaluation(b *testing.B) {
	driveController := NewMockDriveController()
	conditionEvaluator := types.NewDefaultConditionEvaluator()

	condition := types.NewConditionBuilder().
		WithType(types.CondPosition).
		WithParameter("position").
		WithOperator(types.OpGreaterThan).
		WithValue(100.0).
		Build()

	statusProvider := &MockStatusProvider{
		position: 150.0,
		velocity: 25.0,
		force:    5.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conditionEvaluator.Evaluate(condition, statusProvider, map[string]interface{}{})
	}
}

// BenchmarkStatusMonitoring benchmarks status monitoring performance
func BenchmarkStatusMonitoring(b *testing.B) {
	driveController := NewMockDriveController()
	monitor := NewStatusMonitor(driveController, 10*time.Millisecond)

	monitor.Start()
	defer monitor.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.GetStatus()
	}
}

// BenchmarkCommandRegistry benchmarks command registry performance
func BenchmarkCommandRegistry(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	safetyGuard := NewSafetyGuard()

	registry := NewCommandRegistry(driveController, unitConverter, safetyGuard)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.ExecuteCommand(ctx, command)
	}
}

// BenchmarkLargeCommandTable benchmarks execution of large command tables
func BenchmarkLargeCommandTable(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a large command table
	table := manager.CreateTable("large-benchmark", "Large Benchmark", "Large command table benchmark")

	// Add 100 commands
	for i := 0; i < 100; i++ {
		command := types.NewCommandBuilder().
			WithID(i + 1).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(float64(i*10), types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			Build()

		manager.AddCommand(table, command)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.StartExecution(ctx, table)
		// Wait for completion
		for {
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted || status.State == types.StateError {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// BenchmarkConcurrentExecution benchmarks concurrent command table execution
func BenchmarkConcurrentExecution(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table
	table := manager.CreateTable("concurrent-benchmark", "Concurrent Benchmark", "Concurrent execution benchmark")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdWait).
		WithParameter("duration", types.NewTimeValue(0.01, types.TimeUnitSeconds)).
		Build()

	manager.AddCommand(table, command)

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.StartExecution(ctx, table)
			// Wait for completion
			for {
				status := manager.GetExecutionStatus()
				if status.State == types.StateCompleted || status.State == types.StateError {
					break
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory usage during execution
func BenchmarkMemoryUsage(b *testing.B) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a command table with many commands
	table := manager.CreateTable("memory-benchmark", "Memory Benchmark", "Memory usage benchmark")

	// Add 1000 commands
	for i := 0; i < 1000; i++ {
		command := types.NewCommandBuilder().
			WithID(i + 1).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(float64(i), types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			WithParameter("acceleration", types.NewAccelerationValue(100.0, types.AccelerationUnitCountsS2)).
			WithParameter("jerk", types.NewJerkValue(1000.0, types.JerkUnitCountsS3)).
			Build()

		manager.AddCommand(table, command)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.StartExecution(ctx, table)
		// Wait for completion
		for {
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted || status.State == types.StateError {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// MockStatusProvider for testing
type MockStatusProvider struct {
	position float64
	velocity float64
	force    float64
}

func (msp *MockStatusProvider) GetPosition() float64 {
	return msp.position
}

func (msp *MockStatusProvider) GetVelocity() float64 {
	return msp.velocity
}

func (msp *MockStatusProvider) GetForce() float64 {
	return msp.force
}

func (msp *MockStatusProvider) GetDriveState() types.DriveState {
	return types.DriveStateReady
}

func (msp *MockStatusProvider) IsMotionComplete() bool {
	return true
}

func (msp *MockStatusProvider) GetDigitalInput(input int) bool {
	return false
}

func (msp *MockStatusProvider) GetAnalogInput(input int) float64 {
	return 0.0
}

func (msp *MockStatusProvider) GetVariable(name string) interface{} {
	return nil
}

func (msp *MockStatusProvider) GetError() error {
	return nil
}

// TestPerformanceCharacteristics tests various performance characteristics
func TestPerformanceCharacteristics(t *testing.T) {
	// Test command table creation performance
	start := time.Now()
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create 100 command tables
	for i := 0; i < 100; i++ {
		table := manager.CreateTable("perf-test", "Performance Test", "Performance test")
		_ = table
	}

	creationTime := time.Since(start)
	t.Logf("Created 100 command tables in %v", creationTime)

	// Test command execution performance
	table := manager.CreateTable("exec-test", "Execution Test", "Execution performance test")

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
		Build()

	manager.AddCommand(table, command)

	start = time.Now()
	ctx := context.Background()
	manager.StartExecution(ctx, table)

	// Wait for completion
	for {
		status := manager.GetExecutionStatus()
		if status.State == types.StateCompleted || status.State == types.StateError {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}

	executionTime := time.Since(start)
	t.Logf("Executed command in %v", executionTime)

	// Test unit conversion performance
	start = time.Now()
	position := types.NewPositionValue(100.0, types.PositionUnitMillimeters)
	velocity := types.NewVelocityValue(50.0, types.VelocityUnitMillimetersPerSecond)
	force := types.NewForceValue(10.0, types.ForceUnitNewtons)

	for i := 0; i < 1000; i++ {
		unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
		unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
		unitConverter.ConvertForceValue(force, types.ForceUnitCounts)
	}

	conversionTime := time.Since(start)
	t.Logf("Performed 1000 unit conversions in %v", conversionTime)

	// Test safety validation performance
	start = time.Now()
	safetyCommand := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(500.0, types.PositionUnitCounts)).
		WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
		Build()

	for i := 0; i < 1000; i++ {
		safetyGuard.ValidateMotionCommand(safetyCommand, unitConverter)
	}

	validationTime := time.Since(start)
	t.Logf("Performed 1000 safety validations in %v", validationTime)
}

// TestMemoryUsage tests memory usage characteristics
func TestMemoryUsage(t *testing.T) {
	// Test memory usage with large command tables
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator()
	safetyGuard := NewSafetyGuard()

	manager := NewCommandTableManager(driveController, unitConverter, conditionEvaluator, safetyGuard)

	// Create a large command table
	table := manager.CreateTable("memory-test", "Memory Test", "Memory usage test")

	// Add 1000 commands with various parameters
	for i := 0; i < 1000; i++ {
		command := types.NewCommandBuilder().
			WithID(i + 1).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(float64(i), types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			WithParameter("acceleration", types.NewAccelerationValue(100.0, types.AccelerationUnitCountsS2)).
			WithParameter("jerk", types.NewJerkValue(1000.0, types.JerkUnitCountsS3)).
			WithComment("Command " + string(rune(i))).
			Build()

		manager.AddCommand(table, command)
	}

	// Verify the table was created successfully
	if len(table.Commands) != 1000 {
		t.Errorf("Expected 1000 commands, got %d", len(table.Commands))
	}

	// Test execution of the large table
	ctx := context.Background()
	err := manager.StartExecution(ctx, table)
	if err != nil {
		t.Fatalf("Failed to start execution: %v", err)
	}

	// Wait for completion
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Execution timed out")
		case <-ticker.C:
			status := manager.GetExecutionStatus()
			if status.State == types.StateCompleted {
				t.Log("Large command table executed successfully")
				return
			}
			if status.State == types.StateError {
				t.Fatalf("Execution failed with error: %v", status.Error)
			}
		}
	}
}