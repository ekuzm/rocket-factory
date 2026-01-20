package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Category int

const (
	CategoryUnspecified = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Filter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

func (f *Filter) Validate() error {
	if f == nil {
		return fmt.Errorf("filter is required: %w", ErrInvalidFormat)
	}

	for _, partUUID := range f.UUIDs {
		if _, err := uuid.Parse(partUUID); err != nil {
			return fmt.Errorf("part uuid: %w", ErrInvalidFormat)
		}
	}

	return nil
}
