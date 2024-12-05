package fileHandlers

import (
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/files/fileUsecases"
)

type IFileHandlers interface{}

type fileHandlers struct {
	cfg          config.IConfig
	fileUsecases fileUsecases.IFileUsecases
}

func FileHandlers(cfg config.IConfig, fileUsecases fileUsecases.IFileUsecases) IFileHandlers {
	return &fileHandlers{
		cfg:          cfg,
		fileUsecases: fileUsecases,
	}
}
