package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbStore map[string]*mongo.Database // map [dbKey] => mongo.Database
var defaultDBName string

func init() {
	dbStore = make(map[string]*mongo.Database)
}

// errors

// ErrDBNotFound error occurs when the given database name does not found in database store
var ErrDBNotFound = errors.New("database not found")

func SetupDBWithConnectionURI(name string, connectionURI string) (*mongo.Database, error) {
	if len(name) == 0 {
		return nil, errors.New("name could not be blank")
	}

	db, found := dbStore[name]
	if found {
		return db, nil
	}

	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %v", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	// ping until server found, since Connect(...) does not block
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect (ping failed)")
	}

	// store db with key
	db = client.Database(name)
	dbStore[name] = db

	return db, nil
}
