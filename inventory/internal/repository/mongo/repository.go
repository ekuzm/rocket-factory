package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo/entity"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	"github.com/ekuzm/rocket-factory/inventory/pkg/fixtures"
	errs "github.com/ekuzm/rocket-factory/platform/pkg/error"
	"github.com/ekuzm/rocket-factory/platform/pkg/logger"
)

var _ service.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) *repository {
	collection := db.Collection(entity.PartsCollection)

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: entity.PartsCollectionFieldUUID, Value: 1},
			{Key: entity.PartsCollectionFieldName, Value: 1},
			{Key: entity.PartsCollectionFieldCategory, Value: 1},
			{Key: entity.PartsCollectionFieldManufacturerCountry, Value: 1},
			{Key: entity.PartsCollectionFieldTags, Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to create index with name ", indexName)
		panic("failed to create index" + err.Error())
	}

	logger.Debug("Create index with name ", indexName)

	repository := &repository{
		collection: collection,
	}

	if !config.App().Mongo.IsInit() {
		repository.Init()
	}

	return repository
}

func (r *repository) Init() {
	parts := fixtures.GenerateParts()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ids, err := r.Save(ctx, parts)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"parts": parts,
			"err":   err,
		}).Error("Failed to save parts")

		panic("failed to save parts in mongodb " + err.Error())
	}

	logger.Debug("Count of inserted parts: ", len(ids))
}

func (r *repository) Save(ctx context.Context, parts []model.Part) ([]primitive.ObjectID, error) {
	partDocuments := entity.PartsToDocument(parts)

	docs := make([]any, len(partDocuments))

	for i, v := range partDocuments {
		docs[i] = v
	}

	res, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"docs": docs,
			"err":  err,
		}).Error("Failed to insert parts into collection")

		return nil, fmt.Errorf("insert parts into collection: %w", err)
	}

	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for i, v := range res.InsertedIDs {
		ids[i] = v.(primitive.ObjectID)
	}

	return ids, nil
}

func (r *repository) GetByUUID(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	var part entity.PartDocument

	err := r.collection.FindOne(ctx, bson.M{entity.PartsCollectionFieldUUID: uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, errs.ErrNotFound
		}

		return model.Part{}, fmt.Errorf("decode mongodb document into part model: %w", err)
	}

	return entity.PartToModel(part), nil
}

func (r *repository) GetAllByFilter(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	mongoFilter := buildMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, fmt.Errorf("find parts by filter: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Failed to close cursor: %v", cerr)
		}
	}()

	var parts []entity.PartDocument
	if err = cursor.All(ctx, &parts); err != nil {
		return nil, fmt.Errorf("maps documents with part objects: %w", err)
	}

	return entity.PartsToModel(parts), nil
}

func buildMongoFilter(filter model.Filter) bson.M {
	mongoFilter := bson.A{}

	if len(filter.UUIDs) > 0 {
		mongoFilter = append(mongoFilter, bson.M{entity.PartsCollectionFieldUUID: bson.M{"$in": filter.UUIDs}})
	}
	if len(filter.Names) > 0 {
		mongoFilter = append(mongoFilter, bson.M{entity.PartsCollectionFieldName: bson.M{"$in": filter.Names}})
	}
	if len(filter.Categories) > 0 {
		mongoFilter = append(mongoFilter, bson.M{entity.PartsCollectionFieldCategory: bson.M{"$in": entity.CategoriesToDocument(filter.Categories)}})
	}
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter = append(mongoFilter, bson.M{entity.PartsCollectionFieldManufacturerCountry: bson.M{"$in": filter.ManufacturerCountries}})
	}
	if len(filter.Tags) > 0 {
		mongoFilter = append(mongoFilter, bson.M{entity.PartsCollectionFieldTags: bson.M{"$eq": filter.Tags}})
	}

	if len(mongoFilter) == 0 {
		return bson.M{}
	}

	return bson.M{"$and": mongoFilter}
}
