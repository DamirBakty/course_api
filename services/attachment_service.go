package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"web/config"
	"web/models"
	"web/repos"
	"web/schemas"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AttachmentServiceInterface interface {
	UploadFile(file *multipart.FileHeader, courseID, chapterID, lessonID uint) (schemas.UploadResponse, error)
	DownloadFile(id uint) (models.Attachment, *minio.Object, error)
	GetAttachmentsByLessonID(courseID, chapterID, lessonID uint) ([]models.Attachment, error)
	DeleteAttachment(id uint) error
	HasAccessToLesson(userID, lessonID uint) (bool, error)
}

type AttachmentService struct {
	config      *config.AppConfig
	repo        *repos.AttachmentRepository
	lessonRepo  *repos.LessonRepository
	uploadDir   string
	minioClient *minio.Client
}

func NewAttachmentService(config *config.AppConfig, repo *repos.AttachmentRepository, lessonRepo *repos.LessonRepository) (*AttachmentService, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Create the bucket if it doesn't exist
	exists, err := minioClient.BucketExists(context.Background(), config.MinioBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), config.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Use the current directory for file storage (as a fallback)
	uploadDir := "."

	return &AttachmentService{
		config:      config,
		repo:        repo,
		lessonRepo:  lessonRepo,
		uploadDir:   uploadDir,
		minioClient: minioClient,
	}, nil
}

func (s *AttachmentService) UploadFile(file *multipart.FileHeader, courseID, chapterID, lessonID uint) (schemas.UploadResponse, error) {
	// Validate lesson exists and belongs to the specified chapter and course
	if courseID > 0 && chapterID > 0 {
		_, err := s.lessonRepo.GetByID(courseID, chapterID, lessonID)
		if err != nil {
			return schemas.UploadResponse{}, fmt.Errorf("lesson not found or does not belong to the specified chapter and course: %w", err)
		}
	} else {
		// Backward compatibility: if courseID or chapterID is not provided, just check if the lesson exists
		_, err := s.lessonRepo.GetByID(0, 0, lessonID)
		if err != nil {
			return schemas.UploadResponse{}, fmt.Errorf("lesson not found: %w", err)
		}
	}

	// Generate a unique object name for the file in MinIO
	filename := filepath.Base(file.Filename)
	objectName := fmt.Sprintf("lesson-%d/%s", lessonID, filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return schemas.UploadResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Get the file size
	fileSize := file.Size

	// Upload the file to MinIO
	_, err = s.minioClient.PutObject(
		context.Background(),
		s.config.MinioBucket,
		objectName,
		src,
		fileSize,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return schemas.UploadResponse{}, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Create attachment record in the database
	attachment := models.Attachment{
		Name:     filename,
		URL:      objectName, // Store the object name in MinIO instead of a file path
		LessonID: lessonID,
	}

	id, err := s.repo.Create(attachment)
	if err != nil {
		return schemas.UploadResponse{}, fmt.Errorf("failed to create attachment record: %w", err)
	}

	// Generate the URL for the file
	url := fmt.Sprintf("/api/v1/attachments/download/%d", id)

	return schemas.UploadResponse{
		ID:       id,
		Name:     filename,
		URL:      url,
		LessonID: lessonID,
	}, nil
}

func (s *AttachmentService) DownloadFile(id uint) (models.Attachment, *minio.Object, error) {
	// Get attachment from database
	attachment, err := s.repo.GetByID(id)
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("attachment not found: %w", err)
	}

	// The URL field of the Attachment model contains the object name in MinIO
	objectName := attachment.URL

	// Check if the object exists in MinIO
	_, err = s.minioClient.StatObject(context.Background(), s.config.MinioBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("file not found in MinIO: %w", err)
	}

	// Get the object from MinIO
	object, err := s.minioClient.GetObject(context.Background(), s.config.MinioBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	return attachment, object, nil
}

func (s *AttachmentService) GetAttachmentsByLessonID(courseID, chapterID, lessonID uint) ([]models.Attachment, error) {
	// If courseID and chapterID are provided, validate the hierarchy
	if courseID > 0 && chapterID > 0 {
		// Validate that the lesson belongs to the specified chapter and course
		_, err := s.lessonRepo.GetByID(courseID, chapterID, lessonID)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.GetByLessonID(lessonID)
}

func (s *AttachmentService) DeleteAttachment(id uint) error {
	// Get attachment to get the object name
	attachment, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("attachment not found: %w", err)
	}

	// The URL field of the Attachment model contains the object name in MinIO
	objectName := attachment.URL

	// Delete the object from MinIO
	err = s.minioClient.RemoveObject(context.Background(), s.config.MinioBucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		// If the object doesn't exist, we can still proceed with deleting the record
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	// Delete the attachment record from the database
	return s.repo.Delete(id)
}

func (s *AttachmentService) HasAccessToLesson(userID, lessonID uint) (bool, error) {
	// This is a simplified implementation
	// In a real application, you would check if the user has access to the course that contains the lesson
	// For example, if they are enrolled in the course or if they are an admin/teacher

	// For now, we'll just check if the lesson exists
	_, err := s.lessonRepo.GetByID(0, 0, lessonID) // We don't have courseID and chapterID here, but the repository will handle it
	if err != nil {
		if err.Error() == "lesson not found" {
			return false, nil
		}
		return false, err
	}

	// In a real implementation, you would check if the user has access to the course
	// For example:
	// return s.courseService.HasAccess(userID, courseID)

	// For now, we'll just return true
	return true, nil
}
