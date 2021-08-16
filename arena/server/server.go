package server

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Server struct {
	Mongo mongo.Client
}

func NewServer(mongoURI string) *Server {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		zap.L().Fatal("Cannot connect to Mongoutils.DB ",
			zap.String("mongoURI", mongoURI),
			zap.Error(err),
		)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		zap.L().Fatal("Cannot ping Mongoutils.DB ",
			zap.String("mongoURI", mongoURI),
			zap.Error(err),
		)
	}

	a := &Server{
		Mongo: *mongoClient,
	}

	return a
}
