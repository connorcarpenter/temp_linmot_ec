"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           <filename>.py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  12.06.2025
    Version:        0.10
    Description:    Setup Tab for GUI

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
    Functions for LinMot_Start_GUI_

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

import customtkinter as ctk
import tkinter as tk
import pysoem
import time


class Find_EC_Master:
    """
    Description:
    ------------
    This script assists in identifying and selecting the appropriate EtherCAT 
    master adapter.
    """

    def __init__(self):
        self = self

    def on_select(self, result_container):
        if selected_adapter.get():
            selected_name = selected_adapter.get()
            selected_desc = adapter_dict[selected_name]
            result_container['result'] = (selected_name, selected_desc)
            #root.quit()
            root.destroy()

    def on_cancel(self, result_container):
        result_container['result'] = None
        #root.quit()
        root.destroy()

    def create_window(self, adapters, result_container, parent=None):
        global root, selected_adapter, adapter_dict

        if parent:
            root = tk.Toplevel(parent)
        else:
            root = tk.Tk()

        root.title("EtherCAT Master")
        selected_adapter = tk.StringVar(value=None)
        adapter_dict = {adapter.name: adapter.desc.decode('utf-8') for adapter in adapters}

        frame = tk.Frame(root, bd=2, relief="solid")
        frame.pack(padx=10, pady=10, fill="both", expand=True)

        tk.Label(frame, text="Select your EtherCAT Master").pack(pady=5)

        for adapter in adapters:
            adapter_name = adapter.name
            adapter_desc = adapter.desc.decode('utf-8')
            tk.Radiobutton(frame, text=adapter_desc, variable=selected_adapter,
                        value=adapter_name, anchor="w").pack(fill="x", padx=5, pady=2)

        button_frame = tk.Frame(root)
        button_frame.pack(pady=10)

        tk.Button(button_frame, text="Select", command=lambda: self.on_select(result_container)).pack(side="left", padx=5)
        tk.Button(button_frame, text="Cancel", command=lambda: self.on_cancel(result_container)).pack(side="right", padx=5)


    def main(self, parent=None):
        adapters = self.adapter_list()
        if not adapters:
            print("No adapters found.")
            return None

        result_container = {'result': None}

        if tk is not None:
            try:
                self.create_window(adapters, result_container, parent)
                if parent:
                    root.wait_window()
                else:
                    root.mainloop()
            except Exception as e:
                print(f"GUI fallback to CLI due to error: {e}")
                return self.create_cli(adapters)
        else:
            return self.create_cli(adapters)

        return result_container['result']

    def adapter_list(self):
        return pysoem.find_adapters()


    def create_cli(self, adapters):
        print("GUI is not available. Please select an EtherCAT Master from the list below:\n")
        for index, adapter in enumerate(adapters):
            print(f"{index}: {adapter.desc.decode('utf-8')} ({adapter.name})")

        try:
            selection = int(input("\nEnter the number of the adapter you'd like to use: "))
            if 0 <= selection < len(adapters):
                adapter = adapters[selection]
                return (adapter.name, adapter.desc.decode('utf-8'))
            else:
                print("Invalid selection.")
        except (ValueError, IndexError):
            print("Invalid input.")
        return None
    

class Processing_comm_data:
    """
    Description:
    ------------
    This module handles the processing of EtherCAT communication data within 
    the LinMot GUI application. It extracts and updates input data from 
    connected drives in a thread-safe manner, ensuring that the application's 
    internal data structures remain synchronized with real-time drive status. 
    The script plays a key role in maintaining accurate monitoring and control 
    across all connected devices.
    """

    def __init__(self, app):
        self.app = app

    def process_input_data(self, data_length):
        """
        Description:
        Processes EtherCAT communication data for connected drives and updates the application's internal data structures.
        
        Parameters:
            app: The main application instance containing shared data.
            data_length (int): Data block size for each device.

        Functionality:
            Thread-Safe Data Retrieval: Locks access and copies EtherCAT data.
            Device Data Processing: Iterates through devices, extracts their data, unpacks inputs, and updates calculated fields.

        Purpose:
        Keeps drive data updated for real-time monitoring and control.
        """
        # Read Data from Drive
        with self.app.lock:
            all_slave_data = self.app.ec_comm_process.data[:]
        
        # Unpack and update the LMDrive input data per device
        for i in range(self.app.noDev):
            device_data = bytes(all_slave_data[i*data_length:(i+1)*data_length])
            with self.app.lm_drive_lock.gen_wlock():
                self.app.lm_drive_data_dict[i+1].unpack_inputs(device_data)
                self.app.lm_drive_data_dict[i+1].update_calculated_fields()


class Send_Data:
    """
    Description:
    ------------
    This script provides utility functions for motor control and data 
    communication with LinMot drives.
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

    def swichOFF_motor(self, active_drive_number:int):
        """
        Turns the motor OFF by clearing bit 0 (Switch ON) in the control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0001 # Clear bit 0 (Switch ON = 0)
        self.send_data_to_slaves()

    def home_motor(self, active_drive_number:int):
        """
        Starts the homing procedure by setting bit 11 (Home) in the control word.
        
        Parameters:
            active_drive_number (int): The ID of the active drive.
        """
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] |= 0x0800 # Set bit 11 (Home = 1)
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
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] |= 0x0080 # Set bit 7 (Error Acknoledge = 1)
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0001 # Clear bit 0 (Switch ON = 0)
        self.send_data_to_slaves() # Send Data
        time.sleep(max(self.app.cycle_time * 2, 0.001))
        with self.app.lm_drive_lock.gen_wlock():
            self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] &= ~0x0080 # Clear bit 7 (Error Acknoledge = 1)
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
        
        
    def hex_valid(self, value, bit:int=16):
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
            print('Invalid hex string in Control Word')
            return None
        
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

    # Drive Communication Functions --------------------------------------------------------------

    def update_output_drive_data(self, active_drive_number:int, controlWord, header, para_word, execute_mc:bool=True):
        """.
        Updates the drive's motion control (MC) output data.

        Functionality:
            - Processes input data using process_input_data().
            - Validates and updates the control_word if provided.
            - Validates and updates the mc_header with toggled command bits.
            - Processes para_word:
                - If values are valid, assigns them to corresponding mc_para_word fields.
                - Ensures the number of parameters does not exceed the limit.
            - Sends the updated data to the slaves.

        Parameters:
            app: The main application object.
            active_drive_number (int): The ID of the active drive.
            controlWord (str): Hexadecimal control word string.
            header (str): Hexadecimal header string for MC communication.
            para_word (list): List of tuples with format (count, value), for up to 11 parameters.
            execute_mc (bool): If True, sends data immediately to the drive.

        Returns:
            None or str: Returns None if successful, or an error message if input is invalid.
        """
        # Update drive Data
        self.app.pro_comm_data.process_input_data(data_length = self.app.ec_comm_process.InputLength)
        # control_word
        if controlWord and not controlWord == '0':
            controlWord = self.hex_valid(controlWord)
            if controlWord == None:
                return None
            with self.app.lm_drive_lock.gen_wlock():
                self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] = controlWord
            
        if not header == '' or not header == 0:
            # mc_header
            header = self.hex_valid(header)
            if header == None:
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
                            value_1 = pw[1]
                        if pw[0] == 2:
                            value_1, value_2 = self.convert23to16(pw[1])
                        with self.app.lm_drive_lock.gen_wlock():
                            for i in range(pw[0]):
                                self.app.lm_drive_data_dict[active_drive_number].outputs[f'mc_para_word{bit_count:02}'] = locals()[f'value_{i+1}']
                                bit_count += 1
                else:
                    self.app.insert_message(f'Someting went wrong - there is too much data.')
        if execute_mc:
            self.send_data_to_slaves()

    def update_output_cfg(self, active_drive_number:int, cfg_control, cfg_index_out, cfg_value_out, execute_mc:bool=True):
        """
        Updates configuration parameters for the specified drive.

        Functionality:
            - Validates and toggles cfg_control bits.
            - Converts and updates cfg_index_out.
            - Converts and updates cfg_value_out (if provided).
            - Sends the updated configuration to the drive.

        Parameters:
            app: The main application object.
            active_drive_number (int): The ID of the active drive.
            cfg_control (str|int): Hexadecimal string or integer for cfg_control.
            cfg_index_out (str|int): Hexadecimal index value for config output.
            cfg_value_out (str|int|None): Optional hexadecimal value for config output.
            execute_mc (bool): Whether to send the data immediately.

        Returns:
            None or str: Returns None if successful, or an error message if input is invalid.
        """
        # cfg_control
        cfg_control = self.hex_valid(cfg_control)
        if cfg_control == None:
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
        self.app.ec_comm_process.update_queue.put(packed_outputs)


class OtherMessages:
    """
    Description:
    ------------
    This module provides a custom popup message interface for the LinMot 
    EtherCAT GUI application. It is primarily used to display critical 
    communication-related alerts to the user, such as EtherCAT connection 
    errors. The script leverages customtkinter to create modal dialogs that 
    prompt the user for confirmation or cancellation, ensuring safe handling 
    of communication shutdowns and improving user awareness during fault 
    conditions.
    """

    def __init__(self, app):
        self.app = app

    def msg_CommError(self):
        """
        Displays a modal popup window to alert the user of an EtherCAT communication error.

        This function creates a customtkinter-based dialog that prompts the user to either
        terminate the communication properly or cancel the action. It ensures that the
        main application window is disabled while the popup is active and returns the
        user's decision.

        Args:
            app (CTk): The main application instance that owns the popup.

        Returns:
            bool: True if the user chooses to stop communication, False if canceled.
        """
        self.app.wait_window = None  # To store the popup window
        self.app.popup_result = None  # Store user choice

        def on_ok():
            self.app.popup_result = True
            self.app.wait_window.destroy()

        def on_cancel():
            self.app.popup_result = False
            self.app.wait_window.destroy()

        self.app.wait_window = ctk.CTkToplevel(self.app)  # Create popup window
        self.app.wait_window.title("EtherCAT Communication Error")
        self.app.wait_window.geometry("350x135")
        self.app.wait_window.grab_set()  # Make it modal (disable main window)

        # Error message
        message = ("An error has occurred with the EtherCAT communication.\n"
                "Would you like to terminate the communication properly?\n\n"
                "An improper stop may require the master to be restarted.")

        label = ctk.CTkLabel(self.app.wait_window, text=message, font=("Arial", 12), wraplength=380, justify="left")
        label.pack(pady=15, padx=20)

        button_frame = ctk.CTkFrame(self.app.wait_window)
        button_frame.pack(pady=10)

        ok_button = ctk.CTkButton(button_frame, text="Stop Communication", command=on_ok)
        ok_button.pack(side="left", padx=10)

        cancel_button = ctk.CTkButton(button_frame, text="Cancel", command=on_cancel)
        cancel_button.pack(side="right", padx=10)

        self.app.wait_window.wait_window()  # Wait for user response
        return self.app.popup_result



def main():
    print("do nothing")
    
if __name__ == "__main__":
    main()