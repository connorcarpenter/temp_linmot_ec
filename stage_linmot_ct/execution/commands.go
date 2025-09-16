package execution

import (
	"context"
	"fmt"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// Motion command implementations
func (dee *DefaultExecutionEngine) executeMoveAbsolute(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	position, err := extractor.ExtractPosition(command.Parameters, "position")
	if err != nil {
		return fmt.Errorf("missing position parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units if needed
	posValue := dee.unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
	velValue := dee.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := dee.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := dee.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	return dee.driveController.MoveAbsolute(ctx, posValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

func (dee *DefaultExecutionEngine) executeMoveRelative(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	distance, err := extractor.ExtractPosition(command.Parameters, "distance")
	if err != nil {
		return fmt.Errorf("missing distance parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units if needed
	distValue := dee.unitConverter.ConvertPositionValue(distance, types.PositionUnitCounts)
	velValue := dee.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := dee.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := dee.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	return dee.driveController.MoveRelative(ctx, distValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

func (dee *DefaultExecutionEngine) executeMoveIncremental(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	distance, err := extractor.ExtractPosition(command.Parameters, "distance")
	if err != nil {
		return fmt.Errorf("missing distance parameter: %w", err)
	}
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	acceleration, err := extractor.ExtractAcceleration(command.Parameters, "acceleration")
	if err != nil {
		return fmt.Errorf("missing acceleration parameter: %w", err)
	}
	
	jerk, err := extractor.ExtractJerk(command.Parameters, "jerk")
	if err != nil {
		return fmt.Errorf("missing jerk parameter: %w", err)
	}
	
	// Convert to drive units if needed
	distValue := dee.unitConverter.ConvertPositionValue(distance, types.PositionUnitCounts)
	velValue := dee.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	accValue := dee.unitConverter.ConvertAccelerationValue(acceleration, types.AccelerationUnitCountsS2)
	jerkValue := dee.unitConverter.ConvertJerkValue(jerk, types.JerkUnitCountsS3)
	
	return dee.driveController.MoveIncremental(ctx, distValue.Value, velValue.Value, accValue.Value, jerkValue.Value)
}

func (dee *DefaultExecutionEngine) executeJog(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	// Convert to drive units if needed
	velValue := dee.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	
	return dee.driveController.Jog(ctx, velValue.Value)
}

func (dee *DefaultExecutionEngine) executeStop(ctx context.Context, command *types.Command) error {
	return dee.driveController.Stop(ctx)
}

// Wait command implementations
func (dee *DefaultExecutionEngine) executeWait(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	duration, err := extractor.ExtractTime(command.Parameters, "duration")
	if err != nil {
		return fmt.Errorf("missing duration parameter: %w", err)
	}
	
	return dee.driveController.Wait(ctx, duration.Duration())
}

func (dee *DefaultExecutionEngine) executeWaitPosition(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	position, err := extractor.ExtractPosition(command.Parameters, "position")
	if err != nil {
		return fmt.Errorf("missing position parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractPosition(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units if needed
	posValue := dee.unitConverter.ConvertPositionValue(position, types.PositionUnitCounts)
	tolValue := dee.unitConverter.ConvertPositionValue(tolerance, types.PositionUnitCounts)
	
	return dee.driveController.WaitPosition(ctx, posValue.Value, tolValue.Value, timeout.Duration())
}

func (dee *DefaultExecutionEngine) executeWaitVelocity(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	velocity, err := extractor.ExtractVelocity(command.Parameters, "velocity")
	if err != nil {
		return fmt.Errorf("missing velocity parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractVelocity(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units if needed
	velValue := dee.unitConverter.ConvertVelocityValue(velocity, types.VelocityUnitCountsS)
	tolValue := dee.unitConverter.ConvertVelocityValue(tolerance, types.VelocityUnitCountsS)
	
	return dee.driveController.WaitVelocity(ctx, velValue.Value, tolValue.Value, timeout.Duration())
}

func (dee *DefaultExecutionEngine) executeWaitForce(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	force, err := extractor.ExtractForce(command.Parameters, "force")
	if err != nil {
		return fmt.Errorf("missing force parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractForce(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	// Convert to drive units if needed
	forceValue := dee.unitConverter.ConvertForceValue(force, types.ForceUnitCounts)
	tolValue := dee.unitConverter.ConvertForceValue(tolerance, types.ForceUnitCounts)
	
	return dee.driveController.WaitForce(ctx, forceValue.Value, tolValue.Value, timeout.Duration())
}

// I/O command implementations
func (dee *DefaultExecutionEngine) executeSetDigitalOutput(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	output, err := extractor.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing output parameter: %w", err)
	}
	
	value, err := extractor.ExtractBool(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing value parameter: %w", err)
	}
	
	return dee.driveController.SetDigitalOutput(ctx, output, value)
}

func (dee *DefaultExecutionEngine) executeClearDigitalOutput(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	output, err := extractor.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing output parameter: %w", err)
	}
	
	return dee.driveController.ClearDigitalOutput(ctx, output)
}

func (dee *DefaultExecutionEngine) executeSetAnalogOutput(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	output, err := extractor.ExtractInt(command.Parameters, "output")
	if err != nil {
		return fmt.Errorf("missing output parameter: %w", err)
	}
	
	value, err := extractor.ExtractFloat(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing value parameter: %w", err)
	}
	
	return dee.driveController.SetAnalogOutput(ctx, output, value)
}

func (dee *DefaultExecutionEngine) executeWaitDigitalInput(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	input, err := extractor.ExtractInt(command.Parameters, "input")
	if err != nil {
		return fmt.Errorf("missing input parameter: %w", err)
	}
	
	value, err := extractor.ExtractBool(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing value parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	return dee.driveController.WaitDigitalInput(ctx, input, value, timeout.Duration())
}

func (dee *DefaultExecutionEngine) executeWaitAnalogInput(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	input, err := extractor.ExtractInt(command.Parameters, "input")
	if err != nil {
		return fmt.Errorf("missing input parameter: %w", err)
	}
	
	value, err := extractor.ExtractFloat(command.Parameters, "value")
	if err != nil {
		return fmt.Errorf("missing value parameter: %w", err)
	}
	
	tolerance, err := extractor.ExtractFloat(command.Parameters, "tolerance")
	if err != nil {
		return fmt.Errorf("missing tolerance parameter: %w", err)
	}
	
	timeout, err := extractor.ExtractTime(command.Parameters, "timeout")
	if err != nil {
		return fmt.Errorf("missing timeout parameter: %w", err)
	}
	
	return dee.driveController.WaitAnalogInput(ctx, input, value, tolerance, timeout.Duration())
}

// Loop and jump command implementations (simplified for now)
func (dee *DefaultExecutionEngine) executeLoopStart(ctx context.Context, command *types.Command) error {
	// TODO: Implement loop functionality
	return fmt.Errorf("loop commands not yet implemented")
}

func (dee *DefaultExecutionEngine) executeLoopEnd(ctx context.Context, command *types.Command) error {
	// TODO: Implement loop functionality
	return fmt.Errorf("loop commands not yet implemented")
}

func (dee *DefaultExecutionEngine) executeLoopBreak(ctx context.Context, command *types.Command) error {
	// TODO: Implement loop functionality
	return fmt.Errorf("loop commands not yet implemented")
}

func (dee *DefaultExecutionEngine) executeJump(ctx context.Context, command *types.Command) error {
	// TODO: Implement jump functionality
	return fmt.Errorf("jump commands not yet implemented")
}

func (dee *DefaultExecutionEngine) executeJumpIfTrue(ctx context.Context, command *types.Command) error {
	// TODO: Implement jump functionality
	return fmt.Errorf("jump commands not yet implemented")
}

func (dee *DefaultExecutionEngine) executeJumpIfFalse(ctx context.Context, command *types.Command) error {
	// TODO: Implement jump functionality
	return fmt.Errorf("jump commands not yet implemented")
}

// System command implementations
func (dee *DefaultExecutionEngine) executeHome(ctx context.Context, command *types.Command) error {
	return dee.driveController.Home(ctx)
}

func (dee *DefaultExecutionEngine) executeReset(ctx context.Context, command *types.Command) error {
	return dee.driveController.Reset(ctx)
}

func (dee *DefaultExecutionEngine) executeSaveConfiguration(ctx context.Context, command *types.Command) error {
	return dee.driveController.SaveConfiguration(ctx)
}

func (dee *DefaultExecutionEngine) executeLoadConfiguration(ctx context.Context, command *types.Command) error {
	return dee.driveController.LoadConfiguration(ctx)
}

// Force control command implementations
func (dee *DefaultExecutionEngine) executeForceControlOn(ctx context.Context, command *types.Command) error {
	return dee.driveController.ForceControlOn(ctx)
}

func (dee *DefaultExecutionEngine) executeForceControlOff(ctx context.Context, command *types.Command) error {
	return dee.driveController.ForceControlOff(ctx)
}

func (dee *DefaultExecutionEngine) executeSetForce(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	force, err := extractor.ExtractForce(command.Parameters, "force")
	if err != nil {
		return fmt.Errorf("missing force parameter: %w", err)
	}
	
	// Convert to drive units if needed
	forceValue := dee.unitConverter.ConvertForceValue(force, types.ForceUnitCounts)
	
	return dee.driveController.SetForce(ctx, forceValue.Value)
}

// Data acquisition command implementations
func (dee *DefaultExecutionEngine) executeStartOscilloscope(ctx context.Context, command *types.Command) error {
	return dee.driveController.StartOscilloscope(ctx)
}

func (dee *DefaultExecutionEngine) executeStopOscilloscope(ctx context.Context, command *types.Command) error {
	return dee.driveController.StopOscilloscope(ctx)
}

func (dee *DefaultExecutionEngine) executeSaveData(ctx context.Context, command *types.Command) error {
	extractor := types.NewParameterExtractor()
	
	filename, err := extractor.ExtractString(command.Parameters, "filename")
	if err != nil {
		return fmt.Errorf("missing filename parameter: %w", err)
	}
	
	return dee.driveController.SaveData(ctx, filename)
}