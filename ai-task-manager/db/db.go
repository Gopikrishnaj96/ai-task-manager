package db

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
    var err error
    connStr := "postgres://pgadmin:5412@localhost:5432/taskmanager?sslmode=disable"
    DB, err = sqlx.Connect("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }
    fmt.Println("Database connected!")
}
