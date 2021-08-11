package dnd

import (
	"context"
	"sync"
	"time"

	"github.com/wihrt/idle_arena/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	DB = "arena"
	A  = "armor"
	W  = "weapon"
)

type Loader struct {
	API     *Client
	Mongo   *mongo.Client
	Timeout time.Duration
}

func NewLoader(baseUrl string, mongoURI string, timeout time.Duration) *Loader {
	l := &Loader{Timeout: timeout}

	var ctx, cancel = context.WithTimeout(context.Background(), timeout)
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

	l.API = NewClient(baseUrl, timeout)
	l.Mongo = mongoClient

	return l
}

func (l *Loader) LoadIndex(indexName string, wg *sync.WaitGroup) {
	var (
		pathUrl    = "/api/equipment-categories/" + indexName
		internalWg sync.WaitGroup
	)
	defer wg.Done()

	zap.L().Info("Loading index",
		zap.String("name", indexName),
	)

	equipements, err := l.API.GetCategory(pathUrl)
	if err != nil {
		zap.L().Error("Cannot get url",
			zap.String("path", pathUrl),
			zap.Error(err),
		)
	}
	for _, equipment := range equipements.Equipement {
		if utils.StringContains(equipment.URL, []string{"equipment"}) {
			internalWg.Add(1)
			go l.LoadEquipment(indexName, equipment, &internalWg)
		}
	}
	internalWg.Wait()
}

func (l *Loader) LoadEquipment(indexName string, e CategoryItem, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()

	switch indexName {
	case "armor":
		err = l.LoadArmor(e)
	case "weapon":
		err = l.LoadWeapon(e)
	}

	if err != nil {
		zap.L().Error("Cannot load "+indexName,
			zap.String("index", e.Index),
			zap.String("name", e.Name),
			zap.String("path", e.URL),
			zap.Error(err),
		)
	}
}

func (l *Loader) LoadArmor(e CategoryItem) error {
	var ctx, cancel = context.WithTimeout(context.Background(), l.Timeout)
	defer cancel()

	a, err := l.API.GetArmor(e.URL)
	if err != nil {
		zap.L().Error("Cannot get armor",
			zap.String("index", e.Index),
			zap.String("pathUrl", e.URL),
			zap.Error(err),
		)
		return err
	}

	_, err = l.Mongo.Database(DB).Collection(A).InsertOne(ctx, *a)
	if err != nil {
		zap.L().Error("Cannot create armor in MongoDB",
			zap.String("index", a.Index),
			zap.String("name", a.Name),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (l *Loader) LoadWeapon(e CategoryItem) error {
	var ctx, cancel = context.WithTimeout(context.Background(), l.Timeout)
	defer cancel()

	w, err := l.API.GetWeapon(e.URL)
	if err != nil {
		zap.L().Error("Cannot get weapon",
			zap.String("index", e.Index),
			zap.String("pathUrl", e.URL),
			zap.Error(err),
		)
		return err
	}

	l.Mongo.Database(DB).Collection(W).InsertOne(ctx, *w)
	if err != nil {
		zap.L().Error("Cannot create manager in MongoDB",
			zap.String("index", w.Index),
			zap.String("name", w.Name),
			zap.Error(err),
		)
		return err
	}

	return nil
}
