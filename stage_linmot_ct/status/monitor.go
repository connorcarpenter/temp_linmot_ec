package status

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// StatusMonitor provides real-time status monitoring capabilities
type StatusMonitor struct {
	driveController types.DriveController
	statusCache     *StatusCache
	updateInterval  time.Duration
	stopChan        chan struct{}
	mu              sync.RWMutex
}

// StatusCache holds cached status information
type StatusCache struct {
	Position       float64
	Velocity       float64
	Force          float64
	DriveState     types.DriveState
	MotionComplete bool
	LastUpdate     time.Time
	Error          error
}

// String returns a string representation of the status cache
func (sc *StatusCache) String() string {
	return fmt.Sprintf("Status{Position: %f, Velocity: %f, Force: %f, DriveState: %s, MotionComplete: %t, LastUpdate: %s, Error: %v}",
		sc.Position, sc.Velocity, sc.Force, sc.DriveState, sc.MotionComplete, sc.LastUpdate.Format(time.RFC3339), sc.Error)
}

// NewStatusMonitor creates a new status monitor
func NewStatusMonitor(driveController types.DriveController) *StatusMonitor {
	return &StatusMonitor{
		driveController: driveController,
		statusCache:     &StatusCache{},
		updateInterval:  100 * time.Millisecond, // Update every 100ms
		stopChan:        make(chan struct{}),
	}
}

// Start begins monitoring the drive status
func (sm *StatusMonitor) Start(ctx context.Context) error {
	go sm.monitorLoop(ctx)
	return nil
}

// Stop stops monitoring the drive status
func (sm *StatusMonitor) Stop() {
	close(sm.stopChan)
}

// monitorLoop runs the monitoring loop
func (sm *StatusMonitor) monitorLoop(ctx context.Context) {
	ticker := time.NewTicker(sm.updateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-sm.stopChan:
			return
		case <-ticker.C:
			sm.updateStatus(ctx)
		}
	}
}

// updateStatus updates the cached status information
func (sm *StatusMonitor) updateStatus(ctx context.Context) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// Update position
	if position, err := sm.driveController.GetPosition(ctx); err == nil {
		sm.statusCache.Position = position
	} else {
		sm.statusCache.Error = fmt.Errorf("failed to get position: %w", err)
	}
	
	// Update velocity
	if velocity, err := sm.driveController.GetVelocity(ctx); err == nil {
		sm.statusCache.Velocity = velocity
	} else {
		sm.statusCache.Error = fmt.Errorf("failed to get velocity: %w", err)
	}
	
	// Update force
	if force, err := sm.driveController.GetForce(ctx); err == nil {
		sm.statusCache.Force = force
	} else {
		sm.statusCache.Error = fmt.Errorf("failed to get force: %w", err)
	}
	
	// Update drive state
	if driveState, err := sm.driveController.GetDriveState(ctx); err == nil {
		sm.statusCache.DriveState = driveState
	} else {
		sm.statusCache.Error = fmt.Errorf("failed to get drive state: %w", err)
	}
	
	// Update motion complete status
	if motionComplete, err := sm.driveController.IsMotionComplete(ctx); err == nil {
		sm.statusCache.MotionComplete = motionComplete
	} else {
		sm.statusCache.Error = fmt.Errorf("failed to get motion complete status: %w", err)
	}
	
	sm.statusCache.LastUpdate = time.Now()
}

// GetStatus returns the current cached status
func (sm *StatusMonitor) GetStatus() *StatusCache {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	return &StatusCache{
		Position:       sm.statusCache.Position,
		Velocity:       sm.statusCache.Velocity,
		Force:          sm.statusCache.Force,
		DriveState:     sm.statusCache.DriveState,
		MotionComplete: sm.statusCache.MotionComplete,
		LastUpdate:     sm.statusCache.LastUpdate,
		Error:          sm.statusCache.Error,
	}
}

// GetPosition returns the current position
func (sm *StatusMonitor) GetPosition() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.Position
}

// GetVelocity returns the current velocity
func (sm *StatusMonitor) GetVelocity() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.Velocity
}

// GetForce returns the current force
func (sm *StatusMonitor) GetForce() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.Force
}

// GetDriveState returns the current drive state
func (sm *StatusMonitor) GetDriveState() types.DriveState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.DriveState
}

// IsMotionComplete returns whether motion is complete
func (sm *StatusMonitor) IsMotionComplete() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.MotionComplete
}

// GetLastUpdate returns the time of the last status update
func (sm *StatusMonitor) GetLastUpdate() time.Time {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.LastUpdate
}

// GetError returns the last error encountered
func (sm *StatusMonitor) GetError() error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.statusCache.Error
}

// SetUpdateInterval sets the update interval
func (sm *StatusMonitor) SetUpdateInterval(interval time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.updateInterval = interval
}

// IsHealthy checks if the status monitor is healthy
func (sm *StatusMonitor) IsHealthy() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	// Check if we have recent updates
	if time.Since(sm.statusCache.LastUpdate) > 2*sm.updateInterval {
		return false
	}
	
	// Check if there are any errors
	return sm.statusCache.Error == nil
}