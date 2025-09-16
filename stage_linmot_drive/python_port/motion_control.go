package python_port

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-python/gpython/py"
)

// MotionCommand represents the Go wrapper for motion control functionality
type MotionCommand struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// MotionHeader represents different types of motion commands
type MotionHeader string

const (
	AbsoluteVAI     MotionHeader = "Absolute_VAI"
	RelativeVAI     MotionHeader = "Relative_VAI"
	AbsoluteVAJI    MotionHeader = "Absolute_VAJI"
	RelativeVAJI    MotionHeader = "Relative_VAJI"
	IncrActPosRstI  MotionHeader = "Incr_Act_Pos_RstI"
	AbsoluteSin     MotionHeader = "Absolute_Sin"
	RelativeSin     MotionHeader = "Relative_Sin"
)

// ForceControlHeader represents force control motion commands
type ForceControlHeader string

const (
	VAI_GoToPosWithHigherForceCtrlLimit ForceControlHeader = "VAI Go To Pos With Higher Force Ctrl Limit and Target Force"
	VAI_GoToPosWithLowerForceCtrlLimit  ForceControlHeader = "VAI Go To Pos With Lower Force Ctrl Limit and Target Force"
	VAI_IncActPosWithHigherForceCtrlLimit ForceControlHeader = "VAI Inc Act Pos With Higher Force Ctrl Limit and Target Force"
	VAI_IncActPosWithLowerForceCtrlLimit  ForceControlHeader = "VAI Inc Act Pos With Lower Force Ctrl Limit and Target Force"
	VAI_GoToPosFromActPosAndResetForceControl ForceControlHeader = "VAI Go To Pos From Act Pos And Reset Force Control Set I"
	VAI_IncrementActPosAndResetForceControl   ForceControlHeader = "VAI Increment Act Pos And Reset Force Control Set I"
)

// ForceControlCommand represents force control commands
type ForceControlCommand string

const (
	ChangeTargetForce ForceControlCommand = "Change_Target_Force"
	ResetForceCtrl    ForceControlCommand = "Reset_Force_Ctrl"
)

// NewMotionCommand creates a new motion command instance
func NewMotionCommand() (*MotionCommand, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	// In a real implementation, this would be properly initialized
	app := py.Dict{}
	
	// Create LinMot_MotionCommand instance
	motionCommandClass := pyModule.GetAttrString("LinMot_MotionCommand")
	args := py.Tuple{app}
	
	pyInstance, err := motionCommandClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_MotionCommand instance: %w", err)
	}

	return &MotionCommand{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// SendMotionCommand sends a motion command to the specified drive
func (mc *MotionCommand) SendMotionCommand(drive int, header MotionHeader, targetPos, maxV, acc, dcc float64, jerk int, executeMC bool) (int, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Call Python send_motion_command method
	sendMethod := mc.pyInstance.GetAttrString("send_motion_command")
	args := py.Tuple{
		py.Int(drive),
		py.String(string(header)),
		py.Float(targetPos),
		py.Float(maxV),
		py.Float(acc),
		py.Float(dcc),
		py.Int(jerk),
		py.Bool(executeMC),
	}

	result, err := sendMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to send motion command: %w", err)
	}

	if countNibble, ok := py.AsInt(result); ok {
		return countNibble, nil
	}

	return 0, fmt.Errorf("failed to get count nibble from motion command")
}

// MotionFinished checks if motion has finished for the specified drive(s)
func (mc *MotionCommand) MotionFinished(drive interface{}, countNibble interface{}, doNotWait bool, timeout float64) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Call Python motion_finished method
	finishedMethod := mc.pyInstance.GetAttrString("motion_finished")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList, py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt), py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else {
		return false, fmt.Errorf("invalid drive parameter type")
	}

	result, err := finishedMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check motion finished: %w", err)
	}

	return py.IsTrue(result), nil
}

// InTargetPos checks if the drive(s) have reached the target position
func (mc *MotionCommand) InTargetPos(drive interface{}, countNibble interface{}, doNotWait bool, timeout float64) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Call Python in_target_pos method
	inTargetMethod := mc.pyInstance.GetAttrString("in_target_pos")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList, py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt), py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else {
		return false, fmt.Errorf("invalid drive parameter type")
	}

	result, err := inTargetMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check in target position: %w", err)
	}

	return py.IsTrue(result), nil
}

// InPosRange checks if the actual position is within a specified range
func (mc *MotionCommand) InPosRange(drive int, position, tolerance float64, timeout float64, doNotWait bool) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Call Python in_pos_range method
	inRangeMethod := mc.pyInstance.GetAttrString("in_pos_range")
	args := py.Tuple{
		py.Int(drive),
		py.Float(position),
		py.Float(tolerance),
		py.Float(timeout),
		py.Bool(doNotWait),
	}

	result, err := inRangeMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check position range: %w", err)
	}

	return py.IsTrue(result), nil
}

// CommandReceivedByDrive checks if the drive has acknowledged the command
func (mc *MotionCommand) CommandReceivedByDrive(drive, countNibble int, doNotWait bool) (bool, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Call Python command_received_by_drive method
	receivedMethod := mc.pyInstance.GetAttrString("command_received_by_drive")
	args := py.Tuple{
		py.Int(drive),
		py.Int(countNibble),
		py.Bool(doNotWait),
	}

	result, err := receivedMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check command received: %w", err)
	}

	return py.IsTrue(result), nil
}

// ForceControl represents force control functionality
type ForceControl struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewForceControl creates a new force control instance
func NewForceControl() (*ForceControl, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_ForceControl instance
	forceControlClass := pyModule.GetAttrString("LinMot_ForceControl")
	args := py.Tuple{app}
	
	pyInstance, err := forceControlClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_ForceControl instance: %w", err)
	}

	return &ForceControl{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// MotionForceControl sends a force-controlled motion command
func (fc *ForceControl) MotionForceControl(drive int, header ForceControlHeader, targetPos, maxV, acc, dcc, targetForce, forceLimit float64, executeMC bool) (int, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Call Python motion_force_control method
	controlMethod := fc.pyInstance.GetAttrString("motion_force_control")
	args := py.Tuple{
		py.Int(drive),
		py.String(string(header)),
		py.Float(targetPos),
		py.Float(maxV),
		py.Float(acc),
		py.Float(dcc),
		py.Float(targetForce),
		py.Float(forceLimit),
		py.Bool(executeMC),
	}

	result, err := controlMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to send force control motion command: %w", err)
	}

	if countNibble, ok := py.AsInt(result); ok {
		return countNibble, nil
	}

	return 0, fmt.Errorf("failed to get count nibble from force control command")
}

// ForceControl sends a force control command
func (fc *ForceControl) ForceControl(drive int, header ForceControlCommand, force *float64, executeMC bool) (int, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Call Python force_control method
	controlMethod := fc.pyInstance.GetAttrString("force_control")
	
	var args py.Tuple
	if force != nil {
		args = py.Tuple{
			py.Int(drive),
			py.String(string(header)),
			py.Float(*force),
			py.Bool(executeMC),
		}
	} else {
		args = py.Tuple{
			py.Int(drive),
			py.String(string(header)),
			py.None,
			py.Bool(executeMC),
		}
	}

	result, err := controlMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to send force control command: %w", err)
	}

	if countNibble, ok := py.AsInt(result); ok {
		return countNibble, nil
	}

	return 0, fmt.Errorf("failed to get count nibble from force control command")
}

// GetMeasuredForce retrieves the measured force from the specified drive
func (fc *ForceControl) GetMeasuredForce(drive int) (float64, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	// Call Python get_measured_force method
	forceMethod := fc.pyInstance.GetAttrString("get_measured_force")
	args := py.Tuple{py.Int(drive)}

	result, err := forceMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to get measured force: %w", err)
	}

	if force, ok := py.AsFloat(result); ok {
		return force, nil
	}

	return 0, fmt.Errorf("failed to convert measured force to float")
}

// WaitForceTarget waits until the measured force meets the specified condition
func (fc *ForceControl) WaitForceTarget(drive int, exitCondition string, forceTarget, timeout float64, doNotWait bool) (bool, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	// Call Python wait_force_target method
	waitMethod := fc.pyInstance.GetAttrString("wait_force_target")
	args := py.Tuple{
		py.Int(drive),
		py.String(exitCondition),
		py.Float(forceTarget),
		py.Float(timeout),
		py.Bool(doNotWait),
	}

	result, err := waitMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to wait for force target: %w", err)
	}

	return py.IsTrue(result), nil
}

// ForceRange checks if the measured force is within a specified range
func (fc *ForceControl) ForceRange(drive int, force, tolerance, timeout float64, wait bool) (bool, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	// Call Python force_range method
	rangeMethod := fc.pyInstance.GetAttrString("force_range")
	args := py.Tuple{
		py.Int(drive),
		py.Float(force),
		py.Float(tolerance),
		py.Float(timeout),
		py.Bool(wait),
	}

	result, err := rangeMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check force range: %w", err)
	}

	return py.IsTrue(result), nil
}

// SpecialMotionActive checks if special motion is active
func (fc *ForceControl) SpecialMotionActive(drive interface{}, countNibble interface{}, doNotWait bool, timeout float64) (bool, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	// Call Python special_motion_active method
	activeMethod := fc.pyInstance.GetAttrString("special_motion_active")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList, py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt), py.None, py.Bool(doNotWait), py.Float(timeout)}
	} else {
		return false, fmt.Errorf("invalid drive parameter type")
	}

	result, err := activeMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check special motion active: %w", err)
	}

	return py.IsTrue(result), nil
}

// Close cleans up resources
func (mc *MotionCommand) Close() error {
	if mc.pyInstance != nil {
		mc.pyInstance.DecRef()
	}
	if mc.pyModule != nil {
		mc.pyModule.DecRef()
	}
	return nil
}

// Close cleans up resources
func (fc *ForceControl) Close() error {
	if fc.pyInstance != nil {
		fc.pyInstance.DecRef()
	}
	if fc.pyModule != nil {
		fc.pyModule.DecRef()
	}
	return nil
}