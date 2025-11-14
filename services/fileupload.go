package services

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
)

type FileService struct {
	FileServ repositories.FileRepoI
}

type FileServiceI interface {
	UploadFile(*models.FileUpload) error
	GetFile(int) (*models.FileUpload, error)
	DeleteFile(int) error
	UpdateFile(*models.FileUpload) (*models.FileUpload, error)
	GetAllFiles() ([]*models.FileUpload, error)
}

func NewFileService(fileRepo repositories.FileRepoI) FileServiceI {
	return &FileService{FileServ: fileRepo}
}

func (serve *FileService) UploadFile(file *models.FileUpload) error {
	err := serve.FileServ.CreateFile(file)
	if err != nil {
		return err
	}
	return nil
}

func (serve *FileService) GetFile(id int) (*models.FileUpload, error) {
	file, err := serve.FileServ.GetFileById(id)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (serve *FileService) DeleteFile(id int) error {
	err := serve.FileServ.DeleteFile(id)
	if err != nil {
		return err
	}
	return nil
}

func (serve *FileService) UpdateFile(file *models.FileUpload) (*models.FileUpload, error) {
	file, err := serve.FileServ.UpdateFile(file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (serve *FileService) GetAllFiles() ([]*models.FileUpload, error) {
	files, err := serve.FileServ.GetAllFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}
