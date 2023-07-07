package postgres

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/vanyason/list.am-new-adds-scanner/lib/ads"
)

// Function used to check string passed to the sql query. My attempt to prevent sql injection
var isLegal = regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString

// Wrapper over pgx lib that cleans string to prevent sql injection
func sanitize(str string) string {
	return pgx.Identifier{str}.Sanitize()
}

// Wrapper over pgx lib for easy work with the postgres
type Postgres struct {
	conn *pgx.Conn
}

// Create new wrapper to work with the postgres.
// Do not forget to call defer Close() to close connection
func New(userName, password, port, dbName string) (Postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", userName, password, port, dbName)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return Postgres{}, fmt.Errorf("unable to connect to the database: %w", err)
	}

	return Postgres{conn: conn}, nil
}

// Closes connection to the db.
// Do not forget to call it when creating connection!
// It is okay to call this func even if connection failed
func (p *Postgres) Close() {
	if p.conn != nil {
		p.conn.Close(context.Background())
	}
}

// Create table for Ads
// Format:
// id(ts) | link(str) | price(str) | description(str) | at(str) | time(ts)
func (p *Postgres) CreateAdsTable(name string) error {
	sql := `
		CREATE TABLE IF NOT EXISTS` + sanitize(name) + `(
				id 			SERIAL PRIMARY KEY,
				link 		TEXT not null,
				price 		TEXT not null,
				description TEXT not null,
				at 			TEXT unique not null,
				time 		TIMESTAMP
		);`

	_, err := p.conn.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("unable to create table for ads: %w", err)
	}

	return nil
}

// Check if table exists
func (p *Postgres) ExistAdsTable(name string) (bool, error) {
	sql := `SELECT EXISTS (
		SELECT
		FROM information_schema.tables
		WHERE table_name = $1
	);`

	var exists bool
	raw := p.conn.QueryRow(context.Background(), sql, name)
	err := raw.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("unable to check ads table: %w", err)
	}

	return exists, nil
}

// Drop table if exists
func (p *Postgres) DropTable(name string) error {
	sql := "DROP TABLE IF EXISTS " + sanitize(name) + ";"

	_, err := p.conn.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("unable to create table for ads: %w", err)
	}

	return nil
}

// Insert ads to the table create with CreateAdsTable
// Important:
// 0. Insert will be wrapped into txn
// 1. Empty slice doesn`t cause error - just nothing happens
// 1. Ads should be unique by ads.At field. That means that ads duplicates will be removed
// 2. If database already has ad with the same ad.At - new ad will not be inserted
func (p *Postgres) InsertUnique(dbName string, ads []ads.Ad) error {
	if len(ads) == 0 {
		return nil
	}

	return nil
}

// Remove everything from table
func (p *Postgres) AlterTable(dbName string, ads []ads.Ad) error {
	return nil
}
