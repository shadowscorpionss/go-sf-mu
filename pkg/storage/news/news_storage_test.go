package newsstorage

import (
	"context"
	"sf-mu/pkg/models/news"
	"testing"
	"time"
)
const CSTR="postgres://postgres:postgres@localhost:5432/aggregator"

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := New(ctx, CSTR)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStore_AddPost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, CSTR)
	post := news.Post{
		Title:   "тестирования",
		Content: "Пробный текст",
		PubTime: 5,
		Link:    "Линка",
	}
	dataBase.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Row inserted.")
}
