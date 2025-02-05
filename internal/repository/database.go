package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/hollgett/shortURL.git/internal/logger"
	"go.uber.org/zap"
)

type postgresDB struct {
	db *sql.DB
}

// host=localhost user=postgres password=root dbname=USvideos sslmode=disable
func NewDataBase(dbDSN string) (Storage, error) {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		logger.LogInfo("open db error", zap.Error(err))
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS shortener (
"short" VARCHAR(8) NOT NULL UNIQUE,
"original" VARCHAR(2000) NOT NULL
);`); err != nil {
		logger.LogInfo("create table", zap.Error(err))
		return nil, err
	}
	return &postgresDB{
		db: db,
	}, nil
}

func (pg *postgresDB) Save(shortLink, originURL string) error {
	if _, err := pg.db.Exec(`INSERT INTO shortener (short,original) VALUES ($1,$2)`, shortLink, originURL); err != nil {
		logger.LogInfo("error insert db", zap.Error(err))
		return err
	}
	logger.LogInfo("save data")
	return nil
}
func (pg *postgresDB) Find(shortLink string) (string, error) {
	var origURL string
	row := pg.db.QueryRow(`SELECT original FROM shortener WHERE short = $1`, shortLink)
	if err := row.Scan(&origURL); err != nil {
		return "", err
	}
	return origURL, nil
}
func (pg *postgresDB) Close() error {
	logger.LogInfo("close db")
	return pg.db.Close()
}

func (pg *postgresDB) Ping(rCtx context.Context) error {
	ctx, cancel := context.WithTimeout(rCtx, 3*time.Second)
	defer cancel()
	return pg.db.PingContext(ctx)
}
