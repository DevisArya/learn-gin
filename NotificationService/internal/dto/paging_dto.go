package dto

type PaginationResponse struct {
	CurrentPage int
	Limit       int
	TotalRecord int
	TotalPage   int
}
