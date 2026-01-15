package services

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/skygenesisenterprise/aether-vault/package/cli/server/model"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptionService handles encryption and decryption operations
type EncryptionService struct {
	masterKey     []byte
	kdfSalt       []byte
	kdfIterations int
	auditService  *AuditService
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(masterKey string, kdfSalt string, kdfIterations int, auditService *AuditService) *EncryptionService {
	salt := []byte(kdfSalt)
	key := pbkdf2.Key([]byte(masterKey), salt, kdfIterations, 32, sha256.New)

	return &EncryptionService{
		masterKey:     key,
		kdfSalt:       salt,
		kdfIterations: kdfIterations,
		auditService:  auditService,
	}
}

// Encrypt encrypts a file or directory according to the request
func (s *EncryptionService) Encrypt(req *model.EncryptionRequest, userID uuid.UUID) (*model.EncryptionResult, error) {
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
	accessMethods := make([]model.AccessMethod, 0, len(req.AccessMethods))
	for _, methodConfig := range req.AccessMethods {
		method, err := s.createAccessMethod(methodConfig, dataKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create access method %s: %w", methodConfig.Type, err)
		}
		accessMethods = append(accessMethods, *method)
	}

	// Process policies
	policies := make([]model.EncryptionPolicy, 0, len(req.Policies))
	for _, policyConfig := range req.Policies {
		policy, err := s.createEncryptionPolicy(policyConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create policy %s: %w", policyConfig.Type, err)
		}
		policies = append(policies, *policy)
	}

	// Create artifact metadata
	artifact := &model.EncryptedArtifact{
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
	result := &model.EncryptionResult{
		ArtifactID:    artifact.ID,
		FilePath:      req.OutputPath,
		OriginalSize:  artifact.ContentSize,
		EncryptedSize: finalFileInfo.Size(),
		Algorithm:     artifact.Algorithm,
		AccessMethods: s.getAccessMethodNames(accessMethods),
		Policies:      s.getPolicyNames(policies),
		CreatedAt:     artifact.CreatedAt,
	}

	// Log audit event
	if s.auditService != nil {
		metadata := map[string]interface{}{
			"message": fmt.Sprintf("Source: %s, Methods: %v", req.SourcePath, result.AccessMethods),
		}
		s.auditService.LogAction(context.Background(), userID, "artifact_encrypted", "artifact", artifact.ID, true, metadata)
	}

	return result, nil
}

// Decrypt decrypts an artifact according to the request
func (s *EncryptionService) Decrypt(req *model.DecryptionRequest, userID uuid.UUID) (*model.DecryptionResult, error) {
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
		// Log failed attempt
		if s.auditService != nil {
			metadata := map[string]interface{}{
				"message": fmt.Sprintf("Method: %s, Reason: %s", req.AccessMethod.Type, err.Error()),
			}
			s.auditService.LogAction(context.Background(), userID, "artifact_decrypt_failed", "artifact", artifact.ID, false, metadata)
		}
		return nil, err
	}

	// Decrypt data key
	dataKey, err := s.decryptDataKey(accessMethod, req.AccessMethod)
	if err != nil {
		if s.auditService != nil {
			metadata := map[string]interface{}{
				"message": fmt.Sprintf("Method: %s, Reason: %s", req.AccessMethod.Type, err.Error()),
			}
			s.auditService.LogAction(context.Background(), userID, "artifact_decrypt_failed", "artifact", artifact.ID, false, metadata)
		}
		return nil, fmt.Errorf("failed to decrypt data key: %w", err)
	}

	// Validate policies
	err = s.validatePolicies(artifact.Policies, userID)
	if err != nil {
		if s.auditService != nil {
			metadata := map[string]interface{}{
				"message": fmt.Sprintf("Policy validation failed: %s", err.Error()),
			}
			s.auditService.LogAction(context.Background(), userID, "artifact_decrypt_failed", "artifact", artifact.ID, false, metadata)
		}
		return nil, fmt.Errorf("policy validation failed: %w", err)
	}

	// Decrypt content
	contentPath := filepath.Join(tempDir, "content.tar.gz")
	err = s.decryptFile(encryptedContentPath, contentPath, dataKey)
	if err != nil {
		if s.auditService != nil {
			metadata := map[string]interface{}{
				"message": fmt.Sprintf("Content decryption failed: %s", err.Error()),
			}
			s.auditService.LogAction(context.Background(), userID, "artifact_decrypt_failed", "artifact", artifact.ID, false, metadata)
		}
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}

	// Extract archive
	err = s.extractArchive(contentPath, req.OutputPath)
	if err != nil {
		if s.auditService != nil {
			metadata := map[string]interface{}{
				"message": fmt.Sprintf("Archive extraction failed: %s", err.Error()),
			}
			s.auditService.LogAction(context.Background(), userID, "artifact_decrypt_failed", "artifact", artifact.ID, false, metadata)
		}
		return nil, fmt.Errorf("failed to extract archive: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(req.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	// Create result
	result := &model.DecryptionResult{
		ArtifactID:    artifact.ID,
		FilePath:      req.OutputPath,
		OriginalSize:  artifact.ContentSize,
		DecryptedSize: fileInfo.Size(),
		MethodType:    string(accessMethod.Type),
		Success:       true,
		CreatedAt:     time.Now(),
	}

	// Log successful decryption
	if s.auditService != nil {
		metadata := map[string]interface{}{
			"message": fmt.Sprintf("Method: %s, Output: %s", accessMethod.Type, req.OutputPath),
		}
		s.auditService.LogAction(context.Background(), userID, "artifact_decrypted", "artifact", artifact.ID, true, metadata)
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

func (s *EncryptionService) createAccessMethod(config model.AccessMethodConfig, dataKey []byte) (*model.AccessMethod, error) {
	method := &model.AccessMethod{
		ID:        uuid.New(),
		Type:      config.Type,
		Name:      config.Name,
		Config:    s.configToJSON(config.Config),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	switch config.Type {
	case model.AccessMethodTypePassphrase:
		return s.createPassphraseMethod(method, config, dataKey)
	case model.AccessMethodTypeCertificate:
		return s.createCertificateMethod(method, config, dataKey)
	case model.AccessMethodTypeRuntime:
		return s.createRuntimeMethod(method, config, dataKey)
	default:
		return nil, fmt.Errorf("unsupported access method type: %s", config.Type)
	}
}

func (s *EncryptionService) createPassphraseMethod(method *model.AccessMethod, config model.AccessMethodConfig, dataKey []byte) (*model.AccessMethod, error) {
	// For passphrase, we'll store the encrypted data key
	// The actual passphrase will be required during decryption
	encryptedKey, err := s.encryptWithMasterKey(dataKey)
	if err != nil {
		return nil, err
	}
	method.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedKey)
	return method, nil
}

func (s *EncryptionService) createCertificateMethod(method *model.AccessMethod, config model.AccessMethodConfig, dataKey []byte) (*model.AccessMethod, error) {
	// TODO: Implement certificate-based encryption
	return nil, fmt.Errorf("certificate method not yet implemented")
}

func (s *EncryptionService) createRuntimeMethod(method *model.AccessMethod, config model.AccessMethodConfig, dataKey []byte) (*model.AccessMethod, error) {
	// Encrypt with runtime master key
	encryptedKey, err := s.encryptWithMasterKey(dataKey)
	if err != nil {
		return nil, err
	}
	method.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedKey)
	return method, nil
}

func (s *EncryptionService) createEncryptionPolicy(config model.EncryptionPolicyConfig) (*model.EncryptionPolicy, error) {
	policy := &model.EncryptionPolicy{
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

func (s *EncryptionService) configToJSON(config map[string]interface{}) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func (s *EncryptionService) getAccessMethodNames(methods []model.AccessMethod) []string {
	names := make([]string, len(methods))
	for i, method := range methods {
		names[i] = method.Name
	}
	return names
}

func (s *EncryptionService) getPolicyNames(policies []model.EncryptionPolicy) []string {
	names := make([]string, len(policies))
	for i, policy := range policies {
		names[i] = policy.Name
	}
	return names
}

func (s *EncryptionService) createArtifactFile(artifact *model.EncryptedArtifact, encryptedContentPath, tempDir string) error {
	// TODO: Implement artifact file creation with proper format
	// For now, just copy the encrypted content
	return os.Rename(encryptedContentPath, artifact.FilePath)
}

func (s *EncryptionService) loadArtifact(artifactPath string) (*model.EncryptedArtifact, error) {
	// TODO: Implement artifact loading from file
	// For now, return a mock artifact
	return &model.EncryptedArtifact{
		ID:        uuid.New(),
		Algorithm: "AES-256-GCM",
		Version:   "1.0",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *EncryptionService) extractEncryptedContent(artifactPath, outputPath string) error {
	// TODO: Implement content extraction from artifact
	// For now, just copy the file
	return os.Rename(artifactPath, outputPath)
}

func (s *EncryptionService) findAccessMethod(artifact *model.EncryptedArtifact, config model.AccessMethodConfig) (*model.AccessMethod, error) {
	for _, method := range artifact.AccessMethods {
		if method.Type == config.Type && method.IsActive {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no valid access method found for type: %s", config.Type)
}

func (s *EncryptionService) decryptDataKey(method *model.AccessMethod, config model.AccessMethodConfig) ([]byte, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(method.EncryptedKey)
	if err != nil {
		return nil, err
	}

	switch method.Type {
	case model.AccessMethodTypeRuntime:
		return s.decryptWithMasterKey(encryptedKey)
	case model.AccessMethodTypePassphrase:
		// TODO: Implement passphrase-based decryption
		return s.decryptWithMasterKey(encryptedKey)
	default:
		return nil, fmt.Errorf("unsupported method type for data key decryption: %s", method.Type)
	}
}

func (s *EncryptionService) validatePolicies(policies []model.EncryptionPolicy, userID uuid.UUID) error {
	// TODO: Implement policy validation
	return nil
}
