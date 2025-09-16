package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// ExecutionState represents the current state of command table execution
type ExecutionState int

const (
	StateIdle ExecutionState = iota
	StateRunning
	StatePaused
	StateStopped
	StateError
	StateCompleted
)

func (es ExecutionState) String() string {
	switch es {
	case StateIdle:
		return "Idle"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateStopped:
		return "Stopped"
	case StateError:
		return "Error"
	case StateCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

// ExecutionResult represents the result of executing a command
type ExecutionResult struct {
	CommandID   int                    `json:"command_id"`
	Success     bool                   `json:"success"`
	Error       error                  `json:"error,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// ExecutionStatus represents the current status of command table execution
type ExecutionStatus struct {
	State           ExecutionState     `json:"state"`
	CurrentCommand  int                `json:"current_command"`
	TotalCommands   int                `json:"total_commands"`
	Progress        float64            `json:"progress"`
	StartTime       time.Time          `json:"start_time,omitempty"`
	EndTime         time.Time          `json:"end_time,omitempty"`
	ElapsedTime     time.Duration      `json:"elapsed_time"`
	RemainingTime   time.Duration      `json:"remaining_time,omitempty"`
	Error           error              `json:"error,omitempty"`
	Results         []ExecutionResult  `json:"results,omitempty"`
	Variables       map[string]interface{} `json:"variables,omitempty"`
}

// ExecutionEngine defines the interface for executing command tables
type ExecutionEngine interface {
	// Execute starts executing a command table
	Execute(ctx context.Context, table *types.CommandTable) error
	
	// Pause pauses the execution
	Pause() error
	
	// Resume resumes the execution
	Resume() error
	
	// Stop stops the execution
	Stop() error
	
	// GetStatus returns the current execution status
	GetStatus() ExecutionStatus
	
	// WaitForCompletion waits for the execution to complete
	WaitForCompletion(ctx context.Context) error
	
	// IsRunning returns true if the engine is currently executing
	IsRunning() bool
	
	// GetCurrentCommand returns the currently executing command
	GetCurrentCommand() *types.Command
}

// DriveController defines the interface for controlling the LinMot drive
type DriveController interface {
	// Motion commands
	MoveAbsolute(ctx context.Context, position float64, velocity float64, acceleration float64, jerk float64) error
	MoveRelative(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error
	MoveIncremental(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error
	Jog(ctx context.Context, velocity float64) error
	Stop(ctx context.Context) error
	
	// Wait commands
	Wait(ctx context.Context, duration time.Duration) error
	WaitPosition(ctx context.Context, position float64, tolerance float64, timeout time.Duration) error
	WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error
	WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error
	
	// I/O commands
	SetDigitalOutput(ctx context.Context, output int, value bool) error
	ClearDigitalOutput(ctx context.Context, output int) error
	SetAnalogOutput(ctx context.Context, output int, value float64) error
	WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error
	WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error
	
	// System commands
	Home(ctx context.Context) error
	Reset(ctx context.Context) error
	SaveConfiguration(ctx context.Context) error
	LoadConfiguration(ctx context.Context) error
	
	// Force control commands
	ForceControlOn(ctx context.Context) error
	ForceControlOff(ctx context.Context) error
	SetForce(ctx context.Context, force float64) error
	
	// Data acquisition commands
	StartOscilloscope(ctx context.Context) error
	StopOscilloscope(ctx context.Context) error
	SaveData(ctx context.Context, filename string) error
	
	// Status queries
	GetPosition(ctx context.Context) (float64, error)
	GetVelocity(ctx context.Context) (float64, error)
	GetForce(ctx context.Context) (float64, error)
	GetDigitalInput(ctx context.Context, input int) (bool, error)
	GetAnalogInput(ctx context.Context, input int) (float64, error)
	GetDriveState(ctx context.Context) (types.DriveState, error)
	IsMotionComplete(ctx context.Context) (bool, error)
}

// DefaultExecutionEngine provides a default implementation of ExecutionEngine
type DefaultExecutionEngine struct {
	mu              sync.RWMutex
	state           ExecutionState
	currentCommand  int
	commandTable    *types.CommandTable
	results         []ExecutionResult
	startTime       time.Time
	endTime         time.Time
	error           error
	variables       map[string]interface{}
	driveController DriveController
	conditionEvaluator types.ConditionEvaluator
	unitConverter   *types.UnitConverter
	
	// Control channels
	pauseChan  chan struct{}
	resumeChan chan struct{}
	stopChan   chan struct{}
	doneChan   chan struct{}
}

// NewDefaultExecutionEngine creates a new DefaultExecutionEngine
func NewDefaultExecutionEngine(driveController DriveController, conditionEvaluator types.ConditionEvaluator, unitConverter *types.UnitConverter) *DefaultExecutionEngine {
	return &DefaultExecutionEngine{
		state:             StateIdle,
		results:           make([]ExecutionResult, 0),
		variables:         make(map[string]interface{}),
		driveController:   driveController,
		conditionEvaluator: conditionEvaluator,
		unitConverter:     unitConverter,
		pauseChan:         make(chan struct{}),
		resumeChan:        make(chan struct{}),
		stopChan:          make(chan struct{}),
		doneChan:          make(chan struct{}),
	}
}

// Execute starts executing a command table
func (dee *DefaultExecutionEngine) Execute(ctx context.Context, table *types.CommandTable) error {
	dee.mu.Lock()
	defer dee.mu.Unlock()
	
	if dee.state != StateIdle {
		return fmt.Errorf("execution engine is not idle, current state: %s", dee.state)
	}
	
	if table == nil {
		return fmt.Errorf("command table cannot be nil")
	}
	
	if len(table.Commands) == 0 {
		return fmt.Errorf("command table is empty")
	}
	
	// Initialize execution state
	dee.commandTable = table
	dee.currentCommand = 0
	dee.results = make([]ExecutionResult, 0)
	dee.startTime = time.Now()
	dee.endTime = time.Time{}
	dee.error = nil
	dee.variables = make(map[string]interface{})
	
	// Copy table variables
	for k, v := range table.Variables {
		dee.variables[k] = v
	}
	
	dee.state = StateRunning
	
	// Start execution in a goroutine
	go dee.executeCommands(ctx)
	
	return nil
}

// Pause pauses the execution
func (dee *DefaultExecutionEngine) Pause() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()
	
	if dee.state != StateRunning {
		return fmt.Errorf("cannot pause, execution is not running (current state: %s)", dee.state)
	}
	
	dee.state = StatePaused
	dee.pauseChan <- struct{}{}
	
	return nil
}

// Resume resumes the execution
func (dee *DefaultExecutionEngine) Resume() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()
	
	if dee.state != StatePaused {
		return fmt.Errorf("cannot resume, execution is not paused (current state: %s)", dee.state)
	}
	
	dee.state = StateRunning
	dee.resumeChan <- struct{}{}
	
	return nil
}

// Stop stops the execution
func (dee *DefaultExecutionEngine) Stop() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()
	
	if dee.state == StateIdle || dee.state == StateCompleted {
		return fmt.Errorf("cannot stop, execution is not running (current state: %s)", dee.state)
	}
	
	dee.state = StateStopped
	dee.endTime = time.Now()
	
	// Send stop signal in a non-blocking way
	select {
	case dee.stopChan <- struct{}{}:
	default:
		// Channel is full or no one is listening, that's okay
	}
	
	return nil
}

// GetStatus returns the current execution status
func (dee *DefaultExecutionEngine) GetStatus() ExecutionStatus {
	dee.mu.RLock()
	defer dee.mu.RUnlock()
	
	status := ExecutionStatus{
		State:          dee.state,
		CurrentCommand: dee.currentCommand,
		TotalCommands:  len(dee.commandTable.Commands),
		StartTime:      dee.startTime,
		EndTime:        dee.endTime,
		Error:          dee.error,
		Results:        make([]ExecutionResult, len(dee.results)),
		Variables:      make(map[string]interface{}),
	}
	
	// Copy results
	copy(status.Results, dee.results)
	
	// Copy variables
	for k, v := range dee.variables {
		status.Variables[k] = v
	}
	
	// Calculate progress
	if status.TotalCommands > 0 {
		status.Progress = float64(dee.currentCommand) / float64(status.TotalCommands) * 100.0
	}
	
	// Calculate elapsed time
	if !dee.startTime.IsZero() {
		if dee.endTime.IsZero() {
			status.ElapsedTime = time.Since(dee.startTime)
		} else {
			status.ElapsedTime = dee.endTime.Sub(dee.startTime)
		}
	}
	
	// Estimate remaining time
	if status.Progress > 0 && status.ElapsedTime > 0 {
		estimatedTotal := time.Duration(float64(status.ElapsedTime) / (status.Progress / 100.0))
		status.RemainingTime = estimatedTotal - status.ElapsedTime
	}
	
	return status
}

// WaitForCompletion waits for the execution to complete
func (dee *DefaultExecutionEngine) WaitForCompletion(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-dee.doneChan:
		dee.mu.RLock()
		err := dee.error
		dee.mu.RUnlock()
		return err
	}
}

// IsRunning returns true if the engine is currently executing
func (dee *DefaultExecutionEngine) IsRunning() bool {
	dee.mu.RLock()
	defer dee.mu.RUnlock()
	return dee.state == StateRunning
}

// GetCurrentCommand returns the currently executing command
func (dee *DefaultExecutionEngine) GetCurrentCommand() *types.Command {
	dee.mu.RLock()
	defer dee.mu.RUnlock()
	
	if dee.commandTable == nil || dee.currentCommand >= len(dee.commandTable.Commands) {
		return nil
	}
	
	return &dee.commandTable.Commands[dee.currentCommand]
}