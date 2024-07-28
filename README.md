
# Network Scanner (NScanner)
![Version](https://img.shields.io/badge/version-0.1-blue)
![macOS](https://img.shields.io/badge/Tested%20on-macOS-blue?logo=apple)
![Linux](https://img.shields.io/badge/Tested%20on-Linux-blue?logo=linux)
![Not Tested on Windows](https://img.shields.io/badge/tested%20on-Windows-FFD700?logo=windows&label=Not%20Tested)


This project is a network scanner built with [Fiber](https://gofiber.io/), a web framework for 
Go. It serves static files from the `public` directory and performs network scans using 
concurrent ping sweeps and `nmap`.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Building the Project](#building-the-project)
- [Running the Server](#running-the-server)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Development Plan](#development-plan)
  - [Future Enhancements](#future-enhancements)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [nmap](https://nmap.org/download.html) installed on your system
- [Fiber](https://gofiber.io/) (version 2.x)

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/network-scanner.git
    cd network-scanner
    ```

2. Install the required Go packages:
    ```sh
    go get github.com/gofiber/fiber/v2
    go get github.com/gofiber/websocket/v2
    go get github.com/google/gopacket/pcap
    go get github.com/spf13/pflag
    ```

## Building the Project

To build the project, run:

```sh
go build -o network-scanner src/main/network_scanner.go
```

## Running the Server

To start the Fiber server, run:

```sh
sudo ./network-scanner -i <interface>
```

Replace `<interface>` with the network interface you want to scan, e.g., `eth0`, `wlan0`, etc. The `sudo` command is necessary to ensure the scanner has the required permissions to perform network scans.

The server will start on port by default `5234`.

## Usage

### Serving Static Files

The server will serve all static files from the `public` directory. Open the link provided by 
Fiber to open your browser and see the scan.

### Network Scanning

The server performs network scans by running concurrent ping requests to all possible ips inside 
the network. Then `nmap` is used to get more details of the devices encountered. Make sure 
`nmap` is 
installed 
and accessible from your system's PATH.

Each scan is executed hourly until the software is stopped. This interval can be modified by changing the value in `server.go` at line 100:
```go
// Delay between scans
time.Sleep(1 * time.Hour) // Modify this value to change the scan interval
```

## Project Structure

```
project/
│
├── public/                  # Directory for static files
│   └── static/
│       ├── index.html       # Main HTML file
│       ├── reset.css        # CSS reset file
│       ├── script.js        # JavaScript file
│       └── styles.css       # CSS file
├── src/                     # Source directory for other Go packages
│   ├── main/
│   │   └── network_scanner.go  # Main application file
│   ├── network/
│   │   └── network.go       # Network utility functions
│   ├── nmaputil/
│   │   └── nmaputil.go      # Network mapping utilities
│   └── server/
│       └── server.go        # Server functions and configurations
└── README.md                # Project README file
```


## API Endpoints

### Root (`/`)

Serves the `index.html` file and all static files in the `public` directory.

### WebSocket (`/ws`)

WebSocket endpoint for real-time communication.

## Development Plan

### Future Enhancements

1. **AI Integration for Alerts**:
    - Implement basic AI to monitor network traffic and detect suspicious connections.
    - Send push notifications to email or phone numbers when suspicious activity is detected.

2. **Customizable Scan Intervals**:
    - Allow users to configure the scan interval through command-line flags.

3. **Additional Nmap Flags**:
    - Provide more customization options by supporting additional nmap flags.

4. **Improved UI**:
    - Enhance the user interface for better usability and visual appeal.

5. **Cross-Platform Support**:
    - Ensure compatibility and testing on Windows systems.

6. **Logging and Reporting**:
    - Add detailed logging and reporting features to keep track of scan results and alerts.

7. **Containerization**:
    - Create Docker images for easy deployment and management.

8. **Performance Optimization**:
    - Optimize the scanning process to handle large networks more efficiently.

9. **Documentation**:
    - Improve and expand documentation, including usage guides and API references.

## Contributing

### Contribution Guidelines
- **Ideas and Suggestions**: Open an issue to discuss new features or enhancements.
- **Bug Reports**: Report bugs via the issue tracker, providing as much detail as possible.
- **Code Contributions**: Fork the repository, create a new branch, make changes, and open a pull request.

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-name`.
3. Make your changes.
4. Commit your changes: `git commit -m 'Add new feature'`.
5. Push to the branch: `git push origin feature-name`.
6. Open a pull request.
7. 

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.