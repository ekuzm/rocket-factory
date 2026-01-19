package dto

import (
	"github.com/ekuzm/rocket-factory/inventory/internal/model"
)

type GetPartInput struct {
	Uuid string
}

type GetPartOutput struct {
	Part model.Part
}
