package appInfoUsecases

import "go_learn_project_rest_api/modules/appInfo/appInfoRepositories"

type IAppInfoUsecases interface {
}

type appInfoUsecases struct {
	appInfoRepositories appInfoRepositories.IAppInfoRepository
}

func AppInfoUsecases(repo appInfoRepositories.IAppInfoRepository) IAppInfoUsecases {
	return &appInfoUsecases{
		appInfoRepositories: repo,
	}
}
