package python_port

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-python/gpython/py"
)

// Configuration represents drive configuration parameter management
type Configuration struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// ConfigHeader represents different types of configuration commands
type ConfigHeader string

const (
	ReadValueROM  ConfigHeader = "Read_Value_ROM"
	ReadValueRAM  ConfigHeader = "Read_Value_RAM"
	WriteValueROM ConfigHeader = "Write_Value_ROM"
	WriteValueRAM ConfigHeader = "Write_Value_RAM"
	WriteValueRAMAndROM ConfigHeader = "Write_Value_RAM_and_ROM"
)

// NewConfiguration creates a new configuration instance
func NewConfiguration() (*Configuration, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_Cfg instance
	configClass := pyModule.GetAttrString("LinMot_Cfg")
	args := py.Tuple{app}
	
	pyInstance, err := configClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_Cfg instance: %w", err)
	}

	return &Configuration{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// ReadConfig reads a configuration value from the drive
func (cfg *Configuration) ReadConfig(drive int, header ConfigHeader, upid string) (int, error) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Call Python read_cfg method
	readMethod := cfg.pyInstance.GetAttrString("read_cfg")
	args := py.Tuple{
		py.Int(drive),
		py.String(string(header)),
		py.String(upid),
	}

	result, err := readMethod.Call(args)
	if err != nil {
		return 0, fmt.Errorf("failed to read config: %w", err)
	}

	if value, ok := py.AsInt(result); ok {
		return value, nil
	}

	return 0, fmt.Errorf("failed to convert config value to int")
}

// WriteConfig writes a configuration value to the drive
func (cfg *Configuration) WriteConfig(drive int, header ConfigHeader, upid string, value int) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Call Python write_cfg method
	writeMethod := cfg.pyInstance.GetAttrString("write_cfg")
	args := py.Tuple{
		py.Int(drive),
		py.String(string(header)),
		py.String(upid),
		py.Int(value),
	}

	_, err := writeMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Oscilloscope handles oscilloscope data acquisition and storage
type Oscilloscope struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewOscilloscope creates a new oscilloscope instance
func NewOscilloscope() (*Oscilloscope, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_Oszilloscope instance
	oscilloscopeClass := pyModule.GetAttrString("LinMot_Oszilloscope")
	args := py.Tuple{app}
	
	pyInstance, err := oscilloscopeClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_Oszilloscope instance: %w", err)
	}

	return &Oscilloscope{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// SaveOscilloscopeSimple saves raw oscilloscope data to a single CSV file
func (osc *Oscilloscope) SaveOscilloscopeSimple(filename string) error {
	osc.mu.Lock()
	defer osc.mu.Unlock()

	// Call Python save_oszi_simple method
	saveMethod := osc.pyInstance.GetAttrString("save_oszi_simple")
	args := py.Tuple{py.String(filename)}

	_, err := saveMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to save oscilloscope data (simple): %w", err)
	}

	return nil
}

// SaveOscilloscope saves oscilloscope data to separate CSV files per device
func (osc *Oscilloscope) SaveOscilloscope(filename string) error {
	osc.mu.Lock()
	defer osc.mu.Unlock()

	// Call Python save_oszi method
	saveMethod := osc.pyInstance.GetAttrString("save_oszi")
	args := py.Tuple{py.String(filename)}

	_, err := saveMethod.Call(args)
	if err != nil {
		return fmt.Errorf("failed to save oscilloscope data: %w", err)
	}

	return nil
}

// UnpackInputData unpacks binary input data into a structured dictionary
func (osc *Oscilloscope) UnpackInputData(data []byte) (map[string]interface{}, error) {
	osc.mu.RLock()
	defer osc.mu.RUnlock()

	// Call Python _unpack_input_data method
	unpackMethod := osc.pyInstance.GetAttrString("_unpack_input_data")
	args := py.Tuple{py.Bytes(data)}

	result, err := unpackMethod.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack input data: %w", err)
	}

	// Convert Python dict to Go map
	if resultDict, ok := result.(*py.Dict); ok {
		return osc.convertPyDictToGoMap(resultDict), nil
	}

	return nil, fmt.Errorf("failed to convert unpacked data to dictionary")
}

// convertPyDictToGoMap converts a Python dictionary to a Go map
func (osc *Oscilloscope) convertPyDictToGoMap(pyDict *py.Dict) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Get all items from the Python dictionary
	items := pyDict.Items()
	for _, item := range items {
		if key, ok := py.AsString(item.Key); ok {
			var value interface{}
			
			// Convert value based on its type
			if intVal, ok := py.AsInt(item.Value); ok {
				value = intVal
			} else if floatVal, ok := py.AsFloat(item.Value); ok {
				value = floatVal
			} else if strVal, ok := py.AsString(item.Value); ok {
				value = strVal
			} else if boolVal, ok := py.AsBool(item.Value); ok {
				value = boolVal
			} else {
				value = item.Value.String()
			}
			
			result[key] = value
		}
	}
	
	return result
}

// Information provides LinMot drive information
type Information struct {
	pyModule   *py.Object
	pyInstance *py.Object
	mu         sync.RWMutex
}

// NewInformation creates a new information instance
func NewInformation() (*Information, error) {
	// Import the data handling module
	pyModule, err := py.Import("LinMot_Data_Handling_0v09")
	if err != nil {
		return nil, fmt.Errorf("failed to import LinMot_Data_Handling_0v09: %w", err)
	}

	// Create a mock app object for the Python class
	app := py.Dict{}
	
	// Create LinMot_Information instance
	infoClass := pyModule.GetAttrString("LinMot_Information")
	args := py.Tuple{app}
	
	pyInstance, err := infoClass.Call(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create LinMot_Information instance: %w", err)
	}

	return &Information{
		pyModule:   pyModule,
		pyInstance: pyInstance,
	}, nil
}

// DriveInfo represents information about a drive
type DriveInfo struct {
	DriveNumber int    `json:"drive_number"`
	ModelName   string `json:"model_name"`
}

// GetDriveInfo returns information about available drives
func (info *Information) GetDriveInfo() map[int]DriveInfo {
	info.mu.RLock()
	defer info.mu.RUnlock()

	// Get the drive_dict from the Python instance
	driveDictAttr := info.pyInstance.GetAttrString("drive_dict")
	
	// Convert Python dict to Go map
	driveInfo := make(map[int]DriveInfo)
	
	if driveDict, ok := driveDictAttr.(*py.Dict); ok {
		items := driveDict.Items()
		for _, item := range items {
			if articleNum, ok := py.AsInt(item.Key); ok {
				if driveList, ok := item.Value.(*py.List); ok && driveList.Len() >= 2 {
					driveNumber := py.AsInt(driveList.GetItem(0))
					modelName := py.AsString(driveList.GetItem(1))
					
					driveInfo[articleNum] = DriveInfo{
						DriveNumber: driveNumber,
						ModelName:   modelName,
					}
				}
			}
		}
	}
	
	return driveInfo
}

// Close cleans up resources
func (cfg *Configuration) Close() error {
	if cfg.pyInstance != nil {
		cfg.pyInstance.DecRef()
	}
	if cfg.pyModule != nil {
		cfg.pyModule.DecRef()
	}
	return nil
}

// Close cleans up resources
func (osc *Oscilloscope) Close() error {
	if osc.pyInstance != nil {
		osc.pyInstance.DecRef()
	}
	if osc.pyModule != nil {
		osc.pyModule.DecRef()
	}
	return nil
}

// Close cleans up resources
func (info *Information) Close() error {
	if info.pyInstance != nil {
		info.pyInstance.DecRef()
	}
	if info.pyModule != nil {
		info.pyModule.DecRef()
	}
	return nil
}