package dto

type CreateOrderInput struct {
	UserUUID  string
	PartUUIDs []string
}

type CreateOrderOutput struct {
	OrderUUID  string
	TotalPrice float64
}
