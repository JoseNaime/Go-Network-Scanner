package nmaputil

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type NmapScanResult struct {
	IP         net.IP
	Details    string
	MACAddress string
	DeviceType string
}

// RunCommand runs a command and returns its output and any errors
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func ScanIpsWithNmap(ips []net.IP, nmapPath string) []NmapScanResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []NmapScanResult{}

	for _, ip := range ips {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()

			// Prepare the nmap command
			cmd := exec.Command(nmapPath, "-sP", "-sV", ip.String())

			// Capture the output
			var nmapOut bytes.Buffer
			cmd.Stdout = &nmapOut
			cmd.Stderr = &nmapOut

			// Run the command
			if err := cmd.Run(); err != nil {
				mu.Lock()
				fmt.Println("Error executing nmap:", err)
				mu.Unlock()
				return
			}

			output := nmapOut.String()
			if strings.Contains(output, "Host is up") {
				macAddress, deviceType := extractMACAndDeviceType(output)
				mu.Lock()
				results = append(results, NmapScanResult{
					IP:         ip,
					Details:    output,
					MACAddress: macAddress,
					DeviceType: deviceType,
				})
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	return results
}

func extractMACAndDeviceType(output string) (string, string) {
	re := regexp.MustCompile(`MAC Address: ([0-9A-Fa-f:]+) \((.+)\)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}
