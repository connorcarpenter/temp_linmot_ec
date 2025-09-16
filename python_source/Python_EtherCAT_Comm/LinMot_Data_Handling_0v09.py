"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           Data_Handling_ .py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  22.05.2025
    Version:        0.06
    Description:    Functions for Start_...py

    Disclaimer:
    ------------
    This is a demo project developed for NTI AG | LinMot & MagSpring. This
    software is provided "as-is" without any warranties or guarantees.
    NTI AG does not provide support, updates, or maintenance for this software. 
    Use it at your own risk.

    Dependencies:
    -------------
    - Python packages (see import section)

    Description:    
    ------------
    This script is responsible for processing data and managing communication 
    messages within the EtherCAT communication framework. It provides utility 
    functions for motor control and data communication with LinMot drives.
    It also includes the the LMDrive_Data class, which encapsulates all 
    relevant data for a LinMot drive, including communication I/O, motor 
    configuration parameters, and real-time scaled drive status.

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

import time
import os
import csv
import struct
import logging
import queue
#from ieee754_tools import ieee754_bits_to_float

class LinMot_ProCommData:
    """
    Handles EtherCAT communication data processing for LinMot drives.
    This class provides methods to print communication messages (errors and info) 
    from EtherCAT queues and to process input data from EtherCAT, updating the internal 
    state of all connected drives.
    """
    def __init__(self, app):
        self.app = app

    def print_comm_messages(self) -> None:
        """
        Prints communication messages from EtherCAT communication queues.

        Reads and logs error and info messages from the EtherCAT communication process.

        Args:
            None

        Returns:
            None
        """
        while True:
            try:
                msg = self.app.ethercat_comm.error_queue.get_nowait()
                logging.error(f"[EtherCAT] {msg}")
            except queue.Empty:
                break
            except Exception as e:
                logging.error(f"print_comm_messages(error) an error occurred: {e}")
        while True:
            try:
                msg = self.app.ethercat_comm.info_queue.get_nowait()
                logging.info(f"[EtherCAT] {msg}")
            except queue.Empty:
                break
            except Exception as e:
                logging.error(f"print_comm_messages(info) an error occurred: {e}")

    def process_input_data(self) -> None:
        """
        Processes input data from EtherCAT and updates internal device states.

        Locks shared data, unpacks device input, and updates calculated fields.

        Args:
            None

        Returns:
            None
        """
        with self.app.lock:
            all_slave_data = self.app.ethercat_comm.data[:]
        
        # Unpack and update the LMDrive input data per device
        for i in range(self.app.noDev):
            device_data = bytes(all_slave_data[i*self.app.data_length:(i+1)*self.app.data_length])
            with self.app.lm_drive_lock.gen_wlock():
                self.app.lm_drive_data_dict[i+1].unpack_inputs(device_data)
                self.app.lm_drive_data_dict[i+1].update_calculated_fields()


class LinMot_Housekeeping:
    """
    Manages basic drive operations and selection logic.
    This class offers utility functions for selecting active drives, switching motors on/off, 
    and sending homing commands. It ensures that drive operations are performed safely 
    and in the correct sequence.
    """
    def __init__(self, app):
        self.app = app

    def _drive_selection(self, drive:int|list) -> list:
        """
        Selects the active drive(s) based on input.

        Args:
            drive (int or list): The drive number or list of drive numbers.

        Returns:
            list: List of active drive numbers.

        Raises:
            ValueError: If 'drive' is not an integer or list.
        """
        if isinstance(drive, list):
            active_drive_number = drive
        elif isinstance(drive, int):
            active_drive_number = [drive]
        else:
            raise ValueError(f"'drive' must be an integer or a list. But it is a {type(drive)}: {drive}")
        return active_drive_number

    def switch_on_motor(self, drive:int|list) -> None:
        """
        Switch on the specified motor(s) if not already active.

        This method ensures that the given motor(s) are powered on. It waits until
        all specified motors report an "operation enabled" status.

        Args:
            drive (int | list): Drive number(s) of the motor(s) to switch on.

        Returns:
            None
        """
        active_drive_number = self._drive_selection(drive)
        # Switch motor ON if not already active
        self.app.ProCommData.process_input_data() # Receive most current data
        with self.app.lm_drive_lock.gen_rlock():
            motor_started = True
            for d in active_drive_number:
                motor_started = motor_started & self.app.lm_drive_data_dict[d].status['operation_enabled']
        if not motor_started:
            for d in active_drive_number:
                self.app.sendData.switchON_motor(d)
            
        while not motor_started: # Wait for motor to start
            time.sleep(self.app.cycle_time * 5)
            self.app.ProCommData.process_input_data() # Receive most current data
            with self.app.lm_drive_lock.gen_rlock():
                motor_started = True
                for d in active_drive_number:
                    motor_started = motor_started & self.app.lm_drive_data_dict[d].status['operation_enabled']
        logging.info(f'Motor {active_drive_number} switched on')

    def switch_off_motor(self, drive:int|list) -> None:
        """
        Switch off the specified motor(s) if currently active.

        This method ensures that the given motor(s) are powered off. It waits until
        all specified motors report that they are no longer in the "operation enabled" state.

        Args:
            drive (int | list): Drive number(s) of the motor(s) to switch off.

        Returns:
            None
        """
        active_drive_number = self._drive_selection(drive)
        # Switch Off Motor
        self.app.ProCommData.process_input_data()
        with self.app.lm_drive_lock.gen_rlock():
            motor_off = True
            for d in active_drive_number:
                motor_off = motor_off & (not self.app.lm_drive_data_dict[d].status['operation_enabled'])
        if not motor_off:
            for d in active_drive_number:
                self.app.sendData.switchOFF_motor(d)
        
        while not motor_off: # Wait for motor to turn off
            time.sleep(self.app.cycle_time * 5)
            self.app.ProCommData.process_input_data() # Receive most current data
            with self.app.lm_drive_lock.gen_rlock():
                motor_off = True
                for d in active_drive_number:
                    motor_off = motor_off & (not self.app.lm_drive_data_dict[d].status['operation_enabled'])
        logging.info(f'Motor {active_drive_number} switched off')
    
    def home_motor(self, drive:int|list) -> None:
        """
        Send a homing command to the specified motor(s) and wait for completion.

        If multiple motors are selected, all receive the homing command simultaneously.
        Note: Parallel homing can result in high power consumption.

        Args:
            drive (int | list): Drive number(s) of the motor(s) to home.

        Returns:
            None
        """
        active_drive_number = self._drive_selection(drive)
        # Send Home command
        self.app.ProCommData.process_input_data()
        with self.app.lm_drive_lock.gen_rlock():
            homing_started = [False] * len(active_drive_number)
            for i in range(len(active_drive_number)):
                homing_started[i] = ((self.app.lm_drive_data_dict[active_drive_number[i]].outputs['control_word'] & 0x0800) != 0)
        if False in homing_started:
            for i in range(len(active_drive_number)):
                if not homing_started[i]:
                    self.app.sendData.home_motor(active_drive_number[i], execute_now=False)
        self.app.sendData.send_data_to_slaves() # Send all Homing Commands simultaneously

        # Wait for Motor to home
        homing_finished = homing_started[:]
        while False in homing_finished:
            time.sleep(self.app.cycle_time * 3) # Longer wait time in order to make sure that the bits have updated
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                for i in range(len(active_drive_number)):
                    if not homing_finished[i]:
                        homing_finished[i] = self.app.lm_drive_data_dict[active_drive_number[i]].status['homed']

        # End homing procedure
        for i in range(len(active_drive_number)):
            if not homing_started[i]:
                self.app.sendData.end_home_motor(active_drive_number[i])
                logging.info(f'Motor {active_drive_number[i]} homed')


class LinMot_DriveCondition:
    """
    Provides condition-checking utilities for drive status and warning words.
    This class contains methods to check if specific bits or masked values in 
    the drive's status or warning words match given conditions, with optional 
    verification of a command count nibble.
    """
    def __init__(self, app):
        self.app = app
    
    def if_masked_status_word(self, drive:int, bit_mask:str|int, condition:str|int, count_nibble:int|None=None) -> bool:
        """
        Checks if the masked status word of a drive matches a given condition, with optional verification of a count nibble.

        This function processes the latest input data, applies a bit mask to the drive's status word, and compares the result
        to a specified condition. If a count nibble is provided, it further checks whether the lower nibble of the state variable
        matches the given count nibble.

        Args:
            drive (int): The drive number to check.
            bit_mask (str or int): The bit mask to apply to the status word.
            condition (str or int): The value to compare against the masked status word.
            count_nibble (int, optional): If provided, verifies that the lower nibble of the state variable matches this value.

        Returns:
            bool: True if the masked status word equals the condition (and count nibble matches if provided), False otherwise.
        """
        self.app.ProCommData.process_input_data()
        if self.app.lm_drive_data_dict[drive].inputs['status_word'] & bit_mask == condition: # If Status Word Bit Mask == Condition Value 
            if count_nibble is None:
                return True
            else:
                if self.app.lm_drive_data_dict[drive].inputs['state_var'] & 0x000F == count_nibble: # If "count nibble" correct
                    return True
                else:
                    return False
        else:
            return False

    def if_masked_warn_word(self, drive:int, bit_mask:str|int, condition:str|int, count_nibble:int|None=None) -> bool:
        """
        Checks if the masked warning word of a drive matches a given condition, with optional verification of a count nibble.

        This function processes the latest input data, applies a bit mask to the drive's warning word, and compares the result
        to a specified condition. If a count nibble is provided, it further checks whether the lower nibble of the state variable
        matches the given count nibble.

        Args:
            drive (int): The drive number to check.
            bit_mask (str or int): The bit mask to apply to the warning word.
            condition (str or int): The value to compare against the masked warning word.
            count_nibble (int, optional): If provided, verifies that the lower nibble of the state variable matches this value.

        Returns:
            bool: True if the masked warning word equals the condition (and count nibble matches if provided), False otherwise.
        """
        self.app.ProCommData.process_input_data()
        if self.app.lm_drive_data_dict[drive].inputs['warn_word'] & bit_mask == condition: # If Status Word Bit Mask == Condition Value 
            if count_nibble is None:
                return True
            else:
                if self.app.lm_drive_data_dict[drive].inputs['state_var'] & 0x000F == count_nibble: # If "count nibble" correct
                    return True
                else:
                    return False
        else:
            return False


class LinMot_MotionCommand:
    """
    Implements motion command logic for LinMot drives.
    This class allows sending various motion commands (absolute, relative, incremental, etc.) to drives, 
    checking if motions are finished, verifying target positions, and handling command acknowledgment. 
    It also includes input validation and helper methods for command nibble management.
    """
    def __init__(self, app):
        self.app = app

    def _check_input_count_nibble(self, active_drive_number:int|list, count_nibble:int|list|None=None) -> tuple:
        """
        Validates and prepares drive and count nibble inputs for motion verification.

        This function checks the types and lengths of the drive and count_nibble arguments,
        and prepares a list of count nibbles for further verification. If count_nibble is not provided,
        it reads the current count nibble from the drive(s).

        Args:
            active_drive_number (int or list): Drive number or list of drive numbers.
            count_nibble (int or list, optional): Command nibble(s) for verification.

        Returns:
            tuple: (drive_is_list, cn_list)
                drive_is_list (bool): True if input is a list of drives, False otherwise.
                cn_list (list or int): List of count nibbles or a single count nibble.

        Raises:
            TypeError: If input types or lengths are inconsistent.
            ValueError: If active_drive_number is neither an integer nor a list.
        """
        # Check input
        if isinstance(active_drive_number, list):
            drive_is_list = True
            if count_nibble is not None:
                if not isinstance(count_nibble, list) or len(active_drive_number) != len(count_nibble):
                    raise TypeError('active_drive_number({active_drive_number}) and count_nibble({count_nibble}) must both be lists of the same length.')
        elif isinstance(active_drive_number, int):
            drive_is_list = False
            if count_nibble is not None:
                if not isinstance(active_drive_number, int) or not isinstance(count_nibble, int):
                    raise TypeError('active_drive_number({active_drive_number}) and count_nibble({count_nibble}) must both be integers if not lists.')
        else:
            raise ValueError(f"'drive' must be an integer or a list. But it is a {type(active_drive_number)}: {active_drive_number}")
        
        # read required count_nibble if necessary
        if drive_is_list:
            cn_list = [None]*(max(active_drive_number)+1)
        if count_nibble is None:
            with self.app.lm_drive_lock.gen_rlock():
                if drive_is_list:
                    for d in active_drive_number:
                        cn_list[d] = (self.app.lm_drive_data_dict[d].outputs['mc_header']) & 0x000F
                else:
                    cn_list = self.app.lm_drive_data_dict[active_drive_number].outputs['mc_header'] & 0x000F
        else:
            if drive_is_list:
                for i in range(len(active_drive_number)):
                    cn_list[active_drive_number[i]] = count_nibble[i]
            else:
                cn_list = count_nibble
        
        return drive_is_list, cn_list

    def send_motion_command(self, drive:int, header:str, target_pos:float, max_v:float, 
                            acc:float, dcc:float=0, jerk:int=0, execute_mc:bool=True) -> int:
        """
        Sends a motion command to a specified drive with defined motion parameters.

        This function prepares and sends a motion command to the selected drive, using the specified motion profile.
        It supports various motion types (headers), and automatically scales parameters according to the drive's configuration.
        Required parameters depend on the selected motion header.

        Args:
            drive (int): Drive number to which the command is sent.
            header (str): Motion command type. Supported values include:
                'Absolute_VAI', 'Relative_VAI', 'Absolute_VAJI', 'Relative_VAJI',
                'Incr_Act_Pos_RstI', 'Absolute_Sin', 'Relative_Sin'.
            target_pos (float): Target position for the motion.
            max_v (float): Maximum velocity.
            acc (float): Acceleration value.
            dcc (float, optional): Deceleration value. Required for headers that do not combine acceleration.
            jerk (int, optional): Jerk value. Required for headers that require jerk.
            execute_mc (bool, optional): If True, executes the motion command immediately.

        Returns:
            int: Command nibble (lower 4 bits of the motion command header) used for verification.

        Raises:
            ValueError: If an unsupported motion header is provided.
            TypeError: If required arguments 'dcc' or 'jerk' are missing for specific headers.
        """
        # Get active Drive
        active_drive_number = int(drive)

        header_map = {
            "Absolute_VAI": (0x0100, False, False),
            "Relative_VAI": (0x0110, False, False),
            "Absolute_VAJI": (0x3A00, False, True),
            "Relative_VAJI": (0x3A10, False, True),
            "Incr_Act_Pos_RstI": (0x0D90, False, False),
            "Absolute_Sin": (0x0E00, True, False),
            "Relative_Sin": (0x0E10, True, False),
        }
        
        if header not in header_map:
            raise ValueError(f"Unsupported motion header: {header} (DriveNo. {active_drive_number})")
        
        header_code, acc_combined, jerk_necessary = header_map[header]
        if self.app.unit_scale[active_drive_number] is None:
            self.app.unit_scale[active_drive_number] = self.app.sendData.get_unit_scale(active_drive_number)

        pw = [
            [2, float(target_pos) * self.app.unit_scale[active_drive_number]],
            [2, float(max_v) * self.app.unit_scale[active_drive_number] * 100],
            [2, float(acc) * self.app.unit_scale[active_drive_number] * 10],
            [0], [0]
        ]
        if not acc_combined:
            if dcc == 0:
                raise TypeError(f"Missing required argument 'dcc' with header '{header}' (DriveNo. {active_drive_number})")
            else:
                pw[3] = ([2, float(dcc) * self.app.unit_scale[active_drive_number] * 10])
        if jerk_necessary:
            if jerk == 0:
                raise TypeError(f"Missing required argument 'jerk' with header '{header}' (DriveNo. {active_drive_number})")
            else:
                pw[4] = ([2, float(jerk) * self.app.unit_scale[active_drive_number]])

        count_nibble = self.app.sendData.update_output_drive_data(active_drive_number, controlWord=0, header = header_code,
                                        para_word=pw, execute_mc=execute_mc)
        return count_nibble
        
    def motion_finished(self, drive:int|list, count_nibble=None, do_not_wait:bool=False, timeout:float=60*5) -> bool:
        """
        Checks whether the motion has finished for the specified drive(s).

        This function waits until the motion is complete for the given drive or list of drives.
        It verifies that the "motion active" status is cleared and, if provided, that the count nibble matches.
        Can optionally return immediately if do_not_wait is True.

        Args:
            drive (int or list): Drive number or list of drive numbers to check.
            count_nibble (int or list, optional): Command nibble(s) for verification.
            do_not_wait (bool, optional): If True, returns immediately without waiting.
            timeout (float, optional): Maximum time to wait for motion completion in seconds.

        Returns:
            bool: True if motion is complete and count nibble matches (if provided), False otherwise.

        Raises:
            TypeError: If input types for drive and count_nibble are inconsistent.
            ValueError: If drive is neither an integer nor a list.
            TimeoutError: If motion does not complete within the specified timeout.
        """
        start_time = time.time()
        active_drive_number = drive

        # Check input and read required count_nibble if necessary
        drive_is_list, cn_list = self._check_input_count_nibble(active_drive_number, count_nibble)

        # Wait for Motion to finish
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                if drive_is_list:
                    if all(not self.app.lm_drive_data_dict[d].status['motion_active'] for d in active_drive_number):
                        if all((cn_list[d] == (self.app.lm_drive_data_dict[d].inputs['state_var'] & 0x000F)) for d in active_drive_number):
                            return True

                else:
                    if not self.app.lm_drive_data_dict[active_drive_number].status['motion_active']:
                        if cn_list == self.app.lm_drive_data_dict[active_drive_number].inputs['state_var'] & 0x000F:
                            return True

            if do_not_wait:
                return False
            else:
                time.sleep(self.app.cycle_time * 2)

        raise TimeoutError("motion_finished() did not finish within timeout. (DriveNo. {active_drive_number})")
    
    def in_target_pos(self, drive:int|list, count_nibble=None, do_not_wait:bool=False, timeout:float=60*5) -> bool:
        """
        Checks whether the drive(s) have reached the target position.

        This function waits until the drive(s) report the "in target position" status and, if provided,
        that the count nibble matches. Can optionally return immediately if do_not_wait is True.

        Args:
            drive (int or list): Drive number or list of drive numbers to check.
            count_nibble (int or list, optional): Command nibble(s) for verification.
            do_not_wait (bool, optional): If True, returns immediately without waiting.
            timeout (float, optional): Maximum time to wait for position confirmation in seconds.

        Returns:
            bool: True if the drive(s) are in target position and count nibble matches (if provided), False otherwise.

        Raises:
            TypeError: If input types for drive and count_nibble are inconsistent.
            ValueError: If drive is neither an integer nor a list.
            TimeoutError: If target position is not reached within the specified timeout.
        """
        start_time = time.time()
        active_drive_number = drive

        # Check input and read required count_nibble if necessary
        drive_is_list, cn_list = self._check_input_count_nibble(active_drive_number, count_nibble)

        # Wait for Motion to finish
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                if drive_is_list:
                    if all(self.app.lm_drive_data_dict[d].inputs['status_word'] & 0x0400 for d in active_drive_number):
                        if all((cn_list[d] == (self.app.lm_drive_data_dict[d].inputs['state_var'] & 0x000F)) for d in active_drive_number):
                            return True

                else:
                    if self.app.lm_drive_data_dict[active_drive_number].inputs['status_word'] & 0x0400:
                        if cn_list == self.app.lm_drive_data_dict[active_drive_number].inputs['state_var'] & 0x000F:
                            return True

            if do_not_wait:
                return False
            else:
                time.sleep(self.app.cycle_time * 2)

        raise TimeoutError("in_target_pos() did not finish within timeout. (DriveNo. {active_drive_number})")

    def in_pos_range(self, drive:int, position:float, range:float=1, timeout:float=60*5, do_not_wait:bool=False) -> bool:
        """
        Checks whether the actual position of the drive is within a specified range of the required position.

        This function waits until the drive's actual position is within the specified range of the target position.
        Can optionally wait until the condition is met or timeout occurs.

        Args:
            drive (int): Drive number to check.
            position (float): Target position to compare against.
            range (float, optional): Acceptable range around the target position.
            timeout (float, optional): Maximum time to wait for position confirmation in seconds.
            wait (bool, optional): If True, waits until the condition is met or timeout occurs.

        Returns:
            bool: True if the actual position is within the specified range, False otherwise.

        Raises:
            TimeoutError: If the position condition is not met within the specified timeout.
        """
        start_time = time.time()
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                actual_pos = self.app.lm_drive_data_dict[drive].status['actual_position']
            if actual_pos - range <= position <= actual_pos + range:
                return True
            
            if not do_not_wait:
                time.sleep(self.app.cycle_time)
            else:
                return False
        raise TimeoutError("in_pos_range() did not finish within timeout. (DriveNo. {active_drive_number})")
    
    def command_received_by_drive(self, drive:int, count_nibble:int, do_not_wait:bool=False) -> bool:
        """
        Waits until the specified drive acknowledges the received motion command.

        This function checks if the drive's state variable lower nibble matches the given count nibble,
        indicating that the command has been received by the drive.
        Can optionally return immediately if do_not_wait is True.

        Args:
            drive (int): Drive number to check.
            count_nibble (int): Command nibble to verify acknowledgment.
            do_not_wait (bool, optional): If True, returns immediately without waiting.

        Returns:
            bool: True if the command has been received by the drive, False otherwise.

        Raises:
            TimeoutError: If acknowledgment is not received within the expected time.
        """
        for _ in range(8):
            time.sleep(self.app.cycle_time)
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                if self.app.lm_drive_data_dict[drive].inputs['state_var'] & 0x000F == count_nibble:
                    return True
                elif do_not_wait:
                    return False
        raise TimeoutError(f'Motion Command has not been received by drive {drive}.')


class LinMot_ForceControl:
    """
    Implements force-controlled motion and force feedback logic.
    This class provides methods for sending force-controlled motion commands, 
    changing force setpoints, reading measured force, waiting for force conditions, 
    and checking if the measured force is within a specified range. 
    It also includes logic for verifying special motion states.
    """

    
    def __init__(self, app):
        self.app = app
    
    def _check_input_count_nibble(self, active_drive_number, count_nibble) -> tuple:
        """This is a copy of the same function in class MotionCommand!"""
        # Check input
        if isinstance(active_drive_number, list):
            drive_is_list = True
            if count_nibble is not None:
                if not isinstance(count_nibble, list) or len(active_drive_number) != len(count_nibble):
                    raise TypeError('active_drive_number({active_drive_number}) and count_nibble({count_nibble}) must both be lists of the same length.')
        elif isinstance(active_drive_number, int):
            drive_is_list = False
            if count_nibble is not None:
                if not isinstance(active_drive_number, int) or not isinstance(count_nibble, int):
                    raise TypeError('active_drive_number({active_drive_number}) and count_nibble({count_nibble}) must both be integers if not lists.')
        else:
            raise ValueError(f"'drive' must be an integer or a list. But it is a {type(active_drive_number)}: {active_drive_number}")
        
        # read required count_nibble if necessary
        if drive_is_list:
            cn_list = [None]*(max(active_drive_number)+1)
        if count_nibble is None:
            with self.app.lm_drive_lock.gen_rlock():
                if drive_is_list:
                    for d in active_drive_number:
                        cn_list[d] = (self.app.lm_drive_data_dict[d].outputs['mc_header']) & 0x000F
                else:
                    cn_list = self.app.lm_drive_data_dict[active_drive_number].outputs['mc_header'] & 0x000F
        else:
            if drive_is_list:
                for i in range(len(active_drive_number)):
                    cn_list[active_drive_number[i]] = count_nibble[i]
            else:
                cn_list = count_nibble
        
        return drive_is_list, cn_list

    def motion_force_control(self, drive:int, header:str, target_pos:float, max_v:float, 
                            acc:float, dcc:float=0, target_force:float=0, force_limit:float=0, execute_mc:bool=True) -> int:
        """
        Sends a force-controlled motion command to the specified drive.

        This function prepares and sends a force-controlled motion command to the selected drive,
        using the specified motion profile and force parameters. Required parameters depend on the
        selected motion header.

        Args:
            drive (int): Drive number.
            header (str): Motion command type. Supported values include:
                'VAI Go To Pos With Higher Force Ctrl Limit and Target Force',
                'VAI Go To Pos With Lower Force Ctrl Limit and Target Force',
                'VAI Inc Act Pos With Higher Force Ctrl Limit and Target Force',
                'VAI Inc Act Pos With Lower Force Ctrl Limit and Target Force',
                'VAI Go To Pos From Act Pos And Reset Force Control Set I',
                'VAI Increment Act Pos And Reset Force Control Set I'.
            target_pos (float): Target position.
            max_v (float): Maximum velocity.
            acc (float): Acceleration.
            dcc (float, optional): Deceleration. Required for headers that require deceleration.
            target_force (float, optional): Desired force. Required for headers that require target force.
            force_limit (float, optional): Force limit. Required for headers that require force limit.
            execute_mc (bool, optional): If True, executes the motion command immediately.

        Returns:
            int: Count nibble used for verification.

        Raises:
            ValueError: If an unsupported motion header is provided.
            TypeError: If required arguments 'dcc', 'target_force', or 'force_limit' are missing for specific headers.
        """
        # Get active Drive
        active_drive_number = int(drive)
        if self.app.force_scale[drive] is None:
            self.app.force_scale[drive] = self.app.sendData.get_force_scale(active_drive_number)

        header_map = {
            "VAI Go To Pos With Higher Force Ctrl Limit and Target Force": (0x3830, True, True, False),
            "VAI Go To Pos With Lower Force Ctrl Limit and Target Force": (0x3850, True, True, False),
            "VAI Inc Act Pos With Higher Force Ctrl Limit and Target Force": (0x3880, True, True, False),
            "VAI Inc Act Pos With Lower Force Ctrl Limit and Target Force": (0x3890, True, True, False),
            "VAI Go To Pos From Act Pos And Reset Force Control Set I": (0x3860, False, False, True),
            "VAI Increment Act Pos And Reset Force Control Set I": (0x3870, False, False, True)
        }
        
        if header not in header_map:
            raise ValueError(f"Unsupported motion header: {header} (DriveNo. {active_drive_number})")
        
        header_code, target_required, limit_required, dcc_required = header_map[header]
        if self.app.unit_scale[active_drive_number] is None:
            self.app.unit_scale[active_drive_number] = self.app.sendData.get_unit_scale(active_drive_number)

        pw = [
            [2, float(target_pos) * self.app.unit_scale[active_drive_number]],
            [2, float(max_v) * self.app.unit_scale[active_drive_number] * 100],
            [2, float(acc) * self.app.unit_scale[active_drive_number] * 10],
            [0], [0]
        ]
        
        if limit_required:
            if force_limit == 0:
                raise TypeError(f"Missing required argument 'force_limit' with header '{header}'")
            else:
                pw[3] = ([1, int(float(force_limit) / self.app.force_scale[active_drive_number])])
        if dcc_required:
            if dcc == 0:
                raise TypeError(f"Missing required argument 'dcc' with header '{header}'")
            else:
                pw[3] = ([2, float(dcc) * self.app.unit_scale[active_drive_number] * 10])
        if target_required:
            if target_force == 0:
                raise TypeError(f"Missing required argument 'target_force' with header '{header}'")
            else:
                pw[4] = ([1, int(float(target_force) / self.app.force_scale[active_drive_number])])
        

        count_nibble = self.app.sendData.update_output_drive_data(active_drive_number, controlWord=0, header = header_code,
                                        para_word=pw, execute_mc=execute_mc)
        return count_nibble
        
    def force_control(self, drive:int, header:str, force:float|None=None, execute_mc:bool=True) -> int:
        """
        Sends a force control command to the specified drive.

        This function sends a force control command to the drive, such as changing the target force
        or resetting force control. The required arguments depend on the selected header.

        Args:
            drive (int): Drive number.
            header (str): Command type. Supported values include:
                'Change_Target_Force', 'Reset_Force_Ctrl'.
            force (float, optional): Target force. Required for 'Change_Target_Force'.
            execute_mc (bool, optional): If True, executes the command immediately.

        Returns:
            int: Command nibble used for verification.

        Raises:
            ValueError: If an unsupported motion header is provided.
            TypeError: If 'force' is missing when required by the header.
        """
        # Get active Drive
        active_drive_number = int(drive)

        header_map = {
            "Change_Target_Force": (0x3820, True),
            "Reset_Force_Ctrl": (0x3870, False)
        }
        
        if header not in header_map:
            raise ValueError(f"Unsupported motion header: {header} (DriveNo. {active_drive_number})")
        
        header_code, force_necessary = header_map[header]
        if self.app.force_scale[active_drive_number] is None:
            self.app.force_scale[active_drive_number] = self.app.sendData.get_force_scale(active_drive_number)

        if sum([force_necessary]) == 0:
            pw = [[0]]
        else:
            pw = [0] * sum([force_necessary])
        
        if force_necessary:
            if force is None:
                raise TypeError(f"Missing required argument 'force' with header '{header}'")
            else:
                pw[0] = [1, int(float(force) / self.app.force_scale[active_drive_number])]

        count_nibble = self.app.sendData.update_output_drive_data(active_drive_number, controlWord=0, header = header_code,
                                        para_word=pw, execute_mc=execute_mc)
        return count_nibble

    def get_measured_force(self, drive:int) -> float:
        """
        Retrieves the measured force from the specified drive.
        It differentiates between C- and F- Drive Series.

        This function reads the current measured force value 
        from the drive and applies the appropriate scaling.

        Args:
            drive (int): Drive number.

        Returns:
            float: Measured force value.
        """
        if self.app.force_scale[drive] is None and self.app.drive_Nr[drive] <= 2:
            self.app.force_scale[drive] = self.app.sendData.get_force_scale(drive)
        with self.app.lm_drive_lock.gen_rlock():
            force = self.app.lm_drive_data_dict[drive].inputs[self.app.forceControl_Channel]
        if self.app.drive_Nr[drive] <= 2:
            force = self.app.sendData.unsigned_to_signed_16bit(value=force) * self.app.force_scale[drive]
        elif self.app.drive_Nr[drive] <= 4:
            force = self.app.sendData.ieee754_bits_to_float(force, "single")
        else:
            raise TypeError('Drive type not in List (self.drive_dict): Drive {drive}')
        return force
    
    def wait_force_target(self, drive:int, exit_condition:str, force_target:float, timeout:float=60*5, do_not_wait:bool=False) -> bool:
        """
        Waits until the measured force meets the specified exit condition.

        This function waits until the measured force on the drive satisfies the specified comparison
        (e.g., '>=', '<', etc.) with respect to the target force. Can optionally return immediately
        if do_not_wait is True.

        Args:
            drive (int): Drive number.
            exit_condition (str): Comparison operator. Supported values: '>', '>=', '<', '<='.
            force_target (float): Target force value to compare against.
            timeout (float, optional): Maximum time to wait in seconds.
            do_not_wait (bool, optional): If True, returns immediately without waiting.

        Returns:
            bool: True if the condition is met, False otherwise.

        Raises:
            ValueError: If an unsupported exit condition is provided.
            TimeoutError: If the condition is not met within the specified timeout.
        """
        start_time = time.time()
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            force = self.get_measured_force(drive)
            match exit_condition:
                case '>':
                    if force > force_target:
                        return True
                case '>=':
                    if force >= force_target:
                        return True
                case '<':
                    if force < force_target:
                        return True
                case '<=':
                    if force <= force_target:
                        return True
                case '==':
                    raise ValueError(f'Please use greater or smaler than for exit_condition')
                case _:
                    raise ValueError(f'exit_condition not properly defined')
            if do_not_wait:
                return False
            else:
                time.sleep(self.app.cycle_time)
        raise TimeoutError("wait_force_target() did not finish within timeout. (DriveNo. {active_drive_number})")
    
    def force_range(self, drive:int, force:float, range:float, timeout:float=60*5, wait:bool=False) -> bool:
        """
        Checks whether the measured force is within a specified range of the target force.

        This function waits until the measured force is within the specified range of the target force.
        Can optionally wait until the condition is met or timeout occurs.

        Args:
            drive (int): Drive number.
            force (float): Target force value.
            range (float): Acceptable range around the target force.
            timeout (float, optional): Maximum time to wait in seconds.
            wait (bool, optional): If True, waits until the condition is met or timeout occurs.

        Returns:
            bool: True if the measured force is within the specified range, False otherwise.

        Raises:
            TimeoutError: If the condition is not met within the specified timeout.
        """
        start_time = time.time()
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            actual_force = self.get_measured_force(drive)
            if actual_force - range <= force <= actual_force + range:
                return True
            
            if wait:
                time.sleep(self.app.cycle_time)
            else:
                return False
        raise TimeoutError("force_range() did not meet conditions within timeout. (DriveNo. {active_drive_number})")

    def special_motion_active(self, drive:int|list, count_nibble=None, do_not_wait:bool=False, timeout:float=60*5) -> bool:
        """
        Checks whether the "special motion" is active for the specified drive(s).

        This function waits until the special motion status is active for the given drive or list of drives,
        and optionally verifies the count nibble. Can optionally return immediately if do_not_wait is True.

        Args:
            drive (int or list): Drive number or list of drive numbers.
            count_nibble (int or list, optional): Command nibble(s) for verification.
            do_not_wait (bool, optional): If True, returns immediately without waiting.
            timeout (float, optional): Maximum time to wait in seconds.

        Returns:
            bool: True if the special motion is active and verified, False otherwise.

        Raises:
            TypeError: If input types for drive and count_nibble are inconsistent.
            ValueError: If drive is neither an integer nor a list.
            TimeoutError: If the special motion does not activate within the specified timeout.
        """
        start_time = time.time()
        active_drive_number = drive

        # Check input and read required count_nibble if necessary
        drive_is_list, cn_list = self._check_input_count_nibble(active_drive_number, count_nibble)

        # Wait for Motion to finish
        while time.time() - start_time < timeout:
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                if drive_is_list:
                    if all(self.app.lm_drive_data_dict[d].inputs['status_word'] & 0x0200 for d in active_drive_number):
                        if all((cn_list[d] == (self.app.lm_drive_data_dict[d].inputs['state_var'] & 0x000F)) for d in active_drive_number):
                            return True

                else:
                    if self.app.lm_drive_data_dict[active_drive_number].inputs['status_word'] & 0x0200:
                        if cn_list == self.app.lm_drive_data_dict[active_drive_number].inputs['state_var'] & 0x000F:
                            return True

            if do_not_wait:
                return False
            else:
                time.sleep(self.app.cycle_time * 2)

        raise TimeoutError("special_motion_active() did not finish within timeout. (DriveNo. {active_drive_number})")


class LinMot_Oszilloscope:
    """
    Handles oscilloscope data acquisition and storage.
    This class provides methods to save raw or structured oscilloscope data 
    from EtherCAT communication into CSV files, either as a single file or 
    split per device. It also includes a method to unpack binary input data 
    into a structured dictionary.
    """
    def __init__(self, app):
        self.app = app

    def save_oszi_simple(self, filename:str='Oszi_recoding') -> None:
        """
        Saves raw oscilloscope data to a single CSV file.

        This function drains the EtherCAT oscilloscope data queue and writes all collected data
        into a single CSV file. If a file with the same name already exists, it will be deleted first.
        The function also increments the internal file number counter after saving.

        Args:
            filename (str): Name of the output file.

        Returns:
            None

        Side Effects:
            - Deletes any existing file with the same name before saving.
            - Increments self.app.oszi_file_nr after saving.
        """
        # Delete file if it exists
        if os.path.exists(f'{filename}_{self.app.oszi_file_nr}.csv'):
            os.remove(f'{filename}_{self.app.oszi_file_nr}.csv')
            logging.info(f"Existing file '{f'{filename}_{self.app.oszi_file_nr}.csv'}' deleted.")
        
        #Start with save procedure
        data = []
        
        # Drain the queue
        while not self.app.ethercat_comm.data_queue.empty():
            data.append(self.app.ethercat_comm.data_queue.get())

        if not data:
            logging.warning("Queue is empty. Nothing to save.")
            return

        with open(f'{filename}_{self.app.oszi_file_nr}.csv', 'w', newline='') as f:
            writer = csv.writer(f)
            writer.writerows(data)
        
        self.app.oszi_file_nr +=1

        logging.info(f"Saved {len(data)} items to {filename}")

    def save_oszi(self, filename:str='') -> None:
        """
        Saves oscilloscope data to separate CSV files per device.

        This function drains the EtherCAT oscilloscope data queue and writes the data for each device
        into a separate CSV file within a dedicated output directory. Each file contains a header row
        and all samples for the corresponding device. If a file already exists, it will be deleted first.
        The function also increments the internal file number counter after saving.

        Args:
            filename (str, optional): Base name for the output files. If not provided, defaults to 'Oszi_recoding'.

        Returns:
            None

        Side Effects:
            - Creates an output directory for the current save operation.
            - Deletes any existing device files before saving.
            - Logs status and error messages.
            - Increments self.app.oszi_file_nr after saving.
        """
        if filename == '':
            filename = 'Oszi_recoding'

        # Drain oszi queue
        data_with_timestamps = []
        try:
            while True:
                data_with_timestamps.append(self.app.ethercat_comm.data_queue.get_nowait())
        except queue.Empty:
            pass
        except Exception as e:
                logging.error(f"save_oszi() error: {e}")

        if not data_with_timestamps:
            logging.info("Queue is empty. Nothing to save.")
            return
        
        # Create output directory if it doesn't exist
        output_dir = f'{filename}_{self.app.oszi_file_nr}'
        os.makedirs(output_dir, exist_ok=True)
        
        # Unpack and write to separate CSV files for each device
        for device_index in range(self.app.noDev):
            device_filename = os.path.join(output_dir, f'{filename}_device_{device_index + 1}.csv')
            csv_data = []
            header_written = False
            first_row = True
            
            for sample_nr, raw_data in data_with_timestamps:
                # Ensure raw_data is a bytes-like object
                if isinstance(raw_data, list):  # Convert if it's a list
                    raw_data = bytes(raw_data)
                
                # Extract the data for the current device based on its index
                device_data_chunk = raw_data[device_index * self.app.data_length:(device_index + 1) * self.app.data_length]
                unpacked_dict = self._unpack_input_data(device_data_chunk)
                
                # Write the header once and then the data for this device
                if not header_written:
                    csv_data.append(['Sample_Nr'] + ['Timestamp'] + list(unpacked_dict.keys()))
                    header_written = True
                if first_row: # Please refer to the documentation for instructions on how to generate an accurate timestamp.
                    csv_data.append([sample_nr] + [self.app.timestamp_start_oszi] + list(unpacked_dict.values()))
                    first_row = False
                else:
                    csv_data.append([sample_nr] + [''] + list(unpacked_dict.values()))
            
            # Write the CSV data for this device
            try:
                if os.path.exists(device_filename):
                    os.remove(device_filename)
                    logging.info(f"Deleted existing file '{device_filename}'.")

                with open(device_filename, 'w', newline='') as f:
                    writer = csv.writer(f)
                    writer.writerows(csv_data)

                logging.info(f"Saved {len(data_with_timestamps)} entries to '{device_filename}'")

            except Exception as e:
                logging.error(f"Failed to save oszi for device {device_index + 1}: {e}")
        
        # Increment file number for the next time
        self.app.oszi_file_nr += 1

    def _unpack_input_data(self, data:bytes) -> dict:
        """
        Unpacks binary input data into a structured dictionary.

        This function interprets the raw binary input data according to the expected format,
        extracting all relevant fields and returning them as a dictionary with descriptive keys.

        Args:
            data (bytes): Raw binary data to be unpacked.

        Returns:
            dict: Parsed data fields, including state variables, status words, configuration values,
                and all monitoring channels.
        """
        base_format = '<HHHiiiHHi'
        mon_channel_format = 'i' * self.app.no_Monitoring
        full_format = base_format + mon_channel_format

        unpacked = struct.unpack(full_format, data)

        keys = [
            'state_var',
            'status_word',
            'warn_word',
            'demand_pos',
            'actual_pos',
            'demand_curr',
            'cfg_status',
            'cfg_index_in',
            'cfg_value_in'
        ] + [f'mon_ch{i}' for i in range(1, self.app.no_Monitoring + 1)]

        return dict(zip(keys, unpacked))


class LinMot_Cfg:
    """
    Manages configuration parameter commands for drives.
    This class enables reading and writing configuration parameters (via UPID) 
    to LinMot drives, supporting both RAM and ROM operations.
    It includes logic for command execution and acknowledgment handling.
    """
    def __init__(self, app):
        self.app = app

    def _send_parameter_command(self, drive:int, header:str, UPID:str|int, value:int|None=None, execute_mc:bool=True) -> int:
        """
        Sends a configuration parameter command to the drive.

        This function prepares and sends a configuration command (read or write) to the specified drive,
        using the provided header and parameter ID (UPID). For write commands, a value can be specified.
        The command is executed immediately if execute_mc is True.

        Args:
            drive (int): Drive number.
            header (str): Command type. Supported values:
                'Read_Value_ROM', 'Read_Value_RAM', 'Write_Value_ROM',
                'Write_Value_RAM', 'Write_Value_RAM_and_ROM'.
            UPID (str or int): Parameter ID to read or write.
            value (int, optional): Value to write (required for write commands).
            execute_mc (bool, optional): Whether to execute the command immediately.

        Returns:
            int: Status nibble (lower 4 bits of the configuration control word) used for verification.

        Raises:
            ValueError: If an unsupported header is provided.
        """
        # Get active Drive
        active_drive_number = int(drive)

        header_map = {
            "Read_Value_ROM": (0x1000, False),
            "Read_Value_RAM": (0x1100, False),
            "Write_Value_ROM": (0x1200, True),
            "Write_Value_RAM": (0x1300, True),
            "Write_Value_RAM_and_ROM": (0x1400, True)
        }
        cfg_control, value_out = header_map[header]
        cfg_index_out = UPID
        if value_out:
            cfg_value_out = value
        else:
            cfg_value_out = None
        
        # Return Status Value of this command
        return self.app.sendData.update_output_cfg(active_drive_number, cfg_control, cfg_index_out, cfg_value_out, execute_mc=execute_mc)

    def read_cfg(self, drive:int, header:str, UPID:str|int) -> int:
        """
        Reads a configuration value from the drive.

        This function sends a read command to the specified drive and waits for the corresponding value
        to be received. It returns the value of the requested parameter.

        Args:
            drive (int): Drive number.
            header (str): Read command type. Supported values: 'Read_Value_ROM', 'Read_Value_RAM'.
            UPID (str or int): Parameter ID to read.

        Returns:
            int: The value read from the drive for the specified parameter.

        Raises:
            ValueError: If the header is not a supported read command.
            TimeoutError: If the value is not received within the expected time.
        """
        if header == "Read_Value_ROM" or header == "Read_Value_RAM":
            status = self._send_parameter_command(drive=drive, header=header, UPID=UPID, value=None, execute_mc=True)
        else:
            raise ValueError(f"Parameter command '{header}' not valid. (DriveNo. {drive})")
        time.sleep(self.app.cycle_time * 2)
        for i in range(6):
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                status_new = self.app.lm_drive_data_dict[drive].inputs['cfg_status']
                if status_new == status:
                    return self.app.lm_drive_data_dict[drive].inputs['cfg_value_in']
            time.sleep(self.app.cycle_time)
        raise TimeoutError(f'read_cfg() did not receive the corresponding value in time.\n'
                           f'     UPID = {UPID}. DriveNo = {drive}\n'
                           f'     Status Old = {status:04X}h, Status new = {status_new:04X}h')

    def write_cfg(self, drive:int, header:str, UPID:str|int, value:int) -> None:
        """
        Writes a configuration value to the drive and waits for acknowledgment.

        This function sends a write command to the specified drive with the given parameter ID (UPID)
        and value. It waits until the drive acknowledges the command or raises a TimeoutError if
        acknowledgment is not received in time.

        Args:
            drive (int): Drive number.
            header (str): Write command type. Supported values:
                'Write_Value_ROM', 'Write_Value_RAM', 'Write_Value_RAM_and_ROM'.
            UPID (str or int): Parameter ID to write.
            value (int): Value to write.

        Returns:
            None

        Raises:
            ValueError: If the header is not a supported write command.
            TimeoutError: If acknowledgment is not received within the expected time.
        """
        if header == "Write_Value_ROM" or header == "Write_Value_RAM" or header == "Write_Value_RAM_and_ROM":
            status = self._send_parameter_command(drive=drive, header=header, UPID=UPID, value=value, execute_mc=True)
        else:
            raise ValueError(f"Parameter command '{header}' not valid. (DriveNo. {drive})")
        time.sleep(self.app.cycle_time * 2)
        loop_a = 0
        for _ in range(6): # Wait until the drive has received the command
            self.app.ProCommData.process_input_data()
            with self.app.lm_drive_lock.gen_rlock():
                status_new = self.app.lm_drive_data_dict[drive].inputs['cfg_status']
                if status_new == status:
                    return
            time.sleep(self.app.cycle_time)
        raise TimeoutError(f'write_cfg() did not receive the OK from drive.\n'
                           f'     UPID = {UPID}. DriveNo = {drive}\n'
                           f'     Status Old = {status:04X}h, Status new = {status_new:04X}h')


class LinMot_SendData:
    """
    Handles low-level data transmission to LinMot drives.
    This class contains methods for sending control words, motion commands, 
    and configuration data to drives via EtherCAT. It also provides utility 
    functions for bit manipulation, scaling, and packing/unpacking data 
    for communication.
    """
    def __init__(self, app):
        self.app = app

    # Motor Control Functions --------------------------------------------------------

    def switchON_motor(self, active_drive_number:int):
        """
        Turns the motor ON by manipulating bit 0 (Switch ON) in the control word.

        Functionality:
            - Sends the current control word to the drive.
            - Clears bit 0 (Switch ON = 0).
            - Sends the updated control word.
            - Sets bit 0 (Switch ON = 1).
            - Sends the final updated control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        #self.send_data_to_slaves() # Send Current Control Word
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0001 # Clear bit 0 (Switch ON = 0)
        self.send_data_to_slaves() # Send Current Control Word
        time.sleep(max(self.app.cycle_time * 2, 0.001))
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] |= 0x0001 # Set bit 0 (Switch ON = 1)
        self.send_data_to_slaves()# Send Current Control Word

    def switchOFF_motor(self, active_drive_number:int):
        """
        Turns the motor OFF by clearing bit 0 (Switch ON) in the control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0001 # Clear bit 0 (Switch ON = 0)
        self.send_data_to_slaves()

    def home_motor(self, active_drive_number:int, execute_now:bool=True):
        """
        Starts the homing procedure by setting bit 11 (Home) in the control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] |= 0x0800 # Set bit 11 (Home = 1)
        if execute_now:
            self.send_data_to_slaves()

    def end_home_motor(self, active_drive_number:int):
        """
        Ends the homing procedure by clearing bit 11 (Home) in the control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0800 # Clear bit 11 (Home = 0)
        self.send_data_to_slaves()
        
    def error_ack(self, active_drive_number:int):
        """
        Acknowledges and clears error states by manipulating bits 0 and 7 in the control word.

        Functionality:
            - Sets bit 7 (Error Acknowledge = 1).
            - Clears bit 0 (Switch ON = 0).
            - Sends the updated control word.
            - Clears bit 7 (Error Acknowledge = 0).
            - Sends the updated control word again.
            
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] |= 0x0080 # Set bit 7 (Error Acknowledge = 1)
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0001 # Clear bit 0 (Switch ON = 0)
        self.send_data_to_slaves() # Send Data
        time.sleep(max(self.app.cycle_time * 2, 0.001))
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0080 # Clear bit 7 (Error Acknowledge = 1)
        self.send_data_to_slaves() # Send Data

    # Utility Functions -----------------------------------------------------------------

    def get_unit_scale(self, active_drive_number:int):
        """
        Retrieves the unit scaling factor for the selected Drive.

        Parameters:
            active_drive_number (int): The ID of the active drive.

        Returns:
            float: Modulo factor if rotary motor; otherwise, unit scale.
        """
        with self.app.lm_drive_lock.gen_rlock():
            if self.app.lm_drive_data_dict[active_drive_number].config['is_rotary_motor']:
                return self.app.lm_drive_data_dict[active_drive_number].config['modulo_factor']
            else:
                return self.app.lm_drive_data_dict[active_drive_number].config['unit_scale']
            
    def get_force_scale(self, active_drive_number:int):
        """
        Retrieves the force scaling factor for the selected Drive.

        Parameters:
            active_drive_number (int): The ID of the active drive.

        Returns:
            float: Modulo factor if rotary motor; otherwise, unit scale.
        """
        with self.app.lm_drive_lock.gen_rlock():
            if self.app.lm_drive_data_dict[active_drive_number].config['is_rotary_motor']:
                return self.app.lm_drive_data_dict[active_drive_number].config['fc_torque_scale']
            else:
                return self.app.lm_drive_data_dict[active_drive_number].config['fc_force_scale']
        
        
    def hex_valid(self, value:str|int|float, bit:int=16):
        """
        Converts and validates a hexadecimal value from string or numeric input.

        Parameters:
            value (str|int|float): The value to convert.
            bit (int): Base to use for string conversion (default is 16).

        Returns:
            int|None: Converted integer if valid; otherwise None.
        """
        try:
            if isinstance(value, str):
                return int(value, bit)
            elif isinstance(value, int):
                return value
            elif isinstance(value, float):
                return int(value)
            else:
                return None
        except ValueError:
            logging.warning('Invalid hex string in Control Word')
            return None
        except Exception as e:
            logging.error(f'hex_valid() an error occurred: {e}')
        
    def toggle_bits(self, active_drive_number:int, header:int):
        """
        Increments the 4-bit command counter stored in the lowest 4 bits of a 16-bit header.

        Functionality:
            - Extracts the lower 4 bits of state_var (command count).
            - Increments the count (modulo 16).
            - Updates the header with the new command count.

        Parameters:
            app: The main application object.
            active_drive_number (int): The ID of the active drive.
            header (int): A 16-bit command header whose lowest 4 bits represent the command counter.

        Returns:
            int: Updated 16-bit header with incremented 4-bit command counter (modulo 16).
        """
        with self.app.lm_drive_lock.gen_rlock():
            cmd_count_old = self.app.lm_drive_data_dict[active_drive_number].inputs['state_var'] & 0x000F
        if int(cmd_count_old) == 15:
            cmd_count_old = 0
        cmd_count_new = (cmd_count_old + 1) % 16
        return (header & 0xFFF0) | cmd_count_new

    def toggle_bits_cfg(self, active_drive_number:int, header:int):
        """
        Increments and updates the 4-bit config status counter within the header.

        Functionality:
            - Extracts the lower 4 bits of cfg_status (config status count).
            - Increments the count (modulo 16).
            - Updates the header with the new count.

        Parameters:
            app: The main application object.
            active_drive_number (int): The ID of the active drive.
            header (int): The current 16-bit config header.

        Returns:
            int: Updated header with incremented 4-bit config counter.
        """
        with self.app.lm_drive_lock.gen_rlock():
            cmd_count_old = self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_status'] & 0x000F
        cmd_count_new = (cmd_count_old + 1) % 16
        return (header & 0xFFF0) | cmd_count_new

    def convert23to16(self, value:int):
        """
        Splits a 23-bit integer into two 16-bit words.

        Parameters:
            value (int): A 23-bit integer.

        Returns:
            tuple: (lower 16 bits, upper 7 bits padded as a 16-bit integer)
        """
        value_1 = int(value) & 0xFFFF
        value_2 = int(value) >> 16 & 0xFFFF
        return value_1, value_2
    
    def unsigned_to_signed_16bit(self, value:int):
        """
        Converts an unsigned 16-bit integer to signed format.

        Args:
            value (int): Unsigned 16-bit value.

        Returns:
            int: Signed integer.

        Raises:
            ValueError: If value exceeds 16-bit range.
        """
        if value > 0xFFFF:
            raise ValueError("Value exceeds 16-bit unsigned integer range.")
        return value if value < 0x8000 else value - 0x10000

    def ieee754_bits_to_float(self, value:int, precision:str = "single") -> float:
        """
        Convert an IEEE-754 bit pattern (as an integer) to a Python float.
        - precision="single": expects 32-bit pattern (returns the same numerical value as a 32-bit float).
        - precision="double": expects 64-bit pattern.
        """
        if precision == "single":
            if not (0 <= value < (1 << 32)):
                raise ValueError("For 'single', bits must be a 32-bit unsigned integer (0..2^32-1).")
            b = value.to_bytes(4, byteorder="big", signed=False)
            return struct.unpack("!f", b)[0]
        elif precision == "double":
            if not (0 <= value < (1 << 64)):
                raise ValueError("For 'double', bits must be a 64-bit unsigned integer (0..2^64-1).")
            b = value.to_bytes(8, byteorder="big", signed=False)
            return struct.unpack("!d", b)[0]
        else:
            raise ValueError("precision must be 'single' or 'double'")

    # Drive Communication Functions --------------------------------------------------------------

    def update_output_drive_data(self, active_drive_number:int, controlWord, header, para_word:list, execute_mc:bool=True):
        """
        Updates the motion control output data for a specified drive.

        This function processes input data, validates and updates the control word,
        motion command header, and motion parameters. It ensures the data is correctly
        formatted and sent to the EtherCAT communication queue.

        Args:
            active_drive_number (int): The ID of the drive to update.
            controlWord (str or int): The control word to send to the drive.
            header (str or int): The motion command header.
            para_word (list): A list of motion parameters, each defined as [count, value].
            execute_mc (bool, optional): If True, sends the data immediately to the drive. Defaults to True.

        Returns:
            int: The lower 4 bits of the motion command header (count nibble) used for verification.

        Raises:
            ValueError: If controlWord or header is invalid.
        """
        # Update drive Data
        self.app.ProCommData.process_input_data()
        # control_word
        if controlWord and not controlWord == '0':
            controlWord = self.hex_valid(controlWord)
            if controlWord is None:
                return None
            with self.app.lm_drive_lock.gen_wlock():
                self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] = controlWord
            
        if not header == '' or not header == 0:
            # mc_header
            header = self.hex_valid(header)
            if header is None:
                    return None
            header = self.toggle_bits(active_drive_number, header)
            with self.app.lm_drive_lock.gen_wlock():
                self.app.lm_drive_data_dict[active_drive_number].outputs['mc_header'] = header
            
            # para_word
            bit_count = 0
            for pw in para_word:
                if bit_count <= 10:
                    if pw is not None:
                        if pw[0] == 1:
                            value_1 = int(pw[1])
                        if pw[0] == 2:
                            value_1, value_2 = self.convert23to16(pw[1])
                        with self.app.lm_drive_lock.gen_wlock():
                            for i in range(pw[0]):
                                self.app.lm_drive_data_dict[active_drive_number].outputs[f'mc_para_word{bit_count:02}'] = locals()[f'value_{i+1}']
                                bit_count += 1
                else:
                    self.app.insert_message(f'Something went wrong - there is too much data.')
        if execute_mc:
            self.send_data_to_slaves()
        return header & 0x000F

    def update_output_cfg(self, active_drive_number:int, cfg_control, cfg_index_out, cfg_value_out, execute_mc:bool=True):
        """
        Updates configuration parameters for a specified drive.

        This function validates and updates the configuration control word, index,
        and value. It sends the configuration data to the EtherCAT communication queue
        if execution is requested.

        Args:
            active_drive_number (int): The ID of the drive to configure.
            cfg_control (str or int): The configuration control command.
            cfg_index_out (str or int): The index of the configuration parameter.
            cfg_value_out (str or int, optional): The value to write to the configuration parameter.
            execute_mc (bool, optional): If True, sends the data immediately to the drive. Defaults to True.

        Returns:
            int: The lower 4 bits of the configuration control word (count nibble) used for verification.

        Raises:
            ValueError: If cfg_control, cfg_index_out, or cfg_value_out is invalid.
        """
        # cfg_control
        cfg_control = self.hex_valid(cfg_control)
        if cfg_control is None:
            return None
        cfg_control = self.toggle_bits_cfg(active_drive_number, cfg_control)
        # cfg_index_out
        cfg_index_out = self.hex_valid(cfg_index_out)
        # cfg_value_out
        if cfg_value_out is not None:
            cfg_value_out = self.hex_valid(cfg_value_out, bit=32)
        
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['cfg_control'] = cfg_control
            self.app.lm_drive_data_dict[active_drive_number].outputs['cfg_index_out'] = cfg_index_out
            if cfg_value_out is not None:
                self.app.lm_drive_data_dict[active_drive_number].outputs['cfg_value_out'] = cfg_value_out
        
        # Send data to drive
        if execute_mc:
            self.send_data_to_slaves()
        
        # Return Status Value
        return cfg_control & 0x000F
        
    # Send to Drive ----------------------------------------------------------------------------------
        
    def send_data_to_slaves(self):
        """
        Sends output data from all drives to the EtherCAT communication queue.

        Parameters:
            app: The main application object.

        Returns:
            None
        """
        with self.app.lm_drive_lock.gen_wlock():
            packed_outputs = [self.app.lm_drive_data_dict[device].pack_outputs() for device in range(1, self.app.noDev+1)]
        self.app.ethercat_comm.update_queue.put(packed_outputs)


class LinMot_Information:
    """
    Provides LinMot drive information.
    """
    def __init__(self, app):
        self.app = app
        # Drive Information (Article Number)-> UPID 000Ch
        self.drive_dict = {5589: [1, 'C1250-MI-XC-1S'],
                           5597: [2, 'C1250-MI-XC-0S'],
                           6489: [3, 'F1150-DS-UC-3S'],
                           6767: [4, 'F1050-DS-UC-0S']
                           }


def main():
    print("Do Nothing")
    
if __name__ == "__main__":
    main()