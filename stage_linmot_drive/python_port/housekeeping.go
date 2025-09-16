package python_port

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-python/gpython/py"
)

// Housekeeping represents basic drive operations and selection logic
type Housekeeping struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewHousekeeping creates a new housekeeping instance
func NewHousekeeping() (*Housekeeping, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_Housekeeping instance
	housekeepingClass := pyModule.GetAttrString("LinMot_Housekeeping")
	args := py.Tuple{app}
	
	pyInstance, err := housekeepingClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_Housekeeping instance: %w", err)
	}

	return &Housekeeping{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// SwitchOnMotor switches on the specified motor(s)
func (hk *Housekeeping) SwitchOnMotor(drive interface{}) error {
	hk.mu.Lock()
	defer hk.mu.Unlock()

	// Call Python switch_on_motor method
	switchOnMethod := hk.pyInstance.GetAttrString("switch_on_motor")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt)}
	} else {
		return fmt.Errorf("invalid drive parameter type")
	}

	_, err := switchOnMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to switch on motor: %w", err)
	}

	return nil
}

// SwitchOffMotor switches off the specified motor(s)
func (hk *Housekeeping) SwitchOffMotor(drive interface{}) error {
	hk.mu.Lock()
	defer hk.mu.Unlock()

	// Call Python switch_off_motor method
	switchOffMethod := hk.pyInstance.GetAttrString("switch_off_motor")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt)}
	} else {
		return fmt.Errorf("invalid drive parameter type")
	}

	_, err := switchOffMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to switch off motor: %w", err)
	}

	return nil
}

// HomeMotor sends a homing command to the specified motor(s)
func (hk *Housekeeping) HomeMotor(drive interface{}) error {
	hk.mu.Lock()
	defer hk.mu.Unlock()

	// Call Python home_motor method
	homeMethod := hk.pyInstance.GetAttrString("home_motor")
	
	var args py.Tuple
	if driveList, ok := drive.([]int); ok {
		// Convert Go slice to Python list
		pyList := py.List{}
		for _, d := range driveList {
			pyList.Append(py.Int(d))
		}
		args = py.Tuple{pyList}
	} else if driveInt, ok := drive.(int); ok {
		args = py.Tuple{py.Int(driveInt)}
	} else {
		return fmt.Errorf("invalid drive parameter type")
	}

	_, err := homeMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to home motor: %w", err)
	}

	return nil
}

// DriveCondition provides condition-checking utilities for drive status
type DriveCondition struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewDriveCondition creates a new drive condition instance
func NewDriveCondition() (*DriveCondition, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_DriveCondition instance
	driveConditionClass := pyModule.GetAttrString("LinMot_DriveCondition")
	args := py.Tuple{app}
	
	pyInstance, err := driveConditionClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_DriveCondition instance: %w", err)
	}

	return &DriveCondition{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// IfMaskedStatusWord checks if the masked status word matches a condition
func (dc *DriveCondition) IfMaskedStatusWord(drive int, bitMask uint16, condition uint16, countNibble *int) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	// Call Python if_masked_status_word method
	statusMethod := dc.pyInstance.GetAttrString("if_masked_status_word")
	
	var args py.Tuple
	if countNibble != nil {
		args = py.Tuple{
			py.Int(drive),
			py.Int(bitMask),
			py.Int(condition),
			py.Int(*countNibble),
		}
	} else {
		args = py.Tuple{
			py.Int(drive),
			py.Int(bitMask),
			py.Int(condition),
			py.None,
		}
	}

	result, err := statusMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check masked status word: %w", err)
	}

	return py.IsTrue(result), nil
}

// IfMaskedWarnWord checks if the masked warning word matches a condition
func (dc *DriveCondition) IfMaskedWarnWord(drive int, bitMask uint16, condition uint16, countNibble *int) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	// Call Python if_masked_warn_word method
	warnMethod := dc.pyInstance.GetAttrString("if_masked_warn_word")
	
	var args py.Tuple
	if countNibble != nil {
		args = py.Tuple{
			py.Int(drive),
			py.Int(bitMask),
			py.Int(condition),
			py.Int(*countNibble),
		}
	} else {
		args = py.Tuple{
			py.Int(drive),
			py.Int(bitMask),
			py.Int(condition),
			py.None,
		}
	}

	result, err := warnMethod.Call(args)
	if err != nil {
		return false, fmt.Errorf("failed to check masked warning word: %w", err)
	}

	return py.IsTrue(result), nil
}

// SendData handles low-level data transmission to LinMot drives
type SendData struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewSendData creates a new send data instance
func NewSendData() (*SendData, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_SendData instance
	sendDataClass := pyModule.GetAttrString("LinMot_SendData")
	args := py.Tuple{app}
	
	pyInstance, err := sendDataClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_SendData instance: %w", err)
	}

	return &SendData{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// SwitchONMotor turns the motor ON by manipulating the control word
func (sd *SendData) SwitchONMotor(drive int) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python switchON_motor method
	switchOnMethod := sd.pyInstance.GetAttrString("switchON_motor")
	args := py.Tuple{py.Int(drive)}

	_, err := switchOnMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to switch ON motor: %w", err)
	}

	return nil
}

// SwitchOFFMotor turns the motor OFF by clearing the control word
func (sd *SendData) SwitchOFFMotor(drive int) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python switchOFF_motor method
	switchOffMethod := sd.pyInstance.GetAttrString("switchOFF_motor")
	args := py.Tuple{py.Int(drive)}

	_, err := switchOffMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to switch OFF motor: %w", err)
	}

	return nil
}

// HomeMotor starts the homing procedure
func (sd *SendData) HomeMotor(drive int, executeNow bool) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python home_motor method
	homeMethod := sd.pyInstance.GetAttrString("home_motor")
	args := py.Tuple{py.Int(drive), py.Bool(executeNow)}

	_, err := homeMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to home motor: %w", err)
	}

	return nil
}

// EndHomeMotor ends the homing procedure
func (sd *SendData) EndHomeMotor(drive int) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python end_home_motor method
	endHomeMethod := sd.pyInstance.GetAttrString("end_home_motor")
	args := py.Tuple{py.Int(drive)}

	_, err := endHomeMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to end home motor: %w", err)
	}

	return nil
}

// ErrorAck acknowledges and clears error states
func (sd *SendData) ErrorAck(drive int) error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python error_ack method
	errorAckMethod := sd.pyInstance.GetAttrString("error_ack")
	args := py.Tuple{py.Int(drive)}

	_, err := errorAckMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to acknowledge error: %w", err)
	}

	return nil
}

// GetUnitScale retrieves the unit scaling factor for the selected drive
func (sd *SendData) GetUnitScale(drive int) (float64, error) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	// Call Python get_unit_scale method
	scaleMethod := sd.pyInstance.GetAttrString("get_unit_scale")
	args := py.Tuple{py.Int(drive)}

	result, err := scaleMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to get unit scale: %w", err)
	}

	if scale, ok := py.AsFloat(result); ok {
		return scale, nil
	}

	return 0, fmt.Errorf("failed to convert unit scale to float")
}

// GetForceScale retrieves the force scaling factor for the selected drive
func (sd *SendData) GetForceScale(drive int) (float64, error) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	// Call Python get_force_scale method
	scaleMethod := sd.pyInstance.GetAttrString("get_force_scale")
	args := py.Tuple{py.Int(drive)}

	result, err := scaleMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to get force scale: %w", err)
	}

	if scale, ok := py.AsFloat(result); ok {
		return scale, nil
	}

	return 0, fmt.Errorf("failed to convert force scale to float")
}

// SendDataToSlaves sends output data from all drives to the EtherCAT communication queue
func (sd *SendData) SendDataToSlaves() error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Call Python send_data_to_slaves method
	sendMethod := sd.pyInstance.GetAttrString("send_data_to_slaves")

	_, err := sendMethod.Call()
	if err != nil {
		return fmt.Errorf("failed to send data to slaves: %w", err)
	}

	return nil
}

// Close cleans up resources
func (hk *Housekeeping) Close() error {
	if hk.pyInstance != nil {
		hk.pyInstance.DecRef()
	}
	if hk.pyModule != nil {
		hk.pyModule.DecRef()
	}
	return nil
}

// Close cleans up resources
func (dc *DriveCondition) Close() error {
	if dc.pyInstance != nil {
		dc.pyInstance.DecRef()
	}
	if dc.pyModule != nil {
		dc.pyModule.DecRef()
	}
	return nil
}

// Close cleans up resources
func (sd *SendData) Close() error {
	if sd.pyInstance != nil {
		sd.pyInstance.DecRef()
	}
	if sd.pyModule != nil {
		sd.pyModule.DecRef()
	}
	return nil
}