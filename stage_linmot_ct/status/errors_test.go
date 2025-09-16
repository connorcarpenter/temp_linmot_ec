package status

import (
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewErrorTranslator(t *testing.T) {
	translator := NewErrorTranslator()
	
	if translator == nil {
		t.Fatal("Expected non-nil error translator")
	}
	
	if translator.translations == nil {
		t.Error("Expected non-nil translations map")
	}
	
	// Check that some common translations are initialized
	if _, exists := translator.translations["drive must be ready or moving for motion commands"]; !exists {
		t.Error("Expected drive state translation to be initialized")
	}
	
	if _, exists := translator.translations["previous motion is not complete"]; !exists {
		t.Error("Expected motion state translation to be initialized")
	}
}

func TestErrorTranslator_TranslateError(t *testing.T) {
	translator := NewErrorTranslator()
	
	tests := []struct {
		name    string
		err     error
		context map[string]interface{}
		wantCode ErrorCode
	}{
		{
			name:    "Drive not ready error",
			err:     &MockError{message: "drive must be ready or moving for motion commands"},
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: ErrorCodeDriveNotReady,
		},
		{
			name:    "Motion not complete error",
			err:     &MockError{message: "previous motion is not complete"},
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: ErrorCodeMotionNotComplete,
		},
		{
			name:    "Position out of range error",
			err:     &MockError{message: "position 1000001.0 is above maximum limit 1000000.0"},
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: ErrorCodePositionOutOfRange,
		},
		{
			name:    "Emergency stop error",
			err:     &MockError{message: "emergency stop is active"},
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: ErrorCodeEmergencyStop,
		},
		{
			name:    "Unknown error",
			err:     &MockError{message: "some unknown error"},
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: ErrorCodeSystemError,
		},
		{
			name:    "Nil error",
			err:     nil,
			context: map[string]interface{}{"command_type": "MoveAbsolute"},
			wantCode: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translated := translator.TranslateError(tt.err, tt.context)
			
			if tt.err == nil {
				if translated != nil {
					t.Error("Expected nil translation for nil error")
				}
				return
			}
			
			if translated == nil {
				t.Fatal("Expected non-nil translation")
			}
			
			if translated.Code != tt.wantCode {
				t.Errorf("Expected code %s, got %s", tt.wantCode, translated.Code)
			}
			
			if translated.Original != tt.err {
				t.Error("Expected original error to be preserved")
			}
			
			if translated.Details == nil {
				t.Error("Expected details to be set")
			}
		})
	}
}

func TestErrorTranslator_AddTranslation(t *testing.T) {
	translator := NewErrorTranslator()
	
	customTranslation := &TranslatedError{
		Code:     ErrorCodeSystemError,
		Severity: SeverityWarning,
		Message:  "Custom error message",
	}
	
	translator.AddTranslation("custom error", customTranslation)
	
	// Test that the custom translation is used
	err := &MockError{message: "custom error occurred"}
	translated := translator.TranslateError(err, nil)
	
	if translated == nil {
		t.Fatal("Expected non-nil translation")
	}
	
	if translated.Code != ErrorCodeSystemError {
		t.Errorf("Expected code %s, got %s", ErrorCodeSystemError, translated.Code)
	}
}

func TestErrorTranslator_GetTranslations(t *testing.T) {
	translator := NewErrorTranslator()
	
	translations := translator.GetTranslations()
	
	if translations == nil {
		t.Error("Expected non-nil translations map")
	}
	
	if len(translations) == 0 {
		t.Error("Expected some translations to be available")
	}
}

func TestTranslatedError_Error(t *testing.T) {
	err := &TranslatedError{
		Code:     ErrorCodeDriveNotReady,
		Severity: SeverityError,
		Message:  "Drive is not ready",
	}
	
	expected := "[DRIVE_NOT_READY] ERROR: Drive is not ready"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}

func TestTranslatedError_String(t *testing.T) {
	err := &TranslatedError{
		Code:     ErrorCodeDriveNotReady,
		Severity: SeverityError,
		Message:  "Drive is not ready",
		Details: map[string]interface{}{
			"command_type": "MoveAbsolute",
		},
		Timestamp: time.Now(),
		Original:  &MockError{message: "original error"},
	}
	
	str := err.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	// Check that all components are included
	if !contains(str, "DRIVE_NOT_READY") {
		t.Error("Expected error code in string")
	}
	
	if !contains(str, "ERROR") {
		t.Error("Expected severity in string")
	}
	
	if !contains(str, "Drive is not ready") {
		t.Error("Expected message in string")
	}
}

func TestNewStatusShaping(t *testing.T) {
	shaping := NewStatusShaping()
	
	if shaping == nil {
		t.Fatal("Expected non-nil status shaping")
	}
	
	if shaping.translator == nil {
		t.Error("Expected non-nil translator")
	}
}

func TestStatusShaping_ShapeError(t *testing.T) {
	shaping := NewStatusShaping()
	
	err := &MockError{message: "drive must be ready or moving for motion commands"}
	context := map[string]interface{}{
		"command_type": "MoveAbsolute",
	}
	
	shaped := shaping.ShapeError(err, context)
	
	if shaped == nil {
		t.Fatal("Expected non-nil shaped error")
	}
	
	if shaped.Code != ErrorCodeDriveNotReady {
		t.Errorf("Expected code %s, got %s", ErrorCodeDriveNotReady, shaped.Code)
	}
}

func TestStatusShaping_ShapeStatus(t *testing.T) {
	shaping := NewStatusShaping()
	
	status := &StatusCache{
		Position:       1000.0,
		Velocity:       100.0,
		Force:          50.0,
		DriveState:     types.DriveStateReady,
		MotionComplete: true,
		LastUpdate:     time.Now(),
		Error:          nil,
	}
	
	shaped := shaping.ShapeStatus(status)
	
	if shaped == nil {
		t.Fatal("Expected non-nil shaped status")
	}
	
	if shaped["position"] != 1000.0 {
		t.Errorf("Expected position 1000.0, got %v", shaped["position"])
	}
	
	if shaped["velocity"] != 100.0 {
		t.Errorf("Expected velocity 100.0, got %v", shaped["velocity"])
	}
	
	if shaped["force"] != 50.0 {
		t.Errorf("Expected force 50.0, got %v", shaped["force"])
	}
	
	if shaped["drive_state"] != "ready" {
		t.Errorf("Expected drive_state ready, got %v", shaped["drive_state"])
	}
	
	if shaped["motion_complete"] != true {
		t.Errorf("Expected motion_complete true, got %v", shaped["motion_complete"])
	}
	
	if shaped["healthy"] != true {
		t.Errorf("Expected healthy true, got %v", shaped["healthy"])
	}
}

func TestStatusShaping_GetErrorSummary(t *testing.T) {
	shaping := NewStatusShaping()
	
	errors := []*TranslatedError{
		{Code: ErrorCodeDriveNotReady, Severity: SeverityError},
		{Code: ErrorCodeMotionNotComplete, Severity: SeverityWarning},
		{Code: ErrorCodeSystemError, Severity: SeverityError},
		{Code: ErrorCodeEmergencyStop, Severity: SeverityCritical},
	}
	
	summary := shaping.GetErrorSummary(errors)
	
	if summary[SeverityError] != 2 {
		t.Errorf("Expected 2 errors, got %d", summary[SeverityError])
	}
	
	if summary[SeverityWarning] != 1 {
		t.Errorf("Expected 1 warning, got %d", summary[SeverityWarning])
	}
	
	if summary[SeverityCritical] != 1 {
		t.Errorf("Expected 1 critical, got %d", summary[SeverityCritical])
	}
}

func TestStatusShaping_FilterErrorsBySeverity(t *testing.T) {
	shaping := NewStatusShaping()
	
	errors := []*TranslatedError{
		{Code: ErrorCodeDriveNotReady, Severity: SeverityError},
		{Code: ErrorCodeMotionNotComplete, Severity: SeverityWarning},
		{Code: ErrorCodeSystemError, Severity: SeverityError},
		{Code: ErrorCodeEmergencyStop, Severity: SeverityCritical},
	}
	
	// Filter for errors and above
	filtered := shaping.FilterErrorsBySeverity(errors, SeverityError)
	
	if len(filtered) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(filtered))
	}
	
	// Filter for critical only
	filtered = shaping.FilterErrorsBySeverity(errors, SeverityCritical)
	
	if len(filtered) != 1 {
		t.Errorf("Expected 1 critical error, got %d", len(filtered))
	}
}

func TestErrorSeverity_String(t *testing.T) {
	tests := []struct {
		severity ErrorSeverity
		expected string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityError, "ERROR"},
		{SeverityCritical, "CRITICAL"},
		{ErrorSeverity(999), "UNKNOWN"},
	}
	
	for _, tt := range tests {
		if tt.severity.String() != tt.expected {
			t.Errorf("Severity.String() = %s, want %s", tt.severity.String(), tt.expected)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || contains(s[1:], substr))))
}

// MockError for testing
type MockError struct {
	message string
}

func (me *MockError) Error() string {
	return me.message
}