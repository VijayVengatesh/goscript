package influx

import (
	"context"
	"fmt"
	"metrics-agent/internal/config"
	"metrics-agent/internal/models"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var (
	influxURL   = config.LoadConfig().InfluxURL
	influxToken = config.LoadConfig().InfluxToken
	org         = config.LoadConfig().Org
	bucket      = config.LoadConfig().Bucket
)

type InfluxWriteResult struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
	Time        time.Time
	Error       error
}

func SendToInflux(metrics *models.Metrics) {
	// Create InfluxDB client
	client := influxdb2.NewClient(influxURL, influxToken)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)
	fmt.Println("📤 Sending metrics to InfluxDB...", metrics)

	// Create a point
	point := write.NewPoint(
		"vm_metrics",
		map[string]string{
			"user_id":  metrics.UserID,
			"hostname": metrics.Hostname,
		},
		map[string]interface{}{
			"cpu_percent":     metrics.CPUPercent,
			"memory_used":     metrics.MemoryUsed,
			"memory_total":    metrics.MemoryTotal,
			"disk_used":       metrics.DiskUsed,
			"disk_total":      metrics.DiskTotal,
			"uptime_seconds":  metrics.Uptime,
			"metric_get_time": metrics.MetricGetTime,
			"status":          metrics.Status,
		},
		time.Now(),
	)

	// Write the point
	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		fmt.Println("❌ Failed to write to InfluxDB:", err)
	} else {
		fmt.Println("✅ Metrics successfully written to InfluxDB")

		// Convert tags to map for clean printing
		tags := map[string]string{}
		for _, tag := range point.TagList() {
			tags[tag.Key] = tag.Value
		}

		// Convert fields to map for clean printing
		fields := map[string]interface{}{}
		for _, field := range point.FieldList() {
			fields[field.Key] = field.Value
		}

		fmt.Printf("📦 Written data:\n  Measurement: %s\n  Tags: %v\n  Fields: %v\n  Time: %s\n",
			point.Name(), tags, fields, point.Time().Format(time.RFC3339),
		)
	}
}

func SendHealthReportToInflux(report *models.HealthReport) {
	fmt.Println("Sending health report to InfluxDB...", report)
	cfg := config.LoadConfig()

	client := influxdb2.NewClient(cfg.InfluxURL, cfg.InfluxToken)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)

	// Write data to InfluxDB with user_id and hostname as tags
	point := write.NewPoint(
		"vm_health_report",
		map[string]string{
			"user_id":  report.UserID,
			"hostname": report.Hostname,
		},
		map[string]interface{}{
			"cpu":          report.CPUPercent,
			"memory":       report.MemoryPercent,
			"disk":         report.DiskPercent,
			"downtimes":    report.Downtimes,
			"availability": report.Availability,
			"sla_achieved": report.SLA,
		},
		time.Now(),
	)

	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		fmt.Println("❌ Failed to write health report to InfluxDB:", err)
	} else {
		fmt.Println("✅ Health report sent to InfluxDB successfully")
	}
}

func SendSystemSummaryReportToInflux(report *models.SystemSummary) {
	fmt.Println("Sending system summary report to InfluxDB...", report)
	cfg := config.LoadConfig()

	client := influxdb2.NewClient(cfg.InfluxURL, cfg.InfluxToken)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)

	// Write data to InfluxDB with user_id and hostname as tags
	point := write.NewPoint(
		"system_summary",
		map[string]string{
			"user_id":  report.UserID,
			"hostname": report.Hostname,
		},
		map[string]any{
			"ip_address":           report.IPAddress,
			"os":                   report.OS,
			"cpu_model":            report.CPUModel,
			"cpu_cores":            report.CPUCores,
			"ram_mb":               report.RAMMB,
			"disk_count":           report.DiskCount,
			"sys_logs_error_count": report.SysLogsErrorCount,
			"login_count":          report.LoginCount,
			"open_port_count":      report.OpenPortCount,
			"uptime":               report.Uptime,
			"boot_time":            report.BootTime,
			"total_processes":      report.TotalProcesses,
			"nic_count":            report.NICCount,
			"current_user":         report.CurrentUser,
		},
		time.Now(),
	)

	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		fmt.Println(" Failed to write System Summary to InfluxDB:", err)
	} else {
		fmt.Println(" System Summary sent to InfluxDB successfully")
	}
}

func SendLoadAverageToInflux(report *models.LoadAverageMetrics) {
	fmt.Println("Sending load average report to InfluxDB...", report)
	cfg := config.LoadConfig()

	client := influxdb2.NewClient(cfg.InfluxURL, cfg.InfluxToken)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(cfg.Org, cfg.Bucket)

	// Write data to InfluxDB with user_id and hostname as tags
	point := write.NewPoint(
		"load_average",
		map[string]string{
			"user_id":  report.UserID,
			"hostname": report.Hostname,
		},
		map[string]any{
			"load_1m":      report.Load1m,
			"load_5m":      report.Load5m,
			"load_15m":     report.Load15m,
			"load_1m_min":  report.Load1mMin,
			"load_1m_max":  report.Load1mMax,
			"load_5m_min":  report.Load5mMin,
			"load_5m_max":  report.Load5mMax,
			"load_15m_min": report.Load15mMin,
			"load_15m_max": report.Load15mMax,
		},
		time.Now(),
	)

	err := writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		fmt.Println(" Failed to write Load Average to InfluxDB:", err)
	} else {
		fmt.Println(" Load Average sent to InfluxDB successfully")
	}
}
