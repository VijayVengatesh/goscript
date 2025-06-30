package metrics

import (
	"runtime"
	"time"

	"metrics-agent/internal/models"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func Collect(userID string) (*models.Metrics, error) {
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
	utcNow := time.Now().UTC().Format(time.RFC3339)
	// Calculate % usage
	memUsage := (float64(vm.Used) / float64(vm.Total)) * 100
	diskUsage := (float64(diskStat.Used) / float64(diskStat.Total)) * 100
	cpu := cpuPercent[0]

	// Determine status based on thresholds
	status := "up"
	if cpu > 90 || memUsage > 90 || diskUsage > 95 {
		status = "critical"
	} else if cpu > 80 || memUsage > 80 || diskUsage > 85 {
		status = "trouble"
	} else if cpu < 5 && memUsage < 10 && diskUsage < 10 {
		status = "down" // very unlikely; used only for offline or idle systems
	}
	return &models.Metrics{
		UserID:        userID,
		Hostname:      hostStat.Hostname,
		CPUPercent:    cpuPercent[0],
		MemoryUsed:    vm.Used,
		MemoryTotal:   vm.Total,
		DiskUsed:      diskStat.Used,
		DiskTotal:     diskStat.Total,
		Uptime:        hostStat.Uptime,
		MetricGetTime: utcNow,
		Status:        status,
	}, nil
}
