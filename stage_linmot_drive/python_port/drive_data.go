package python_port

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/go-python/gpython/py"
)

// LMDriveData represents the Go wrapper for the Python LMDrive_Data class
type LMDriveData struct {
	pyInstance *py.Object
	mu         sync.RWMutex
	
	// Configuration
	Config DriveConfig
	
	// Status
	Status DriveStatus
	
	// Inputs (raw data from drive)
	Inputs DriveInputs
	
	// Outputs (data to send to drive)
	Outputs DriveOutputs
}

// DriveConfig contains drive configuration parameters
type DriveConfig struct {
	IsRotaryMotor        bool    `json:"is_rotary_motor"`
	PosScaleNumerator    float64 `json:"pos_scale_numerator"`
	PosScaleDenominator  float64 `json:"pos_scale_denominator"`
	UnitScale            float64 `json:"unit_scale"`
	ModuloFactor         float64 `json:"modulo_factor"`
	FcForceScale         float64 `json:"fc_force_scale"`
	FcTorqueScale        float64 `json:"fc_torque_scale"`
	DriveName            string  `json:"drive_name"`
	DriveType            string  `json:"drive_type"`
}

// DriveStatus contains calculated status values
type DriveStatus struct {
	OperationEnabled   bool    `json:"operation_enabled"`
	SwitchOnLocked     bool    `json:"switch_on_locked"`
	Homed              bool    `json:"homed"`
	MotionActive       bool    `json:"motion_active"`
	Jogging            bool    `json:"jogging"`
	Warning            bool    `json:"warning"`
	Error              bool    `json:"error"`
	ErrorCode          int     `json:"error_code"`
	DemandPosition     float64 `json:"demand_position"`
	ActualPosition     float64 `json:"actual_position"`
	DifferencePosition float64 `json:"difference_position"`
	ActualCurrent      float64 `json:"actual_current"`
	NrOfRevolutions    int     `json:"nr_of_revolutions"`
}

// DriveInputs contains raw input data from the drive
type DriveInputs struct {
	StateVar     uint16 `json:"state_var"`
	StatusWord   uint16 `json:"status_word"`
	WarnWord     uint16 `json:"warn_word"`
	DemandPos    int32  `json:"demand_pos"`
	ActualPos    int32  `json:"actual_pos"`
	DemandCurr   int16  `json:"demand_curr"`
	CfgStatus    uint16 `json:"cfg_status"`
	CfgIndexIn   uint16 `json:"cfg_index_in"`
	CfgValueIn   uint32 `json:"cfg_value_in"`
	MonChannels  []int32 `json:"mon_channels"`
}

// DriveOutputs contains output data to send to the drive
type DriveOutputs struct {
	ControlWord     uint16 `json:"control_word"`
	McHeader        uint16 `json:"mc_header"`
	McParaWords     [10]uint16 `json:"mc_para_words"`
	CfgControl      uint16 `json:"cfg_control"`
	CfgIndexOut     uint16 `json:"cfg_index_out"`
	CfgValueOut     uint32 `json:"cfg_value_out"`
	ParChannels     []uint16 `json:"par_channels"`
}

// NewLMDriveData creates a new LMDrive data instance
func NewLMDriveData(numMonChannels, numParChannels int) (*LMDriveData, error) {
	// Import the LinMot EtherCAT communication module
	pyModule, err := py.Import("LinMot_EtherCAT_Comm_0v82e")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_EtherCAT_Comm_0v82e: %w", err)
	}

	// Create LMDrive_Data instance
	lmDriveDataClass := pyModule.GetAttrString("LMDrive_Data")
	args := py.Tuple{
		py.Int(numMonChannels),
		py.Int(numParChannels),
	}

	pyInstance, err := lmDriveDataClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LMDrive_Data instance: %w", err)
	}

	dd := &LMDriveData{
		pyInstance: pyInstance,
		Config: DriveConfig{
			IsRotaryMotor:       false,
			PosScaleNumerator:   10000.0,
			PosScaleDenominator: 1.0,
			UnitScale:           10000.0,
			ModuloFactor:        360000,
			FcForceScale:        0.1,
			FcTorqueScale:       0.00057295779513082,
			DriveName:           "LMDrive",
			DriveType:           "0",
		},
		Inputs: DriveInputs{
			MonChannels: make([]int32, numMonChannels),
		},
		Outputs: DriveOutputs{
			ControlWord: 0x003E,
			ParChannels: make([]uint16, numParChannels),
		},
	}

	return dd, nil
}

// UnpackInputs unpacks binary input data into the inputs structure
func (dd *LMDriveData) UnpackInputs(data []byte) error {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	// Call Python unpack_inputs method
	unpackMethod := dd.pyInstance.GetAttrString("unpack_inputs")
	_, err := unpackMethod.Call(py.Tuple{py.Bytes(data)})
	if err != nil {
		return fmt.Errorf("failed to unpack inputs: %w", err)
	}

	// Update Go struct from Python object
	return dd.updateFromPython()
}

// PackOutputs packs the outputs structure into binary data
func (dd *LMDriveData) PackOutputs() ([]byte, error) {
	dd.mu.RLock()
	defer dd.mu.RUnlock()

	// Update Python object from Go struct
	if err := dd.updateToPython(); err != nil {
		return nil, err
	}

	// Call Python pack_outputs method
	packMethod := dd.pyInstance.GetAttrString("pack_outputs")
	result, err := packMethod.Call()
	if err != nil {
		return nil, fmt.Errorf("failed to pack outputs: %w", err)
	}

	if dataBytes, ok := py.AsBytes(result); ok {
		return dataBytes, nil
	}

	return nil, fmt.Errorf("failed to convert packed outputs to bytes")
}

// UpdateCalculatedFields updates the calculated status fields
func (dd *LMDriveData) UpdateCalculatedFields() error {
	dd.mu.Lock()
	defer dd.mu.Unlock()

	// Call Python update_calculated_fields method
	updateMethod := dd.pyInstance.GetAttrString("update_calculated_fields")
	_, err := updateMethod.Call()
	if err != nil {
		return fmt.Errorf("failed to update calculated fields: %w", err)
	}

	// Update Go struct from Python object
	return dd.updateFromPython()
}

// updateFromPython updates the Go struct from the Python object
func (dd *LMDriveData) updateFromPython() error {
	// Get config
	configAttr := dd.pyInstance.GetAttrString("config")
	if configDict, ok := configAttr.(*py.Dict); ok {
		dd.updateConfigFromPython(configDict)
	}

	// Get status
	statusAttr := dd.pyInstance.GetAttrString("status")
	if statusDict, ok := statusAttr.(*py.Dict); ok {
		dd.updateStatusFromPython(statusDict)
	}

	// Get inputs
	inputsAttr := dd.pyInstance.GetAttrString("inputs")
	if inputsDict, ok := inputsAttr.(*py.Dict); ok {
		dd.updateInputsFromPython(inputsDict)
	}

	// Get outputs
	outputsAttr := dd.pyInstance.GetAttrString("outputs")
	if outputsDict, ok := outputsAttr.(*py.Dict); ok {
		dd.updateOutputsFromPython(outputsDict)
	}

	return nil
}

// updateToPython updates the Python object from the Go struct
func (dd *LMDriveData) updateToPython() error {
	// Update config
	configAttr := dd.pyInstance.GetAttrString("config")
	if configDict, ok := configAttr.(*py.Dict); ok {
		dd.updateConfigToPython(configDict)
	}

	// Update outputs
	outputsAttr := dd.pyInstance.GetAttrString("outputs")
	if outputsDict, ok := outputsAttr.(*py.Dict); ok {
		dd.updateOutputsToPython(outputsDict)
	}

	return nil
}

func (dd *LMDriveData) updateConfigFromPython(configDict *py.Dict) {
	if val, ok := py.AsBool(configDict.GetItem(py.String("is_rotary_motor"))); ok {
		dd.Config.IsRotaryMotor = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("pos_scale_numerator"))); ok {
		dd.Config.PosScaleNumerator = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("pos_scale_denominator"))); ok {
		dd.Config.PosScaleDenominator = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("unit_scale"))); ok {
		dd.Config.UnitScale = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("modulo_factor"))); ok {
		dd.Config.ModuloFactor = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("fc_force_scale"))); ok {
		dd.Config.FcForceScale = val
	}
	if val, ok := py.AsFloat(configDict.GetItem(py.String("fc_torque_scale"))); ok {
		dd.Config.FcTorqueScale = val
	}
	if val, ok := py.AsString(configDict.GetItem(py.String("drive_name"))); ok {
		dd.Config.DriveName = val
	}
	if val, ok := py.AsString(configDict.GetItem(py.String("drive_type"))); ok {
		dd.Config.DriveType = val
	}
}

func (dd *LMDriveData) updateConfigToPython(configDict *py.Dict) {
	configDict.SetItem(py.String("is_rotary_motor"), py.Bool(dd.Config.IsRotaryMotor))
	configDict.SetItem(py.String("pos_scale_numerator"), py.Float(dd.Config.PosScaleNumerator))
	configDict.SetItem(py.String("pos_scale_denominator"), py.Float(dd.Config.PosScaleDenominator))
	configDict.SetItem(py.String("unit_scale"), py.Float(dd.Config.UnitScale))
	configDict.SetItem(py.String("modulo_factor"), py.Float(dd.Config.ModuloFactor))
	configDict.SetItem(py.String("fc_force_scale"), py.Float(dd.Config.FcForceScale))
	configDict.SetItem(py.String("fc_torque_scale"), py.Float(dd.Config.FcTorqueScale))
	configDict.SetItem(py.String("drive_name"), py.String(dd.Config.DriveName))
	configDict.SetItem(py.String("drive_type"), py.String(dd.Config.DriveType))
}

func (dd *LMDriveData) updateStatusFromPython(statusDict *py.Dict) {
	if val, ok := py.AsBool(statusDict.GetItem(py.String("operation_enabled"))); ok {
		dd.Status.OperationEnabled = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("switch_on_locked"))); ok {
		dd.Status.SwitchOnLocked = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("homed"))); ok {
		dd.Status.Homed = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("motion_active"))); ok {
		dd.Status.MotionActive = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("jogging"))); ok {
		dd.Status.Jogging = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("warning"))); ok {
		dd.Status.Warning = val
	}
	if val, ok := py.AsBool(statusDict.GetItem(py.String("error"))); ok {
		dd.Status.Error = val
	}
	if val, ok := py.AsInt(statusDict.GetItem(py.String("error_code"))); ok {
		dd.Status.ErrorCode = val
	}
	if val, ok := py.AsFloat(statusDict.GetItem(py.String("demand_position"))); ok {
		dd.Status.DemandPosition = val
	}
	if val, ok := py.AsFloat(statusDict.GetItem(py.String("actual_position"))); ok {
		dd.Status.ActualPosition = val
	}
	if val, ok := py.AsFloat(statusDict.GetItem(py.String("difference_position"))); ok {
		dd.Status.DifferencePosition = val
	}
	if val, ok := py.AsFloat(statusDict.GetItem(py.String("actual_current"))); ok {
		dd.Status.ActualCurrent = val
	}
	if val, ok := py.AsInt(statusDict.GetItem(py.String("nr_of_revolutions"))); ok {
		dd.Status.NrOfRevolutions = val
	}
}

func (dd *LMDriveData) updateInputsFromPython(inputsDict *py.Dict) {
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("state_var"))); ok {
		dd.Inputs.StateVar = uint16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("status_word"))); ok {
		dd.Inputs.StatusWord = uint16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("warn_word"))); ok {
		dd.Inputs.WarnWord = uint16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("demand_pos"))); ok {
		dd.Inputs.DemandPos = int32(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("actual_pos"))); ok {
		dd.Inputs.ActualPos = int32(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("demand_curr"))); ok {
		dd.Inputs.DemandCurr = int16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("cfg_status"))); ok {
		dd.Inputs.CfgStatus = uint16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("cfg_index_in"))); ok {
		dd.Inputs.CfgIndexIn = uint16(val)
	}
	if val, ok := py.AsInt(inputsDict.GetItem(py.String("cfg_value_in"))); ok {
		dd.Inputs.CfgValueIn = uint32(val)
	}

	// Update monitoring channels
	for i := 0; i < len(dd.Inputs.MonChannels); i++ {
		key := fmt.Sprintf("mon_ch%d", i+1)
		if val, ok := py.AsInt(inputsDict.GetItem(py.String(key))); ok {
			dd.Inputs.MonChannels[i] = int32(val)
		}
	}
}

func (dd *LMDriveData) updateOutputsFromPython(outputsDict *py.Dict) {
	if val, ok := py.AsInt(outputsDict.GetItem(py.String("control_word"))); ok {
		dd.Outputs.ControlWord = uint16(val)
	}
	if val, ok := py.AsInt(outputsDict.GetItem(py.String("mc_header"))); ok {
		dd.Outputs.McHeader = uint16(val)
	}
	
	// Update parameter words
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("mc_para_word%02d", i)
		if val, ok := py.AsInt(outputsDict.GetItem(py.String(key))); ok {
			dd.Outputs.McParaWords[i] = uint16(val)
		}
	}
	
	if val, ok := py.AsInt(outputsDict.GetItem(py.String("cfg_control"))); ok {
		dd.Outputs.CfgControl = uint16(val)
	}
	if val, ok := py.AsInt(outputsDict.GetItem(py.String("cfg_index_out"))); ok {
		dd.Outputs.CfgIndexOut = uint16(val)
	}
	if val, ok := py.AsInt(outputsDict.GetItem(py.String("cfg_value_out"))); ok {
		dd.Outputs.CfgValueOut = uint32(val)
	}

	// Update parameter channels
	for i := 0; i < len(dd.Outputs.ParChannels); i++ {
		key := fmt.Sprintf("par_ch%d", i+1)
		if val, ok := py.AsInt(outputsDict.GetItem(py.String(key))); ok {
			dd.Outputs.ParChannels[i] = uint16(val)
		}
	}
}

func (dd *LMDriveData) updateOutputsToPython(outputsDict *py.Dict) {
	outputsDict.SetItem(py.String("control_word"), py.Int(dd.Outputs.ControlWord))
	outputsDict.SetItem(py.String("mc_header"), py.Int(dd.Outputs.McHeader))
	
	// Update parameter words
	for i, val := range dd.Outputs.McParaWords {
		key := fmt.Sprintf("mc_para_word%02d", i)
		outputsDict.SetItem(py.String(key), py.Int(val))
	}
	
	outputsDict.SetItem(py.String("cfg_control"), py.Int(dd.Outputs.CfgControl))
	outputsDict.SetItem(py.String("cfg_index_out"), py.Int(dd.Outputs.CfgIndexOut))
	outputsDict.SetItem(py.String("cfg_value_out"), py.Int(dd.Outputs.CfgValueOut))

	// Update parameter channels
	for i, val := range dd.Outputs.ParChannels {
		key := fmt.Sprintf("par_ch%d", i+1)
		outputsDict.SetItem(py.String(key), py.Int(val))
	}
}

// GetStatusString returns a human-readable status string
func (dd *LMDriveData) GetStatusString() string {
	dd.mu.RLock()
	defer dd.mu.RUnlock()

	return fmt.Sprintf(
		"Operation_Enabled: %t | SwitchOn_Locked: %t | Homed: %t | Motion_Active: %t | "+
		"Jogging: %t | Warning: %t | Error: %t | Error_Code: %d | "+
		"Demand_Position: %.4f | Actual_Position: %.4f | Difference_Position: %.4f | "+
		"Actual_Current: %.4f",
		dd.Status.OperationEnabled, dd.Status.SwitchOnLocked, dd.Status.Homed,
		dd.Status.MotionActive, dd.Status.Jogging, dd.Status.Warning, dd.Status.Error,
		dd.Status.ErrorCode, dd.Status.DemandPosition, dd.Status.ActualPosition,
		dd.Status.DifferencePosition, dd.Status.ActualCurrent,
	)
}

// Close cleans up resources
func (dd *LMDriveData) Close() error {
	if dd.pyInstance != nil {
		dd.pyInstance.DecRef()
	}
	return nil
}

// Utility functions for data conversion

// UnsignedToSigned16Bit converts an unsigned 16-bit integer to signed format
func UnsignedToSigned16Bit(value uint16) int16 {
	if value < 0x8000 {
		return int16(value)
	}
	return int16(value - 0x10000)
}

// IEEE754BitsToFloat converts IEEE-754 bit pattern to float
func IEEE754BitsToFloat(value uint32, precision string) (float64, error) {
	switch precision {
	case "single":
		if value > math.MaxUint32 {
			return 0, fmt.Errorf("value exceeds 32-bit range")
		}
		bits := make([]byte, 4)
		binary.BigEndian.PutUint32(bits, value)
		return float64(math.Float32frombits(binary.BigEndian.Uint32(bits))), nil
	case "double":
		if value > math.MaxUint64 {
			return 0, fmt.Errorf("value exceeds 64-bit range")
		}
		bits := make([]byte, 8)
		binary.BigEndian.PutUint64(bits, uint64(value))
		return math.Float64frombits(binary.BigEndian.Uint64(bits)), nil
	default:
		return 0, fmt.Errorf("precision must be 'single' or 'double'")
	}
}