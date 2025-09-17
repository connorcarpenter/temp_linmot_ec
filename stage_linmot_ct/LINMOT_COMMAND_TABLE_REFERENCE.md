# LinMot-Talk Command Table Reference

This document provides a comprehensive reference for LinMot-Talk Command Table functionality, based on section 4.3 of the LinMot-Talk User Manual. This reference is essential for implementing the `stage_linmot_ct` module.

## Overview

The LinMot Command Table allows users to define up to 255 motion commands that can be executed sequentially or based on external triggers. This feature enables complex motion sequences without continuous communication from the controller, making it ideal for applications requiring precise and repetitive motion control.

## Command Table Structure

### Basic Command Format

Each command in the table consists of:
- **Command Number**: Sequential identifier (1-255)
- **Command Type**: Specific operation to perform
- **Parameters**: Command-specific values and settings
- **Conditions**: Optional execution conditions
- **Next Command**: Pointer to next command in sequence

### Command Categories

#### 1. Motion Commands

**Move Absolute (MA)**
- Moves to a specific absolute position
- Parameters:
  - Position (mm or counts)
  - Velocity (mm/s or counts/s)
  - Acceleration (mm/s² or counts/s²)
  - Deceleration (mm/s² or counts/s²)
  - Jerk (mm/s³ or counts/s³)

**Move Relative (MR)**
- Moves by a relative distance from current position
- Parameters:
  - Distance (mm or counts)
  - Velocity (mm/s or counts/s)
  - Acceleration (mm/s² or counts/s²)
  - Deceleration (mm/s² or counts/s²)
  - Jerk (mm/s³ or counts/s³)

**Move Incremental (MI)**
- Moves by a fixed increment
- Parameters:
  - Increment (mm or counts)
  - Velocity (mm/s or counts/s)
  - Acceleration (mm/s² or counts/s²)
  - Deceleration (mm/s² or counts/s²)

**Jog (JO)**
- Continuous motion in specified direction
- Parameters:
  - Direction (+1 or -1)
  - Velocity (mm/s or counts/s)
  - Acceleration (mm/s² or counts/s²)

**Stop (ST)**
- Stops motion immediately
- Parameters:
  - Deceleration (mm/s² or counts/s²)

#### 2. Control Commands

**Wait (WA)**
- Pauses execution for specified time
- Parameters:
  - Time (ms)

**Wait for Position (WP)**
- Waits until position condition is met
- Parameters:
  - Position (mm or counts)
  - Tolerance (mm or counts)
  - Timeout (ms)

**Wait for Velocity (WV)**
- Waits until velocity condition is met
- Parameters:
  - Velocity (mm/s or counts/s)
  - Tolerance (mm/s or counts/s)
  - Timeout (ms)

**Wait for Force (WF)**
- Waits until force condition is met
- Parameters:
  - Force (N or counts)
  - Tolerance (N or counts)
  - Timeout (ms)

#### 3. I/O Commands

**Set Digital Output (DO)**
- Sets digital output state
- Parameters:
  - Output number (1-8)
  - State (0 or 1)

**Clear Digital Output (CO)**
- Clears digital output state
- Parameters:
  - Output number (1-8)

**Set Analog Output (AO)**
- Sets analog output value
- Parameters:
  - Output number (1-2)
  - Value (V or mA)

**Wait for Digital Input (DI)**
- Waits for digital input condition
- Parameters:
  - Input number (1-8)
  - Expected state (0 or 1)
  - Timeout (ms)

**Wait for Analog Input (AI)**
- Waits for analog input condition
- Parameters:
  - Input number (1-2)
  - Value (V or mA)
  - Tolerance (V or mA)
  - Timeout (ms)

#### 4. Loop Commands

**Loop Start (LS)**
- Marks beginning of loop
- Parameters:
  - Loop counter variable
  - Maximum iterations

**Loop End (LE)**
- Marks end of loop
- Parameters:
  - Loop counter variable
  - Next command if loop continues

**Loop Break (LB)**
- Exits loop conditionally
- Parameters:
  - Condition to check
  - Next command if break occurs

#### 5. Jump Commands

**Jump (JP)**
- Unconditional jump to another command
- Parameters:
  - Target command number

**Jump if True (JT)**
- Conditional jump based on condition
- Parameters:
  - Condition to evaluate
  - Target command number if true
  - Target command number if false

**Jump if False (JF)**
- Conditional jump based on condition
- Parameters:
  - Condition to evaluate
  - Target command number if true
  - Target command number if false

#### 6. System Commands

**Home (HO)**
- Performs homing sequence
- Parameters:
  - Homing method (1-4)
  - Velocity (mm/s or counts/s)
  - Acceleration (mm/s² or counts/s²)
  - Timeout (ms)

**Reset (RE)**
- Resets drive to initial state
- Parameters:
  - Reset type (1-3)

**Save Configuration (SC)**
- Saves current configuration
- Parameters:
  - Configuration slot (1-4)

**Load Configuration (LC)**
- Loads saved configuration
- Parameters:
  - Configuration slot (1-4)

#### 7. Force Control Commands

**Force Control On (FC)**
- Enables force control mode
- Parameters:
  - Force setpoint (N or counts)
  - Force limit (N or counts)
  - Position limit (mm or counts)

**Force Control Off (FO)**
- Disables force control mode
- Parameters:
  - Transition time (ms)

**Set Force (SF)**
- Sets force setpoint
- Parameters:
  - Force (N or counts)
  - Ramp time (ms)

#### 8. Data Acquisition Commands

**Start Oscilloscope (SO)**
- Starts data acquisition
- Parameters:
  - Sample rate (Hz)
  - Number of samples
  - Channels to record

**Stop Oscilloscope (SP)**
- Stops data acquisition
- Parameters:
  - Save data flag (0 or 1)

**Save Data (SD)**
- Saves acquired data
- Parameters:
  - Filename
  - Format (CSV, BIN)

## Command Execution Flow

### Sequential Execution
Commands are executed in numerical order (1, 2, 3, ...) unless modified by jump commands.

### Conditional Execution
Commands can be executed based on:
- Digital input states
- Analog input values
- Position conditions
- Velocity conditions
- Force conditions
- Timer conditions

### Loop Execution
Loops allow repetition of command sequences:
- Fixed number of iterations
- Conditional loops based on input states
- Nested loops (up to 8 levels)

### Error Handling
- Command execution errors
- Timeout conditions
- Hardware faults
- Communication errors

## Parameter Units

### Position
- **Primary**: mm (millimeters)
- **Alternative**: counts (encoder counts)
- **Conversion**: Position[mm] = Position[counts] / ScalingFactor

### Velocity
- **Primary**: mm/s (millimeters per second)
- **Alternative**: counts/s (counts per second)
- **Conversion**: Velocity[mm/s] = Velocity[counts/s] / ScalingFactor

### Acceleration
- **Primary**: mm/s² (millimeters per second squared)
- **Alternative**: counts/s² (counts per second squared)
- **Conversion**: Acceleration[mm/s²] = Acceleration[counts/s²] / ScalingFactor

### Force
- **Primary**: N (Newtons)
- **Alternative**: counts (force counts)
- **Conversion**: Force[N] = Force[counts] / ForceScalingFactor

### Time
- **Primary**: ms (milliseconds)
- **Alternative**: s (seconds)
- **Conversion**: Time[s] = Time[ms] / 1000

## Command Table Management

### Creating Commands
1. Select command type
2. Set parameters
3. Define conditions (optional)
4. Set next command pointer
5. Validate command

### Modifying Commands
1. Load existing command
2. Modify parameters
3. Update conditions
4. Validate changes
5. Save command

### Executing Commands
1. Start command table execution
2. Monitor execution status
3. Handle errors
4. Stop execution when complete

### Debugging Commands
1. Single-step execution
2. Breakpoint setting
3. Variable monitoring
4. Error logging

## Best Practices

### Command Design
- Use descriptive command names
- Group related commands together
- Implement proper error handling
- Use appropriate timeouts

### Performance Optimization
- Minimize command table size
- Use efficient loop structures
- Avoid unnecessary waits
- Optimize parameter ranges

### Safety Considerations
- Implement position limits
- Use force limits
- Add emergency stop conditions
- Validate all parameters

### Maintenance
- Document command sequences
- Use version control
- Regular testing
- Backup configurations

## Integration with EtherCAT

The Command Table integrates with EtherCAT through:
- **PDO Mapping**: Process Data Objects for real-time data exchange
- **SDO Access**: Service Data Objects for configuration
- **Cyclic Communication**: Regular status updates
- **Event Handling**: Asynchronous event processing

## Error Codes

Common error codes and their meanings:
- **0x0000**: No error
- **0x0001**: Invalid command number
- **0x0002**: Invalid parameter
- **0x0003**: Timeout error
- **0x0004**: Hardware fault
- **0x0005**: Communication error
- **0x0006**: Position limit exceeded
- **0x0007**: Force limit exceeded
- **0x0008**: Velocity limit exceeded

## References

- LinMot-Talk User Manual, Section 4.3
- LinMot C1250-EC-XC-0S DC Drive Interface Manual
- EtherCAT Communication Protocol Specification
- LinMot Motion Control Software Documentation