package newsstorage

import (
	"context"
	"sf-mu/pkg/models/news"
	"sf-mu/pkg/storage/postgres"
)

type NewsStore struct {
	postgres.Store
}

// creates newsstorage
func New(ctx context.Context, constr string) (*NewsStore, error) {
	db, err := postgres.NewDb(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := NewsStore{
		postgres.Store{Db: db},
	}
	return &s, nil
}

// creates news posts
func (p *NewsStore) PostsCreation(posts []news.Post) error {
	for _, post := range posts {
		err := p.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

// Adds Post
func (s *NewsStore) AddPost(p news.Post) error {

	err := s.Db.QueryRow(context.Background(), `
		INSERT INTO gonews (title, content, pubtime, link)
		VALUES ($1, $2, $3, $4);
		`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	).Scan()
	return err
}

// searches by title
func (p *NewsStore) PostSearchILIKE(pattern string, limit, offset int) ([]news.Post, news.Pagination, error) {
	pattern = "%" + pattern + "%"

	pagination := news.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	row := p.Db.QueryRow(context.Background(), "SELECT count(*) FROM gonews WHERE title ILIKE $1;", pattern)
	err := row.Scan(&pagination.NumOfPages)

	if pagination.NumOfPages%limit > 0 {
		pagination.NumOfPages = pagination.NumOfPages/limit + 1
	} else {
		pagination.NumOfPages /= limit
	}

	if err != nil {
		return nil, news.Pagination{}, err
	}

	rows, err := p.Db.Query(context.Background(), "SELECT id, title, content, pubtime, link FROM gonews WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;", pattern, limit, offset)
	if err != nil {
		return nil, news.Pagination{}, err
	}
	defer rows.Close()
	var posts []news.Post
	for rows.Next() {
		var p news.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, news.Pagination{}, err
		}
		posts = append(posts, p)
	}
	return posts, pagination, rows.Err()
}

// Posts with pagination
func (s *NewsStore) Posts(limit, offset int) ([]news.Post, error) {
	pagination := news.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	rows, err := s.Db.Query(context.Background(), `
	SELECT id, title, content, pubtime, link FROM gonews
	ORDER BY pubtime DESC LIMIT $1 OFFSET $2
	`,
		pagination.Limit, pagination.Page,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []news.Post

	for rows.Next() {
		var p news.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)

	}
	// WARN: don't forget check error
	return posts, rows.Err()
}

// Post Details by id
func (p *NewsStore) PostDetails(id int) (news.Post, error) {
	row := p.Db.QueryRow(context.Background(), `
	SELECT * FROM gonews 
    WHERE id =$1;
	`, id)
	var post news.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PubTime,
		&post.Link)
	if err != nil {
		return news.Post{}, err
	}
	return post, nil
}

// Creates Gonews Table 
func (p *NewsStore) CreateGonewsTable() error {
	_, err := p.Db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS gonews (
                id SERIAL PRIMARY KEY,
                title TEXT NOT NULL DEFAULT 'empty',
                 content TEXT NOT NULL DEFAULT 'empty',
                pubtime BIGINT NOT NULL DEFAULT extract (epoch from now()),
                link TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// Drops Gonews Table 
func (p *NewsStore) DropGonewsTable() error {
	_, err := p.Db.Exec(context.Background(), "DROP TABLE IF EXISTS gonews;")
	if err != nil {
		return err
	}
	return nil
}
