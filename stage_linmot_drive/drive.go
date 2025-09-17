package stage_linmot_drive

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive/python_port"
)

// Drive represents the main LinMot drive interface
type Drive struct {
	// Core communication
	ethercatComm *python_port.EtherCATCommunication
	
	// Data handling
	driveData map[int]*python_port.LMDriveData
	
	// Control classes
	motionCommand *python_port.MotionCommand
	forceControl  *python_port.ForceControl
	housekeeping  *python_port.Housekeeping
	driveCondition *python_port.DriveCondition
	sendData     *python_port.SendData
	configuration *python_port.Configuration
	oscilloscope *python_port.Oscilloscope
	information  *python_port.Information
	
	// Configuration
	numDevices    int
	numMonitoring int
	numParameter  int
	cycleTime     float64
	
	// Synchronization
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
}

// DriveConfig represents the configuration for a LinMot drive
type DriveConfig struct {
	AdapterID           string  `json:"adapter_id"`
	NumDevices          int     `json:"num_devices"`
	CycleTime           float64 `json:"cycle_time"`
	NumMonitoring       int     `json:"num_monitoring"`
	NumParameter        int     `json:"num_parameter"`
	ActivateLMDriveData bool    `json:"activate_lm_drive_data"`
	MpLogging           int     `json:"mp_logging"`
}

// NewDrive creates a new LinMot drive instance
func NewDrive(config DriveConfig) (*Drive, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create EtherCAT communication
	ethercatComm, err := python_port.NewEtherCATCommunication(
		config.AdapterID,
		config.NumDevices,
		config.CycleTime,
		config.NumMonitoring,
		config.NumParameter,
		config.ActivateLMDriveData,
		config.MpLogging,
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create EtherCAT communication: %w", err)
	}
	
	// Create control classes
	motionCommand, err := python_port.NewMotionCommand()
	if err != nil {
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create motion command: %w", err)
	}
	
	forceControl, err := python_port.NewForceControl()
	if err != nil {
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create force control: %w", err)
	}
	
	housekeeping, err := python_port.NewHousekeeping()
	if err != nil {
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create housekeeping: %w", err)
	}
	
	driveCondition, err := python_port.NewDriveCondition()
	if err != nil {
		housekeeping.Close()
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create drive condition: %w", err)
	}
	
	sendData, err := python_port.NewSendData()
	if err != nil {
		driveCondition.Close()
		housekeeping.Close()
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create send data: %w", err)
	}
	
	configuration, err := python_port.NewConfiguration()
	if err != nil {
		sendData.Close()
		driveCondition.Close()
		housekeeping.Close()
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create configuration: %w", err)
	}
	
	oscilloscope, err := python_port.NewOscilloscope()
	if err != nil {
		configuration.Close()
		sendData.Close()
		driveCondition.Close()
		housekeeping.Close()
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create oscilloscope: %w", err)
	}
	
	information, err := python_port.NewInformation()
	if err != nil {
		oscilloscope.Close()
		configuration.Close()
		sendData.Close()
		driveCondition.Close()
		housekeeping.Close()
		forceControl.Close()
		motionCommand.Close()
		ethercatComm.Close()
		cancel()
		return nil, fmt.Errorf("failed to create information: %w", err)
	}
	
	// Create drive data instances
	driveData := make(map[int]*python_port.LMDriveData)
	for i := 1; i <= config.NumDevices; i++ {
		data, err := python_port.NewLMDriveData(config.NumMonitoring, config.NumParameter)
		if err != nil {
			// Clean up already created instances
			for _, dd := range driveData {
				dd.Close()
			}
			information.Close()
			oscilloscope.Close()
			configuration.Close()
			sendData.Close()
			driveCondition.Close()
			housekeeping.Close()
			forceControl.Close()
			motionCommand.Close()
			ethercatComm.Close()
			cancel()
			return nil, fmt.Errorf("failed to create drive data for device %d: %w", i, err)
		}
		driveData[i] = data
	}
	
	return &Drive{
		ethercatComm:   ethercatComm,
		driveData:      driveData,
		motionCommand:  motionCommand,
		forceControl:   forceControl,
		housekeeping:   housekeeping,
		driveCondition: driveCondition,
		sendData:       sendData,
		configuration:  configuration,
		oscilloscope:   oscilloscope,
		information:    information,
		numDevices:     config.NumDevices,
		numMonitoring:  config.NumMonitoring,
		numParameter:   config.NumParameter,
		cycleTime:      config.CycleTime,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

// Start begins the EtherCAT communication
func (d *Drive) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if d.isActive {
		return fmt.Errorf("drive is already active")
	}
	
	if err := d.ethercatComm.Start(); err != nil {
		return fmt.Errorf("failed to start EtherCAT communication: %w", err)
	}
	
	d.isActive = true
	return nil
}

// Stop stops the EtherCAT communication
func (d *Drive) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if !d.isActive {
		return nil
	}
	
	if err := d.ethercatComm.Stop(); err != nil {
		return fmt.Errorf("failed to stop EtherCAT communication: %w", err)
	}
	
	d.isActive = false
	d.cancel()
	return nil
}

// IsActive returns whether the drive is currently active
func (d *Drive) IsActive() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.isActive
}

// GetStatus returns the current status of all drives
func (d *Drive) GetStatus() (map[int]python_port.DriveStatus, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return nil, fmt.Errorf("drive is not active")
	}
	
	// Get input data from EtherCAT
	inputData, err := d.ethercatComm.GetInputData()
	if err != nil {
		return nil, fmt.Errorf("failed to get input data: %w", err)
	}
	
	// Process data for each device
	status := make(map[int]python_port.DriveStatus)
	dataLength := 18 + 8 + (4 * d.numMonitoring)
	
	for i := 1; i <= d.numDevices; i++ {
		startIdx := (i - 1) * dataLength
		endIdx := startIdx + dataLength
		
		if endIdx > len(inputData) {
			return nil, fmt.Errorf("insufficient input data for device %d", i)
		}
		
		deviceData := inputData[startIdx:endIdx]
		
		// Unpack and update drive data
		if err := d.driveData[i].UnpackInputs(deviceData); err != nil {
			return nil, fmt.Errorf("failed to unpack inputs for device %d: %w", i, err)
		}
		
		if err := d.driveData[i].UpdateCalculatedFields(); err != nil {
			return nil, fmt.Errorf("failed to update calculated fields for device %d: %w", i, err)
		}
		
		status[i] = d.driveData[i].Status
	}
	
	return status, nil
}

// SwitchOnMotor switches on the specified motor(s)
func (d *Drive) SwitchOnMotor(drive interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.housekeeping.SwitchOnMotor(drive)
}

// SwitchOffMotor switches off the specified motor(s)
func (d *Drive) SwitchOffMotor(drive interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.housekeeping.SwitchOffMotor(drive)
}

// HomeMotor sends a homing command to the specified motor(s)
func (d *Drive) HomeMotor(drive interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.housekeeping.HomeMotor(drive)
}

// MoveToPosition moves the motor to a specific position
func (d *Drive) MoveToPosition(drive int, position, maxVelocity, acceleration, deceleration float64, jerk int) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return 0, fmt.Errorf("drive is not active")
	}
	
	return d.motionCommand.SendMotionCommand(
		drive,
		python_port.AbsoluteVAI,
		position,
		maxVelocity,
		acceleration,
		deceleration,
		jerk,
		true,
	)
}

// MoveByOffset moves the motor by a relative offset
func (d *Drive) MoveByOffset(drive int, offset, maxVelocity, acceleration, deceleration float64, jerk int) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return 0, fmt.Errorf("drive is not active")
	}
	
	return d.motionCommand.SendMotionCommand(
		drive,
		python_port.RelativeVAI,
		offset,
		maxVelocity,
		acceleration,
		deceleration,
		jerk,
		true,
	)
}

// WaitForMotionFinished waits until motion is finished
func (d *Drive) WaitForMotionFinished(drive interface{}, countNibble interface{}, timeout time.Duration) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return false, fmt.Errorf("drive is not active")
	}
	
	return d.motionCommand.MotionFinished(drive, countNibble, false, timeout.Seconds())
}

// WaitForTargetPosition waits until the drive reaches the target position
func (d *Drive) WaitForTargetPosition(drive interface{}, countNibble interface{}, timeout time.Duration) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return false, fmt.Errorf("drive is not active")
	}
	
	return d.motionCommand.InTargetPos(drive, countNibble, false, timeout.Seconds())
}

// ReadConfig reads a configuration value from the drive
func (d *Drive) ReadConfig(drive int, header python_port.ConfigHeader, upid string) (int, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return 0, fmt.Errorf("drive is not active")
	}
	
	return d.configuration.ReadConfig(drive, header, upid)
}

// WriteConfig writes a configuration value to the drive
func (d *Drive) WriteConfig(drive int, header python_port.ConfigHeader, upid string, value int) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.configuration.WriteConfig(drive, header, upid, value)
}

// GetErrorMessages retrieves error messages from the communication
func (d *Drive) GetErrorMessages() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return nil, fmt.Errorf("drive is not active")
	}
	
	return d.ethercatComm.GetErrorMessages()
}

// GetInfoMessages retrieves info messages from the communication
func (d *Drive) GetInfoMessages() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return nil, fmt.Errorf("drive is not active")
	}
	
	return d.ethercatComm.GetInfoMessages()
}

// SetOscilloscopeRecording enables or disables oscilloscope data recording
func (d *Drive) SetOscilloscopeRecording(enable bool) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.ethercatComm.SetOscilloscopeRecording(enable)
}

// SaveOscilloscopeData saves oscilloscope data to files
func (d *Drive) SaveOscilloscopeData(filename string) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if !d.isActive {
		return fmt.Errorf("drive is not active")
	}
	
	return d.oscilloscope.SaveOscilloscope(filename)
}

// GetDriveInfo returns information about available drives
func (d *Drive) GetDriveInfo() map[int]python_port.DriveInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	return d.information.GetDriveInfo()
}

// Close cleans up all resources
func (d *Drive) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Stop if active
	if d.isActive {
		d.ethercatComm.Stop()
		d.isActive = false
	}
	
	// Close all components
	var lastErr error
	
	// Close drive data
	for _, dd := range d.driveData {
		if err := dd.Close(); err != nil {
			lastErr = err
		}
	}
	
	// Close control classes
	if err := d.motionCommand.Close(); err != nil {
		lastErr = err
	}
	if err := d.forceControl.Close(); err != nil {
		lastErr = err
	}
	if err := d.housekeeping.Close(); err != nil {
		lastErr = err
	}
	if err := d.driveCondition.Close(); err != nil {
		lastErr = err
	}
	if err := d.sendData.Close(); err != nil {
		lastErr = err
	}
	if err := d.configuration.Close(); err != nil {
		lastErr = err
	}
	if err := d.oscilloscope.Close(); err != nil {
		lastErr = err
	}
	if err := d.information.Close(); err != nil {
		lastErr = err
	}
	
	// Close EtherCAT communication
	if err := d.ethercatComm.Close(); err != nil {
		lastErr = err
	}
	
	// Cancel context
	d.cancel()
	
	return lastErr
}