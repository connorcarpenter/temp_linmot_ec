package commands

import (
	"context"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewCommandRegistry(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	
	registry := NewCommandRegistry(driveController, unitConverter)
	
	if registry == nil {
		t.Fatal("NewCommandRegistry returned nil")
	}
	
	// Check that motion commands are registered
	motionCommands := []types.CommandType{
		types.CmdMoveAbsolute,
		types.CmdMoveRelative,
		types.CmdMoveIncremental,
		types.CmdJog,
		types.CmdStop,
	}
	
	for _, cmdType := range motionCommands {
		if !registry.IsCommandSupported(cmdType) {
			t.Errorf("Motion command %s not registered", cmdType)
		}
	}
	
	// Check that control commands are registered
	controlCommands := []types.CommandType{
		types.CmdWait,
		types.CmdWaitPosition,
		types.CmdWaitVelocity,
		types.CmdWaitForce,
	}
	
	for _, cmdType := range controlCommands {
		if !registry.IsCommandSupported(cmdType) {
			t.Errorf("Control command %s not registered", cmdType)
		}
	}
	
	// Check total count
	expectedCount := len(motionCommands) + len(controlCommands)
	if registry.GetCommandCount() != expectedCount {
		t.Errorf("Expected %d commands, got %d", expectedCount, registry.GetCommandCount())
	}
}

func TestCommandRegistry_RegisterExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	// Create a custom executor
	customExecutor := NewMotionCommandExecutor(driveController, unitConverter)
	
	// Register a new command type (using a non-standard type for testing)
	registry.RegisterExecutor(types.CmdHome, customExecutor)
	
	if !registry.IsCommandSupported(types.CmdHome) {
		t.Error("Custom command not registered")
	}
	
	if registry.GetCommandCount() != 10 { // 9 original + 1 custom
		t.Errorf("Expected 10 commands, got %d", registry.GetCommandCount())
	}
}

func TestCommandRegistry_GetExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	// Test existing command
	executor, err := registry.GetExecutor(types.CmdMoveAbsolute)
	if err != nil {
		t.Fatalf("GetExecutor failed for existing command: %v", err)
	}
	
	if executor == nil {
		t.Fatal("GetExecutor returned nil executor")
	}
	
	// Test non-existing command
	_, err = registry.GetExecutor(types.CmdHome)
	if err == nil {
		t.Fatal("Expected error for non-existing command, got nil")
	}
	
	if err.Error() != "no executor registered for command type: HO" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestCommandRegistry_ExecuteCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	// Create a move absolute command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdMoveAbsolute).
		WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
		WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
		WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
		WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
		Build()
	
	// Execute the command
	err := registry.ExecuteCommand(context.Background(), command)
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v", err)
	}
	
	// Verify the drive controller was called
	expectedPosition := 10000.0 // 10.0 mm * 1000 counts/mm
	if driveController.position != expectedPosition {
		t.Errorf("Expected position %f, got %f", expectedPosition, driveController.position)
	}
}

func TestCommandRegistry_ExecuteCommand_UnsupportedCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	// Create an unsupported command
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdHome).
		Build()
	
	// Execute the command
	err := registry.ExecuteCommand(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error for unsupported command, got nil")
	}
	
	if err.Error() != "no executor registered for command type: HO" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestCommandRegistry_ValidateCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid move absolute command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(10.0, types.PositionUnitMM)).
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
				WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
				Build(),
			wantErr: false,
		},
		{
			name: "Invalid move absolute command (missing position)",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("velocity", types.NewVelocityValue(5.0, types.VelocityUnitMMS)).
				WithParameter("acceleration", types.NewAccelerationValue(10.0, types.AccelerationUnitMMS2)).
				WithParameter("jerk", types.NewJerkValue(100.0, types.JerkUnitMMS3)).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid wait command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWait).
				WithParameter("duration", types.NewTimeValue(100.0, types.TimeUnitMS)).
				Build(),
			wantErr: false,
		},
		{
			name: "Invalid wait command (missing duration)",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdWait).
				Build(),
			wantErr: true,
		},
		{
			name: "Unsupported command type",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdHome).
				Build(),
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.ValidateCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRegistry_GetCommandInfo(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	tests := []struct {
		commandType types.CommandType
		wantErr     bool
	}{
		{
			commandType: types.CmdMoveAbsolute,
			wantErr:     false,
		},
		{
			commandType: types.CmdWait,
			wantErr:     false,
		},
		{
			commandType: types.CmdHome, // Unsupported
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			description, parameters, err := registry.GetCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if description == "" {
					t.Error("Expected non-empty description")
				}
				if parameters == nil {
					t.Error("Expected non-nil parameters")
				}
			}
		})
	}
}

func TestCommandRegistry_GetSupportedCommandTypes(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	commandTypes := registry.GetSupportedCommandTypes()
	
	if len(commandTypes) != 9 { // 5 motion + 4 control commands
		t.Errorf("Expected 9 command types, got %d", len(commandTypes))
	}
	
	// Check that all expected types are present
	expectedTypes := map[types.CommandType]bool{
		types.CmdMoveAbsolute:    true,
		types.CmdMoveRelative:    true,
		types.CmdMoveIncremental: true,
		types.CmdJog:             true,
		types.CmdStop:            true,
		types.CmdWait:            true,
		types.CmdWaitPosition:    true,
		types.CmdWaitVelocity:    true,
		types.CmdWaitForce:       true,
	}
	
	for _, cmdType := range commandTypes {
		if !expectedTypes[cmdType] {
			t.Errorf("Unexpected command type: %s", cmdType)
		}
	}
}

func TestCommandRegistry_ListCommandInfo(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter)
	
	info := registry.ListCommandInfo()
	
	if len(info) != 9 { // 5 motion + 4 control commands
		t.Errorf("Expected 9 command info entries, got %d", len(info))
	}
	
	// Check that all entries have valid information
	for cmdType, cmdInfo := range info {
		if cmdInfo.Type != cmdType {
			t.Errorf("Command type mismatch: expected %s, got %s", cmdType, cmdInfo.Type)
		}
		if cmdInfo.Description == "" {
			t.Errorf("Empty description for command type: %s", cmdType)
		}
		if cmdInfo.Parameters == nil {
			t.Errorf("Nil parameters for command type: %s", cmdType)
		}
	}
}

func TestCommandInfo_String(t *testing.T) {
	info := CommandInfo{
		Type:        types.CmdMoveAbsolute,
		Name:        "MoveAbsolute",
		Description: "Move to absolute position",
		Parameters:  []string{"position", "velocity", "acceleration", "jerk"},
	}
	
	str := info.String()
	expected := "MoveAbsolute: Move to absolute position (parameters: [position velocity acceleration jerk])"
	
	if str != expected {
		t.Errorf("String() = %v, want %v", str, expected)
	}
}