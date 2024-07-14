package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
)

func GetActiveIPsInSubnet(ipNet *net.IPNet) []net.IP {
	var ips []net.IP
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)
		ips = append(ips, ipCopy)
	}

	activeIps := concurrentPingSweep(ips)
	return activeIps
}

func concurrentPingSweep(ips []net.IP) []net.IP {
	var activeIPs []net.IP
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, ip := range ips {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()
			if ping(ip.String()) {
				mu.Lock()
				activeIPs = append(activeIPs, ip)
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	return activeIPs
}

func GetInterfaceIPv4Addr(ifaceName string) (net.IP, *net.IPNet, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
			return ipNet.IP, ipNet, nil
		}
	}
	return nil, nil, fmt.Errorf("no suitable IPv4 address found for interface %s", ifaceName)
}

func ping(ip string) bool {
	out, err := exec.Command("ping", "-c", "1", "-W", "5", ip).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "1 received")
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
