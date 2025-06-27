package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"metrics-agent/internal/config"
	"metrics-agent/internal/models"
)

func SendToMetricsAPI(metrics *models.Metrics) {
	jsonData, _ := json.MarshalIndent(metrics, "", "  ")
	fmt.Println("📤 Sending payload:\n", string(jsonData))

	url := config.LoadConfig().APIEndpoint + "/metrics"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("✅ Metrics sent! Status:", resp.Status)
	} else {
		fmt.Printf("❌ Failed to send metrics.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}

func SendToHealthReportAPI(healthReport *models.HealthReport) {
	jsonData, _ := json.MarshalIndent(healthReport, "", "  ")
	fmt.Println("📤 Sending payload:\n", string(jsonData))
	url := config.LoadConfig().APIEndpoint + "/healthReport"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Error sending healthReport:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("✅ healthReport sent! Status:", resp.Status)
	} else {
		fmt.Printf("❌ Failed to send healthReport.\nStatus: %s\nResponse: %s\n", resp.Status, string(body))
	}
}
