package dto

import (
	"math"
)

// Response 统一的API响应结构体
type Response struct {
	Success bool        `json:"success"`           // 请求是否成功
	Code    int         `json:"code"`              // 状态码
	Message string      `json:"message,omitempty"` // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Error   string      `json:"error,omitempty"`   // 错误信息（仅在Success为false时返回）
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"pageSize"`
	TotalPage int64       `json:"totalPage"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, err string) Response {
	return Response{
		Success: false,
		Code:    code,
		Message: message,
		Error:   err,
	}
}

// EmptySuccessResponse 创建不带数据的成功响应
func EmptySuccessResponse(message string) Response {
	return Response{
		Success: true,
		Code:    200,
		Message: message,
	}
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) PaginatedResponse {
	// 计算总页数（向上取整）
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	return PaginatedResponse{
		Data:      data,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}
}
