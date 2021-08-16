package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/wihrt/idle_arena/manager"
	"github.com/wihrt/idle_arena/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (s *Server) NewManager(w http.ResponseWriter, r *http.Request) {

	var m manager.Manager

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("Cannot read body",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(body, &m)
	if err != nil {
		zap.L().Error("Cannot unmarshal body",
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	_, err = s.getManager(m.ManagerID)
	if err == mongo.ErrNoDocuments {
		_, err := s.createManager(&m)
		if err != nil {
			zap.L().Error("Cannot create manager",
				zap.String("managerID", m.ManagerID),
				zap.Error(err),
			)
			w.WriteHeader(500)
			return
		}
	}
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot get manager",
			zap.String("managerID", m.ManagerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
}

// Public functions
func (s *Server) GetManager(w http.ResponseWriter, r *http.Request) {

	var (
		m           = &manager.Manager{}
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-1]
	)

	m, err := s.getManager(managerID)
	if err != nil && err != mongo.ErrNoDocuments {
		zap.L().Error("Cannot search manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(404)
		return
	}

	data, err := json.Marshal(m)
	if err != nil {
		zap.L().Error("Cannot marshal manager",
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

func (s *Server) DeleteManager(w http.ResponseWriter, r *http.Request) {

	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-1]
	)

	deleted, err := s.deleteManager(managerID)
	if err != nil {
		zap.L().Error("Cannot delete manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	if deleted {
		w.WriteHeader(200)
	} else {
		zap.L().Warn("No manager deleted",
			zap.String("managerID", managerID),
		)
		w.WriteHeader(204)
	}
}

// Private functions
func (s *Server) getManager(managerID string) (*manager.Manager, error) {
	var (
		m           manager.Manager
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := s.Mongo.Database(utils.DB).Collection(utils.M).FindOne(ctx, bson.M{"manager_id": managerID}).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", utils.DB),
				zap.String("collection", utils.M),
				zap.Error(err),
			)
			return &m, mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot search managers",
			zap.String("database", utils.DB),
			zap.String("collection", utils.M),
			zap.Error(err),
		)
		return &m, err
	}
	return &m, nil
}

func (s *Server) createManager(m *manager.Manager) (*manager.Manager, error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	_, err := s.Mongo.Database(utils.DB).Collection(utils.M).InsertOne(ctx, *m)
	if err != nil {
		zap.L().Error("Cannot create manager in Mongoutils.DB ",
			zap.String("managerID", m.ManagerID),
			zap.Error(err),
		)
		return m, err
	}

	return m, nil
}

func (s *Server) updateManager(m *manager.Manager) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res := s.Mongo.Database(utils.DB).Collection(utils.M).FindOneAndUpdate(ctx, bson.M{"manager_id": m.ManagerID}, bson.M{"$set": *m})

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("managerID", m.ManagerID),
				zap.Error(res.Err()),
			)
			return mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot update manager",
			zap.String("managerID", m.ManagerID),
			zap.Error(res.Err()),
		)
		return res.Err()
	}

	return nil

}

func (s *Server) deleteManager(managerID string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := s.Mongo.Database(utils.DB).Collection(utils.M).DeleteOne(ctx, bson.M{"manager_id": managerID})
	if err != nil {
		zap.L().Error("Cannot delete manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		return false, err
	}

	return res.DeletedCount > 0, nil
}
