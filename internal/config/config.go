package config

import (
	"strconv"

	"github.com/A1essandr0/umf-server/internal/utils"
)

var (
	DEVELOPMENT_MODE = utils.GetEnv("DEVELOPMENT_MODE", "development")
	WEB_PORT = utils.GetEnv("WEB_PORT", "0.0.0.0:10007")
	PRODUCTION_CORS = "*" // TODO

	// time to live in hours, 0 means no ttl
	DEFAULT_TTL, _ = strconv.Atoi(utils.GetEnv("DEFAULT_TTL", "0"))

	REDIS_ADDR = utils.GetEnv("REDIS_ADDR", "localhost:6379")
	REDIS_PWD = utils.GetEnv("REDIS_PWD", "pwd")
	REDIS_DB_NUM, _ = strconv.Atoi(utils.GetEnv("REDIS_DB_NUM", "0"))

	// TODO proper dsn forming
	DB_DSN = utils.GetEnv("DB_DSN", "host=localhost user=umf_user dbname=umf port=5432 sslmode=disable")

	USE_TLS = utils.GetEnv("USE_TLS", "true")
	CERT_FILE = utils.GetEnv("CERT_FILE", "")
	CERT_KEY_FILE = utils.GetEnv("CERT_KEY_FILE", "")

	DEFAULT_RECORDS_AMOUNT_TO_GET = 5
	HASH_LENGTH, _ = strconv.Atoi(utils.GetEnv("HASH_LENGTH", "8"))
)


