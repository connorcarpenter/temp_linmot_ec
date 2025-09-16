package stage_linmot_ct

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/execution"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// CommandTableManager provides high-level management of command tables
type CommandTableManager struct {
	executionEngine *execution.DefaultExecutionEngine
	unitConverter   *types.UnitConverter
	validator       CommandTableValidator
}

// CommandTableValidator defines the interface for validating command tables
type CommandTableValidator interface {
	ValidateCommand(cmd *types.Command) error
	ValidateTable(table *types.CommandTable) error
	CheckDependencies(table *types.CommandTable) error
	CheckCircularReferences(table *types.CommandTable) error
}

// NewCommandTableManager creates a new command table manager
func NewCommandTableManager(
	executionEngine *execution.DefaultExecutionEngine,
	unitConverter *types.UnitConverter,
	validator CommandTableValidator,
) *CommandTableManager {
	return &CommandTableManager{
		executionEngine: executionEngine,
		unitConverter:   unitConverter,
		validator:       validator,
	}
}

// CreateTable creates a new command table
func (ctm *CommandTableManager) CreateTable(id, name, description string) *types.CommandTable {
	return &types.CommandTable{
		ID:          id,
		Name:        name,
		Description: description,
		Commands:    make([]types.Command, 0),
		Variables:   make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// AddCommand adds a command to the table
func (ctm *CommandTableManager) AddCommand(table *types.CommandTable, cmd *types.Command) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	if cmd == nil {
		return fmt.Errorf("command is nil")
	}

	// Validate the command
	if ctm.validator != nil {
		if err := ctm.validator.ValidateCommand(cmd); err != nil {
			return fmt.Errorf("command validation failed: %w", err)
		}
	}

	// Check if command ID already exists
	for _, existingCmd := range table.Commands {
		if existingCmd.ID == cmd.ID {
			return fmt.Errorf("command with ID %d already exists", cmd.ID)
		}
	}

	// Set command timestamps
	cmd.CreatedAt = time.Now()
	cmd.UpdatedAt = time.Now()

	// Add command to table
	table.Commands = append(table.Commands, *cmd)
	table.UpdatedAt = time.Now()

	return nil
}

// RemoveCommand removes a command from the table
func (ctm *CommandTableManager) RemoveCommand(table *types.CommandTable, cmdID int) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	// Find and remove the command
	for i, cmd := range table.Commands {
		if cmd.ID == cmdID {
			table.Commands = append(table.Commands[:i], table.Commands[i+1:]...)
			table.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("command with ID %d not found", cmdID)
}

// UpdateCommand updates a command in the table
func (ctm *CommandTableManager) UpdateCommand(table *types.CommandTable, cmdID int, cmd *types.Command) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	if cmd == nil {
		return fmt.Errorf("command is nil")
	}

	// Validate the command
	if ctm.validator != nil {
		if err := ctm.validator.ValidateCommand(cmd); err != nil {
			return fmt.Errorf("command validation failed: %w", err)
		}
	}

	// Find and update the command
	for i, existingCmd := range table.Commands {
		if existingCmd.ID == cmdID {
			cmd.ID = cmdID // Ensure ID matches
			cmd.UpdatedAt = time.Now()
			table.Commands[i] = *cmd
			table.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("command with ID %d not found", cmdID)
}

// GetCommand gets a command from the table
func (ctm *CommandTableManager) GetCommand(table *types.CommandTable, cmdID int) (*types.Command, error) {
	if table == nil {
		return nil, fmt.Errorf("command table is nil")
	}

	for _, cmd := range table.Commands {
		if cmd.ID == cmdID {
			return &cmd, nil
		}
	}

	return nil, fmt.Errorf("command with ID %d not found", cmdID)
}

// ValidateTable validates the entire command table
func (ctm *CommandTableManager) ValidateTable(table *types.CommandTable) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	if ctm.validator == nil {
		return nil // No validator configured
	}

	// Validate individual commands
	for _, cmd := range table.Commands {
		if err := ctm.validator.ValidateCommand(&cmd); err != nil {
			return fmt.Errorf("command %d validation failed: %w", cmd.ID, err)
		}
	}

	// Validate table-level dependencies
	if err := ctm.validator.CheckDependencies(table); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	// Check for circular references
	if err := ctm.validator.CheckCircularReferences(table); err != nil {
		return fmt.Errorf("circular reference check failed: %w", err)
	}

	return nil
}

// LoadTable loads a command table from JSON data
func (ctm *CommandTableManager) LoadTable(data []byte) (*types.CommandTable, error) {
	var table types.CommandTable
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, fmt.Errorf("failed to unmarshal command table: %w", err)
	}

	// Validate the loaded table
	if err := ctm.ValidateTable(&table); err != nil {
		return nil, fmt.Errorf("loaded table validation failed: %w", err)
	}

	return &table, nil
}

// SaveTable saves a command table to JSON data
func (ctm *CommandTableManager) SaveTable(table *types.CommandTable) ([]byte, error) {
	if table == nil {
		return nil, fmt.Errorf("command table is nil")
	}

	// Validate before saving
	if err := ctm.ValidateTable(table); err != nil {
		return nil, fmt.Errorf("table validation failed: %w", err)
	}

	data, err := json.MarshalIndent(table, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command table: %w", err)
	}

	return data, nil
}

// StartExecution starts execution of a command table
func (ctm *CommandTableManager) StartExecution(ctx context.Context, table *types.CommandTable) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}

	// Validate table before execution
	if err := ctm.ValidateTable(table); err != nil {
		return fmt.Errorf("table validation failed: %w", err)
	}

	// Convert units if needed
	if ctm.unitConverter != nil {
		// This would convert all commands to the target unit system
		// Implementation depends on the specific unit system requirements
	}

	return ctm.executionEngine.Start(ctx, table)
}

// PauseExecution pauses execution
func (ctm *CommandTableManager) PauseExecution() error {
	return ctm.executionEngine.Pause()
}

// ResumeExecution resumes execution
func (ctm *CommandTableManager) ResumeExecution() error {
	return ctm.executionEngine.Resume()
}

// StopExecution stops execution
func (ctm *CommandTableManager) StopExecution() error {
	return ctm.executionEngine.Stop()
}

// GetExecutionStatus returns the current execution status
func (ctm *CommandTableManager) GetExecutionStatus() execution.ExecutionStatus {
	return ctm.executionEngine.GetStatus()
}

// GetCurrentCommand returns the current command being executed
func (ctm *CommandTableManager) GetCurrentCommand() *types.Command {
	return ctm.executionEngine.GetCurrentCommand()
}

// GetVariables returns the current variables
func (ctm *CommandTableManager) GetVariables() map[string]interface{} {
	return ctm.executionEngine.GetVariables()
}

// SetVariable sets a variable value
func (ctm *CommandTableManager) SetVariable(name string, value interface{}) error {
	return ctm.executionEngine.SetVariable(name, value)
}

// GetExecutionHistory returns the execution history
func (ctm *CommandTableManager) GetExecutionHistory() []execution.ExecutionStep {
	return ctm.executionEngine.GetExecutionHistory()
}

// ClearHistory clears the execution history
func (ctm *CommandTableManager) ClearHistory() error {
	return ctm.executionEngine.ClearHistory()
}

// CommandTableBuilder provides a fluent interface for building command tables
type CommandTableBuilder struct {
	manager *CommandTableManager
	table   *types.CommandTable
}

// NewCommandTableBuilder creates a new command table builder
func NewCommandTableBuilder(manager *CommandTableManager) *CommandTableBuilder {
	return &CommandTableBuilder{
		manager: manager,
		table:   manager.CreateTable("", "", ""),
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
func (ctb *CommandTableBuilder) WithCommand(cmd *types.Command) *CommandTableBuilder {
	ctb.table.Commands = append(ctb.table.Commands, *cmd)
	return ctb
}

// WithVariable adds a variable to the table
func (ctb *CommandTableBuilder) WithVariable(name string, value interface{}) *CommandTableBuilder {
	ctb.table.Variables[name] = value
	return ctb
}

// Build returns the constructed command table
func (ctb *CommandTableBuilder) Build() *types.CommandTable {
	ctb.table.CreatedAt = time.Now()
	ctb.table.UpdatedAt = time.Now()
	return ctb.table
}

// CommandTableRepository provides persistence for command tables
type CommandTableRepository interface {
	Save(table *types.CommandTable) error
	Load(id string) (*types.CommandTable, error)
	Delete(id string) error
	List() ([]*types.CommandTable, error)
	Exists(id string) (bool, error)
}

// InMemoryCommandTableRepository provides an in-memory implementation of CommandTableRepository
type InMemoryCommandTableRepository struct {
	tables map[string]*types.CommandTable
}

// NewInMemoryCommandTableRepository creates a new in-memory command table repository
func NewInMemoryCommandTableRepository() *InMemoryCommandTableRepository {
	return &InMemoryCommandTableRepository{
		tables: make(map[string]*types.CommandTable),
	}
}

// Save saves a command table
func (repo *InMemoryCommandTableRepository) Save(table *types.CommandTable) error {
	if table == nil {
		return fmt.Errorf("command table is nil")
	}
	repo.tables[table.ID] = table
	return nil
}

// Load loads a command table
func (repo *InMemoryCommandTableRepository) Load(id string) (*types.CommandTable, error) {
	table, exists := repo.tables[id]
	if !exists {
		return nil, fmt.Errorf("command table with ID %s not found", id)
	}
	return table, nil
}

// Delete deletes a command table
func (repo *InMemoryCommandTableRepository) Delete(id string) error {
	if _, exists := repo.tables[id]; !exists {
		return fmt.Errorf("command table with ID %s not found", id)
	}
	delete(repo.tables, id)
	return nil
}

// List lists all command tables
func (repo *InMemoryCommandTableRepository) List() ([]*types.CommandTable, error) {
	tables := make([]*types.CommandTable, 0, len(repo.tables))
	for _, table := range repo.tables {
		tables = append(tables, table)
	}
	return tables, nil
}

// Exists checks if a command table exists
func (repo *InMemoryCommandTableRepository) Exists(id string) (bool, error) {
	_, exists := repo.tables[id]
	return exists, nil
}

// CommandTableService provides high-level service operations for command tables
type CommandTableService struct {
	manager   *CommandTableManager
	repository CommandTableRepository
}

// NewCommandTableService creates a new command table service
func NewCommandTableService(
	manager *CommandTableManager,
	repository CommandTableRepository,
) *CommandTableService {
	return &CommandTableService{
		manager:   manager,
		repository: repository,
	}
}

// CreateTable creates a new command table
func (cts *CommandTableService) CreateTable(id, name, description string) (*types.CommandTable, error) {
	table := cts.manager.CreateTable(id, name, description)
	if err := cts.repository.Save(table); err != nil {
		return nil, fmt.Errorf("failed to save table: %w", err)
	}
	return table, nil
}

// GetTable gets a command table by ID
func (cts *CommandTableService) GetTable(id string) (*types.CommandTable, error) {
	return cts.repository.Load(id)
}

// UpdateTable updates a command table
func (cts *CommandTableService) UpdateTable(table *types.CommandTable) error {
	if err := cts.manager.ValidateTable(table); err != nil {
		return fmt.Errorf("table validation failed: %w", err)
	}
	return cts.repository.Save(table)
}

// DeleteTable deletes a command table
func (cts *CommandTableService) DeleteTable(id string) error {
	return cts.repository.Delete(id)
}

// ListTables lists all command tables
func (cts *CommandTableService) ListTables() ([]*types.CommandTable, error) {
	return cts.repository.List()
}

// ExecuteTable executes a command table
func (cts *CommandTableService) ExecuteTable(ctx context.Context, id string) error {
	table, err := cts.repository.Load(id)
	if err != nil {
		return fmt.Errorf("failed to load table: %w", err)
	}
	return cts.manager.StartExecution(ctx, table)
}

// GetExecutionStatus returns the current execution status
func (cts *CommandTableService) GetExecutionStatus() execution.ExecutionStatus {
	return cts.manager.GetExecutionStatus()
}

// ControlExecution provides control over execution
func (cts *CommandTableService) ControlExecution() *ExecutionController {
	return &ExecutionController{
		manager: cts.manager,
	}
}

// ExecutionController provides control over command table execution
type ExecutionController struct {
	manager *CommandTableManager
}

// Start starts execution
func (ec *ExecutionController) Start(ctx context.Context, table *types.CommandTable) error {
	return ec.manager.StartExecution(ctx, table)
}

// Pause pauses execution
func (ec *ExecutionController) Pause() error {
	return ec.manager.PauseExecution()
}

// Resume resumes execution
func (ec *ExecutionController) Resume() error {
	return ec.manager.ResumeExecution()
}

// Stop stops execution
func (ec *ExecutionController) Stop() error {
	return ec.manager.StopExecution()
}

// GetStatus returns the execution status
func (ec *ExecutionController) GetStatus() execution.ExecutionStatus {
	return ec.manager.GetExecutionStatus()
}

// GetCurrentCommand returns the current command
func (ec *ExecutionController) GetCurrentCommand() *types.Command {
	return ec.manager.GetCurrentCommand()
}

// GetVariables returns the current variables
func (ec *ExecutionController) GetVariables() map[string]interface{} {
	return ec.manager.GetVariables()
}

// SetVariable sets a variable
func (ec *ExecutionController) SetVariable(name string, value interface{}) error {
	return ec.manager.SetVariable(name, value)
}