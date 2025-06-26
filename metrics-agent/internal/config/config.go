package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"metrics-agent/internal/models"
)

func getConfigPath() string {
	var baseDir string
	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
	default:
		baseDir = "/etc"
	}
	return filepath.Join(baseDir, "metrics-agent", "config.json")
}

func LoadUserID() (string, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cfg models.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}
	return cfg.UserID, nil
}

func PromptAndSaveUserID() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your User ID (or Device Key): ")
	userID, _ := reader.ReadString('\n')
	userID = string(bytes.TrimSpace([]byte(userID)))

	config := models.Config{UserID: userID}
	data, _ := json.MarshalIndent(config, "", "  ")

	os.MkdirAll(filepath.Dir(getConfigPath()), 0700)
	_ = os.WriteFile(getConfigPath(), data, 0600)

	fmt.Println("✔ User ID stored successfully")
	return userID
}
