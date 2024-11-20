package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	errNoGateway        = errors.New("no gateway found")
	errCantParseGateway = errors.New("can't parse gateway")
	errInvalidMac       = errors.New("invalid mac address")
	errNoMac            = errors.New("no mac address provided")
)

func main() {
	if len(os.Args) < 2 {
		panic(errNoMac)
	}

	inp := os.Args[1]
	if len(inp) < 2 {
		panic(errInvalidMac)
	}

	ip, err := discoverGateway()
	if err != nil {
		panic(err)
	}

	nmap, err := nmap(ip)
	if err != nil {
		panic(err)
	}

	results := searchIP(nmap, inp)
	if len(results) == 0 {
		fmt.Println("No matching MAC addresses found")
	} else {
		fmt.Printf("Found %d result(s):\n", len(results))
		for result := range results {
			fmt.Println(result)
		}
	}
}

func discoverGateway() (net.IP, error) {
	routeCmd := exec.Command("route", "print", "0.0.0.0")
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	ip, err := parseGateway(string(output))
	if err != nil {
		return nil, err
	}

	return ip, nil
}

func parseGateway(route string) (net.IP, error) {
	lines := strings.Split(route, "\n")
	sep := 0
	for idx, line := range lines {
		if sep == 3 {
			if len(lines) <= idx+2 {
				return nil, errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 5 {
				return nil, errCantParseGateway
			}

			ip := net.ParseIP(fields[2])
			if ip == nil {
				return nil, errCantParseGateway
			}

			return ip, nil
		}
		if strings.HasPrefix(line, "=======") {
			sep++
			continue
		}
	}

	return nil, errNoGateway
}

func nmap(ip net.IP) (string, error) {
	nmapCmd := exec.Command("nmap", "-sn", "-n", ip.String()+"/24")
	nmapCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := nmapCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func searchIP(nmapOutput, macInput string) map[string]string {
	results := make(map[string]string)
	lines := strings.Split(nmapOutput, "\n")
	var currentIP string

	for _, line := range lines {
		if strings.HasPrefix(line, "Nmap scan report for") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				currentIP = fields[4]
			}
		} else if strings.Contains(line, "MAC Address:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				macAddr := fields[2]
				vendor := strings.Join(fields[3:], " ")

				if strings.Contains(strings.ToLower(macAddr), strings.ToLower(macInput)) {
					results[fmt.Sprintf("%s = %s %s", macAddr, currentIP, vendor)] = currentIP
				}
			}
		}
	}
	return results
}
