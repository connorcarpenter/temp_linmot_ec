package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// ExecutionState represents the current state of command execution
type ExecutionState int

const (
	StateIdle        ExecutionState = iota
	StateRunning
	StatePaused
	StateStopped
	StateError
	StateCompleted
)

// String returns the string representation of the execution state
func (es ExecutionState) String() string {
	switch es {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StatePaused:
		return "paused"
	case StateStopped:
		return "stopped"
	case StateError:
		return "error"
	case StateCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

// ExecutionStatus represents the current status of command execution
type ExecutionStatus struct {
	State           ExecutionState           `json:"state" yaml:"state"`
	CurrentCommand  int                      `json:"current_command" yaml:"current_command"`
	TotalCommands   int                      `json:"total_commands" yaml:"total_commands"`
	Progress        float64                  `json:"progress" yaml:"progress"`
	StartTime       time.Time                `json:"start_time" yaml:"start_time"`
	LastUpdateTime  time.Time                `json:"last_update_time" yaml:"last_update_time"`
	Error           error                    `json:"error,omitempty" yaml:"error,omitempty"`
	Variables       map[string]interface{}   `json:"variables" yaml:"variables"`
	ExecutionTime   time.Duration            `json:"execution_time" yaml:"execution_time"`
	RemainingTime   time.Duration            `json:"remaining_time,omitempty" yaml:"remaining_time,omitempty"`
}

// ExecutionEngine defines the interface for executing command tables
type ExecutionEngine interface {
	Start(ctx context.Context, table *types.CommandTable) error
	Pause() error
	Resume() error
	Stop() error
	GetStatus() ExecutionStatus
	GetCurrentCommand() *types.Command
	GetVariables() map[string]interface{}
	SetVariable(name string, value interface{}) error
	GetExecutionHistory() []ExecutionStep
	ClearHistory() error
}

// ExecutionStep represents a single step in command execution
type ExecutionStep struct {
	CommandID     int                    `json:"command_id" yaml:"command_id"`
	CommandType   types.CommandType      `json:"command_type" yaml:"command_type"`
	StartTime     time.Time              `json:"start_time" yaml:"start_time"`
	EndTime       time.Time              `json:"end_time" yaml:"end_time"`
	Duration      time.Duration          `json:"duration" yaml:"duration"`
	Success       bool                   `json:"success" yaml:"success"`
	Error         error                  `json:"error,omitempty" yaml:"error,omitempty"`
	Variables     map[string]interface{} `json:"variables" yaml:"variables"`
	Result        interface{}            `json:"result,omitempty" yaml:"result,omitempty"`
}

// DefaultExecutionEngine provides a default implementation of ExecutionEngine
type DefaultExecutionEngine struct {
	mu                sync.RWMutex
	state             ExecutionState
	currentCommand    int
	commandTable      *types.CommandTable
	variables         map[string]interface{}
	startTime         time.Time
	lastUpdateTime    time.Time
	error             error
	cancelFunc        context.CancelFunc
	executionHistory  []ExecutionStep
	commandExecutors  map[types.CommandType]CommandExecutor
	conditionEvaluator types.ConditionEvaluator
	unitConverter     *types.UnitConverter
	safetyGuard       SafetyGuard
}

// CommandExecutor defines the interface for executing individual commands
type CommandExecutor interface {
	Execute(ctx context.Context, cmd *types.Command, vars map[string]interface{}) (interface{}, error)
	CanExecute(cmd *types.Command) bool
	GetRequiredParameters(cmdType types.CommandType) []string
	ValidateParameters(cmd *types.Command) error
}

// SafetyGuard defines the interface for safety checks
type SafetyGuard interface {
	CheckPreconditions(ctx context.Context, cmd *types.Command) error
	ValidateParameters(cmd *types.Command) error
	CheckLimits(cmd *types.Command) error
	HandleError(err error) error
}

// NewDefaultExecutionEngine creates a new default execution engine
func NewDefaultExecutionEngine(
	commandExecutors map[types.CommandType]CommandExecutor,
	conditionEvaluator types.ConditionEvaluator,
	unitConverter *types.UnitConverter,
	safetyGuard SafetyGuard,
) *DefaultExecutionEngine {
	return &DefaultExecutionEngine{
		state:             StateIdle,
		currentCommand:    0,
		variables:         make(map[string]interface{}),
		executionHistory:  make([]ExecutionStep, 0),
		commandExecutors:  commandExecutors,
		conditionEvaluator: conditionEvaluator,
		unitConverter:     unitConverter,
		safetyGuard:       safetyGuard,
	}
}

// Start starts execution of the command table
func (dee *DefaultExecutionEngine) Start(ctx context.Context, table *types.CommandTable) error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state != StateIdle {
		return fmt.Errorf("execution engine is not idle, current state: %s", dee.state.String())
	}

	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	if len(table.Commands) == 0 {
		return fmt.Errorf("command table is empty")
	}

	// Create a new context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	dee.cancelFunc = cancel

	// Initialize execution state
	dee.commandTable = table
	dee.currentCommand = 0
	dee.variables = make(map[string]interface{})
	// Copy table variables
	for k, v := range table.Variables {
		dee.variables[k] = v
	}
	dee.startTime = time.Now()
	dee.lastUpdateTime = time.Now()
	dee.error = nil
	dee.executionHistory = make([]ExecutionStep, 0)
	dee.state = StateRunning

	// Start execution in a goroutine
	go dee.executeCommands(ctx)

	return nil
}

// Pause pauses execution
func (dee *DefaultExecutionEngine) Pause() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state != StateRunning {
		return fmt.Errorf("execution engine is not running, current state: %s", dee.state.String())
	}

	dee.state = StatePaused
	dee.lastUpdateTime = time.Now()

	return nil
}

// Resume resumes execution
func (dee *DefaultExecutionEngine) Resume() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state != StatePaused {
		return fmt.Errorf("execution engine is not paused, current state: %s", dee.state.String())
	}

	dee.state = StateRunning
	dee.lastUpdateTime = time.Now()

	return nil
}

// Stop stops execution
func (dee *DefaultExecutionEngine) Stop() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state == StateIdle || dee.state == StateStopped {
		return fmt.Errorf("execution engine is not running, current state: %s", dee.state.String())
	}

	dee.state = StateStopped
	dee.lastUpdateTime = time.Now()

	if dee.cancelFunc != nil {
		dee.cancelFunc()
	}

	return nil
}

// GetStatus returns the current execution status
func (dee *DefaultExecutionEngine) GetStatus() ExecutionStatus {
	dee.mu.RLock()
	defer dee.mu.RUnlock()

	progress := 0.0
	if dee.commandTable != nil && len(dee.commandTable.Commands) > 0 {
		progress = float64(dee.currentCommand) / float64(len(dee.commandTable.Commands)) * 100.0
	}

	executionTime := time.Since(dee.startTime)
	remainingTime := time.Duration(0)
	if dee.state == StateRunning && progress > 0 {
		estimatedTotal := time.Duration(float64(executionTime) / progress * 100.0)
		remainingTime = estimatedTotal - executionTime
	}

	return ExecutionStatus{
		State:           dee.state,
		CurrentCommand:  dee.currentCommand,
		TotalCommands:   len(dee.commandTable.Commands),
		Progress:        progress,
		StartTime:       dee.startTime,
		LastUpdateTime:  dee.lastUpdateTime,
		Error:           dee.error,
		Variables:       dee.variables,
		ExecutionTime:   executionTime,
		RemainingTime:   remainingTime,
	}
}

// GetCurrentCommand returns the current command being executed
func (dee *DefaultExecutionEngine) GetCurrentCommand() *types.Command {
	dee.mu.RLock()
	defer dee.mu.RUnlock()

	if dee.commandTable == nil || dee.currentCommand >= len(dee.commandTable.Commands) {
		return nil
	}

	return &dee.commandTable.Commands[dee.currentCommand]
}

// GetVariables returns the current variables
func (dee *DefaultExecutionEngine) GetVariables() map[string]interface{} {
	dee.mu.RLock()
	defer dee.mu.RUnlock()

	// Return a copy to prevent external modification
	vars := make(map[string]interface{})
	for k, v := range dee.variables {
		vars[k] = v
	}
	return vars
}

// SetVariable sets a variable value
func (dee *DefaultExecutionEngine) SetVariable(name string, value interface{}) error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state == StateIdle {
		return fmt.Errorf("execution engine is not running")
	}

	dee.variables[name] = value
	dee.lastUpdateTime = time.Now()

	return nil
}

// GetExecutionHistory returns the execution history
func (dee *DefaultExecutionEngine) GetExecutionHistory() []ExecutionStep {
	dee.mu.RLock()
	defer dee.mu.RUnlock()

	// Return a copy to prevent external modification
	history := make([]ExecutionStep, len(dee.executionHistory))
	copy(history, dee.executionHistory)
	return history
}

// ClearHistory clears the execution history
func (dee *DefaultExecutionEngine) ClearHistory() error {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	if dee.state == StateRunning {
		return fmt.Errorf("cannot clear history while execution is running")
	}

	dee.executionHistory = make([]ExecutionStep, 0)
	return nil
}

// executeCommands executes the commands in the command table
func (dee *DefaultExecutionEngine) executeCommands(ctx context.Context) {
	defer func() {
		dee.mu.Lock()
		if dee.state == StateRunning {
			dee.state = StateCompleted
		}
		dee.lastUpdateTime = time.Now()
		dee.mu.Unlock()
	}()

	for {
		// Check for cancellation
		select {
		case <-ctx.Done():
			dee.mu.Lock()
			dee.state = StateStopped
			dee.error = ctx.Err()
			dee.lastUpdateTime = time.Now()
			dee.mu.Unlock()
			return
		default:
		}

		dee.mu.RLock()
		state := dee.state
		currentCmd := dee.currentCommand
		commandTable := dee.commandTable
		dee.mu.RUnlock()

		// Check if execution should continue
		if state != StateRunning {
			if state == StatePaused {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return
		}

		// Check if we've reached the end of the command table
		if currentCmd >= len(commandTable.Commands) {
			dee.mu.Lock()
			dee.state = StateCompleted
			dee.lastUpdateTime = time.Now()
			dee.mu.Unlock()
			return
		}

		// Get the current command
		cmd := &commandTable.Commands[currentCmd]

		// Execute the command
		err := dee.executeCommand(ctx, cmd)
		if err != nil {
			dee.mu.Lock()
			dee.state = StateError
			dee.error = err
			dee.lastUpdateTime = time.Now()
			dee.mu.Unlock()
			return
		}

		// Move to next command
		dee.mu.Lock()
		dee.currentCommand++
		dee.lastUpdateTime = time.Now()
		dee.mu.Unlock()
	}
}

// executeCommand executes a single command
func (dee *DefaultExecutionEngine) executeCommand(ctx context.Context, cmd *types.Command) error {
	startTime := time.Now()

	// Create execution step
	step := ExecutionStep{
		CommandID:   cmd.ID,
		CommandType: cmd.Type,
		StartTime:   startTime,
		Variables:   make(map[string]interface{}),
	}

	// Copy current variables
	for k, v := range dee.variables {
		step.Variables[k] = v
	}

	// Check preconditions
	if dee.safetyGuard != nil {
		if err := dee.safetyGuard.CheckPreconditions(ctx, cmd); err != nil {
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime)
			step.Success = false
			step.Error = err
			dee.addExecutionStep(step)
			return err
		}
	}

	// Validate parameters
	if dee.safetyGuard != nil {
		if err := dee.safetyGuard.ValidateParameters(cmd); err != nil {
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime)
			step.Success = false
			step.Error = err
			dee.addExecutionStep(step)
			return err
		}
	}

	// Check conditions
	if len(cmd.Conditions) > 0 {
		allConditionsMet := true
		for _, condition := range cmd.Conditions {
			if dee.conditionEvaluator != nil {
				met, err := dee.conditionEvaluator.Evaluate(ctx, &condition, dee.variables)
				if err != nil {
					step.EndTime = time.Now()
					step.Duration = step.EndTime.Sub(step.StartTime)
					step.Success = false
					step.Error = err
					dee.addExecutionStep(step)
					return err
				}
				if !met {
					allConditionsMet = false
					break
				}
			}
		}
		if !allConditionsMet {
			// Skip this command
			step.EndTime = time.Now()
			step.Duration = step.EndTime.Sub(step.StartTime)
			step.Success = true
			step.Error = nil
			dee.addExecutionStep(step)
			return nil
		}
	}

	// Execute the command
	executor, ok := dee.commandExecutors[cmd.Type]
	if !ok {
		err := fmt.Errorf("no executor found for command type: %s", cmd.Type.String())
		step.EndTime = time.Now()
		step.Duration = step.EndTime.Sub(step.StartTime)
		step.Success = false
		step.Error = err
		dee.addExecutionStep(step)
		return err
	}

	result, err := executor.Execute(ctx, cmd, dee.variables)
	step.EndTime = time.Now()
	step.Duration = step.EndTime.Sub(step.StartTime)
	step.Success = err == nil
	step.Error = err
	step.Result = result

	// Add execution step to history
	dee.addExecutionStep(step)

	if err != nil {
		return err
	}

	// Handle next command logic
	if cmd.NextCommand > 0 {
		dee.mu.Lock()
		dee.currentCommand = cmd.NextCommand - 1 // Convert to 0-based index
		dee.mu.Unlock()
	}

	return nil
}

// addExecutionStep adds an execution step to the history
func (dee *DefaultExecutionEngine) addExecutionStep(step ExecutionStep) {
	dee.mu.Lock()
	defer dee.mu.Unlock()

	dee.executionHistory = append(dee.executionHistory, step)
}

// ExecutionEngineFactory creates execution engines for different use cases
type ExecutionEngineFactory struct {
	defaultExecutors map[types.CommandType]CommandExecutor
	defaultConditionEvaluator types.ConditionEvaluator
	defaultUnitConverter *types.UnitConverter
	defaultSafetyGuard SafetyGuard
}

// NewExecutionEngineFactory creates a new execution engine factory
func NewExecutionEngineFactory() *ExecutionEngineFactory {
	return &ExecutionEngineFactory{
		defaultExecutors: make(map[types.CommandType]CommandExecutor),
	}
}

// RegisterCommandExecutor registers a command executor for a specific command type
func (eef *ExecutionEngineFactory) RegisterCommandExecutor(cmdType types.CommandType, executor CommandExecutor) {
	eef.defaultExecutors[cmdType] = executor
}

// SetDefaultConditionEvaluator sets the default condition evaluator
func (eef *ExecutionEngineFactory) SetDefaultConditionEvaluator(evaluator types.ConditionEvaluator) {
	eef.defaultConditionEvaluator = evaluator
}

// SetDefaultUnitConverter sets the default unit converter
func (eef *ExecutionEngineFactory) SetDefaultUnitConverter(converter *types.UnitConverter) {
	eef.defaultUnitConverter = converter
}

// SetDefaultSafetyGuard sets the default safety guard
func (eef *ExecutionEngineFactory) SetDefaultSafetyGuard(guard SafetyGuard) {
	eef.defaultSafetyGuard = guard
}

// CreateExecutionEngine creates a new execution engine with the registered components
func (eef *ExecutionEngineFactory) CreateExecutionEngine() *DefaultExecutionEngine {
	return NewDefaultExecutionEngine(
		eef.defaultExecutors,
		eef.defaultConditionEvaluator,
		eef.defaultUnitConverter,
		eef.defaultSafetyGuard,
	)
}