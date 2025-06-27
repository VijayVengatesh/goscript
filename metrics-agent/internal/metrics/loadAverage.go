package metrics

import (
	"fmt"
	"metrics-agent/internal/models"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

// GetLoadAverage returns system load average for Linux/macOS and CPU usage average for Windows.
func GetLoadAverage(UserId string) (*models.LoadAverageMetrics, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	if runtime.GOOS == "windows" {
		// Collect multiple samples to estimate min/max for Windows
		var load1m, load5m, load15m []float64

		for i := 0; i < 3; i++ {
			l, err := cpu.Percent(2*time.Second, false)
			if err != nil {
				return nil, err
			}
			load1m = append(load1m, l[0])
			time.Sleep(1 * time.Second)
		}
		for i := 0; i < 3; i++ {
			l, err := cpu.Percent(2*time.Second, false)
			if err != nil {
				return nil, err
			}
			load5m = append(load5m, l[0])
			time.Sleep(1 * time.Second)
		}
		for i := 0; i < 3; i++ {
			l, err := cpu.Percent(2*time.Second, false)
			if err != nil {
				return nil, err
			}
			load15m = append(load15m, l[0])
			time.Sleep(1 * time.Second)
		}

		return &models.LoadAverageMetrics{
			UserID:     UserId,
			Hostname:   hostInfo.Hostname,
			Load1m:     average(load1m),
			Load1mMin:  min(load1m),
			Load1mMax:  max(load1m),
			Load5m:     average(load5m),
			Load5mMin:  min(load5m),
			Load5mMax:  max(load5m),
			Load15m:    average(load15m),
			Load15mMin: min(load15m),
			Load15mMax: max(load15m),
		}, nil

	} else {
		avg, err := load.Avg()
		if err != nil {
			return nil, fmt.Errorf("failed to get load average: %v", err)
		}

		// On Unix, there's no direct min/max from OS; send same as average
		return &models.LoadAverageMetrics{
			UserID:     UserId,
			Hostname:   hostInfo.Hostname,
			Load1m:     round(avg.Load1),
			Load1mMin:  round(avg.Load1),
			Load1mMax:  round(avg.Load1),
			Load5m:     round(avg.Load5),
			Load5mMin:  round(avg.Load5),
			Load5mMax:  round(avg.Load5),
			Load15m:    round(avg.Load15),
			Load15mMin: round(avg.Load15),
			Load15mMax: round(avg.Load15),
		}, nil
	}
}
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return round(sum / float64(len(values)))
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return round(min)
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return round(max)
}
