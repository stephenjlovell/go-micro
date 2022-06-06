package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stephenjlovell/go-micro/authentication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting authentication service...")

	conn := connetToDB()
	if conn == nil {
		log.Panic("could not connect to db.")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

const retryCount = 10

func connetToDB() *sql.DB {
	dataSource := os.Getenv("DATA_SOURCE")
	var err error
	var conn *sql.DB
	for i := 0; i < retryCount; i++ {
		conn, err = openDB(dataSource)
		if err != nil {
			log.Printf("unable to connect: %s\n", err.Error())
			time.Sleep(1 * time.Second)
			log.Printf("%ds - retrying connection...\n", i+1)
		} else {
			log.Println("connected to DB...")
			return conn
		}
	}
	return nil
}
