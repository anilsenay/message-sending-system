package model

type DbFilters[T any] struct {
	Limit  *int
	Offset *int
	Order  *string
	Model  *T
}
