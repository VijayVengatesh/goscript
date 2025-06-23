package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Config struct {
	UserID string `json:"user_id"`
}

type Metrics struct {
	UserID        string  `json:"user_id"`
	Hostname      string  `json:"hostname"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsed    uint64  `json:"memory_used"`
	MemoryTotal   uint64  `json:"memory_total"`
	DiskUsed      uint64  `json:"disk_used"`
	DiskTotal     uint64  `json:"disk_total"`
	Uptime        uint64  `json:"uptime_seconds"`
	MetricGetTime string  `json:"metric_get_time"` // ➕ added UTC timestamp

}

func getConfigPath() string {
	var baseDir string
	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
	case "darwin":
		baseDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default:
		baseDir = "/etc"
	}
	return filepath.Join(baseDir, "metrics-agent", "config.json")
}

func loadUserID() (string, error) {
	path := getConfigPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg.UserID, nil
}

func promptAndSaveUserID() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your User ID (or Device Key): ")
	userID, _ := reader.ReadString('\n')
	userID = string(bytes.TrimSpace([]byte(userID)))

	config := Config{UserID: userID}
	data, _ := json.MarshalIndent(config, "", "  ")

	os.MkdirAll(filepath.Dir(getConfigPath()), 0700)
	_ = ioutil.WriteFile(getConfigPath(), data, 0600)

	fmt.Println("✔ User ID stored successfully")
	return userID
}

func collectMetrics(userID string) (*Metrics, error) {
	var path string
	if runtime.GOOS == "windows" {
		path = "C:\\"
	} else {
		path = "/"
	}

	vm, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	diskStat, _ := disk.Usage(path)
	hostStat, _ := host.Info()
	utcNow := time.Now().UTC().Format(time.RFC3339) // e.g., "2025-06-18T06:15:04Z"

	return &Metrics{
		UserID:        userID,
		Hostname:      hostStat.Hostname,
		CPUPercent:    cpuPercent[0],
		MemoryUsed:    vm.Used,
		MemoryTotal:   vm.Total,
		DiskUsed:      diskStat.Used,
		DiskTotal:     diskStat.Total,
		Uptime:        hostStat.Uptime,
		MetricGetTime: utcNow,
	}, nil
}

func sendToAPI(metrics *Metrics) {
	jsonData, _ := json.Marshal(metrics)
	resp, err := http.Post("https://cloudops-api.idevopz.com/metrics", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("✅ Metrics sent! Status:", resp.Status)
}

func main() {
	userID, err := loadUserID()
	if err != nil || userID == "" {
		userID = promptAndSaveUserID()
	}

	for {
		metrics, err := collectMetrics(userID)
		if err == nil {
			sendToAPI(metrics)
		}
		time.Sleep(10 * time.Second)
	}
}
