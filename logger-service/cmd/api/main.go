package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stephenjlovell/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const webPort = "80"
const rpcPort = "5001"
const grpcPort = "50001"
const mongoUrl = "mongodb:27017"

var client *mongo.Client

type Config struct {
	Models data.Models
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient
	// create a context in order to disconnect properly
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	go app.serve()
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl).SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error connecting:", err)
		return nil, err
	}
	return conn, nil
}
