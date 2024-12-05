package fileUsecases

import "go_learn_project_rest_api/config"

type IFileUsecases interface{}

type fileUsecases struct {
	cfg config.IConfig
}

func FileUsecases(cfg config.IConfig) IFileUsecases {
	return &fileUsecases{
		cfg: cfg,
	}
}
