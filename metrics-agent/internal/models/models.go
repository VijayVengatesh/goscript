package models

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
	MetricGetTime string  `json:"metric_get_time"`
}
