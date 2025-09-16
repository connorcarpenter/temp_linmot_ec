package commands

import (
	"context"
	"testing"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewLoopCommandExecutor(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	if executor == nil {
		t.Fatal("Expected non-nil executor")
	}

	if executor.driveController != driveController {
		t.Error("DriveController not set correctly")
	}

	if executor.unitConverter != unitConverter {
		t.Error("UnitConverter not set correctly")
	}
}

func TestLoopCommandExecutor_ExecuteLoopStart(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoopStart).
		WithParameter("count", 5).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteLoopStart_InvalidCount(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoopStart).
		WithParameter("count", 0). // Invalid count
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "loop count must be positive, got 0" {
		t.Errorf("Expected invalid count error, got: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteLoopEnd(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoopEnd).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteLoopBreak(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdLoopBreak).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteJump(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdJump).
		WithParameter("target_id", 10).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteJump_InvalidTargetID(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdJump).
		WithParameter("target_id", 0). // Invalid target ID
		Build()

	err := executor.Execute(context.Background(), command)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "target command ID must be positive, got 0" {
		t.Errorf("Expected invalid target ID error, got: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteJumpIfTrue(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdJumpIfTrue).
		WithParameter("target_id", 10).
		WithParameter("condition", true).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ExecuteJumpIfFalse(t *testing.T) {
	driveController := NewMockDriveController()
	unitConverter := types.NewUnitConverter()
	executor := NewLoopCommandExecutor(driveController, unitConverter)

	command := types.NewCommandBuilder().
		WithID(1).
		WithType(types.CmdJumpIfFalse).
		WithParameter("target_id", 10).
		WithParameter("condition", false).
		Build()

	err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

func TestLoopCommandExecutor_ValidateLoopCommand(t *testing.T) {
	executor := NewLoopCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		name    string
		command *types.Command
		wantErr bool
	}{
		{
			name: "Valid loop start command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdLoopStart).
				WithParameter("count", 5).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing count parameter for loop start",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdLoopStart).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid loop end command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdLoopEnd).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid loop break command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdLoopBreak).
				Build(),
			wantErr: false,
		},
		{
			name: "Valid jump command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJump).
				WithParameter("target_id", 10).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing target_id parameter for jump",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJump).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid jump if true command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJumpIfTrue).
				WithParameter("target_id", 10).
				WithParameter("condition", true).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing condition parameter for jump if true",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJumpIfTrue).
				WithParameter("target_id", 10).
				Build(),
			wantErr: true,
		},
		{
			name: "Valid jump if false command",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJumpIfFalse).
				WithParameter("target_id", 10).
				WithParameter("condition", false).
				Build(),
			wantErr: false,
		},
		{
			name: "Missing condition parameter for jump if false",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdJumpIfFalse).
				WithParameter("target_id", 10).
				Build(),
			wantErr: true,
		},
		{
			name: "Unsupported command type",
			command: types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				Build(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.ValidateLoopCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoopCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoopCommandExecutor_GetLoopCommandInfo(t *testing.T) {
	executor := NewLoopCommandExecutor(NewMockDriveController(), types.NewUnitConverter())

	tests := []struct {
		commandType types.CommandType
		wantName    string
		wantParams  []string
		wantErr     bool
	}{
		{
			commandType: types.CmdLoopStart,
			wantName:    "LoopStart",
			wantParams:  []string{"count"},
			wantErr:     false,
		},
		{
			commandType: types.CmdLoopEnd,
			wantName:    "LoopEnd",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdLoopBreak,
			wantName:    "LoopBreak",
			wantParams:  []string{},
			wantErr:     false,
		},
		{
			commandType: types.CmdJump,
			wantName:    "Jump",
			wantParams:  []string{"target_id"},
			wantErr:     false,
		},
		{
			commandType: types.CmdJumpIfTrue,
			wantName:    "JumpIfTrue",
			wantParams:  []string{"target_id", "condition"},
			wantErr:     false,
		},
		{
			commandType: types.CmdJumpIfFalse,
			wantName:    "JumpIfFalse",
			wantParams:  []string{"target_id", "condition"},
			wantErr:     false,
		},
		{
			commandType: types.CmdMoveAbsolute, // Unsupported
			wantName:    "",
			wantParams:  nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.commandType.String(), func(t *testing.T) {
			name, params, err := executor.GetLoopCommandInfo(tt.commandType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoopCommandInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if name != tt.wantName {
				t.Errorf("GetLoopCommandInfo() name = %v, want %v", name, tt.wantName)
			}
			if len(params) != len(tt.wantParams) {
				t.Errorf("GetLoopCommandInfo() params length = %v, want %v", len(params), len(tt.wantParams))
			}
		})
	}
}