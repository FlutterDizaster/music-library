package models

//go:generate easyjson -all library.go

// Library contains filtered songs and pagination data.
// @Description Response with music library and pagination.
type Library struct {
	Songs      []Song     `json:"songs"`      // Songs list
	Pagination Pagination `json:"pagination"` // Pagination
}
