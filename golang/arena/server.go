package arena

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/wihrt/idle_arena/fight"
	"github.com/wihrt/idle_arena/gladiator"
	"github.com/wihrt/idle_arena/manager"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	DB      = "arena"
	M       = "managers"
	G       = "gladiators"
	APIBase = "api/v1"
)

type ArenaServer struct {
	Mongo mongo.Client
}

func NewArenaServer(mongoURI string) *ArenaServer {
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

	a := &ArenaServer{
		Mongo: *mongoClient,
	}

	return a
}

func (a *ArenaServer) NewManager(w http.ResponseWriter, r *http.Request) {

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

	_, err = a.getManager(m.ManagerID)
	if err == mongo.ErrNoDocuments {
		_, err := a.createManager(&m)
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
func (a *ArenaServer) GetManager(w http.ResponseWriter, r *http.Request) {

	var (
		m           = &manager.Manager{}
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-1]
	)

	m, err := a.getManager(managerID)
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

func (a *ArenaServer) DeleteManager(w http.ResponseWriter, r *http.Request) {

	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-1]
	)

	deleted, err := a.deleteManager(managerID)
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

func (a *ArenaServer) NewGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-2]
	)

	m, err := a.getManager(managerID)
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

	g, err := gladiator.NewGladiator(1, managerID)
	if err != nil {
		zap.L().Error("Cannot create gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}
	_, err = a.createGladiator(g)
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
	err = a.updateManager(m)
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

func (a *ArenaServer) GetGladiators(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-2]
	)

	_, err := a.getManager(managerID)
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

	g, err := a.getGladiators(managerID)
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

func (a *ArenaServer) GetGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-3]
		gladiatorID = splittedURL[len(splittedURL)-1]
	)

	_, err := a.getManager(managerID)
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

	g, err := a.getGladiator(managerID, gladiatorID)
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

func (a *ArenaServer) FightGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-4]
		gladiatorID = splittedURL[len(splittedURL)-2]
	)

	_, err := a.getManager(managerID)
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

	g, err := a.getGladiator(managerID, gladiatorID)
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

	fightResult, err := fight.ResolveFight(g)
	if err != nil {
		zap.L().Error("Cannot fight gladiator",
			zap.String("managerID", managerID),
			zap.String("gladiatorID", gladiatorID),
			zap.Error(err),
		)
		w.WriteHeader(500)
		return
	}

	if fightResult.FightWon {
		err := a.updateGladiator(g)
		if err != nil {
			zap.L().Error("Cannot update gladiator",
				zap.String("managerID", managerID),
				zap.String("gladiatorID", gladiatorID),
				zap.Error(err),
			)
			w.WriteHeader(500)
			return
		}
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

func (a *ArenaServer) DeleteGladiator(w http.ResponseWriter, r *http.Request) {
	var (
		splittedURL = strings.Split(r.RequestURI, "/")
		managerID   = splittedURL[len(splittedURL)-3]
		gladiatorID = splittedURL[len(splittedURL)-1]
	)

	_, err := a.getManager(managerID)
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

	deleted, err := a.deleteGladiator(managerID, gladiatorID)
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

// Private functions
func (a *ArenaServer) getManager(managerID string) (*manager.Manager, error) {
	var (
		m           manager.Manager
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := a.Mongo.Database(DB).Collection(M).FindOne(ctx, bson.M{"manager_id": managerID}).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", DB),
				zap.String("collection", M),
				zap.Error(err),
			)
			return &m, mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot search managers",
			zap.String("database", "arenaServer"),
			zap.String("collection", "managers"),
			zap.Error(err),
		)
		return &m, err
	}
	return &m, nil
}

func (a *ArenaServer) createManager(m *manager.Manager) (*manager.Manager, error) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	_, err := a.Mongo.Database(DB).Collection(M).InsertOne(ctx, *m)
	if err != nil {
		zap.L().Error("Cannot create manager in MongoDB",
			zap.String("managerID", m.ManagerID),
			zap.Error(err),
		)
		return m, err
	}

	return m, nil
}

func (a *ArenaServer) updateManager(m *manager.Manager) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res := a.Mongo.Database(DB).Collection(M).FindOneAndUpdate(ctx, bson.M{"manager_id": m.ManagerID}, bson.M{"$set": *m})

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

func (a *ArenaServer) deleteManager(managerID string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := a.Mongo.Database(DB).Collection(M).DeleteOne(ctx, bson.M{"manager_id": managerID})
	if err != nil {
		zap.L().Error("Cannot delete manager",
			zap.String("managerID", managerID),
			zap.Error(err),
		)
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (a *ArenaServer) getGladiators(managerID string) ([]gladiator.Gladiator, error) {

	var (
		g           []gladiator.Gladiator
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	cursor, err := a.Mongo.Database(DB).Collection(G).Find(ctx, bson.M{"manager_id": managerID})
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

func (a *ArenaServer) getGladiator(managerID string, gladiatorID string) (*gladiator.Gladiator, error) {

	var (
		g           gladiator.Gladiator
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)
	defer cancel()

	err := a.Mongo.Database(DB).Collection(G).FindOne(ctx, bson.M{"manager_id": managerID, "gladiator_id": gladiatorID}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			zap.L().Warn("No documents found",
				zap.String("database", DB),
				zap.String("collection", G),
				zap.Error(err),
			)
			return &g, mongo.ErrNoDocuments
		}
		zap.L().Error("Cannot search gladiators",
			zap.String("database", DB),
			zap.String("collection", G),
			zap.Error(err),
		)
		return &g, err
	}
	return &g, nil
}

func (a *ArenaServer) createGladiator(g *gladiator.Gladiator) (*gladiator.Gladiator, error) {

	var (
		ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	)

	defer cancel()

	_, err := a.Mongo.Database(DB).Collection(G).InsertOne(ctx, *g)
	if err != nil {
		zap.L().Error("Cannot create gladiatro in MongoDB",
			zap.String("gladiatorID", g.GladiatorID),
			zap.Error(err),
		)
		return g, err
	}

	return g, nil
}

func (a *ArenaServer) updateGladiator(g *gladiator.Gladiator) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res := a.Mongo.Database(DB).Collection(G).FindOneAndUpdate(ctx, bson.M{"manager_id": g.ManagerID, "gladiator_id": g.GladiatorID}, bson.M{"$set": *g})

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

func (a *ArenaServer) deleteGladiator(managerID string, gladiatorID string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, err := a.Mongo.Database(DB).Collection(G).DeleteOne(ctx, bson.M{"manager_id": managerID, "gladiator_id": gladiatorID})
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
