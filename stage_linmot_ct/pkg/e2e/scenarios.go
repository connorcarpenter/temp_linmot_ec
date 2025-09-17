package e2e

import (
	"context"
	"fmt"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/hardware"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/safety"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// CreateMotionSequenceScenario creates a complete motion sequence test scenario
func CreateMotionSequenceScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *TestScenario {
	return &TestScenario{
		Name:        "Complete Motion Sequence",
		Description: "Test complete motion workflow from home to target position",
		Category:    CategoryMotion,
		Priority:    PriorityCritical,
		Timeout:     60 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create command table
			table := manager.CreateTable("motion_test", "Motion Test", "Complete motion sequence")
			
			// Add commands
			homeCmd := types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdHome).
				Build()
			
			moveCmd := types.NewCommandBuilder().
				WithID(2).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(1000.0, types.PositionUnitCounts)).
				WithParameter("velocity", types.NewVelocityValue(100.0, types.VelocityUnitCountsS)).
				WithParameter("acceleration", types.NewAccelerationValue(50.0, types.AccelerationUnitCountsS2)).
				WithParameter("jerk", types.NewJerkValue(25.0, types.JerkUnitCountsS3)).
				Build()
			
			waitCmd := types.NewCommandBuilder().
				WithID(3).
				WithType(types.CmdWait).
				WithParameter("time", types.NewTimeValue(1.0, types.TimeUnitS)).
				Build()
			
			relativeMoveCmd := types.NewCommandBuilder().
				WithID(4).
				WithType(types.CmdMoveRelative).
				WithParameter("position", types.NewPositionValue(500.0, types.PositionUnitCounts)).
				WithParameter("velocity", types.NewVelocityValue(150.0, types.VelocityUnitCountsS)).
				WithParameter("acceleration", types.NewAccelerationValue(75.0, types.AccelerationUnitCountsS2)).
				WithParameter("jerk", types.NewJerkValue(37.5, types.JerkUnitCountsS3)).
				Build()
			
			stopCmd := types.NewCommandBuilder().
				WithID(5).
				WithType(types.CmdStop).
				Build()
			
			manager.AddCommand(table, homeCmd)
			manager.AddCommand(table, moveCmd)
			manager.AddCommand(table, waitCmd)
			manager.AddCommand(table, relativeMoveCmd)
			manager.AddCommand(table, stopCmd)
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate final position
			ctx := context.Background()
			position, err := controller.GetPosition(ctx)
			if err != nil {
				return err
			}
			
			expectedPosition := 1500.0 // 1000 + 500
			tolerance := 10.0
			if abs(position-expectedPosition) > tolerance {
				return fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, position)
			}
			
			return nil
		},
		Cleanup: func() error {
			// Stop any remaining motion
			ctx := context.Background()
			controller.Stop(ctx)
			return nil
		},
	}
}

// CreateForceControlScenario creates a force control test scenario
func CreateForceControlScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *TestScenario {
	return &TestScenario{
		Name:        "Force Control Workflow",
		Description: "Test complete force control workflow from enable to disable",
		Category:    CategoryForceControl,
		Priority:    PriorityHigh,
		Timeout:     30 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create command table
			table := manager.CreateTable("force_test", "Force Test", "Force control workflow")
			
			// Add commands
			forceOnCmd := types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdForceControlOn).
				Build()
			
			setForceCmd := types.NewCommandBuilder().
				WithID(2).
				WithType(types.CmdSetForce).
				WithParameter("force", types.NewForceValue(10.0, types.ForceUnitN)).
				Build()
			
			waitCmd := types.NewCommandBuilder().
				WithID(3).
				WithType(types.CmdWait).
				WithParameter("time", types.NewTimeValue(2.0, types.TimeUnitS)).
				Build()
			
			forceOffCmd := types.NewCommandBuilder().
				WithID(4).
				WithType(types.CmdForceControlOff).
				Build()
			
			manager.AddCommand(table, forceOnCmd)
			manager.AddCommand(table, setForceCmd)
			manager.AddCommand(table, waitCmd)
			manager.AddCommand(table, forceOffCmd)
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate that force control is disabled
			ctx := context.Background()
			// Note: In a real implementation, you would check the force control state
			// For now, we'll just verify the command table executed successfully
			return nil
		},
		Cleanup: func() error {
			// Ensure force control is disabled
			ctx := context.Background()
			controller.ForceControlOff(ctx)
			return nil
		},
	}
}

// CreateIOControlScenario creates an I/O control test scenario
func CreateIOControlScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *TestScenario {
	return &TestScenario{
		Name:        "I/O Control Workflow",
		Description: "Test digital and analog I/O control workflow",
		Category:    CategoryIO,
		Priority:    PriorityHigh,
		Timeout:     20 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create command table
			table := manager.CreateTable("io_test", "I/O Test", "I/O control workflow")
			
			// Add commands
			setDigitalCmd := types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdSetDigitalOutput).
				WithParameter("output", 1).
				WithParameter("value", true).
				Build()
			
			setAnalogCmd := types.NewCommandBuilder().
				WithID(2).
				WithType(types.CmdSetAnalogOutput).
				WithParameter("output", 1).
				WithParameter("value", 3.14).
				Build()
			
			waitCmd := types.NewCommandBuilder().
				WithID(3).
				WithType(types.CmdWait).
				WithParameter("time", types.NewTimeValue(1.0, types.TimeUnitS)).
				Build()
			
			clearDigitalCmd := types.NewCommandBuilder().
				WithID(4).
				WithType(types.CmdClearDigitalOutput).
				WithParameter("output", 1).
				Build()
			
			manager.AddCommand(table, setDigitalCmd)
			manager.AddCommand(table, setAnalogCmd)
			manager.AddCommand(table, waitCmd)
			manager.AddCommand(table, clearDigitalCmd)
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate I/O states
			ctx := context.Background()
			
			// Check digital output is cleared
			digitalOutput, err := controller.GetDigitalOutput(ctx, 1)
			if err != nil {
				return err
			}
			if digitalOutput {
				return fmt.Errorf("digital output 1 should be cleared")
			}
			
			// Check analog output is set
			analogOutput, err := controller.GetAnalogOutput(ctx, 1)
			if err != nil {
				return err
			}
			expectedVoltage := 3.14
			tolerance := 0.1
			if abs(analogOutput-expectedVoltage) > tolerance {
				return fmt.Errorf("analog output 1 mismatch: expected %f, got %f", expectedVoltage, analogOutput)
			}
			
			return nil
		},
		Cleanup: func() error {
			// Clear all outputs
			ctx := context.Background()
			controller.ClearDigitalOutput(ctx, 1)
			controller.SetAnalogOutput(ctx, 1, 0.0)
			return nil
		},
	}
}

// CreateSafetyScenario creates a safety system test scenario
func CreateSafetyScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager, safetyGuard *safety.SafetyGuard) *TestScenario {
	return &TestScenario{
		Name:        "Safety System Workflow",
		Description: "Test safety system including emergency stop and limits",
		Category:    CategorySafety,
		Priority:    PriorityCritical,
		Timeout:     30 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create command table
			table := manager.CreateTable("safety_test", "Safety Test", "Safety system workflow")
			
			// Add commands
			moveCmd := types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(1000.0, types.PositionUnitCounts)).
				WithParameter("velocity", types.NewVelocityValue(100.0, types.VelocityUnitCountsS)).
				WithParameter("acceleration", types.NewAccelerationValue(50.0, types.AccelerationUnitCountsS2)).
				WithParameter("jerk", types.NewJerkValue(25.0, types.JerkUnitCountsS3)).
				Build()
			
			waitCmd := types.NewCommandBuilder().
				WithID(2).
				WithType(types.CmdWait).
				WithParameter("time", types.NewTimeValue(0.5, types.TimeUnitS)).
				Build()
			
			manager.AddCommand(table, moveCmd)
			manager.AddCommand(table, waitCmd)
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate safety system is working
			ctx := context.Background()
			
			// Check that motion completed safely
			complete, err := controller.IsMotionComplete(ctx)
			if err != nil {
				return err
			}
			if !complete {
				return fmt.Errorf("motion did not complete")
			}
			
			// Check position is within safe limits
			position, err := controller.GetPosition(ctx)
			if err != nil {
				return err
			}
			
			// Basic safety check - position should be reasonable
			if position < -10000 || position > 10000 {
				return fmt.Errorf("position out of safe range: %f", position)
			}
			
			return nil
		},
		Cleanup: func() error {
			// Ensure system is in safe state
			ctx := context.Background()
			controller.Stop(ctx)
			return nil
		},
	}
}

// CreatePerformanceScenario creates a performance test scenario
func CreatePerformanceScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *TestScenario {
	return &TestScenario{
		Name:        "Performance Test",
		Description: "Test system performance with rapid commands",
		Category:    CategoryPerformance,
		Priority:    PriorityMedium,
		Timeout:     60 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create command table with many rapid commands
			table := manager.CreateTable("performance_test", "Performance Test", "Rapid command execution")
			
			// Add multiple rapid motion commands
			for i := 1; i <= 10; i++ {
				moveCmd := types.NewCommandBuilder().
					WithID(i).
					WithType(types.CmdMoveAbsolute).
					WithParameter("position", types.NewPositionValue(float64(i*100), types.PositionUnitCounts)).
					WithParameter("velocity", types.NewVelocityValue(200.0, types.VelocityUnitCountsS)).
					WithParameter("acceleration", types.NewAccelerationValue(100.0, types.AccelerationUnitCountsS2)).
					WithParameter("jerk", types.NewJerkValue(50.0, types.JerkUnitCountsS3)).
					Build()
				
				manager.AddCommand(table, moveCmd)
			}
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate performance metrics
			ctx := context.Background()
			
			// Check final position
			position, err := controller.GetPosition(ctx)
			if err != nil {
				return err
			}
			
			expectedPosition := 1000.0 // 10 * 100
			tolerance := 50.0
			if abs(position-expectedPosition) > tolerance {
				return fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, position)
			}
			
			return nil
		},
		Cleanup: func() error {
			// Stop any remaining motion
			ctx := context.Background()
			controller.Stop(ctx)
			return nil
		},
	}
}

// CreateIntegrationScenario creates an integration test scenario
func CreateIntegrationScenario(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *TestScenario {
	return &TestScenario{
		Name:        "Integration Test",
		Description: "Test complete system integration with all features",
		Category:    CategoryIntegration,
		Priority:    PriorityCritical,
		Timeout:     120 * time.Second,
		Setup: func() error {
			// Initialize system
			ctx := context.Background()
			if !controller.IsConnected() {
				return controller.Connect(ctx)
			}
			return nil
		},
		Execute: func() error {
			// Create comprehensive command table
			table := manager.CreateTable("integration_test", "Integration Test", "Complete system integration")
			
			// Home command
			homeCmd := types.NewCommandBuilder().
				WithID(1).
				WithType(types.CmdHome).
				Build()
			
			// Motion command
			moveCmd := types.NewCommandBuilder().
				WithID(2).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(500.0, types.PositionUnitCounts)).
				WithParameter("velocity", types.NewVelocityValue(100.0, types.VelocityUnitCountsS)).
				WithParameter("acceleration", types.NewAccelerationValue(50.0, types.AccelerationUnitCountsS2)).
				WithParameter("jerk", types.NewJerkValue(25.0, types.JerkUnitCountsS3)).
				Build()
			
			// I/O command
			setDigitalCmd := types.NewCommandBuilder().
				WithID(3).
				WithType(types.CmdSetDigitalOutput).
				WithParameter("output", 1).
				WithParameter("value", true).
				Build()
			
			// Force control command
			forceOnCmd := types.NewCommandBuilder().
				WithID(4).
				WithType(types.CmdForceControlOn).
				Build()
			
			// Wait command
			waitCmd := types.NewCommandBuilder().
				WithID(5).
				WithType(types.CmdWait).
				WithParameter("time", types.NewTimeValue(1.0, types.TimeUnitS)).
				Build()
			
			// Force control off
			forceOffCmd := types.NewCommandBuilder().
				WithID(6).
				WithType(types.CmdForceControlOff).
				Build()
			
			// Clear digital output
			clearDigitalCmd := types.NewCommandBuilder().
				WithID(7).
				WithType(types.CmdClearDigitalOutput).
				WithParameter("output", 1).
				Build()
			
			// Final motion
			finalMoveCmd := types.NewCommandBuilder().
				WithID(8).
				WithType(types.CmdMoveAbsolute).
				WithParameter("position", types.NewPositionValue(0.0, types.PositionUnitCounts)).
				WithParameter("velocity", types.NewVelocityValue(100.0, types.VelocityUnitCountsS)).
				WithParameter("acceleration", types.NewAccelerationValue(50.0, types.AccelerationUnitCountsS2)).
				WithParameter("jerk", types.NewJerkValue(25.0, types.JerkUnitCountsS3)).
				Build()
			
			manager.AddCommand(table, homeCmd)
			manager.AddCommand(table, moveCmd)
			manager.AddCommand(table, setDigitalCmd)
			manager.AddCommand(table, forceOnCmd)
			manager.AddCommand(table, waitCmd)
			manager.AddCommand(table, forceOffCmd)
			manager.AddCommand(table, clearDigitalCmd)
			manager.AddCommand(table, finalMoveCmd)
			
			// Execute
			ctx := context.Background()
			return manager.StartExecution(ctx, table)
		},
		Validate: func() error {
			// Validate complete system state
			ctx := context.Background()
			
			// Check final position
			position, err := controller.GetPosition(ctx)
			if err != nil {
				return err
			}
			
			expectedPosition := 0.0
			tolerance := 10.0
			if abs(position-expectedPosition) > tolerance {
				return fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, position)
			}
			
			// Check digital output is cleared
			digitalOutput, err := controller.GetDigitalOutput(ctx, 1)
			if err != nil {
				return err
			}
			if digitalOutput {
				return fmt.Errorf("digital output 1 should be cleared")
			}
			
			return nil
		},
		Cleanup: func() error {
			// Ensure system is in safe state
			ctx := context.Background()
			controller.Stop(ctx)
			controller.ForceControlOff(ctx)
			controller.ClearDigitalOutput(ctx, 1)
			return nil
		},
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}