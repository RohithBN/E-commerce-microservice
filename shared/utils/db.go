package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {


	connStr := os.Getenv(("NEONDB_URI"))
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")
	dbPool = pool
	return dbPool, nil

}

func GetDB() (*pgxpool.Pool,error) {
	if dbPool == nil {
		var err error
		dbPool, err = ConnectDB()
		if err != nil {
			panic(err)
		}
	}
	return dbPool,nil
}
