package dto

// ListDto struct defines http response of users
type ListDto[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"totalItem"`
}
