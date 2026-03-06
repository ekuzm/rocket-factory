package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartDocument struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          uuid.UUID          `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      Category           `bson:"category"`
	Dimensions    *Dimensions        `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer      `bson:"manufacturer,omitempty"`
	Tags          []string           `bson:"tags"`
	Metadata      map[string]*Value  `bson:"metadata"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty"`
}

type Category string

const (
	CategoryUnspecified Category = "unspecified"
	CategoryEngine      Category = "engine"
	CategoryFuel        Category = "fuel"
	CategoryPorthole    Category = "porthole"
	CategoryWing        Category = "wing"
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Value struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
}

const (
	PartsCollectionFieldID                  = "id"
	PartsCollectionFieldUUID                = "uuid"
	PartsCollectionFieldName                = "name"
	PartsCollectionFieldDescription         = "description"
	PartsCollectionFieldPrice               = "price"
	PartsCollectionFieldStockQuantity       = "stock_quantity"
	PartsCollectionFieldCategory            = "category"
	PartsCollectionFieldDimensions          = "dimensions"
	PartsCollectionFieldDimensionsLength    = "dimensions.length"
	PartsCollectionFieldDimensionsWidth     = "dimensions.width"
	PartsCollectionFieldDimensionsHeight    = "dimensions.height"
	PartsCollectionFieldDimensionsWeight    = "dimensions.weight"
	PartsCollectionFieldManufacturer        = "manufacturer"
	PartsCollectionFieldManufacturerName    = "manufacturer.name"
	PartsCollectionFieldManufacturerCountry = "manufacturer.country"
	PartsCollectionFieldManufacturerWebsite = "manufacturer.website"
	PartsCollectionFieldTags                = "tags"
	PartsCollectionFieldMetadata            = "metadata"
	PartsCollectionFieldCreatedAt           = "created_at"
	PartsCollectionFieldUpdatedAt           = "updated_at"
)
