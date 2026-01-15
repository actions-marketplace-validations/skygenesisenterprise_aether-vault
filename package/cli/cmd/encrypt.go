package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/services"
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
	encryptOutputPath     string
	encryptDescription    string
	encryptCompress       bool
	encryptRemoveOriginal bool
	encryptAccessMethods  []string
	encryptPolicies       []string
	encryptPassphrase     bool
	encryptRuntime        bool
	encryptCertificate    string
	encryptTTL            string
	encryptEnvironment    string
	encryptInstance       string
	encryptRegion         string
)

func init() {
	// Command will be added by root.go

	// File options
	encryptCmd.Flags().StringVarP(&encryptOutputPath, "output", "o", "", "Output file path (default: source.ava)")
	encryptCmd.Flags().StringVar(&encryptDescription, "description", "", "Description for the encrypted artifact")
	encryptCmd.Flags().BoolVar(&encryptCompress, "compress", true, "Compress content before encryption")
	encryptCmd.Flags().BoolVar(&encryptRemoveOriginal, "remove-original", true, "Remove original files/folders after successful encryption")

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

	// Check if we should use password protection instead of direct .ava conversion
	if encryptPassphrase || containsPasswordMethod(encryptAccessMethods) {
		return runPasswordProtectionFlow(sourcePath)
	}

	// Set default output path for normal encryption
	if encryptOutputPath == "" {
		encryptOutputPath = sourcePath + ".ava"
	}

	// Create execution context without server (for encryption operations)
	_, err := context.NewWithoutServer(nil)
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

	encryptionService := services.NewEncryptionService(masterKey, kdfSalt, kdfIterations)

	// Perform encryption
	result, err := encryptionService.Encrypt(req, userID)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Display result
	displayEncryptionResult(result)

	// Restrict access to original files if requested
	if encryptRemoveOriginal {
		fmt.Printf("🔒 Restricting access to original files/directories...\n")
		if err := restrictOriginalFiles(sourcePath); err != nil {
			fmt.Printf("⚠️  Warning: Failed to restrict original files: %v\n", err)
		} else {
			fmt.Printf("✅ Original files access restricted - only decryption can restore them\n")
		}
	}

	return nil
}

func buildEncryptionRequest(sourcePath string) (*services.EncryptionRequest, error) {
	// Build access methods
	accessMethods, err := buildAccessMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to build access methods: %w", err)
	}

	// For passphrase access method, prompt for password interactively
	if encryptPassphrase {
		password, err := promptForPassword()
		if err != nil {
			return nil, fmt.Errorf("failed to read password: %w", err)
		}

		// Update access method config with the password
		for i := range accessMethods {
			if accessMethods[i].Type == services.AccessMethodTypePassphrase {
				if accessMethods[i].Config == nil {
					accessMethods[i].Config = make(map[string]interface{})
				}
				accessMethods[i].Config["password"] = password
				break
			}
		}
	}

	// Build policies
	policies, err := buildPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to build policies: %w", err)
	}

	// Create request
	req := &services.EncryptionRequest{
		SourcePath:    sourcePath,
		OutputPath:    encryptOutputPath,
		AccessMethods: accessMethods,
		Policies:      policies,
		Description:   encryptDescription,
		Compression:   encryptCompress,
	}

	return req, nil
}

func buildAccessMethods() ([]services.AccessMethodConfig, error) {
	var methods []services.AccessMethodConfig

	// Add passphrase method if requested
	if encryptPassphrase {
		methods = append(methods, services.AccessMethodConfig{
			Type: services.AccessMethodTypePassphrase,
			Name: "passphrase",
			Config: map[string]interface{}{
				"iterations": 100000,
				"salt":       "", // Will be generated
			},
		})
	}

	// Add runtime method if requested
	if encryptRuntime {
		methods = append(methods, services.AccessMethodConfig{
			Type: services.AccessMethodTypeRuntime,
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

		methods = append(methods, services.AccessMethodConfig{
			Type: services.AccessMethodTypeCertificate,
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

func buildPolicies() ([]services.EncryptionPolicyConfig, error) {
	var policies []services.EncryptionPolicyConfig

	// Add TTL policy if specified
	if encryptTTL != "" {
		policies = append(policies, services.EncryptionPolicyConfig{
			Type: services.PolicyTypeTTL,
			Name: "ttl",
			Rules: map[string]interface{}{
				"duration": encryptTTL,
			},
		})
	}

	// Add environment policy if specified
	if encryptEnvironment != "" {
		policies = append(policies, services.EncryptionPolicyConfig{
			Type: services.PolicyTypeEnvironment,
			Name: "environment",
			Rules: map[string]interface{}{
				"environment": encryptEnvironment,
			},
		})
	}

	// Add instance policy if specified
	if encryptInstance != "" {
		policies = append(policies, services.EncryptionPolicyConfig{
			Type: services.PolicyTypeInstance,
			Name: "instance",
			Rules: map[string]interface{}{
				"instance": encryptInstance,
			},
		})
	}

	// Add region policy if specified
	if encryptRegion != "" {
		policies = append(policies, services.EncryptionPolicyConfig{
			Type: services.PolicyTypeRegion,
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

func parsePolicyString(policyStr string) (services.EncryptionPolicyConfig, error) {
	parts := strings.SplitN(policyStr, "=", 2)
	if len(parts) != 2 {
		return services.EncryptionPolicyConfig{}, fmt.Errorf("policy must be in format 'type=value'")
	}

	policyType := services.PolicyType(parts[0])
	value := parts[1]

	var rules map[string]interface{}
	var name string

	switch policyType {
	case services.PolicyTypeTTL:
		rules = map[string]interface{}{
			"duration": value,
		}
		name = "ttl"
	case services.PolicyTypeEnvironment:
		rules = map[string]interface{}{
			"environment": value,
		}
		name = "environment"
	case services.PolicyTypeInstance:
		rules = map[string]interface{}{
			"instance": value,
		}
		name = "instance"
	case services.PolicyTypeRegion:
		rules = map[string]interface{}{
			"region": value,
		}
		name = "region"
	default:
		return services.EncryptionPolicyConfig{}, fmt.Errorf("unsupported policy type: %s", policyType)
	}

	return services.EncryptionPolicyConfig{
		Type:  policyType,
		Name:  name,
		Rules: rules,
	}, nil
}

func displayEncryptionResult(result *services.EncryptionResult) {
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

// promptForPassword prompts user for a password securely
func promptForPassword() (string, error) {
	fmt.Print("🔐 Enter password: ")

	// Try to read password securely first
	password, err := readPassword()
	if err != nil {
		// Fallback to standard input if terminal not available
		fmt.Println("\n⚠️  Secure input not available, using standard input")
		fmt.Print("🔐 Enter password (will be visible): ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil && err.Error() != "unexpected newline" {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		password = input
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Confirm password for security
	fmt.Print("🔐 Confirm password: ")
	confirmPassword, err := readPassword()
	if err != nil {
		// Fallback to standard input if terminal not available
		fmt.Println("\n⚠️  Secure input not available, using standard input")
		fmt.Print("🔐 Confirm password (will be visible): ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil && err.Error() != "unexpected newline" {
			return "", fmt.Errorf("failed to read confirmation: %w", err)
		}
		confirmPassword = input
	}

	if password != confirmPassword {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}

// restrictOriginalFiles restricts access to original files/folders after encryption
func restrictOriginalFiles(sourcePath string) error {
	// Check if source is a directory or file
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}

	if fileInfo.IsDir() {
		// Restrict directory and all contents
		return restrictDirectory(sourcePath)
	} else {
		// Restrict single file
		return restrictFile(sourcePath)
	}
}

// restrictDirectory restricts access to a directory and all its contents
func restrictDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the directory itself for now, handle it after walking
		if path == dirPath && info.IsDir() {
			return nil
		}

		return restrictFile(path)
	})

	return nil
}

// containsPasswordMethod checks if any access method is password/passphrase related
func containsPasswordMethod(methods []string) bool {
	for _, method := range methods {
		if method == "passphrase" || method == "password" {
			return true
		}
	}
	return false
}

// runPasswordProtectionFlow handles the password protection workflow
func runPasswordProtectionFlow(sourcePath string) error {
	fmt.Printf("🔐 Password Protection Mode\n")
	fmt.Printf("📁 Resource: %s\n", sourcePath)
	fmt.Println()

	// Prompt for password to assign to the resource
	password, err := promptForResourcePassword()
	if err != nil {
		return fmt.Errorf("password prompt failed: %w", err)
	}

	// Create execution context without server
	_, err = context.NewWithoutServer(nil)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// Set output path for encrypted file
	if encryptOutputPath == "" {
		encryptOutputPath = sourcePath + ".ava"
	}

	// Build encryption request with password
	req, err := buildPasswordProtectionRequest(sourcePath, password)
	if err != nil {
		return fmt.Errorf("failed to build encryption request: %w", err)
	}

	// Get user ID
	userID := uuid.New()

	// Create encryption service
	masterKey, kdfSalt, kdfIterations, err := getEncryptionConfig()
	if err != nil {
		return fmt.Errorf("failed to get encryption configuration: %w", err)
	}

	encryptionService := services.NewEncryptionService(masterKey, kdfSalt, kdfIterations)

	// Perform encryption
	fmt.Printf("🔒 Encrypting resource with password protection...\n")
	result, err := encryptionService.Encrypt(req, userID)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	// Display result
	displayEncryptionResult(result)

	// Restrict access to original files
	fmt.Printf("🔒 Restricting access to original files/directories...\n")
	if err := restrictOriginalFiles(sourcePath); err != nil {
		fmt.Printf("⚠️  Warning: Failed to restrict original files: %v\n", err)
	} else {
		fmt.Printf("✅ Original files access restricted - password required to access\n")
	}

	// Create password-protected launcher
	fmt.Printf("🚀 Creating password-protected launcher...\n")
	if err := createPasswordLauncher(encryptOutputPath, password); err != nil {
		fmt.Printf("⚠️  Warning: Failed to create launcher: %v\n", err)
	} else {
		fmt.Printf("✅ Password launcher created - double-click to access with password\n")
	}

	return nil
}

// promptForResourcePassword prompts for password to assign to the resource
func promptForResourcePassword() (string, error) {
	fmt.Println("🔑 Set password for this resource:")
	fmt.Printf("   This password will be required to access the encrypted content\n")
	fmt.Println()

	maxAttempts := 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			fmt.Printf("\n❌ Passwords didn't match. Attempts remaining: %d\n", maxAttempts-attempt+1)
		}

		fmt.Print("🔐 Enter password to assign: ")
		password, err := readPassword()
		if err != nil {
			// Fallback to standard input if terminal not available
			fmt.Println("\n⚠️  Secure input not available, using standard input")
			fmt.Print("🔐 Enter password to assign (will be visible): ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil && err.Error() != "unexpected newline" {
				return "", fmt.Errorf("failed to read password: %w", err)
			}
			password = input
		}

		if password == "" {
			if attempt < maxAttempts {
				fmt.Printf("\n⚠️  Password cannot be empty. Please try again.\n")
				continue
			}
			return "", fmt.Errorf("password cannot be empty")
		}

		if len(password) < 8 {
			if attempt < maxAttempts {
				fmt.Printf("\n⚠️  Password should be at least 8 characters. Please try again.\n")
				continue
			}
			return "", fmt.Errorf("password should be at least 8 characters")
		}

		fmt.Print("🔐 Confirm password: ")
		confirmPassword, err := readPassword()
		if err != nil {
			// Fallback to standard input if terminal not available
			fmt.Println("\n⚠️  Secure input not available, using standard input")
			fmt.Print("🔐 Confirm password (will be visible): ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil && err.Error() != "unexpected newline" {
				return "", fmt.Errorf("failed to read confirmation: %w", err)
			}
			confirmPassword = input
		}

		if password == confirmPassword {
			fmt.Println()
			fmt.Printf("✅ Password set successfully\n")
			return password, nil
		}
	}

	return "", fmt.Errorf("maximum attempts (%d) reached. Please try again.", maxAttempts)
}

// buildPasswordProtectionRequest builds encryption request with password method
func buildPasswordProtectionRequest(sourcePath, password string) (*services.EncryptionRequest, error) {
	// Create access method config with password
	accessMethods := []services.AccessMethodConfig{
		{
			Type: services.AccessMethodTypePassphrase,
			Name: "password",
			Config: map[string]interface{}{
				"password":   password,
				"iterations": 100000,
			},
		},
	}

	// Create request
	req := &services.EncryptionRequest{
		SourcePath:    sourcePath,
		OutputPath:    encryptOutputPath,
		AccessMethods: accessMethods,
		Policies:      []services.EncryptionPolicyConfig{},
		Description:   encryptDescription,
		Compression:   encryptCompress,
	}

	return req, nil
}

// createPasswordLauncher creates an executable script that prompts for password when double-clicked
func createPasswordLauncher(encryptedFilePath, password string) error {
	// Create launcher script path
	launcherPath := encryptedFilePath + ".sh"

	// Create script content
	scriptContent := fmt.Sprintf(`#!/bin/bash
# Aether Vault Password Launcher
# This script will prompt for password to decrypt the encrypted file

ENCRYPTED_FILE="%s"
OUTPUT_DIR="%s_decrypted"

echo "🔐 Aether Vault - Encrypted File"
echo "📁 File: $ENCRYPTED_FILE"
echo ""

# Prompt for password
echo "🔑 Enter password to decrypt:"
read -s password
echo ""

# Attempt to decrypt using vault CLI
if vault decrypt "$ENCRYPTED_FILE" --passphrase --output "$OUTPUT_DIR" <<< "$password"; then
    echo "✅ Decryption successful!"
    echo "📁 Files are available in: $OUTPUT_DIR"
    
    # Ask if user wants to open the decrypted folder
    read -p "🚀 Open decrypted folder? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if command -v xdg-open > /dev/null 2>&1; then
            xdg-open "$OUTPUT_DIR"
        elif command -v open > /dev/null 2>&1; then
            open "$OUTPUT_DIR"
        elif command -v explorer > /dev/null 2>&1; then
            explorer "$OUTPUT_DIR"
        else
            echo "📂 Decrypted files available at: $OUTPUT_DIR"
        fi
    fi
else
    echo "❌ Decryption failed! Invalid password or corrupted file."
    echo "Please try again with the correct password."
    exit 1
fi
`, encryptedFilePath, strings.TrimSuffix(encryptedFilePath, filepath.Ext(encryptedFilePath)))

	// Write launcher script
	if err := os.WriteFile(launcherPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create launcher script: %w", err)
	}

	// For desktop environments, also create a .desktop file on Linux
	if runtime.GOOS == "linux" {
		desktopPath := strings.TrimSuffix(encryptedFilePath, filepath.Ext(encryptedFilePath)) + ".desktop"
		desktopContent := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=Encrypted File - %s
Comment=Click to decrypt and access encrypted content
Exec=%s
Icon=folder-locked
Terminal=true
Categories=Utility;Security;
Keywords=encrypted,password,decrypt,vault;
`, filepath.Base(encryptedFilePath), launcherPath)

		if err := os.WriteFile(desktopPath, []byte(desktopContent), 0644); err != nil {
			fmt.Printf("⚠️  Warning: Failed to create desktop file: %v\n", err)
		}
	}

	return nil
}

// restrictFile restricts access to a single file by removing read/write permissions
func restrictFile(filePath string) error {
	// Remove read and write permissions for owner, group, and others
	// Only keep execute permission if it was executable
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	var newMode os.FileMode = 0 // No permissions

	// Keep execute permission if it was executable
	if fileInfo.Mode().Perm()&0111 != 0 {
		newMode = 0111 // Execute only
	}

	if err := os.Chmod(filePath, newMode); err != nil {
		return fmt.Errorf("failed to change permissions for %s: %w", filePath, err)
	}

	return nil
}

// NewEncryptCommand creates the encrypt command
func NewEncryptCommand() *cobra.Command {
	return encryptCmd
}
