package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	_ "github.com/joho/godotenv/autoload"
)

type DBClient struct {
	*pg.DB
}

var DB *DBClient

type dbLogger struct {
	Verbose   bool
	EmptyLine bool
}

func (h dbLogger) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		fmt.Printf("%s executing a query:\n%s\n", evt.Err, q)
	} else if h.Verbose {
		if h.EmptyLine {
			fmt.Println()
		}
		fmt.Println(string(q))
	}

	return ctx, nil
}

func (h dbLogger) AfterQuery(ctx context.Context, evt *pg.QueryEvent) error {
	if evt.Err != nil {
		q, _ := evt.FormattedQuery()
		fmt.Printf("%s executing a query:\n%s\n", evt.Err, q)
	}
	return nil
}

// NewClient create a db client
func CreateDBClient() *DBClient {
	c := pg.Connect(&pg.Options{
		Addr:     os.Getenv("POSTGRES_ADDR"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	})
	c.AddQueryHook(dbLogger{
		Verbose:   false,
		EmptyLine: true,
	})
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
		(*DepositCallBack)(nil),
	} {
		if err := DB.Model(model).DropTable(&orm.DropTableOptions{
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
		(*DepositCallBack)(nil),
	} {
		if err := DB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			Temp:          false,
			FKConstraints: true,
		}); err != nil {
			return fmt.Errorf("createSchema error: %v", err)
		}
	}
	return nil
}
