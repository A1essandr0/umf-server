package models

type Config struct {
	DEVELOPMENT_MODE string `mapstructure:"development_mode" json:"development_mode"`
	WEB_PORT string `mapstructure:"web_port" json:"web_port"`
	PRODUCTION_CORS string `mapstructure:"production_cors" json:"production_cors"`

	KVSTORE_TYPE string `mapstructure:"kvstore_type" json:"kvstore_type"`
	// time to live in hours, 0 means no ttl
	DEFAULT_TTL int `mapstructure:"default_ttl" json:"default_ttl"`

	REDIS_ADDR string `mapstructure:"redis_addr" json:"redis_addr"`
	REDIS_PWD string `mapstructure:"redis_pwd" json:"redis_pwd"`
	REDIS_DB_NUM int `mapstructure:"redis_db_num" json:"redis_db_num"`

	DBSTORE_TYPE string `mapstructure:"dbstore_type" json:"dbstore_type"`
	DB_DSN string `mapstructure:"db_dsn" json:"db_dsn"`
	APPLY_MIGRATIONS bool `mapstructure:"apply_migrations" json:"apply_migrations"`

	USE_TLS bool `mapstructure:"use_tls" json:"use_tls"`
	CERT_FILE string `mapstructure:"cert_file" json:"cert_file"`
	CERT_KEY_FILE string `mapstructure:"cert_key_file" json:"cert_key_file"`

	DEFAULT_RECORDS_AMOUNT_TO_GET int `mapstructure:"default_records_amount_to_get" json:"default_records_amount_to_get"`
	HASH_LENGTH int `mapstructure:"hash_length" json:"hash_length"`
}