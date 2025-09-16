LinMot EtherCAT Software
==================================================
These scripts demonstrate how to communicate with LinMot Drives using Python over EtherCAT. They are intended for demonstration and evaluation purposes and are provided "as-is" without official support.

This folder includes two software examples:

-----------------------------------------

1. Python_EtherCAT_Comm
-----------------------------------------
This script provides a command-line interface for establishing and managing EtherCAT communication with LinMot drives. It is designed to demonstrate the core functionality of drive communication, motion control, and data acquisition.

Key Features:
 - Real-time EtherCAT communication using the pysoem library
 - Basic motor control: switch on, home, move
 - Oscilloscope data recording with CSV export
 - Modular architecture with multiprocessing for robust performance
 - Tested with LinMot C1250-MI and F1150-DS drives on Windows and Linux

Includes:
 - Script with all necessary functions
 - 3 different scripts:
    - Basic Motion Control (LinMot_Start_MC_)
    - Basic Motion Control with 2 motors (LinMot_Start_2Motor_)
    - Basic Force Control (LinMot_Start_FC_)
 - Documentation:
   - Quick Start Guide
   - Software Documentation
 - Example of pre-recorded Oscilloscope data (modified)


2. Python_GUI_EtherCAT_Comm
-----------------------------------------
This script features a graphical user interface (GUI) for manual interaction with LinMot drives. It is ideal for quick testing, visualization, and demonstration of motion profiles and drive status.
Note: This is only an application example and is currently in a pilot/testing phase. Use at your own risk.

Key Features:
 - Real-time EtherCAT communication using the pysoem library
 - Tabbed GUI interface for configuration, control, motion, and oscilloscope visualization
 - Real-time monitoring and manual command execution
 - Visual oscilloscope with CSV export
 - Tested with LinMot C1250-MI drives on Windows
 - Pilot-phase software: may contain bugs or incomplete features

Includes:
 - GUI Script with all necessary functions
 - Quick Start Guide
