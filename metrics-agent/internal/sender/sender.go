package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"metrics-agent/internal/models"
)

func SendToAPI(metrics *models.Metrics) {
	jsonData, _ := json.MarshalIndent(metrics, "", "  ")
	fmt.Println("📤 Sending payload:\n", string(jsonData))

	resp, err := http.Post("https://cloudops-api.idevopz.com/metrics", "application/json", bytes.NewBuffer(jsonData))
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
