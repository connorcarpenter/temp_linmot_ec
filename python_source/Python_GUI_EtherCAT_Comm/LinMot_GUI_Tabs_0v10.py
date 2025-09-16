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
    Setup and functions for each GUI Tab

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
from tkinter import filedialog
import threading
import re
import time
import configparser
import matplotlib.pyplot as plt
from matplotlib.backends.backend_tkagg import FigureCanvasTkAgg
import csv
import queue
import struct
import ctypes
import LinMot_DataHandling_0v10 as lmdh




class Tab_SetupModule:
    """
    Description:
    ------------
    This module defines the EC Master Setup tab for the LinMot EtherCAT GUI 
    application. It provides the user interface and logic for configuring 
    EtherCAT communication parameters, including adapter selection, number 
    of connected drives, cycle time, and channel settings. The tab also 
    supports saving, loading, and resetting configuration profiles, and 
    dynamically generates drive-specific setup panels. It plays a central 
    role in initializing the communication environment and preparing the 
    system for real-time interaction with LinMot drives.
    """

    def __init__(self, app):
        self.app = app

    def setup_tab(self):
        """
        Initializes the "EC Master Setup" tab in the GUI.

        This includes:
        - General settings for adapter ID, number of devices, and cycle time.
        - Advanced settings for monitoring and parameter channels.
        - Drive-specific configuration tabs.
        - Buttons for starting/stopping communication and saving/loading/resetting configurations.

        Args:
            app (CTk): The main application instance.
        """
        #-----------EC Settings-----------
        # Left side widgets
        left_frame = ctk.CTkFrame(self.app.tabview.tab("EC Master Setup"))
        left_frame.grid(row=0, column=0, padx=0, pady=0, sticky="nw")
        
        # Left side top widgets
        left_frame_top = ctk.CTkFrame(left_frame)
        left_frame_top.grid(row=0, column=0, padx=20, pady=20, sticky="nw")
        
        # Title
        ctk.CTkLabel(left_frame_top, text="General Settings", font=ctk.CTkFont(size=15, weight="bold")).grid(
            row=0, column=0, pady=(5, 10), sticky="w")
        
        # Adapter ID
        adapter_id_button = ctk.CTkButton(left_frame_top, text="Configure Adapter ID", command=lambda: self.configure_adapter_id())
        adapter_id_button.grid(row=1, column=0, padx=0, pady=5)
        
        self.app.adapter_id_label = ctk.CTkLabel(left_frame_top, text="Adapter: Not Selected")
        self.app.adapter_id_label.grid(row=1, column=1, padx=5, pady=5)
        self.app.adapter_id_label.configure(wraplength=150)
        
        # Number of devices input (positive integers only)
        num_devices_label = ctk.CTkLabel(left_frame_top, text="Number of devices")
        num_devices_label.grid(row=2, column=0, padx=0, pady=5, sticky="w")
        
        self.app.num_devices_entry = ctk.CTkEntry(left_frame_top)
        self.app.num_devices_entry.insert(0, "1")
        self.app.num_devices_entry.grid(row=2, column=1, padx=5, pady=5)
        self.app.num_devices_entry.bind("<KeyRelease>", lambda event: self.update_device_settings_tab())
        
        # Cycle time input (floating point numbers only between 0.0005 and 0.1)
        cycle_time_label = ctk.CTkLabel(left_frame_top, text="Cycle time [ms]")
        cycle_time_label.grid(row=3, column=0, padx=0, pady=5, sticky="w")
        
        self.app.cycle_time_entry = ctk.CTkEntry(left_frame_top)
        self.app.cycle_time_entry.grid(row=3, column=1, padx=5, pady=5)
        self.app.cycle_time_entry.bind("<FocusOut>", lambda event: self.validate_float_input(event))
        
        
        # Left side bottom widgets
        left_bottom_frame = ctk.CTkFrame(left_frame)
        left_bottom_frame.grid(row=1, column=0, padx=20, pady=10, sticky="nw")
        
        # Title
        logo_label_l2 = ctk.CTkLabel(left_bottom_frame, text="Advanced Settings", font=ctk.CTkFont(size=15, weight="bold"))
        logo_label_l2.grid(row=0, column=0, padx=0, pady=(5, 0), sticky="w")
        
        logo_label_l2 = ctk.CTkLabel(left_bottom_frame, text="(Valid for all Drives)", font=ctk.CTkFont(size=12, weight="normal"))
        logo_label_l2.grid(row=1, column=0, padx=0, pady=(0, 5), sticky="w")
        
        def validate_mpCh_input(input_text):
            if input_text == "":  # Allow empty input for clearing
                return True
            if input_text.isdigit() and 0 <= int(input_text) <= 4:  # Check if input is a number between 0 and 4
                return True
            return False
        
        # Number of monitoring channels (0...4 integers only)
        num_monit_ch = ctk.CTkLabel(left_bottom_frame, text="Number of Monitoring Channels")
        num_monit_ch.grid(row=2, column=0, padx=0, pady=1, sticky="w")
        
        self.app.num_monit_ch_entry = ctk.CTkEntry(left_bottom_frame, validate="key", validatecommand=(self.app.register(validate_mpCh_input), "%P"))
        self.app.num_monit_ch_entry.insert(0, 0)
        self.app.num_monit_ch_entry.grid(row=2, column=1, padx=5, pady=1)
        
        # Number of parameter channels (0...4 integers only)
        num_para_ch = ctk.CTkLabel(left_bottom_frame, text="Number of Parameter Channels")
        num_para_ch.grid(row=3, column=0, padx=0, pady=1, sticky="w")
        
        self.app.num_para_ch_entry = ctk.CTkEntry(left_bottom_frame, validate="key", validatecommand=(self.app.register(validate_mpCh_input), "%P"))
        self.app.num_para_ch_entry.insert(0, 0)
        self.app.num_para_ch_entry.grid(row=3, column=1, padx=5, pady=1)
        
        # Left side bottom widgets
        left_bottom_b_frame = ctk.CTkFrame(left_frame)
        left_bottom_b_frame.grid(row=3, column=0, padx=20, pady=10, sticky="s")
        
        # Start Communication Button
        self.app.start_comm_button = ctk.CTkButton(left_bottom_b_frame, text="Start Communication", anchor='w', command=self.app.start_communication)
        self.app.start_comm_button.grid(row=0, column=0, padx=0, pady=0, sticky="w")
        
        # Stop Communication Button
        self.app.stop_comm_button = ctk.CTkButton(left_bottom_b_frame, text="Stop Communication", anchor='e', command=self.app.stop_communication)
        self.app.stop_comm_button.grid(row=0, column=1, padx=20, pady=0, sticky="e")
        
        #-----------Drive Settings-----------
        # Right side widgets
        self.app.device_tab_view = ctk.CTkTabview(self.app.tabview.tab("EC Master Setup"))
        self.app.device_tab_view.grid(row=0, column=1, padx=20, pady=20, sticky="nsew")
        self.update_device_settings_tab()
        
        #-----------Config-----------
        # Main botom widgets
        config_frame = ctk.CTkFrame(self.app.tabview.tab("EC Master Setup"))
        config_frame.grid(row=1, column=1, padx=0, pady=0, sticky="nw")
        
        # Save Config
        self.app.save_config_button = ctk.CTkButton(config_frame, text="Save Config", anchor='ew', command=lambda: self.save_config())
        self.app.save_config_button.grid(row=0, column=0, padx=0, pady=0, sticky="w")
        
        # Load Config
        self.app.load_config_button = ctk.CTkButton(config_frame, text="Load Config", anchor='ew', command=lambda: self.load_config())
        self.app.load_config_button.grid(row=0, column=1, padx=5, pady=0, sticky="w")
        
        # Reset Config
        self.app.reset_config_button = ctk.CTkButton(config_frame, text="Reset Config", anchor='ew', command=lambda: self.reset_config())
        self.app.reset_config_button.grid(row=0, column=3, padx=5, pady=0, sticky="w")
        
    def configure_adapter_id(self):
        """
        Launches the adapter selection tool to detect and configure the EtherCAT adapter.

        Updates the application with the selected adapter ID and description.

        Args:
            app (CTk): The main application instance.
        """
        try:
            self.Find_Master = lmdh.Find_EC_Master()
            result = self.Find_Master.main(parent=self.app)
            if result:
                self.app.adapter_id, self.app.adapter_desc = result
                print(f"Received value: {self.app.adapter_id} ({self.app.adapter_desc})")
                self.app.adapter_id_label.configure(text=f"Adapter: {self.app.adapter_desc}")
                self.app.insert_message(text=f'Adapter ID configured: {self.app.adapter_id} ({self.app.adapter_desc})')
            else:
                self.app.adapter_id_label.configure(text="No adapter selected.")
        except Exception as e:
            self.app.adapter_id_label.configure(text="Error selecting adapter.")
            self.app.insert_message(text=f"Error: {str(e)}")
            print(f"Exception occurred: {e}")


        
    def validate_integer_input(self, event=None):
        """
        Validates that the input for number of devices contains only digits.

        Args:
            app (CTk): The main application instance.
            event (tk.Event, optional): The triggering event. Defaults to None.
        """
        widget = event.widget
        value = widget.get()
        if not value.isdigit() and value != "":
            widget.delete(0, ctk.END)
            widget.insert(0, ''.join(filter(str.isdigit, value)))

    def validate_float_input(self, event=None):
        """
        Validates and clamps the cycle time input to a valid float between 0.5 and 500 ms.

        If the input is invalid, it resets to a default value.

        Args:
            app (CTk): The main application instance.
            event (tk.Event, optional): The triggering event. Defaults to None.
        """
        value = self.app.cycle_time_entry.get()  # Get the current input
        try:
            # Attempt to convert the value to a float
            float_value = float(value)
            if 0.5 <= float_value <= 500:
                # Valid input, do nothing
                return
            else:
                # Out of range: Clamp the value and give feedback
                clamped_value = max(0.5, min(float_value, 500))
                self.app.cycle_time_entry.delete(0, ctk.END)
                self.app.cycle_time_entry.insert(0, f"{clamped_value:.1f}")
                self.app.insert_message(text= f"Cycle time clamped to {clamped_value:.1f} (valid range: 0.5 - 500)")
        except ValueError:
            # Invalid input (not a float): Reset to default and give feedback
            self.app.cycle_time_entry.delete(0, ctk.END)
            self.app.cycle_time_entry.insert(0, "20.0")
            self.app.insert_message(text= 'Invalid input for cycle time. Reset to default (20.0).')
        
        
    def update_device_settings_tab(self, unlock=False):
        """
        Updates the drive configuration tabs based on the number of devices entered.

        Adds or removes tabs dynamically. Optionally clears and resets tabs if `unlock` is True.

        Args:
            app (CTk): The main application instance.
            unlock (bool, optional): Whether to reset all tabs before updating. Defaults to False.
        """
        try:
            num_devices = int(self.app.num_devices_entry.get())
        except ValueError:
            num_devices = 1
            
        # Renew all Tabs first if start_communication has been pressed.
        if unlock:
            tab_names = self.app.device_tab_view._segmented_button._value_list  # Access the internal list of tab names
            for _ in range(min(num_devices, len(tab_names))):
                self.app.device_tab_view.delete(tab_names[0])
            
        tab_names = self.app.device_tab_view._segmented_button._value_list  # Access the internal list of tab names
        if len(tab_names) > num_devices:
            #Remove Tabs
            for i in range(num_devices, len(tab_names)):
                self.app.device_tab_view.delete(tab_names[num_devices])
        else:
            #Add new Tabs
            for i in range(len(tab_names)+1, num_devices + 1):
                tab_name_device = f"Drive {i}"
                if tab_name_device not in tab_names:
                    self.drive_info_tab(tab_name_device)
        
        
    def drive_info_tab(self, tab_name):
        """
        Creates a new tab for a drive with editable fields for name, motor type,
        and monitoring/parameter channels.

        Args:
            app (CTk): The main application instance.
            tab_name (str): The name of the drive tab (e.g., "Drive 1").
        """
        tab_number = int(re.search(r'\d+', tab_name).group())
        self.app.device_tab_view.add(tab_name)
        ctk.CTkLabel(self.app.device_tab_view.tab(tab_name), text=f'Optional settings for {tab_name}').pack(padx=5, pady=5)
        
        # Frame for Device Name
        device_name_frame = ctk.CTkFrame(self.app.device_tab_view.tab(tab_name))
        device_name_frame.pack(padx=10, pady=2, anchor="w")
        
        # Device Name
        Device_Name_label = ctk.CTkLabel(device_name_frame, text='Drive Name')
        Device_Name_label.pack(side="left", padx=(10, 2), pady=5)
        entry = ctk.CTkEntry(device_name_frame, state="disabled", fg_color="gray")
        entry.insert(0, f'{tab_name}')
        entry.pack(side="left", padx=2, pady=5)
        self.app.Device_Name_entry[f'D{tab_number}_DeviceName'] = entry

        # Frame for Device Name
        motor_type_frame = ctk.CTkFrame(self.app.device_tab_view.tab(tab_name))
        motor_type_frame.pack(padx=10, pady=2, anchor="w")
        
        # Motor Type
        Device_Name_label = ctk.CTkLabel(motor_type_frame, text='Motor Type')
        Device_Name_label.pack(side="left", padx=(10, 2), pady=5)
        entry = ctk.CTkSegmentedButton(motor_type_frame, values=["Linear", "Rotary"], state="disabled")
        entry.set("Linear")
        entry.pack(side="left", padx=2, pady=5)
        self.app.seg_button_1[f'D{tab_number}_SegBut'] = entry
        
        # Create a parent frame for Monitoring and Parameter Channels
        parent_frame = ctk.CTkFrame(self.app.device_tab_view.tab(tab_name))
        parent_frame.pack(padx=10, pady=10, fill="both", expand=True)
        
        # Create frames for Monitoring and Parameter Channels side by side
        monitoring_frame = ctk.CTkFrame(parent_frame)
        monitoring_frame.grid(row=0, column=0, padx=10, pady=5, sticky="nsew")
        
        parameter_frame = ctk.CTkFrame(parent_frame)
        parameter_frame.grid(row=0, column=1, padx=10, pady=5, sticky="nsew")
        
        # Label for Monitoring Channels
        ctk.CTkLabel(monitoring_frame, text="Monitoring Channels", font=("Arial", 12, "bold")).pack(padx=5, pady=5)
        
        # Label for Parameter Channels
        ctk.CTkLabel(parameter_frame, text="Parameter Channels", font=("Arial", 12, "bold")).pack(padx=5, pady=5)
        
        # Create Monitoring Channels dynamically
        num_monitoring_channels = 4  # Number of Monitoring Channels
        
        for i in range(1, num_monitoring_channels + 1):
            mch_frame = ctk.CTkFrame(monitoring_frame)
            mch_frame.pack(padx=1, pady=2, anchor="w")
            
            mch_label = ctk.CTkLabel(mch_frame, text=f'M Ch {i}')
            mch_label.grid(row=0, column=0, padx=(0, 2), pady=5)
            
            upid_entry = ctk.CTkEntry(mch_frame, placeholder_text="UPID", width=70, state="disabled", fg_color="gray")
            upid_entry.grid(row=0, column=1, padx=2, pady=5)
            
            name_entry = ctk.CTkEntry(mch_frame, placeholder_text="Name", state="disabled", fg_color="gray")
            name_entry.grid(row=0, column=2, padx=2, pady=5)
            
            # Store entries in the dictionary
            self.app.monitoring_entries[f'D{tab_number}_MCh{i}_UPID'] = upid_entry
            self.app.monitoring_entries[f'D{tab_number}_MCh{i}_Name'] = name_entry
        
        # Create Parameter Channels dynamically
        num_parameter_channels = 4  # Number of Parameter Channels
        
        for i in range(1, num_parameter_channels + 1):
            pch_frame = ctk.CTkFrame(parameter_frame)
            pch_frame.pack(padx=1, pady=2, anchor="w")
            
            pch_label = ctk.CTkLabel(pch_frame, text=f'P Ch {i}')
            pch_label.grid(row=0, column=0, padx=(0, 2), pady=5)
            
            upid_entry = ctk.CTkEntry(pch_frame, placeholder_text="UPID", width=70, state="disabled", fg_color="gray")
            upid_entry.grid(row=0, column=1, padx=2, pady=5)
            
            name_entry = ctk.CTkEntry(pch_frame, placeholder_text="Name", state="disabled", fg_color="gray")
            name_entry.grid(row=0, column=2, padx=2, pady=5)
            
            # Store entries in the dictionary
            self.app.parameter_entries[f'D{tab_number}_PCh{i}_UPID'] = upid_entry
            self.app.parameter_entries[f'D{tab_number}_PCh{i}_Name'] = name_entry
            
    def unlock_drive_info(self, lock=False):
        """
        Enables or disables editing of drive configuration fields.

        Args:
            app (CTk): The main application instance.
            lock (bool, optional): If True, disables editing. If False, enables it. Defaults to False.
        """
        for tab_name in self.app.device_tab_view._segmented_button._value_list:
            tab_number = int(re.search(r'\d+', tab_name).group())
            
            # Unlock Device Name Entry
            if not lock:
                self.app.Device_Name_entry[f'D{tab_number}_DeviceName'].configure(state="normal", fg_color="white")
            else:
                self.app.Device_Name_entry[f'D{tab_number}_DeviceName'].configure(state="disabled", fg_color="gray")
            
            # Unlock Motor Type Selection
            if not lock:
                self.app.seg_button_1[f'D{tab_number}_SegBut'].configure(state="normal")
            else:
                self.app.seg_button_1[f'D{tab_number}_SegBut'].configure(state="disabled")
            
            # Unlock Monitoring Channels
            if not lock:
                for i in range(1, int(self.app.num_monit_ch_entry.get())+1):
                    self.app.monitoring_entries[f"D{tab_number}_MCh{i}_UPID"].configure(state="normal", fg_color="white")
                    self.app.monitoring_entries[f'D{tab_number}_MCh{i}_Name'].configure(state="normal", fg_color="white")
            else:
                for i in range(1, int(self.app.num_monit_ch_entry.get()) + 1):
                    self.app.monitoring_entries[f"D{tab_number}_MCh{i}_UPID"].configure(state="disabled", fg_color="gray")
                    self.app.monitoring_entries[f'D{tab_number}_MCh{i}_Name'].configure(state="disabled", fg_color="gray")
            
            # Unlock Parameter Channels
            if not lock:
                for i in range(1, int(self.app.num_para_ch_entry.get())+1):
                    self.app.parameter_entries[f'D{tab_number}_PCh{i}_UPID'].configure(state="normal", fg_color="white")
                    self.app.parameter_entries[f'D{tab_number}_PCh{i}_Name'].configure(state="normal", fg_color="white")
            else:
                for i in range(1, int(self.app.num_para_ch_entry.get())+1):
                    self.app.parameter_entries[f'D{tab_number}_PCh{i}_UPID'].configure(state="disabled", fg_color="gray")
                    self.app.parameter_entries[f'D{tab_number}_PCh{i}_Name'].configure(state="disabled", fg_color="gray")
            

    def save_config(self):
        """
        Saves the current communication and drive configuration to an INI file.

        Args:
            app (CTk): The main application instance.
        """
        config = configparser.ConfigParser()
        config_file_name = filedialog.asksaveasfilename(defaultextension=".ini", filetypes=[("INI files", "*.ini"), ("All files", "*.*")])
        
        if not config_file_name:
            return
        
        config['CommunicationSettings'] = {
            'Adapter_ID': self.app.adapter_id,
            'Adapter_Desc': self.app.adapter_desc,
            'Number_of_Devices': self.app.num_devices_entry.get(),
            'Cycle_Time': self.app.cycle_time_entry.get(),
            'No_Monitoring_Ch': self.app.num_monit_ch_entry.get(),
            'No_Parameter_Ch': self.app.num_para_ch_entry.get()
        }
        
        for drive_name in self.app.device_tab_view._segmented_button._value_list:
            i = int(re.search(r'\d+', drive_name).group())
            config[f'Drive.{i}'] = {
                'Drive_Tab': drive_name,
                'Drive_Name': self.app.Device_Name_entry[f'D{i}_DeviceName'].get(),
                'Motor_Type': self.app.seg_button_1[f'D{i}_SegBut'].get(),
                f'D{i}_MCh{1}_UPID': self.app.monitoring_entries[f'D{i}_MCh1_UPID'].get(),
                f'D{i}_MCh{1}_Name': self.app.monitoring_entries[f'D{i}_MCh1_Name'].get(),
                f'D{i}_MCh{2}_UPID': self.app.monitoring_entries[f'D{i}_MCh2_UPID'].get(),
                f'D{i}_MCh{2}_Name': self.app.monitoring_entries[f'D{i}_MCh2_Name'].get(),
                f'D{i}_MCh{3}_UPID': self.app.monitoring_entries[f'D{i}_MCh3_UPID'].get(),
                f'D{i}_MCh{3}_Name': self.app.monitoring_entries[f'D{i}_MCh3_Name'].get(),
                f'D{i}_MCh{4}_UPID': self.app.monitoring_entries[f'D{i}_MCh4_UPID'].get(),
                f'D{i}_MCh{4}_Name': self.app.monitoring_entries[f'D{i}_MCh4_Name'].get(),
                f'D{i}_PCh{1}_UPID': self.app.parameter_entries[f'D{i}_PCh1_UPID'].get(),
                f'D{i}_PCh{1}_Name': self.app.parameter_entries[f'D{i}_PCh1_Name'].get(),
                f'D{i}_PCh{2}_UPID': self.app.parameter_entries[f'D{i}_PCh2_UPID'].get(),
                f'D{i}_PCh{2}_Name': self.app.parameter_entries[f'D{i}_PCh2_Name'].get(),
                f'D{i}_PCh{3}_UPID': self.app.parameter_entries[f'D{i}_PCh3_UPID'].get(),
                f'D{i}_PCh{3}_Name': self.app.parameter_entries[f'D{i}_PCh3_Name'].get(),
                f'D{i}_PCh{4}_UPID': self.app.parameter_entries[f'D{i}_PCh4_UPID'].get(),
                f'D{i}_PCh{4}_Name': self.app.parameter_entries[f'D{i}_PCh4_Name'].get()
            }
        
        
        with open(config_file_name, 'w', encoding='utf-8') as configfile:
            config.write(configfile)
        
        self.app.insert_message(f"Data saved to {config_file_name}")

    def load_config(self):
        """
        Loads communication and drive configuration from an INI file.

        Updates the GUI fields and unlocks drive settings accordingly.

        Args:
            app (CTk): The main application instance.
        """

        config = configparser.ConfigParser()
        config_file_name = filedialog.askopenfilename(filetypes=[("INI files", "*.ini"), ("All files", "*.*")])
        
        if not config_file_name:
            return
        
        config.read(config_file_name, encoding='utf-8')
        
        if 'CommunicationSettings' in config:
            self.app.adapter_id = config['CommunicationSettings'].get('Adapter_ID', '')
            self.app.adapter_desc = config['CommunicationSettings'].get('Adapter_Desc', '')
            self.app.adapter_id_label.configure(text=f"Adapter: {self.app.adapter_desc}")
            
            self.app.num_devices_entry.delete(0, 'end')
            self.app.num_devices_entry.insert(0, config['CommunicationSettings'].get('Number_of_Devices', ''))
            
            self.update_device_settings_tab() # Unlock number of Devices
            
            self.app.cycle_time_entry.delete(0, 'end')
            self.app.cycle_time_entry.insert(0, config['CommunicationSettings'].get('Cycle_Time', ''))
            
            self.app.num_monit_ch_entry.delete(0, 'end')
            self.app.num_monit_ch_entry.insert(0, config['CommunicationSettings'].get('No_Monitoring_Ch', ''))
            
            self.app.num_para_ch_entry.delete(0, 'end')
            self.app.num_para_ch_entry.insert(0, config['CommunicationSettings'].get('No_Parameter_Ch', ''))
            
            self.unlock_drive_info() # Unlock Monitoring and Parameter CH
            
        for section in config.sections():
            if section.startswith("Drive."):
                i = int(section.split(".")[1])
                
                self.app.Device_Name_entry[f'D{i}_DeviceName'].delete(0, 'end')
                self.app.Device_Name_entry[f'D{i}_DeviceName'].insert(0, config[section].get('Drive_Name', ''))
                
                self.app.seg_button_1[f'D{i}_SegBut'].set(config[section].get('Motor_Type', ''))
                
                for ch in range(1, 5):
                    self.app.monitoring_entries[f'D{i}_MCh{ch}_UPID'].delete(0, 'end')
                    self.app.monitoring_entries[f'D{i}_MCh{ch}_UPID'].insert(0, config[section].get(f'D{i}_MCh{ch}_UPID', ''))
                    
                    self.app.monitoring_entries[f'D{i}_MCh{ch}_Name'].delete(0, 'end')
                    self.app.monitoring_entries[f'D{i}_MCh{ch}_Name'].insert(0, config[section].get(f'D{i}_MCh{ch}_Name', ''))
                    
                    self.app.parameter_entries[f'D{i}_PCh{ch}_UPID'].delete(0, 'end')
                    self.app.parameter_entries[f'D{i}_PCh{ch}_UPID'].insert(0, config[section].get(f'D{i}_PCh{ch}_UPID', ''))
                    
                    self.app.parameter_entries[f'D{i}_PCh{ch}_Name'].delete(0, 'end')
                    self.app.parameter_entries[f'D{i}_PCh{ch}_Name'].insert(0, config[section].get(f'D{i}_PCh{ch}_Name', ''))
        
        self.app.insert_message(f"Configuration loaded from {config_file_name}")
            
    def reset_config(self):
        """
        Resets all configuration fields to their default values.

        Clears adapter ID, cycle time, monitoring/parameter channels, and drive-specific settings.

        Args:
            app (CTk): The main application instance.
        """
        self.app.adapter_id = ''
        self.app.adapter_desc = ''
        self.app.adapter_id_label.configure(text="Adapter: Not Selected")
        
        self.app.num_devices_entry.delete(0, 'end')
        self.app.num_devices_entry.insert(0, '1')
        
        self.update_device_settings_tab(unlock=False)
        
        self.app.cycle_time_entry.delete(0, 'end')
        self.app.cycle_time_entry.insert(0, '')
        
        self.app.num_monit_ch_entry.delete(0, 'end')
        self.app.num_monit_ch_entry.insert(0, '0')
        
        self.app.num_para_ch_entry.delete(0, 'end')
        self.app.num_para_ch_entry.insert(0, '0')
        
        for drive_name in self.app.device_tab_view._segmented_button._value_list:
            i = int(re.search(r'\d+', drive_name).group())
            
            self.app.Device_Name_entry[f'D{i}_DeviceName'].delete(0, 'end')
            self.app.Device_Name_entry[f'D{i}_DeviceName'].insert(0, '')
            
            self.app.seg_button_1[f'D{i}_SegBut'].set('')
            
            for ch in range(1, 5):
                self.app.monitoring_entries[f'D{i}_MCh{ch}_UPID'].delete(0, 'end')
                self.app.monitoring_entries[f'D{i}_MCh{ch}_UPID'].insert(0, '')
                
                self.app.monitoring_entries[f'D{i}_MCh{ch}_Name'].delete(0, 'end')
                self.app.monitoring_entries[f'D{i}_MCh{ch}_Name'].insert(0, '')
                
                self.app.parameter_entries[f'D{i}_PCh{ch}_UPID'].delete(0, 'end')
                self.app.parameter_entries[f'D{i}_PCh{ch}_UPID'].insert(0, '')
                
                self.app.parameter_entries[f'D{i}_PCh{ch}_Name'].delete(0, 'end')
                self.app.parameter_entries[f'D{i}_PCh{ch}_Name'].insert(0, '')
        
        self.unlock_drive_info(lock=True)


class Tab_Status:
    """
    Description:
    ------------
    This module defines the Control Status tab for the LinMot EtherCAT GUI 
    application. It provides a real-time interface for monitoring and 
    interacting with connected LinMot drives. The tab displays key drive status 
    indicators such as operational state, homing status, motion activity, 
    warnings, and errors. It also includes controls for switching drives on or 
    off, initiating homing sequences, and acknowledging errors. The module 
    dynamically generates UI components for each drive and integrates with 
    the EtherCAT communication backend to reflect live data updates.
    """

    def __init__(self, app):
        self.app = app

    def monitor_status(self, active_tab):
        """
        Updates the status display and control logic for the currently selected drive tab.

        This function reads the latest drive data from the EtherCAT communication backend
        and updates the GUI fields accordingly. It also handles user interactions such as
        switching the motor on/off, initiating homing, and acknowledging errors.

        Args:
            app (CTk): The main application instance containing shared state and UI elements.
            active_tab (str): The name of the currently active drive tab.
        """
            
        for drive_name in self.app.device_tab_view._segmented_button._value_list:
            active_drive_number = int(re.search(r'\d+', drive_name).group())
            
            if active_tab == drive_name:
                no_Monitoring = int(self.app.num_monit_ch_entry.get())
                # 'Operation Enabled' = 'opeed_Entry'
                self.app.status_entries[active_drive_number]['opeed_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['operation_enabled'])
                # 'Switch On Locked' = 'swied_Entry'
                self.app.status_entries[active_drive_number]['swied_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['switch_on_locked'])
                # 'Homed' = 'homed_Entry'
                self.app.status_entries[active_drive_number]['homed_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['homed'])
                # 'Motion Active' = 'motve_Entry'
                self.app.status_entries[active_drive_number]['motve_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['motion_active'])
                # 'Jogging' = 'jogng_Entry'
                self.app.status_entries[active_drive_number]['jogng_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['jogging'])
                # 'Warning' = 'warng_Entry'
                self.app.status_entries[active_drive_number]['warng_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['warning'])
                # 'Error' = 'error_Entry'
                self.app.status_entries[active_drive_number]['error_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['error'])
                # 'Error Code' = 'errde_Entry'
                self.app.status_entries[active_drive_number]['errde_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].status['error_code'])
                # 'State Variable' = 'stale_Entry'
                self.app.status_entries[active_drive_number]['stale_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].inputs['state_var']:04X}h")
                # 'Status Word' = 'stard_Entry'
                self.app.status_entries[active_drive_number]['stard_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].inputs['status_word']:04X}h")
                # 'Warn Word' = 'warrd_Entry'
                self.app.status_entries[active_drive_number]['warrd_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].inputs['warn_word']:04X}h")
                # 'Demand Position' = 'demon_Entry'
                self.app.status_entries[active_drive_number]['demon_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].status['demand_position']} mm")
                # 'Actual Position' = 'acton_Entry'
                self.app.status_entries[active_drive_number]['acton_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].status['actual_position']} mm")
                # 'Diference Position' = 'difon_Entry'
                self.app.status_entries[active_drive_number]['difon_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].status['difference_position']} mm")
                # 'Nr of Revolutions' = 'nrns_Entry'
                #self.app.status_entries[active_drive_number]['nrns_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].status['nr_of_revolutions']} mm") # Does not do anything!
                # 'Actual Current' = 'actnt_Entry'
                self.app.status_entries[active_drive_number]['actnt_Entry'].cget("textvariable").set(f"{self.app.lm_drive_data_dict[active_drive_number].status['actual_current']} mm")
                # 'Monitoring Chanel 1' = 'mon1_Entry'
                if no_Monitoring >= 1:
                    self.app.status_entries[active_drive_number]['mon1_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['mon_ch1'])
                # 'Monitoring Chanel 2' = 'mon2_Entry' 
                if no_Monitoring >= 2:
                    self.app.status_entries[active_drive_number]['mon2_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['mon_ch2'])
                # 'Monitoring Chanel 3' = 'mon3_Entry' 
                if no_Monitoring >= 3:
                    self.app.status_entries[active_drive_number]['mon3_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['mon_ch3'])
                # 'Monitoring Chanel 4' = 'mon4_Entry'
                if no_Monitoring >= 4:
                    self.app.status_entries[active_drive_number]['mon4_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['mon_ch4'])
                # 'CFG Status' = 'cfgus_Entry'
                self.app.status_entries[active_drive_number]['cfgus_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_status'])
                # 'CFG Index in' = 'cfgin_Entry'
                self.app.status_entries[active_drive_number]['cfgin_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_index_in'])
                # 'CFG Value in' = 'cfginV_Entry'
                self.app.status_entries[active_drive_number]['cfginV_Entry'].cget("textvariable").set(self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_value_in'])
                
                # Swich on Motor
                motor_started = self.app.lm_drive_data_dict[active_drive_number].status['operation_enabled']
                if self.app.switch_on[active_drive_number].get() and not motor_started:
                    self.app.sendData.switchON_motor(active_drive_number)
                if not self.app.switch_on[active_drive_number].get() and motor_started:
                    self.app.sendData.swichOFF_motor(active_drive_number)
                    
                # Home
                if self.app.home_switch[active_drive_number].get():
                    homing_started = (self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] & 0x0800) != 0
                    if not homing_started:
                        self.app.sendData.home_motor(active_drive_number)
                        
                if self.app.error_ack[active_drive_number].get():
                    self.app.sendData.error_ack(active_drive_number)
                    self.app.switch_on[active_drive_number].deselect()
                    self.app.error_ack[active_drive_number].deselect()
                    
            # End Homing sequence
            homing_started = (self.app.lm_drive_data_dict[active_drive_number].outputs['control_word'] & 0x0800) != 0
            if homing_started and self.app.lm_drive_data_dict[active_drive_number].status['homed']:
                self.app.sendData.end_home_motor(active_drive_number)
                self.app.home_switch[active_drive_number].deselect()
                
        

    def drive_status_tab(self, drive_names):
        """
        Initializes the "Control Status" tab with a tabbed interface for each connected drive.

        For each drive, this function creates a dedicated tab containing:
        - Control switches (e.g., switch on, homing, error acknowledge)
        - Status indicators (e.g., operation enabled, motion active, error code)
        - Monitoring and configuration fields

        Args:
            app (CTk): The main application instance containing shared state and UI elements.
            drive_names (list of str): A list of drive tab names (e.g., ["Drive 1", "Drive 2"]).
        """

        control_status_tab = self.app.tabview.tab("Control Status")
        self.app.drive_tabview = ctk.CTkTabview(control_status_tab)
        self.app.drive_tabview.grid(row=0, column=0, padx=20, pady=10, sticky="nsew")
        
        self.app.switch_on = {}
        self.app.home_switch = {}
        self.app.error_ack = {}

        for drive_name in drive_names:
            drive_number = int(re.search(r'\d+', drive_name).group()) #active_drive_number
            self.app.drive_tabview.add(drive_name)
            drive_contol_panel_tab = self.app.drive_tabview.tab(drive_name)

            # Create Control Panel frame inside the tab
            drive_contol_panel_frame = ctk.CTkFrame(drive_contol_panel_tab)
            drive_contol_panel_frame.grid(row=0, column=0, padx=20, pady=10, sticky="nsew")
            
            #-----------Control-----------
            ctk.CTkLabel(drive_contol_panel_frame, text="Control", font=ctk.CTkFont(size=12, weight="bold")).grid(
                row=0, column=0, padx=10, pady=(5, 5), sticky="w")
            # Create Control frame inside the Control Panel frame
            drive_contol_frame = ctk.CTkFrame(drive_contol_panel_frame)
            drive_contol_frame.grid(row=1, column=0, padx=10, pady=5, sticky="nsew")
            drive_contol_frame.grid_columnconfigure(1, minsize=20)
            
            # Switch on
            switch_on_label = ctk.CTkLabel(drive_contol_frame, text="Switch on")
            switch_on_label.grid(row=1, column=1, padx=0, pady=5, sticky="ew")
            switch_on_off_label = ctk.CTkLabel(drive_contol_frame, text="OFF")
            switch_on_off_label.grid(row=2, column=0, padx=0, pady=5, sticky="e")
            entry = ctk.CTkSwitch(drive_contol_frame, text="", onvalue=True, offvalue=False)
            entry.grid(row=2, column=1, padx=5, pady=5, sticky="we")
            self.app.switch_on[drive_number] = entry
            switch_on_on_label = ctk.CTkLabel(drive_contol_frame, text="ON")
            switch_on_on_label.grid(row=2, column=2, padx=0, pady=5, sticky="w")
            
            # Arrow pointing down
            arrow_label = ctk.CTkLabel(drive_contol_frame, text="↓")
            arrow_label.grid(row=3, column=1, padx=0, pady=5)
            
            # Home switch
            home_switch_label = ctk.CTkLabel(drive_contol_frame, text="Homing")
            home_switch_label.grid(row=4, column=1, padx=0, pady=5, sticky="ew")
            home_switch_off_label = ctk.CTkLabel(drive_contol_frame, text="OFF")
            home_switch_off_label.grid(row=5, column=0, padx=0, pady=5, sticky="e")
            entry = ctk.CTkSwitch(drive_contol_frame, text="", onvalue=True, offvalue=False)
            entry.grid(row=5, column=1, padx=5, pady=5, sticky="ew")
            self.app.home_switch[drive_number] = entry
            home_switch_on_label = ctk.CTkLabel(drive_contol_frame, text="ON")
            home_switch_on_label.grid(row=5, column=2, padx=0, pady=5, sticky="w")
            
            # Arrow pointing down
            arrow_label = ctk.CTkLabel(drive_contol_frame, text=" ")
            arrow_label.grid(row=6, column=1, padx=0, pady=5)
            
            # Error Acknowledge
            error_ack_label = ctk.CTkLabel(drive_contol_frame, text="Error Acknoledge")
            error_ack_label.grid(row=7, column=1, padx=0, pady=5, sticky="ew")
            entry = ctk.CTkSwitch(drive_contol_frame, text="", onvalue=True, offvalue=False)
            entry.grid(row=8, column=1, padx=0, pady=5)
            self.app.error_ack[drive_number] = entry
            
            #-----------Status 1-----------
            ctk.CTkLabel(drive_contol_panel_frame, text="Status", font=ctk.CTkFont(size=12, weight="bold")).grid(
                row=0, column=1, padx=10, pady=(5, 5), sticky="w")
            # Create Status frame inside the Control Panel frame
            drive_status_frame1 = ctk.CTkFrame(drive_contol_panel_frame)
            drive_status_frame1.grid(row=1, column=1, padx=10, pady=5, sticky="nsew")
            
            self.app.status_values1 = {
                'Operation Enabled': ctk.IntVar(value=True),
                'Switch On Locked': ctk.IntVar(value=False),
                'Homed': ctk.IntVar(value=True),
                'Motion Active': ctk.IntVar(value=False),
                'Jogging': ctk.IntVar(value=True),
                'Warning': ctk.IntVar(value=False),
                'Error': ctk.IntVar(value=True),
                'Error Code': ctk.IntVar(value=0),
                'State Variable': ctk.StringVar(value='0000h'),
                'Status Word': ctk.StringVar(value='0000h'),
                'Warn Word': ctk.StringVar(value='0000h')
            }
            
            # Dictionary to store the entries with their dynamically generated names
            self.app.status_entries[drive_number] = {}

            for i, (label, var) in enumerate(self.app.status_values1.items()):
                ctk.CTkLabel(drive_status_frame1, text=label).grid(row=i, column=0, padx=10, pady=3, sticky="w")
                
                # Generate the dynamic name
                dynamic_name = (label[:3] + label[-2:]).replace(" ", "").lower() + "_Entry"
                
                # Create the entry and assign it a dynamic name
                entry = ctk.CTkEntry(drive_status_frame1, textvariable=var, state="readonly")
                entry.grid(row=i, column=1, padx=5, pady=3, sticky="e")
                
                # Store the entry in the dictionary
                self.app.status_entries[drive_number][dynamic_name] = entry
                
            
            #-----------Status 2-----------
            ctk.CTkLabel(drive_contol_panel_frame, text=" ", font=ctk.CTkFont(size=12, weight="bold")).grid(
                row=0, column=2, padx=10, pady=(5, 5), sticky="w")
            # Create Status frame inside the Control Panel frame
            drive_status_frame1 = ctk.CTkFrame(drive_contol_panel_frame)
            drive_status_frame1.grid(row=1, column=2, padx=10, pady=5, sticky="nsew")
            
            self.app.status_values2 = {
                'Demand Position': ctk.StringVar(value='0.0 mm'),
                'Actual Position': ctk.StringVar(value='0.0 mm'),
                'Difference Position': ctk.StringVar(value='0.0 mm'),
                'Nr of Revolutions': ctk.IntVar(value=0),
                'Actual Current': ctk.StringVar(value='0.0 mm')
            }
            # Add monitoring channels conditionally
            no_Monitoring = int(self.app.num_monit_ch_entry.get())
            if no_Monitoring >= 1:
                self.app.status_values2['Monitoring Chanel 1'] = ctk.DoubleVar(value=0.0)
            if no_Monitoring >= 2:
                self.app.status_values2['Monitoring Chanel 2'] = ctk.DoubleVar(value=0.0)
            if no_Monitoring >= 3:
                self.app.status_values2['Monitoring Chanel 3'] = ctk.DoubleVar(value=0.0)
            if no_Monitoring >= 4:
                self.app.status_values2['Monitoring Chanel 4'] = ctk.DoubleVar(value=0.0)
            self.app.status_values2['CFG Status'] = ctk.IntVar(value=0)
            self.app.status_values2['CFG Index in'] = ctk.IntVar(value=0)
            self.app.status_values2['CFG Value in'] = ctk.IntVar(value=0)

            for i, (label, var) in enumerate(self.app.status_values2.items()):
                ctk.CTkLabel(drive_status_frame1, text=label).grid(row=i, column=0, padx=10, pady=3, sticky="w")
                
                # Generate the dynamic name
                if label == 'CFG Value in':
                    dynamic_name = (label[:3] + label[-2:]).replace(" ", "").lower() + "V_Entry"
                else:
                    dynamic_name = (label[:3] + label[-2:]).replace(" ", "").lower() + "_Entry"
                
                # Create the entry and assign it a dynamic name
                entry = ctk.CTkEntry(drive_status_frame1, textvariable=var, state="readonly")
                entry.grid(row=i, column=1, padx=5, pady=3, sticky="e")
                
                # Store the entry in the dictionary
                self.app.status_entries[drive_number][dynamic_name] = entry


class Tab_SimpleMotion:
    """
    Description:
    ------------
    This module implements the Simple Motion tab within the LinMot EtherCAT 
    GUI application. It provides a user-friendly interface for sending motion 
    commands, configuring drive parameters, and accessing configuration 
    registers for connected LinMot drives. The tab supports both direct motion 
    control and low-level parameter manipulation, integrating seamlessly with 
    the application's communication backend and drive data structures. It is 
    designed to facilitate quick testing, diagnostics, and parameter tuning 
    in a structured and interactive environment.
    """

    def __init__(self, app):
        self.app = app

    def simple_motion_tab(self, drive_names):
        """
        Initializes the Simple Motion tab in the GUI for controlling LinMot drives.

        This tab includes:
        - A segmented button for selecting drives.
        - A motion command panel for setting position, speed, acceleration, etc.
        - A drive command interface for sending control words and parameter words.
        - A parameter access section for reading/writing UPIDs.

        Args:
            app (CTk): The main application instance.
            drive_names (list): List of drive names to populate the segmented button.
        """
        motion_simple_tab = self.app.tabview.tab("Simple Motion")
        # Top row: Horizontal selection of drives
        self.app.drive_selection = ctk.CTkSegmentedButton(motion_simple_tab, values=drive_names, command=lambda event=None: self.delete_cfg_output()) #delete_cfg_output(app)
        self.app.drive_selection.pack(pady=10, padx=10)
        
        # Frames container
        frames_container = ctk.CTkFrame(motion_simple_tab)
        frames_container.pack(fill="both", expand=True, padx=10, pady=10)
        frames_container.grid_columnconfigure(0, minsize=300)
        
        # Frame 1 ----------------------------------
        frame1 = ctk.CTkFrame(frames_container, width=300)
        frame1.grid(row=0, column=0, padx=10, pady=10, sticky="nsew")

        # Configure the grid columns
        frame1.grid_columnconfigure(0, weight=1)
        frame1.grid_columnconfigure(1, weight=1)

        # Mode label and dropdown spanning both columns
        mode_label = ctk.CTkLabel(frame1, text="Motion Command Mode", font=ctk.CTkFont(size=12, weight="bold"))
        mode_label.grid(row=0, column=0, columnspan=2, pady=(10, 0), padx=10, sticky="w")
        self.app.motion_mode_dropdown = ctk.CTkOptionMenu(frame1, values=["Absolute_VAI", "Relative_VAI", "Absolute_VAJI", "Relative_VAJI", "Incr_Act_Pos_RstI", "Absolute_Sin", "Relative_Sin"])
        self.app.motion_mode_dropdown.grid(row=1, column=0, columnspan=2, pady=3, padx=(10, 10), sticky="ew")

        # Position label and entry
        label = ctk.CTkLabel(frame1, text='Position [mm]')
        label.grid(row=2, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.target_pos = ctk.CTkEntry(frame1, placeholder_text="Enter value")
        self.app.target_pos.grid(row=2, column=1, pady=3, padx=(10, 10), sticky="ew")
        # Speed label and entry
        label = ctk.CTkLabel(frame1, text='Speed [m/s]')
        label.grid(row=3, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.max_v = ctk.CTkEntry(frame1, placeholder_text="Enter value")
        self.app.max_v.grid(row=3, column=1, pady=5, padx=(10, 10), sticky="ew")
        # Acceleration label and entry
        label = ctk.CTkLabel(frame1, text='Acceleration [m/s^2]')
        label.grid(row=4, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.acc = ctk.CTkEntry(frame1, placeholder_text="Enter value")
        self.app.acc.grid(row=4, column=1, pady=3, padx=(10, 10), sticky="ew")
        # Deceleration label and entry
        label = ctk.CTkLabel(frame1, text='Deceleration [m/s^2]')
        label.grid(row=5, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.dcc = ctk.CTkEntry(frame1, placeholder_text="Enter value")
        self.app.dcc.grid(row=5, column=1, pady=3, padx=(10, 10), sticky="ew")
        # Jerk label and entry
        label = ctk.CTkLabel(frame1, text='Jerk [m/s^3]')
        label.grid(row=6, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.jerk = ctk.CTkEntry(frame1, placeholder_text="Enter value")
        self.app.jerk.grid(row=6, column=1, pady=3, padx=(10, 10), sticky="ew")

        # Button spanning both columns
        self.app.send_motion = ctk.CTkButton(frame1, text='Send motion command', command=lambda: self.send_motion_command())
        self.app.send_motion.grid(row=7, column=0, columnspan=2, pady=10, padx=(10, 10), sticky="ew")
        
        
        # Frame 2 ------------------------------------
        frame2 = ctk.CTkFrame(frames_container, width=300)
        frame2.grid(row=0, column=1, padx=10, pady=10, sticky="nsew")

        # Configure the grid columns
        frame2.grid_columnconfigure(0, weight=2)  # Column 1
        frame2.grid_columnconfigure(1, weight=1)  # Column 2
        frame2.grid_columnconfigure(2, weight=2)  # Column 3
        
        label = ctk.CTkLabel(frame2, text="Drive Command Interface", font=ctk.CTkFont(size=12, weight="bold"))
        label.grid(row=0, column=0, columnspan=3, pady=(10, 0), padx=10, sticky="w")
        
        # Control Word label and entry
        label = ctk.CTkLabel(frame2, text='Control Word')
        label.grid(row=1, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.control_word = ctk.CTkEntry(frame2, placeholder_text="Hex Input")
        self.app.control_word.grid(row=1, column=2, pady=3, padx=(10, 10), sticky="ew")
        
        # MC Header label and entry
        label = ctk.CTkLabel(frame2, text='MC Header')
        label.grid(row=2, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.mc_header = ctk.CTkEntry(frame2, placeholder_text="Hex Input")
        self.app.mc_header.grid(row=2, column=2, pady=3, padx=(10, 10), sticky="ew")

        # MC Para Word 00 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 00")
        label.grid(row=3, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word00_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word00_dropdown.grid(row=3, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word00_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word00_input.grid(row=3, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 01 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 01")
        label.grid(row=4, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word01_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word01_dropdown.grid(row=4, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word01_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word01_input.grid(row=4, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")
        
        # MC Para Word 02 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 02")
        label.grid(row=5, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word02_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word02_dropdown.grid(row=5, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word02_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word02_input.grid(row=5, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 03 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 03")
        label.grid(row=6, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word03_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word03_dropdown.grid(row=6, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word03_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word03_input.grid(row=6, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 04 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 04")
        label.grid(row=7, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word04_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word04_dropdown.grid(row=7, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word04_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word04_input.grid(row=7, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 05 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 05")
        label.grid(row=8, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word05_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word05_dropdown.grid(row=8, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word05_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word05_input.grid(row=8, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 06 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 06")
        label.grid(row=9, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word06_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word06_dropdown.grid(row=9, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word06_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word06_input.grid(row=9, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 07 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 07")
        label.grid(row=10, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word07_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word07_dropdown.grid(row=10, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word07_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word07_input.grid(row=10, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 08 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 08")
        label.grid(row=11, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word08_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word08_dropdown.grid(row=11, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word08_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word08_input.grid(row=11, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")

        # MC Para Word 09 label, dropdown, and entry
        label = ctk.CTkLabel(frame2, text="MC Para Word 09")
        label.grid(row=12, column=0, pady=(10, 0), padx=2, sticky="w")
        self.app.mc_para_word09_dropdown = ctk.CTkOptionMenu(frame2, values=["Int16", "Int32"])
        self.app.mc_para_word09_dropdown.grid(row=12, column=1, pady=5, padx=(10, 10), sticky="ew")
        self.app.mc_para_word09_input = ctk.CTkEntry(frame2, placeholder_text="Int Input")
        self.app.mc_para_word09_input.grid(row=12, column=2, columnspan=2, pady=5, padx=(10, 10), sticky="ew")
        
        # Description
        label = ctk.CTkLabel(frame2, text="Entries left empty will not be changed / send to the drive")
        label.grid(row=13, column=0, columnspan=3, pady=(10, 0), padx=10, sticky="w")
        
        # Send drive command button spanning both columns
        self.app.send_drive_command = ctk.CTkButton(frame2, text='Send drive command', command=lambda: self.send_drive_command())
        self.app.send_drive_command.grid(row=14, column=0, columnspan=3, pady=10, padx=(10, 10), sticky="ew")
        
        
        # Frame 3 ------------------------------------
        frame3 = ctk.CTkFrame(frames_container, corner_radius=10)
        frame3.grid(row=0, column=2, padx=10, pady=10, sticky="nsew")
        
        label = ctk.CTkLabel(frame3, text="Parameter Access Settings", font=ctk.CTkFont(size=12, weight="bold"))
        label.grid(row=0, column=0, columnspan=3, pady=(10, 0), padx=10, sticky="w")

        # Parameter Access Mode
        label = ctk.CTkLabel(frame3, text='Mode')
        label.grid(row=1, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.parameter_mode_dropdown = ctk.CTkOptionMenu(frame3, values=["Read_Value_ROM", "Read_Value_RAM", "Write_Value_ROM", "Write_Value_RAM", "Write_Value_RAM_and_ROM"])
        self.app.parameter_mode_dropdown.grid(row=1, column=1, pady=3, padx=(10, 10), sticky="ew")
        
        # UPID Out
        label = ctk.CTkLabel(frame3, text='UPID')
        label.grid(row=2, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.upid_out = ctk.CTkEntry(frame3, placeholder_text="Hex Input")
        self.app.upid_out.grid(row=2, column=1, pady=3, padx=(10, 10), sticky="ew")
        
        # UPID Value Out
        label = ctk.CTkLabel(frame3, text='UPID Value')
        label.grid(row=3, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.upid_value_out = ctk.CTkEntry(frame3, placeholder_text="Int Input")
        self.app.upid_value_out.grid(row=3, column=1, pady=3, padx=(10, 10), sticky="ew")
        
        # Send parameter command button spanning both columns
        self.app.send_parameter_command = ctk.CTkButton(frame3, text='Send parameter access command', command=lambda: self.send_parameter_command())
        self.app.send_parameter_command.grid(row=4, column=0, columnspan=2, pady=10, padx=(10, 10), sticky="ew")

        # Recieve Data
        label = ctk.CTkLabel(frame3, text="Recieved Parameter", font=ctk.CTkFont(size=12))
        label.grid(row=5, column=0, columnspan=3, pady=(10, 0), padx=10, sticky="w")
        
        # CFG Ctrl Word
        label = ctk.CTkLabel(frame3, text='CFG Ctrl Word')
        label.grid(row=6, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.sm_cfg_ctrl_word = ctk.CTkEntry(frame3, placeholder_text="Sent value", state='readonly')
        self.app.sm_cfg_ctrl_word.grid(row=6, column=1, pady=3, padx=(10, 10), sticky="ew")
        # CFG Status
        label = ctk.CTkLabel(frame3, text='CFG Status')
        label.grid(row=7, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.sm_cfg_status = ctk.CTkEntry(frame3, placeholder_text="Recieved value after request", state='readonly')
        self.app.sm_cfg_status.grid(row=7, column=1, pady=3, padx=(10, 10), sticky="ew")
        # CFG Index in
        label = ctk.CTkLabel(frame3, text='CFG Index in')
        label.grid(row=8, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.sm_cfg_index_in = ctk.CTkEntry(frame3, placeholder_text="Recieved value after request", state='readonly')
        self.app.sm_cfg_index_in.grid(row=8, column=1, pady=3, padx=(10, 10), sticky="ew")
        # CFG Value in
        label = ctk.CTkLabel(frame3, text='CFG Value in')
        label.grid(row=9, column=0, pady=(10, 0), padx=10, sticky="w")
        self.app.sm_cfg_value_in = ctk.CTkEntry(frame3, placeholder_text="Recieved value after request", state='readonly')
        self.app.sm_cfg_value_in.grid(row=9, column=1, pady=3, padx=(10, 10), sticky="ew")
        
        
        # GUI Apperence
        frames_container.grid_columnconfigure(0, weight=2)
        frames_container.grid_columnconfigure(1, weight=1)
        frames_container.grid_columnconfigure(2, weight=2)
        

    def send_motion_command(self):
        """
        Sends a motion command to the selected drive using user-defined parameters.

        The command includes position, speed, acceleration, deceleration, and jerk.
        It supports various motion types such as Absolute_VAI, Relative_VAI, etc.

        Args:
            app (CTk): The main application instance.
        """
        # Get active Drive
        if self.app.drive_selection.get() == '':
            self.message_select_false()
            return
        active_drive_number = int(re.search(r'\d+', self.app.drive_selection.get()).group())
        # Assign Motion command
        header = self.app.motion_mode_dropdown.get()
        header_map = {
            "Absolute_VAI": (0x0100, False, False),
            "Relative_VAI": (0x0110, False, False),
            "Absolute_VAJI": (0x3A00, False, True),
            "Relative_VAJI": (0x3A10, False, True),
            "Incr_Act_Pos_RstI": (0x0D90, False, False),
            "Absolute_Sin": (0x0E00, True, False),
            "Relative_Sin": (0x0E10, True, False)
        }
        if header not in header_map:
            raise ValueError(f"Unsupported motion header: {header}")
        
        header_code, acc_combined, jerk_necessary = header_map[header]
        unit_scale = self.app.sendData.get_unit_scale(active_drive_number)

        # Issue: No error handling if the user enters invalid (non-numeric) input.
        
        pw = [
            [2, float(self.app.target_pos.get()) * unit_scale],
            [2, float(self.app.max_v.get()) * unit_scale * 100],
            [2, float(self.app.acc.get()) * unit_scale * 10],
            [0], [0]
        ]
        if not acc_combined:
            if self.app.dcc.get() == 0:
                raise TypeError(f"Missing required argument 'dcc' with header '{header}'")
            else:
                pw[3] = ([2, float(self.app.dcc.get()) * unit_scale * 10])
        if jerk_necessary:
            if self.app.jerk.get() == 0:
                raise TypeError(f"Missing required argument 'jerk' with header '{header}'")
            else:
                pw[4] = ([2, float(self.app.jerk.get()) * unit_scale])

        self.app.sendData.update_output_drive_data(active_drive_number, controlWord = 0, header = header_code, para_word=pw, execute_mc=True)
            
    def send_drive_command(self):
        """
        Sends a drive command to the selected drive.

        This includes a control word, motion control header, and up to 10 motion control parameter words.
        The parameter words can be 16-bit or 32-bit integers.

        Args:
            app (CTk): The main application instance.
        """
        if self.app.drive_selection.get() == '':
            a = self.message_select_false()
            return
        active_drive_number = int(re.search(r'\d+', self.app.drive_selection.get()).group())
        controlWord = self.app.control_word.get()
        header = self.app.mc_header.get()
        pw = [None]*10
        for i in range(0, 10):
            attribute_name = f"mc_para_word{i:02}_input"  # Construct the attribute name
            para_input = getattr(self.app, attribute_name)  # Get the attribute
            if not para_input.get() == '':
                attribute_name = f"mc_para_word{i:02}_dropdown"
                para_bit = getattr(self.app, attribute_name)
                if para_bit.get() == 'Int16':
                    bit = 1
                elif para_bit.get() == 'Int32':
                    bit = 2
                else:
                    raise ValueError(f'Bit-size not found for mc_para_word{i:02}')
                pw[i] = [bit, para_input.get()]
        self.app.sendData.update_output_drive_data(active_drive_number, controlWord, header, para_word=pw)
        
    def send_parameter_command(self):
        """
        Sends a parameter access command to the selected drive.

        Supports reading and writing UPID values in RAM or ROM. Updates the GUI with the received response.

        Args:
            app (CTk): The main application instance.
        """
        # Get active Drive
        if self.app.drive_selection.get() == '':
            a = self.message_select_false()
            return
        active_drive_number = int(re.search(r'\d+', self.app.drive_selection.get()).group())
        # Assign value_out
        value_out = False
        match self.app.parameter_mode_dropdown.get():
            case "Read_Value_ROM":
                cfg_control = 0x1000
            case "Read_Value_RAM":
                cfg_control = 0x1100
            case "Write_Value_ROM":
                cfg_control = 0x1200
                value_out = True
            case "Write_Value_RAM":
                cfg_control = 0x1300
                value_out = True
            case "Write_Value_RAM_and_ROM":
                cfg_control = 0x1400
                value_out = True
                
        cfg_index_out = self.app.upid_out.get()
        if value_out:
            try:
                cfg_value_out = int(self.app.upid_value_out.get())
            except ValueError:
                self.app.insert_message("Invalid UPID Value input.")
                return
        else:
            cfg_value_out = None
        self.app.sendData.update_output_cfg(active_drive_number, cfg_control, cfg_index_out, cfg_value_out)
        sm_cfg_str_ctrl_word = f"{self.app.lm_drive_data_dict[active_drive_number].outputs['cfg_control']:04X}h"

        # Update Recieved Parameter
        time.sleep(self.app.cycle_time * 4)
        self.app.pro_comm_data.process_input_data(data_length = self.app.ec_comm_process.InputLength)
        with self.app.lm_drive_lock.gen_rlock():
            sm_cfg_str = {'status': f"{self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_status']:04X}h",
                        'index': f"{self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_index_in']:04X}h",
                        'value': str(self.app.lm_drive_data_dict[active_drive_number].inputs['cfg_value_in'])
                        }
        # sm_cfg_ctrl_word
        self.app.sm_cfg_ctrl_word.configure(state='normal') # Temporarily make it editable
        self.app.sm_cfg_ctrl_word.delete(0, 'end') # Clear the current value
        self.app.sm_cfg_ctrl_word.insert(0, sm_cfg_str_ctrl_word)
        self.app.sm_cfg_ctrl_word.configure(state='readonly') # Set it back to read-only
        # sm_cfg_status
        self.app.sm_cfg_status.configure(state='normal') # Temporarily make it editable
        self.app.sm_cfg_status.delete(0, 'end') # Clear the current value
        self.app.sm_cfg_status.insert(0, sm_cfg_str['status'])
        self.app.sm_cfg_status.configure(state='readonly') # Set it back to read-only
        # sm_cfg_index_in
        self.app.sm_cfg_index_in.configure(state='normal')
        self.app.sm_cfg_index_in.delete(0, 'end')
        self.app.sm_cfg_index_in.insert(0, sm_cfg_str['index'])
        self.app.sm_cfg_index_in.configure(state='readonly')
        # sm_cfg_value_in
        self.app.sm_cfg_value_in.configure(state='normal')
        self.app.sm_cfg_value_in.delete(0, 'end')
        self.app.sm_cfg_value_in.insert(0, sm_cfg_str['value'])
        self.app.sm_cfg_value_in.configure(state='readonly')
            
    def message_select_false(self):
        """
        Displays an error message if no drive is selected or required values are missing.

        Args:
            app (CTk): The main application instance.
        """
        self.app.insert_message(f'Someting went wrong - please make sure that you have selected a drive and entered all required values.')
        
    def delete_cfg_output(self):
        """
        Clears the configuration output fields in the parameter access section.

        This is typically called when a new drive is selected.

        Args:
            app (CTk): The main application instance.
        """
        # sm_cfg_ctrl_word
        self.app.sm_cfg_ctrl_word.configure(state='normal') # Temporarily make it editable
        self.app.sm_cfg_ctrl_word.delete(0, 'end') # Clear the current value
        self.app.sm_cfg_ctrl_word.configure(state='readonly') # Set it back to read-only
        # sm_cfg_status
        self.app.sm_cfg_status.configure(state='normal') # Temporarily make it editable
        self.app.sm_cfg_status.delete(0, 'end') # Clear the current value
        self.app.sm_cfg_status.configure(state='readonly') # Set it back to read-only
        #sm_cfg_index_in
        self.app.sm_cfg_index_in.configure(state='normal')
        self.app.sm_cfg_index_in.delete(0, 'end')
        self.app.sm_cfg_index_in.configure(state='readonly')
        #sm_cfg_value_in
        self.app.sm_cfg_value_in.configure(state='normal')
        self.app.sm_cfg_value_in.delete(0, 'end')
        self.app.sm_cfg_value_in.configure(state='readonly')


class Tab_MotionProfile:
    """
    Description:
    ------------
    This module defines the Motion Profile tab for the LinMot EtherCAT GUI 
    application. It provides a structured interface for executing predefined 
    motion sequences on one or more connected LinMot drives. The tab is designed 
    to demonstrate motion control capabilities using various velocity, 
    acceleration, and position profiles. It integrates with the application's 
    communication backend and drive data structures to ensure synchronized 
    execution and real-time feedback.
    """

    def __init__(self, app):
        self.app = app

    def setup_motion_profile_tab(self):
        """
        Initializes the Motion Profile tab in the GUI.

        This function sets up the layout, description, and control buttons for executing
        predefined motion sequences on one or two LinMot drives.

        Args:
            app (CTk): The main application instance.
        """
        # Get the Motion Profile tab
        motion_profile_tab = self.app.tabview.tab("Motion Profile")
        
        # Configure grid layout
        motion_profile_tab.grid_rowconfigure(0, weight=0)
        motion_profile_tab.grid_rowconfigure(1, weight=1)
        motion_profile_tab.grid_columnconfigure(0, weight=1)
        motion_profile_tab.grid_columnconfigure(1, weight=1)
        
        # Description label
        description = ("These are two motion examples for LinMot Drives. "
                    "The movements are designed for linear motors with a stroke of more than 50 mm. "
                    "The buttons differentiate between the number of motors that will be moved. "
                    "Details on how the motion is coded, can be found in the GUI script.")
        
        description_label = ctk.CTkLabel(motion_profile_tab, text=description, wraplength=600, justify="left")
        description_label.grid(row=0, column=0, columnspan=2, padx=20, pady=10, sticky="w")
        
        # Buttons for motion profiles
        one_motor_button = ctk.CTkButton(motion_profile_tab, text="1 Motor Motion Profile (GUI save)", command=lambda: self.one_motor_motion())
        one_motor_button.grid(row=1, column=0, padx=20, pady=20, sticky="ew")
        
        two_motor_button = ctk.CTkButton(motion_profile_tab, text="2 Motor Motion Profile (GUI blocking)", command=lambda: self.two_motor_motion()) #two_motor_motion(app))
        two_motor_button.grid(row=1, column=1, padx=20, pady=20, sticky="ew")

        # Indicator lights
        self.app.one_motor_indicator = ctk.CTkLabel(motion_profile_tab, text="", width=20, height=20, fg_color="orange")
        self.app.one_motor_indicator.grid(row=2, column=0, pady=10)
        
        self.app.two_motor_indicator = ctk.CTkLabel(motion_profile_tab, text="", width=20, height=20, fg_color="orange")
        self.app.two_motor_indicator.grid(row=2, column=1, pady=10)

    def trigger_motion(self, indicator):
        """
        Temporarily changes the color of a motion indicator to green and resets it after 1 second.

        Args:
            app (CTk): The main application instance.
            indicator (CTkLabel): The indicator label to update.
        """
        indicator.configure(fg_color="green")
        threading.Thread(target=self.reset_indicator, args=(self.app, indicator), daemon=True).start()

    def reset_indicator(self, indicator):
        """
        Resets the motion indicator color to orange after a delay.

        Args:
            app (CTk): The main application instance.
            indicator (CTkLabel): The indicator label to reset.
        """
        time.sleep(1)
        indicator.configure(fg_color="orange")

    def one_motor_motion(self):
        """
        Executes a predefined motion sequence for a single LinMot drive.

        This function performs a series of absolute position moves with varying
        velocity and acceleration profiles. It also manages GUI updates and
        ensures the drive is ready before starting.

        Args:
            app (CTk): The main application instance.
        """
        self.app.insert_message('1 Motor motion sequence started')
        self.app.one_motor_indicator.configure(fg_color="green")
        self.app.specific_update_interval[1] = self.app.cycle_time
        self.app.start_fast_update_thread()

        active_drive_number = 1
        sleep_time_cycle = max(self.app.cycle_time, 0.001)

        def start_motion_sequence(self, sleep_time):
            """Starts the motion sequence if the motor is ready."""
            with self.app.lm_drive_lock.gen_rlock():
                operation_enabled = self.app.lm_drive_data_dict[active_drive_number].status['operation_enabled']
                homed = self.app.lm_drive_data_dict[active_drive_number].status['homed']
                error = self.app.lm_drive_data_dict[active_drive_number].status['error']
                ma = self.app.lm_drive_data_dict[active_drive_number].status['motion_active']

            if operation_enabled and homed and not error and not ma:
                #print('Step 1')
                self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=50, max_v=0.01, acc=0.1, dcc=0.1, jerk=0)
                self.motion_finished_forGUI(sleep_time, active_drive_number, target_pos=50, next_step=step_2)
            else:
                self.app.insert_message('Motor not ready to start motion. Ensure operation is enabled and motor is homed.')
                finalize_motion(self)

        def step_2(self, sleep_time):
            #print('Step 2')
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=0, max_v=0.01, acc=0.1, dcc=0.1, jerk=0)
            self.motion_finished_forGUI(sleep_time, active_drive_number, target_pos=0, next_step=step_3)

        def step_3(self, sleep_time):
            #print('Step 3')
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=25, max_v=1, acc=10, dcc=10, jerk=0)
            self.motion_finished_forGUI(sleep_time, active_drive_number, target_pos=25, next_step=step_4)

        def step_4(self, sleep_time):
            #print('Step 4')
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=50, max_v=0.01, acc=0.1, dcc=0.1, jerk=0)
            self.motion_finished_forGUI(sleep_time, active_drive_number, target_pos=50, next_step=step_5)

        def step_5(self, sleep_time):
            #print('Step 5')
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=0, max_v=2, acc=20, dcc=20, jerk=0)
            self.motion_finished_forGUI(sleep_time, active_drive_number, target_pos=0, next_step=finalize_motion)  # Avoid passing `sleep_time`

        def finalize_motion(self, sleep_time=0):
            #print('Step final')
            """Executes GUI updates after the motion sequence is finished."""
            self.app.insert_message('1 Motor motion sequence ended')
            self.app.stop_fast_update()
            self.app.one_motor_indicator.configure(fg_color="orange")

        # Start motion sequence directly without threading
        start_motion_sequence(self, sleep_time_cycle)




    def two_motor_motion(self):
        """
        Executes a coordinated motion sequence for two LinMot drives.

        This function performs synchronized and alternating motion commands
        between two drives, demonstrating multi-axis control. It checks drive
        readiness and updates GUI indicators accordingly.

        Args:
            app (CTk): The main application instance.
        """
        if self.app.noDev < 2:
            self.app.insert_message("At least 2 drives are required for this motion profile.")
            return
        self.app.insert_message(f'2 Motor motion sequence started')
        self.app.two_motor_indicator.configure(fg_color="green")
        self.app.update()
        self.app.specific_update_interval[1] = self.app.cycle_time
        self.app.start_fast_update_thread()
        
        with self.app.lm_drive_lock.gen_rlock():
            operation_enabled_1 = self.app.lm_drive_data_dict[1].status['operation_enabled']
            homed_1 = self.app.lm_drive_data_dict[1].status['homed']
            error_1 = self.app.lm_drive_data_dict[1].status['error']
            ma_1 = self.app.lm_drive_data_dict[1].status['motion_active']
            operation_enabled_2 = self.app.lm_drive_data_dict[2].status['operation_enabled']
            homed_2 = self.app.lm_drive_data_dict[2].status['homed']
            error_2 = self.app.lm_drive_data_dict[2].status['error']
            ma_2 = self.app.lm_drive_data_dict[2].status['motion_active']
        sleep_time_cycle = max(self.app.cycle_time, 0.001)

        if operation_enabled_1 and operation_enabled_2 and homed_1 and homed_2 and not error_1 and not error_2 and not ma_1 and not ma_2:
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=50, max_v=0.01, acc=0.1, dcc=0.1, jerk=0, execute_mc=False)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=50, max_v=0.01, acc=0.1, dcc=0.1, jerk=0, execute_mc=False)
            self.app.sendData.send_data_to_slaves()
            self.motion_finished(sleep_time_cycle, active_drive_number = [1, 2], target_pos=[50, 50])
            self.send_motion_command(drive=1, header='Absolute_Sin', target_pos=0, max_v=0.1, acc=1, dcc=1, jerk=0, execute_mc=False)
            self.send_motion_command(drive=2, header='Absolute_Sin', target_pos=0, max_v=0.1, acc=1, dcc=1, jerk=0, execute_mc=False)
            self.app.sendData.send_data_to_slaves()
            self.motion_finished(sleep_time_cycle, active_drive_number = [1, 2], target_pos=[0, 0])
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=25, max_v=1, acc=10, dcc=10, jerk=0, execute_mc=False)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=50, max_v=1, acc=0.1, dcc=0.1, jerk=0, execute_mc=False)
            self.app.sendData.send_data_to_slaves()
            self.motion_finished(sleep_time_cycle, active_drive_number = [1, 2], target_pos=[25, 50])
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=50, max_v=0.01, acc=0.1, dcc=0.1, jerk=0, execute_mc=False)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=25, max_v=0.01, acc=0.1, dcc=0.1, jerk=0, execute_mc=False)
            self.app.sendData.send_data_to_slaves()
            self.motion_finished(sleep_time_cycle, active_drive_number = [1, 2], target_pos=[50, 25])
            self.send_motion_command(drive=1, header='Absolute_VAI', target_pos=0, max_v=0.03, acc=0.1, dcc=0.1, jerk=0)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=0, max_v=1, acc=1, dcc=1, jerk=0)
            self.motion_finished(sleep_time_cycle, active_drive_number = 2, target_pos=0)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=50, max_v=1, acc=1, dcc=1, jerk=0)
            self.motion_finished(sleep_time_cycle, active_drive_number = 2, target_pos=50)
            self.send_motion_command(drive=2, header='Absolute_VAI', target_pos=0, max_v=1, acc=1, dcc=1, jerk=0)
            self.motion_finished(sleep_time_cycle, active_drive_number = [1, 2], target_pos=[0, 0])
        else:
            self.app.insert_message(f'Motor not ready to start the motion. Please make sure that the operaion is enabled and motor is homed.')
        self.app.insert_message(f'2 Motor motion sequence ended')

        self.app.stop_fast_update()
        self.app.two_motor_indicator.configure(fg_color="orange")


    def send_motion_command(self, drive:int, header:str, target_pos:float, max_v:float, 
                            acc:float, dcc:float=0, jerk:int=0, execute_mc:bool=True):
        """
        Sends a motion command to a specified drive.

        This method formats and sends a motion command to the connected drive using parameters like
        position, velocity, acceleration, deceleration, and jerk.
        The header specifies the motion type (absolute, relative, etc.).

        Parameters:
            drive (int): The drive number to send the motion command to.
            header (str): The type of motion (e.g., "Absolute_VAI").
            target_pos (float): The target position for the motion.
            max_v (float): The maximum velocity for the motion.
            acc (float): The acceleration value.
            dcc (float, optional): The deceleration value. Defaults to 0 if not provided.
            jerk (float, optional): The jerk value. Defaults to 0 if not provided.
            execute_mc (bool, optional): When True this will execute the motion command imedately 
                                        after caclulation. Otherwise it will just write it to the LMDrive_Data.
                                        Defaults to True if not provided.

        Raises:
            ValueError: If the header is not recognized.
            TypeError: If argument is missing for the given Motion Command.
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
            "Relative_Sin": (0x0E10, True, False)
        }
        
        if header not in header_map:
            raise ValueError(f"Unsupported motion header: {header}")
        
        header_code, acc_combined, jerk_necessary = header_map[header]
        unit_scale = self.app.sendData.get_unit_scale(active_drive_number)

        pw = [
            [2, float(target_pos) * unit_scale],
            [2, float(max_v) * unit_scale * 100],
            [2, float(acc) * unit_scale * 10],
            [0], [0]
        ]
        if not acc_combined:
            if dcc == 0:
                raise TypeError(f"Missing required argument 'dcc' with header '{header}'")
            else:
                pw[3] = ([2, float(dcc) * unit_scale * 10])
        if jerk_necessary:
            if jerk == 0:
                raise TypeError(f"Missing required argument 'jerk' with header '{header}'")
            else:
                pw[4] = ([2, float(jerk) * unit_scale])

        self.app.sendData.update_output_drive_data(active_drive_number, controlWord = 0, header = header_code,
                                        para_word=pw, execute_mc=execute_mc)


    def motion_finished(self, sleep_time_cycle:float, active_drive_number, target_pos, timeout:float=60*10):
        """
        Waits for motion completion of one or more drives within a specified timeout.

        This function continuously monitors the motion status of the specified drive(s) and checks whether they have 
        reached their respective target positions. 

        It periodically checks the drive's status and sleeps between checks.

        Parameters:
            sleep_time_cycle (float): Interval in seconds between each motion status check.
            active_drive_number (int or list of int): Drive(s) to monitor. If a list is provided, all drives must complete motion.
            target_pos (int or list of int): Target position(s) corresponding to the drive(s).
                        Must match the order and length of active_drive_number.
            timeout (float): Maximum duration in seconds to wait for motion completion. Defaults to 300 seconds (5 minutes).

        Returns:
            bool: True if all specified drives complete motion and reach their target positions within the timeout.

        Raises:
            TypeError: input types are incorrect or if list lengths of drives and positions do not match.
            TimeoutError: If motion does not complete within the specified timeout period.
        """
        start_time = time.time()

        if isinstance(active_drive_number, list):
            if not isinstance(target_pos, list) or len(active_drive_number) != len(target_pos):
                raise TypeError('active_drive_number and target_pos must both be lists of the same length.')
        elif isinstance(active_drive_number, int) and not isinstance(target_pos, int):
            raise TypeError('active_drive_number and target_pos must both be integers if not lists.')

        while time.time() - start_time < timeout:
            with self.app.lm_drive_lock.gen_rlock():
                if isinstance(active_drive_number, list):
                    if all(not self.app.lm_drive_data_dict[d].status['motion_active'] for d in active_drive_number):
                        if all(
                            target_pos[i] - self.app.target_range <= self.app.lm_drive_data_dict[active_drive_number[i]].status['actual_position'] <= target_pos[i] + self.app.target_range
                            for i in range(len(active_drive_number))
                            ):
                            return True
                else:
                    if not self.app.lm_drive_data_dict[active_drive_number].status['motion_active']:
                        actual_pos = self.app.lm_drive_data_dict[active_drive_number].status['actual_position']
                        if target_pos - self.app.target_range <= actual_pos <= target_pos + self.app.target_range:
                            return True
            time.sleep(sleep_time_cycle * 2)

        raise TimeoutError(f"Motion did not finish within timeout.")


    def motion_finished_forGUI(self, sleep_time, drive_number, target_pos, next_step):
        """
        Schedules non-blocking checks for motion completion using `app.after()`.

        This function is used in GUI contexts to avoid freezing the interface
        while waiting for motion to complete. It triggers the next step once
        the drive(s) reach the target position.

        Args:
            app (CTk): The main application instance.
            sleep_time (float): Time interval between checks (in seconds).
            drive_number (int or list of int): Drive(s) to monitor.
            target_pos (int or list of int): Target position(s).
            next_step (Callable): Function to call once motion is complete.

        Raises:
            TypeError: If input types or list lengths are mismatched.
        """
        if isinstance(drive_number, list):
            if not isinstance(target_pos, list) or len(drive_number) != len(target_pos):
                raise TypeError('active_drive_number and target_pos must both be lists of the same length.')
        elif not isinstance(drive_number, int) or not isinstance(target_pos, int):
            raise TypeError('active_drive_number and target_pos must both be integers if not lists.')
        
        def check_motion():
            """Check if the motor motion has completed."""
            if isinstance(drive_number, list):
                with self.app.lm_drive_lock.gen_rlock():
                    all_inactive = all(not self.app.lm_drive_data_dict[i].status['motion_active'] for i in drive_number)
                    all_in_position = all(
                        target_pos[i] - self.app.target_range <= self.app.lm_drive_data_dict[i].status['actual_position'] <= target_pos[i] + self.app.target_range
                        for i in drive_number)
            elif isinstance(drive_number, int):
                with self.app.lm_drive_lock.gen_rlock():
                    all_inactive = not self.app.lm_drive_data_dict[drive_number].status['motion_active']
                    all_in_position = (
                        target_pos - self.app.target_range <= self.app.lm_drive_data_dict[drive_number].status['actual_position'] <= target_pos + self.app.target_range)
            else:
                raise TypeError('active_drive_number must be an integer or list')

            if all_inactive and all_in_position:
                next_step(self, sleep_time)
                return True

            # Schedule the next check
            if not self.app.shutting_down:
                after_id1 = self.app.after(int(sleep_time * 1000), check_motion)
                self.app.after_ids.append(after_id1)

        if not self.app.shutting_down:
            after_id2 = self.app.after(int(sleep_time * 2000), check_motion)
            self.app.after_ids.append(after_id2)


class Tab_Oscilloscope:
    """
    Description:
    ------------
    This module defines the Oscilloscope tab for the LinMot EtherCAT GUI 
    application. It enables real-time visualization of drive signals 
    (e.g., position, current, monitoring channels) with support for recording, 
    saving to CSV, and multi-axis plotting using matplotlib.
    """

    def __init__(self, app):
        self.app = app
        self.recording = False
        self.lock = threading.Lock()
        self.lock2 = threading.Lock()

    def setup_ui(self, drive_names):
        """
        Initializes the Oscilloscope tab layout and UI components.

        This includes:
        - A control panel with buttons for starting/stopping recording and saving data.
        - Checkboxes for selecting which signals to visualize.
        - A scrollable chart panel with matplotlib canvases for each drive.

        Args:
            self: The OscilloscopeTab instance.
        """
        self.drive_names = drive_names
        self.data = {drive: [] for drive in drive_names}

        self.osci_frame = self.app.tabview.tab("Oscilloscope")

        # Left control panel
        self.control_panel = ctk.CTkFrame(self.osci_frame)
        self.control_panel.grid(row=0, column=0, padx=10, pady=10, sticky="ns")

        self.start_button = ctk.CTkButton(self.control_panel, text="Start Recording", command=self.start_recording)
        self.start_button.pack(pady=5)

        self.stop_button = ctk.CTkButton(self.control_panel, text="Stop Recording", command=self.stop_recording)
        self.stop_button.pack(pady=5)

        self.save_button = ctk.CTkButton(self.control_panel, text="Save Data", command=self.save_data)
        self.save_button.pack(pady=5)

        self.checkboxes = {}
        default_checked = ["demand_pos", "actual_pos", "difference_pos", "demand_curr"]
        optional_checked = ["state_var", "status_word", "warn_word", "mon_ch1", "mon_ch2", "mon_ch3", "mon_ch4"]

        for key in default_checked + optional_checked:
            var = ctk.BooleanVar(value=key in default_checked)
            chk = ctk.CTkCheckBox(self.control_panel, text=key, variable=var)
            chk.pack(anchor="w", padx=10)
            self.checkboxes[key] = var

        # Right scrollable chart panel
        self.chart_panel = ctk.CTkScrollableFrame(self.osci_frame)
        self.chart_panel.grid(row=0, column=1, padx=10, pady=10, sticky="nsew")

        self.osci_frame.grid_rowconfigure(0, weight=1)
        self.osci_frame.grid_columnconfigure(1, weight=1)

        self.figures = {}
        self.axes = {}
        self.canvases = {}
        for drive in self.drive_names:
            fig, ax = plt.subplots(figsize=(8, 3))  # Fixed height
            self.figures[drive] = fig
            self.axes[drive] = ax
            canvas = FigureCanvasTkAgg(fig, master=self.chart_panel)
            canvas.get_tk_widget().pack(fill="x", expand=False, pady=10)
            self.canvases[drive] = canvas

    def on_close(self):
        """
        Stopps running recordings.

        Args:
            self: The OscilloscopeTab instance.
        """
        self.stop_recording()

    def start_recording(self):
        """
        Starts recording cyclic drive data.

        Clears previous data, activates the EtherCAT data queue, and launches
        background threads for data collection and batch processing.

        Args:
            self: The OscilloscopeTab instance.
        """
        self.data = {drive: [] for drive in self.drive_names} # Clear previous data
        self.recording = True
        self.app.ec_comm_process.data_queue_ON.set()
        self.app.start_oszi.set()
        self.thread = threading.Thread(target=self.record_data, daemon=True)
        self.thread.start()
        self.batch_thread = threading.Thread(target=self.process_batches, daemon=True)
        self.batch_thread.start()


    def stop_recording(self):
        """
        Stops the ongoing data recording process.

        Deactivates the EtherCAT data queue, joins the recording thread,
        clears the queue, and notifies the user via the message area.

        Args:
            self: The OscilloscopeTab instance.
        """
        self.app.ec_comm_process.data_queue_ON.clear()
        self.recording = False
        if self.thread and self.thread.is_alive():
                self.thread.join(timeout=1)

        while not self.app.ec_comm_process.data_queue.empty():
            try:
                self.app.ec_comm_process.data_queue.get_nowait()
            except queue.Empty:
                break
        self.app.insert_message('Oszi Queue has been emptied. Recording Stopped.')
        self.app.start_oszi.clear()

    def save_data(self):
        """
        Saves the recorded oscilloscope data to a CSV file.

        Prompts the user to select a file path and writes selected signal
        data for each drive to the file. Displays a success or error message.

        Args:
            self: The OscilloscopeTab instance.
        """
        file_path = filedialog.asksaveasfilename(
            defaultextension=".csv",
            filetypes=[("CSV files", "*.csv")],
            title="Save data as"
        )
        if not file_path:
            return  # User cancelled

        try:
            with open(file_path, "w", newline="") as f:
                writer = csv.writer(f)
                header = ["Drive", "Sample_Nr"] + [k for k, v in self.checkboxes.items() if v.get()]
                writer.writerow(header)
                with self.lock2:
                    for drive, rows in self.data.items():
                        for row in rows:
                            writer.writerow([drive] + row)
            self.app.insert_message(f"Data successfully saved to {file_path}")
        except Exception as e:
            self.app.insert_message(f"Error saving data: {e}")

    def record_data(self):
        """
        Collects real-time drive data from the EtherCAT communication queue.

        Unpacks binary input data for each drive, extracts selected signal values
        based on user checkboxes, and buffers the data for later batch processing.

        Runs continuously in a background thread while recording is active.

        Args:
            self: The OscilloscopeTab instance.
        """
        self.buffer = {drive: [] for drive in self.drive_names}
        self.last_update_time = {drive: time.time() for drive in self.drive_names}
        while self.recording:
            try:
                latest_data = None
                latest_data = self.app.ec_comm_process.data_queue.get_nowait()

                if latest_data is None:
                    time.sleep(self.app.cycle_time)
                    continue

                sample_nr, all_slave_data = latest_data
                for i, drive in enumerate(self.drive_names):
                    data_length = self.app.ec_comm_process.InputLength
                    device_data = bytes(all_slave_data[i * data_length:(i + 1) * data_length])
                    drive_data = self.oszi_Drive_Data(num_mon_channels=int(self.app.num_monit_ch_entry.get()))
                    drive_data.unpack_inputs(device_data)
                    drive_data.update_calculated_fields()
                    row = [sample_nr]
                    for key, var in self.checkboxes.items():
                        if var.get():
                            row.append(drive_data.status.get(key, None))
                    
                    with self.lock:
                        self.buffer[drive].append(row)

            except queue.Empty:
                pass
            time.sleep(self.app.cycle_time)


    def process_batches(self):
        """
        Processes buffered data and updates the GUI charts.

        Periodically transfers buffered data to the main data store and triggers
        chart updates using `after()` to avoid blocking the GUI thread.

        Args:
            self: The OscilloscopeTab instance.
        """
        while self.recording:
            time.sleep(1)
            with self.lock:
                for drive in self.drive_names:
                    if self.buffer[drive]:
                        with self.lock2:
                            self.data[drive].extend(self.buffer[drive])
                        self.buffer[drive].clear()
                        if self.recording and not self.app.shutting_down and self.app.winfo_exists():
                            after_id = self.app.after(0, self.update_chart, drive)
                            self.app.after_ids.append(after_id)
    
    def update_chart(self, drive):
        """
        Updates the matplotlib chart for a specific drive.

        Clears and redraws the chart with new data using three y-axes to separate
        distance, current, and raw signals. Applies legends and labels based on
        selected signals.

        Args:
            drive (str): The name of the drive whose chart should be updated.
        """
        fig = self.figures[drive]
        fig.clf()  # Clear the entire figure to remove all axes
        ax = fig.add_subplot(111)
        ax2 = ax.twinx()
        ax3 = ax.twinx()
        ax3.spines["right"].set_position(("outward", 60))

        self.axes[drive] = ax  # Update the reference
        self.canvases[drive].figure = fig  # Ensure canvas uses the updated figure

        distance_fields = {"demand_pos", "actual_pos", "difference_pos"}
        current_fields = {"demand_curr"}
        raw_fields = set(self.checkboxes.keys()) - distance_fields - current_fields

        selected_keys = [k for k, v in self.checkboxes.items() if v.get()]
        key_to_axis = {}

        for key in selected_keys:
            if key in distance_fields:
                key_to_axis[key] = ax
            elif key in current_fields:
                key_to_axis[key] = ax2
            else:
                key_to_axis[key] = ax3

        for idx, key in enumerate(selected_keys):
            with self.lock2:
                y = [row[1 + idx] for row in self.data[drive] if len(row) > 1 + idx]
            x = list(range(len(y)))
            axis = key_to_axis[key]
            color = "brown" if key == "demand_curr" else None
            axis.plot(x, y, label=key, color=color)

        ax.set_title(drive)
        ax.set_xlabel("Sample No.")
        ax.set_ylabel("Distance")
        ax2.set_ylabel("Current")
        ax3.set_ylabel("Raw Data")

        if any(k in distance_fields for k in selected_keys):
            ax.legend(loc="upper left")
        if any(k in current_fields for k in selected_keys):
            ax2.legend(loc="upper right")
        if any(k in raw_fields for k in selected_keys):
            ax3.legend(loc="lower right")

        self.canvases[drive].draw()

    class oszi_Drive_Data:
        def __init__(self, num_mon_channels):
            self.num_mon_ch = num_mon_channels  # Number of monitoring channels
            
            self.config = { 
                'pos_scale_numerator': 10000.0,             # Increments/Ticks per motor revolution (10000 for LinMot linear motors) 
                'pos_scale_denominator': 1.0,               # Units (mm/degrees/...) per motor revolution (1 for LinMot linear motors) 
                'unit_scale': 10000.0,                      # RO (read only): Unit scale. Is calculated automatically
                'modulo_factor': 360000,                    # Modulo increments (only used if isRotaryMotor = 1) 
                'fc_force_scale': 0.1,                      # Force scale
                'fc_torque_scale': 0.00057295779513082      # Torque scale
            }
            
            self.status = {
                'demand_pos': 0.0,
                'actual_pos': 0.0,
                'difference_pos': 0.0,
                'demand_curr': 0.0,
                'nr_of_revolutions': 0,
                'state_var': 0x0000,
                'status_word': 0x0000,
                'warn_word': 0x0000,
            }
            for i in range(1, self.num_mon_ch + 1):
                self.status[f'mon_ch{i}'] = 0x0000
            
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
            for i in range(1, self.num_mon_ch + 1):
                self.inputs[f'mon_ch{i}'] = 0x0000
        
        def update_calculated_fields(self):
            """
            Updates calculated fields based on current input values and configuration.
            """
            # Update `unit_scale` in config
            self.config['unit_scale'] = self.config['pos_scale_numerator'] / self.config['pos_scale_denominator']

            # Calculate scaled positions and current
            self.status['demand_pos'] = ctypes.c_int32(self.inputs['demand_pos']).value / self.config['unit_scale']
            self.status['actual_pos'] = ctypes.c_int32(self.inputs['actual_pos']).value / self.config['unit_scale']
            self.status['difference_pos'] = round(self.status['demand_pos'] - self.status['actual_pos'], 4)
            self.status['demand_curr'] = ctypes.c_int16(self.inputs['demand_curr']).value / 1000.0
            
            # Move other values
            self.status['state_var'] = self.inputs['state_var']
            self.status['status_word'] = self.inputs['status_word']
            self.status['warn_word'] = self.inputs['warn_word']
            
            for i in range(1, self.num_mon_ch + 1):
                self.status[f'mon_ch{i}'] = self.inputs[f'mon_ch{i}']
            
        def unpack_inputs(self, data):
            """
            Unpack input data from a binary structure, adjusting for the number of monitoring channels.
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
                
            # After unpacking, update calculated fields
            self.update_calculated_fields()






def main():
    print("do nothing")
    
if __name__ == "__main__":
    main()