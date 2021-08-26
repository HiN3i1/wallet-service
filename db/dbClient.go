package db

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	_ "github.com/joho/godotenv/autoload"
)

type DBClient struct {
	*pg.DB
}

var DB *DBClient

type dbLogger struct{}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {
}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	// fmt.Println(q.FormattedQuery())
}

// NewClient create a db client
func CreateDBClient() *DBClient {
	c := pg.Connect(&pg.Options{
		Addr:     os.Getenv("POSTGRES_ADDR"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	})
	c.AddQueryHook(dbLogger{})
	DB = &DBClient{c}
	return DB
}

// Disconnect the client
func (m *DBClient) Disconnect() {
	m.Close()
}

// GetDBClient
func GetDBClient() *DBClient {
	return DB
}

// CleanTable drops tables.
func CleanTable() error {
	for _, model := range []interface{}{
		(*Wallet)(nil),
		(*Customer)(nil),
		(*SubWallet)(nil),
	} {
		if err := DB.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		}); err != nil {
			return fmt.Errorf("createSchema error: %v", err)
		}
	}
	return nil
}

// InitTable creates tables.
func InitTable() error {
	for _, model := range []interface{}{
		(*Wallet)(nil),
		(*Customer)(nil),
		(*SubWallet)(nil),
	} {
		if err := DB.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			Temp:          false,
			FKConstraints: true,
		}); err != nil {
			return fmt.Errorf("createSchema error: %v", err)
		}
	}
	return nil
}
