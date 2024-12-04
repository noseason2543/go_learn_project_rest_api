package appInfoRepositories

import "github.com/jmoiron/sqlx"

type IAppInfoRepository interface{}

type appInfoRepository struct {
	db *sqlx.DB
}

func AppInfoRepository(db *sqlx.DB) IAppInfoRepository {
	return &appInfoRepository{
		db: db,
	}
}
