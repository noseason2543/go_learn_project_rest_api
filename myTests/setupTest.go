package myTests

import (
	"encoding/json"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/modules/servers"
	"go_learn_project_rest_api/pkgs/databases"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := databases.DbConnection(cfg.Db())

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
