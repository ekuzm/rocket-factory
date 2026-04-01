package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
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
)

var _ service.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func New(ctx context.Context, db *mongo.Database) *repository {
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

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		slog.Error("parts index create failed", "error", err)

		panic("failed to create index" + err.Error())
	}

	slog.Debug("parts index created")

	repository := &repository{
		collection: collection,
	}

	if !config.App().Mongo.IsInit() {
		repository.Init(ctx)
	}

	return repository
}

func (r *repository) Init(ctx context.Context) {
	parts := fixtures.GenerateParts()

	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	if _, err := r.Save(ctx, parts); err != nil {
		slog.Error("parts seed failed", "error", err)

		panic("failed to save parts in mongodb " + err.Error())
	}

	slog.Debug("parts seeded")
}

func (r *repository) Save(ctx context.Context, parts []model.Part) ([]primitive.ObjectID, error) {
	partDocuments := entity.PartsToDocument(parts)

	docs := make([]any, len(partDocuments))

	for i, v := range partDocuments {
		docs[i] = v
	}

	res, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		slog.Error("parts insert failed", "error", err)

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

		slog.Error("part decode failed", "partUUID", uuid, "error", err)

		return model.Part{}, fmt.Errorf("decode mongodb document into part model: %w", err)
	}

	return entity.PartToModel(part), nil
}

func (r *repository) GetAllByFilter(ctx context.Context, filter model.Filter) ([]model.Part, error) {
	mongoFilter := buildMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		slog.Error("parts find failed", "error", err)

		return nil, fmt.Errorf("find parts by filter: %w", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			slog.Warn("cursor close failed", "error", cerr)
		}
	}()

	var parts []entity.PartDocument
	if err = cursor.All(ctx, &parts); err != nil {
		slog.Error("parts decode failed", "error", err)

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
