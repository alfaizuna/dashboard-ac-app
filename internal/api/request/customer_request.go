package request

type CustomerCreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Phone   string `json:"phone" validate:"required,min=10,max=15"`
	Address string `json:"address" validate:"required,min=10,max=500"`
	Email   string `json:"email" validate:"required,email"`
}

type CustomerUpdateRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Address *string `json:"address,omitempty" validate:"omitempty,min=10,max=500"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
}

type CustomerSearchRequest struct {
	*PaginationRequest
	Name  string `json:"name" query:"name"`
	Phone string `json:"phone" query:"phone"`
	Email string `json:"email" query:"email"`
}