package nmaputil

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
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

func CheckNmap() (string, bool) {
	nmapPath, err := exec.LookPath("nmap")
	if err != nil {
		return "", false
	}
	return nmapPath, true
}

func DownloadNmap(binDir string) (string, error) {
	nmapURL := "https://nmap.org/dist/nmap-7.80.tgz" // Adjust the URL to the appropriate version and platform
	nmapPath := filepath.Join(binDir, "nmap")

	fmt.Println("Downloading nmap...")
	cmd := exec.Command("wget", "-O", "nmap.tgz", nmapURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("tar", "-xzf", "nmap.tgz", "-C", binDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("rm", "nmap.tgz")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = exec.Command("chmod", "+x", nmapPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return nmapPath, nil
}
