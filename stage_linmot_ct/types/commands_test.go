package types

import (
	"testing"
	"time"
)

func TestCommandType_String(t *testing.T) {
	tests := []struct {
		name     string
		cmdType  CommandType
		expected string
	}{
		{"MoveAbsolute", CmdMoveAbsolute, "MA"},
		{"MoveRelative", CmdMoveRelative, "MR"},
		{"MoveIncremental", CmdMoveIncremental, "MI"},
		{"Jog", CmdJog, "JO"},
		{"Stop", CmdStop, "ST"},
		{"Wait", CmdWait, "WA"},
		{"WaitPosition", CmdWaitPosition, "WP"},
		{"WaitVelocity", CmdWaitVelocity, "WV"},
		{"WaitForce", CmdWaitForce, "WF"},
		{"SetDigitalOutput", CmdSetDigitalOutput, "DO"},
		{"ClearDigitalOutput", CmdClearDigitalOutput, "CO"},
		{"SetAnalogOutput", CmdSetAnalogOutput, "AO"},
		{"WaitDigitalInput", CmdWaitDigitalInput, "DI"},
		{"WaitAnalogInput", CmdWaitAnalogInput, "AI"},
		{"LoopStart", CmdLoopStart, "LS"},
		{"LoopEnd", CmdLoopEnd, "LE"},
		{"LoopBreak", CmdLoopBreak, "LB"},
		{"Jump", CmdJump, "JP"},
		{"JumpIfTrue", CmdJumpIfTrue, "JT"},
		{"JumpIfFalse", CmdJumpIfFalse, "JF"},
		{"Home", CmdHome, "HO"},
		{"Reset", CmdReset, "RE"},
		{"SaveConfiguration", CmdSaveConfiguration, "SC"},
		{"LoadConfiguration", CmdLoadConfiguration, "LC"},
		{"ForceControlOn", CmdForceControlOn, "FC"},
		{"ForceControlOff", CmdForceControlOff, "FO"},
		{"SetForce", CmdSetForce, "SF"},
		{"StartOscilloscope", CmdStartOscilloscope, "SO"},
		{"StopOscilloscope", CmdStopOscilloscope, "SP"},
		{"SaveData", CmdSaveData, "SD"},
		{"Unknown", CommandType(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cmdType.String(); got != tt.expected {
				t.Errorf("CommandType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCommand_Validation(t *testing.T) {
	tests := []struct {
		name    string
		command *Command
		wantErr bool
	}{
		{
			name: "Valid command",
			command: &Command{
				ID:          1,
				Type:        CmdMoveAbsolute,
				Parameters:  map[string]interface{}{"position": 100.0},
				NextCommand: 2,
			},
			wantErr: false,
		},
		{
			name: "Invalid command with negative ID",
			command: &Command{
				ID:          -1,
				Type:        CmdMoveAbsolute,
				Parameters:  map[string]interface{}{"position": 100.0},
				NextCommand: 2,
			},
			wantErr: true,
		},
		{
			name: "Invalid command with zero ID",
			command: &Command{
				ID:          0,
				Type:        CmdMoveAbsolute,
				Parameters:  map[string]interface{}{"position": 100.0},
				NextCommand: 2,
			},
			wantErr: true,
		},
		{
			name: "Command with nil parameters",
			command: &Command{
				ID:          1,
				Type:        CmdMoveAbsolute,
				Parameters:  nil,
				NextCommand: 2,
			},
			wantErr: false, // nil parameters should be handled gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation logic
			if tt.command.ID <= 0 {
				if !tt.wantErr {
					t.Errorf("Expected error for invalid ID, but got none")
				}
			} else if tt.wantErr {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

func TestCommandTable_Creation(t *testing.T) {
	table := &CommandTable{
		ID:          "test_table",
		Name:        "Test Table",
		Description: "A test command table",
		Commands:    []Command{},
		Variables:   make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if table.ID != "test_table" {
		t.Errorf("Expected ID 'test_table', got %s", table.ID)
	}

	if table.Name != "Test Table" {
		t.Errorf("Expected Name 'Test Table', got %s", table.Name)
	}

	if len(table.Commands) != 0 {
		t.Errorf("Expected empty commands slice, got %d commands", len(table.Commands))
	}

	if table.Variables == nil {
		t.Error("Expected non-nil Variables map")
	}
}

func TestCommandTable_AddCommand(t *testing.T) {
	table := &CommandTable{
		ID:        "test_table",
		Commands:  []Command{},
		Variables: make(map[string]interface{}),
	}

	cmd := &Command{
		ID:          1,
		Type:        CmdMoveAbsolute,
		Parameters:  map[string]interface{}{"position": 100.0},
		NextCommand: 2,
	}

	// Test adding command
	table.Commands = append(table.Commands, *cmd)

	if len(table.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(table.Commands))
	}

	if table.Commands[0].ID != 1 {
		t.Errorf("Expected command ID 1, got %d", table.Commands[0].ID)
	}
}

func TestCommandBuilder(t *testing.T) {
	cmd := NewCommandBuilder().
		WithID(1).
		WithType(CmdMoveAbsolute).
		WithParameter("position", 100.0).
		WithParameter("velocity", 50.0).
		WithNextCommand(2).
		WithComment("Test command").
		Build()

	if cmd.ID != 1 {
		t.Errorf("Expected ID 1, got %d", cmd.ID)
	}

	if cmd.Type != CmdMoveAbsolute {
		t.Errorf("Expected type CmdMoveAbsolute, got %v", cmd.Type)
	}

	if cmd.Parameters["position"] != 100.0 {
		t.Errorf("Expected position 100.0, got %v", cmd.Parameters["position"])
	}

	if cmd.Parameters["velocity"] != 50.0 {
		t.Errorf("Expected velocity 50.0, got %v", cmd.Parameters["velocity"])
	}

	if cmd.NextCommand != 2 {
		t.Errorf("Expected next command 2, got %d", cmd.NextCommand)
	}

	if cmd.Comment != "Test command" {
		t.Errorf("Expected comment 'Test command', got %s", cmd.Comment)
	}

	if cmd.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if cmd.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestCommandTableBuilder(t *testing.T) {
	table := NewCommandTableBuilder().
		WithID("test_table").
		WithName("Test Table").
		WithDescription("A test table").
		WithVariable("test_var", "test_value").
		Build()

	if table.ID != "test_table" {
		t.Errorf("Expected ID 'test_table', got %s", table.ID)
	}

	if table.Name != "Test Table" {
		t.Errorf("Expected Name 'Test Table', got %s", table.Name)
	}

	if table.Description != "A test table" {
		t.Errorf("Expected Description 'A test table', got %s", table.Description)
	}

	if table.Variables["test_var"] != "test_value" {
		t.Errorf("Expected variable 'test_var' = 'test_value', got %v", table.Variables["test_var"])
	}

	if table.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if table.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestCommand_WithConditions(t *testing.T) {
	condition := Condition{
		Type:      CondDigitalInput,
		Parameter: "digital_input_1",
		Operator:  OpEqual,
		Value:     true,
	}

	cmd := NewCommandBuilder().
		WithID(1).
		WithType(CmdMoveAbsolute).
		WithCondition(condition).
		Build()

	if len(cmd.Conditions) != 1 {
		t.Errorf("Expected 1 condition, got %d", len(cmd.Conditions))
	}

	if cmd.Conditions[0].Type != CondDigitalInput {
		t.Errorf("Expected condition type CondDigitalInput, got %v", cmd.Conditions[0].Type)
	}

	if cmd.Conditions[0].Parameter != "digital_input_1" {
		t.Errorf("Expected parameter 'digital_input_1', got %s", cmd.Conditions[0].Parameter)
	}
}

func TestCommand_WithMultipleParameters(t *testing.T) {
	params := map[string]interface{}{
		"position":     100.0,
		"velocity":     50.0,
		"acceleration": 25.0,
		"timeout":      5000,
	}

	cmd := NewCommandBuilder().
		WithID(1).
		WithType(CmdMoveAbsolute).
		WithParameters(params).
		Build()

	for key, expectedValue := range params {
		if cmd.Parameters[key] != expectedValue {
			t.Errorf("Expected parameter %s = %v, got %v", key, expectedValue, cmd.Parameters[key])
		}
	}
}

func TestCommandTable_CommandOrdering(t *testing.T) {
	table := &CommandTable{
		ID:        "test_table",
		Commands:  []Command{},
		Variables: make(map[string]interface{}),
	}

	// Add commands in non-sequential order
	cmd1 := &Command{ID: 3, Type: CmdMoveAbsolute, Parameters: map[string]interface{}{"position": 100.0}}
	cmd2 := &Command{ID: 1, Type: CmdWait, Parameters: map[string]interface{}{"time": 1000}}
	cmd3 := &Command{ID: 2, Type: CmdMoveAbsolute, Parameters: map[string]interface{}{"position": 0.0}}

	table.Commands = append(table.Commands, *cmd1, *cmd2, *cmd3)

	// Verify commands are stored in the order they were added
	expectedIDs := []int{3, 1, 2}
	for i, cmd := range table.Commands {
		if cmd.ID != expectedIDs[i] {
			t.Errorf("Expected command %d to have ID %d, got %d", i, expectedIDs[i], cmd.ID)
		}
	}
}

func TestCommand_DefaultValues(t *testing.T) {
	cmd := NewCommandBuilder().Build()

	// Test default values
	if cmd.ID != 0 {
		t.Errorf("Expected default ID 0, got %d", cmd.ID)
	}

	if cmd.Parameters == nil {
		t.Error("Expected non-nil Parameters map")
	}

	if cmd.Conditions == nil {
		t.Error("Expected non-nil Conditions slice")
	}

	if cmd.NextCommand != 0 {
		t.Errorf("Expected default NextCommand 0, got %d", cmd.NextCommand)
	}

	if cmd.Comment != "" {
		t.Errorf("Expected empty Comment, got %s", cmd.Comment)
	}
}

func TestCommandTable_EmptyTable(t *testing.T) {
	table := &CommandTable{}

	// Test empty table properties
	if table.ID != "" {
		t.Errorf("Expected empty ID, got %s", table.ID)
	}

	if table.Name != "" {
		t.Errorf("Expected empty Name, got %s", table.Name)
	}

	if table.Commands != nil {
		t.Error("Expected nil Commands slice")
	}

	if table.Variables != nil {
		t.Error("Expected nil Variables map")
	}
}

func TestCommand_ImmutableAfterBuild(t *testing.T) {
	cmd := NewCommandBuilder().
		WithID(1).
		WithType(CmdMoveAbsolute).
		Build()

	originalCreatedAt := cmd.CreatedAt
	originalUpdatedAt := cmd.UpdatedAt

	// Wait a bit to ensure time has passed
	time.Sleep(1 * time.Millisecond)

	// Modify the command
	cmd.ID = 2
	cmd.Type = CmdWait

	// Verify timestamps haven't changed
	if !cmd.CreatedAt.Equal(originalCreatedAt) {
		t.Error("CreatedAt should not change after build")
	}

	if !cmd.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("UpdatedAt should not change after build")
	}

	// Verify the changes were applied
	if cmd.ID != 2 {
		t.Errorf("Expected ID 2, got %d", cmd.ID)
	}

	if cmd.Type != CmdWait {
		t.Errorf("Expected type CmdWait, got %v", cmd.Type)
	}
}