package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/model"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/services"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt [artifact]",
	Short: "Decrypt encrypted artifacts",
	Long: `Decrypt encrypted artifacts using the specified access method.

The decrypt command extracts the original content from encrypted artifacts
after validating access methods and policies.

Examples:
  # Decrypt with passphrase
  vault decrypt ./config.tar.gz.ava --access-method passphrase

  # Decrypt with runtime access
  vault decrypt ./secrets.tar.gz.ava --access-method runtime

  # Decrypt with certificate
  vault decrypt ./backup.tar.gz.ava --access-method certificate:cert.pem

  # Decrypt to specific output
  vault decrypt ./data.ava --output ./extracted-data --force`,
	RunE: runDecryptCommand,
}

var (
	decryptOutputPath   string
	decryptForce        bool
	decryptAccessMethod string
	decryptPassphrase   bool
	decryptRuntime      bool
	decryptCertificate  string
)

func init() {
	// Command will be added by root.go

	// Output options
	decryptCmd.Flags().StringVarP(&decryptOutputPath, "output", "o", "", "Output directory path (default: extracted-artifact-name)")
	decryptCmd.Flags().BoolVar(&decryptForce, "force", false, "Force overwrite existing files")

	// Access method options
	decryptCmd.Flags().StringVar(&decryptAccessMethod, "access-method", "", "Access method (passphrase, runtime, certificate:file)")
	decryptCmd.Flags().BoolVar(&decryptPassphrase, "passphrase", false, "Use passphrase access method")
	decryptCmd.Flags().BoolVar(&decryptRuntime, "runtime", false, "Use runtime access method")
	decryptCmd.Flags().StringVar(&decryptCertificate, "certificate", "", "Certificate file for certificate-based access")
}

func runDecryptCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("artifact path is required")
	}

	artifactPath := args[0]

	// Validate artifact path
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return fmt.Errorf("artifact file does not exist: %s", artifactPath)
	}

	// Create execution context
	_, err := context.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// Build decryption request
	req, err := buildDecryptionRequest(artifactPath)
	if err != nil {
		return fmt.Errorf("failed to build decryption request: %w", err)
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

	// Perform decryption
	result, err := encryptionService.Decrypt(req, userID)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	// Display result
	displayDecryptionResult(result)

	return nil
}

func buildDecryptionRequest(artifactPath string) (*model.DecryptionRequest, error) {
	// Build access method
	accessMethod, err := buildAccessMethod()
	if err != nil {
		return nil, fmt.Errorf("failed to build access method: %w", err)
	}

	// Set default output path
	outputPath := decryptOutputPath
	if outputPath == "" {
		// Extract base name without .ava extension
		baseName := strings.TrimSuffix(artifactPath, ".ava")
		outputPath = baseName + "-extracted"
	}

	// Create request
	req := &model.DecryptionRequest{
		ArtifactPath: artifactPath,
		OutputPath:   outputPath,
		AccessMethod: accessMethod,
		Force:        decryptForce,
	}

	return req, nil
}

func buildAccessMethod() (model.AccessMethodConfig, error) {
	// Determine which access method to use
	methodCount := 0
	if decryptPassphrase {
		methodCount++
	}
	if decryptRuntime {
		methodCount++
	}
	if decryptCertificate != "" {
		methodCount++
	}
	if decryptAccessMethod != "" {
		methodCount++
	}

	if methodCount == 0 {
		return model.AccessMethodConfig{}, fmt.Errorf("at least one access method is required")
	}

	if methodCount > 1 {
		return model.AccessMethodConfig{}, fmt.Errorf("only one access method can be specified")
	}

	// Build the appropriate access method
	if decryptPassphrase {
		// Prompt for passphrase
		passphrase, err := promptForPassphrase()
		if err != nil {
			return model.AccessMethodConfig{}, fmt.Errorf("failed to read passphrase: %w", err)
		}

		return model.AccessMethodConfig{
			Type: model.AccessMethodTypePassphrase,
			Name: "passphrase",
			Config: map[string]interface{}{
				"passphrase": passphrase,
				"iterations": 100000,
			},
		}, nil
	}

	if decryptRuntime {
		return model.AccessMethodConfig{
			Type:   model.AccessMethodTypeRuntime,
			Name:   "runtime",
			Config: map[string]interface{}{},
		}, nil
	}

	if decryptCertificate != "" {
		certInfo, err := getCertificateInfo(decryptCertificate)
		if err != nil {
			return model.AccessMethodConfig{}, fmt.Errorf("failed to load certificate: %w", err)
		}

		// Prompt for certificate private key if needed
		privateKeyFile, err := promptForPrivateKey()
		if err != nil {
			return model.AccessMethodConfig{}, fmt.Errorf("failed to read private key: %w", err)
		}

		return model.AccessMethodConfig{
			Type: model.AccessMethodTypeCertificate,
			Name: "certificate-" + decryptCertificate,
			Config: map[string]interface{}{
				"certificate_file": decryptCertificate,
				"private_key_file": privateKeyFile,
				"key_id":           certInfo.KeyID,
			},
		}, nil
	}

	if decryptAccessMethod != "" {
		return parseAccessMethodString(decryptAccessMethod)
	}

	return model.AccessMethodConfig{}, fmt.Errorf("no valid access method specified")
}

func promptForPassphrase() (string, error) {
	fmt.Print("Enter passphrase: ")

	// Read passphrase securely without echoing
	passphrase, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}

	if passphrase == "" {
		return "", fmt.Errorf("passphrase cannot be empty")
	}

	// Confirm passphrase for security
	fmt.Print("Confirm passphrase: ")
	confirmPass, err := readPassword()
	if err != nil {
		return "", fmt.Errorf("failed to read confirmation: %w", err)
	}

	if passphrase != confirmPass {
		return "", fmt.Errorf("passphrases do not match")
	}

	return passphrase, nil
}

// readPassword reads password from stdin without echoing
func readPassword() (string, error) {
	// Try to use terminal password reading if available
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err == nil {
		fmt.Println() // Add newline after password input
		return string(password), nil
	}

	// Fallback to regular input (less secure)
	fmt.Print()
	var passphrase string
	_, err = fmt.Scanln(&passphrase)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return passphrase, nil
}

func promptForPrivateKey() (string, error) {
	fmt.Print("Enter private key file (optional): ")

	var privateKeyFile string
	_, err := fmt.Scanln(&privateKeyFile)
	if err != nil {
		// Empty input is valid, might be in agent
		return "", nil
	}

	// Validate file exists if provided
	if privateKeyFile != "" {
		if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
			return "", fmt.Errorf("private key file does not exist: %s", privateKeyFile)
		}
	}

	return privateKeyFile, nil
}

func displayDecryptionResult(result *model.DecryptionResult) {
	if result.Success {
		fmt.Println("✅ Decryption completed successfully")
		fmt.Println()

		// Display artifact info
		fmt.Printf("📦 Artifact ID: %s\n", result.ArtifactID)
		fmt.Printf("📁 Output: %s\n", result.FilePath)
		fmt.Printf("🔑 Method used: %s\n", result.MethodType)
		fmt.Printf("📊 Original size: %d bytes\n", result.OriginalSize)
		fmt.Printf("📊 Decrypted size: %d bytes\n", result.DecryptedSize)

		if result.OriginalSize > 0 && result.DecryptedSize > 0 {
			ratio := float64(result.DecryptedSize) / float64(result.OriginalSize) * 100
			fmt.Printf("📏 Size ratio: %.1f%%\n", ratio)
		}

		// Display creation time
		fmt.Println()
		fmt.Printf("⏰ Completed: %s\n", result.CreatedAt.Format("2006-01-02 15:04:05"))

	} else {
		fmt.Println("❌ Decryption failed")
		fmt.Println()

		// Display error information
		fmt.Printf("📦 Artifact ID: %s\n", result.ArtifactID)
		fmt.Printf("❌ Reason: %s\n", result.Reason)
		fmt.Printf("⏰ Attempted: %s\n", result.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

// ValidateDecryptionRequest validates the decryption request before processing
func validateDecryptionRequest(req *model.DecryptionRequest) error {
	// Check if artifact file exists
	if _, err := os.Stat(req.ArtifactPath); os.IsNotExist(err) {
		return fmt.Errorf("artifact file does not exist: %s", req.ArtifactPath)
	}

	// Check if artifact has correct extension
	if !strings.HasSuffix(req.ArtifactPath, ".ava") {
		return fmt.Errorf("artifact file must have .ava extension: %s", req.ArtifactPath)
	}

	// Check if output directory exists (and force is not set)
	if _, err := os.Stat(req.OutputPath); err == nil && !req.Force {
		return fmt.Errorf("output directory already exists: %s (use --force to overwrite)", req.OutputPath)
	}

	// Validate access method
	if req.AccessMethod.Type == "" {
		return fmt.Errorf("access method type is required")
	}

	return nil
}

// DisplayArtifactInfo displays information about an encrypted artifact
func displayArtifactInfo(artifactPath string) error {
	// TODO: Implement artifact info display
	// This would read the artifact metadata and display it
	fmt.Printf("📦 Artifact: %s\n", artifactPath)
	fmt.Println("   Type: Aether Vault Encrypted Artifact")
	fmt.Println("   Version: 1.0")
	fmt.Println("   Algorithm: AES-256-GCM")

	return nil
}

// ConfirmDecryption asks for user confirmation before decryption
func confirmDecryption(req *model.DecryptionRequest) bool {
	fmt.Printf("🔓 Ready to decrypt artifact: %s\n", req.ArtifactPath)
	fmt.Printf("📁 Output directory: %s\n", req.OutputPath)
	fmt.Printf("🔑 Access method: %s\n", req.AccessMethod.Type)
	fmt.Println()

	fmt.Print("Do you want to continue? (y/N): ")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// NewDecryptCommand creates the decrypt command
func NewDecryptCommand() *cobra.Command {
	return decryptCmd
}
