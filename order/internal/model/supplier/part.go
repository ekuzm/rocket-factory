package supplier

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Category      Category
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     time.Time
	UpdateAt      *time.Time
}

type Dimensions struct {
	Length float64
	Width  float64
	Weight float64
	Height float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Filter struct {
	UUIDs                 uuid.UUIDs
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}
