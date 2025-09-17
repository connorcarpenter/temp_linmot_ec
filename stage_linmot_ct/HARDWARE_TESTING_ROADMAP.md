# Hardware Testing Roadmap

## Immediate Implementation Plan

### Phase 1: Hardware Abstraction Layer (Week 1-2)

#### 1.1 Create Hardware Interface

```go
// pkg/hardware/interface.go
package hardware

import (
    "context"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

type HardwareController interface {
    types.DriveController
    Connect(ctx context.Context) error
    Disconnect() error
    GetHardwareInfo() (*HardwareInfo, error)
    IsConnected() bool
    GetConnectionStatus() ConnectionStatus
}

type HardwareInfo struct {
    Model           string
    SerialNumber    string
    FirmwareVersion string
    EtherCATAddress int
    Capabilities    []string
    SafetyLimits    *SafetyLimits
}

type ConnectionStatus struct {
    Connected    bool
    LastSeen     time.Time
    ErrorCount   int
    Latency      time.Duration
}
```

#### 1.2 Implement Real Hardware Controller

```go
// pkg/hardware/real_controller.go
package hardware

import (
    "context"
    "fmt"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive"
)

type RealDriveController struct {
    ethercatMaster *drive.EtherCATMaster
    driveAddress   int
    connected      bool
    hardwareInfo   *HardwareInfo
    lastSeen       time.Time
}

func NewRealDriveController(master *drive.EtherCATMaster, address int) *RealDriveController {
    return &RealDriveController{
        ethercatMaster: master,
        driveAddress:   address,
        connected:      false,
    }
}

func (rdc *RealDriveController) Connect(ctx context.Context) error {
    // Implementation to connect to real hardware
    // This would use the stage_linmot_drive module
    return rdc.ethercatMaster.Connect(ctx)
}

func (rdc *RealDriveController) GetHardwareInfo() (*HardwareInfo, error) {
    // Query hardware for model, serial, firmware, etc.
    return rdc.hardwareInfo, nil
}
```

### Phase 2: Hardware Test Framework (Week 3-4)

#### 2.1 Create Hardware Test Suite

```go
// pkg/hil/hardware_test.go
package hil

import (
    "context"
    "testing"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/hardware"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/safety"
)

type HardwareTestSuite struct {
    controller   hardware.HardwareController
    safetyGuard *safety.SafetyGuard
    testConfig  *TestConfig
    results     *TestResults
}

type TestConfig struct {
    EtherCATMaster    string
    DriveAddress      int
    SafetyLimits      *safety.SafetyLimits
    TestTimeout       time.Duration
    RetryAttempts     int
    LogLevel          string
    EnableSafetyTests bool
}

type TestResults struct {
    PassedTests    int
    FailedTests    int
    SkippedTests   int
    TotalDuration  time.Duration
    TestDetails    []*TestDetail
}

type TestDetail struct {
    TestName    string
    Status      TestStatus
    Duration    time.Duration
    Error       error
    Metrics     map[string]interface{}
}

type TestStatus int

const (
    TestPassed TestStatus = iota
    TestFailed
    TestSkipped
    TestError
)
```

#### 2.2 Implement Core Hardware Tests

```go
// pkg/hil/motion_tests.go
package hil

import (
    "context"
    "testing"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

func (hts *HardwareTestSuite) TestBasicMotion(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), hts.testConfig.TestTimeout)
    defer cancel()

    // Test absolute motion
    err := hts.controller.MoveAbsolute(ctx, 1000.0, 100.0, 50.0, 25.0)
    if err != nil {
        t.Fatalf("MoveAbsolute failed: %v", err)
    }

    // Wait for motion to complete
    for {
        complete, err := hts.controller.IsMotionComplete(ctx)
        if err != nil {
            t.Fatalf("IsMotionComplete failed: %v", err)
        }
        if complete {
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    // Verify position
    position, err := hts.controller.GetPosition(ctx)
    if err != nil {
        t.Fatalf("GetPosition failed: %v", err)
    }

    expectedPosition := 1000.0
    tolerance := 10.0
    if abs(position - expectedPosition) > tolerance {
        t.Errorf("Position mismatch: expected %f, got %f", expectedPosition, position)
    }
}

func (hts *HardwareTestSuite) TestForceControl(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), hts.testConfig.TestTimeout)
    defer cancel()

    // Enable force control
    err := hts.controller.ForceControlOn(ctx)
    if err != nil {
        t.Fatalf("ForceControlOn failed: %v", err)
    }

    // Set force setpoint
    err = hts.controller.SetForce(ctx, 10.0) // 10N
    if err != nil {
        t.Fatalf("SetForce failed: %v", err)
    }

    // Monitor force for a period
    startTime := time.Now()
    for time.Since(startTime) < 5*time.Second {
        force, err := hts.controller.GetForce(ctx)
        if err != nil {
            t.Fatalf("GetForce failed: %v", err)
        }

        // Validate force is within reasonable range
        if force < 5.0 || force > 15.0 {
            t.Errorf("Force out of range: %f N", force)
        }

        time.Sleep(100 * time.Millisecond)
    }

    // Disable force control
    err = hts.controller.ForceControlOff(ctx)
    if err != nil {
        t.Fatalf("ForceControlOff failed: %v", err)
    }
}
```

### Phase 3: End-to-End Test Suite (Week 5-6)

#### 3.1 Create E2E Test Framework

```go
// pkg/e2e/test_suite.go
package e2e

import (
    "context"
    "fmt"
    "time"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct"
    "github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

type E2ETestSuite struct {
    hardwareController  hardware.HardwareController
    commandTableManager *stage_linmot_ct.CommandTableManager
    testScenarios       []*TestScenario
    resultCollector     *ResultCollector
}

type TestScenario struct {
    Name        string
    Description string
    Setup       func() error
    Execute     func() error
    Validate    func() error
    Cleanup     func() error
    Timeout     time.Duration
    Priority    Priority
}

type Priority int

const (
    PriorityCritical Priority = iota
    PriorityHigh
    PriorityMedium
    PriorityLow
)
```

#### 3.2 Implement E2E Test Scenarios

```go
// pkg/e2e/scenarios.go
package e2e

func (e2e *E2ETestSuite) CreateMotionSequenceScenario() *TestScenario {
    return &TestScenario{
        Name:        "Complete Motion Sequence",
        Description: "Test complete motion workflow from home to target position",
        Setup: func() error {
            // Initialize system
            return e2e.hardwareController.Connect(context.Background())
        },
        Execute: func() error {
            // Create command table
            table := e2e.commandTableManager.CreateTable("motion_test", "Motion Test", "Complete motion sequence")
            
            // Add commands
            homeCmd := types.NewCommandBuilder().
                WithID(1).
                WithType(types.CmdHome).
                Build()
            
            moveCmd := types.NewCommandBuilder().
                WithID(2).
                WithType(types.CmdMoveAbsolute).
                WithParameter("position", types.NewPositionValue(1000.0, types.PositionUnitCounts)).
                WithParameter("velocity", types.NewVelocityValue(100.0, types.VelocityUnitCountsS)).
                Build()
            
            waitCmd := types.NewCommandBuilder().
                WithID(3).
                WithType(types.CmdWait).
                WithParameter("time", types.NewTimeValue(1.0, types.TimeUnitS)).
                Build()
            
            e2e.commandTableManager.AddCommand(table, homeCmd)
            e2e.commandTableManager.AddCommand(table, moveCmd)
            e2e.commandTableManager.AddCommand(table, waitCmd)
            
            // Execute
            return e2e.commandTableManager.StartExecution(context.Background(), table)
        },
        Validate: func() error {
            // Validate final position
            position, err := e2e.hardwareController.GetPosition(context.Background())
            if err != nil {
                return err
            }
            
            expectedPosition := 1000.0
            tolerance := 10.0
            if abs(position - expectedPosition) > tolerance {
                return fmt.Errorf("position mismatch: expected %f, got %f", expectedPosition, position)
            }
            
            return nil
        },
        Cleanup: func() error {
            return e2e.hardwareController.Disconnect()
        },
        Timeout: 30 * time.Second,
        Priority: PriorityCritical,
    }
}
```

### Phase 4: Automated Test Execution (Week 7-8)

#### 4.1 Create Test Executor

```go
// pkg/e2e/executor.go
package e2e

import (
    "context"
    "fmt"
    "time"
    "sync"
)

type TestExecutor struct {
    testSuite *E2ETestSuite
    config    *ExecutorConfig
    reporter  *TestReporter
    monitor   *TestMonitor
}

type ExecutorConfig struct {
    ParallelExecution   bool
    MaxConcurrency      int
    RetryFailedTests    bool
    StopOnFirstFailure  bool
    GenerateReport      bool
    ReportFormat        string
    TestTimeout         time.Duration
}

func (te *TestExecutor) RunAllTests(ctx context.Context) (*TestResults, error) {
    var wg sync.WaitGroup
    results := &TestResults{
        TestDetails: make([]*TestDetail, 0),
    }
    
    semaphore := make(chan struct{}, te.config.MaxConcurrency)
    
    for _, scenario := range te.testSuite.testScenarios {
        wg.Add(1)
        go func(scenario *TestScenario) {
            defer wg.Done()
            
            if te.config.ParallelExecution {
                semaphore <- struct{}{}
                defer func() { <-semaphore }()
            }
            
            result := te.runTestScenario(ctx, scenario)
            results.TestDetails = append(results.TestDetails, result)
            
            if result.Status == TestFailed && te.config.StopOnFirstFailure {
                te.monitor.StopAllTests()
            }
        }(scenario)
    }
    
    wg.Wait()
    
    // Generate report
    if te.config.GenerateReport {
        err := te.reporter.GenerateReport(results)
        if err != nil {
            return results, fmt.Errorf("failed to generate report: %v", err)
        }
    }
    
    return results, nil
}
```

## Implementation Checklist

### Week 1-2: Hardware Abstraction
- [ ] Create `HardwareController` interface
- [ ] Implement `RealDriveController` using `stage_linmot_drive`
- [ ] Add hardware discovery and connection logic
- [ ] Implement hardware info querying
- [ ] Add connection status monitoring

### Week 3-4: Hardware Test Framework
- [ ] Create `HardwareTestSuite` structure
- [ ] Implement basic motion tests
- [ ] Implement force control tests
- [ ] Implement I/O tests
- [ ] Implement safety tests
- [ ] Add test result collection and reporting

### Week 5-6: E2E Test Suite
- [ ] Create `E2ETestSuite` framework
- [ ] Implement motion sequence scenario
- [ ] Implement force control workflow
- [ ] Implement I/O control workflow
- [ ] Implement error recovery workflow
- [ ] Implement safety system workflow

### Week 7-8: Automated Execution
- [ ] Create `TestExecutor` with parallel execution
- [ ] Implement test result reporting
- [ ] Add test monitoring and alerting
- [ ] Create CI/CD integration
- [ ] Add performance metrics collection

## Success Criteria

1. **Hardware Connectivity**: Successfully connect to real LinMot C1250-EC hardware
2. **Test Coverage**: >80% of hardware features tested
3. **Test Reliability**: >95% test pass rate on stable hardware
4. **Performance**: Tests complete within reasonable timeframes
5. **Safety**: All safety tests pass without hardware damage
6. **Documentation**: Complete test documentation and examples

## Risk Mitigation

1. **Hardware Safety**: Implement emergency stop and safety limits
2. **Test Isolation**: Ensure tests don't interfere with each other
3. **Error Recovery**: Implement robust error handling and recovery
4. **Data Validation**: Validate all test results and measurements
5. **Backup Testing**: Maintain mock-based tests as fallback