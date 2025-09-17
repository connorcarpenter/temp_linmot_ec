package hardware

import (
	"context"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// HardwareController extends the DriveController interface with hardware-specific methods
type HardwareController interface {
	types.DriveController
	Connect(ctx context.Context) error
	Disconnect() error
	GetHardwareInfo() (*HardwareInfo, error)
	IsConnected() bool
	GetConnectionStatus() *ConnectionStatus
	Ping() error
}

// HardwareInfo contains information about the connected hardware
type HardwareInfo struct {
	Model           string
	SerialNumber    string
	FirmwareVersion string
	EtherCATAddress int
	Capabilities    []string
	SafetyLimits    *SafetyLimits
	LastUpdated     time.Time
}

// SafetyLimits defines hardware-specific safety constraints
type SafetyLimits struct {
	MaxPosition     float64
	MinPosition     float64
	MaxVelocity     float64
	MaxAcceleration float64
	MaxJerk         float64
	MaxForce        float64
	EmergencyStop   bool
}

// ConnectionStatus provides real-time connection information
type ConnectionStatus struct {
	Connected    bool
	LastSeen     time.Time
	ErrorCount   int
	Latency      time.Duration
	Throughput   float64
	Quality      ConnectionQuality
}

// ConnectionQuality represents the quality of the connection
type ConnectionQuality int

const (
	QualityUnknown ConnectionQuality = iota
	QualityExcellent
	QualityGood
	QualityFair
	QualityPoor
	QualityCritical
)

// String returns a string representation of ConnectionQuality
func (cq ConnectionQuality) String() string {
	switch cq {
	case QualityUnknown:
		return "Unknown"
	case QualityExcellent:
		return "Excellent"
	case QualityGood:
		return "Good"
	case QualityFair:
		return "Fair"
	case QualityPoor:
		return "Poor"
	case QualityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// HardwareCapability represents a specific hardware capability
type HardwareCapability int

const (
	CapabilityMotion HardwareCapability = iota
	CapabilityForceControl
	CapabilityDigitalIO
	CapabilityAnalogIO
	CapabilityDataAcquisition
	CapabilitySafetyLimits
	CapabilityEmergencyStop
	CapabilityHome
	CapabilityReset
)

// String returns a string representation of HardwareCapability
func (hc HardwareCapability) String() string {
	switch hc {
	case CapabilityMotion:
		return "Motion"
	case CapabilityForceControl:
		return "ForceControl"
	case CapabilityDigitalIO:
		return "DigitalIO"
	case CapabilityAnalogIO:
		return "AnalogIO"
	case CapabilityDataAcquisition:
		return "DataAcquisition"
	case CapabilitySafetyLimits:
		return "SafetyLimits"
	case CapabilityEmergencyStop:
		return "EmergencyStop"
	case CapabilityHome:
		return "Home"
	case CapabilityReset:
		return "Reset"
	default:
		return "Unknown"
	}
}

// HardwareError represents hardware-specific errors
type HardwareError struct {
	Code        string
	Message     string
	Severity    ErrorSeverity
	Recoverable bool
	Timestamp   time.Time
}

// ErrorSeverity represents the severity of a hardware error
type ErrorSeverity int

const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// String returns a string representation of ErrorSeverity
func (es ErrorSeverity) String() string {
	switch es {
	case SeverityInfo:
		return "Info"
	case SeverityWarning:
		return "Warning"
	case SeverityError:
		return "Error"
	case SeverityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// Error implements the error interface
func (he *HardwareError) Error() string {
	return he.Message
}

// HardwareTestResult represents the result of a hardware test
type HardwareTestResult struct {
	TestName    string
	Passed      bool
	Duration    time.Duration
	Error       error
	Metrics     map[string]interface{}
	Timestamp   time.Time
}

// HardwareTestSuite defines the interface for hardware testing
type HardwareTestSuite interface {
	RunBasicMotionTests(ctx context.Context) ([]*HardwareTestResult, error)
	RunForceControlTests(ctx context.Context) ([]*HardwareTestResult, error)
	RunIOTests(ctx context.Context) ([]*HardwareTestResult, error)
	RunSafetyTests(ctx context.Context) ([]*HardwareTestResult, error)
	RunPerformanceTests(ctx context.Context) ([]*HardwareTestResult, error)
	RunAllTests(ctx context.Context) ([]*HardwareTestResult, error)
}