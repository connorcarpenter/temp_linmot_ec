package commands

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// CommandExecutor defines the interface for command execution
type CommandExecutor interface {
	Execute(ctx context.Context, command *types.Command) error
	Validate(command *types.Command) error
	GetCommandInfo(commandType types.CommandType) (string, []string, error)
}

// CommandRegistry manages command executors
type CommandRegistry struct {
	executors map[types.CommandType]CommandExecutor
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(driveController types.DriveController, unitConverter *types.UnitConverter) *CommandRegistry {
	registry := &CommandRegistry{
		executors: make(map[types.CommandType]CommandExecutor),
	}
	
	// Register motion command executor
	motionExecutor := NewMotionCommandExecutor(driveController, unitConverter)
	registry.RegisterExecutor(types.CmdMoveAbsolute, motionExecutor)
	registry.RegisterExecutor(types.CmdMoveRelative, motionExecutor)
	registry.RegisterExecutor(types.CmdMoveIncremental, motionExecutor)
	registry.RegisterExecutor(types.CmdJog, motionExecutor)
	registry.RegisterExecutor(types.CmdStop, motionExecutor)
	
	// Register control command executor
	controlExecutor := NewControlCommandExecutor(driveController, unitConverter)
	registry.RegisterExecutor(types.CmdWait, controlExecutor)
	registry.RegisterExecutor(types.CmdWaitPosition, controlExecutor)
	registry.RegisterExecutor(types.CmdWaitVelocity, controlExecutor)
	registry.RegisterExecutor(types.CmdWaitForce, controlExecutor)
	
	// Register I/O command executor
	ioExecutor := NewIOCommandExecutor(driveController, unitConverter)
	registry.RegisterExecutor(types.CmdSetDigitalOutput, ioExecutor)
	registry.RegisterExecutor(types.CmdClearDigitalOutput, ioExecutor)
	registry.RegisterExecutor(types.CmdSetAnalogOutput, ioExecutor)
	registry.RegisterExecutor(types.CmdWaitDigitalInput, ioExecutor)
	registry.RegisterExecutor(types.CmdWaitAnalogInput, ioExecutor)
	
	// Register loop/jump command executor
	loopExecutor := NewLoopCommandExecutor(driveController, unitConverter)
	registry.RegisterExecutor(types.CmdLoopStart, loopExecutor)
	registry.RegisterExecutor(types.CmdLoopEnd, loopExecutor)
	registry.RegisterExecutor(types.CmdLoopBreak, loopExecutor)
	registry.RegisterExecutor(types.CmdJump, loopExecutor)
	registry.RegisterExecutor(types.CmdJumpIfTrue, loopExecutor)
	registry.RegisterExecutor(types.CmdJumpIfFalse, loopExecutor)
	
	// Register system command executor
	systemExecutor := NewSystemCommandExecutor(driveController, unitConverter)
	registry.RegisterExecutor(types.CmdHome, systemExecutor)
	registry.RegisterExecutor(types.CmdReset, systemExecutor)
	registry.RegisterExecutor(types.CmdSaveConfiguration, systemExecutor)
	registry.RegisterExecutor(types.CmdLoadConfiguration, systemExecutor)
	
	return registry
}

// RegisterExecutor registers a command executor for a specific command type
func (cr *CommandRegistry) RegisterExecutor(commandType types.CommandType, executor CommandExecutor) {
	cr.executors[commandType] = executor
}

// GetExecutor returns the executor for a specific command type
func (cr *CommandRegistry) GetExecutor(commandType types.CommandType) (CommandExecutor, error) {
	executor, exists := cr.executors[commandType]
	if !exists {
		return nil, fmt.Errorf("no executor registered for command type: %s", commandType)
	}
	return executor, nil
}

// ExecuteCommand executes a command using the appropriate executor
func (cr *CommandRegistry) ExecuteCommand(ctx context.Context, command *types.Command) error {
	executor, err := cr.GetExecutor(command.Type)
	if err != nil {
		return err
	}
	
	return executor.Execute(ctx, command)
}

// ValidateCommand validates a command using the appropriate executor
func (cr *CommandRegistry) ValidateCommand(command *types.Command) error {
	executor, err := cr.GetExecutor(command.Type)
	if err != nil {
		return err
	}
	
	return executor.Validate(command)
}

// GetCommandInfo returns information about a command type
func (cr *CommandRegistry) GetCommandInfo(commandType types.CommandType) (string, []string, error) {
	executor, err := cr.GetExecutor(commandType)
	if err != nil {
		return "", nil, err
	}
	
	return executor.GetCommandInfo(commandType)
}

// GetSupportedCommandTypes returns all supported command types
func (cr *CommandRegistry) GetSupportedCommandTypes() []types.CommandType {
	types := make([]types.CommandType, 0, len(cr.executors))
	for commandType := range cr.executors {
		types = append(types, commandType)
	}
	return types
}

// IsCommandSupported checks if a command type is supported
func (cr *CommandRegistry) IsCommandSupported(commandType types.CommandType) bool {
	_, exists := cr.executors[commandType]
	return exists
}

// GetCommandCount returns the number of registered command types
func (cr *CommandRegistry) GetCommandCount() int {
	return len(cr.executors)
}

// ListCommandInfo returns information about all supported commands
func (cr *CommandRegistry) ListCommandInfo() map[types.CommandType]CommandInfo {
	info := make(map[types.CommandType]CommandInfo)
	
	// Map command types to their full names
	commandNames := map[types.CommandType]string{
		types.CmdMoveAbsolute:  "MoveAbsolute",
		types.CmdMoveRelative:  "MoveRelative", 
		types.CmdMoveIncremental: "MoveIncremental",
		types.CmdJog:          "Jog",
		types.CmdStop:         "Stop",
		types.CmdWait:         "Wait",
		types.CmdWaitPosition: "WaitPosition",
		types.CmdWaitVelocity: "WaitVelocity",
		types.CmdWaitForce:    "WaitForce",
		types.CmdSetDigitalOutput: "SetDigitalOutput",
		types.CmdClearDigitalOutput: "ClearDigitalOutput",
		types.CmdSetAnalogOutput: "SetAnalogOutput",
		types.CmdWaitDigitalInput: "WaitDigitalInput",
		types.CmdWaitAnalogInput: "WaitAnalogInput",
		types.CmdLoopStart: "LoopStart",
		types.CmdLoopEnd: "LoopEnd",
		types.CmdLoopBreak: "LoopBreak",
		types.CmdJump: "Jump",
		types.CmdJumpIfTrue: "JumpIfTrue",
		types.CmdJumpIfFalse: "JumpIfFalse",
		types.CmdHome: "Home",
		types.CmdReset: "Reset",
		types.CmdSaveConfiguration: "SaveConfiguration",
		types.CmdLoadConfiguration: "LoadConfiguration",
	}
	
	for commandType, executor := range cr.executors {
		description, parameters, err := executor.GetCommandInfo(commandType)
		if err == nil {
			name := commandNames[commandType]
			if name == "" {
				name = fmt.Sprintf("CommandType_%d", int(commandType)) // fallback to type string
			}
			info[commandType] = CommandInfo{
				Type:        commandType,
				Name:        name,
				Description: description,
				Parameters:  parameters,
			}
		}
	}
	
	return info
}

// CommandInfo contains information about a command type
type CommandInfo struct {
	Type        types.CommandType `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  []string          `json:"parameters"`
}

// String returns a string representation of CommandInfo
func (ci CommandInfo) String() string {
	return fmt.Sprintf("%s: %s (parameters: %v)", ci.Name, ci.Description, ci.Parameters)
}