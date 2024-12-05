package appInfoUsecases

import (
	"go_learn_project_rest_api/modules/appInfo"
	"go_learn_project_rest_api/modules/appInfo/appInfoRepositories"
)

type IAppInfoUsecases interface {
	FindCategory(*appInfo.CategoryFilter) ([]*appInfo.Category, error)
	InsertCategory([]*appInfo.Category) error
	DeleteCategory(int) error
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

func (u *appInfoUsecases) InsertCategory(req []*appInfo.Category) error {
	if err := u.appInfoRepositories.InsertCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appInfoUsecases) DeleteCategory(id int) error {
	if err := u.appInfoRepositories.DeleteCategory(id); err != nil {
		return err
	}
	return nil
}
