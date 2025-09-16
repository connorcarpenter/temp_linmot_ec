"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           <filename>.py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  10.06.2025
    Version:        0.36
    Description:    Creates a graphical user interface, where LinMot Drives can  
                    be connected and interacted with.

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
    This script serves as the main graphical user interface (GUI) for a demo 
    project focused on LinMot drive communication via EtherCAT, developed for 
    NTI AG LinMot & MagSpring. It uses the customtkinter framework to create 
    a multi-tabbed application that allows users to configure, monitor, and 
    control LinMot drives in real time.

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
import multiprocessing
import threading
from readerwriterlock import rwlock
import queue
import time
import LinMot_GUI_Tabs_0v10 as tabs
import LinMot_DataHandling_0v10 as lmdh
import LinMot_EtherCAT_Comm_0v82e as commEC


class EtherCATApp(ctk.CTk):
    def __init__(self):
        """
        Initialize the EtherCAT GUI application.

        Sets up the main window, initializes global variables, configures the tabbed interface,
        and prepares the message area and communication indicator.
        """
        super().__init__()
        self.title("EtherCAT UI.py")
        self.geometry(f"{1100}x{850}")
        self.iconbitmap("linmot.ico")
        
        # Initialize Global Variables
        self.adapter_id = None
        self.adapter_desc = None
        self.comm_running = False
        self.status_entries = {}
        self.lm_drive_data_dict = {}
        self.lm_drive_lock = rwlock.RWLockFairD()
        self.add_drive_data = {}
        self.numDev = 0

        # User defined Values
        self.mp_logging: int = 50 # Logging level for multiprocessing EtherCat communication
        self.target_range = 1
        
        # Global Variables
        self.cycle_time = 0
        self.specific_update_interval = [False, 0.5] # [Enabled:bool, Update Interval (seconds):float]
        self.is_updating2 = False
        # UI state
        self.monitoring_entries = {}
        self.parameter_entries = {}
        self.Device_Name_entry = {}
        self.seg_button_1 = {}

        self.after_ids = []
        self.shutting_down = False
        
        # Oscilloscope
        self.start_oszi = threading.Event()
        self.save_oszi_while_running = threading.Event()
        self.csv_queue = queue.Queue()  # Queue to store data for writing
        
        # GUI Layout
        self.grid_rowconfigure(0, weight=1)
        self.grid_rowconfigure(1, weight=0)
        self.grid_columnconfigure(0, weight=1)
        
        self.tabview = ctk.CTkTabview(self)
        self.tabview.grid(row=0, column=0, padx=20, pady=10, sticky="nsew")
        for tab in ["EC Master Setup", "Control Status", "Simple Motion", "Motion Profile", "Oscilloscope", "UI Options"]:
            self.tabview.add(tab)
        
        self.drive_status_tab = tabs.Tab_Status(self)
        self.ec_setup_tab = tabs.Tab_SetupModule(self)
        self.simple_motion_tab = tabs.Tab_SimpleMotion(self)
        self.motion_profile_tab = tabs.Tab_MotionProfile(self)
        self.oscilloscope_tab = tabs.Tab_Oscilloscope(self)
        self.pro_comm_data = lmdh.Processing_comm_data(self)
        self.sendData = lmdh.Send_Data(self)
        self.msg_other = lmdh.OtherMessages(self)
        
        # Setup the "Setup" tab layout
        self.ec_setup_tab.setup_tab()
        self.ui_options()
        
        # Messages section
        title_label = ctk.CTkLabel(self, text="Messages", font=("Arial", 10))
        title_label.grid(row=1, column=0, padx=20, pady=(0, 0), sticky="w")
        self.text_field = ctk.CTkTextbox(self, height=80, width=80)
        self.text_field.grid(row=2, column=0, padx=(20, 100), pady=5, sticky="ew")
        self.text_field.configure(state="normal")
        self.text_field.insert(ctk.END, "Message here:\n")
        self.text_field.configure(state="disabled")
        self.text_field.see(ctk.END)
        
        # Create an indicator light on the bottom right corner
        self.indicator_light_comm = ctk.CTkLabel(self, text="No Comm", font=("Arial", 8), width=40, height=30, corner_radius=15, bg_color="red")
        self.indicator_light_comm.grid(row=2, column=0, padx=(0, 20), pady=5, sticky="se")
        
        # Bind closing event
        self.protocol("WM_DELETE_WINDOW", self.on_closing)
        
    def insert_message(self, text):
        """
        Insert a message into the GUI's message text field.

        Args:
            text (str): The message to display.
        """
        self.text_field.configure(state="normal")
        self.text_field.insert(ctk.END, f"{text}\n")
        self.text_field.yview_moveto(1.0)
        self.text_field.configure(state="disabled")

    def on_closing(self):
        """
        Handle the window close event.

        Cancels any running threads or scheduled tasks, stops EtherCAT communication if active,
        and safely destroys the GUI window.
        """
        self.shutting_down = True
        # Check if EtherCAT communication is running and stop it
        if hasattr(self, "ec_comm_process") and self.ec_comm_process.comm_proc and self.ec_comm_process.comm_proc.is_alive():
            self.stop_communication()
        time.sleep((int(self.update_freq_entry.get())/1000)*2)

        for after_id in getattr(self, 'after_ids', []):
            try:
                self.after_cancel(after_id)
            except Exception as e:
                print(f"Failed to cancel after callback {after_id}: {e}")
        self.after_ids.clear()
        
        self.protocol("WM_DELETE_WINDOW", self.destroy)
        #self.destroy()
    
    def ui_options(self):
        """
        Configure the UI Options tab.

        Sets up appearance mode, UI scaling, and update frequency settings.
        """
        # Left side widgets
        left_frame_UI = ctk.CTkFrame(self.tabview.tab("UI Options"))
        left_frame_UI.grid(row=0, column=0, padx=20, pady=20, sticky="nw")
        
        # Title for UI Settings
        left_frame_title = ctk.CTkLabel(left_frame_UI, text="UI Settings", font=("Arial", 14, "bold"))
        left_frame_title.grid(row=0, column=0, padx=20, pady=(10, 10))
        
        # Apperance
        self.appearance_mode_label = ctk.CTkLabel(left_frame_UI, text="Appearance Mode:", anchor="w")
        self.appearance_mode_label.grid(row=1, column=0, padx=20, pady=(10, 0))
        self.appearance_mode_optionemenu = ctk.CTkOptionMenu(left_frame_UI, values=["Light", "Dark", "System"],
                                                                       command=self.change_appearance_mode_event)
        self.appearance_mode_optionemenu.grid(row=2, column=0, padx=20, pady=(10, 10))
        self.scaling_label = ctk.CTkLabel(left_frame_UI, text="UI Scaling:", anchor="w")
        self.scaling_label.grid(row=3, column=0, padx=20, pady=(10, 0))
        self.scaling_optionemenu = ctk.CTkOptionMenu(left_frame_UI, values=["80%", "90%", "100%", "110%", "120%"],
                                                               command=self.change_scaling_event)
        self.scaling_optionemenu.grid(row=4, column=0, padx=20, pady=(10, 20))
        
        # Update Frequency
        # Middle Frame for Software Settings
        middle_frame = ctk.CTkFrame(self.tabview.tab("UI Options"))
        middle_frame.grid(row=0, column=1, padx=20, pady=20, sticky="nw")

        # Title for Software Settings
        middle_frame_title = ctk.CTkLabel(middle_frame, text="Software Settings", font=("Arial", 14, "bold"))
        middle_frame_title.grid(row=0, column=0, padx=20, pady=(10, 10))

        # Update Frequency input field
        self.update_freq_label = ctk.CTkLabel(middle_frame, text="Update Frequency [ms]:", anchor="w")
        self.update_freq_label.grid(row=1, column=0, padx=20, pady=(10, 0))

        self.update_freq_entry = ctk.CTkEntry(middle_frame)
        self.update_freq_entry.insert(0, "500")  # Default value
        self.update_freq_entry.grid(row=2, column=0, padx=20, pady=(10, 20))
        
    def start_communication(self):
        """
        Start EtherCAT communication with connected LinMot drives.

        Initializes communication parameters, validates user input, starts the communication process,
        and sets up drive data structures and GUI tabs.
        """
        self.indicator_light_comm.configure(bg_color="orange", text="Connecting...")
        self.update()
        # Get all data
        try:
            adapter_id = self.adapter_id
            if not adapter_id:
                raise ValueError("Adapter ID not defined!")
        except ValueError as e:
            self.insert_message(str(e))  # Error Message!
            return
        try:
            self.noDev = int(self.num_devices_entry.get())
            self.cycle_time = float(self.cycle_time_entry.get()) / 1000
        except ValueError as e:
            self.insert_message('Invalid number of devices or cycle time! Communication cannot be started.')
            return
        try:
            no_Monitoring = int(self.num_monit_ch_entry.get())
            no_Parameter = int(self.num_para_ch_entry.get())
        except ValueError as e:
            self.insert_message('Invalid no. of Monitoring and Parameter Channels! Communication cannot be started.')
            return
        self.lock = multiprocessing.Lock()  # Lock for synchronizing access to the data array
        Activate_LMDrive_Data = False # Do not provide LM_Drive data, but only the raw values (beter performance)

        # Make EtherCAT Setup Read Only
        self.num_devices_entry.configure(state="readonly")
        self.cycle_time_entry.configure(state="readonly")
        self.num_monit_ch_entry.configure(state="readonly")
        self.num_para_ch_entry.configure(state="readonly")
        
        # Create an instance of the EtherCATCommunication class
        self.ec_comm_process = commEC.EtherCATCommunication(
            adapter_id, self.noDev, self.cycle_time, self.lock,
            no_Monitoring, no_Parameter, Activate_LMDrive_Data, self.mp_logging)
        self.ec_comm_process.start()

        if self.ec_comm_process.comm_proc and self.ec_comm_process.comm_proc.is_alive():
            j = 1
            while bool(j):
                EC_is_running = not self.ec_comm_process.stop_event.wait(timeout=1)
                #print(f'Wait for the master to establish communication with the drive.')
                if not EC_is_running:
                    time.sleep(0.2)
                    j += 1
                    if j > 30:
                        EC_is_running = False
                        j = 0
                else:
                    j = 0
            
            if EC_is_running:
                self.indicator_light_comm.configure(bg_color="green", text="Comm Active")
                self.update()
                self.insert_message('EtherCAT communication process is running.')
                self.comm_running = True
                # Setup Dictionary for every Drive
                for i in range(self.noDev):
                    self.lm_drive_data_dict[i+1] = commEC.LMDrive_Data(num_mon_channels=no_Monitoring, num_par_channels=no_Parameter)
                self.updating_values_in_active_window()  # Updating values only in active window
                
                # Setup all Tabs
                self.ec_setup_tab.unlock_drive_info()
                self.drive_status_tab.drive_status_tab(self.device_tab_view._segmented_button._value_list)
                self.simple_motion_tab.simple_motion_tab(self.device_tab_view._segmented_button._value_list)
                self.motion_profile_tab.setup_motion_profile_tab()
                self.oscilloscope_tab.setup_ui(self.device_tab_view._segmented_button._value_list)
            else:
                self.insert_message('EtherCAT communication process is NOT running!')
                self.stop_communication()
        
    def stop_communication(self):
        """
        Stop EtherCAT communication and clean up resources.

        Terminates background threads, stops oscilloscope recording if active,
        and resets the communication indicator.
        """
        self.indicator_light_comm.configure(bg_color="red", text="No Comm")
        self.update()
        self.stop_fast_update()

        if self.start_oszi.is_set():
            try:
                self.oscilloscope_tab.on_close()
            except Exception as e:
                self.insert_message(f"Error stopping oscilloscope: {e}")
        try:
            if self.mp_logging != 0:
                while not self.ec_comm_process.error_queue.empty(): self.insert_message(f'Error COMM: {self.ec_comm_process.error_queue.get()}')
                while not self.ec_comm_process.info_queue.empty(): self.insert_message(f'Info COMM: {self.ec_comm_process.info_queue.get()}')
            # Ensure that the EtherCAT communication process is stopped properly
            self.ec_comm_process.stop()
            self.insert_message('EtherCAT communication has been stopped.')
        except AttributeError:
            self.insert_message(f"Commmunication does not exist -> 'Stop Communication' not possible")
        except Exception as e:
            self.insert_message(f"Failed to stop commmunication: {e}")
        finally:
            self.comm_running = False
        
    def change_appearance_mode_event(self, new_appearance_mode: str):
        """
        Change the appearance mode of the GUI.

        Args:
            new_appearance_mode (str): The selected appearance mode ("Light", "Dark", or "System").
        """
        ctk.set_appearance_mode(new_appearance_mode)
        
    def change_scaling_event(self, new_scaling: str):
        """
        Change the UI scaling factor.

        Args:
            new_scaling (str): The selected scaling percentage (e.g., "100%").
        """
        new_scaling_float = int(new_scaling.replace("%", "")) / 100
        ctk.set_widget_scaling(new_scaling_float)
        
    def updating_values_in_active_window(self):
        """
        Periodically update drive status values in the active GUI tab.

        Args:
            update_interval (int, optional): Update interval in milliseconds. Defaults to 1000.
        """
        self.is_updating = False
        
        def update_field():
            # Flag to prevent overlapping calls
            if self.is_updating:
                return  # Prevent overlapping updates
            self.is_updating = True
            
            if not self.specific_update_interval[0]:
                self.pro_comm_data.process_input_data(data_length = self.ec_comm_process.InputLength)

            # Add Comm logging continuously to text_field
            if self.mp_logging != 0:
                while not self.ec_comm_process.info_queue.empty():
                    self.insert_message(f'Info COMM: {self.ec_comm_process.info_queue.get()}')
                while not self.ec_comm_process.error_queue.empty():
                    self.insert_message(f'Error COMM: {self.ec_comm_process.error_queue.get()}')
                        
            # Update Values in Active Window
            active_tab0 = self.tabview.get()
            match active_tab0:
                case 'Control Status':
                    self.drive_status_tab.monitor_status(self.drive_tabview.get())
            self.is_updating = False
            
            # Check if communication is still running
            if self.ec_comm_process.stop_event.is_set() and self.comm_running:
                user_choice = self.msg_other.msg_CommError()
                if user_choice:
                    self.comm_running = False
                    self.stop_communication()

        # Schedule the next update
        update_field()
        if self.comm_running and not self.shutting_down and self.winfo_exists():
            after_id0 = self.after(int(self.update_freq_entry.get()), self.updating_values_in_active_window)
            self.after_ids.append(after_id0)

    
    def fast_update_drive_data(self):
        """
        Continuously update drive data in a separate thread while fast update is enabled.
        """
        if self.is_updating2:
            return  # Prevent overlapping updates
        self.is_updating2 = True
        
        while self.specific_update_interval[0]:  # Keep running while it's True
            self.pro_comm_data.process_input_data(data_length=self.ec_comm_process.InputLength)

            # Simulate data processing
            time.sleep(max(self.specific_update_interval[1]/2, 0.001))  # Wait for the next update

        self.is_updating2 = False  # Reset flag when stopping

    def start_fast_update_thread(self):
        """
        Start the fast update thread if it is not already running.
        """
        if not self.is_updating2:
            self.specific_update_interval[0] = True  # Enable fast updates
            self.update_thread = threading.Thread(target=self.fast_update_drive_data, daemon=True)
            self.update_thread.start()

    def stop_fast_update(self):
        """
        Stop the fast update loop.
        """
        self.specific_update_interval[0] = False  # Disable fast updates
            
    
# Run the application
if __name__ == "__main__":
    app = EtherCATApp()
    app.mainloop()
