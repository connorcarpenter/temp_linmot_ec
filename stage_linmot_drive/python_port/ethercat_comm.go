package python_port

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-python/gpython"
	"github.com/go-python/gpython/py"
)

// EtherCATCommunication represents the Go wrapper for the Python EtherCAT communication class
type EtherCATCommunication struct {
	adapterID           string
	numDevices          int
	cycleTime           float64
	numMonitoring       int
	numParameter        int
	activateLMDriveData bool
	mpLogging           int
	cpuAffinity         []int
	realtime            bool
	realtimePriority    int

	// Python objects
	pyModule    *py.Object
	pyInstance  *py.Object
	pyLock      *py.Object
	stopEvent   *py.Object
	errorQueue  *py.Object
	infoQueue   *py.Object
	updateQueue *py.Object
	dataQueue   *py.Object

	// Go synchronization
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
}

// NewEtherCATCommunication creates a new EtherCAT communication instance
func NewEtherCATCommunication(adapterID string, numDevices int, cycleTime float64, 
	numMonitoring, numParameter int, activateLMDriveData bool, mpLogging int) (*EtherCATCommunication, error) {
	
	ctx, cancel := context.WithCancel(context.Background())
	
	ec := &EtherCATCommunication{
		adapterID:           adapterID,
		numDevices:          numDevices,
		cycleTime:           cycleTime,
		numMonitoring:       numMonitoring,
		numParameter:        numParameter,
		activateLMDriveData: activateLMDriveData,
		mpLogging:           mpLogging,
		ctx:                 ctx,
		cancel:              cancel,
	}

	// Initialize Python interpreter
	if err := ec.initializePython(); err != nil {
		return nil, fmt.Errorf("failed to initialize Python: %w", err)
	}

	return ec, nil
}

// initializePython sets up the Python interpreter and imports the LinMot module
func (ec *EtherCATCommunication) initializePython() error {
	// Import the LinMot EtherCAT communication module
	pyModule, err := py.Import("LinMot_EtherCAT_Comm_0v82e")
	if err != nil {
		return fmt.Errorf("failed to import LinMot_EtherCAT_Comm_0v82e: %w", err)
	}
	ec.pyModule = pyModule

	// Create multiprocessing lock
	mpModule, err := py.Import("multiprocessing")
	if err != nil {
		return fmt.Errorf("failed to import multiprocessing: %w", err)
	}
	
	lockClass := mpModule.GetAttrString("Lock")
	ec.pyLock, err = lockClass.Call()
	if err != nil {
		return fmt.Errorf("failed to create multiprocessing lock: %w", err)
	}

	// Create EtherCATCommunication instance
	commClass := ec.pyModule.GetAttrString("EtherCATCommunication")
	
	args := py.Tuple{
		py.String(ec.adapterID),
		py.Int(ec.numDevices),
		py.Float(ec.cycleTime),
		ec.pyLock,
		py.Int(ec.numMonitoring),
		py.Int(ec.numParameter),
		py.Bool(ec.activateLMDriveData),
		py.Int(ec.mpLogging),
	}

	ec.pyInstance, err = commClass.Call(args)
	if err != nil {
		return fmt.Errorf("failed to create EtherCATCommunication instance: %w", err)
	}

	// Get queue references
	ec.stopEvent = ec.pyInstance.GetAttrString("stop_event")
	ec.errorQueue = ec.pyInstance.GetAttrString("error_queue")
	ec.infoQueue = ec.pyInstance.GetAttrString("info_queue")
	ec.updateQueue = ec.pyInstance.GetAttrString("update_queue")
	ec.dataQueue = ec.pyInstance.GetAttrString("data_queue")

	return nil
}

// Start begins the EtherCAT communication process
func (ec *EtherCATCommunication) Start() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.isActive {
		return fmt.Errorf("EtherCAT communication already active")
	}

	// Call the Python start method
	startMethod := ec.pyInstance.GetAttrString("start")
	_, err := startMethod.Call()
	if err != nil {
		return fmt.Errorf("failed to start EtherCAT communication: %w", err)
	}

	ec.isActive = true
	return nil
}

// Stop stops the EtherCAT communication process
func (ec *EtherCATCommunication) Stop() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if !ec.isActive {
		return nil
	}

	// Call the Python stop method
	stopMethod := ec.pyInstance.GetAttrString("stop")
	_, err := stopMethod.Call()
	if err != nil {
		return fmt.Errorf("failed to stop EtherCAT communication: %w", err)
	}

	ec.isActive = false
	ec.cancel()
	return nil
}

// IsActive returns whether the communication is currently active
func (ec *EtherCATCommunication) IsActive() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.isActive
}

// GetErrorMessages retrieves error messages from the error queue
func (ec *EtherCATCommunication) GetErrorMessages() ([]string, error) {
	var messages []string
	
	// Check if queue is empty
	emptyMethod := ec.errorQueue.GetAttrString("empty")
	isEmpty, err := emptyMethod.Call()
	if err != nil {
		return nil, fmt.Errorf("failed to check error queue: %w", err)
	}

	if py.IsTrue(isEmpty) {
		return messages, nil
	}

	// Drain the queue
	for {
		getMethod := ec.errorQueue.GetAttrString("get_nowait")
		msg, err := getMethod.Call()
		if err != nil {
			// Queue is empty or error occurred
			break
		}
		
		if msgStr, ok := py.AsString(msg); ok {
			messages = append(messages, msgStr)
		}
	}

	return messages, nil
}

// GetInfoMessages retrieves info messages from the info queue
func (ec *EtherCATCommunication) GetInfoMessages() ([]string, error) {
	var messages []string
	
	// Check if queue is empty
	emptyMethod := ec.infoQueue.GetAttrString("empty")
	isEmpty, err := emptyMethod.Call()
	if err != nil {
		return nil, fmt.Errorf("failed to check info queue: %w", err)
	}

	if py.IsTrue(isEmpty) {
		return messages, nil
	}

	// Drain the queue
	for {
		getMethod := ec.infoQueue.GetAttrString("get_nowait")
		msg, err := getMethod.Call()
		if err != nil {
			// Queue is empty or error occurred
			break
		}
		
		if msgStr, ok := py.AsString(msg); ok {
			messages = append(messages, msgStr)
		}
	}

	return messages, nil
}

// SendUpdateData sends output data to the drives
func (ec *EtherCATCommunication) SendUpdateData(data [][]byte) error {
	// Convert Go data to Python list of bytes
	pyDataList := py.List{}
	for _, deviceData := range data {
		pyDataList.Append(py.Bytes(deviceData))
	}

	// Put data in update queue
	putMethod := ec.updateQueue.GetAttrString("put")
	_, err := putMethod.Call(py.Tuple{pyDataList})
	if err != nil {
		return fmt.Errorf("failed to send update data: %w", err)
	}

	return nil
}

// GetInputData retrieves input data from all devices
func (ec *EtherCATCommunication) GetInputData() ([]byte, error) {
	// Get the data array from the Python instance
	dataAttr := ec.pyInstance.GetAttrString("data")
	
	// Convert Python array to Go bytes
	if dataBytes, ok := py.AsBytes(dataAttr); ok {
		return dataBytes, nil
	}
	
	return nil, fmt.Errorf("failed to get input data")
}

// SetOscilloscopeRecording enables or disables oscilloscope data recording
func (ec *EtherCATCommunication) SetOscilloscopeRecording(enable bool) error {
	dataQueueON := ec.pyInstance.GetAttrString("data_queue_ON")
	
	if enable {
		setMethod := dataQueueON.GetAttrString("set")
		_, err := setMethod.Call()
		if err != nil {
			return fmt.Errorf("failed to enable oscilloscope recording: %w", err)
		}
	} else {
		clearMethod := dataQueueON.GetAttrString("clear")
		_, err := clearMethod.Call()
		if err != nil {
			return fmt.Errorf("failed to disable oscilloscope recording: %w", err)
		}
	}
	
	return nil
}

// GetOscilloscopeData retrieves oscilloscope data from the queue
func (ec *EtherCATCommunication) GetOscilloscopeData() ([]OscilloscopeSample, error) {
	var samples []OscilloscopeSample
	
	// Check if queue is empty
	emptyMethod := ec.dataQueue.GetAttrString("empty")
	isEmpty, err := emptyMethod.Call()
	if err != nil {
		return nil, fmt.Errorf("failed to check data queue: %w", err)
	}

	if py.IsTrue(isEmpty) {
		return samples, nil
	}

	// Drain the queue
	for {
		getMethod := ec.dataQueue.GetAttrString("get_nowait")
		sample, err := getMethod.Call()
		if err != nil {
			// Queue is empty or error occurred
			break
		}
		
		// Convert Python tuple (sample_nr, data) to Go struct
		if sampleTuple, ok := sample.(*py.Tuple); ok && sampleTuple.Len() == 2 {
			sampleNr := py.AsInt(sampleTuple.GetItem(0))
			data := py.AsBytes(sampleTuple.GetItem(1))
			
			samples = append(samples, OscilloscopeSample{
				SampleNumber: sampleNr,
				Data:        data,
			})
		}
	}

	return samples, nil
}

// Close cleans up resources
func (ec *EtherCATCommunication) Close() error {
	if ec.isActive {
		if err := ec.Stop(); err != nil {
			return err
		}
	}
	
	// Clean up Python objects
	if ec.pyInstance != nil {
		ec.pyInstance.DecRef()
	}
	if ec.pyModule != nil {
		ec.pyModule.DecRef()
	}
	if ec.pyLock != nil {
		ec.pyLock.DecRef()
	}
	
	return nil
}

// OscilloscopeSample represents a single oscilloscope data sample
type OscilloscopeSample struct {
	SampleNumber int
	Data        []byte
}