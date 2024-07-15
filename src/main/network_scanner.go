package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"go-network-scanner-nmap/src/nmaputil"
	"go-network-scanner-nmap/src/server"
	"os"
	"syscall"
)

func main() {
	if !checkRootPermissions() {
		fmt.Println("Network scanner requires of root privileges. eg: sudo ./network-scanner -i <interface>")
		return
	}

	binDir := "bin"

	// Check if interface added
	ifaceName := pflag.StringP("interface", "i", "wlan0", "Interface to use for scanning (eg. eth0, wlan0, en0, etc.)")
	pflag.Parse()

	// Check if nmap exits
	nmapPath, nmapFound := nmaputil.CheckNmap()

	if !nmapFound {
		fmt.Println("nmap not found on the system.")
		var err error
		nmapPath, err = nmaputil.DownloadNmap(binDir)
		if err != nil {
			fmt.Println("Error downloading nmap:", err)
			return
		}
		fmt.Println("nmap downloaded and installed.")
	} else {
		fmt.Println("nmap is already installed on the system.")
	}

	if err := os.Chmod(nmapPath, 0755); err != nil {
		fmt.Println("Error setting executable nmap permissions:", err)
		return
	}

	server.StartServer(nmapPath, ifaceName)
}

func checkRootPermissions() bool {
	return syscall.Getuid() == 0
}
