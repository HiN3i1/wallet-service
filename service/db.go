package service

import (
	"github.com/HiN3i1/wallet-service/db"
)

// CreateDBClient initialize database connection.
func CreateDBClient() {
	db.CreateDBClient()
	db.InitTable()
}
