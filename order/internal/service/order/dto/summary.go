package dto

import "github.com/google/uuid"

type Summary struct {
	OrderUUID  uuid.UUID
	TotalPrice float64
}
