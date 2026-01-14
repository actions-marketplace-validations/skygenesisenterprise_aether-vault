package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Security    SecurityConfig
	JWT         JWTConfig
	VaultServer VaultServerConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	Environment  string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Host       string
	User       string
	Password   string
	DBName     string
	Port       int
	SSLMode    string
	Type       string // "postgres" or "sqlite"
	SQLitePath string // for SQLite
}

type SecurityConfig struct {
	EncryptionKey string
	KDFIterations int
}

type JWTConfig struct {
	Secret     string
	Expiration int // hours
}

type VaultServerConfig struct {
	Host        string
	Port        int
	DevMode     bool
	LogLevel    string
	ConfigPath  string
	StoragePath string
	RootToken   string
}

func LoadConfig() (*Config, error) {
	// Determine storage path based on dev mode
	devMode := getEnvAsBool("CLI_SERVER_DEV", false)
	var storagePath string
	if devMode {
		home, err := os.UserHomeDir()
		if err != nil {
			storagePath = "./data"
		} else {
			storagePath = filepath.Join(home, ".aether", "vault", "dev")
		}
	} else {
		configPath := getEnv("CLI_SERVER_CONFIG", "")
		if configPath != "" {
			storagePath = filepath.Join(filepath.Dir(configPath), "data")
		} else {
			storagePath = "./data"
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("CLI_SERVER_HOST", "127.0.0.1"),
			Port:         getEnvAsInt("CLI_SERVER_PORT", 8081),
			Environment:  getEnv("CLI_SERVER_ENV", "development"),
			ReadTimeout:  getEnvAsInt("CLI_SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("CLI_SERVER_WRITE_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			Host:       getEnv("CLI_DB_HOST", ""),
			User:       getEnv("CLI_DB_USER", ""),
			Password:   getEnv("CLI_DB_PASSWORD", ""),
			DBName:     getEnv("CLI_DB_NAME", ""),
			Port:       getEnvAsInt("CLI_DB_PORT", 5432),
			SSLMode:    getEnv("CLI_DB_SSLMODE", "disable"),
			Type:       getEnv("CLI_DB_TYPE", "sqlite"),
			SQLitePath: getEnv("CLI_DB_SQLITE_PATH", "./aether-vault-cli.db"),
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("CLI_ENCRYPTION_KEY", "default-encryption-key-change-in-production"),
			KDFIterations: getEnvAsInt("CLI_KDF_ITERATIONS", 100000),
		},
		JWT: JWTConfig{
			Secret:     getEnv("CLI_JWT_SECRET", "default-jwt-secret-change-in-production"),
			Expiration: getEnvAsInt("CLI_JWT_EXPIRATION", 24),
		},
		VaultServer: VaultServerConfig{
			Host:        getEnv("CLI_SERVER_HOST", "127.0.0.1"),
			Port:        getEnvAsInt("CLI_SERVER_PORT", 8081),
			DevMode:     devMode,
			LogLevel:    getEnv("CLI_SERVER_LOG_LEVEL", "info"),
			ConfigPath:  getEnv("CLI_SERVER_CONFIG", ""),
			StoragePath: storagePath,
			RootToken:   getEnv("CLI_SERVER_ROOT_TOKEN", "dev-token"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
