package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

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

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	GCPBucket() string
}

func (a *app) Url() string { return fmt.Sprintf("%s:%d", a.host, a.port) }

func (a *app) Name() string { return a.name }

func (a *app) Version() string { return a.version }

func (a *app) ReadTimeout() time.Duration { return a.readTimeout }

func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }

func (a *app) BodyLimit() int { return a.bodyLimit }

func (a *app) FileLimit() int { return a.fileLimit }

func (a *app) GCPBucket() string { return a.gcpBucket }

func (c *config) App() IAppConfig { return c.app }

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

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
}

func (c *config) Db() IDbConfig {
	return c.db
}

func (db *db) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.host,
		db.port,
		db.username,
		db.password,
		db.database,
		db.sslMode,
	)
}

func (db *db) MaxOpenConns() int { return db.maxConnection }

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

type IJwtConfig interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(t int)
	SetJwtRefreshExpires(t int)
}

func (jwt *jwt) SecretKey() []byte { return []byte(jwt.secretKey) }

func (jwt *jwt) AdminKey() []byte { return []byte(jwt.adminKey) }

func (jwt *jwt) ApiKey() []byte { return []byte(jwt.apiKey) }

func (jwt *jwt) AccessExpiresAt() int { return jwt.accessExpiresAt }

func (jwt *jwt) RefreshExpiresAt() int { return jwt.refreshExpiresAt }

func (jwt *jwt) SetJwtAccessExpires(t int) { jwt.accessExpiresAt = t }

func (jwt *jwt) SetJwtRefreshExpires(t int) { jwt.refreshExpiresAt = t }

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int
	refreshExpiresAt int
}
