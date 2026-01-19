package model

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Category      Category
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
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

type Category string

const (
	CategoryUnknown  = "UNKNOWN"
	CategoryEngine   = "ENGINE"
	CategoryFuel     = "FUEL"
	CategoryPorthole = "PORTHOLE"
	CategoryWing     = "WING"
)

type Filter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

func (f Filter) Validate() error {
	for _, partUUID := range f.UUIDs {
		if _, err := uuid.Parse(partUUID); err != nil {
			return err
		}
	}

	return nil
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}
