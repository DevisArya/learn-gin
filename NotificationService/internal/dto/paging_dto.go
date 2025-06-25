package dto

type PaginationResponse struct {
	CurrentPage uint32
	Limit       uint32
	TotalRecord uint32
	TotalPage   uint32
}
