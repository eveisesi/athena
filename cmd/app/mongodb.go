package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/eveisesi/athena/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func makeMongoDB(cfg config) (*mongo.Database, error) {

	q := url.Values{}
	q.Set("authMechanism", cfg.Mongo.AuthMech)
	q.Set("maxIdleTimeMS", strconv.FormatInt(int64(time.Second*10), 10))
	q.Set("connectTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("serverSelectionTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	q.Set("socketTimeoutMS", strconv.FormatInt(int64(time.Second*4), 10))
	c := &url.URL{
		Scheme:   "mongodb",
		Host:     fmt.Sprintf("%s:%d", cfg.Mongo.Host, cfg.Mongo.Port),
		User:     url.UserPassword(cfg.Mongo.User, cfg.Mongo.Pass),
		Path:     fmt.Sprintf("/%s", cfg.Mongo.Name),
		RawQuery: q.Encode(),
	}

	mc, err := mongodb.Connect(context.TODO(), c)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo, sleep and continue")
	}

	mdb := mc.Database(cfg.Mongo.Name)

	return mdb, nil

}
