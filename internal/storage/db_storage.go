package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func newInDataBaseStorage(storageConfig storageConfig) *inDataBaseStorage {
	cred := strings.Split(storageConfig.Parameter, ":")
	if len(cred) != 2 {
		panic(errors.New("database parametr must be user:password"))
	}

	dataSource := fmt.Sprintf("host=localhost user=%s password=%s dbname=postgres sslmode=disable", cred[0], cred[1])

	storage, err := sql.Open("pgx", dataSource)
	if err != nil {
		panic(err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// if err = storage.PingContext(ctx); err != nil {
	// 	panic(err)
	// }
	return &inDataBaseStorage{storage}
}

func (dbs *inDataBaseStorage) InsertOneUrl(ctx context.Context, url *Url) error {
	query := `INSERT INTO urlstore (original_url, short_url) VALUES ($1, $2)`
	_, err := dbs.ExecContext(ctx, query , url.OriginalUrl, url.ShortUrl)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (dbs *inDataBaseStorage) InsertMultipleUrl(ctx context.Context, urls *[]*Url) error {
	tx, err := dbs.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO urlstore (original_url, short_url) VALUES ($1, $2)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, v := range *urls {
		_, err := stmt.ExecContext(ctx, v.OriginalUrl, v.ShortUrl)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (dbs *inDataBaseStorage) SelectOneUrl(ctx context.Context, shortUrl string) (string, error) {
	var original_url string

	query := `SELECT original_url FROM urlstore WHERE short_url=$1`

	row := dbs.QueryRowContext(ctx, query, shortUrl)
	row.Scan(&original_url)
	return original_url, nil
}

func (dbs *inDataBaseStorage) DeletePairUrl(ctx context.Context, shortUrl string) error {
	query := "DELETE FROM urlstore WHER short_url=$1"

	_, err := dbs.ExecContext(ctx, query, shortUrl)
	if err != nil {
		return err
	}

	return nil
}


func (dbs *inDataBaseStorage) initDB() {
	query := `CREATE TABLE IF NOT EXISTS urlstore (
				"ID" SERIAL PRIMARY KEY,
				"original_url" VARCHAR(30),
				"short_url" VARCHAR(8)
			)`
	
	_, err := dbs.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}