package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
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
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: config.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

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
	if courseID > 0 && chapterID > 0 {
		_, err := s.lessonRepo.GetByID(courseID, chapterID, lessonID)
		if err != nil {
			return schemas.UploadResponse{}, fmt.Errorf("lesson not found or does not belong to the specified chapter and course: %w", err)
		}
	} else {
		_, err := s.lessonRepo.GetByID(0, 0, lessonID)
		if err != nil {
			return schemas.UploadResponse{}, fmt.Errorf("lesson not found: %w", err)
		}
	}

	filename := filepath.Base(file.Filename)
	objectName := fmt.Sprintf("lesson-%d/%s", lessonID, filename)

	src, err := file.Open()
	if err != nil {
		return schemas.UploadResponse{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileSize := file.Size

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

	attachment := models.Attachment{
		Name:     filename,
		URL:      objectName,
		LessonID: lessonID,
	}

	id, err := s.repo.Create(attachment)
	if err != nil {
		return schemas.UploadResponse{}, fmt.Errorf("failed to create attachment record: %w", err)
	}

	presignedURL, err := s.generatePresignedURL(objectName)
	if err != nil {

		url := fmt.Sprintf("/api/v1/attachments/download/%d", id)
		return schemas.UploadResponse{
			ID:       id,
			Name:     filename,
			URL:      url,
			LessonID: lessonID,
		}, nil
	}

	return schemas.UploadResponse{
		ID:       id,
		Name:     filename,
		URL:      presignedURL,
		LessonID: lessonID,
	}, nil
}

func (s *AttachmentService) DownloadFile(id uint) (models.Attachment, *minio.Object, error) {

	attachment, err := s.repo.GetByID(id)
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("attachment not found: %w", err)
	}

	objectName := attachment.URL

	_, err = s.minioClient.StatObject(context.Background(), s.config.MinioBucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("file not found in MinIO: %w", err)
	}

	object, err := s.minioClient.GetObject(context.Background(), s.config.MinioBucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return models.Attachment{}, nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	return attachment, object, nil
}

func (s *AttachmentService) GetAttachmentsByLessonID(courseID, chapterID, lessonID uint) ([]models.Attachment, error) {
	if courseID > 0 && chapterID > 0 {
		_, err := s.lessonRepo.GetByID(courseID, chapterID, lessonID)
		if err != nil {
			return nil, err
		}
	}

	attachments, err := s.repo.GetByLessonID(lessonID)
	if err != nil {
		return nil, err
	}

	for i := range attachments {
		objectName := attachments[i].URL

		presignedURL, err := s.generatePresignedURL(objectName)
		if err != nil {
			fmt.Printf("Failed to generate pre-signed URL for %s: %v\n", objectName, err)
			continue
		}

		attachments[i].URL = presignedURL
	}

	return attachments, nil
}

func (s *AttachmentService) generatePresignedURL(objectName string) (string, error) {
	expiry := time.Hour * 24

	presignedURL, err := s.minioClient.PresignedGetObject(context.Background(), s.config.MinioBucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return presignedURL.String(), nil
}

func (s *AttachmentService) DeleteAttachment(id uint) error {
	attachment, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("attachment not found: %w", err)
	}

	objectName := attachment.URL

	err = s.minioClient.RemoveObject(context.Background(), s.config.MinioBucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return s.repo.Delete(id)
}

func (s *AttachmentService) HasAccessToLesson(userID, lessonID uint) (bool, error) {
	_, err := s.lessonRepo.GetByID(0, 0, lessonID)
	if err != nil {
		if err.Error() == "lesson not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
