package utils

import (
	"net/http"
	"time"
)

// IsServerUp checks if a server is up (returns "up" or "down").
func IsServerUp(url string) string {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "down"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "up"
	}
	return "down"
}
