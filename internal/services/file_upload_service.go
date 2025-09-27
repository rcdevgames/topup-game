package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	appConfig "topup-backend/internal/config"
)

// FileUploadService handles file upload operations
type FileUploadService struct {
	cfg      *appConfig.Config
	s3Client *s3.Client
	useS3    bool
}

// FileUploadResult represents file upload result
type FileUploadResult struct {
	URL      string `json:"url"`
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
}

// NewFileUploadService creates a new file upload service
func NewFileUploadService(cfg *appConfig.Config) *FileUploadService {
	service := &FileUploadService{
		cfg:   cfg,
		useS3: false,
	}

	// Try to initialize S3 client if AWS credentials are provided
	if cfg.AWS.S3Bucket != "" && cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		if err := service.initS3Client(); err != nil {
			fmt.Printf("Warning: Failed to initialize S3 client: %v, falling back to local storage\n", err)
		} else {
			service.useS3 = true
			fmt.Println("✅ S3 client initialized successfully")
		}
	} else {
		fmt.Println("ℹ️  S3 credentials not configured, using local storage")
	}

	// Ensure local upload directory exists
	if err := service.ensureUploadDir(); err != nil {
		fmt.Printf("Warning: Failed to create upload directory: %v\n", err)
	}

	return service
}

// initS3Client initializes AWS S3 client
func (s *FileUploadService) initS3Client() error {
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.cfg.AWS.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     s.cfg.AWS.AccessKeyID,
				SecretAccessKey: s.cfg.AWS.SecretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s.s3Client = s3.NewFromConfig(cfg)
	return nil
}

// ensureUploadDir ensures upload directory exists
func (s *FileUploadService) ensureUploadDir() error {
	dirs := []string{
		s.cfg.Server.UploadDir,
		filepath.Join(s.cfg.Server.UploadDir, "products"),
		filepath.Join(s.cfg.Server.UploadDir, "categories"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// UploadProductImage uploads product image
func (s *FileUploadService) UploadProductImage(file multipart.File, filename string) (*FileUploadResult, error) {
	return s.uploadFile(file, filename, "products")
}

// UploadCategoryIcon uploads category icon
func (s *FileUploadService) UploadCategoryIcon(file multipart.File, filename string) (*FileUploadResult, error) {
	return s.uploadFile(file, filename, "categories")
}

// uploadFile uploads a file to either S3 or local storage
func (s *FileUploadService) uploadFile(file multipart.File, filename, folder string) (*FileUploadResult, error) {
	// Validate file type
	if err := s.validateFileType(filename); err != nil {
		return nil, err
	}

	// Generate unique filename
	uniqueFilename := s.generateUniqueFilename(filename)

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	if s.useS3 {
		return s.uploadToS3(fileContent, uniqueFilename, folder)
	}

	return s.uploadToLocal(fileContent, uniqueFilename, folder)
}

// uploadToS3 uploads file to AWS S3
func (s *FileUploadService) uploadToS3(content []byte, filename, folder string) (*FileUploadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s/%s", folder, filename)

	// Upload to S3
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.AWS.S3Bucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(string(content)),
		ContentType: aws.String(s.getContentType(filename)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.cfg.AWS.S3Bucket, s.cfg.AWS.Region, key)

	return &FileUploadResult{
		URL:      url,
		FileName: filename,
		Size:     int64(len(content)),
	}, nil
}

// uploadToLocal uploads file to local storage
func (s *FileUploadService) uploadToLocal(content []byte, filename, folder string) (*FileUploadResult, error) {
	// Create file path
	filePath := filepath.Join(s.cfg.Server.UploadDir, folder, filename)

	// Write file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file to local storage: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("%s/cdn/%s/%s", s.cfg.Server.BaseURL, folder, filename)

	return &FileUploadResult{
		URL:      url,
		FileName: filename,
		Size:     int64(len(content)),
	}, nil
}

// DeleteFile deletes a file from storage
func (s *FileUploadService) DeleteFile(fileURL string) error {
	if s.useS3 {
		return s.deleteFromS3(fileURL)
	}

	return s.deleteFromLocal(fileURL)
}

// deleteFromS3 deletes file from AWS S3
func (s *FileUploadService) deleteFromS3(fileURL string) error {
	// Extract key from URL
	key := s.extractS3KeyFromURL(fileURL)
	if key == "" {
		return fmt.Errorf("invalid S3 URL: %s", fileURL)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.AWS.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// deleteFromLocal deletes file from local storage
func (s *FileUploadService) deleteFromLocal(fileURL string) error {
	// Extract filename from URL
	filename := s.extractFilenameFromURL(fileURL)
	if filename == "" {
		return fmt.Errorf("invalid local file URL: %s", fileURL)
	}

	// Find file in upload directories
	dirs := []string{"products", "categories"}
	for _, dir := range dirs {
		filePath := filepath.Join(s.cfg.Server.UploadDir, dir, filename)
		if _, err := os.Stat(filePath); err == nil {
			return os.Remove(filePath)
		}
	}

	return fmt.Errorf("file not found: %s", filename)
}

// validateFileType validates uploaded file type
func (s *FileUploadService) validateFileType(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	return nil
}

// generateUniqueFilename generates a unique filename
func (s *FileUploadService) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)

	// Clean filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)

	// Add UUID
	id := uuid.New().String()[:8]

	return fmt.Sprintf("%s_%s%s", name, id, ext)
}

// getContentType returns content type based on file extension
func (s *FileUploadService) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}

// extractS3KeyFromURL extracts S3 key from S3 URL
func (s *FileUploadService) extractS3KeyFromURL(url string) string {
	// Expected format: https://bucket.s3.region.amazonaws.com/key
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return ""
	}

	// Join everything after the domain
	return strings.Join(parts[3:], "/")
}

// extractFilenameFromURL extracts filename from local URL
func (s *FileUploadService) extractFilenameFromURL(url string) string {
	// Expected format: http://localhost:8080/cdn/folder/filename
	parts := strings.Split(url, "/")
	if len(parts) < 1 {
		return ""
	}

	return parts[len(parts)-1]
}
