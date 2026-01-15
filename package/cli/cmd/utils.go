package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/services"
)

// parseAccessMethodString parses an access method string in format "type[:config]"
func parseAccessMethodString(methodStr string) (services.AccessMethodConfig, error) {
	parts := strings.SplitN(methodStr, ":", 2)
	if len(parts) == 0 {
		return services.AccessMethodConfig{}, fmt.Errorf("invalid method format")
	}

	methodType := services.AccessMethodType(parts[0])
	var config map[string]interface{}

	switch methodType {
	case services.AccessMethodTypePassphrase:
		config = map[string]interface{}{
			"iterations": 100000,
		}
	case services.AccessMethodTypeRuntime:
		config = map[string]interface{}{}
	case services.AccessMethodTypeCertificate:
		if len(parts) < 2 {
			return services.AccessMethodConfig{}, fmt.Errorf("certificate method requires certificate file")
		}
		certFile := parts[1]
		certInfo, err := getCertificateInfo(certFile)
		if err != nil {
			return services.AccessMethodConfig{}, fmt.Errorf("failed to load certificate: %w", err)
		}
		config = map[string]interface{}{
			"certificate_file": certFile,
			"key_id":           certInfo.KeyID,
		}
	default:
		return services.AccessMethodConfig{}, fmt.Errorf("unsupported access method type: %s", methodType)
	}

	return services.AccessMethodConfig{
		Type:   methodType,
		Name:   string(methodType),
		Config: config,
	}, nil
}

// CertificateInfo holds certificate information
type CertificateInfo struct {
	KeyID    string
	Subject  string
	Issuer   string
	NotAfter string
}

// getCertificateInfo extracts information from a certificate file
func getCertificateInfo(certFile string) (*CertificateInfo, error) {
	// Read certificate file
	certData, err := os.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Parse certificate
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Generate key ID from certificate public key
	keyID := fmt.Sprintf("%x", sha256.Sum256(cert.RawSubjectPublicKeyInfo))

	return &CertificateInfo{
		KeyID:    keyID[:16], // Use first 16 characters as ID
		Subject:  cert.Subject.String(),
		Issuer:   cert.Issuer.String(),
		NotAfter: cert.NotAfter.Format("2006-01-02"),
	}, nil
}

// getEncryptionConfig retrieves encryption configuration from environment or config file
func getEncryptionConfig() (masterKey, kdfSalt string, kdfIterations int, err error) {
	// Try environment variables first
	if masterKey = os.Getenv("AETHER_MASTER_KEY"); masterKey != "" {
		if kdfSalt = os.Getenv("AETHER_KDF_SALT"); kdfSalt != "" {
			if iterStr := os.Getenv("AETHER_KDF_ITERATIONS"); iterStr != "" {
				if iterations, err := strconv.Atoi(iterStr); err == nil {
					return masterKey, kdfSalt, iterations, nil
				}
			}
		}
	}

	// Fallback to config file
	configPath := filepath.Join(os.Getenv("HOME"), ".aether", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		// TODO: Implement YAML config parsing
		// For now, generate secure defaults
		return generateSecureConfig()
	}

	// Generate secure defaults
	return generateSecureConfig()
}

// generateSecureConfig generates secure encryption parameters
func generateSecureConfig() (masterKey, kdfSalt string, kdfIterations int, err error) {
	// Generate random master key
	masterKeyBytes := make([]byte, 32)
	if _, err := rand.Read(masterKeyBytes); err != nil {
		return "", "", 0, fmt.Errorf("failed to generate master key: %w", err)
	}
	masterKey = base64.StdEncoding.EncodeToString(masterKeyBytes)

	// Generate random salt
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", "", 0, fmt.Errorf("failed to generate salt: %w", err)
	}
	kdfSalt = base64.StdEncoding.EncodeToString(saltBytes)

	kdfIterations = 100000 // PBKDF2 recommended iterations

	// Save to config file for persistence
	configDir := filepath.Join(os.Getenv("HOME"), ".aether")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		// Continue even if we can't save
		return masterKey, kdfSalt, kdfIterations, nil
	}

	configFile := filepath.Join(configDir, "encryption_key.yaml")
	configContent := fmt.Sprintf(`master_key: %s
kdf_salt: %s
kdf_iterations: %d
created: %s
`, masterKey, kdfSalt, kdfIterations, time.Now().Format(time.RFC3339))

	if err := os.WriteFile(configFile, []byte(configContent), 0600); err != nil {
		// Continue even if we can't save
		return masterKey, kdfSalt, kdfIterations, nil
	}

	return masterKey, kdfSalt, kdfIterations, nil
}
