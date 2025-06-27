package metrics

import (
	"fmt"
	"runtime"
	"time"

	"metrics-agent/internal/models"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	totalChecks    int
	totalDowntimes int
	slaSuccess     int
)

func GenerateHealthReport(userId string) (*models.HealthReport, error) {
	totalChecks++

	// CPU usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	// Memory usage
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %v", err)
	}

	// Disk usage
	path := "/"
	if runtime.GOOS == "windows" {
		path = "C:\\"
	}
	diskStat, err := disk.Usage(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage: %v", err)
	}

	// Host uptime and hostname
	hostStat, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}

	uptime := time.Duration(hostStat.Uptime) * time.Second
	downtime := false

	// Threshold checks (adjustable if needed)
	const threshold = 95.0
	if cpuPercent[0] >= threshold || vm.UsedPercent >= threshold || diskStat.UsedPercent >= threshold {
		downtime = true
	}

	if uptime < 24*time.Hour {
		downtime = true
	}

	if downtime {
		totalDowntimes++
	} else {
		slaSuccess++
	}

	availability := (float64(totalChecks-totalDowntimes) / float64(totalChecks)) * 100
	sla := (float64(slaSuccess) / float64(totalChecks)) * 100

	report := &models.HealthReport{
		UserID:        userId,
		Hostname:      hostStat.Hostname,
		Availability:  fmt.Sprintf("%.1f %%", availability),
		CPUPercent:    round(cpuPercent[0]),
		MemoryPercent: round(vm.UsedPercent),
		DiskPercent:   round(diskStat.UsedPercent),
		Downtimes:     totalDowntimes,
		SLA:           fmt.Sprintf("%.2f %%", sla),
	}

	return report, nil
}

func round(value float64) float64 {
	return float64(int(value*10)) / 10 // 71.79 -> 71.7
}
