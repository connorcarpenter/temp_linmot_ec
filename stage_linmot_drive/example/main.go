package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive"
	"github.com/Smart-Vision-Works/svw_mono/stage_linmot_drive/python_port"
)

func main() {
	// Configure the drive
	config := stage_linmot_drive.DriveConfig{
		AdapterID:           "eth0", // Replace with actual adapter
		NumDevices:          1,
		CycleTime:           0.050, // 50ms cycle time
		NumMonitoring:       4,
		NumParameter:        4,
		ActivateLMDriveData: false,
		MpLogging:           20,
	}

	// Create drive instance
	drive, err := stage_linmot_drive.NewDrive(config)
	if err != nil {
		log.Fatalf("Failed to create drive: %v", err)
	}
	defer drive.Close()

	// Start EtherCAT communication
	if err := drive.Start(); err != nil {
		log.Fatalf("Failed to start drive: %v", err)
	}
	defer drive.Stop()

	fmt.Println("LinMot drive started successfully")

	// Wait for communication to establish
	time.Sleep(2 * time.Second)

	// Get initial status
	status, err := drive.GetStatus()
	if err != nil {
		log.Printf("Failed to get status: %v", err)
	} else {
		for deviceID, deviceStatus := range status {
			fmt.Printf("Device %d Status: %+v\n", deviceID, deviceStatus)
		}
	}

	// Switch on motor
	fmt.Println("Switching on motor...")
	if err := drive.SwitchOnMotor(1); err != nil {
		log.Printf("Failed to switch on motor: %v", err)
	} else {
		fmt.Println("Motor switched on")
	}

	// Wait a bit
	time.Sleep(1 * time.Second)

	// Home motor
	fmt.Println("Homing motor...")
	if err := drive.HomeMotor(1); err != nil {
		log.Printf("Failed to home motor: %v", err)
	} else {
		fmt.Println("Motor homing started")
	}

	// Wait for homing to complete
	time.Sleep(5 * time.Second)

	// Move to position
	fmt.Println("Moving to position 5mm...")
	countNibble, err := drive.MoveToPosition(1, 5.0, 0.01, 0.1, 0.1, 10000)
	if err != nil {
		log.Printf("Failed to move to position: %v", err)
	} else {
		fmt.Printf("Move command sent with count nibble: %d\n", countNibble)
	}

	// Wait for motion to complete
	fmt.Println("Waiting for motion to complete...")
	finished, err := drive.WaitForMotionFinished(1, countNibble, 30*time.Second)
	if err != nil {
		log.Printf("Failed to wait for motion: %v", err)
	} else if finished {
		fmt.Println("Motion completed successfully")
	} else {
		fmt.Println("Motion did not complete within timeout")
	}

	// Get final status
	status, err = drive.GetStatus()
	if err != nil {
		log.Printf("Failed to get final status: %v", err)
	} else {
		for deviceID, deviceStatus := range status {
			fmt.Printf("Final Device %d Status: Position=%.4f, Enabled=%t, Homed=%t\n", 
				deviceID, deviceStatus.ActualPosition, deviceStatus.OperationEnabled, deviceStatus.Homed)
		}
	}

	// Switch off motor
	fmt.Println("Switching off motor...")
	if err := drive.SwitchOffMotor(1); err != nil {
		log.Printf("Failed to switch off motor: %v", err)
	} else {
		fmt.Println("Motor switched off")
	}

	// Get drive information
	driveInfo := drive.GetDriveInfo()
	fmt.Println("Available drive types:")
	for articleNum, info := range driveInfo {
		fmt.Printf("  Article %d: %s (Drive %d)\n", articleNum, info.ModelName, info.DriveNumber)
	}

	fmt.Println("Example completed successfully")
}