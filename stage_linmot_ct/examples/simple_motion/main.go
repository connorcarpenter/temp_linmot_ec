package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/execution"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/safety"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// MockDriveController implements the DriveController interface for testing
type MockDriveController struct {
	position      float64
	velocity      float64
	force         float64
	driveState    types.DriveState
	motionComplete bool
	digitalOutputs map[int]bool
	analogOutputs  map[int]float64
	digitalInputs  map[int]bool
	analogInputs   map[int]float64
}

func NewMockDriveController() *MockDriveController {
	return &MockDriveController{
		driveState:     types.DriveStateReady,
		motionComplete: true,
		digitalOutputs: make(map[int]bool),
		analogOutputs:  make(map[int]float64),
		digitalInputs:  make(map[int]bool),
		analogInputs:   make(map[int]float64),
	}
}

func (mdc *MockDriveController) MoveAbsolute(ctx context.Context, position, velocity, acceleration, jerk float64) error {
	fmt.Printf("Moving to position %.2f at velocity %.2f\n", position, velocity)
	mdc.position = position
	mdc.motionComplete = false
	
	// Simulate motion completion after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		mdc.motionComplete = true
	}()
	
	return nil
}

func (mdc *MockDriveController) MoveRelative(ctx context.Context, distance, velocity, acceleration, jerk float64) error {
	fmt.Printf("Moving relative distance %.2f at velocity %.2f\n", distance, velocity)
	mdc.position += distance
	mdc.motionComplete = false
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		mdc.motionComplete = true
	}()
	
	return nil
}

func (mdc *MockDriveController) Jog(ctx context.Context, velocity float64) error {
	fmt.Printf("Jogging at velocity %.2f\n", velocity)
	mdc.velocity = velocity
	return nil
}

func (mdc *MockDriveController) Stop(ctx context.Context) error {
	fmt.Println("Stopping motion")
	mdc.velocity = 0
	mdc.motionComplete = true
	return nil
}

func (mdc *MockDriveController) GetPosition() (float64, error) {
	return mdc.position, nil
}

func (mdc *MockDriveController) GetVelocity() (float64, error) {
	return mdc.velocity, nil
}

func (mdc *MockDriveController) GetForce() (float64, error) {
	return mdc.force, nil
}

func (mdc *MockDriveController) GetDriveState() (types.DriveState, error) {
	return mdc.driveState, nil
}

func (mdc *MockDriveController) IsMotionComplete() (bool, error) {
	return mdc.motionComplete, nil
}

func (mdc *MockDriveController) SetDigitalOutput(ctx context.Context, output int, value bool) error {
	mdc.digitalOutputs[output] = value
	fmt.Printf("Set digital output %d to %t\n", output, value)
	return nil
}

func (mdc *MockDriveController) ClearDigitalOutput(ctx context.Context, output int) error {
	mdc.digitalOutputs[output] = false
	fmt.Printf("Cleared digital output %d\n", output)
	return nil
}

func (mdc *MockDriveController) SetAnalogOutput(ctx context.Context, output int, value float64) error {
	mdc.analogOutputs[output] = value
	fmt.Printf("Set analog output %d to %.3f\n", output, value)
	return nil
}

func (mdc *MockDriveController) GetDigitalOutput(ctx context.Context, output int) (bool, error) {
	return mdc.digitalOutputs[output], nil
}

func (mdc *MockDriveController) GetAnalogOutput(ctx context.Context, output int) (float64, error) {
	return mdc.analogOutputs[output], nil
}

func (mdc *MockDriveController) WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error {
	fmt.Printf("Waiting for digital input %d to be %t (timeout: %v)\n", input, value, timeout)
	time.Sleep(50 * time.Millisecond) // Simulate wait
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value, tolerance float64, timeout time.Duration) error {
	fmt.Printf("Waiting for analog input %d to be %.3f±%.3f (timeout: %v)\n", input, value, tolerance, timeout)
	time.Sleep(50 * time.Millisecond) // Simulate wait
	return nil
}

func (mdc *MockDriveController) Home(ctx context.Context) error {
	fmt.Println("Homing drive")
	mdc.position = 0
	mdc.motionComplete = false
	
	go func() {
		time.Sleep(200 * time.Millisecond)
		mdc.motionComplete = true
	}()
	
	return nil
}

func (mdc *MockDriveController) Reset(ctx context.Context) error {
	fmt.Println("Resetting drive")
	mdc.driveState = types.DriveStateReady
	return nil
}

func (mdc *MockDriveController) SaveConfiguration(ctx context.Context) error {
	fmt.Println("Saving configuration")
	return nil
}

func (mdc *MockDriveController) LoadConfiguration(ctx context.Context) error {
	fmt.Println("Loading configuration")
	return nil
}

func (mdc *MockDriveController) ForceControlOn(ctx context.Context) error {
	fmt.Println("Force control enabled")
	return nil
}

func (mdc *MockDriveController) ForceControlOff(ctx context.Context) error {
	fmt.Println("Force control disabled")
	return nil
}

func (mdc *MockDriveController) SetForce(ctx context.Context, force float64) error {
	mdc.force = force
	fmt.Printf("Set force to %.2f\n", force)
	return nil
}

func (mdc *MockDriveController) StartOscilloscope(ctx context.Context) error {
	fmt.Println("Started oscilloscope")
	return nil
}

func (mdc *MockDriveController) StopOscilloscope(ctx context.Context) error {
	fmt.Println("Stopped oscilloscope")
	return nil
}

func (mdc *MockDriveController) SaveData(ctx context.Context, filename string) error {
	fmt.Printf("Saved data to %s\n", filename)
	return nil
}

func (mdc *MockDriveController) GetAnalogInput(ctx context.Context, input int) (float64, error) {
	return mdc.analogInputs[input], nil
}

func main() {
	fmt.Println("=== Simple Motion Sequence Example ===")
	
	// Create components
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	conditionEvaluator := types.NewDefaultConditionEvaluator(driveController)
	safetyGuard := safety.NewSafetyGuard()
	
	// Create execution engine
	executionEngine := execution.NewDefaultExecutionEngine(
		driveController, unitConverter, conditionEvaluator, safetyGuard,
	)
	
	// Create command table manager
	manager := stage_linmot_ct.NewCommandTableManager(executionEngine, unitConverter, nil)
	
	// Create a command table
	table := manager.CreateTable("simple-motion", "Simple Motion", "A basic motion sequence")
	
	// Add commands to the table
	commands := []*types.Command{
		// Home the drive
		types.NewCommandBuilder().
			WithID(1).
			WithType(types.CmdHome).
			WithComment("Home the drive").
			Build(),
		
		// Move to position 100
		types.NewCommandBuilder().
			WithID(2).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(100.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(50.0, types.VelocityUnitCountsS)).
			WithParameter("acceleration", types.NewAccelerationValue(100.0, types.AccelerationUnitCountsS2)).
			WithComment("Move to position 100").
			Build(),
		
		// Wait for 1 second
		types.NewCommandBuilder().
			WithID(3).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(1.0, types.TimeUnitS)).
			WithComment("Wait for 1 second").
			Build(),
		
		// Move to position 200
		types.NewCommandBuilder().
			WithID(4).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(200.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(75.0, types.VelocityUnitCountsS)).
			WithComment("Move to position 200").
			Build(),
		
		// Wait for 0.5 seconds
		types.NewCommandBuilder().
			WithID(5).
			WithType(types.CmdWait).
			WithParameter("duration", types.NewTimeValue(0.5, types.TimeUnitS)).
			WithComment("Wait for 0.5 seconds").
			Build(),
		
		// Return to home position
		types.NewCommandBuilder().
			WithID(6).
			WithType(types.CmdMoveAbsolute).
			WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
			WithParameter("velocity", types.NewVelocityValue(25.0, types.VelocityUnitCountsS)).
			WithComment("Return to home position").
			Build(),
	}
	
	// Add all commands to the table
	for _, cmd := range commands {
		err := manager.AddCommand(table, cmd)
		if err != nil {
			log.Fatalf("Failed to add command: %v", err)
		}
	}
	
	// Validate the table
	err := manager.ValidateTable(table)
	if err != nil {
		log.Fatalf("Table validation failed: %v", err)
	}
	
	fmt.Printf("Created command table with %d commands\n", len(table.Commands))
	
	// Execute the table
	ctx := context.Background()
	err = manager.StartExecution(ctx, table)
	if err != nil {
		log.Fatalf("Failed to start execution: %v", err)
	}
	
	fmt.Println("Execution started...")
	
	// Monitor execution
	for {
		status := manager.GetExecutionStatus()
		fmt.Printf("Execution state: %s, Current command: %d\n", status.State, status.CurrentCommand)
		
		switch status.State {
		case execution.StateCompleted:
			fmt.Println("✅ Execution completed successfully!")
			return
		case execution.StateError:
			log.Fatalf("❌ Execution failed: %v", status.Error)
		case execution.StateStopped:
			fmt.Println("⏹️ Execution stopped")
			return
		case execution.StatePaused:
			fmt.Println("⏸️ Execution paused")
		case execution.StateRunning:
			fmt.Println("▶️ Execution running...")
		}
		
		time.Sleep(200 * time.Millisecond)
	}
}