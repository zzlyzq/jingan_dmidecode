# System Information & RAID Utility
This utility is designed to gather and display system information, RAID configuration, and IPMI (Intelligent Platform Management Interface) details in an efficient and user-friendly manner.
## Features
- Gather RAID information and configurations.
- Retrieve and display detailed RAID controller and drive information.
- Collect IPMI information, including IP, MAC address, user list, and more.
- Cross-platform support with specific handling for Windows systems.
## Prerequisites
Before you run the utility, ensure that you have the following installed on your system:
- `storcli` or `storcli64` (for RAID operations)
- `ipmicfg` or equivalent IPMI tool supported by your hardware.
## Getting Started
To get started with this tool, clone the repository to your local machine:

```bash
git clone <https://github.com/your-username/system-info-raid-utility.git>
cd system-info-raid-utility
Usage
To run the utility and display the system information, use the following command:

go run main.go
Ensure that the storcli and ipmi tools are in your system's PATH or modify the utility to point to their locations.
```

# Contributing
Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

# Fork the Project
Create your Feature Branch (git checkout -b feature/AmazingFeature)
Commit your Changes (git commit -m 'Add some AmazingFeature')
Push to the Branch (git push origin feature/AmazingFeature)
Open a Pull Request

# License
Distributed under the MIT License. See LICENSE for more information.
