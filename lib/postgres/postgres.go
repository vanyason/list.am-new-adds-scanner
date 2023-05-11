package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type Postgres struct {
	sql.DB
	cfg Config
	db  *sql.DB
}

func New(cfg Config) Postgres {
	return Postgres{cfg: cfg}
}

func (cfg Config) string() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode)
}

func (p *Postgres) Connect() error {
	return nil
}

func (p *Postgres) Disconnect() error {
	return nil
}
