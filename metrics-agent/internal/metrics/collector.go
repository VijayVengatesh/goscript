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
	}, nil
}
