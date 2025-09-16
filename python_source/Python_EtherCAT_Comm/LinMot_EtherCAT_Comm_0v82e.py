"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           EtherCAT_Comm_ .py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  22.05.2025
    Version:        0.82
    Description:    Establishes EtherCAT communication with LinMot Drives
                    and provides basic testing functionalities.

    Disclaimer:
    ------------
    This is a demo project developed for NTI AG | LinMot & MagSpring. This
    software is provided "as-is" without any warranties or guarantees.
    NTI AG does not provide support, updates, or maintenance for this software. 
    Use it at your own risk.

    Dependencies:
    -------------
    Libraries (see import)
    Information for EtherCAT protocoll: LinMot Shop Article Number: 0185-1079

    Description:
    ------
    This Python script manages the EtherCAT communication with LinMot Drives 
    using the pysoem library. It includes a robust multiprocessing-based 
    communication handler class, which sets up the EtherCAT master, configures 
    PDO mappings, and continuously exchanges process data with connected LinMot 
    slave devices. The system is designed exclusively for use with LinMot Drives 
    and supports optional advanced handling through LMDrive_Data structures.
    A key feature of the project is its real-time communication loop that 
    ensures operational state monitoring, fault detection, and safe shutdown 
    capabilities. It includes built-in logging, exception handling, and data 
    sharing mechanisms for inter-process communication.

    Important Note: Closing the Script Properly
    -------------------------------------------
    The script uses continuous communication with the EtherCAT master.
    To ensure proper closure of the program, you must invoke the `stop` function.
    This can be done by implementing a proper cleanup process within the application code.

    If the script is not closed properly, older instances may continue running
    in the background, leading to data conflicts, degraded performance, or
    communication instability with the LinMot drives. Always ensure the script
    terminates cleanly before starting a new instance.
    
    Usage:
    ------
    Instructions to Run <filename>.py
    
    1. Ensure Python 3.12 or later is installed.
    2. Ensure Dependencies are Installed:
        Make sure you have all the necessary Python packages installed.
    3. Importing the Script:
        This script must be imported into other Python projects.
        Refer to Start.py for an example of how to integrate <filename>.py 
        into your own project setup.

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

import pysoem
import time
import multiprocessing as mp
import logging
import datetime
import psutil
import os
import queue
import struct
import ctypes


#----------------------------------------------------------------------------------------------------
# Error Handling
class SlaveCountError(Exception):
    pass

class SdoError(Exception):
    pass


#----------------------------------------------------------------------------------------------------
class EtherCATCommunication:
    """
    EtherCATCommunication handles real-time communication with EtherCAT slaves.

    It sets up the master, manages cyclic communication, provides data synchronization,
    and supports multiprocessing with optional detailed LMDrive-specific data handling.

    Attributes:
        adapter_id (str): The ID of the network adapter used for EtherCAT communication.
        noDev (int): Number of EtherCAT devices expected to connect.
        cycle_time (float): Communication cycle time (in seconds).
        mp_log (int): Logging level for multiprocessing (Info: 20; Error: 40).
        master (pysoem.Master): EtherCAT master instance.
        stop_event (mp.Event): Event used to stop communication gracefully.
        no_Monitoring (int): Number of monitoring input channels (0...4).
        no_Parameter (int): Number of parameter output channels (0...4).
        InputLength (int): Length of each device's input data (calculated automatically).
        data (mp.Array): Shared memory array for input data.
        lock (mp.Lock): Lock to synchronize shared data access.
        data_queue (mp.Queue): Optional queue for recording each communication cycle.
        data_queue_ON (mp.Event): Event flag to control saving of cycle data into the queue.
        slave_name (manager.list): Shared list storing names of connected slaves.
        update_queue (mp.Queue): Queue for sending output updates to slaves (only the latest entry is used).
        error_queue (mp.Queue): Queue for error messages (level 40).
        info_queue (mp.Queue): Queue for informational messages (level 20).
        comm_proc (mp.Process): Background communication process.
        cpu_affinity (list): List of CPU cores to which the process is bound (only active on Linux).
        REALTIME (bool): Flag indicating if real-time scheduling is enabled (only active on Linux).
        realtime_priority (int): Priority level for real-time scheduling (only active when REALTIME enabled).
        MAX_CYCLE_OVERRUN (int): Maximum tolerated consecutive cycle time overruns.
        MAX_SLAVE_COMM_ATTEMPTS (int): Maximum failed state checks before stopping.
        Activate_LMDrive_Data (bool): If True, enables LMDrive_Data management instead of simple raw data.
        lm_drive_data_dict (dict): Shared dictionary holding LMDrive_Data objects per slave
                                   (File: _20250314a_LMDrive_Data_v3)
        
    """
    
    def __init__(self, adapter_id:str, noDev:int, cycle_time:float, lock, no_Monitoring:int=0, no_Parameter:int=0, 
                 Activate_LMDrive_Data:bool=False, mp_logging:int=0, cpu_affinity:list=None, realtime:bool=False, 
                 realtime_priority:int=99):
        """
        Initialize EtherCAT communication.

        Args:
            adapter_id (str): Adapter to open.
            noDev (int): Number of devices expected.
            cycle_time (float): Time between communication cycles.
            lock (mp.Lock): Multiprocessing lock for safe shared access.
            no_Monitoring (int): Number of monitoring input channels (default 0). Valid for all Drives.
            no_Parameter (int): Number of parameter output channels (default 0). Valid for all Drives.
            Activate_LMDrive_Data (bool): If True, use LMDrive_Data structures. Might reduce performance!
            mp_logging (int, optional): Logging level (default 0 = no logging).
            cpu_affinity (list, optional): List of CPU cores to which the process is bound. Default is None. (only active on Linux).
            realtime (bool, optional): Flag indicating if real-time scheduling is enabled. Default is False.
            realtime_priority (int, optional): Priority level for real-time scheduling. Default is 99.

        """
        # Basic attributes
        self.adapter_id = adapter_id
        self.noDev = noDev
        self.cycle_time = cycle_time
        self.mp_log = mp_logging
        self.master = None
        
        # Communication control
        self.stop_event = mp.Event()
        self.stop_event.set() # Start with stop_event set
        self.no_Parameter = no_Parameter
        self.no_Monitoring = no_Monitoring
        self.InputLength = 18 + 8 + (4 * self.no_Monitoring)
        
        # Shared resources
        self.data = mp.Array('i', noDev*self.InputLength) # Raw input data
        self.lock = lock
        
        # Queues
        self.data_queue = mp.Queue() # Queue for data
        self.data_queue_ON = mp.Event() # Disabled by default
        manager = mp.Manager()
        self.slave_name = manager.list([None] * noDev)
        self.update_queue = mp.Queue()
        self.error_queue = mp.Queue()
        self.info_queue = mp.Queue()
        
        # Process reference
        self.comm_proc = None

        # Realtime
        self.REALTIME = realtime
        self.cpu_affinity = cpu_affinity
        self.realtime_priority = realtime_priority
        
        # Constants
        self.MAX_CYCLE_OVERRUN: int = 20
        self.MAX_SLAVE_COMM_ATTEMPTS: int = 10
        
        if os.name == 'nt': # 'nt' is for Windows
            self.BUFFER_TIME = 0.000403
        elif os.name == 'posix': # 'posix' is for Linux and other Unix-like systems
            self.BUFFER_TIME = 0.000053
        else:
            raise ValueError("Unsupported operating system")

        
        # Optional LMDrive data activation
        self.Activate_LMDrive_Data = Activate_LMDrive_Data
        if self.Activate_LMDrive_Data:
            self.lm_drive_data_dict = manager.dict({i+1: LMDrive_Data(no_Monitoring, no_Parameter) for i in range(self.noDev)})
        
    def check_values(self):
        """
        Validate input arguments.
        
        Raises:
            ValueError: If any provided value is out of expected range.
        """
        if not(0 < self.noDev):
            raise ValueError(f"noDev {self.noDev} is out of range! Must be greater than 0.")
        if not(0.0001 <= self.cycle_time <= 1):
            raise ValueError(f"cycle_time {self.cycle_time} is out of range! Must be between 0.00025s and 1s.")
        if not(0 <= self.no_Monitoring <= 4):
            raise ValueError(f"no_Monitoring {self.no_Monitoring} is out of range! Must be between {0} and {4}.")
        if not(0 <= self.no_Parameter <= 4):
            raise ValueError(f"no_Parameter {self.no_Parameter} is out of range! Must be between {0} and {4}.")

    def setup_comm(self):
        """
        Setup EtherCAT master and configure connected slaves.

        Returns:
            list: Configured slave objects or None if setup failed.
        """
        try:
            # Open master
            self.master = pysoem.Master()
            self.master.open(self.adapter_id)
            if self.master.config_init() != self.noDev:
                raise SlaveCountError(f'Expected {self.noDev} devices, but found {self.master.config_init()}')
            
            # Set master into INIT_STATE state
            self.master.state = pysoem.INIT_STATE
            self.master.write_state()
            self.master.state_check(pysoem.INIT_STATE, timeout=50000)
            time.sleep(0.1)

            # Set slaves into PRE-OPERATIONAL state
            self.master.state = pysoem.PREOP_STATE
            self.master.write_state()
            self.master.state_check(pysoem.PREOP_STATE, timeout=50000)
            
            # Configure mappings for each slave
            for i, slave in enumerate(self.master.slaves, start=0):

                # Try to read slave name (from SDO 0x1008:0)
                try:
                    name_bytes = slave.sdo_read(0x1008, 0)
                    self.slave_name[i] = name_bytes.decode('utf-8')
                except (pysoem.SdoError, AttributeError, UnicodeDecodeError) as e:
                    self.slave_name[i] = f"Unnamed_{i}"
                    self.error_queue.put(f'{datetime.datetime.now()} - Slave {i} name could not be read: {e}') if self.mp_log >= 40 else None
                    
                # Configure PDOs (Inputs and Outputs)
                try:
                    # Clear PDOs first
                    slave.sdo_write(0x1C12, 0x00, b'\x00')
                    slave.sdo_write(0x1C13, 0x00, b'\x00')
                    slave.sdo_write(0x1A20, 0x00, b'\x00')
                    slave.sdo_write(0x1620, 0x00, b'\x00')
                    
                    # Configure output PDO mappings
                    slave.sdo_write(0x1C12, 1, (0x1700).to_bytes(2, 'little')) # Default Output
                    slave.sdo_write(0x1C12, 2, (0x1708).to_bytes(2, 'little')) # Config Module Outputs
                    for p in range(self.no_Parameter):
                        slave.sdo_write(0x1C12, 3 + p, (0x1728 + p).to_bytes(2, 'little'))
                    slave.sdo_write(0x1C12, 0, bytes([2 + self.no_Parameter]))
                    
                    # Configure input PDO mappings
                    slave.sdo_write(0x1C13, 1, (0x1B00).to_bytes(2, 'little')) # Default Input
                    slave.sdo_write(0x1C13, 2, (0x1B08).to_bytes(2, 'little')) # Config Module Inputs
                    for m in range(self.no_Monitoring):
                        slave.sdo_write(0x1C13, 3 + m, (0x1B28 + m).to_bytes(2, 'little'))
                    slave.sdo_write(0x1C13, 0, bytes([2 + self.no_Monitoring]))
                    
                except pysoem.pysoem.SdoError as e:
                    raise SdoError(f'{e}\n    ErrorNote: Please try again later. These startup SDO errors are sometimes transient.')
                except Exception as e:
                    raise SdoError(f'Unexpected SDO setup error: {e}')
            
            # Map I/O
            self.master.config_map()
            # Set SAFE-OP
            self.master.state = pysoem.SAFEOP_STATE
            self.master.write_state()
            self.master.state_check(pysoem.SAFEOP_STATE, 50000)
            # Set OPERATIONAL
            self.master.state = pysoem.OP_STATE
            self.master.write_state()
            self.master.state_check(pysoem.OP_STATE, 50000)
            
            return self.master.slaves
            
        except Exception as e:
            if hasattr(self, "master") and self.master is not None:
                try:
                    self.master.close()
                except Exception:
                    pass
            self.error_queue.put(f'{datetime.datetime.now()} - Setup failed: {e}') if self.mp_log >= 40 else None
            return None

    def comm_process(self):
        """
        Main communication process that continuously manages EtherCAT communication cycles.

        It ensures that all slaves stay in the operational state (OP_STATE), handles data exchange 
        with the slaves, processes optional monitoring and parameter channels, and manages 
        overrun protection. In case of communication problems or user interruption, 
        it safely stops communication.
        """
        # Setup the EtherCAT communication
        slaves = self.setup_comm()
        if (slaves is None) or (None in slaves):
            self.stop_event.set()
            self.error_queue.put(f'{datetime.datetime.now()} - Communication could not be established with the slaves/drives. Slaves={slaves}') if self.mp_log >= 40 else None
            return
        
        # If LMDrive Data is activated, fill the drive_type for each drive
        if self.Activate_LMDrive_Data:
            with self.lock:
                for i in range(self.noDev):
                    drive_type = self.slave_name[i]
                    lm_data = self.lm_drive_data_dict[i+1]
                    lm_data.config['drive_type'] = drive_type
                    self.lm_drive_data_dict[i+1] = lm_data
        
        self.info_queue.put(f'Communication setup was successful.') if self.mp_log >= 20 else None
        overrun_count = 0
        self.data_queue_ON.clear() # Default to no oscilloscope recording
        sampleNr = 0
        self.stop_event.clear() # Allow communication loop to run
        lock_timeout = max(self.cycle_time-0.010, 0.004) # Lock acquisition timeout
        start_time = time.perf_counter()
        
        try:
            while not self.stop_event.is_set():
                # Exchange process data with the slaves
                self.master.send_processdata()
                self.master.receive_processdata(2000)

                # Aggregate input data from all slaves
                all_data = b''.join(slave.input for slave in slaves)

                # Optionally send raw data to data_queue (for oscilloscope-style logging)
                if self.data_queue_ON.is_set():
                    try:
                        self.data_queue.put_nowait((sampleNr, all_data))
                    except queue.Full:
                        self.error_queue.put('data_queue is full. Skipping this cycle.') if self.mp_log >= 40 else None
                    sampleNr += 1
                
                # Process Data
                if self.Activate_LMDrive_Data:
                    # Store data into LMDrive_Data objects
                    if self.lock.acquire(timeout=lock_timeout):
                        try:
                            for i in range(self.noDev):
                                device_data = bytes(all_data[i*self.InputLength:(i+1)*self.InputLength])
                                lm_data = self.lm_drive_data_dict[i+1]
                                lm_data.unpack_inputs(device_data)
                                lm_data.update_calculated_fields()
                                self.lm_drive_data_dict[i+1] = lm_data
                                
                                # Prepare outputs to send back
                                packed_data = lm_data.pack_outputs()
                                slaves[i].output = packed_data
                        finally:
                            self.lock.release()
                else:
                    # Store data into shared memory array
                    if self.lock.acquire(timeout=lock_timeout):
                        try:
                            self.data[:] = all_data[:]
                        finally:
                            self.lock.release()

                    # Process latest update_queue data
                    if not self.update_queue.empty():
                        try:
                            while not self.update_queue.empty(): # Empty queue to get the latest value from queue
                                new_rx_data = self.update_queue.get_nowait()
                            if isinstance(new_rx_data, list) and len(new_rx_data) == len(slaves):
                                for i, rx_data_instance in enumerate(new_rx_data):
                                    slaves[i].output = rx_data_instance
                        except Exception as e:
                            self.error_queue.put(f'{datetime.datetime.now()} - An unexpected error occurred while sending data: {e}') if self.mp_log >= 40 else None
                    
                # Timing: keep fixed cycle period
                elapsed_time = time.perf_counter() - start_time
                sleep_time = self.cycle_time - elapsed_time - self.BUFFER_TIME
                if sleep_time > 0:
                    time.sleep(sleep_time)
                    start_time = time.perf_counter()
                    overrun_count = 0
                else:
                    start_time = time.perf_counter()
                    overrun_count += 1
                    self.error_queue.put(f'{datetime.datetime.now()} - Cycle time overrun: '
                                         f'No. {overrun_count} with {(sleep_time*1000):.2f}ms') if self.mp_log >= 40 else None
                    if overrun_count > self.MAX_CYCLE_OVERRUN:
                        raise RuntimeError(f'Cycle time repeatedly overrun ({self.MAX_CYCLE_OVERRUN}), stopping communication.')
        except KeyboardInterrupt:
            self.info_queue.put('Communication interrupted by user.') if self.mp_log >= 20 else None
            self.stop_event.set()
        except Exception as e:
            self.error_queue.put(f'{datetime.datetime.now()} - Unexpected error: {e}') if self.mp_log >= 40 else None
        finally:
            # Ensure safe stopping
            self.stop_event.set()
            self.info_queue.put('Set the master to SAFEOP_STATE and close it.') if self.mp_log >= 20 else None
            try:
                self.master.state = pysoem.SAFEOP_STATE
                self.master.write_state()
                self.master.state_check(pysoem.SAFEOP_STATE, timeout=50000)
                time.sleep(0.1)
                self.master.close()
            except Exception as e:
                self.error_queue.put(f'{datetime.datetime.now()} - Error during closing master: {e}') if self.mp_log >= 40 else None
            self.info_queue.put('Comm function stopped') if self.mp_log >= 20 else None
    
    def start(self):
        """
        Starts the EtherCAT communication by launching a separate process.

        Verifies the configuration values first, then initializes and starts the communication
        subprocess using Python's multiprocessing.
        """
        try:
            # Validate configuration before starting
            self.check_values()

            # Force 'spawn' method for multiprocessing
            mp.set_start_method('spawn', force=True)

            self.comm_proc = mp.Process(target=self.comm_process, daemon=True) # Deamon Thread may not be necessery
            self.comm_proc.start()

            # After starting the process, set CPU affinity and priority
            time.sleep(0.1)  # slight delay to make sure process is started

            if not self.comm_proc.is_alive():
                raise RuntimeError("Communication process failed to start - Multiprocessing not alive.")
            
            if os.name == 'posix': # Check if the OS is Linux (POSIX-compliant)
                p = psutil.Process(self.comm_proc.pid)
                if self.REALTIME:
                    p.cpu_affinity(self.cpu_affinity)  # Bind process to certain CPU cores

                    # Set real-time priority
                    os.system(f"sudo chrt -f {self.realtime_priority} -p {self.comm_proc.pid}")
                
                    logging.info(f"The EtherCAT communication process started with PID {self.comm_proc.pid}, "
                            f"affinity {self.cpu_affinity}, real-time priority {self.realtime_priority}")
                else:
                    if self.cpu_affinity == None:
                        self.cpu_affinity = list(range(psutil.cpu_count()))
                        if 0 in self.cpu_affinity and len(self.cpu_affinity) > 1:
                            self.cpu_affinity.remove(0) # Remove core 0 from the list of available cores
                        p.cpu_affinity(self.cpu_affinity) # Bind process to any core but 0
                    else:
                        p.cpu_affinity(self.cpu_affinity)
                    
                    logging.info(f"The EtherCAT communication process started with PID {self.comm_proc.pid}, "
                            f"affinity {self.cpu_affinity}")
                
            else:
                logging.info(f"The EtherCAT communication process started with PID {self.comm_proc.pid}")
            
        except Exception as e:
            logging.error(f"Failed to start communication process: {e}")
            self.error_queue.put(f"{datetime.datetime.now()} - Failed to start communication process: {e}") if self.mp_log >= 40 else None
            self.stop_event.set()
            self.stop() # Ensure clean shutdown if startup fails
    
    def stop(self):
        """
        Stops the EtherCAT communication process.

        Sends a stop event, waits for the communication process to terminate,
        and clears queues if necessary to unblock the process.
        """
        if self.comm_proc:
            logging.info("Setting stop event.")
            self.stop_event.set()
            self.comm_proc.join(timeout=2)
            
            if self.comm_proc.is_alive():
                logging.warning('The communication process did not terminate in time. Clearing queues.')
                # Empty queues to avoid deadlocks
                for queue_name, queue_obj in [("error_queue", self.error_queue),
                                              ("info_queue", self.info_queue),
                                              ("update_queue", self.update_queue),
                                              ("data_queue", self.data_queue)]:
                    try:
                        if not queue_obj.empty():
                            logging.info(f'Clearing {queue_name} with {queue_obj.qsize()} entries.')
                            while not queue_obj.empty():
                                queue_obj.get_nowait()
                    except Exception as e:
                        logging.warning(f'Error clearing {queue_name}: {e}')

                self.comm_proc.join()
                
            if self.comm_proc.is_alive():
                logging.error('The communication process is still alive —> forcing termination.')
                self.comm_proc.terminate()
                self.comm_proc.join()

            if not self.comm_proc.is_alive():
                logging.info('The EtherCAT communication process stopped successfully.')
            else:
                logging.error('Communication process still alive — forcefull termination failed!')
            self.comm_proc = None


class LMDrive_Data:
    """
    Class representing LinMot drive data including communication I/O,
    motor configuration parameters, and real-time scaled drive status.

    This class handles unpacking and packing of binary data structures,
    maintains configuration settings, and calculates readable status
    values from raw input data.
    """
    def __init__(self, num_mon_channels, num_par_channels):
        """
        Initialize the LMDrive_Data instance.

        Args:
            num_mon_channels (int): Number of monitoring (input) channels.
            num_par_channels (int): Number of parameter (output) channels.
        """
        self.num_mon_ch = num_mon_channels
        self.num_par_ch = num_par_channels
        
        # Configuration parameters for the drive/motor
        self.config = {
            'is_rotary_motor': False,
            'pos_scale_numerator': 10000.0,
            'pos_scale_denominator': 1.0,
            'unit_scale': 10000.0,
            'modulo_factor': 360000,
            'fc_force_scale': 0.1,
            'fc_torque_scale': 0.00057295779513082,
            'drive_name': "LMDrive",
            'drive_type': "0" #"Undefined"
        }
        
        # Status values calculated from inputs or representing drive state
        self.status = {
            'operation_enabled': False,
            'switch_on_locked': False,
            'homed': False,
            'motion_active': False,
            'jogging': False,
            'warning': False,
            'error': False,
            'error_code': 0x00,
            'demand_position': 0.0,
            'actual_position': 0.0,
            'difference_position': 0.0,
            'actual_current': 0.0,
            'nr_of_revolutions': 0
        }
        
        # Output/control data sent to the drive
        self.outputs = {
            'control_word': 0x003E,
            'mc_header': 0x0000,
            'mc_para_word00': 0x0000,
            'mc_para_word01': 0x0000,
            'mc_para_word02': 0x0000,
            'mc_para_word03': 0x0000,
            'mc_para_word04': 0x0000,
            'mc_para_word05': 0x0000,
            'mc_para_word06': 0x0000,
            'mc_para_word07': 0x0000,
            'mc_para_word08': 0x0000,
            'mc_para_word09': 0x0000,
            'cfg_control': 0x0000,
            'cfg_index_out': 0x0000,
            'cfg_value_out': 0x00000000,
        }
        # Add dynamic parameter channels to outputs
        for i in range(1, self.num_par_ch + 1):
            self.outputs[f'par_ch{i}'] = 0x0000
        
        # Input/status data received from the drive
        self.inputs = {
            'state_var': 0x0000,
            'status_word': 0x0000,
            'warn_word': 0x0000,
            'demand_pos': 0x00000000,
            'actual_pos': 0x00000000,
            'demand_curr': 0x0000,
            'cfg_status': 0x0000,
            'cfg_index_in': 0x0000,
            'cfg_value_in': 0x00000000,
        }
        # Add dynamic monitoring channels to inputs
        for i in range(1, self.num_mon_ch + 1):
            self.inputs[f'mon_ch{i}'] = 0x0000
        
        
    def update_calculated_fields(self):
        """
        Updates internal status values using current input data and config.
        Converts raw binary values into readable, scaled physical units.
        """
        # Update `unit_scale` in config
        self.config['unit_scale'] = self.config['pos_scale_numerator'] / self.config['pos_scale_denominator']

        # Update status fields based on inputs
        self.status['operation_enabled'] = bool(self.inputs['status_word'] & 0x0001)  # Bit 0
        self.status['switch_on_locked'] = bool(self.inputs['status_word'] & 0x0040)   # Bit 6
        self.status['homed'] = bool(self.inputs['status_word'] & 0x0800)             # Bit 11
        self.status['motion_active'] = bool(self.inputs['status_word'] & 0x2000)     # Bit 13
        self.status['warning'] = bool(self.inputs['status_word'] & 0x0080)           # Bit 7
        self.status['error'] = bool(self.inputs['status_word'] & 0x0008)             # Bit 3

        # Check error state and set error code
        if self.inputs['state_var'] & 0xFF00 == 0x0400:  # Error state
            self.status['error_code'] = self.inputs['state_var'] & 0x00FF
        else:
            self.status['error_code'] = 0x00

        # Calculate scaled positions and current
        self.status['demand_position'] = ctypes.c_int32(self.inputs['demand_pos']).value / self.config['unit_scale']
        self.status['actual_position'] = ctypes.c_int32(self.inputs['actual_pos']).value / self.config['unit_scale']
        self.status['difference_position'] = round(self.status['demand_position'] - self.status['actual_position'], 4)
        self.status['actual_current'] = ctypes.c_int16(self.inputs['demand_curr']).value / 1000.0
        
    def unpack_inputs(self, data):
        """
        Unpacks binary input data into readable values.

        Args:
            data (bytes): Binary input data stream from the drive.
        """
        base_format = '<HHHiiiHHi'  # Format for fixed fields
        mon_channel_format = 'i' * self.num_mon_ch  # Format for dynamic monitoring channels
        full_format = base_format + mon_channel_format  # Combine formats
        
        unpacked_data = struct.unpack(full_format, data)
        
        (
            self.inputs['state_var'],
            self.inputs['status_word'],
            self.inputs['warn_word'],
            self.inputs['demand_pos'],
            self.inputs['actual_pos'],
            self.inputs['demand_curr'],
            self.inputs['cfg_status'],
            self.inputs['cfg_index_in'],
            self.inputs['cfg_value_in'],
            *mon_channels
        ) = unpacked_data

        # Assign monitoring channels dynamically
        for i, value in enumerate(mon_channels, start=1):
            self.inputs[f'mon_ch{i}'] = value

    def unpack_outputs(self, data):
        """
        Unpacks binary output data into the outputs dictionary.

        Args:
            data (bytes): Binary output data stream to be interpreted.
        """
        base_format_par = '<HHHHHHHHHHHHHHi'  # Format for fixed fields
        par_channel_format = 'i' * self.num_par_ch  # Format for dynamic monitoring channels
        full_format_par = base_format_par + par_channel_format  # Combine formats
        
        unpacked_par_data = struct.unpack(full_format_par, data)
        (
            self.outputs['control_word'],
            self.outputs['mc_header'],
            self.outputs['mc_para_word00'],
            self.outputs['mc_para_word01'],
            self.outputs['mc_para_word02'],
            self.outputs['mc_para_word03'],
            self.outputs['mc_para_word04'],
            self.outputs['mc_para_word05'],
            self.outputs['mc_para_word06'],
            self.outputs['mc_para_word07'],
            self.outputs['mc_para_word08'],
            self.outputs['mc_para_word09'],
            self.outputs['cfg_control'],
            self.outputs['cfg_index_out'],
            self.outputs['cfg_value_out'],
            *par_channels
        ) = unpacked_par_data

        # Assign monitoring channels dynamically
        for i, value in enumerate(par_channels, start=1):
            self.outputs[f'par_ch{i}'] = value
    
    def pack_outputs(self):
        """
        Packs the current `outputs` dictionary into a binary structure.

        Returns:
            bytes: Binary representation of all output fields.
        """
        # Define the fixed structure for outputs
        base_format = '<HHHHHHHHHHHHHHi'
        par_channel_format = 'H' * self.num_par_ch  # Dynamically add parameter channels
        full_format = base_format + par_channel_format

        # Prepare data for packing
        data_to_pack = [
            self.outputs['control_word'],
            self.outputs['mc_header'],
            self.outputs['mc_para_word00'],
            self.outputs['mc_para_word01'],
            self.outputs['mc_para_word02'],
            self.outputs['mc_para_word03'],
            self.outputs['mc_para_word04'],
            self.outputs['mc_para_word05'],
            self.outputs['mc_para_word06'],
            self.outputs['mc_para_word07'],
            self.outputs['mc_para_word08'],
            self.outputs['mc_para_word09'],
            self.outputs['cfg_control'],
            self.outputs['cfg_index_out'],
            self.outputs['cfg_value_out'],
        ]

        # Add parameter channels dynamically
        for i in range(1, self.num_par_ch + 1):
            data_to_pack.append(self.outputs[f'par_ch{i}'])

        # Pack the data
        return struct.pack(full_format, *data_to_pack)

    def __str__(self):
        """
        Returns a human-readable summary of key status fields.

        Returns:
            str: Formatted status string.
        """
        status_str = (
            f"Operation_Enabled: {self.status['operation_enabled']}, "
            f"SwitchOn_Locked: {self.status['switch_on_locked']}, "
            f"Homed: {self.status['homed']}, "
            f"Motion_Active: {self.status['motion_active']}, "
            f"Jogging: {self.status['jogging']}, "
            f"Warning: {self.status['warning']}, "
            f"Error: {self.status['error']}, "
            f"Error_Code: {self.status['error_code']}, "
            f"Demand_Position: {self.status['demand_position']}, "
            f"Actual_Position: {self.status['actual_position']}, "
            f"Difference_Position: {self.status['difference_position']}, "
            f"Actual_Current: {self.status['actual_current']}"
        )
        mon_channels_str = ""
        for i in range(1, self.num_mon_ch + 1):
            key = f"mon_ch{i}"
            mon_channels_str += f", MonCh{i}: {self.inputs.get(key)}"

        return status_str + mon_channels_str


    def __getstate__(self):
        """
        Returns the internal state for pickling.
        
        Returns:
            dict: Dictionary representing the instance state.
        """
        return self.__dict__.copy()  # Shallow copy of the instance's dictionary

    def __setstate__(self, state):
        """
        Restores state from a pickled dictionary.

        Args:
            state (dict): State dictionary to restore.
        """
        self.__dict__.update(state)  # Restore state


if __name__ == "__main__":
    # Multiprocessing start method
    mp.set_start_method('spawn', force=True)
    logging.basicConfig(format='%(levelname)s:%(message)s', level=logging.DEBUG)
    print('Do Nothing')

