package dto

// PaginationRequest 分页请求的基础结构
type PaginationRequest struct {
	Page     int `json:"page"`      // 当前页码，默认为1
	PageSize int `json:"page_size"` // 每页数量，默认为10
}

// GetOffset 获取数据库查询的偏移量
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取数据库查询的限制数量
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}
