package censorsstorage

import (
	"context"
	"log"
	"sf-mu/pkg/models/censors"
	"sf-mu/pkg/storage/postgres"
)

type CsStore struct {
	postgres.Store
}

// New censor storage
func New(ctx context.Context, constr string) (*CsStore, error) {
	db, err := postgres.NewDb(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := CsStore{
		postgres.Store{Db: db},
	}
	return &s, nil
}

func (p *CsStore) AllList() ([]censors.BlackList, error) {
	rows, err := p.Db.Query(context.Background(), "SELECT id, ban_word FROM black_list")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []censors.BlackList
	for rows.Next() {
		var c censors.BlackList
		err = rows.Scan(&c.ID, &c.BanWord)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (p *CsStore) AddList(c censors.BlackList) error {
	_, err := p.Db.Exec(context.Background(),
		"INSERT INTO black_list (ban_word) VALUES ($1);", c.BanWord)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Create black list Table
func (p *CsStore) CreateBlackListTable() error {
	_, err := p.Db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS black_list (
			id SERIAL PRIMARY KEY,
			ban_word TEXT
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// Drop black list table
func (p *CsStore) DropBlackListTable() error {
	_, err := p.Db.Exec(context.Background(), "DROP TABLE IF EXISTS black_list;")
	if err != nil {
		return err
	}
	return nil
}
