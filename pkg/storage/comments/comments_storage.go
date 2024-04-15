package commentsstorage

import (
	"context"

	"sf-mu/pkg/models/comments"
	"sf-mu/pkg/storage/postgres"
)

type ComStore struct {
	postgres.Store
}

// creates new  commentsstorage
func New(ctx context.Context, constr string) (*ComStore, error) {
	db, err := postgres.NewDb(ctx, constr)
	if err != nil {
		return nil, err
	}

	s := &ComStore{
		postgres.Store{Db: db},
	}
	return s, nil
}

// returns all comments by newsID
func (p *ComStore) AllComments(newsID int) ([]comments.Comment, error) {
	rows, err := p.Db.Query(context.Background(), "SELECT id, news_id, content, pubtime FROM comments WHERE news_id = $1;", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comms []comments.Comment
	for rows.Next() {
		var c comments.Comment
		err = rows.Scan(&c.ID, &c.NewsID, &c.Content, &c.PubTime)
		if err != nil {
			return nil, err
		}
		comms = append(comms, c)
	}
	return comms, rows.Err()
}

// Adds Comment
func (p *ComStore) AddComment(c comments.Comment) error {
	_, err := p.Db.Exec(context.Background(),
		"INSERT INTO comments (news_id,content) VALUES ($1,$2);", c.NewsID, c.Content)
	if err != nil {
		return err
	}
	return nil
}

// Deletes Comment 
func (p *ComStore) DeleteComment(c comments.Comment) error {
	_, err := p.Db.Exec(context.Background(),
		"DELETE FROM comments WHERE id=$1;", c.ID)
	if err != nil {
		return err
	}
	return nil
}

// Creates Comment Table 
func (p *ComStore) CreateCommentTable() error {
	_, err := p.Db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS comments (
                id SERIAL PRIMARY KEY,
                news_id INT,
                content TEXT NOT NULL DEFAULT 'empty',
                pubtime BIGINT NOT NULL DEFAULT extract (epoch from now())
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// Drops Comment Table 
func (p *ComStore) DropCommentTable() error {
	_, err := p.Db.Exec(context.Background(), "DROP TABLE IF EXISTS comments;")
	if err != nil {
		return err
	}
	return nil
}
