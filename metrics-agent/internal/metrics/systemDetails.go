package metrics

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"metrics-agent/internal/models"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	stdnet "net"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func formatDuration(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	mins := d / time.Minute
	d -= mins * time.Minute
	secs := d / time.Second
	return fmt.Sprintf("%d day(s) %d hr(s) %d min(s) %d sec(s)", days, hours, mins, secs)
}

func GetSystemSummary(userID string) (*models.SystemSummary, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Partitions(true)
	netInterfaces, _ := net.Interfaces()
	procs, _ := process.Processes()
	currentUser, _ := user.Current()
	systemLogErrorCount, err := GetSystemErrorLogCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get system error log count: %v", err)
	}
	loginCount, err := getLoginCount()
	if err != nil {
		fmt.Println("Error getting login count:", err)
	} else {
		fmt.Println("Login Count:", loginCount)
	}

	portCount, err := getOpenPortCount()
	if err != nil {
		fmt.Println("Error getting open port count:", err)
	} else {
		fmt.Println("Open Port Count:", portCount)
	}
	ip := getIP()
	if ip == "" {
		ip = "N/A"
	}
	summary := &models.SystemSummary{
		UserID:            userID,
		Hostname:          hostInfo.Hostname,
		IPAddress:         ip,
		OS:                fmt.Sprintf("%s %s (%s)", hostInfo.Platform, hostInfo.PlatformVersion, runtime.GOARCH),
		CPUModel:          cpuInfo[0].ModelName,
		CPUCores:          runtime.NumCPU(),
		RAMMB:             float64(memInfo.Total) / (1024 * 1024),
		DiskCount:         len(diskInfo),
		SysLogsErrorCount: systemLogErrorCount, // Placeholder for sys logs error count, if applicable
		LoginCount:        loginCount,
		OpenPortCount:     portCount,
		Uptime:            formatDuration(hostInfo.Uptime),
		BootTime:          hostInfo.BootTime,
		TotalProcesses:    len(procs),
		NICCount:          len(netInterfaces),
		CurrentUser:       currentUser.Username,
	}

	return summary, nil
}

func getIP() string {
	var ipList []string

	ifaces, _ := stdnet.Interfaces()
	for _, iface := range ifaces {
		if iface.Flags&stdnet.FlagUp != 0 && iface.Flags&stdnet.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				var ip stdnet.IP
				switch v := addr.(type) {
				case *stdnet.IPNet:
					ip = v.IP
				case *stdnet.IPAddr:
					ip = v.IP
				}
				if ip == nil || ip.IsLoopback() {
					continue
				}
				ipList = append(ipList, ip.String())
			}
		}
	}

	if len(ipList) == 0 {
		return "N/A"
	}
	return ipList[1] // Return the first non-loopback IP address found
}

func GetSystemErrorLogCount() (int, error) {
	switch runtime.GOOS {
	case "linux":
		return getLinuxSyslogErrorCount()
	case "darwin":
		return getMacSyslogErrorCount()
	case "windows":
		return getWindowsEventLogErrorCount()
	default:
		return 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func getLinuxSyslogErrorCount() (int, error) {
	files := []string{"/var/log/syslog", "/var/log/messages"}
	var count int

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			f, err := os.Open(file)
			if err != nil {
				return 0, err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(strings.ToLower(line), "error") {
					count++
				}
			}
			return count, nil
		}
	}
	return 0, fmt.Errorf("no syslog file found")
}
func getMacSyslogErrorCount() (int, error) {
	cmd := exec.Command("log", "show", "--predicate", "eventMessage contains 'error'", "--style", "syslog", "--last", "1h")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines := bytes.Split(output, []byte("\n"))
	return len(lines) - 1, nil
}
func getWindowsEventLogErrorCount() (int, error) {
	cmd := exec.Command("powershell", "-Command", `
		(Get-WinEvent -LogName System | Where-Object { $_.LevelDisplayName -eq 'Error' } | Measure-Object).Count
	`)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to execute command: %v", err)
	}

	countStr := strings.TrimSpace(out.String())
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse count: %v", err)
	}

	return count, nil
}

func getLoginCount() (int, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		// Unix-based systems: use 'who'
		out, err := exec.Command("who").Output()
		if err != nil {
			return 0, err
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		return len(lines), nil
	case "windows":
		// Windows: use PowerShell for event ID 4624 (successful login)
		cmd := exec.Command("powershell", "-NoProfile", "-Command",
			`Try {
				$events = Get-WinEvent -FilterHashtable @{LogName='Security'; Id=4624}
				$events.Count
			} Catch {
				Write-Error $_.Exception.Message
				Exit 1
			}`)

		var out, stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return 0, fmt.Errorf("PowerShell error: %s", stderr.String())
		}

		countStr := strings.TrimSpace(out.String())
		return parseCount(countStr)
	default:
		return 0, errors.New("unsupported OS")
	}
}

// Get open port count (works on Linux/macOS/Windows)
func getOpenPortCount() (int, error) {
	switch runtime.GOOS {
	case "linux", "darwin":
		// Try 'ss', fallback to 'netstat'
		cmd := exec.Command("ss", "-tuln")
		out, err := cmd.Output()
		if err != nil {
			cmd = exec.Command("netstat", "-tuln")
			out, err = cmd.Output()
			if err != nil {
				return 0, err
			}
		}
		lines := strings.Split(string(out), "\n")
		count := 0
		for _, line := range lines {
			if strings.Contains(line, "LISTEN") {
				count++
			}
		}
		return count, nil
	case "windows":
		// Windows: use netstat
		cmd := exec.Command("netstat", "-an")
		out, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		lines := strings.Split(string(out), "\n")
		count := 0
		for _, line := range lines {
			if strings.Contains(line, "LISTENING") {
				count++
			}
		}
		return count, nil
	default:
		return 0, errors.New("unsupported OS")
	}
}

// Helper to parse count string to int
func parseCount(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}
