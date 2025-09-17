# stage_linmot_app — Implementation Plan (TDD)

## Implementation Status

### ✅ Completed (Phase 1)
- **stage_linmot_drive**: Complete Go port of Python LinMot library
  - CPython bridge implementation using `go-python/gpython`
  - EtherCAT communication wrapper
  - Drive data structures and motion control
  - Force control and configuration management
  - Comprehensive error handling and type mapping
  - Unit tests and example usage
  - Complete documentation

### 🚧 In Progress (Phase 2)
- **stage_linmot_ct**: High-level motion layer
  - Command table replacement with clear verbs
  - Unit conversion and scaling
  - Safety guards and preconditions
  - Status shaping and error translation

### 📋 Pending (Phase 3)
- **stage_linmot_app**: Application layer
  - gRPC server implementation
  - Configuration management
  - Telemetry and logging
  - Health checks and monitoring

- **stage_linmot_proto**: Protocol definitions
  - Complete proto definitions
  - Generated Go stubs
  - API versioning

## Phase 2: Command Table Layer (stage_linmot_ct)

### Goals
- Replace LinMot Command Tables with clear, blocking verbs
- Provide unit-explicit API (mm, mm/s, ms)
- Implement safety guards and preconditions
- Handle error translation and status shaping

### Implementation Tasks

#### 2.1 Core API Design
- [ ] Define verb interface (Home, MoveTo, MoveBy, Jog, Stop)
- [ ] Implement status structure with unit-explicit fields
- [ ] Create error taxonomy with machine-readable codes
- [ ] Design configuration snapshot interface

#### 2.2 Motion Verbs
- [ ] `Home(ctx context.Context, drive int) error`
- [ ] `MoveTo(ctx context.Context, drive int, position_mm float64, timeout_ms int64, tolerance_mm float64) error`
- [ ] `MoveBy(ctx context.Context, drive int, delta_mm float64, timeout_ms int64, tolerance_mm float64) error`
- [ ] `Jog(ctx context.Context, drive int, velocity_mm_s, accel_mm_s2 float64, timeout_ms int64) error`
- [ ] `Stop(ctx context.Context, drive int) error`

#### 2.3 Status and Monitoring
- [ ] `GetStatus(ctx context.Context, drive int) (Status, error)`
- [ ] `StreamStatus(ctx context.Context, drive int) (<-chan Status, error)`
- [ ] Implement status field calculations and scaling

#### 2.4 Safety and Validation
- [ ] Precondition checking (enabled, !fault, homed, within limits)
- [ ] Interlock validation
- [ ] Timeout handling and cancellation
- [ ] Error state recovery

#### 2.5 Unit Conversion
- [ ] Position scaling (counts ↔ mm)
- [ ] Velocity scaling (counts/s ↔ mm/s)
- [ ] Force scaling (counts ↔ N)
- [ ] Time conversion (ms ↔ cycles)

## Phase 3: Application Layer (stage_linmot_app)

### Goals
- Implement gRPC server with motion control and diagnostics services
- Handle configuration loading and validation
- Provide telemetry, logging, and health monitoring
- Integrate with command table layer

### Implementation Tasks

#### 3.1 gRPC Server
- [ ] Implement MotionControl service
- [ ] Implement Diagnostics service
- [ ] Add request validation and error handling
- [ ] Implement streaming status endpoint

#### 3.2 Configuration Management
- [ ] Load and validate `defaults.yaml`
- [ ] Create immutable configuration snapshots
- [ ] Handle configuration updates and validation
- [ ] Environment-specific configuration support

#### 3.3 Telemetry and Monitoring
- [ ] Structured logging with context
- [ ] Metrics collection (latency, jitter, fault counts)
- [ ] Health check endpoints
- [ ] Audit logging for motion commands

#### 3.4 Integration
- [ ] Wire command table layer
- [ ] Implement graceful shutdown
- [ ] Add signal handling
- [ ] Create main application entry point

## Phase 4: Protocol Definitions (stage_linmot_proto)

### Goals
- Complete protobuf definitions for all services
- Generate Go stubs and client code
- Implement API versioning strategy
- Add comprehensive documentation

### Implementation Tasks

#### 4.1 Service Definitions
- [ ] Complete MotionControl service
- [ ] Complete Diagnostics service
- [ ] Add Status service for monitoring
- [ ] Define error codes and messages

#### 4.2 Message Definitions
- [ ] Motion control messages with units
- [ ] Status and telemetry messages
- [ ] Configuration messages
- [ ] Error and diagnostic messages

#### 4.3 Code Generation
- [ ] Set up protoc/buf build system
- [ ] Generate Go stubs
- [ ] Create client libraries
- [ ] Add validation and documentation

## Testing Strategy

### Unit Tests
- [ ] Command table layer unit tests
- [ ] Application layer unit tests
- [ ] Protocol validation tests
- [ ] Error handling tests

### Integration Tests
- [ ] End-to-end motion control tests
- [ ] gRPC client-server tests
- [ ] Configuration loading tests
- [ ] Error recovery tests

### Hardware Tests
- [ ] Manual testing with LinMot hardware
- [ ] Performance and timing tests
- [ ] Stress testing and fault injection
- [ ] Real-world scenario validation

## Deployment Strategy

### Development
- [ ] Local development environment setup
- [ ] Docker containerization
- [ ] CI/CD pipeline configuration
- [ ] Automated testing

### Production
- [ ] BalenaOS deployment configuration
- [ ] Hardware-specific configuration
- [ ] Monitoring and alerting setup
- [ ] Backup and recovery procedures

## Success Criteria

### Phase 2 (Command Table Layer)
- ✅ All motion verbs implemented and tested
- ✅ Unit conversion working correctly
- ✅ Safety guards preventing unsafe operations
- ✅ Clear error messages and status reporting

### Phase 3 (Application Layer)
- ✅ gRPC server responding to all requests
- ✅ Configuration loading and validation working
- ✅ Telemetry and monitoring operational
- ✅ Graceful shutdown and error recovery

### Phase 4 (Protocol Layer)
- ✅ Complete protobuf definitions
- ✅ Generated code working correctly
- ✅ API versioning implemented
- ✅ Client libraries available

### Overall System
- ✅ Motion control working reliably
- ✅ Real-time performance meeting requirements
- ✅ Error handling and recovery robust
- ✅ Documentation complete and accurate