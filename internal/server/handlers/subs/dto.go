package subs

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"       example:"Netflix"`
	Price       int    `json:"price,omitempty"    example:"1000"`
	UserID      string `json:"user_id"            example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date"         example:"01-2024"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
}

type SubscriptionResponse struct {
	ID          string `json:"id"                 example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name"       example:"Netflix"`
	Price       int    `json:"price,omitempty"    example:"1000"`
	UserID      string `json:"user_id"            example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date"         example:"01-2024"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name"       example:"Netflix Premium"`
	Price       int    `json:"price,omitempty"    example:"1500"`
	UserID      string `json:"user_id"            example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date"         example:"01-2024"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
}

type SummaryRequest struct {
	ServiceName string `json:"service_name,omitempty" example:"Netflix"`
	UserID      string `json:"user_id,omitempty"      example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date"             example:"01-2024"`
	EndDate     string `json:"end_date"               example:"12-2024"`
}

type SummaryResponse struct {
	Summary int `json:"summary" example:"12000"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"error description"`
}
