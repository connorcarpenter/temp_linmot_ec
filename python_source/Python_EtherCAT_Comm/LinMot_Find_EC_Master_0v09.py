"""
==============================================================================
    Project:        Demo Project for LinMot Drive Communication with EtherCAT
    File:           Find_EC_Master_ .py
    Author:         AP
    Created:        22.08.2024
    Last Modified:  22.05.2025
    Version:        0.09
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
    This script assists in identifying and selecting the appropriate EtherCAT 
    master adapter.

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
import sys
try:
    import tkinter as tk
except ImportError:
    tk = None



def on_select(result_container):
    if selected_adapter.get():
        selected_name = selected_adapter.get()
        selected_desc = adapter_dict[selected_name]
        result_container['result'] = (selected_name, selected_desc)
        #root.quit()
        root.destroy()

def on_cancel(result_container):
    result_container['result'] = None
    #root.quit()
    root.destroy()

def create_window(adapters, result_container, parent=None):
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

    tk.Button(button_frame, text="Select", command=lambda: on_select(result_container)).pack(side="left", padx=5)
    tk.Button(button_frame, text="Cancel", command=lambda: on_cancel(result_container)).pack(side="right", padx=5)


def main(parent=None):
    adapters = adapter_list()
    if not adapters:
        print("No adapters found.")
        return None

    result_container = {'result': None}

    if tk is not None:
        try:
            create_window(adapters, result_container, parent)
            if parent:
                root.wait_window()
            else:
                root.mainloop()
        except Exception as e:
            print(f"GUI fallback to CLI due to error: {e}")
            return create_cli(adapters)
    else:
        return create_cli(adapters)

    return result_container['result']







def adapter_list():
    return pysoem.find_adapters()


def create_cli(adapters):
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


if __name__ == "__main__":
    selected_adapter_details = main()
    if selected_adapter_details:
        adapter_name, adapter_desc = selected_adapter_details
        print(f"{adapter_name} || {adapter_desc}")
    else:
        print("None")
    input('Press Enter to exit:')
    sys.exit(0)