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

type HealthReport struct {
	UserID        string  `json:"user_id"`
	Hostname      string  `json:"hostname"`
	Availability  string  `json:"availability"`
	CPUPercent    float64 `json:"cpu"`
	MemoryPercent float64 `json:"memory"`
	DiskPercent   float64 `json:"disk"`
	Downtimes     int     `json:"downtimes"`
	SLA           string  `json:"sla_achieved"`
}

type SystemSummary struct {
	UserID            string  `json:"user_id"`
	Hostname          string  `json:"hostname"`
	IPAddress         string  `json:"ip_address"`
	OS                string  `json:"os"`
	CPUModel          string  `json:"cpu_model"`
	CPUCores          int     `json:"cpu_cores"`
	RAMMB             float64 `json:"ram_mb"`
	DiskCount         int     `json:"disk_count"`
	SysLogsErrorCount int     `json:"sys_logs_errors"`
	Uptime            string  `json:"uptime"`
	BootTime          uint64  `json:"boot_time"`
	TotalProcesses    int     `json:"total_processes"`
	NICCount          int     `json:"nic_count"`
	LoginCount        int     `json:"login_count"`
	OpenPortCount     int     `json:"open_port_count"`
	CurrentUser       string  `json:"current_user"`
}

// LoadAverageMetrics holds system load metrics.
type LoadAverageMetrics struct {
	UserID   string  `json:"user_id"`
	Hostname string  `json:"hostname"`
	Load1m   float64 `json:"load_1m"`
	Load5m   float64 `json:"load_5m"`
	Load15m  float64 `json:"load_15m"`

	Load1mMin float64
	Load1mMax float64

	Load5mMin float64
	Load5mMax float64

	Load15mMin float64
	Load15mMax float64
}
