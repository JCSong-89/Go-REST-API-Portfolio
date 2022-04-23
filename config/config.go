package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Cookie   Cookie
	Store    Store
	Session  Session
	Logger   Logger
	Jaeger   Jaeger
	Metrics  Metrics
	Minio    Minio
}

type ServerConfig struct {
	AppVersion       string
	Port             string
	Pproport         string
	Mode             string
	JwtSecretKey     string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	SSL              bool
	CtxDefultTimeout time.Duration
	CSRF             bool
	Debug            bool
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultDb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

type MongoDB struct {
	MongoURI string
}

type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

type Session struct {
	Prefix string
	Name   string
	Expire int
}

type Metrics struct {
	URL         string
	ServiceName string
}

type Store struct {
	ImagesFolder string
}

type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

type Minio struct {
	Endpoint       string
	MinioAccessKey string
	MinioSecretKey string
	UseSSL         bool
	MinioEndpoint  string
}