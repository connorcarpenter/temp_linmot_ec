package e2e

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/hardware"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_ct/types"
)

// E2ETestSuite provides comprehensive end-to-end testing capabilities
type E2ETestSuite struct {
	hardwareController  hardware.HardwareController
	commandTableManager *stage_linmot_ct.CommandTableManager
	testScenarios       []*TestScenario
	resultCollector     *ResultCollector
	mu                  sync.RWMutex
}

// TestScenario defines a complete end-to-end test scenario
type TestScenario struct {
	Name        string
	Description string
	Setup       func() error
	Execute     func() error
	Validate    func() error
	Cleanup     func() error
	Timeout     time.Duration
	Priority    Priority
	Category    Category
}

// Priority represents the priority of a test scenario
type Priority int

const (
	PriorityCritical Priority = iota
	PriorityHigh
	PriorityMedium
	PriorityLow
)

// String returns a string representation of Priority
func (p Priority) String() string {
	switch p {
	case PriorityCritical:
		return "Critical"
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

// Category represents the category of a test scenario
type Category int

const (
	CategoryMotion Category = iota
	CategoryForceControl
	CategoryIO
	CategorySafety
	CategoryPerformance
	CategoryIntegration
)

// String returns a string representation of Category
func (c Category) String() string {
	switch c {
	case CategoryMotion:
		return "Motion"
	case CategoryForceControl:
		return "ForceControl"
	case CategoryIO:
		return "IO"
	case CategorySafety:
		return "Safety"
	case CategoryPerformance:
		return "Performance"
	case CategoryIntegration:
		return "Integration"
	default:
		return "Unknown"
	}
}

// ResultCollector collects and manages test results
type ResultCollector struct {
	results []*TestResult
	mu      sync.RWMutex
}

// TestResult contains the result of a test scenario execution
type TestResult struct {
	Scenario    *TestScenario
	Status      TestStatus
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Error       error
	Metrics     map[string]interface{}
	Logs        []string
}

// TestStatus represents the status of a test execution
type TestStatus int

const (
	StatusNotStarted TestStatus = iota
	StatusRunning
	StatusPassed
	StatusFailed
	StatusSkipped
	StatusError
)

// String returns a string representation of TestStatus
func (ts TestStatus) String() string {
	switch ts {
	case StatusNotStarted:
		return "NotStarted"
	case StatusRunning:
		return "Running"
	case StatusPassed:
		return "Passed"
	case StatusFailed:
		return "Failed"
	case StatusSkipped:
		return "Skipped"
	case StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// NewE2ETestSuite creates a new end-to-end test suite
func NewE2ETestSuite(controller hardware.HardwareController, manager *stage_linmot_ct.CommandTableManager) *E2ETestSuite {
	return &E2ETestSuite{
		hardwareController:  controller,
		commandTableManager: manager,
		testScenarios:       make([]*TestScenario, 0),
		resultCollector:     &ResultCollector{results: make([]*TestResult, 0)},
	}
}

// AddScenario adds a test scenario to the suite
func (e2e *E2ETestSuite) AddScenario(scenario *TestScenario) {
	e2e.mu.Lock()
	defer e2e.mu.Unlock()
	e2e.testScenarios = append(e2e.testScenarios, scenario)
}

// RunScenario executes a specific test scenario
func (e2e *E2ETestSuite) RunScenario(ctx context.Context, scenario *TestScenario) (*TestResult, error) {
	result := &TestResult{
		Scenario: scenario,
		Status:   StatusNotStarted,
		Metrics:  make(map[string]interface{}),
		Logs:     make([]string, 0),
	}

	// Set up timeout
	scenarioCtx, cancel := context.WithTimeout(ctx, scenario.Timeout)
	defer cancel()

	// Execute scenario
	result.Status = StatusRunning
	result.StartTime = time.Now()

	// Setup
	if scenario.Setup != nil {
		if err := scenario.Setup(); err != nil {
			result.Status = StatusError
			result.Error = fmt.Errorf("setup failed: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result, result.Error
		}
	}

	// Execute
	if scenario.Execute != nil {
		if err := scenario.Execute(); err != nil {
			result.Status = StatusFailed
			result.Error = fmt.Errorf("execution failed: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result, result.Error
		}
	}

	// Validate
	if scenario.Validate != nil {
		if err := scenario.Validate(); err != nil {
			result.Status = StatusFailed
			result.Error = fmt.Errorf("validation failed: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result, result.Error
		}
	}

	// Cleanup
	if scenario.Cleanup != nil {
		if err := scenario.Cleanup(); err != nil {
			result.Status = StatusError
			result.Error = fmt.Errorf("cleanup failed: %v", err)
			result.EndTime = time.Now()
			result.Duration = result.EndTime.Sub(result.StartTime)
			return result, result.Error
		}
	}

	result.Status = StatusPassed
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// RunAllScenarios executes all test scenarios
func (e2e *E2ETestSuite) RunAllScenarios(ctx context.Context) ([]*TestResult, error) {
	e2e.mu.RLock()
	scenarios := make([]*TestScenario, len(e2e.testScenarios))
	copy(scenarios, e2e.testScenarios)
	e2e.mu.RUnlock()

	var results []*TestResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, scenario := range scenarios {
		wg.Add(1)
		go func(scenario *TestScenario) {
			defer wg.Done()

			result, err := e2e.RunScenario(ctx, scenario)
			if err != nil {
				result.Error = err
			}

			mu.Lock()
			results = append(results, result)
			e2e.resultCollector.AddResult(result)
			mu.Unlock()
		}(scenario)
	}

	wg.Wait()
	return results, nil
}

// RunScenariosByCategory executes test scenarios filtered by category
func (e2e *E2ETestSuite) RunScenariosByCategory(ctx context.Context, category Category) ([]*TestResult, error) {
	e2e.mu.RLock()
	scenarios := make([]*TestScenario, 0)
	for _, scenario := range e2e.testScenarios {
		if scenario.Category == category {
			scenarios = append(scenarios, scenario)
		}
	}
	e2e.mu.RUnlock()

	var results []*TestResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, scenario := range scenarios {
		wg.Add(1)
		go func(scenario *TestScenario) {
			defer wg.Done()

			result, err := e2e.RunScenario(ctx, scenario)
			if err != nil {
				result.Error = err
			}

			mu.Lock()
			results = append(results, result)
			e2e.resultCollector.AddResult(result)
			mu.Unlock()
		}(scenario)
	}

	wg.Wait()
	return results, nil
}

// RunScenariosByPriority executes test scenarios filtered by priority
func (e2e *E2ETestSuite) RunScenariosByPriority(ctx context.Context, priority Priority) ([]*TestResult, error) {
	e2e.mu.RLock()
	scenarios := make([]*TestScenario, 0)
	for _, scenario := range e2e.testScenarios {
		if scenario.Priority == priority {
			scenarios = append(scenarios, scenario)
		}
	}
	e2e.mu.RUnlock()

	var results []*TestResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, scenario := range scenarios {
		wg.Add(1)
		go func(scenario *TestScenario) {
			defer wg.Done()

			result, err := e2e.RunScenario(ctx, scenario)
			if err != nil {
				result.Error = err
			}

			mu.Lock()
			results = append(results, result)
			e2e.resultCollector.AddResult(result)
			mu.Unlock()
		}(scenario)
	}

	wg.Wait()
	return results, nil
}

// GetResults returns all test results
func (e2e *E2ETestSuite) GetResults() []*TestResult {
	return e2e.resultCollector.GetResults()
}

// GetResultsByStatus returns test results filtered by status
func (e2e *E2ETestSuite) GetResultsByStatus(status TestStatus) []*TestResult {
	return e2e.resultCollector.GetResultsByStatus(status)
}

// GetResultsByCategory returns test results filtered by category
func (e2e *E2ETestSuite) GetResultsByCategory(category Category) []*TestResult {
	return e2e.resultCollector.GetResultsByCategory(category)
}

// GetResultsByPriority returns test results filtered by priority
func (e2e *E2ETestSuite) GetResultsByPriority(priority Priority) []*TestResult {
	return e2e.resultCollector.GetResultsByPriority(priority)
}

// GetSummary returns a summary of all test results
func (e2e *E2ETestSuite) GetSummary() *TestSummary {
	return e2e.resultCollector.GetSummary()
}

// AddResult adds a test result to the collector
func (rc *ResultCollector) AddResult(result *TestResult) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.results = append(rc.results, result)
}

// GetResults returns all test results
func (rc *ResultCollector) GetResults() []*TestResult {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	results := make([]*TestResult, len(rc.results))
	copy(results, rc.results)
	return results
}

// GetResultsByStatus returns test results filtered by status
func (rc *ResultCollector) GetResultsByStatus(status TestStatus) []*TestResult {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	
	var results []*TestResult
	for _, result := range rc.results {
		if result.Status == status {
			results = append(results, result)
		}
	}
	return results
}

// GetResultsByCategory returns test results filtered by category
func (rc *ResultCollector) GetResultsByCategory(category Category) []*TestResult {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	
	var results []*TestResult
	for _, result := range rc.results {
		if result.Scenario.Category == category {
			results = append(results, result)
		}
	}
	return results
}

// GetResultsByPriority returns test results filtered by priority
func (rc *ResultCollector) GetResultsByPriority(priority Priority) []*TestResult {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	
	var results []*TestResult
	for _, result := range rc.results {
		if result.Scenario.Priority == priority {
			results = append(results, result)
		}
	}
	return results
}

// GetSummary returns a summary of all test results
func (rc *ResultCollector) GetSummary() *TestSummary {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	
	summary := &TestSummary{
		TotalTests:    len(rc.results),
		PassedTests:   0,
		FailedTests:   0,
		SkippedTests:  0,
		ErrorTests:    0,
		TotalDuration: 0,
		Categories:    make(map[Category]int),
		Priorities:    make(map[Priority]int),
		Statuses:      make(map[TestStatus]int),
	}

	for _, result := range rc.results {
		summary.TotalDuration += result.Duration
		
		switch result.Status {
		case StatusPassed:
			summary.PassedTests++
		case StatusFailed:
			summary.FailedTests++
		case StatusSkipped:
			summary.SkippedTests++
		case StatusError:
			summary.ErrorTests++
		}
		
		summary.Categories[result.Scenario.Category]++
		summary.Priorities[result.Scenario.Priority]++
		summary.Statuses[result.Status]++
	}

	return summary
}

// TestSummary provides a summary of test execution results
type TestSummary struct {
	TotalTests    int
	PassedTests   int
	FailedTests   int
	SkippedTests  int
	ErrorTests    int
	TotalDuration time.Duration
	Categories    map[Category]int
	Priorities    map[Priority]int
	Statuses      map[TestStatus]int
}

// GetPassRate returns the pass rate as a percentage
func (ts *TestSummary) GetPassRate() float64 {
	if ts.TotalTests == 0 {
		return 0.0
	}
	return float64(ts.PassedTests) / float64(ts.TotalTests) * 100.0
}

// GetAverageDuration returns the average test duration
func (ts *TestSummary) GetAverageDuration() time.Duration {
	if ts.TotalTests == 0 {
		return 0
	}
	return ts.TotalDuration / time.Duration(ts.TotalTests)
}

// String returns a string representation of the test summary
func (ts *TestSummary) String() string {
	return fmt.Sprintf("Test Summary: %d total, %d passed (%.1f%%), %d failed, %d skipped, %d errors, avg duration: %v",
		ts.TotalTests, ts.PassedTests, ts.GetPassRate(), ts.FailedTests, ts.SkippedTests, ts.ErrorTests, ts.GetAverageDuration())
}