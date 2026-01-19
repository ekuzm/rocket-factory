package dto

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

type ListPartsInput struct {
	Filter *model.Filter
}

type ListPartsOutput struct {
	Parts []model.Part
}
