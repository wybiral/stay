package db

import (
	"testing"
)

func TestPage(t *testing.T) {
	page := NewPage()
	page.Set(0, true)
	if page.Get(0) != true {
		t.Errorf("page.Set(0, true) failed")
	}
	page.Set(0, false)
	if page.Get(0) != false {
		t.Errorf("page.Set(0, false) failed")
	}
	page.Set(64, true)
	if page.Get(64) != true {
		t.Errorf("page.Set(64, true) failed")
	}
}

func TestColumn(t *testing.T) {
	col := NewColumn()
	col.Set(100000, true)
	if col.Get(100000) != true {
		t.Errorf("col.Set(100000, true) failed")
	}
}

func TestDatabase(t *testing.T) {
	db := NewDatabase()
	db.Add("davy", "likes:coffee")
	if db.Get("davy", "likes:coffee") != true {
		t.Errorf("db.Add failed")
	}
	db.Remove("davy", "likes:coffee")
	if db.Get("davy", "likes:coffee") != false {
		t.Errorf("db.Remove failed")
	}
	db.Add("davy", "likes:coffee")
	db.Add("davy", "likes:programming")
	db.Add("nina", "likes:coffee")
	query := db.Query("likes:coffee")
	if query.Count() != 2 {
		t.Errorf("query.Count failed")
	}
	query = db.Query("likes:coffee").And(db.Query("likes:programming"))
	if query.Count() != 1 {
		t.Errorf("query.And failed")
	}
	keys := db.Keys(db.Query("likes:coffee"))
	temp := <-keys
	if temp != "davy" {
		t.Errorf("db.Keys failed on first")
	}
	temp = <-keys
	if temp != "nina" {
		t.Errorf("db.Keys failed on second")
	}
}
