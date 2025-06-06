package dto

// PaginationRequest 分页请求的基础结构
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"gte=1"`
	PageSize int `form:"page_size,default=10" binding:"gte=1,lte=100"`
}

func (p *PaginationRequest) Limit() int {
	if p.PageSize <= 0 {
		return 10 // 默认值
	}
	return p.PageSize
}

func (p *PaginationRequest) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit() // 使用Limit确保一致性
}
