package request

type ServiceCreateRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=100"`
	Price    float64 `json:"price" validate:"required,min=0"`
	Duration int     `json:"duration" validate:"required,min=1"` // in minutes
}

type ServiceUpdateRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Price    *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Duration *int     `json:"duration,omitempty" validate:"omitempty,min=1"`
}

type ServiceSearchRequest struct {
	*PaginationRequest
	Name     string  `json:"name" query:"name"`
	MinPrice float64 `json:"min_price" query:"min_price"`
	MaxPrice float64 `json:"max_price" query:"max_price"`
}