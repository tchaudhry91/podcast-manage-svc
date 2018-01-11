package podcastmg

import (
	"flag"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"path"
	"testing"
)

var dbDialect = flag.String("dbDialect", "sqlite3", "Databse dialect to use, default=sqlite3")
var dbConnectionString = flag.String(
	"dbConnString",
	path.Join(os.TempDir(), "gorm.db"),
	"Connection string for the db to use, default tmp/sqlite")
var store DBStore

func init() {
	store = DBStore{
		*dbDialect,
		*dbConnectionString,
		nil,
	}
	if *dbDialect == "sqlite3" {
		os.Remove(*dbConnectionString)
	}
}

func TestDBConnection(t *testing.T) {
	err := store.Connect()
	defer store.Close()
	if err != nil {
		t.Errorf("Could not connect to DB:%s", err.Error())
	}

	// Empty Store
	badStore := DBStore{}
	err = badStore.Connect()
	if err == nil {
		t.Errorf("Should have errored on invalid store, but did not")
	}
}

func TestDBClosure(t *testing.T) {
	err := store.Connect()
	if err != nil {
		t.Errorf("Could not connect to DB:%s", err.Error())
	}
	err = store.Close()
	if err != nil {
		t.Errorf("Could not close the DB succesfully:%s", err.Error())
	}
	err = store.Database.CreateTable("test").Error
	if err == nil {
		t.Errorf("Should return error trying to operate on closed DB")
	}

	// Empty Store
	badStore := DBStore{}
	err = badStore.Close()
	if err == nil {
		t.Errorf("Should return error trying to close invalid store, but did not")
	}
}

func TestDBMigration(t *testing.T) {
	tables := []string{
		"podcasts",
		"podcast_items",
		"users",
		"subscriptions",
	}
	store.Connect()
	err := store.Migrate()
	if err != nil {
		t.Errorf("Failed to migrate DB:%v", err)
	}
	for _, table := range tables {
		if ok := store.Database.HasTable(table); !ok {
			t.Errorf("Migration failed. Missing table:%s", table)
		}
	}
	store.DropExistingTables()
	for _, table := range tables {
		if ok := store.Database.HasTable(table); ok {
			t.Errorf("Drop Tables failed. Found table:%s", table)
		}
	}

	store.Close()
	err = store.Migrate()
	if err == nil {
		t.Errorf("Migrating closed store should return error, but did not")
	}
}
