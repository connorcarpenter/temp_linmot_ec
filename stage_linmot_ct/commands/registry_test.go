package commands

import (
	"context"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewCommandRegistry(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
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
	expectedCount := 30 // 5 motion + 4 control + 5 I/O + 6 loop/jump + 4 system + 3 force + 3 data
	if registry.GetCommandCount() != expectedCount {
		t.Errorf("Expected %d commands, got %d", expectedCount, registry.GetCommandCount())
	}
}

func TestCommandRegistry_RegisterExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
	// Create a custom executor
	customExecutor := NewMotionCommandExecutor(driveController, unitConverter, nil)
	
	// Register a new command type (using a non-standard type for testing)
	customCommandType := types.CommandType(999)
	registry.RegisterExecutor(customCommandType, customExecutor)
	
	if !registry.IsCommandSupported(customCommandType) {
		t.Error("Custom command not registered")
	}
	
	if registry.GetCommandCount() != 31 { // 30 original + 1 custom
		t.Errorf("Expected 31 commands, got %d", registry.GetCommandCount())
	}
}

func TestCommandRegistry_GetExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
	// Test existing command
	executor, err := registry.GetExecutor(types.CmdMoveAbsolute)
	if err != nil {
		t.Fatalf("GetExecutor failed for existing command: %v", err)
	}
	
	if executor == nil {
		t.Fatal("GetExecutor returned nil executor")
	}
	
	// Test non-existing command (use a command type that doesn't exist)
	nonExistentCommand := types.CommandType(999)
	_, err = registry.GetExecutor(nonExistentCommand)
	if err == nil {
		t.Fatal("Expected error for non-existing command, got nil")
	}
	
	if err.Error() != "no executor registered for command type: UNKNOWN" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestCommandRegistry_ExecuteCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
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
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
	// Create an unsupported command (use a command type that doesn't exist)
	unsupportedCommandType := types.CommandType(999)
	command := types.NewCommandBuilder().
		WithID(1).
		WithType(unsupportedCommandType).
		Build()
	
	// Execute the command
	err := registry.ExecuteCommand(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error for unsupported command, got nil")
	}
	
	if err.Error() != "no executor registered for command type: UNKNOWN" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestCommandRegistry_ValidateCommand(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
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
				WithType(types.CommandType(999)).
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
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
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
			commandType: types.CommandType(999), // Unsupported
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
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
	commandTypes := registry.GetSupportedCommandTypes()
	
	if len(commandTypes) != 30 { // 5 motion + 4 control + 5 I/O + 6 loop/jump + 4 system + 3 force + 3 data
		t.Errorf("Expected 30 command types, got %d", len(commandTypes))
	}
	
	// Check that all expected types are present
	expectedTypes := map[types.CommandType]bool{
		// Motion commands
		types.CmdMoveAbsolute:    true,
		types.CmdMoveRelative:    true,
		types.CmdMoveIncremental: true,
		types.CmdJog:             true,
		types.CmdStop:            true,
		// Control commands
		types.CmdWait:            true,
		types.CmdWaitPosition:    true,
		types.CmdWaitVelocity:    true,
		types.CmdWaitForce:       true,
		// I/O commands
		types.CmdSetDigitalOutput:   true,
		types.CmdClearDigitalOutput: true,
		types.CmdSetAnalogOutput:    true,
		types.CmdWaitDigitalInput:   true,
		types.CmdWaitAnalogInput:    true,
		// Loop/Jump commands
		types.CmdLoopStart:    true,
		types.CmdLoopEnd:      true,
		types.CmdLoopBreak:    true,
		types.CmdJump:         true,
		types.CmdJumpIfTrue:   true,
		types.CmdJumpIfFalse:  true,
		// System commands
		types.CmdHome:               true,
		types.CmdReset:              true,
		types.CmdSaveConfiguration:  true,
		types.CmdLoadConfiguration:  true,
		// Force Control commands
		types.CmdForceControlOn:     true,
		types.CmdForceControlOff:    true,
		types.CmdSetForce:           true,
		// Data Acquisition commands
		types.CmdStartOscilloscope:  true,
		types.CmdStopOscilloscope:   true,
		types.CmdSaveData:           true,
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
	registry := NewCommandRegistry(driveController, unitConverter, nil)
	
	info := registry.ListCommandInfo()
	
	if len(info) != 30 { // 5 motion + 4 control + 5 I/O + 6 loop/jump + 4 system + 3 force + 3 data
		t.Errorf("Expected 30 command info entries, got %d", len(info))
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