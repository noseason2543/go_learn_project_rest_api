package entities

import (
	"go_learn_project_rest_api/pkgs/logger"

	"github.com/gofiber/fiber/v3"
)

type IResponse interface {
	SuccessResponse(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	Context    fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (res *Response) SuccessResponse(code int, data any) IResponse {
	res.StatusCode = code
	res.Data = &data
	logger.InitLogger(res.Context, &res.Data).Print().Save()
	return res
}

func (res *Response) Error(code int, traceId, msg string) IResponse {
	res.IsError = true
	res.StatusCode = code
	res.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Msg:     msg,
	}
	logger.InitLogger(res.Context, &res.ErrorRes).Print().Save()
	return res
}

func (res *Response) Res() error {
	return res.Context.Status(res.StatusCode).JSON(func() any {
		if res.IsError {
			return &res.ErrorRes
		}
		return &res.Data
	}())
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
