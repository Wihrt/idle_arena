package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/wihrt/idle_arena/fight"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (s *Server) NewGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-2]
	)

	m, err := s.getManager(managerID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No manager found",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot find manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	g, err := gladiator.NewGladiator(1, m, &s.Mongo)
	if err != nil {
		zap.L().Error("Cannot create gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}
	_, err = s.createGladiator(g)
	if err != nil {
		zap.L().Error("Cannot create gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	m.Gladiators = append(m.Gladiators, g.GladiatorID)
	err = s.updateManager(m)
	if err != nil {
		zap.L().Error("Cannot update manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("Cannot marshal gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write data",
			zap.Error(err),
		)
	}
}

func (s *Server) GetGladiators(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-2]
	)

	_, err := s.getManager(managerID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No manager found",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot find manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	g, err := s.getGladiators(managerID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("Cannot marshal gladiators",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write data",
			zap.Error(err),
		)
	}
}

func (s *Server) GetGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-3]
		gladiatorID = splittedURL[len(splittedURL)-1]
	)

	_, err := s.getManager(managerID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No manager found",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot find manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	g, err := s.getGladiator(managerID, gladiatorID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No gladiator found",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot search gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(g)
	if err != nil {
		zap.L().Error("Cannot marshal gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write data",
			zap.Error(err),
		)
	}
}

func (s *Server) FightGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-4]
		gladiatorID = splittedURL[len(splittedURL)-2]
		settings    fight.Settings
	)

	_, err := s.getManager(managerID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No manager found",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot find manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	g, err := s.getGladiator(managerID, gladiatorID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No gladiator found",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot search gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("Cannot read body",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(body, &settings)
	if err != nil {
		zap.L().Error("Cannot unmarshal body",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	fightResult, err := fight.ResolveFight(g, &s.Mongo, &settings)
	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	err = s.updateGladiator(g)
	if err != nil {
		zap.L().Error("Cannot update gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	data, err := json.Marshal(fightResult)
	if err != nil {
		zap.L().Error("Cannot marshal fight result",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		zap.L().Error("Cannot write data",
			zap.Error(err),
		)
	}
}

func (s *Server) DeleteGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-3]
		gladiatorID = splittedURL[len(splittedURL)-1]
	)

	_, err := s.getManager(managerID)
	if err == mongo.ErrNoDocuments {
		zap.L().Error("No manager found",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(404)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot find manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	deleted, err := s.deleteGladiator(managerID, gladiatorID)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if deleted {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204)
	}
}

func (s *Server) getGladiators(managerID string) ([]gladiator.Gladiator, error) {

	var (
		g           []gladiator.Gladiator
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	cursor, err := s.Mongo.Database(utils.DB).Collection(utils.G).Find(ctx, bson.M{"manager_id": managerID})
	if err != nil {
		zap.L().Error("Cannot search gladiators",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		return g, err
	}

	err = cursor.All(ctx, &g)
	if err != nil {
		zap.L().Error("Cannot decode gladiators",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		return g, err
	}
	return g, nil
}

func (s *Server) getGladiator(managerID string, gladiatorID string) (*gladiator.Gladiator, error) {

	var (
		g           gladiator.Gladiator
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := s.Mongo.Database(utils.DB).Collection(utils.G).FindOne(ctx, bson.M{"manager_id": managerID, "gladiator_id": gladiatorID}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", utils.DB),
				zap.String("collection", utils.G),
				zap.Error(err),
			)
			return &g, mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot search gladiators",
			zap.String("database", utils.DB),
			zap.String("collection", utils.G),
			zap.Error(err),
		)
		return &g, err
	}
	return &g, nil
}

func (s *Server) createGladiator(g *gladiator.Gladiator) (*gladiator.Gladiator, error) {

	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)

	defer cancel()

	_, err := s.Mongo.Database(utils.DB).Collection(utils.G).InsertOne(ctx, *g)
	if err != nil {
		zap.L().Error("Cannot create gladiatro in Mongoutils.DB ",
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		return g, err
	}

	return g, nil
}

func (s *Server) updateGladiator(g *gladiator.Gladiator) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res := s.Mongo.Database(utils.DB).Collection(utils.G).FindOneAndUpdate(ctx, bson.M{"manager_id": g.ManagerID, "gladiator_id": g.GladiatorID}, bson.M{"$set": *g})

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("managerID", g.ManagerID),
				zap.String("gladiatorID", g.GladiatorID),
				zap.Error(res.Err()),
			)
			return mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot update gladiator",
			zap.String("managerID", g.ManagerID),
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(res.Err()),
		)
		return res.Err()
	}
	return nil
}

func (s *Server) deleteGladiator(managerID string, gladiatorID string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := s.Mongo.Database(utils.DB).Collection(utils.G).DeleteOne(ctx, bson.M{"manager_id": managerID, "gladiator_id": gladiatorID})
	if err != nil {
		zap.L().Error("Cannot delete gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		return false, err
	}

	return res.DeletedCount > 0, nil
}
