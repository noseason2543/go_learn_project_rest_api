package appInfoUsecases

import (
	"go_learn_project_rest_api/modules/appInfo"
	"go_learn_project_rest_api/modules/appInfo/appInfoRepositories"
)

type IAppInfoUsecases interface {
	FindCategory(*appInfo.CategoryFilter) ([]*appInfo.Category, error)
}

type appInfoUsecases struct {
	appInfoRepositories appInfoRepositories.IAppInfoRepository
}

func AppInfoUsecases(repo appInfoRepositories.IAppInfoRepository) IAppInfoUsecases {
	return &appInfoUsecases{
		appInfoRepositories: repo,
	}
}

func (u *appInfoUsecases) FindCategory(req *appInfo.CategoryFilter) ([]*appInfo.Category, error) {
	category, err := u.appInfoRepositories.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}
