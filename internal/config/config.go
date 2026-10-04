package config

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Port           string
	DBPath         string
	MediaDir       string
	DataDir        string
	SessionSecret  string
	TrustedProxies []string
}

func LoadConfig() (*Config, error) {
	// 1. Attempt to load .env if present
	loadDotEnv(".env")

	port := getEnv("PORT", "4907")
	dataDir := getEnv("DATA_DIR", "./data")
	mediaDir := getEnv("MEDIA_DIR", "./media")
	dbPath := getEnv("DB_PATH", filepath.Join(dataDir, "flan.db"))

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = loadOrGenerateSessionSecret(dataDir)
	}

	trustedProxiesStr := os.Getenv("TRUSTED_PROXIES")
	var trustedProxies []string
	if trustedProxiesStr != "" {
		for _, p := range strings.Split(trustedProxiesStr, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				trustedProxies = append(trustedProxies, trimmed)
			}
		}
	}

	return &Config{
		Port:           port,
		DBPath:         dbPath,
		MediaDir:       mediaDir,
		DataDir:        dataDir,
		SessionSecret:  sessionSecret,
		TrustedProxies: trustedProxies,
	}, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func loadDotEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if _, exists := os.LookupEnv(k); !exists {
				_ = os.Setenv(k, v)
			}
		}
	}
}

func loadOrGenerateSessionSecret(dataDir string) string {
	secretPath := filepath.Join(dataDir, ".session_secret")
	if data, err := os.ReadFile(secretPath); err == nil {
		sec := strings.TrimSpace(string(data))
		if len(sec) >= 32 {
			return sec
		}
	}

	// Generate random 32-byte secret
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback deterministic if rand fails
		return "flan-media-server-default-secret-key-32bytes!!"
	}
	secretHex := hex.EncodeToString(bytes)

	_ = os.MkdirAll(dataDir, 0755)
	_ = os.WriteFile(secretPath, []byte(secretHex), 0600)
	return secretHex
}
