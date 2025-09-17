# Hardware Testing Guide

## Overview

The `stage_linmot_ct` module includes comprehensive hardware testing capabilities for LinMot C1250-EC servo drives over EtherCAT. This guide explains how to run the hardware tests.

## Quick Start

### 1. Check Hardware Availability

```bash
# Navigate to the stage_linmot_ct directory
cd stage_linmot_ct

# Check if hardware is available
go run examples/hardware_testing/main.go
```

### 2. Run Specific Test Categories

The hardware testing framework supports several test categories:

- **motion**: Basic motion tests (absolute, relative, incremental, jog, stop)
- **force_control**: Force control tests (enable/disable, setpoint, monitoring)
- **io**: I/O tests (digital and analog inputs/outputs)
- **safety**: Safety tests (emergency stop, limits, error recovery)
- **performance**: Performance tests (latency, throughput, jitter)
- **integration**: Complete system integration tests

### 3. Using the CLI Tool

The CLI tool provides a command-line interface for running hardware tests:

```bash
# Show help
go run cmd/hardware-test/main.go -help

# Check hardware availability
go run cmd/hardware-test/main.go -check-hardware

# Run all tests
go run cmd/hardware-test/main.go -run-hil -run-e2e

# Run only motion tests
go run cmd/hardware-test/main.go -category=motion -run-hil

# Run tests with custom hardware address
go run cmd/hardware-test/main.go -hardware-address=2 -ethercat-master=eth1

# Generate HTML report
go run cmd/hardware-test/main.go -generate-report -report-format=HTML
```

## Test Categories

### Motion Tests
- **Absolute Motion**: Move to specific position
- **Relative Motion**: Move relative to current position
- **Incremental Motion**: Move by specific increment
- **Jog Motion**: Continuous motion at specified velocity
- **Stop Motion**: Stop current motion

### Force Control Tests
- **Force Control Enable/Disable**: Test force control activation
- **Force Setpoint**: Test force setpoint setting
- **Force Monitoring**: Test force measurement and monitoring

### I/O Tests
- **Digital Outputs**: Test digital output control
- **Digital Inputs**: Test digital input reading
- **Analog Outputs**: Test analog output control
- **Analog Inputs**: Test analog input reading

### Safety Tests
- **Emergency Stop**: Test emergency stop functionality
- **Safety Limits**: Test position, velocity, and force limits
- **Error Recovery**: Test error recovery mechanisms

### Performance Tests
- **Latency**: Measure command execution latency
- **Throughput**: Measure command execution throughput
- **Jitter**: Measure timing jitter

### Integration Tests
- **Complete Workflow**: Test complete motion sequences
- **Multi-Feature**: Test combinations of features
- **Error Handling**: Test error scenarios and recovery

## Configuration Options

### Hardware Configuration
- `-hardware-address`: EtherCAT slave address (1-255)
- `-ethercat-master`: EtherCAT master interface (e.g., eth0)
- `-timeout`: Test timeout duration

### Test Configuration
- `-category`: Test category to run
- `-parallel`: Run tests in parallel
- `-max-concurrency`: Maximum concurrent tests
- `-stop-on-failure`: Stop on first test failure

### Reporting Configuration
- `-generate-report`: Generate test report
- `-report-format`: Report format (HTML, JSON, PDF, XML)
- `-report-output`: Report output file

## Example Usage

### Basic Hardware Check
```bash
go run cmd/hardware-test/main.go -check-hardware
```

Output:
```
Stage LinMot CT Hardware Testing
================================
Test Category: all
Hardware Address: 1
EtherCAT Master: eth0
Test Timeout: 5m0s
Log Level: INFO
Parallel Execution: true
Max Concurrency: 3
Stop on Failure: false

Checking hardware availability...
Hardware Info:
  Model: LinMot C1250-EC
  Serial Number: MOCK123456
  Firmware Version: 1.0.0
  EtherCAT Address: 1
  Capabilities: [Motion ForceControl DigitalIO AnalogIO]
Hardware is available and responding!
```

### Run Motion Tests
```bash
go run cmd/hardware-test/main.go -category=motion -run-hil
```

Output:
```
Stage LinMot CT Hardware Testing
================================
Test Category: motion
Hardware Address: 1
EtherCAT Master: eth0
Test Timeout: 5m0s
Log Level: INFO
Parallel Execution: true
Max Concurrency: 3
Stop on Failure: false

Running Hardware-in-the-Loop Tests:
==================================
HIL Test Results: 5 tests executed
  PASS - Absolute Motion (100ms)
  PASS - Relative Motion (100ms)
  PASS - Incremental Motion (100ms)
  PASS - Jog Motion (1s)
  PASS - Stop Motion (100ms)
  Passed: 5, Failed: 0
```

### Run All Tests with Report
```bash
go run cmd/hardware-test/main.go -run-hil -run-e2e -generate-report
```

## Troubleshooting

### Common Issues

1. **Module Import Errors**
   - Ensure you're running from the correct directory
   - Check that all dependencies are installed

2. **Hardware Connection Issues**
   - Verify EtherCAT master is running
   - Check hardware address is correct
   - Ensure hardware is powered and connected

3. **Test Failures**
   - Check hardware state and configuration
   - Verify safety limits are appropriate
   - Review test logs for specific error messages

### Debug Mode

Run tests with debug logging:
```bash
go run cmd/hardware-test/main.go -log-level=DEBUG -run-hil
```

## Integration with CI/CD

The hardware testing framework can be integrated with CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Hardware Tests
  run: |
    go run stage_linmot_ct/cmd/hardware-test/main.go -run-hil -run-e2e -generate-report
```

## Next Steps

1. **Connect Real Hardware**: Replace the mock controller with actual hardware interface
2. **Customize Tests**: Add custom test scenarios for specific use cases
3. **Integrate CI/CD**: Set up automated testing in your CI/CD pipeline
4. **Monitor Performance**: Use the built-in performance monitoring features

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review the test logs for error details
3. Consult the API documentation for specific interfaces
4. Check the best practices guide for usage recommendations