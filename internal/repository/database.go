package repository

import (
	"context"
	"errors"
	"time"

	"github.com/hollgett/shortURL.git/internal/logger"
	"github.com/hollgett/shortURL.git/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	tableScheme = `CREATE TABLE IF NOT EXISTS shortener (
"id_short" VARCHAR(8) NOT NULL UNIQUE,
"original" VARCHAR(2000) NOT NULL
);`
	saveReq      = `INSERT INTO shortener (id_short,original) VALUES ($1,$2)`
	findReq      = `SELECT original FROM shortener WHERE id_short = $1`
	saveBatchReq = `INSERT INTO shortener (id_short,original) VALUES (:short,:origin)`
)

type postgresDB struct {
	db       *sqlx.DB
	stmtSave *sqlx.Stmt
	stmtFind *sqlx.Stmt
}

// host=localhost user=postgres password=root dbname=USvideos sslmode=disable
func NewDataBase(ctx context.Context, dbDSN string) (Storage, error) {
	db, err := sqlx.Open("pgx", dbDSN)
	if err != nil {
		logger.LogInfo("open db error", zap.Error(err))
		return nil, err
	}
	db.SetMaxOpenConns(3)
	if _, err := db.ExecContext(ctx, tableScheme); err != nil {
		logger.LogInfo("create table", zap.Error(err))
		return nil, err
	}
	stmtSave, err := db.PreparexContext(ctx, saveReq)
	if err != nil {
		logger.LogInfo("create stmt save error", zap.Error(err))
		return nil, err
	}
	stmtFind, err := db.PreparexContext(ctx, findReq)
	if err != nil {
		logger.LogInfo("create stmt find error", zap.Error(err))
		return nil, err
	}
	return &postgresDB{
		db:       db,
		stmtSave: stmtSave,
		stmtFind: stmtFind,
	}, nil
}

func (pg *postgresDB) Save(ctx context.Context, shortLink, originURL string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := pg.stmtSave.ExecContext(ctx, shortLink, originURL); err != nil {
		logger.LogInfo("error insert db", zap.Error(err))
		return err
	}
	logger.LogInfo("save data")
	return nil
}
func (pg *postgresDB) Find(ctx context.Context, shortLink string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var origURL string
	row := pg.stmtFind.QueryRowxContext(ctx, shortLink)
	if err := row.Scan(&origURL); err != nil {
		return "", err
	}
	return origURL, nil
}
func (pg *postgresDB) Close() error {
	errSave := pg.stmtSave.Close()
	errFind := pg.stmtFind.Close()
	err := pg.db.Close()
	logger.LogInfo("close db")
	return errors.Join(err, errSave, errFind)
}

func (pg *postgresDB) Ping(rCtx context.Context) error {
	ctx, cancel := context.WithTimeout(rCtx, 3*time.Second)
	defer cancel()
	return pg.db.PingContext(ctx)
}

func (pg *postgresDB) SaveBatch(data []models.DBBatch) error {
	tx, err := pg.db.BeginTxx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamedContext(context.TODO(), saveBatchReq)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, v := range data {
		_, err := stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
