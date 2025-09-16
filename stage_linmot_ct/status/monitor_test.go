package status

import (
	"context"
	"testing"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func TestNewStatusMonitor(t *testing.T) {
	mockDrive := &MockDriveController{}
	monitor := NewStatusMonitor(mockDrive)
	
	if monitor == nil {
		t.Fatal("Expected non-nil status monitor")
	}
	
	if monitor.driveController != mockDrive {
		t.Error("Expected drive controller to be set correctly")
	}
	
	if monitor.statusCache == nil {
		t.Error("Expected non-nil status cache")
	}
	
	if monitor.updateInterval != 100*time.Millisecond {
		t.Errorf("Expected update interval 100ms, got %v", monitor.updateInterval)
	}
}

func TestStatusMonitor_StartStop(t *testing.T) {
	mockDrive := &MockDriveController{}
	monitor := NewStatusMonitor(mockDrive)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Start monitoring
	err := monitor.Start(ctx)
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}
	
	// Wait a bit for updates
	time.Sleep(150 * time.Millisecond)
	
	// Stop monitoring
	monitor.Stop()
	
	// Check that we got some updates
	status := monitor.GetStatus()
	if status.LastUpdate.IsZero() {
		t.Error("Expected status to be updated")
	}
}

func TestStatusMonitor_GetStatus(t *testing.T) {
	mockDrive := &MockDriveController{
		position:       1000.0,
		velocity:       100.0,
		force:          50.0,
		driveState:     types.DriveStateReady,
		motionComplete: true,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	
	// Manually update status
	monitor.updateStatus(context.Background())
	
	status := monitor.GetStatus()
	
	if status.Position != 1000.0 {
		t.Errorf("Expected position 1000.0, got %f", status.Position)
	}
	
	if status.Velocity != 100.0 {
		t.Errorf("Expected velocity 100.0, got %f", status.Velocity)
	}
	
	if status.Force != 50.0 {
		t.Errorf("Expected force 50.0, got %f", status.Force)
	}
	
	if status.DriveState != types.DriveStateReady {
		t.Errorf("Expected drive state Ready, got %v", status.DriveState)
	}
	
	if !status.MotionComplete {
		t.Error("Expected motion to be complete")
	}
	
	if status.Error != nil {
		t.Errorf("Expected no error, got %v", status.Error)
	}
}

func TestStatusMonitor_GetPosition(t *testing.T) {
	mockDrive := &MockDriveController{
		position: 2000.0,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	position := monitor.GetPosition()
	if position != 2000.0 {
		t.Errorf("Expected position 2000.0, got %f", position)
	}
}

func TestStatusMonitor_GetVelocity(t *testing.T) {
	mockDrive := &MockDriveController{
		velocity: 200.0,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	velocity := monitor.GetVelocity()
	if velocity != 200.0 {
		t.Errorf("Expected velocity 200.0, got %f", velocity)
	}
}

func TestStatusMonitor_GetForce(t *testing.T) {
	mockDrive := &MockDriveController{
		force: 100.0,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	force := monitor.GetForce()
	if force != 100.0 {
		t.Errorf("Expected force 100.0, got %f", force)
	}
}

func TestStatusMonitor_GetDriveState(t *testing.T) {
	mockDrive := &MockDriveController{
		driveState: types.DriveStateMoving,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	driveState := monitor.GetDriveState()
	if driveState != types.DriveStateMoving {
		t.Errorf("Expected drive state Moving, got %v", driveState)
	}
}

func TestStatusMonitor_IsMotionComplete(t *testing.T) {
	mockDrive := &MockDriveController{
		motionComplete: false,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	motionComplete := monitor.IsMotionComplete()
	if motionComplete {
		t.Error("Expected motion to not be complete")
	}
}

func TestStatusMonitor_GetLastUpdate(t *testing.T) {
	mockDrive := &MockDriveController{}
	monitor := NewStatusMonitor(mockDrive)
	
	// Before update
	lastUpdate := monitor.GetLastUpdate()
	if !lastUpdate.IsZero() {
		t.Error("Expected zero time before update")
	}
	
	// After update
	monitor.updateStatus(context.Background())
	lastUpdate = monitor.GetLastUpdate()
	if lastUpdate.IsZero() {
		t.Error("Expected non-zero time after update")
	}
}

func TestStatusMonitor_GetError(t *testing.T) {
	mockDrive := &MockDriveController{
		error: context.DeadlineExceeded,
	}
	
	monitor := NewStatusMonitor(mockDrive)
	monitor.updateStatus(context.Background())
	
	err := monitor.GetError()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStatusMonitor_SetUpdateInterval(t *testing.T) {
	mockDrive := &MockDriveController{}
	monitor := NewStatusMonitor(mockDrive)
	
	newInterval := 50 * time.Millisecond
	monitor.SetUpdateInterval(newInterval)
	
	if monitor.updateInterval != newInterval {
		t.Errorf("Expected update interval %v, got %v", newInterval, monitor.updateInterval)
	}
}

func TestStatusMonitor_IsHealthy(t *testing.T) {
	mockDrive := &MockDriveController{}
	monitor := NewStatusMonitor(mockDrive)
	
	// Before any updates
	if monitor.IsHealthy() {
		t.Error("Expected monitor to be unhealthy before updates")
	}
	
	// After update
	monitor.updateStatus(context.Background())
	if !monitor.IsHealthy() {
		t.Error("Expected monitor to be healthy after update")
	}
	
	// With error
	mockDrive.error = context.DeadlineExceeded
	monitor.updateStatus(context.Background())
	if monitor.IsHealthy() {
		t.Error("Expected monitor to be unhealthy with error")
	}
}

func TestStatusCache_String(t *testing.T) {
	cache := &StatusCache{
		Position:       1000.0,
		Velocity:       100.0,
		Force:          50.0,
		DriveState:     types.DriveStateReady,
		MotionComplete: true,
		LastUpdate:     time.Now(),
		Error:          nil,
	}
	
	str := cache.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}

// MockDriveController for testing
type MockDriveController struct {
	position       float64
	velocity       float64
	force          float64
	driveState     types.DriveState
	motionComplete bool
	error          error
}

func (mdc *MockDriveController) GetPosition(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.position, nil
}

func (mdc *MockDriveController) GetVelocity(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.velocity, nil
}

func (mdc *MockDriveController) GetForce(ctx context.Context) (float64, error) {
	if mdc.error != nil {
		return 0, mdc.error
	}
	return mdc.force, nil
}

func (mdc *MockDriveController) GetDriveState(ctx context.Context) (types.DriveState, error) {
	if mdc.error != nil {
		return types.DriveState(0), mdc.error
	}
	return mdc.driveState, nil
}

func (mdc *MockDriveController) IsMotionComplete(ctx context.Context) (bool, error) {
	if mdc.error != nil {
		return false, mdc.error
	}
	return mdc.motionComplete, nil
}

// Implement other required methods as no-ops for testing
func (mdc *MockDriveController) Stop(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) MoveAbsolute(ctx context.Context, position float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) MoveRelative(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) MoveIncremental(ctx context.Context, distance float64, velocity float64, acceleration float64, jerk float64) error {
	return nil
}

func (mdc *MockDriveController) Jog(ctx context.Context, velocity float64) error {
	return nil
}

func (mdc *MockDriveController) Wait(ctx context.Context, duration time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitPosition(ctx context.Context, position float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitVelocity(ctx context.Context, velocity float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitForce(ctx context.Context, force float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) SetDigitalOutput(ctx context.Context, output int, value bool) error {
	return nil
}

func (mdc *MockDriveController) ClearDigitalOutput(ctx context.Context, output int) error {
	return nil
}

func (mdc *MockDriveController) SetAnalogOutput(ctx context.Context, output int, value float64) error {
	return nil
}

func (mdc *MockDriveController) WaitDigitalInput(ctx context.Context, input int, value bool, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) WaitAnalogInput(ctx context.Context, input int, value float64, tolerance float64, timeout time.Duration) error {
	return nil
}

func (mdc *MockDriveController) Home(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) Reset(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SaveConfiguration(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) LoadConfiguration(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) ForceControlOn(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) ForceControlOff(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SetForce(ctx context.Context, force float64) error {
	return nil
}

func (mdc *MockDriveController) StartOscilloscope(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) StopOscilloscope(ctx context.Context) error {
	return nil
}

func (mdc *MockDriveController) SaveData(ctx context.Context, filename string) error {
	return nil
}

func (mdc *MockDriveController) GetDigitalInput(ctx context.Context, input int) (bool, error) {
	return false, nil
}

func (mdc *MockDriveController) GetAnalogInput(ctx context.Context, input int) (float64, error) {
	return 0.0, nil
}

func (mdc *MockDriveController) GetDigitalOutput(ctx context.Context, output int) (bool, error) {
	return false, nil
}

func (mdc *MockDriveController) GetAnalogOutput(ctx context.Context, output int) (float64, error) {
	return 0.0, nil
}