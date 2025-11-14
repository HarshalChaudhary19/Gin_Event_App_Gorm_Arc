package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type FileRepoI interface {
	CreateFile(*models.FileUpload) error
	UpdateFile(*models.FileUpload) (*models.FileUpload, error)
	DeleteFile(int) error
	GetFileById(int) (*models.FileUpload, error)
	GetAllFiles() ([]*models.FileUpload, error)
}

type FileRepo struct {
	DB *gorm.DB
}

func InitFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{
		DB: db,
	}
}

func (repo *FileRepo) CreateFile(file *models.FileUpload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Create(file).Error
}

func (repo *FileRepo) UpdateFile(file *models.FileUpload) (*models.FileUpload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id=?", file.ID).Model(&models.FileUpload{}).Updates(models.FileUpload{ //Basic Update Query hai ye
		ID:          file.ID,
		FileName:    file.FileName,
		ContentType: file.ContentType,
		Size:        file.Size,
		Data:        file.Data,
	}).Error
	return file, err
}

func (repo *FileRepo) DeleteFile(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Delete(&models.FileUpload{}, id).Error
}

func (repo *FileRepo) GetFileById(id int) (*models.FileUpload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var file *models.FileUpload
	err := repo.DB.WithContext(ctx).First(&file, id).Error
	return file, err
}

func (repo *FileRepo) GetAllFiles() ([]*models.FileUpload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var file []*models.FileUpload
	err := repo.DB.WithContext(ctx).Find(&file).Error
	return file, err
}
