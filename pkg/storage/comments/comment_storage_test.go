package commentsstorage

import (
	"context"
	"sf-mu/pkg/models/comments"
	"testing"
	"time"
)
const CSTR="postgres://postgres:postgres@localhost:5432/comm"
func TestStore_AddComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, CSTR)
	comment := comments.Comment{
		NewsID:  2,
		Content: "Текст проверки",
	}
	dataBase.AddComment(comment)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Row insert.")
}

func TestStore_DeleteComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, CSTR)
	comment := comments.Comment{
		ID: 1,
	}
	dataBase.DeleteComment(comment)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Row deleted.")
}
