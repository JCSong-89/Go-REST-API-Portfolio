package postrgres

import (
	"Go-REST-API-Portfolio/config"
	"fmt"
	"github.com/jmoiron/sqlx"
)

/*
DB Opt
*/
const (
	MAX_OPEN_CONNECTIONS    = 60
	CONNECTION_MAX_LIFETIME = 120
	MAX_IDLE_CONNECTIONS    = 30
	CONNECTION_MAX_IDLETIME = 20
)

/*
return db instance
*/
func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dbEnvListStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable user=%s password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlPassword,
	) // 문자열을 반환

	db, err := sqlx.Connect(c.Postgres.PgDriver, dbEnvListStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	db.SetConnMaxLifetime(CONNECTION_MAX_LIFETIME)
	db.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)
	db.SetConnMaxIdleTime(CONNECTION_MAX_IDLETIME)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}