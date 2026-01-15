package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/model"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [source]",
	Short: "Encrypt files or directories with access controls",
	Long: `Encrypt files or directories using AES-256-GCM with multiple access methods.

The encrypt command creates encrypted artifacts that can only be decrypted
with the specified access methods and policies.

Examples:
  # Encrypt with passphrase
  vault encrypt ./config.tar.gz --access-method passphrase

  # Encrypt with multiple access methods
  vault encrypt ./secrets.tar.gz \
    --access-method passphrase \
    --access-method runtime \
    --policy ttl=24h

  # Encrypt with compression
  vault encrypt ./backup.tar.gz --compress --output backup.ava

  # Encrypt directory
  vault encrypt ./data --compress --output data.ava`,
	RunE: runEncryptCommand,
}

var (
	encryptOutputPath    string
	encryptDescription   string
	encryptCompress      bool
	encryptAccessMethods []string
	encryptPolicies      []string
	encryptPassphrase    bool
	encryptRuntime       bool
	encryptCertificate   string
	encryptTTL           string
	encryptEnvironment   string
	encryptInstance      string
	encryptRegion        string
)

func init() {
	// Command will be added by root.go

	// File options
	encryptCmd.Flags().StringVarP(&encryptOutputPath, "output", "o", "", "Output file path (default: source.ava)")
	encryptCmd.Flags().StringVar(&encryptDescription, "description", "", "Description for the encrypted artifact")
	encryptCmd.Flags().BoolVar(&encryptCompress, "compress", true, "Compress content before encryption")

	// Access methods
	encryptCmd.Flags().StringSliceVar(&encryptAccessMethods, "access-method", []string{}, "Access methods (passphrase, runtime, certificate:file)")
	encryptCmd.Flags().BoolVar(&encryptPassphrase, "passphrase", false, "Add passphrase access method")
	encryptCmd.Flags().BoolVar(&encryptRuntime, "runtime", false, "Add runtime access method")
	encryptCmd.Flags().StringVar(&encryptCertificate, "certificate", "", "Certificate file for certificate-based access")

	// Policy options
	encryptCmd.Flags().StringSliceVar(&encryptPolicies, "policy", []string{}, "Policies (ttl=24h, environment=prod, instance=id, region=name)")
	encryptCmd.Flags().StringVar(&encryptTTL, "ttl", "", "Time-to-live policy")
	encryptCmd.Flags().StringVar(&encryptEnvironment, "environment", "", "Environment policy")
	encryptCmd.Flags().StringVar(&encryptInstance, "instance", "", "Instance policy")
	encryptCmd.Flags().StringVar(&encryptRegion, "region", "", "Region policy")
}

func runEncryptCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("source path is required")
	}

	sourcePath := args[0]

	// Validate source path
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}

	// Set default output path
	if encryptOutputPath == "" {
		encryptOutputPath = sourcePath + ".ava"
	}

	// Create execution context
	_, err := context.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// Build encryption request
	req, err := buildEncryptionRequest(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to build encryption request: %w", err)
	}

	// Get user ID from context (for now, use a fixed UUID)
	userID := uuid.New()

	// Create encryption service with secure configuration
	masterKey, kdfSalt, kdfIterations, err := getEncryptionConfig()
	if err != nil {
		return fmt.Errorf("failed to get encryption configuration: %w", err)
	}

	auditService := services.NewAuditService()
	encryptionService := services.NewEncryptionService(masterKey, kdfSalt, kdfIterations, auditService)

	// Perform encryption
	result, err := encryptionService.Encrypt(req, userID)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Display result
	displayEncryptionResult(result)

	return nil
}

func buildEncryptionRequest(sourcePath string) (*model.EncryptionRequest, error) {
	// Build access methods
	accessMethods, err := buildAccessMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to build access methods: %w", err)
	}

	// Build policies
	policies, err := buildPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to build policies: %w", err)
	}

	// Create request
	req := &model.EncryptionRequest{
		SourcePath:    sourcePath,
		OutputPath:    encryptOutputPath,
		AccessMethods: accessMethods,
		Policies:      policies,
		Description:   encryptDescription,
		Compression:   encryptCompress,
	}

	return req, nil
}

func buildAccessMethods() ([]model.AccessMethodConfig, error) {
	var methods []model.AccessMethodConfig

	// Add passphrase method if requested
	if encryptPassphrase {
		methods = append(methods, model.AccessMethodConfig{
			Type: model.AccessMethodTypePassphrase,
			Name: "passphrase",
			Config: map[string]interface{}{
				"iterations": 100000,
				"salt":       "", // Will be generated
			},
		})
	}

	// Add runtime method if requested
	if encryptRuntime {
		methods = append(methods, model.AccessMethodConfig{
			Type: model.AccessMethodTypeRuntime,
			Name: "runtime",
			Config: map[string]interface{}{
				"instance": encryptInstance,
				"region":   encryptRegion,
			},
		})
	}

	// Add certificate method if specified
	if encryptCertificate != "" {
		certInfo, err := getCertificateInfo(encryptCertificate)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}

		methods = append(methods, model.AccessMethodConfig{
			Type: model.AccessMethodTypeCertificate,
			Name: "certificate-" + filepath.Base(encryptCertificate),
			Config: map[string]interface{}{
				"certificate_file": encryptCertificate,
				"key_id":           certInfo.KeyID,
			},
		})
	}

	// Add methods from --access-method flag
	for _, methodStr := range encryptAccessMethods {
		method, err := parseAccessMethodString(methodStr)
		if err != nil {
			return nil, fmt.Errorf("invalid access method '%s': %w", methodStr, err)
		}
		methods = append(methods, method)
	}

	// Ensure at least one access method
	if len(methods) == 0 {
		return nil, fmt.Errorf("at least one access method is required")
	}

	return methods, nil
}

func buildPolicies() ([]model.EncryptionPolicyConfig, error) {
	var policies []model.EncryptionPolicyConfig

	// Add TTL policy if specified
	if encryptTTL != "" {
		policies = append(policies, model.EncryptionPolicyConfig{
			Type: model.PolicyTypeTTL,
			Name: "ttl",
			Rules: map[string]interface{}{
				"duration": encryptTTL,
			},
		})
	}

	// Add environment policy if specified
	if encryptEnvironment != "" {
		policies = append(policies, model.EncryptionPolicyConfig{
			Type: model.PolicyTypeEnvironment,
			Name: "environment",
			Rules: map[string]interface{}{
				"environment": encryptEnvironment,
			},
		})
	}

	// Add instance policy if specified
	if encryptInstance != "" {
		policies = append(policies, model.EncryptionPolicyConfig{
			Type: model.PolicyTypeInstance,
			Name: "instance",
			Rules: map[string]interface{}{
				"instance": encryptInstance,
			},
		})
	}

	// Add region policy if specified
	if encryptRegion != "" {
		policies = append(policies, model.EncryptionPolicyConfig{
			Type: model.PolicyTypeRegion,
			Name: "region",
			Rules: map[string]interface{}{
				"region": encryptRegion,
			},
		})
	}

	// Add policies from --policy flag
	for _, policyStr := range encryptPolicies {
		policy, err := parsePolicyString(policyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid policy '%s': %w", policyStr, err)
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

func parsePolicyString(policyStr string) (model.EncryptionPolicyConfig, error) {
	parts := strings.SplitN(policyStr, "=", 2)
	if len(parts) != 2 {
		return model.EncryptionPolicyConfig{}, fmt.Errorf("policy must be in format 'type=value'")
	}

	policyType := model.PolicyType(parts[0])
	value := parts[1]

	var rules map[string]interface{}
	var name string

	switch policyType {
	case model.PolicyTypeTTL:
		rules = map[string]interface{}{
			"duration": value,
		}
		name = "ttl"
	case model.PolicyTypeEnvironment:
		rules = map[string]interface{}{
			"environment": value,
		}
		name = "environment"
	case model.PolicyTypeInstance:
		rules = map[string]interface{}{
			"instance": value,
		}
		name = "instance"
	case model.PolicyTypeRegion:
		rules = map[string]interface{}{
			"region": value,
		}
		name = "region"
	default:
		return model.EncryptionPolicyConfig{}, fmt.Errorf("unsupported policy type: %s", policyType)
	}

	return model.EncryptionPolicyConfig{
		Type:  policyType,
		Name:  name,
		Rules: rules,
	}, nil
}

func displayEncryptionResult(result *model.EncryptionResult) {
	fmt.Println("✅ Encryption completed successfully")
	fmt.Println()

	// Display artifact info
	fmt.Printf("📦 Artifact ID: %s\n", result.ArtifactID)
	fmt.Printf("📁 File: %s\n", result.FilePath)
	fmt.Printf("🔐 Algorithm: %s\n", result.Algorithm)
	fmt.Printf("📊 Original size: %d bytes\n", result.OriginalSize)
	fmt.Printf("📊 Encrypted size: %d bytes\n", result.EncryptedSize)

	if result.OriginalSize > 0 {
		ratio := float64(result.EncryptedSize) / float64(result.OriginalSize) * 100
		fmt.Printf("📏 Size ratio: %.1f%%\n", ratio)
	}

	// Display access methods
	if len(result.AccessMethods) > 0 {
		fmt.Println()
		fmt.Println("🔑 Access methods:")
		for _, method := range result.AccessMethods {
			fmt.Printf("  • %s\n", method)
		}
	}

	// Display policies
	if len(result.Policies) > 0 {
		fmt.Println()
		fmt.Println("📋 Policies:")
		for _, policy := range result.Policies {
			fmt.Printf("  • %s\n", policy)
		}
	}

	// Display creation time
	fmt.Println()
	fmt.Printf("⏰ Created: %s\n", result.CreatedAt.Format("2006-01-02 15:04:05"))
}

// NewEncryptCommand creates the encrypt command
func NewEncryptCommand() *cobra.Command {
	return encryptCmd
}
