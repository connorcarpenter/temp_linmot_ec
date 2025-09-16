package status

import (
	"fmt"
	"strings"
	"time"
)

// ErrorSeverity represents the severity level of an error
type ErrorSeverity int

const (
	SeverityInfo ErrorSeverity = iota
	SeverityWarning
	SeverityError
	SeverityCritical
)

// String returns a string representation of the error severity
func (es ErrorSeverity) String() string {
	switch es {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// ErrorCode represents a standardized error code
type ErrorCode string

const (
	// Drive errors
	ErrorCodeDriveNotReady     ErrorCode = "DRIVE_NOT_READY"
	ErrorCodeDriveError        ErrorCode = "DRIVE_ERROR"
	ErrorCodeDriveTimeout      ErrorCode = "DRIVE_TIMEOUT"
	ErrorCodeDriveCommunication ErrorCode = "DRIVE_COMM_ERROR"
	
	// Motion errors
	ErrorCodeMotionNotComplete ErrorCode = "MOTION_NOT_COMPLETE"
	ErrorCodeMotionTimeout     ErrorCode = "MOTION_TIMEOUT"
	ErrorCodeMotionAborted     ErrorCode = "MOTION_ABORTED"
	
	// Safety errors
	ErrorCodeSafetyLimitExceeded ErrorCode = "SAFETY_LIMIT_EXCEEDED"
	ErrorCodeEmergencyStop       ErrorCode = "EMERGENCY_STOP"
	ErrorCodePositionOutOfRange  ErrorCode = "POSITION_OUT_OF_RANGE"
	ErrorCodeVelocityExceeded    ErrorCode = "VELOCITY_EXCEEDED"
	ErrorCodeForceExceeded       ErrorCode = "FORCE_EXCEEDED"
	
	// Parameter errors
	ErrorCodeInvalidParameter   ErrorCode = "INVALID_PARAMETER"
	ErrorCodeMissingParameter   ErrorCode = "MISSING_PARAMETER"
	ErrorCodeParameterOutOfRange ErrorCode = "PARAMETER_OUT_OF_RANGE"
	
	// System errors
	ErrorCodeSystemError        ErrorCode = "SYSTEM_ERROR"
	ErrorCodeConfigurationError ErrorCode = "CONFIGURATION_ERROR"
	ErrorCodeInitializationError ErrorCode = "INITIALIZATION_ERROR"
)

// TranslatedError represents a translated and standardized error
type TranslatedError struct {
	Code      ErrorCode
	Severity  ErrorSeverity
	Message   string
	Details   map[string]interface{}
	Timestamp time.Time
	Original  error
}

// Error implements the error interface
func (te *TranslatedError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", te.Code, te.Severity, te.Message)
}

// String returns a detailed string representation
func (te *TranslatedError) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s: %s", te.Code, te.Severity, te.Message))
	
	if te.Details != nil && len(te.Details) > 0 {
		sb.WriteString("\nDetails:")
		for key, value := range te.Details {
			sb.WriteString(fmt.Sprintf("\n  %s: %v", key, value))
		}
	}
	
	if te.Original != nil {
		sb.WriteString(fmt.Sprintf("\nOriginal error: %v", te.Original))
	}
	
	sb.WriteString(fmt.Sprintf("\nTimestamp: %s", te.Timestamp.Format(time.RFC3339)))
	
	return sb.String()
}

// ErrorTranslator provides error translation capabilities
type ErrorTranslator struct {
	translations map[string]*TranslatedError
}

// NewErrorTranslator creates a new error translator
func NewErrorTranslator() *ErrorTranslator {
	et := &ErrorTranslator{
		translations: make(map[string]*TranslatedError),
	}
	
	// Initialize with common error translations
	et.initializeTranslations()
	
	return et
}

// initializeTranslations sets up common error translations
func (et *ErrorTranslator) initializeTranslations() {
	// Drive state errors
	et.translations["drive must be ready or moving for motion commands"] = &TranslatedError{
		Code:      ErrorCodeDriveNotReady,
		Severity:  SeverityError,
		Message:   "Drive is not in a ready state for motion commands",
		Timestamp: time.Now(),
	}
	
	et.translations["drive is in error state"] = &TranslatedError{
		Code:      ErrorCodeDriveError,
		Severity:  SeverityCritical,
		Message:   "Drive is in an error state and requires attention",
		Timestamp: time.Now(),
	}
	
	// Motion errors
	et.translations["previous motion is not complete"] = &TranslatedError{
		Code:      ErrorCodeMotionNotComplete,
		Severity:  SeverityWarning,
		Message:   "Previous motion operation is still in progress",
		Timestamp: time.Now(),
	}
	
	// Safety errors
	et.translations["position"] = &TranslatedError{
		Code:      ErrorCodePositionOutOfRange,
		Severity:  SeverityError,
		Message:   "Position is outside the allowed safety range",
		Timestamp: time.Now(),
	}
	
	et.translations["velocity"] = &TranslatedError{
		Code:      ErrorCodeVelocityExceeded,
		Severity:  SeverityError,
		Message:   "Velocity exceeds the maximum allowed limit",
		Timestamp: time.Now(),
	}
	
	et.translations["force"] = &TranslatedError{
		Code:      ErrorCodeForceExceeded,
		Severity:  SeverityError,
		Message:   "Force exceeds the maximum allowed limit",
		Timestamp: time.Now(),
	}
	
	et.translations["emergency stop is active"] = &TranslatedError{
		Code:      ErrorCodeEmergencyStop,
		Severity:  SeverityCritical,
		Message:   "Emergency stop is active - all motion is halted",
		Timestamp: time.Now(),
	}
	
	// Parameter errors
	et.translations["missing required"] = &TranslatedError{
		Code:      ErrorCodeMissingParameter,
		Severity:  SeverityError,
		Message:   "Required parameter is missing from command",
		Timestamp: time.Now(),
	}
	
	et.translations["invalid parameter"] = &TranslatedError{
		Code:      ErrorCodeInvalidParameter,
		Severity:  SeverityError,
		Message:   "Parameter value is invalid or out of range",
		Timestamp: time.Now(),
	}
}

// TranslateError translates a raw error into a standardized error
func (et *ErrorTranslator) TranslateError(err error, context map[string]interface{}) *TranslatedError {
	if err == nil {
		return nil
	}
	
	errorMsg := err.Error()
	
	// Look for exact matches first
	if translation, exists := et.translations[errorMsg]; exists {
		translated := *translation
		translated.Timestamp = time.Now()
		translated.Original = err
		translated.Details = context
		return &translated
	}
	
	// Look for partial matches
	for pattern, translation := range et.translations {
		if strings.Contains(errorMsg, pattern) {
			translated := *translation
			translated.Timestamp = time.Now()
			translated.Original = err
			translated.Details = context
			return &translated
		}
	}
	
	// Default translation for unknown errors
	return &TranslatedError{
		Code:      ErrorCodeSystemError,
		Severity:  SeverityError,
		Message:   "An unknown error occurred",
		Details:   context,
		Timestamp: time.Now(),
		Original:  err,
	}
}

// AddTranslation adds a custom error translation
func (et *ErrorTranslator) AddTranslation(pattern string, translation *TranslatedError) {
	et.translations[pattern] = translation
}

// GetTranslations returns all available translations
func (et *ErrorTranslator) GetTranslations() map[string]*TranslatedError {
	return et.translations
}

// StatusShaping provides status shaping capabilities
type StatusShaping struct {
	translator *ErrorTranslator
}

// NewStatusShaping creates a new status shaping instance
func NewStatusShaping() *StatusShaping {
	return &StatusShaping{
		translator: NewErrorTranslator(),
	}
}

// ShapeError shapes an error for user consumption
func (ss *StatusShaping) ShapeError(err error, context map[string]interface{}) *TranslatedError {
	return ss.translator.TranslateError(err, context)
}

// ShapeStatus shapes the drive status for user consumption
func (ss *StatusShaping) ShapeStatus(status *StatusCache) map[string]interface{} {
	shaped := map[string]interface{}{
		"position":       status.Position,
		"velocity":       status.Velocity,
		"force":          status.Force,
		"drive_state":    status.DriveState.String(),
		"motion_complete": status.MotionComplete,
		"last_update":    status.LastUpdate.Format(time.RFC3339),
		"healthy":        status.Error == nil,
	}
	
	if status.Error != nil {
		shaped["error"] = ss.ShapeError(status.Error, map[string]interface{}{
			"position":       status.Position,
			"velocity":       status.Velocity,
			"force":          status.Force,
			"drive_state":    status.DriveState.String(),
			"motion_complete": status.MotionComplete,
		})
	}
	
	return shaped
}

// GetErrorSummary returns a summary of errors by severity
func (ss *StatusShaping) GetErrorSummary(errors []*TranslatedError) map[ErrorSeverity]int {
	summary := make(map[ErrorSeverity]int)
	
	for _, err := range errors {
		summary[err.Severity]++
	}
	
	return summary
}

// FilterErrorsBySeverity filters errors by minimum severity level
func (ss *StatusShaping) FilterErrorsBySeverity(errors []*TranslatedError, minSeverity ErrorSeverity) []*TranslatedError {
	var filtered []*TranslatedError
	
	for _, err := range errors {
		if err.Severity >= minSeverity {
			filtered = append(filtered, err)
		}
	}
	
	return filtered
}