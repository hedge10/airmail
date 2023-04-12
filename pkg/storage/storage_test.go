package storage

import (
	"context"
	"os"
	"testing"

	"github.com/hedge10/airmail/pkg/conf"
	"github.com/hedge10/airmail/pkg/mail"
	"github.com/qiniu/qmgo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func setup() *conf.Config {
	os.Setenv("AM_MONGODB_USERNAME", "root")
	os.Setenv("AM_MONGODB_PASSWORD", "root")
	os.Setenv("AM_MONGODB_DB", "test")
	defer os.Unsetenv("AM_MONGODB_USERNAME")
	defer os.Unsetenv("AM_MONGODB_PASSWORD")
	defer os.Unsetenv("AM_MONGODB_DB")

	c, _ := conf.New()

	return c
}

func teardown(c *qmgo.QmgoClient, ctx context.Context) {
	err := c.Remove(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	c.Close(ctx)
}

func TestConnect(t *testing.T) {
	config := setup()

	s, conn_err := Connect(config)
	assert.Nil(t, conn_err)

	cnt, _ := s.Client.Find(s.Context, bson.M{}).Count()
	assert.Equal(t, int64(0), cnt)

	// Manually doing the teardown process as we have no results to delete
	s.Client.Close(s.Context)
}

func TestInsert(t *testing.T) {
	config := setup()
	s, err := Connect(config)
	assert.Nil(t, err)

	s.Message = &mail.Email{
		Subject: "Demo subject",
		Message: "Demo message",
		From: mail.Party{
			Name:  "Mickey",
			Email: "mickey@disney.com",
		},
		To: []mail.Party{
			{
				Name:  "Mickey",
				Email: "mickey@disney.com",
			},
		},
	}

	e := s.Insert()
	assert.Nil(t, e)

	cnt, _ := s.Client.Find(s.Context, bson.M{"subject": "Demo subject"}).Count()
	assert.Equal(t, int64(1), cnt)

	teardown(s.Client, s.Context)
}
