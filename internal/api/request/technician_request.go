package request

type TechnicianCreateRequest struct {
	Name           string `json:"name" validate:"required,min=2,max=100"`
	Phone          string `json:"phone" validate:"required,min=10,max=15"`
	Specialization string `json:"specialization" validate:"required,min=2,max=100"`
}

type TechnicianUpdateRequest struct {
	Name           *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone          *string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	Specialization *string `json:"specialization,omitempty" validate:"omitempty,min=2,max=100"`
}

type TechnicianSearchRequest struct {
	*PaginationRequest
	Name           string `json:"name" query:"name"`
	Specialization string `json:"specialization" query:"specialization"`
}