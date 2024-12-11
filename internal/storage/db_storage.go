package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = storage.PingContext(ctx); err != nil {
		panic(err)
	}
	return &inDataBaseStorage{storage}
}

func (dbs *inDataBaseStorage) InsertOneUrl(ctx context.Context, url *Url, userID uint) error {
	query := `INSERT INTO urlstore (original_url, short_url, user_id) 
			VALUES ($1, $2, $3)
			ON CONFLICT (original_url) DO NOTHING`

	res, err := dbs.ExecContext(ctx, query, url.OriginalUrl, url.ShortUrl, userID)
	if err != nil {
		return err
	}
	if r, _ := res.RowsAffected(); r == 0 {
		query = `SELECT short_url FROM urlstore
				WHERE original_url=$1`
		var shortUrl string
		dbs.QueryRowContext(ctx, query, url.OriginalUrl).Scan(&shortUrl)
		return newUrlError(shortUrl)
	}
	return nil
}

func (dbs *inDataBaseStorage) InsertMultipleUrl(ctx context.Context, urls *[]*Url) error {
	tx, err := dbs.Begin()
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	stmtInsert, err := tx.PrepareContext(ctx, `INSERT INTO urlstore (original_url, short_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING`)
	if err != nil {
		return nil
	}
	defer stmtInsert.Close()

	for _, v := range *urls {
		v.ShortUrl = generateShortUrl()
		_, err := stmtInsert.ExecContext(ctx, v.OriginalUrl, v.ShortUrl)
		if err != nil {
			return nil
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

func (dbs *inDataBaseStorage) SelectAllPairUrl(ctx context.Context, userID uint) (*[]Url, error) {
	urls := make([]Url, 0)
	query := `SELECT original_url, short_url FROM urlstore WHERE user_id=$1`
	rows, err := dbs.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u Url
		err = rows.Scan(&u.OriginalUrl, &u.ShortUrl)
		if err != nil {
			continue
		}

		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &urls, nil
}

func (dbs *inDataBaseStorage) DeletePairUrl(ctx context.Context, shortUrl string) error {
	query := "DELETE FROM urlstore WHER short_url=$1"

	_, err := dbs.ExecContext(ctx, query, shortUrl)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *inDataBaseStorage) CreateUser(ctx context.Context) (uint, error) {
	query := `INSERT INTO users DEFAULT VALUES`
	_, err := dbs.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	query = `SELECT id FROM users
			ORDER BY id DESC
			LIMIT 1
			`
	var lastUserID uint
	row := dbs.QueryRowContext(ctx, query)
	row.Scan(&lastUserID)
	fmt.Println(lastUserID)

	return lastUserID, nil
}

func (dbs *inDataBaseStorage) initDB() {
	query := `CREATE TABLE IF NOT EXISTS users (
		"id" SERIAL PRIMARY KEY
	)`

	_, err := dbs.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	query = `CREATE TABLE IF NOT EXISTS urlstore (
				"ID" SERIAL PRIMARY KEY,
				"original_url" VARCHAR(30) UNIQUE,
				"short_url" VARCHAR(8),
				"user_id" SERIAL REFERENCES users ("id")
			)`

	_, err = dbs.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}
