package storage

import (
	"context"
	"fmt"

	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/mail"
	"github.com/qiniu/qmgo"
)

type Storage struct {
	Context context.Context
	Client  *qmgo.QmgoClient
	Message *mail.Email
}

func Connect(c *conf.Config) (*Storage, error) {
	ctx := context.Background()
	uri := fmt.Sprintf("mongodb://%s:%d", c.MongoDbHost, c.MongoDbPort)
	credentials := qmgo.Credential{
		Username: c.MongoDbUsername,
		Password: c.MongoDbPassword,
	}
	client, err := qmgo.Open(ctx, &qmgo.Config{
		Auth:     &credentials,
		Uri:      uri,
		Database: c.MongoDbDatabase,
		Coll:     c.MongoDbCollection,
	})

	return &Storage{Client: client, Context: ctx}, err
}

func (s *Storage) Insert() error {
	_, err := s.Client.InsertOne(s.Context, s.Message)

	return err
}
