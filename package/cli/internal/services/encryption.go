package services

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptionService handles encryption and decryption operations
type EncryptionService struct {
	masterKey     []byte
	kdfSalt       []byte
	kdfIterations int
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(masterKey string, kdfSalt string, kdfIterations int) *EncryptionService {
	salt := []byte(kdfSalt)
	key := pbkdf2.Key([]byte(masterKey), salt, kdfIterations, 32, sha256.New)

	return &EncryptionService{
		masterKey:     key,
		kdfSalt:       salt,
		kdfIterations: kdfIterations,
	}
}

// AccessMethodType represents the type of access method
type AccessMethodType string

const (
	AccessMethodTypePassphrase  AccessMethodType = "passphrase"
	AccessMethodTypeCertificate AccessMethodType = "certificate"
	AccessMethodTypeRuntime     AccessMethodType = "runtime"
	AccessMethodTypePolicy      AccessMethodType = "policy"
)

// PolicyType represents the type of policy
type PolicyType string

const (
	PolicyTypeTTL         PolicyType = "ttl"
	PolicyTypeEnvironment PolicyType = "environment"
	PolicyTypeInstance    PolicyType = "instance"
	PolicyTypeRegion      PolicyType = "region"
	PolicyTypeMultiFactor PolicyType = "multi_factor"
)

// AccessMethod represents an access method for decrypting an artifact
type AccessMethod struct {
	ID           uuid.UUID        `json:"id"`
	ArtifactID   uuid.UUID        `json:"artifact_id"`
	Type         AccessMethodType `json:"type"`
	Name         string           `json:"name"`
	Config       string           `json:"config"` // JSON config
	EncryptedKey string           `json:"encrypted_key"`
	KeyID        string           `json:"key_id"` // For certificate methods
	IsActive     bool             `json:"is_active"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// EncryptionPolicy represents a policy for encrypted artifacts
type EncryptionPolicy struct {
	ID         uuid.UUID  `json:"id"`
	ArtifactID uuid.UUID  `json:"artifact_id"`
	Name       string     `json:"name"`
	Type       PolicyType `json:"type"`
	Rules      string     `json:"rules"` // JSON rules
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// EncryptedArtifact represents an encrypted artifact
type EncryptedArtifact struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	FilePath     string    `json:"file_path"`
	OriginalPath string    `json:"original_path"`
	Algorithm    string    `json:"algorithm"`
	Version      string    `json:"version"`
	DataKeyHash  string    `json:"data_key_hash"`
	ContentSize  int64     `json:"content_size"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	AccessMethods []AccessMethod     `json:"access_methods"`
	Policies      []EncryptionPolicy `json:"policies"`
}

// AccessMethodConfig represents configuration for an access method
type AccessMethodConfig struct {
	Type   AccessMethodType       `json:"type"`
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// EncryptionPolicyConfig represents configuration for an encryption policy
type EncryptionPolicyConfig struct {
	Type  PolicyType             `json:"type"`
	Name  string                 `json:"name"`
	Rules map[string]interface{} `json:"rules"`
}

// EncryptionRequest represents a request to encrypt data
type EncryptionRequest struct {
	SourcePath    string                   `json:"source_path"`
	OutputPath    string                   `json:"output_path"`
	AccessMethods []AccessMethodConfig     `json:"access_methods"`
	Policies      []EncryptionPolicyConfig `json:"policies"`
	Description   string                   `json:"description"`
	Compression   bool                     `json:"compression"`
}

// DecryptionRequest represents a request to decrypt data
type DecryptionRequest struct {
	ArtifactPath string             `json:"artifact_path"`
	OutputPath   string             `json:"output_path"`
	AccessMethod AccessMethodConfig `json:"access_method"`
	Force        bool               `json:"force"`
}

// EncryptionResult represents the result of an encryption operation
type EncryptionResult struct {
	ArtifactID    uuid.UUID `json:"artifact_id"`
	FilePath      string    `json:"file_path"`
	OriginalSize  int64     `json:"original_size"`
	EncryptedSize int64     `json:"encrypted_size"`
	Algorithm     string    `json:"algorithm"`
	AccessMethods []string  `json:"access_methods"`
	Policies      []string  `json:"policies"`
	CreatedAt     time.Time `json:"created_at"`
}

// DecryptionResult represents the result of a decryption operation
type DecryptionResult struct {
	ArtifactID    uuid.UUID `json:"artifact_id"`
	FilePath      string    `json:"file_path"`
	OriginalSize  int64     `json:"original_size"`
	DecryptedSize int64     `json:"decrypted_size"`
	MethodType    string    `json:"method_type"`
	Success       bool      `json:"success"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
}

// Encrypt encrypts a file or directory according to the request
func (s *EncryptionService) Encrypt(req *EncryptionRequest, userID uuid.UUID) (*EncryptionResult, error) {
	// Generate data key
	dataKey, err := s.generateDataKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate data key: %w", err)
	}

	// Create temporary directory for processing
	tempDir, err := os.MkdirTemp("", "aether-encrypt-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Create content archive
	contentPath := filepath.Join(tempDir, "content.tar.gz")
	if req.Compression {
		err = s.createCompressedArchive(req.SourcePath, contentPath)
	} else {
		err = s.createArchive(req.SourcePath, contentPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}

	// Encrypt content
	encryptedContentPath := filepath.Join(tempDir, "content.enc")
	err = s.encryptFile(contentPath, encryptedContentPath, dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}

	// Process access methods
	accessMethods := make([]AccessMethod, 0, len(req.AccessMethods))
	for _, methodConfig := range req.AccessMethods {
		method, err := s.createAccessMethod(methodConfig, dataKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create access method %s: %w", methodConfig.Type, err)
		}
		accessMethods = append(accessMethods, *method)
	}

	// Process policies
	policies := make([]EncryptionPolicy, 0, len(req.Policies))
	for _, policyConfig := range req.Policies {
		policy, err := s.createEncryptionPolicy(policyConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create policy %s: %w", policyConfig.Type, err)
		}
		policies = append(policies, *policy)
	}

	// Create artifact metadata
	artifact := &EncryptedArtifact{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          filepath.Base(req.SourcePath),
		Description:   req.Description,
		FilePath:      req.OutputPath,
		OriginalPath:  req.SourcePath,
		Algorithm:     "AES-256-GCM",
		Version:       "1.0",
		DataKeyHash:   s.hashDataKey(dataKey),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		AccessMethods: accessMethods,
		Policies:      policies,
	}

	// Get file info
	fileInfo, err := os.Stat(contentPath)
	if err == nil {
		artifact.ContentSize = fileInfo.Size()
	}

	// Create final artifact file
	err = s.createArtifactFile(artifact, encryptedContentPath, tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create artifact file: %w", err)
	}

	// Get final file size
	finalFileInfo, err := os.Stat(req.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get final file info: %w", err)
	}

	// Create result
	result := &EncryptionResult{
		ArtifactID:    artifact.ID,
		FilePath:      req.OutputPath,
		OriginalSize:  artifact.ContentSize,
		EncryptedSize: finalFileInfo.Size(),
		Algorithm:     artifact.Algorithm,
		AccessMethods: s.getAccessMethodNames(accessMethods),
		Policies:      s.getPolicyNames(policies),
		CreatedAt:     artifact.CreatedAt,
	}

	return result, nil
}

// Decrypt decrypts an artifact according to the request
func (s *EncryptionService) Decrypt(req *DecryptionRequest, userID uuid.UUID) (*DecryptionResult, error) {
	// Load artifact
	artifact, err := s.loadArtifact(req.ArtifactPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load artifact: %w", err)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "aether-decrypt-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract encrypted content
	encryptedContentPath := filepath.Join(tempDir, "content.enc")
	err = s.extractEncryptedContent(req.ArtifactPath, encryptedContentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract encrypted content: %w", err)
	}

	// Find and validate access method
	accessMethod, err := s.findAccessMethod(artifact, req.AccessMethod)
	if err != nil {
		return nil, err
	}

	// Decrypt data key
	dataKey, err := s.decryptDataKey(accessMethod, req.AccessMethod)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
	}

	// Validate policies
	err = s.validatePolicies(artifact.Policies, userID)
	if err != nil {
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}

	// Decrypt content
	contentPath := filepath.Join(tempDir, "content.tar.gz")
	err = s.decryptFile(encryptedContentPath, contentPath, dataKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}

	// Extract archive
	err = s.extractArchive(contentPath, req.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract archive: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(req.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	// Create result
	result := &DecryptionResult{
		ArtifactID:    artifact.ID,
		FilePath:      req.OutputPath,
		OriginalSize:  artifact.ContentSize,
		DecryptedSize: fileInfo.Size(),
		MethodType:    string(accessMethod.Type),
		Success:       true,
		CreatedAt:     time.Now(),
	}

	return result, nil
}

// Helper methods

func (s *EncryptionService) generateDataKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *EncryptionService) hashDataKey(dataKey []byte) string {
	hash := sha256.Sum256(dataKey)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (s *EncryptionService) encryptFile(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return s.encryptStream(inputFile, outputFile, key)
}

func (s *EncryptionService) decryptFile(inputPath, outputPath string, key []byte) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return s.decryptStream(inputFile, outputFile, key)
}

func (s *EncryptionService) encryptStream(src io.Reader, dst io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Write nonce
	if _, err = dst.Write(nonce); err != nil {
		return err
	}

	// Encrypt and write
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	_, err = dst.Write(gcm.Seal(nil, nonce, data, nil))
	return err
}

func (s *EncryptionService) decryptStream(src io.Reader, dst io.Writer, key []byte) error {
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	_, err = dst.Write(plaintext)
	return err
}

func (s *EncryptionService) createArchive(sourcePath, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself if we're archiving a directory
		if path == sourcePath && info.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Make path relative to source
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.Open(path)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *EncryptionService) createCompressedArchive(sourcePath, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == sourcePath && info.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.Open(path)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *EncryptionService) extractArchive(archivePath, outputPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file

	// Try to detect if it's gzipped
	if strings.HasSuffix(archivePath, ".gz") {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(outputPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *EncryptionService) createAccessMethod(config AccessMethodConfig, dataKey []byte) (*AccessMethod, error) {
	method := &AccessMethod{
		ID:        uuid.New(),
		Type:      config.Type,
		Name:      config.Name,
		Config:    s.configToJSON(config.Config),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	switch config.Type {
	case AccessMethodTypePassphrase:
		return s.createPassphraseMethod(method, config, dataKey)
	case AccessMethodTypeCertificate:
		return s.createCertificateMethod(method, config, dataKey)
	case AccessMethodTypeRuntime:
		return s.createRuntimeMethod(method, config, dataKey)
	default:
		return nil, fmt.Errorf("unsupported access method type: %s", config.Type)
	}
}

func (s *EncryptionService) createPassphraseMethod(method *AccessMethod, config AccessMethodConfig, dataKey []byte) (*AccessMethod, error) {
	// Generate salt for PBKDF2
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	iterations := 100000
	if iter, exists := config.Config["iterations"].(int); exists {
		iterations = iter
	} else if iterFloat, exists := config.Config["iterations"].(float64); exists {
		iterations = int(iterFloat)
	}

	// For passphrase method, we'll encrypt the data key with a temporary key
	// The actual passphrase will be used to decrypt during decryption
	tempKey := make([]byte, 32)
	if _, err := rand.Read(tempKey); err != nil {
		return nil, fmt.Errorf("failed to generate temp key: %w", err)
	}

	// Encrypt data key with temp key
	encryptedKey, err := s.encryptDataKeyWithKey(dataKey, tempKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data key: %w", err)
	}

	// Store the encrypted key and salt
	method.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedKey)

	// Update config to include salt and iterations
	if method.Config == "" {
		method.Config = "{}"
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(method.Config), &configMap); err != nil {
		configMap = make(map[string]interface{})
	}

	configMap["salt"] = base64.StdEncoding.EncodeToString(salt)
	configMap["iterations"] = iterations

	updatedConfig, _ := json.Marshal(configMap)
	method.Config = string(updatedConfig)

	return method, nil
}

func (s *EncryptionService) createCertificateMethod(method *AccessMethod, config AccessMethodConfig, dataKey []byte) (*AccessMethod, error) {
	// TODO: Implement certificate-based encryption
	return nil, fmt.Errorf("certificate method not yet implemented")
}

func (s *EncryptionService) createRuntimeMethod(method *AccessMethod, config AccessMethodConfig, dataKey []byte) (*AccessMethod, error) {
	// Encrypt with runtime master key
	encryptedKey, err := s.encryptWithMasterKey(dataKey)
	if err != nil {
		return nil, err
	}
	method.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedKey)
	return method, nil
}

func (s *EncryptionService) createEncryptionPolicy(config EncryptionPolicyConfig) (*EncryptionPolicy, error) {
	policy := &EncryptionPolicy{
		ID:        uuid.New(),
		Type:      config.Type,
		Name:      config.Name,
		Rules:     s.configToJSON(config.Rules),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return policy, nil
}

func (s *EncryptionService) encryptWithMasterKey(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *EncryptionService) decryptWithMasterKey(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (s *EncryptionService) encryptDataKeyWithKey(dataKey []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, dataKey, nil), nil
}

func (s *EncryptionService) configToJSON(config map[string]interface{}) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func (s *EncryptionService) getAccessMethodNames(methods []AccessMethod) []string {
	names := make([]string, len(methods))
	for i, method := range methods {
		names[i] = method.Name
	}
	return names
}

func (s *EncryptionService) getPolicyNames(policies []EncryptionPolicy) []string {
	names := make([]string, len(policies))
	for i, policy := range policies {
		names[i] = policy.Name
	}
	return names
}

func (s *EncryptionService) createArtifactFile(artifact *EncryptedArtifact, encryptedContentPath, tempDir string) error {
	// Create artifact with Aether Vault format
	outputFile, err := os.Create(artifact.FilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write AVA header
	header := []byte("AVA\000")
	if _, err := outputFile.Write(header); err != nil {
		return err
	}

	// Write metadata JSON
	metadata, err := json.Marshal(artifact)
	if err != nil {
		return err
	}

	// Write metadata length (4 bytes)
	metadataLen := uint32(len(metadata))
	if err := binary.Write(outputFile, binary.BigEndian, metadataLen); err != nil {
		return err
	}

	// Write metadata
	if _, err := outputFile.Write(metadata); err != nil {
		return err
	}

	// Copy encrypted content
	contentFile, err := os.Open(encryptedContentPath)
	if err != nil {
		return err
	}
	defer contentFile.Close()

	_, err = io.Copy(outputFile, contentFile)
	return err
}

func (s *EncryptionService) loadArtifact(artifactPath string) (*EncryptedArtifact, error) {
	file, err := os.Open(artifactPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Check AVA header
	header := make([]byte, 4)
	if _, err := file.Read(header); err != nil {
		return nil, err
	}

	if string(header) != "AVA\000" {
		return nil, fmt.Errorf("invalid artifact format")
	}

	// Read metadata length
	var metadataLen uint32
	if err := binary.Read(file, binary.BigEndian, &metadataLen); err != nil {
		return nil, err
	}

	// Read metadata
	metadata := make([]byte, metadataLen)
	if _, err := file.Read(metadata); err != nil {
		return nil, err
	}

	// Parse metadata
	var artifact EncryptedArtifact
	if err := json.Unmarshal(metadata, &artifact); err != nil {
		return nil, err
	}

	return &artifact, nil
}

func (s *EncryptionService) extractEncryptedContent(artifactPath, outputPath string) error {
	file, err := os.Open(artifactPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Skip AVA header (4 bytes)
	if _, err := file.Seek(4, io.SeekStart); err != nil {
		return err
	}

	// Read metadata length
	var metadataLen uint32
	if err := binary.Read(file, binary.BigEndian, &metadataLen); err != nil {
		return err
	}

	// Skip metadata
	if _, err := file.Seek(int64(4+4+metadataLen), io.SeekStart); err != nil {
		return err
	}

	// Copy encrypted content to output
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, file)
	return err
}

func (s *EncryptionService) findAccessMethod(artifact *EncryptedArtifact, config AccessMethodConfig) (*AccessMethod, error) {
	for _, method := range artifact.AccessMethods {
		if method.Type == config.Type && method.IsActive {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no valid access method found for type: %s", config.Type)
}

func (s *EncryptionService) decryptDataKey(method *AccessMethod, config AccessMethodConfig) ([]byte, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(method.EncryptedKey)
	if err != nil {
		return nil, err
	}

	switch method.Type {
	case AccessMethodTypeRuntime:
		return s.decryptWithMasterKey(encryptedKey)
	case AccessMethodTypePassphrase:
		// Get passphrase from config
		passphrase, ok := config.Config["passphrase"].(string)
		if !ok || passphrase == "" {
			return nil, fmt.Errorf("passphrase is required for passphrase-based decryption")
		}

		// Derive key from passphrase using PBKDF2
		var salt []byte
		if saltStr, exists := config.Config["salt"].(string); exists && saltStr != "" {
			salt, err = base64.StdEncoding.DecodeString(saltStr)
			if err != nil {
				return nil, fmt.Errorf("failed to decode salt: %w", err)
			}
		} else {
			// Use a default salt (in production, this should be stored with the method)
			salt = []byte("aether-vault-default-salt")
		}

		iterations := 100000
		if iter, exists := config.Config["iterations"].(int); exists {
			iterations = iter
		} else if iterFloat, exists := config.Config["iterations"].(float64); exists {
			iterations = int(iterFloat)
		}

		// Derive decryption key from passphrase
		derivedKey := pbkdf2.Key([]byte(passphrase), salt, iterations, 32, sha256.New)

		// Decrypt data key with derived key
		block, err := aes.NewCipher(derivedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create cipher: %w", err)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, fmt.Errorf("failed to create GCM: %w", err)
		}

		nonceSize := gcm.NonceSize()
		if len(encryptedKey) < nonceSize {
			return nil, fmt.Errorf("encrypted key too short")
		}

		nonce, ciphertext := encryptedKey[:nonceSize], encryptedKey[nonceSize:]
		dataKey, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data key: %w", err)
		}

		return dataKey, nil
	default:
		return nil, fmt.Errorf("unsupported method type for data key decryption: %s", method.Type)
	}
}

func (s *EncryptionService) validatePolicies(policies []EncryptionPolicy, userID uuid.UUID) error {
	// TODO: Implement policy validation
	return nil
}
