package main

import (
	"fmt"
	"metrics-agent/internal/influx"
	"time"

	"metrics-agent/internal/config"
	"metrics-agent/internal/metrics"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println("✅ Running in:", cfg.Env)
	fmt.Println("📡 API Endpoint:", cfg.APIEndpoint)
	fmt.Println("📊 InfluxDB URL:", cfg.InfluxURL)

	userID, err := config.LoadUserID()
	if err != nil || userID == "" {
		userID = config.PromptAndSaveUserID()
	}

	metricTicker := time.NewTicker(1000 * time.Second)
	healthTicker := time.NewTicker(1000 * time.Second)
	systemSummeryTicker := time.NewTicker(10 * time.Second)
	loadAverageTicker := time.NewTicker(10 * time.Second)
	// Stop tickers on exit
	defer metricTicker.Stop()
	defer healthTicker.Stop()

	// Run metric collection every 10s
	go func() {
		for range metricTicker.C {
			metric, err := metrics.Collect(userID)
			if err != nil {
				fmt.Println("❌ Failed to collect metrics:", err)
				continue
			}
			fmt.Println("Collected metrics:", metric)
			influx.SendToInflux(metric)
		}
	}()

	// Run health report generation every 10s
	go func() {
		for range healthTicker.C {
			healthReport, err := metrics.GenerateHealthReport(userID)
			if err != nil {
				fmt.Println("❌ Failed to generate health report:", err)
				continue
			}
			fmt.Println("Generated health report:", healthReport)
			influx.SendHealthReportToInflux(healthReport)
		}
	}()
	go func() {
		for range systemSummeryTicker.C {
			systemSummary, err := metrics.GetSystemSummary(userID)
			if err != nil {
				fmt.Println(" Failed to get system summary:", err)
				continue
			}
			fmt.Println("Collected system summary:", systemSummary)
			influx.SendSystemSummaryReportToInflux(systemSummary)
		}
	}()
	go func() {
		for range loadAverageTicker.C {
			loadAverage, err := metrics.GetLoadAverage(userID)
			if err != nil {
				fmt.Println(" Failed to get system summary:", err)
				continue
			}
			fmt.Println("Collected load average:", loadAverage)
			influx.SendLoadAverageToInflux(loadAverage)
		}
	}()
	// Prevent main() from exiting
	select {}
}
