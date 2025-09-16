package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// executeCommands executes the command table
func (dee *DefaultExecutionEngine) executeCommands(ctx context.Context) {
	defer func() {
		dee.mu.Lock()
		if dee.state == StateRunning || dee.state == StatePaused {
			dee.state = StateCompleted
			dee.endTime = time.Now()
		}
		dee.mu.Unlock()
		close(dee.doneChan)
	}()
	
	for dee.currentCommand < len(dee.commandTable.Commands) {
		// Check for stop signal
		select {
		case <-dee.stopChan:
			dee.mu.Lock()
			dee.state = StateStopped
			dee.endTime = time.Now()
			dee.mu.Unlock()
			return
		case <-ctx.Done():
			dee.mu.Lock()
			dee.state = StateError
			dee.error = ctx.Err()
			dee.endTime = time.Now()
			dee.mu.Unlock()
			return
		default:
		}
		
		// Check for pause
		select {
		case <-dee.pauseChan:
			dee.mu.Lock()
			dee.state = StatePaused
			dee.mu.Unlock()
			
			// Wait for resume or stop
			select {
			case <-dee.resumeChan:
				dee.mu.Lock()
				dee.state = StateRunning
				dee.mu.Unlock()
			case <-dee.stopChan:
				dee.mu.Lock()
				dee.state = StateStopped
				dee.endTime = time.Now()
				dee.mu.Unlock()
				return
			case <-ctx.Done():
				dee.mu.Lock()
				dee.state = StateError
				dee.error = ctx.Err()
				dee.endTime = time.Now()
				dee.mu.Unlock()
				return
			}
		default:
		}
		
		// Execute current command
		command := dee.commandTable.Commands[dee.currentCommand]
		result := dee.executeCommand(ctx, &command)
		
		// Store result
		dee.mu.Lock()
		dee.results = append(dee.results, result)
		dee.mu.Unlock()
		
		// Check if command failed
		if !result.Success {
			dee.mu.Lock()
			dee.state = StateError
			dee.error = result.Error
			dee.endTime = time.Now()
			dee.mu.Unlock()
			return
		}
		
		// Move to next command
		dee.mu.Lock()
		dee.currentCommand++
		dee.mu.Unlock()
	}
}

// executeCommand executes a single command
func (dee *DefaultExecutionEngine) executeCommand(ctx context.Context, command *types.Command) ExecutionResult {
	startTime := time.Now()
	result := ExecutionResult{
		CommandID: command.ID,
		Success:   false,
		Timestamp: startTime,
	}
	
	// Validate command
	if err := command.Validate(); err != nil {
		result.Error = fmt.Errorf("command validation failed: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	
	// Check conditions
	if len(command.Conditions) > 0 {
		// Convert []Condition to []*Condition
		conditions := make([]*types.Condition, len(command.Conditions))
		for i := range command.Conditions {
			conditions[i] = &command.Conditions[i]
		}
		
		canExecute, err := dee.evaluateConditions(ctx, conditions)
		if err != nil {
			result.Error = fmt.Errorf("condition evaluation failed: %w", err)
			result.Duration = time.Since(startTime)
			return result
		}
		
		if !canExecute {
			result.Success = true // Command skipped due to conditions
			result.Duration = time.Since(startTime)
			return result
		}
	}
	
	// Execute command based on type
	var err error
	switch command.Type {
	case types.CmdMoveAbsolute:
		err = dee.executeMoveAbsolute(ctx, command)
	case types.CmdMoveRelative:
		err = dee.executeMoveRelative(ctx, command)
	case types.CmdMoveIncremental:
		err = dee.executeMoveIncremental(ctx, command)
	case types.CmdJog:
		err = dee.executeJog(ctx, command)
	case types.CmdStop:
		err = dee.executeStop(ctx, command)
	case types.CmdWait:
		err = dee.executeWait(ctx, command)
	case types.CmdWaitPosition:
		err = dee.executeWaitPosition(ctx, command)
	case types.CmdWaitVelocity:
		err = dee.executeWaitVelocity(ctx, command)
	case types.CmdWaitForce:
		err = dee.executeWaitForce(ctx, command)
	case types.CmdSetDigitalOutput:
		err = dee.executeSetDigitalOutput(ctx, command)
	case types.CmdClearDigitalOutput:
		err = dee.executeClearDigitalOutput(ctx, command)
	case types.CmdSetAnalogOutput:
		err = dee.executeSetAnalogOutput(ctx, command)
	case types.CmdWaitDigitalInput:
		err = dee.executeWaitDigitalInput(ctx, command)
	case types.CmdWaitAnalogInput:
		err = dee.executeWaitAnalogInput(ctx, command)
	case types.CmdLoopStart:
		err = dee.executeLoopStart(ctx, command)
	case types.CmdLoopEnd:
		err = dee.executeLoopEnd(ctx, command)
	case types.CmdLoopBreak:
		err = dee.executeLoopBreak(ctx, command)
	case types.CmdJump:
		err = dee.executeJump(ctx, command)
	case types.CmdJumpIfTrue:
		err = dee.executeJumpIfTrue(ctx, command)
	case types.CmdJumpIfFalse:
		err = dee.executeJumpIfFalse(ctx, command)
	case types.CmdHome:
		err = dee.executeHome(ctx, command)
	case types.CmdReset:
		err = dee.executeReset(ctx, command)
	case types.CmdSaveConfiguration:
		err = dee.executeSaveConfiguration(ctx, command)
	case types.CmdLoadConfiguration:
		err = dee.executeLoadConfiguration(ctx, command)
	case types.CmdForceControlOn:
		err = dee.executeForceControlOn(ctx, command)
	case types.CmdForceControlOff:
		err = dee.executeForceControlOff(ctx, command)
	case types.CmdSetForce:
		err = dee.executeSetForce(ctx, command)
	case types.CmdStartOscilloscope:
		err = dee.executeStartOscilloscope(ctx, command)
	case types.CmdStopOscilloscope:
		err = dee.executeStopOscilloscope(ctx, command)
	case types.CmdSaveData:
		err = dee.executeSaveData(ctx, command)
	default:
		err = fmt.Errorf("unsupported command type: %s", command.Type)
	}
	
	result.Success = err == nil
	result.Error = err
	result.Duration = time.Since(startTime)
	
	return result
}

// evaluateConditions evaluates command conditions
func (dee *DefaultExecutionEngine) evaluateConditions(ctx context.Context, conditions []*types.Condition) (bool, error) {
	if dee.conditionEvaluator == nil {
		return true, nil // No evaluator means all conditions pass
	}
	
	for _, condition := range conditions {
		canEvaluate := dee.conditionEvaluator.CanEvaluate(condition)
		
		if !canEvaluate {
			return false, nil
		}
		
		result, err := dee.conditionEvaluator.Evaluate(ctx, condition, dee.variables)
		if err != nil {
			return false, fmt.Errorf("condition evaluation failed: %w", err)
		}
		
		if !result {
			return false, nil
		}
	}
	
	return true, nil
}