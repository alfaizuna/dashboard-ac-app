package request

// BaseRequest represents common request fields
type BaseRequest struct {
	ID string `json:"id" validate:"omitempty,uuid"`
}