"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           Start_2Motor_ .py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  16.09.2025
    Version:        0.85
    Description:    It establishes EtherCAT communication with the LinMot drive 
                    and offers basic motion command testing capabilities for 
                    two drives.

    Disclaimer:
    ------------
    This is a demo project developed for NTI AG | LinMot & MagSpring. This
    software is provided "as-is" without any warranties or guarantees.
    NTI AG does not provide support, updates, or maintenance for this software. 
    Use it at your own risk.

    Dependencies:
    -------------
    - See Documentation Chapter 2: Quick Start Guide

    Description:    
    ------------
    This script establishes EtherCAT communication with LinMot Drives and provides 
    basic testing and motion functionalities via the `EtherCATCommunication` class.
    
    It configures the EtherCAT master, manages the communication cycle, and 
    handles process data exchange with LinMot slave devices exclusively.
    
    See Documentation Chapter 3 for a more detailled description.

    Key Features:
    -------------
    - Configures and starts EtherCAT communication.
    - Displays and logs process data from the drives.
    - Allows basic drive control (enable, home, move) with two motors.
    - Supports oscilloscope recording and monitoring.
    - Ensures proper shutdown to prevent resource conflicts.
    
    Usage:
    ------
    See Chapters 2 and 3 of the documentation for more detailed descriptions 
    and important information.

    1. Ensure dependencies are installed.
        Make sure all the required Python packages are installed.

    2. Adjust configuration variables.
        This script has primarily been tested on Windows.
        Other environments may not be compatible.
        At a minimum, adjust the following variables in the main_test class:
        - adapter_id: Use Find_EC_Master.py to identify the correct adapter ID.
        - noDev: Specify the number of LinMot slaves connected to the master.
        - cycle_time: On Windows, the typical minimum cycle time is >10 ms.
                      On Linux, the typical minimum cycle time is >1 ms.

    3. Run the script.
        Open the Windows Command Prompt.
        Navigate to the directory containing filename.py.
        Run the script using python filename.py.

    4. Observe and interact.
        The script will start EtherCAT communication and display the received data.
        It also offers basic control commands to interact with the drives.

        When the received data is outputted, confirm that the displayed values 
        are meaningful.

        The script will prompt you before executing each action.
            Before confirming with Enter, ensure that the action will not 
            cause issues (e.g., motor crash).
            If you notice any abnormal behavior, such as unexpected motor 
            motion, immediately press Ctrl+C to stop the script.

        The script will stop automatically once execution is complete.
        You can also stop the script manually at any time using Ctrl+C.

    License:
    --------
    This software is developed by "NTI AG | LinMot & MagSpring".
    
    Copyright (c) 2025 NTI AG | LinMot & MagSpring
    
    Permission is hereby granted, free of charge, to any person obtaining a
    copy of this software and associated documentation files (the "Software"),
    to deal in the Software without restriction, including without limitation
    the rights to use, copy, modify, merge, publish, distribute, sublicense,
    and/or sell copies of the Software, and to permit persons to whom the
    Software is furnished to do so.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
    THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
    FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
    DEALINGS IN THE SOFTWARE.
==============================================================================
"""

from sys import exception
import time
import multiprocessing as mp
import logging
from readerwriterlock import rwlock
import datetime
import LinMot_EtherCAT_Comm_0v82e as ECComm
import LinMot_Data_Handling_0v09 as DataH


#----------------------------------------------------------------------------------------------------
# Main Execution
class LinMot_EtherCATDemo():
    """
    Main class to set up, start, and monitor EtherCAT communication and basic motor operations for LinMot drives.

    This class manages the initialization and configuration of EtherCAT communication, 
    monitors and prints drive data, performs basic motor control (such as switching on, homing, and motion), 
    handles data acquisition and oscilloscope recording, and ensures safe shutdown and error handling.

    Attributes:
        adapter_id (str): The EtherCAT adapter identifier.
        noDev (int): Number of expected EtherCAT devices.
        cycle_time (float): Communication cycle time in seconds.
        no_Monitoring (int): Number of monitoring channels to receive.
        no_Parameter (int): Number of parameter channels to send.
        Activate_LMDrive_Data (bool): Flag to activate/deactivate LMDrive data.
        mp_logging (int): Logging level for multiprocessing communication.
        lock (mp.Lock): Lock to synchronize shared data access.
        lm_drive_data_dict (dict): Stores LMDrive data per device.
        data_length (int): Length of data expected per device.
        oszi_file_nr (int): Counter for oscilloscope file naming.
        timestamp_start_oszi (datetime.datetime): Timestamp when oscilloscope recording starts.
        unit_scale (list): List of unit scaling factors for each device.
        lm_drive_lock (rwlock.RWLockFairD): Read/write lock for LMDrive data access.
        sendData (LinMot_SendData): Instance for data handling.
        ProCommData (LinMot_ProCommData): Instance for process communication data.
        DC (LinMot_DriveCondition): Instance for drive condition monitoring.
        HK (LinMot_Housekeeping): Instance for housekeeping and basic drive control.
        MC (LinMot_MotionCommand): Instance for motion command handling.
        Oszi (LinMot_Oszilloscope): Instance for oscilloscope data handling.
        cfg (LinMot_Cfg): Instance for drive configuration handling.

    Methods:
        __init__():
            Initializes the LinMot_EtherCATDemo class with EtherCAT and LinMot drive configuration.
        start():
            Starts EtherCAT communication and runs the main test process.
        await_connection(deadline_s=30.0):
            Waits for the EtherCAT master to establish communication with the drive.
        initialize_drive_dict():
            Initializes the LMDrive data dictionary for each connected device.
        loop_print_data(max_cycles, t_sleep=1):
            Prints the communication data in a loop for a specified number of cycles.
        force_control():
            Executes a simple motion sequence and force control for the connected motor.
        wait_for_user(message, require_confirmation=False, timeout=None):
            Prompts the user for input or confirmation before proceeding, with optional timeout.
    """
    
    def __init__(self):
        """
        Initializes the MainTest class with EtherCAT and LinMot drive configuration.

        Sets up all configuration parameters, initializes locks, and creates instances
        of data handling, communication, and control classes required for EtherCAT
        communication and LinMot drive operations.
        """
        # Configuration parameters - Setup
        self.adapter_id = '\\Device\\NPF_{BE37777D-028C-4B31-A3CF-863CD3040A49}' # Replace with actual adapter ID
        self.noDev: int = 2 # Number of expected EtherCAT devices
        self.cycle_time: float = 0.050 # Cycle time in seconds
        self.no_Monitoring: int = 2 # How many Monitoring Channels do you want to receive.
        self.no_Parameter: int = 0 # How many Parameter Channels do you want to send
        self.Activate_LMDrive_Data: bool = False # This script works only when set to False
        self.mp_logging: int = 50 # Logging level for multiprocessing
        self.lock = mp.Lock() # Lock for synchronizing access to the data array
        # Optional Parameters for Realtime, which can be sent to ECComm
        #cpu_affinity = [2]
        #realtime = True
        #realtime_priority = 99

        # User defined Parameters 
        #self.forceControl_Chanel = 'mon_ch2' # Has to be mon_ch 1 to 4
        
        # Parameters
        self.lm_drive_data_dict = {}
        self.data_length = 0
        self.oszi_file_nr = 0
        self.timestamp_start_oszi = 0
        self.unit_scale = [None] * (self.noDev+1)
        #self.force_scale= [None] * (self.noDev+1)
        
        self.lm_drive_lock = rwlock.RWLockFairD()

        # Init Functions
        self.sendData = DataH.LinMot_SendData(self) # Create Data Handling Class
        self.ProCommData = DataH.LinMot_ProCommData(self) # Create Process Communication Data Class
        self.DC = DataH.LinMot_DriveCondition(self) # Create Drive Condition Class
        self.HK = DataH.LinMot_Housekeeping(self) # Create ForceControll Class
        self.MC = DataH.LinMot_MotionCommand(self) # Create Motion Command Class
        self.Oszi = DataH.LinMot_Oszilloscope(self) # Create Oscilloscope Class
        self.cfg = DataH.LinMot_Cfg(self) # Create Oscilloscope Class
        #self.FC = DataH.LinMot_ForceControl(self) # Create ForceControll Class

        # Constant

    def start(self):
        """
        Starts EtherCAT communication and runs the main test process.

        This method initializes EtherCAT communication, establishes a connection with
        the drives, and handles error checking. It runs the communication process,
        processes data, and performs motor operations like switching on, homing, and
        motion control.

        Raises:
            RuntimeError: If communication cannot be established or if the script is incorrectly configured.
            KeyboardInterrupt: If the process is interrupted manually.
        """
        print('\n=== LinMot start motion control example with 2 motors ===')
        print('=========================================================\n')

        user_input = self.wait_for_user(
            '> Have you read and understood the Quick Start Guide for this Software?\n'
            '> Has the Lin Mot Talk Wizard been completed?',
            require_confirmation=True
            )
        print('')
        
        if not user_input:
            print('Interrupted by user')
            quit()
        
        # Create EtherCATCommunication instance with configured parameters
        self.ethercat_comm = ECComm.EtherCATCommunication(self.adapter_id, self.noDev, 
                                                          self.cycle_time, self.lock, 
                                                          self.no_Monitoring, self.no_Parameter, 
                                                          self.Activate_LMDrive_Data, self.mp_logging)

        try:
            # Try to establish communication
            print('--- Connecting with Drive ---')
            self.ethercat_comm.start()
            self.await_connection()
            print('EtherCAT communication established. Press Ctrl+C to stop. \n')
            time.sleep(0.1)

            self.initialize_drive_dict()
            
            # These functions can be run in parallel using threading if needed
            self.simple_motion()

        except KeyboardInterrupt:
            logging.info('Keyboard interrupt received, stopping EtherCAT communication.')
            self.ethercat_comm.stop_event.set()
            
        except Exception as e:
            logging.error(f'Exception occurred: {e}')
            self.ethercat_comm.stop_event.set()
        
        finally:
            print('--- Stop EtherCAT communication. ---')
            # Print all Communication Messages
            while not self.ethercat_comm.error_queue.empty(): print(f'[EtherCAT] Error: {self.ethercat_comm.error_queue.get()}')
            while not self.ethercat_comm.info_queue.empty(): print(f'[EtherCAT] Info: {self.ethercat_comm.info_queue.get()}')
            # Ensure that the EtherCAT communication process is stopped properly
            logging.info("EtherCAT communication stopped.")
            self.ethercat_comm.stop()
            input("\n-> Press Enter to exit;")


    def await_connection(self):
        """
        Waits for the EtherCAT master to establish communication with the drive.

        Checks the communication process and waits until a connection is confirmed or
        a timeout occurs.

        Raises:
            RuntimeError: If communication cannot be established or the master setup fails.
        """
        if self.ethercat_comm.comm_proc and self.ethercat_comm.comm_proc.is_alive():
            # Wait until communication is confirmed
            attempt = 1
            EC_is_running = False
            while bool(attempt):
                EC_is_running = not self.ethercat_comm.stop_event.wait(timeout=1)
                print(f'Wait for the master to establish communication with the drive...')
                if not EC_is_running:
                    time.sleep(1)
                    attempt += 1
                    if attempt > 40:
                        EC_is_running = False
                        attempt = 0
                else:
                    attempt = 0
            if not EC_is_running:
                while not self.ethercat_comm.error_queue.empty(): print(f'[EtherCAT] Error: {self.ethercat_comm.error_queue.get()}')
                while not self.ethercat_comm.info_queue.empty(): print(f'[EtherCAT] Info: {self.ethercat_comm.info_queue.get()}')
                raise RuntimeError(f'Communication could not be established')
        else:
            raise RuntimeError(f'Master could not be set up with this port. Please make sure the selected port is correct.')
        
    def initialize_drive_dict(self):
        """
        Initializes the LMDrive data dictionary for each connected device.

        Creates LMDrive_Data instances for each device and sets the expected data length.

        Raises:
            RuntimeError: If Activate_LMDrive_Data is set to True.
        """
        if self.Activate_LMDrive_Data:
            raise RuntimeError(f'This script works only when Activate_LMDrive_Data is set to False')
        
        # Create LMDrive_Data
        for i in range(self.noDev):
            self.lm_drive_data_dict[i+1] = ECComm.LMDrive_Data(
                num_mon_channels=self.no_Monitoring, 
                num_par_channels=self.no_Parameter
                )
        self.data_length = self.ethercat_comm.InputLength
            
    def loop_print_data(self, max_cycles:int, t_sleep:float=1):
        """
        Prints the communication data in a loop.

        Fetches data from the EtherCAT communication process, processes it, and prints
        the data for each connected device.

        Args:
            max_cycles (int): The number of cycles to print the data.
            t_sleep (float, optional): Time to sleep between cycles in seconds. Defaults to 1.
        """
        cycle = 1
        while not self.ethercat_comm.stop_event.is_set() and cycle <= max_cycles:
            self.ProCommData.print_comm_messages()
            
            self.ProCommData.process_input_data()
            with self.lm_drive_lock.gen_rlock():
                for i in range(self.noDev):
                    print(f'[Status] Device {i+1}:')
                    print(f'\t{str(self.lm_drive_data_dict[i+1]).replace(",", " |")}')
            print(' ')
            cycle += 1

            if max_cycles != 1: # Do not sleep if max_cycles=1
                time.sleep(t_sleep)
    
    def simple_motion(self):
        """
        Executes a simple motion sequence for the connected motors.

        Handles switching on multipe motors, homing the motors, and moving the motors to a
        target position. Also manages motor control commands like enabling operation,
        homing, and stopping, and interacts with an oscilloscope for data logging.

        Raises:
            RuntimeError: If motion cannot be completed.
        """
        print('--- Status Drive(s) ---')
        # Print current motor Data
        self.loop_print_data(max_cycles=1)

        print('--- Enable Operation Motor 1 ---')
        # Switch on Motor 1
        self.wait_for_user(f'-> Please press Enter to switch on motor 1+2.') # Request user input
        self.HK.switch_on_motor(drive=[1,2])
        self.loop_print_data(max_cycles=1)
        
        # Home Motor 1
        self.wait_for_user("-> Please press Enter to home motor 1+2.") # Request user input
        self.HK.home_motor(drive=1)
        self.HK.home_motor(drive=2)
        self.loop_print_data(max_cycles=1)

        # Wait to make sure that everything is updated
        time.sleep(self.cycle_time * 3)
        
        # Start recording oscilloscope data
        self.ethercat_comm.data_queue_ON.set()
        self.timestamp_start_oszi = datetime.datetime.now() # Get current time, when data recording has been started
        
        print('--- Motion Motor 1 ---')
        # Move motor in positive direction
        self.wait_for_user("-> Press Enter to authorize movement sequence.")
        print('Sending quick move command to 3 mm...')
        self.MC.send_motion_command(drive=1, header='Absolute_VAI', target_pos=3, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=True)
        self.MC.motion_finished(drive=1)
        logging.info('Motion Finished (Motor 1): 3 mm\n')

        time.sleep(1)

        print('Sending quick move command to 3 mm...')
        cn0 = self.MC.send_motion_command(drive=2, header='Absolute_VAI', target_pos=3, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=True)
        self.MC.in_target_pos(drive=1, count_nibble=cn0)
        logging.info('In Target Position (Motor 2): 3 mm\n')

        time.sleep(1)

        print('Sending quick move command to 0 mm...')
        self.MC.send_motion_command(drive=1, header='Absolute_VAI', target_pos=0, 
                                    max_v=2, acc=1, dcc=1, jerk=10000, execute_mc=False)
        self.MC.send_motion_command(drive=2, header='Absolute_VAI', target_pos=0, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=False)
        self.sendData.send_data_to_slaves()
        self.MC.in_target_pos(drive=[1,2])
        logging.info('In Target Position (both motors): 0 mm\n')

        time.sleep(1)

        print('Sending quick move command to max. 5 mm...')
        self.MC.send_motion_command(drive=1, header='Absolute_VAI', target_pos=3, 
                                    max_v=2, acc=1, dcc=1, jerk=10000, execute_mc=False)
        self.MC.send_motion_command(drive=2, header='Absolute_VAI', target_pos=5, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=False)
        self.sendData.send_data_to_slaves()
        self.MC.in_target_pos(drive=[1,2])
        logging.info('In Target Position (both motors)\n')

        time.sleep(1)

        print('Sending quick move command to 0 mm...')
        cn1 = self.MC.send_motion_command(drive=1, header='Absolute_VAI', target_pos=0, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=False)
        cn2 = self.MC.send_motion_command(drive=2, header='Absolute_VAI', target_pos=0, 
                                    max_v=0.1, acc=0.1, dcc=0.1, jerk=10000, execute_mc=False)
        self.sendData.send_data_to_slaves()
        self.MC.motion_finished(drive=[1,2], count_nibble=[cn1,cn2])
        logging.info('Motion Finished (both motors): 0 mm\n')

        time.sleep(1)

        print('--- Save Oscilloscope ---')
        # Stop and save oscilloscope data
        self.ethercat_comm.data_queue_ON.clear()
        self.Oszi.save_oszi(filename = 'Queue_Save')
        print('')
        
        print('--- Disable Operation Motor 1 ---')
        self.wait_for_user(f'-> Please press Enter to switch off motor 1+2.') # Request user input
        # Switch off Motor
        self.HK.switch_off_motor(drive=[1,2])
        self.loop_print_data(max_cycles=1)
        print('')



    def wait_for_user(self, message, require_confirmation=False):
        """
        Prompts the user for input or confirmation before proceeding.

        Args:
            message (str): The message to display to the user.
            require_confirmation (bool, optional): Whether to require a 'yes' or 'no' confirmation. Defaults to False.

        Returns:
            bool: True if the user confirms or presses Enter, False otherwise.
        """
        if require_confirmation:
            print(message)
            response = input("-> Please type 'yes' or 'no' to continue: ")
            while response.lower() not in ["yes", "no"]:
                response = input("-> Please type 'yes' or 'no' to continue: ")
            if response.lower() == "yes":
                return True
            else:
                return False
        else:
            input(message)
            return True


def main():
    """
    Entry point for the LinMot EtherCAT demo application.

    Sets up multiprocessing, configures logging, creates the MainTest application
    instance, and starts the main process.
    """
    mp.set_start_method('spawn', force=True) # Important!
    logging.basicConfig(format='%(levelname)s:%(message)s', level=logging.DEBUG)
    app = LinMot_EtherCATDemo()
    app.start()

if __name__ == "__main__":
    main()
