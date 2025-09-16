package types

import (
	"context"
	"time"
)

// CommandType represents the type of command in the command table
type CommandType int

const (
	// Motion Commands
	CmdMoveAbsolute    CommandType = iota // MA - Move to absolute position
	CmdMoveRelative                       // MR - Move by relative distance
	CmdMoveIncremental                    // MI - Move by fixed increment
	CmdJog                                // JO - Continuous motion
	CmdStop                               // ST - Stop motion immediately
	
	// Control Commands
	CmdWait                               // WA - Wait for specified time
	CmdWaitPosition                       // WP - Wait for position condition
	CmdWaitVelocity                       // WV - Wait for velocity condition
	CmdWaitForce                          // WF - Wait for force condition
	
	// I/O Commands
	CmdSetDigitalOutput                   // DO - Set digital output
	CmdClearDigitalOutput                 // CO - Clear digital output
	CmdSetAnalogOutput                    // AO - Set analog output
	CmdWaitDigitalInput                   // DI - Wait for digital input
	CmdWaitAnalogInput                    // AI - Wait for analog input
	
	// Loop Commands
	CmdLoopStart                          // LS - Start loop
	CmdLoopEnd                            // LE - End loop
	CmdLoopBreak                          // LB - Break loop
	
	// Jump Commands
	CmdJump                               // JP - Unconditional jump
	CmdJumpIfTrue                         // JT - Jump if condition true
	CmdJumpIfFalse                        // JF - Jump if condition false
	
	// System Commands
	CmdHome                               // HO - Home motor
	CmdReset                              // RE - Reset drive
	CmdSaveConfiguration                  // SC - Save configuration
	CmdLoadConfiguration                  // LC - Load configuration
	
	// Force Control Commands
	CmdForceControlOn                     // FC - Enable force control
	CmdForceControlOff                    // FO - Disable force control
	CmdSetForce                           // SF - Set force setpoint
	
	// Data Acquisition Commands
	CmdStartOscilloscope                  // SO - Start data acquisition
	CmdStopOscilloscope                   // SP - Stop data acquisition
	CmdSaveData                           // SD - Save acquired data
)

// String returns the string representation of the command type
func (ct CommandType) String() string {
	switch ct {
	case CmdMoveAbsolute:
		return "MA"
	case CmdMoveRelative:
		return "MR"
	case CmdMoveIncremental:
		return "MI"
	case CmdJog:
		return "JO"
	case CmdStop:
		return "ST"
	case CmdWait:
		return "WA"
	case CmdWaitPosition:
		return "WP"
	case CmdWaitVelocity:
		return "WV"
	case CmdWaitForce:
		return "WF"
	case CmdSetDigitalOutput:
		return "DO"
	case CmdClearDigitalOutput:
		return "CO"
	case CmdSetAnalogOutput:
		return "AO"
	case CmdWaitDigitalInput:
		return "DI"
	case CmdWaitAnalogInput:
		return "AI"
	case CmdLoopStart:
		return "LS"
	case CmdLoopEnd:
		return "LE"
	case CmdLoopBreak:
		return "LB"
	case CmdJump:
		return "JP"
	case CmdJumpIfTrue:
		return "JT"
	case CmdJumpIfFalse:
		return "JF"
	case CmdHome:
		return "HO"
	case CmdReset:
		return "RE"
	case CmdSaveConfiguration:
		return "SC"
	case CmdLoadConfiguration:
		return "LC"
	case CmdForceControlOn:
		return "FC"
	case CmdForceControlOff:
		return "FO"
	case CmdSetForce:
		return "SF"
	case CmdStartOscilloscope:
		return "SO"
	case CmdStopOscilloscope:
		return "SP"
	case CmdSaveData:
		return "SD"
	default:
		return "UNKNOWN"
	}
}

// Command represents a single command in the command table
type Command struct {
	ID          int                    `json:"id" yaml:"id"`
	Type        CommandType            `json:"type" yaml:"type"`
	Parameters  map[string]interface{} `json:"parameters" yaml:"parameters"`
	Conditions  []Condition            `json:"conditions" yaml:"conditions"`
	NextCommand int                    `json:"next_command" yaml:"next_command"`
	LineNumber  int                    `json:"line_number" yaml:"line_number"`
	Comment     string                 `json:"comment" yaml:"comment"`
	CreatedAt   time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" yaml:"updated_at"`
}

// CommandTable represents a collection of commands that can be executed
type CommandTable struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description string                 `json:"description" yaml:"description"`
	Commands    []Command              `json:"commands" yaml:"commands"`
	Variables   map[string]interface{} `json:"variables" yaml:"variables"`
	CreatedAt   time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" yaml:"updated_at"`
}

// CommandExecutor defines the interface for executing commands
type CommandExecutor interface {
	Execute(ctx context.Context, cmd *Command, vars map[string]interface{}) error
	CanExecute(cmd *Command) bool
	GetRequiredParameters(cmdType CommandType) []string
	ValidateParameters(cmd *Command) error
}

// CommandTableManager defines the interface for managing command tables
type CommandTableManager interface {
	CreateTable(id, name, description string) *CommandTable
	AddCommand(table *CommandTable, cmd *Command) error
	RemoveCommand(table *CommandTable, cmdID int) error
	UpdateCommand(table *CommandTable, cmdID int, cmd *Command) error
	GetCommand(table *CommandTable, cmdID int) (*Command, error)
	ValidateTable(table *CommandTable) error
	LoadTable(data []byte) (*CommandTable, error)
	SaveTable(table *CommandTable) ([]byte, error)
}

// CommandValidator defines the interface for validating commands
type CommandValidator interface {
	ValidateCommand(cmd *Command) error
	ValidateTable(table *CommandTable) error
	CheckDependencies(table *CommandTable) error
	CheckCircularReferences(table *CommandTable) error
}

// CommandBuilder provides a fluent interface for building commands
type CommandBuilder struct {
	cmd *Command
}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{
		cmd: &Command{
			Parameters: make(map[string]interface{}),
			Conditions: make([]Condition, 0),
		},
	}
}

// WithID sets the command ID
func (cb *CommandBuilder) WithID(id int) *CommandBuilder {
	cb.cmd.ID = id
	return cb
}

// WithType sets the command type
func (cb *CommandBuilder) WithType(cmdType CommandType) *CommandBuilder {
	cb.cmd.Type = cmdType
	return cb
}

// WithParameter adds a parameter to the command
func (cb *CommandBuilder) WithParameter(key string, value interface{}) *CommandBuilder {
	cb.cmd.Parameters[key] = value
	return cb
}

// WithParameters sets multiple parameters
func (cb *CommandBuilder) WithParameters(params map[string]interface{}) *CommandBuilder {
	for k, v := range params {
		cb.cmd.Parameters[k] = v
	}
	return cb
}

// WithCondition adds a condition to the command
func (cb *CommandBuilder) WithCondition(condition Condition) *CommandBuilder {
	cb.cmd.Conditions = append(cb.cmd.Conditions, condition)
	return cb
}

// WithNextCommand sets the next command to execute
func (cb *CommandBuilder) WithNextCommand(next int) *CommandBuilder {
	cb.cmd.NextCommand = next
	return cb
}

// WithComment adds a comment to the command
func (cb *CommandBuilder) WithComment(comment string) *CommandBuilder {
	cb.cmd.Comment = comment
	return cb
}

// Build returns the constructed command
func (cb *CommandBuilder) Build() *Command {
	cb.cmd.CreatedAt = time.Now()
	cb.cmd.UpdatedAt = time.Now()
	return cb.cmd
}

// CommandTableBuilder provides a fluent interface for building command tables
type CommandTableBuilder struct {
	table *CommandTable
}

// NewCommandTableBuilder creates a new command table builder
func NewCommandTableBuilder() *CommandTableBuilder {
	return &CommandTableBuilder{
		table: &CommandTable{
			Commands:  make([]Command, 0),
			Variables: make(map[string]interface{}),
		},
	}
}

// WithID sets the table ID
func (ctb *CommandTableBuilder) WithID(id string) *CommandTableBuilder {
	ctb.table.ID = id
	return ctb
}

// WithName sets the table name
func (ctb *CommandTableBuilder) WithName(name string) *CommandTableBuilder {
	ctb.table.Name = name
	return ctb
}

// WithDescription sets the table description
func (ctb *CommandTableBuilder) WithDescription(desc string) *CommandTableBuilder {
	ctb.table.Description = desc
	return ctb
}

// WithCommand adds a command to the table
func (ctb *CommandTableBuilder) WithCommand(cmd *Command) *CommandTableBuilder {
	ctb.table.Commands = append(ctb.table.Commands, *cmd)
	return ctb
}

// WithVariable adds a variable to the table
func (ctb *CommandTableBuilder) WithVariable(name string, value interface{}) *CommandTableBuilder {
	ctb.table.Variables[name] = value
	return ctb
}

// Build returns the constructed command table
func (ctb *CommandTableBuilder) Build() *CommandTable {
	ctb.table.CreatedAt = time.Now()
	ctb.table.UpdatedAt = time.Now()
	return ctb.table
}