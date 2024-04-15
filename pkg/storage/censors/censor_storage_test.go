package censorsstorage

import (
	"context"
	"sf-mu/pkg/models/censors"

	"testing"
	"time"
)

const CSTR="postgres://postgres:postgres@localhost:5432/comm"

func TestStore_AddList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, CSTR )
	str := censors.BlackList{
		BanWord: "ups",
	}
	dataBase.AddList(str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Row created.")
}

func TestStore_AllList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, CSTR)
	if err != nil {
		t.Fatal(err)
	}

	result, err := dataBase.AllList()
	if err != nil {
		t.Fatal(err)
	}

	//
	if len(result) == 0 {
		t.Errorf("table \"blacklist\" is empty.")
	} else {
		t.Logf("Table \"blacklist\" contains %d rows.", len(result))
	}

}
