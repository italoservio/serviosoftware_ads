package db

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type DB struct {
	Client *mongo.Client
}

func NewDB(uri string) (*DB, error) {
	db := &DB{}

	if _, err := db.Connection(uri); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Connection(uri string) (*mongo.Client, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}

	err = ping(client)
	if err != nil {
		return nil, err
	}

	log.Println("banco de dados conectado")

	d.Client = client

	return d.Client, nil
}

func connect(uri string) (*mongo.Client, error) {
	if uri == "" {
		return nil, errors.New("string de conexao com o banco de dados nao fornecida")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ping(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	return nil
}

func (d *DB) Disconnect() error {
	if err := d.Client.Disconnect(context.TODO()); err != nil {
		return err
	}

	log.Println("conexao com o banco de dados encerrada")
	return nil
}
