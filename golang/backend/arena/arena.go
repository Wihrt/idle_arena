package arena

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/wihrt/idle_arena/arena/fight"
	"github.com/wihrt/idle_arena/arena/gladiator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var ErrNoManagerFound = errors.New("no manager found")

type Arena struct {
	Mongo mongo.Client
}

func NewArena(mongoURI string) *Arena {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		zap.L().Fatal("Cannot connect to MongoDB",
			zap.String("mongoURI", mongoURI),
			zap.Error(err),
		)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		zap.L().Fatal("Cannot ping MongoDB",
			zap.String("mongoURI", mongoURI),
			zap.Error(err),
		)
	}

	a := &Arena{
		Mongo: *mongoClient,
	}

	return a
}

func (a *Arena) getManager(mReq *Manager) (*Manager, error) {
	var (
		m           Manager
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := a.Mongo.Database("arena").Collection("managers").FindOne(ctx, bson.M{"user_id": m.UserID, "guild_id": m.GuildID}).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", "arena"),
				zap.String("collection", "managers"),
				zap.String("user_id", m.UserID),
				zap.String("guild_id", m.GuildID),
				zap.Error(err),
			)
			return &m, mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot search managers",
			zap.String("database", "arena"),
			zap.String("collection", "managers"),
			zap.String("user_id", m.UserID),
			zap.String("guild_id", m.GuildID),
			zap.Error(err),
		)
		return &m, mongo.ErrNoDocuments
	}
	return &m, nil
}

func (a *Arena) getGladiator(m *Manager) (*gladiator.Gladiator, error) {
	var (
		g           gladiator.Gladiator
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := a.Mongo.Database("arena").Collection("gladiators").FindOne(ctx, bson.M{"_id": m.Gladiators[0]}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", "arena"),
				zap.String("collection", "gladiators"),
				zap.String("id", m.Gladiators[0]),
				zap.Error(err),
			)
			return &g, mongo.ErrNoDocuments
		}
		zap.L().Warn("Cannot search gladiators",
			zap.String("database", "arena"),
			zap.String("collection", "gladiators"),
			zap.String("id", m.Gladiators[0]),
			zap.Error(err),
		)
		return &g, mongo.ErrNoDocuments
	}
	return &g, nil
}

func (a *Arena) GetGladiator(w http.ResponseWriter, r *http.Request) {
	mReq := decodeManager(r.Body)
	m, err := a.getManager(mReq)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
			return
		}

	}
	g, err := a.getGladiator(m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
			return
		}
	}

	w.WriteHeader(200)
	data, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("Cannot encode data")
	}
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write data")
	}
}

func (a *Arena) HireGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	mReq := decodeManager(r.Body)
	m, err := a.getManager(mReq)
	if err == mongo.ErrNoDocuments {
		_, err := a.Mongo.Database("arena").Collection("managers").InsertOne(ctx, *mReq)
		if err != nil {
			zap.L().Error("Cannot create manager in MongoDB",
				zap.String("guild_id", mReq.GuildID),
				zap.String("user_id", mReq.UserID),
				zap.Error(err),
			)
			w.WriteHeader(500)
			return
		}
		m = mReq
	}
	if err != nil {
		zap.L().Error("Cannot search in MongoDB",
			zap.String("user_id", mReq.UserID),
			zap.String("guild_id", mReq.GuildID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	if len(m.Gladiators) >= 1 {
		w.WriteHeader(401)
		return
	}

	g := gladiator.NewGladiator(1)
	resInsert, err := a.Mongo.Database("arena").Collection("gladiators").InsertOne(ctx, g)
	if err != nil {
		zap.L().Error("Cannot add gladiator to arena",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}
	m.Gladiators = append(m.Gladiators, resInsert.InsertedID.(primitive.ObjectID).String())
	resUpdate := a.Mongo.Database("arena").Collection("managers").FindOneAndReplace(ctx, bson.M{"guild_id": m.GuildID, "user_id": m.UserID}, m)
	if resUpdate.Err() != nil {
		zap.L().Error("Cannot update manager",
			zap.String("guild_id", m.GuildID),
			zap.String("user_id", m.UserID),
			zap.Error(resUpdate.Err()),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	data, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("Cannot decode gladiator")
	}
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write gladiator data")
	}
}

func (a *Arena) FightGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	mReq := decodeManager(r.Body)
	m, err := a.getManager(mReq)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
		}
	}
	g, err := a.getGladiator(m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
			return
		}
	}

	fight.ResolveFight(g)
	resUpdate := a.Mongo.Database("arena").Collection("managers").FindOneAndReplace(ctx, bson.M{"_id": m.Gladiators[0]}, m)
	if resUpdate.Err() != nil {
		zap.L().Error("Cannot update manager",
			zap.String("guild_id", m.GuildID),
			zap.String("user_id", m.UserID),
			zap.Error(resUpdate.Err()),
		)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

func (a *Arena) FireGladiator(w http.ResponseWriter, r *http.Request) {

	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	mReq := decodeManager(r.Body)
	m, err := a.getManager(mReq)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
		}
	}

	for _, id := range m.Gladiators {
		res, err := a.Mongo.Database("arena").Collection("gladiators").DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			zap.L().Error("Cannot delete gladiator",
				zap.String("id", id),
				zap.Error(err),
			)
		}
		if res.DeletedCount > 0 {
			zap.L().Info("Gladiator deleted",
				zap.String("ID", id))
		} else {
			zap.L().Warn("No gladiator deleted")
		}
	}

	resUpdate := a.Mongo.Database("arena").Collection("managers").FindOneAndReplace(ctx, bson.M{"guild_id": m.GuildID, "user_id": m.UserID}, m)
	if resUpdate.Err() != nil {
		zap.L().Error("Cannot update manager",
			zap.String("guild_id", m.GuildID),
			zap.String("user_id", m.UserID),
			zap.Error(resUpdate.Err()),
		)
	}

	w.WriteHeader(200)
}
