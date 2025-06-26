package main

import (
	"fmt"
	"time"

	"metrics-agent/internal/config"
	"metrics-agent/internal/metrics"
	"metrics-agent/internal/sender"
)

func main() {
	userID, err := config.LoadUserID()
	if err != nil || userID == "" {
		userID = config.PromptAndSaveUserID()
	}

	for {
		metric, err := metrics.Collect(userID)
		if err == nil {
			sender.SendToAPI(metric)
		} else {
			fmt.Println("❌ Failed to collect metrics:", err)
		}
		time.Sleep(10 * time.Second)
	}
}
