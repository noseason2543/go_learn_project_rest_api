package config

import (
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func convertEnvStringToInt(env map[string]string, field string) int {
	data, err := strconv.Atoi(env[field])
	if err != nil {
		log.Fatalf("load %v failed: %v", field, err)
	}
	return data
}

func convertEnvStringToTimeDuration(env map[string]string, field string) time.Duration {
	data, err := strconv.Atoi(env[field])
	if err != nil {
		log.Fatalf("load %v failed: %v", field, err)
	}
	return time.Duration(int64(data) * int64(math.Pow10(9)))
}

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}

	return &config{
		app: &app{
			host:         envMap["APP_HOST"],
			port:         convertEnvStringToInt(envMap, "APP_PORT"),
			name:         envMap["APP_NAME"],
			version:      envMap["APP_VERSION"],
			readTimeout:  convertEnvStringToTimeDuration(envMap, "APP_READ_TIMEOUT"),
			writeTimeout: convertEnvStringToTimeDuration(envMap, "APP_WRITE_TIMEOUT"),
			bodyLimit:    convertEnvStringToInt(envMap, "APP_BODY_LIMIT"),
			fileLimit:    convertEnvStringToInt(envMap, "APP_FILE_LIMIT"),
			gcpBucket:    envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host:          envMap["DB_HOST"],
			port:          convertEnvStringToInt(envMap, "DB_PORT"),
			protocol:      envMap["DB_PROTOCOL"],
			username:      envMap["DB_USERNAME"],
			password:      envMap["DB_PASSWORD"],
			database:      envMap["DB_DATABASE"],
			sslMode:       envMap["DB_SSL_MODE"],
			maxConnection: convertEnvStringToInt(envMap, "DB_MAX_CONNECTIONS"),
		},
		jwt: &jwt{
			adminKey:         envMap["APP_ADMIN_KEY"],
			secretKey:        envMap["JWT_SECRET_KEY"],
			apiKey:           envMap["APP_API_KEY"],
			accessExpiresAt:  convertEnvStringToInt(envMap, "JWT_ACCESS_EXPIRES"),
			refreshExpiresAt: convertEnvStringToInt(envMap, "JWT_REFRESH_EXPIRES"),
		},
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type IAppConfig interface{}

func (c *config) App() IAppConfig {
	return nil
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
	fileLimit    int
	gcpBucket    string
}

type IDbConfig interface{}

func (c *config) Db() IDbConfig {
	return nil
}

type db struct {
	host          string
	port          int
	protocol      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

type IJwtConfig interface{}

func (c *config) Jwt() IJwtConfig {
	return nil
}

type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int
	refreshExpiresAt int
}
